// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package store

import (
	"bufio"
	"github.com/nmeum/cpod/util"
	"github.com/nmeum/go-feedparser"
	"os"
)

// Podcast represents a Podcast loaded from the store.
type Podcast struct {
	// URL to the feed.
	URL string

	// Feed itself.
	Feed feedparser.Feed

	// Error if parsing failed.
	Error error
}

// Store represents a storage backend.
type Store struct {
	// path describes the URL file location.
	path string

	// urls contains all URLs which are part of the URL file.
	urls []string
}

// Load returns and creates a new store with the URL file located
// at the give filepath.
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
		s.urls = append(s.urls, scanner.Text())
	}

	err = scanner.Err()
	return
}

// Add appends a new URL to the store. It doesn't check if the
// given data is a valid URL and it doesn't check if the URL
// is already a part of the store either.
func (s *Store) Add(url string) {
	s.urls = append(s.urls, url)
}

// Contains returns true if the url is already a part of the
// store. If it isn't it returns false.
func (s *Store) Contains(url string) bool {
	for _, u := range s.urls {
		if u == url {
			return true
		}
	}

	return false
}

// Fetch fetches all feeds form the urls and returns a channel
// which contains all podcasts.
func (s *Store) Fetch() <-chan Podcast {
	out := make(chan Podcast)
	go func() {
		for _, url := range s.urls {
			resp, err := util.Get(url)
			if err != nil {
				continue
			}

			reader := resp.Body
			defer reader.Close()

			f, err := feedparser.Parse(reader)
			out <- Podcast{url, f, err}
		}

		close(out)
	}()

	return out
}

// Save writes the URL file to the store path.
func (s *Store) Save() error {
	file, err := os.OpenFile(s.path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	defer file.Close()
	for _, url := range s.urls {
		if _, err := file.WriteString(url + "\n"); err != nil {
			return err
		}
	}

	return nil
}
