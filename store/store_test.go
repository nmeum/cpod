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
		"http://feeds.thisamericanlife.org/talpodcast",
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
	url := "http://feeds.thisamericanlife.org/talpodcast"
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
