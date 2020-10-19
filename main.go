package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path"
	"sort"
	"strconv"

	"github.com/gdamore/tcell"
	"github.com/jroyal/gohar/har"
	"github.com/rivo/tview"
)

var currentEntry *har.Entry
var generalTextView *tview.TextView
var requestTextView *tview.TextView
var responseTextView *tview.TextView

func createNetworkTable(app *tview.Application, harFile *har.HarFile, fileName string) *tview.Table {
	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false).
		SetSeparator(' ').
		SetEvaluateAllRows(true)

	// Create Table headers
	name := tview.NewTableCell("Name")
	status := tview.NewTableCell("Status")
	rtype := tview.NewTableCell("Type")
	initiator := tview.NewTableCell("Initiator")
	size := tview.NewTableCell("Size")
	time := tview.NewTableCell("Time")
	table.SetCell(0, 0, name)
	table.SetCell(0, 1, status)
	table.SetCell(0, 2, rtype)
	table.SetCell(0, 3, initiator)
	table.SetCell(0, 4, size)
	table.SetCell(0, 5, time)

	for row, entry := range harFile.Log.Entries {
		row++ // bump the row up to skip past the header row

		url, _ := url.Parse(entry.Request.URL)
		urlBase := path.Base(url.Path)
		if urlBase == "/" {
			urlBase = entry.Request.URL
		}

		domain := tview.NewTableCell(urlBase).SetMaxWidth(75).SetExpansion(1)
		status := tview.NewTableCell(strconv.Itoa(entry.Response.Status))
		rType := tview.NewTableCell(entry.ResourceType)
		initiator := tview.NewTableCell(entry.Initiator.URL).SetMaxWidth(50)
		size := tview.NewTableCell(strconv.Itoa(entry.Response.TransferSize))
		time := tview.NewTableCell(fmt.Sprintf("%.0f ms", entry.Time))
		table.SetCell(row, 0, domain)
		table.SetCell(row, 1, status)
		table.SetCell(row, 2, rType)
		table.SetCell(row, 3, initiator)
		table.SetCell(row, 4, size)
		table.SetCell(row, 5, time)

	}
	table.SetTitle(fileName).SetBorder(true)
	table.Select(1, 0).SetFixed(1, 0).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
	}).SetSelectedFunc(func(row int, column int) {
		table.GetCell(row, column).SetTextColor(tcell.ColorRed)
	}).SetSelectionChangedFunc(func(row int, column int) {
		if row == 0 {
			return
		}
		currentEntry = &harFile.Log.Entries[row-1]
		updateDetailPanes()
	})
	return table
}

func createTextPane(title string) *tview.TextView {
	textView := tview.NewTextView().
		SetScrollable(true).
		SetDynamicColors(true)
	textView.SetBorder(true).SetTitle(title)
	return textView
}

func updateDetailPanes() {
	generalTextView.Clear()
	requestTextView.Clear()
	responseTextView.Clear()

	text := fmt.Sprintf("[yellow]Request URL:[white] %s\n", currentEntry.Request.URL)
	text += fmt.Sprintf("[yellow]Request Method:[white] %s\n", currentEntry.Request.Method)
	text += fmt.Sprintf("[yellow]Status Code:[white] %d %s\n", currentEntry.Response.Status, currentEntry.Response.StatusText)
	text += fmt.Sprintf("[yellow]Remote Address:[white] %s", currentEntry.ServerIPAddress)
	fmt.Fprint(generalTextView, text)

	fmt.Fprint(requestTextView, getHeaderText(currentEntry.Request.Headers))
	fmt.Fprint(responseTextView, getHeaderText(currentEntry.Response.Headers))
}

func getHeaderText(headers []har.Headers) string {
	sort.Slice(headers, func(i, j int) bool {
		return headers[i].Name < headers[j].Name
	})
	text := ""
	for _, header := range headers {
		text += fmt.Sprintf("[yellow]%s:[white] %s\n", header.Name, header.Value)
	}
	return text
}

func main() {

	f, err := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	fileName := "test.almightyzero.com.har"
	harFile := har.Load(fileName)
	app := tview.NewApplication()

	networkTable := createNetworkTable(app, &harFile, fileName)

	currentEntry = &harFile.Log.Entries[0]
	generalTextView = createTextPane("General")
	requestTextView = createTextPane("Request Headers")
	responseTextView = createTextPane("Response Headers")

	detailsFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(generalTextView, 0, 1, false).
		AddItem(tview.NewFlex().
			SetDirection(tview.FlexColumn).
			AddItem(requestTextView, 0, 1, false).
			AddItem(responseTextView, 0, 1, false), 0, 4, false)

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(networkTable, 0, 1, true).
		AddItem(detailsFlex, 0, 3, false)

	if err := app.SetRoot(flex, true).SetFocus(networkTable).Run(); err != nil {
		panic(err)
	}
}
