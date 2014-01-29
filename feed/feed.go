package feed

import (
	"encoding/xml"
	"errors"
	"github.com/nmeum/cpod/feed/atom"
	"github.com/nmeum/cpod/feed/rss"
	"io/ioutil"
	"net/http"
)

type Feed struct {
	Title string
	Type  string
	Link  string
	Items []Item
}

type Item struct {
	Title      string
	Link       string
	Date       string
	Attachment string
}

func Parse(url string) (f Feed, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var rssFeed rss.Feed
	var atomFeed atom.Feed

	if err := xml.Unmarshal(body, &rssFeed); err == nil {
		f = convertRss(rssFeed)
		f.Type = "rss"
	} else if err := xml.Unmarshal(body, &atomFeed); err == nil {
		f = convertAtom(atomFeed)
		f.Type = "atom"
	} else {
		err = errors.New("Unknown feed type")
	}

	return
}

func convertRss(r rss.Feed) (f Feed) {
	f.Title = r.Title
	f.Link = r.Link

	for _, i := range r.Items {
		item := Item{
			Title:      i.Title,
			Link:       i.Link,
			Date:       i.PubDate,
			Attachment: i.Enclosure.Url,
		}

		f.Items = append(f.Items, item)
	}

	return
}

func convertAtom(a atom.Feed) (f Feed) {
	f.Title = a.Title
	f.Link = findLink(a.Links).Href

	for _, e := range a.Entries {
		item := Item{
			Title:      e.Title,
			Link:       findLink(e.Links).Href,
			Date:       e.Published,
			Attachment: findAttachment(e.Links).Href,
		}

		f.Items = append(f.Items, item)
	}

	return
}
