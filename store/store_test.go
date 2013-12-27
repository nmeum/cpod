package store

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	testFeed := &Feed{
		Latest: 42,
		Title: "Foo",
		Type: "rss",
		Url: "http://example.com/rss.xml",
	}

	store, err := Load("testdata/testLoad.json")
	if err != nil {
		t.Fatal(err)
	}

	if &store.Feeds[0] == testFeed {
		t.Fatalf("Expected %q - got %q", testFeed, &store.Feeds[0])
	}
}

func TestAdd(t *testing.T) {
	testFeed := Feed{
		Title: "Foobar",
		Type: "atom",
		Url: "http://example.io/feed.xml",
	}

	store := new(Store)
	store.Add("Foobar", "atom", "http://example.io/feed.xml")

	if store.Feeds[0] != testFeed {
		t.Fatalf("Expected %q - got %q", testFeed, store.Feeds)
	}
}

func TestSave(t *testing.T) {
	store := Store{path: "testdata/testSave.json"}
	feed := Feed{
		Latest: 1337,
		Title: "Test Feed",
		Type: "atom",
		Url: "http://example.com/testFeed.atom",
	}

	store.Feeds = append(store.Feeds, feed)
	if err := store.Save(); err != nil {
		t.Fatal(err)
	}

	loaded, err := Load("testdata/testSave.json")
	if err != nil {
		t.Fatal(err)
	}

	if loaded.Feeds[0] != feed {
		t.Fatalf("Expected %q - got %q", loaded.Feeds[0], feed)
	}

	os.Remove("testdata/testSave.json")
}
