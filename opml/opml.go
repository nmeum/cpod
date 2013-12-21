package opml

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

type Opml struct {
	XMLName  xml.Name  `xml:"opml"`
	Version  string    `xml"version,attr"`
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

func Load(path string) (o *Opml, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	if err = xml.Unmarshal(data, &o); err != nil {
		return
	}

	return
}
