package har

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func Load(f string) HarFile {
	file, err := os.Open(f)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	var har HarFile
	err = json.Unmarshal(b, &har)
	if err != nil {
		log.Fatal(err)
	}

	return har
}

type HarFile struct {
	Log Log `json:"log"`
}

type Log struct {
	Version string   `json:"version"`
	Creator Creator  `json:"creator"`
	Pages   []Page   `json:"pages"`
	Entries []Entrie `json:"entries"`
}

type Creator struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Page struct {
	StartedDateTime time.Time   `json:"startedDateTime"`
	ID              string      `json:"id"`
	Title           string      `json:"title"`
	PageTimings     PageTimings `json:"pageTimings"`
}

type PageTimings struct {
	OnContentLoad float64 `json:"onContentLoad"`
	OnLoad        float64 `json:"onLoad"`
}

type Entrie struct {
	Initiator       Initiator `json:"_initiator,omitempty"`
	Priority        string    `json:"_priority"`
	ResourceType    string    `json:"_resourceType"`
	Cache           Cache     `json:"cache"`
	Connection      string    `json:"connection,omitempty"`
	Pageref         string    `json:"pageref"`
	Request         Request   `json:"request"`
	Response        Response  `json:"response"`
	ServerIPAddress string    `json:"serverIPAddress"`
	StartedDateTime time.Time `json:"startedDateTime"`
	Time            float64   `json:"time"`
	Timings         Timings   `json:"timings"`
}

type Cache struct {
}
type Headers struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type Cookies struct {
	Name     string      `json:"name"`
	Value    string      `json:"value"`
	Expires  interface{} `json:"expires"`
	HTTPOnly bool        `json:"httpOnly"`
	Secure   bool        `json:"secure"`
}
type Request struct {
	Method      string        `json:"method"`
	URL         string        `json:"url"`
	HTTPVersion string        `json:"httpVersion"`
	Headers     []Headers     `json:"headers"`
	QueryString []interface{} `json:"queryString"`
	Cookies     []Cookies     `json:"cookies"`
	HeadersSize int           `json:"headersSize"`
	BodySize    int           `json:"bodySize"`
}
type Content struct {
	Size     int    `json:"size"`
	MimeType string `json:"mimeType"`
	Text     string `json:"text"`
}
type Response struct {
	Status       int           `json:"status"`
	StatusText   string        `json:"statusText"`
	HTTPVersion  string        `json:"httpVersion"`
	Headers      []Headers     `json:"headers"`
	Cookies      []interface{} `json:"cookies"`
	Content      Content       `json:"content"`
	RedirectURL  string        `json:"redirectURL"`
	HeadersSize  int           `json:"headersSize"`
	BodySize     int           `json:"bodySize"`
	TransferSize int           `json:"_transferSize"`
	Error        interface{}   `json:"_error"`
}
type Timings struct {
	Blocked         float64 `json:"blocked"`
	DNS             float64 `json:"dns"`
	Ssl             float64 `json:"ssl"`
	Connect         float64 `json:"connect"`
	Send            float64 `json:"send"`
	Wait            float64 `json:"wait"`
	Receive         float64 `json:"receive"`
	BlockedQueueing float64 `json:"_blocked_queueing"`
}
type Initiator struct {
	Type       string `json:"type"`
	URL        string `json:"url"`
	LineNumber int    `json:"lineNumber"`
	Stack      Stack  `json:"stack"`
}
type CallFrames struct {
	FunctionName string `json:"functionName"`
	ScriptID     string `json:"scriptId"`
	URL          string `json:"url"`
	LineNumber   int    `json:"lineNumber"`
	ColumnNumber int    `json:"columnNumber"`
}
type Parent struct {
	Description string       `json:"description"`
	CallFrames  []CallFrames `json:"callFrames"`
}
type Stack struct {
	CallFrames []CallFrames `json:"callFrames"`
	Parent     Parent       `json:"parent"`
}
