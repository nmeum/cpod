package parser

import (
	"errors"
	"github.com/nmeum/cpod/feed"
	"github.com/nmeum/cpod/feed/atom"
	"github.com/nmeum/cpod/feed/rss"
	"io/ioutil"
	"net/http"
	"sort"
)

type FeedFunc func([]byte) (feed.Feed, error)

var parsers = []FeedFunc{
	rss.Parse,
	atom.Parse,
}

func Parse(url string) (f feed.Feed, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	for _, p := range parsers {
		f, err = p(body)
		if err == nil {
			sort.Sort(feed.ByDate(f.Items))
			return
		}
	}

	if err != nil {
		err = errors.New("unknown feed type")
	}

	return
}
