package feed

import (
	"encoding/xml"
	"errors"
	"github.com/nmeum/cpod/feed/atom"
	"github.com/nmeum/cpod/feed/rss"
	"io/ioutil"
	"net/http"
	"time"
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
	Date       time.Time
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
		f, err = convertRss(rssFeed)
	} else if err := xml.Unmarshal(body, &atomFeed); err == nil {
		f, err = convertAtom(atomFeed)
	} else {
		err = errors.New("unknown feed type")
	}

	return
}

func convertRss(r rss.Feed) (f Feed, err error) {
	f.Title = r.Title
	f.Type = "rss"
	f.Link = r.Link

	for _, i := range r.Items {
		item := Item{
			Title:      i.Title,
			Link:       i.Link,
			Attachment: i.Enclosure.URL,
		}

		item.Date, err = parseDate(i.PubDate)
		if err != nil {
			return
		}

		f.Items = append(f.Items, item)
	}

	return
}

func convertAtom(a atom.Feed) (f Feed, err error) {
	f.Title = a.Title
	f.Type = "atom"
	f.Link = findLink(a.Links).Href

	for _, e := range a.Entries {
		item := Item{
			Title:      e.Title,
			Link:       findLink(e.Links).Href,
			Attachment: findAttachment(e.Links).Href,
		}

		item.Date, err = parseDate(e.Published)
		if err != nil {
			return
		}

		f.Items = append(f.Items, item)
	}

	return
}

func parseDate(date string) (t time.Time, err error) {
	formats := []string{
		time.RFC1123Z, time.RFC1123, time.RFC822Z,
		time.RFC822, time.ANSIC, time.RFC3339,
		time.RFC850, time.RubyDate, time.UnixDate,
		"2 January 2006 15:04:05 -0700", "2 January 2006 15:04:05 MST",
		"2 Jan 2006 15:04:05 -0700", "2 Jan 2006 15:04:05 MST",
		"Mon, 2 Jan 2006 15:04:05 -0700", "Mon, 2 Jan 2006 15:04:05 MST",
		"2006-01-02T15:04:05", "2006-01-02T15:04:05Z",
	}

	for _, format := range formats {
		t, err = time.Parse(format, date)
		if err == nil {
			return
		}
	}

	return
}
