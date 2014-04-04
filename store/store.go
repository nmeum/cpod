package store

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Store struct {
	path     string
	Podcasts []*Podcast // TODO doesn't need to be a pointer
}

type Podcast struct {
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

	if err = json.Unmarshal(data, &s.Podcasts); err != nil {
		return
	}

	return
}

func (s *Store) Add(title, ftype, url string) {
	cast := &Podcast{
		Title: title,
		Type:  ftype,
		Url:   url,
	}

	s.Podcasts = append(s.Podcasts, cast)
}

func (s *Store) Save() (err error) {
	file, err := os.Create(s.path)
	if err != nil {
		return
	}

	defer file.Close()
	data, err := json.MarshalIndent(s.Podcasts, "", "\t")
	if err != nil {
		return
	}

	if _, err = file.Write(data); err != nil {
		return
	}

	return
}
