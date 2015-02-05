package store

import (
	"bufio"
	"github.com/nmeum/cpod/util"
	"github.com/nmeum/go-feedparser"
	"github.com/nmeum/go-feedparser/atom"
	"github.com/nmeum/go-feedparser/rss"
	"os"
)

var Parsers = []feedparser.FeedFunc{
	rss.Parse,
	atom.Parse,
}

type Store struct {
	path string
	URLs []string
}

func Load(path string) (s *Store, err error) {
	s = new(Store)
	s.path = path

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

func (s *Store) Add(url string) {
	s.URLs = append(s.URLs, url)
}

func (s *Store) Contains(url string) bool {
	for _, u := range s.URLs {
		if u == url {
			return true
		}
	}

	return false
}

func (s *Store) Fetch() <-chan feedparser.Feed {
	out := make(chan feedparser.Feed)
	go func() {
		for _, url := range s.URLs {
			resp, err := util.Get(url)
			if err != nil {
				continue
			}

			reader := resp.Body
			defer reader.Close()

			f, err := feedparser.Parse(reader, Parsers)
			if err == nil {
				out <- f
			}
		}

		close(out)
	}()

	return out
}

func (s *Store) Save() error {
	file, err := os.OpenFile(s.path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, url := range s.URLs {
		if _, err := file.WriteString(url + "\n"); err != nil {
			return err
		}
	}

	return nil
}
