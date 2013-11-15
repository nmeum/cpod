package opml

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

type Opml struct {
	Head struct {
		Title   string `xml:"title"`
		Created string `xml:"dateCreated"`
	} `xml:"head"`
	Body struct {
		Outline []struct {
			Text   string `xml:"text,attr"`
			Type   string `xml:"type,attr"`
			XmlUrl string `xml:"xmlUrl,attr"`
		} `xml:"outline"`
	} `xml:"body"`
	XMLName xml.Name `xml:"opml"`
	Version string   `xml:"version,attr"`
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
