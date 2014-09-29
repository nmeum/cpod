package atom

import (
	"encoding/xml"
	"github.com/nmeum/cpod/feed"
)

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Title   string   `xml:"title"`
	Links   []Link   `xml:"link"`
	Entries []Entry  `xml:"entry"`
}

type Entry struct {
	Published string `xml:"published"`
	Title     string `xml:"title"`
	Links     []Link `xml:"link"`
}

type Link struct {
	Type string `xml:"type,attr"`
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
}

func Parse(data []byte) (f feed.Feed, err error) {
	var atom Feed
	if err = xml.Unmarshal(data, &atom); err != nil {
		return
	}

	f = feed.Feed{
		Title: atom.Title,
		Type:  "atom",
		Link:  findLink(atom.Links).Href,
	}

	for _, e := range atom.Entries {
		item := feed.Item{
			Title:      e.Title,
			Link:       findLink(e.Links).Href,
			Attachment: findAttachment(e.Links).Href,
		}

		item.Date, err = feed.ParseDate(e.Published)
		if err != nil {
			return
		}

		f.Items = append(f.Items, item)
	}

	return
}

func findLink(links []Link) Link {
	var score int
	var match Link

	for _, link := range links {
		switch {
		case link.Rel == "alternate" && link.Type == "text/html":
			return link
		case score < 3 && link.Type == "text/html":
			score = 3
			match = link
		case score < 2 && link.Rel == "self":
			score = 2
			match = link
		case score < 1 && link.Rel == "":
			score = 1
			match = link
		case &match == nil:
			match = link
		}
	}

	return match
}

func findAttachment(links []Link) Link {
	for _, link := range links {
		if link.Rel == "enclosure" {
			return link
		}
	}

	return Link{}
}
