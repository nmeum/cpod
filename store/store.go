package store

import (
	"bufio"
	"github.com/nmeum/cpod/feed"
	"os"
)

type Store struct {
	URLs []string
}

func Load(path string) (s *Store, err error) {
	s = new(Store)

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

func (s *Store) Fetch() <-chan feed.Feed {
	out := make(chan feed.Feed)
	go func() {
		for _, url := range s.URLs {
			f, err := feed.Parse(url)
			if err == nil {
				out <- f
			}
		}

		close(out)
	}()

	return out
}
