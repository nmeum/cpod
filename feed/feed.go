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
	Link  string
	Items []Item
}

type Item struct {
	Title      string
	Link       string
	Date       string
	Attachment string
}

func Parse(url string) (f *Feed, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var rssFeed *rss.Feed
	var atomFeed *atom.Feed

	if err := xml.Unmarshal(body, &rssFeed); err == nil {
		f = convertRss(rssFeed)
	} else if err := xml.Unmarshal(body, &atomFeed); err == nil {
		f = convertAtom(atomFeed)
	} else {
		err = errors.New("Unknown feed type")
	}

	return
}
