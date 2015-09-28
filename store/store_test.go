// Copyright (C) 2013-2015 Sören Tempel
//
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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	store, err := Load("testdata/testLoad.txt")
	if err != nil {
		t.Fatal(err)
	}

	urls := []string{
		"http://feed.thisamericanlife.org/talpodcast",
		"http://www.npr.org/rss/podcast.php?id=510294",
	}

	for _, url := range urls {
		if !store.Contains(url) {
			t.Fail()
		}
	}
}

func TestAdd(t *testing.T) {
	url := "http://example.com"
	store := new(Store)

	store.Add(url)
	if !store.Contains(url) {
		t.Fail()
	}
}

func TestContains(t *testing.T) {
	url := "http://foo.com"
	store := &Store{"", []string{url}}

	if !store.Contains(url) {
		t.Fail()
	}

	if store.Contains("http://foo.bar") {
		t.Fail()
	}
}

func TestFetch(t *testing.T) {
	url := "http://feed.thisamericanlife.org/talpodcast"
	store := &Store{"", []string{url}}

	channel := store.Fetch()
	podcast := <-channel

	if podcast.Error != nil {
		t.Fatal(podcast.Error)
	}

	feed := podcast.Feed
	if url != podcast.URL {
		t.Fatalf("Expected %q - got %q", podcast.URL, url)
	}

	expected := "This American Life"
	if feed.Title != expected {
		t.Fatalf("Expected %q - got %q", expected, feed.Title)
	}
}

func TestSave(t *testing.T) {
	url := "http://example.io"
	fp := filepath.Join(os.TempDir(), "testSave")

	store := &Store{fp, []string{url}}
	if err := store.Save(); err != nil {
		t.Fatal(err)
	}

	defer os.Remove(fp)
	data, err := ioutil.ReadFile(fp)
	if err != nil {
		t.Fatal(err)
	}

	expected := url + "\n"
	if string(data) != expected {
		t.Fatalf("Expected %q - got %q", string(data), expected)
	}
}
