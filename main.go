package main

import (
	"fmt"
	"os"
	"log"
	"strconv"

	"github.com/jroyal/gohar/har"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func setTableHeaders(table *tview.Table) {
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
}

func main() {

	f, err := os.OpenFile("debug.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	harFile := har.Load("test.almightyzero.com.har")
	app := tview.NewApplication()
	table := tview.NewTable().
		SetBorders(false).
		SetSelectable(true, false).
		SetSeparator(' ')
	

	
	setTableHeaders(table)
	for row, entry := range harFile.Log.Entries {
		row++ // bump the row up to skip past the header row

		domain := tview.NewTableCell(entry.Request.URL).SetMaxWidth(50)
		status := tview.NewTableCell(strconv.Itoa(entry.Response.Status))
		rType := tview.NewTableCell(entry.ResourceType)
		initiator := tview.NewTableCell(entry.Initiator.URL)
		size := tview.NewTableCell(strconv.Itoa(entry.Response.TransferSize))
		time := tview.NewTableCell(fmt.Sprintf("%f", entry.Time))
		table.SetCell(row, 0, domain)
		table.SetCell(row, 1, status)
		table.SetCell(row, 2, rType)
		table.SetCell(row, 3, initiator)
		table.SetCell(row, 4, size)
		table.SetCell(row, 5, time)

	}
	table.Select(0, 0).SetFixed(1, 0).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape {
			app.Stop()
		}
	}).SetSelectedFunc(func(row int, column int) {
		table.GetCell(row, column).SetTextColor(tcell.ColorRed)
	})
	if err := app.SetRoot(table, true).SetFocus(table).Run(); err != nil {
		panic(err)
	}
}
