package opml

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"time"
)

type Opml struct {
	XMLName  xml.Name  `xml:"opml"`
	Version  string    `xml:"version,attr"`
	Head     Head      `xml:"head"`
	Outlines []Outline `xml:"body>outline"`
}

type Head struct {
	Title   string `xml:"title"`
	Created string `xml:"dateCreated"`
}

type Outline struct {
	Text   string `xml:"text,attr"`
	Type   string `xml:"type,attr"`
	XmlUrl string `xml:"xmlUrl,attr"`
}

func New(title string) (o *Opml) {
	o = &Opml{Version: "2.0"}
	o.Head = Head{
		Title:   title,
		Created: time.Now().Format(time.RFC1123Z),
	}

	return
}

func Load(path string) (o *Opml, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}

	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	if err = xml.Unmarshal(data, &o); err != nil {
		return
	}

	return
}

func (o *Opml) Add(text string, ftype string, url string) {
	outline := Outline{
		Text: text,
		Type: ftype,
		XmlUrl: url,
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

	if _, err = file.Write([]byte(xml.Header)); err != nil {
		return
	}

	if _, err = file.Write(data); err != nil {
		return
	}

	return
}
