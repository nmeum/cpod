package store

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Store struct {
	path  string
	file  *os.File
	Feeds []Feed `json:"feeds"`
}

type Feed struct {
	Latest int64  `json:"latest"`
	Title  string `json:"title"`
	Url    string `json:"url"`
}

// TODO seperate Load() and New()
func Load(path string) (s *Store, err error) {
	s = &Store{path: path}

	s.file, err = os.Open(s.path)
	if os.IsNotExist(err) {
		s.file, err = os.Create(s.path)
		return
	} else if err != nil {
		return
	}

	data, err := ioutil.ReadAll(s.file)
	if err != nil {
		return
	}

	if err = json.Unmarshal(data, &s); err != nil {
		return
	}

	return
}

func (s *Store) Close() {
	s.file.Close()
}

func (s *Store) Add(title string, url string) {
	feed := Feed{Title: title, Url: url}
	s.Feeds = append(s.Feeds, feed)
}

func (s *Store) Save() (err error) {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return
	}

	if _, err = s.file.Write(data); err != nil {
		return
	}

	return
}
