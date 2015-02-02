package store

import (
	"bufio"
	"github.com/nmeum/go-feedparser"
	"github.com/nmeum/go-feedparser/atom"
	"github.com/nmeum/go-feedparser/rss"
	"net/http"
	"os"
)

type Store struct {
	URLs    []string
	parsers []feedparser.FeedFunc
}

func Load(path string) (s *Store, err error) {
	s = new(Store)
	s.parsers = []feedparser.FeedFunc{rss.Parse, atom.Parse}

	file, err := os.Open(path)
	if err != nil {
		return
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		s.URLs = append(s.URLs, scanner.Text())
	}

	err = scanner.Err()
	return
}

func (s *Store) Fetch() <-chan feedparser.Feed {
	out := make(chan feedparser.Feed)
	go func() {
		for _, url := range s.URLs {
			resp, err := http.Get(url)
			if err != nil {
				continue
			}

			reader := resp.Body
			defer reader.Close()

			f, err := feedparser.Parse(reader, s.parsers)
			if err == nil {
				out <- f
			}
		}

		close(out)
	}()

	return out
}
