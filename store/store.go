package store

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Store struct {
	path  string
	Feeds []Feed
}

type Feed struct {
	Latest int64  `json:"latest"`
	Title  string `json:"title"`
	Type   string `json:"type"`
	Url    string `json:"url"`
}

func Load(path string) (s *Store, err error) {
	s = &Store{path: path}

	data, err := ioutil.ReadFile(s.path)
	if err != nil {
		return
	}

	if err = json.Unmarshal(data, &s.Feeds); err != nil {
		return
	}

	return
}

func (s *Store) Add(title string, ftype string, url string) {
	feed := Feed{
		Title: title,
		Type:  ftype,
		Url:   url,
	}

	s.Feeds = append(s.Feeds, feed)
}

func (s *Store) Save() (err error) {
	file, err := os.Create(s.path)
	if err != nil {
		return
	}

	defer file.Close()
	data, err := json.MarshalIndent(s.Feeds, "", "\t")
	if err != nil {
		return
	}

	if _, err = file.Write(data); err != nil {
		return
	}

	return
}
