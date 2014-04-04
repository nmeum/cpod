package store

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	testCast := Podcast{
		Latest: 42,
		Title:  "Foo",
		Type:   "rss",
		Url:    "http://example.com/rss.xml",
	}

	store, err := Load("testdata/testLoad.json")
	if err != nil {
		t.Fatal(err)
	}

	if *store.Podcasts[0] != testCast {
		t.Fatalf("Expected %q - got %q", testCast, *store.Podcasts[0])
	}
}

func TestAdd(t *testing.T) {
	testCast := Podcast{
		Title: "Foobar",
		Type:  "atom",
		Url:   "http://example.io/podcast.xml",
	}

	store := new(Store)
	store.Add(testCast.Title, testCast.Type, testCast.Url)

	if *store.Podcasts[0] != testCast {
		t.Fatalf("Expected %q - got %q", testCast, *store.Podcasts[0])
	}
}

func TestSave(t *testing.T) {
	store := Store{path: "testdata/testSave.json"}
	cast := Podcast{
		Title: "Test Podcast",
		Type:  "atom",
		Url:   "http://example.com/testPodcast.atom",
	}

	store.Add(cast.Title, cast.Type, cast.Url)
	if err := store.Save(); err != nil {
		t.Fatal(err)
	}

	loaded, err := Load("testdata/testSave.json")
	if err != nil {
		t.Fatal(err)
	}

	if *loaded.Podcasts[0] != cast {
		t.Fatalf("Expected %q - got %q", cast, *loaded.Podcasts[0])
	}

	os.Remove("testdata/testSave.json")
}
