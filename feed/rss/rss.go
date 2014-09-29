package rss

import (
	"encoding/xml"
	"github.com/nmeum/cpod/feed"
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
	URL  string `xml:"url,attr"`
}

func Parse(data []byte) (f feed.Feed, err error) {
	var rss Feed
	if err = xml.Unmarshal(data, &rss); err != nil {
		return
	}

	f = feed.Feed{
		Title: rss.Title,
		Type:  "rss",
		Link:  rss.Link,
	}

	for _, i := range rss.Items {
		item := feed.Item{
			Title:      i.Title,
			Link:       i.Link,
			Attachment: i.Enclosure.URL,
		}

		item.Date, err = feed.ParseDate(i.PubDate)
		if err != nil {
			return
		}

		f.Items = append(f.Items, item)
	}

	return
}
