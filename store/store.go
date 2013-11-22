package store

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Store struct {
	path  string
	Feeds []Feed `json:"feeds"`
}

type Feed struct {
	Latest int64  `json:"latest"`
	Title  string `json:"title"`
	Url    string `json:"url"`
}

func Load(path string) (s *Store, err error) {
	s = &Store{path: path}

	file, err := os.Open(s.path)
	if err != nil {
		return
	}

	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	if err = json.Unmarshal(data, &s); err != nil {
		return
	}

	return
}

func (s *Store) Add(title string, url string) {
	feed := Feed{Title: title, Url: url}
	s.Feeds = append(s.Feeds, feed)
}

func (s *Store) Save() (err error) {
	file, err := os.Create(s.path)
	if err != nil {
		return
	}

	defer file.Close()
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return
	}

	if _, err = file.Write(data); err != nil {
		return
	}

	return
}
