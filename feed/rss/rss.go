package rss

import (
	"encoding/xml"
)

type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Title   string   `xml:"channel>title"`
	Link    string   `xml:"channel>link"`
	Items   []Item   `xml:"channel>item"`
}

type Item struct {
	PubDate   string    `xml:"pubDate"`
	Title     string    `xml:"title"`
	Link      string    `xml:"link"`
	Enclosure Enclosure `xml:"enclosure"`
}

type Enclosure struct {
	Type string `xml:"type,attr"`
	Url  string `xml:"url,attr"`
}
