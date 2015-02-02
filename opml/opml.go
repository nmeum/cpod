package opml

import (
	"encoding/xml"
	"github.com/nmeum/go-feedparser/util"
	"io/ioutil"
	"os"
	"time"
)

// OPML version supported by this library.
const version = "2.0"

type Opml struct {
	XMLName  xml.Name  `xml:"opml"`
	Version  string    `xml:"version,attr"`
	Title    string    `xml:"head>title"`
	Created  string    `xml:"head>dateCreated"`
	Outlines []Outline `xml:"body>outline"`
}

type Outline struct {
	Text string `xml:"text,attr"`
	Type string `xml:"type,attr"`
	URL  string `xml:"xmlUrl,attr"`
}

func Create(title string) (o *Opml) {
	o = &Opml{
		Version: version,
		Title:   title,
		Created: time.Now().Format(time.RFC1123Z),
	}

	return
}

func Load(path string) (o *Opml, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	if err = util.Unmarshal(data, &o); err != nil {
		return
	}

	return
}

func (o *Opml) Add(text, ftype, url string) {
	outline := Outline{
		Text: text,
		Type: ftype,
		URL:  url,
	}

	o.Outlines = append(o.Outlines, outline)
}

func (o *Opml) Save(path string) (err error) {
	file, err := os.Create(path)
	if err != nil {
		return
	}

	defer file.Close()
	data, err := xml.MarshalIndent(o, "", "  ")
	if err != nil {
		return
	}

	if _, err = file.WriteString(xml.Header); err != nil {
		return
	}

	if _, err = file.Write(data); err != nil {
		return
	}

	return
}
