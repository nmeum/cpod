package store

import (
	"testing"
)

func TestLoad(t *testing.T) {
	store, err := Load("testdata/testLoad.txt")
	if err != nil {
		t.Fatal(err)
	}

	if store.URLs[0] != "http://feeds.thisamericanlife.org/talpodcast" {
		t.Fatalf("Expected %q - got %q", "http://feeds.thisamericanlife.org/talpodcast", store.URLs[0])
	}

	if store.URLs[1] != "http://www.npr.org/rss/podcast.php?id=510294" {
		t.Fatalf("Expected %q - got %q", "http://www.npr.org/rss/podcast.php?id=510294", store.URLs[1])
	}
}

func TestFetch(t *testing.T) {
	store := &Store{[]string{"http://feeds.thisamericanlife.org/talpodcast"}}
	channel := store.Fetch()

	feed := <-channel
	if feed.Title != "This American Life" {
		t.Fatalf("Expected %q - got %q", "This American Life", feed.Title)
	}
}
