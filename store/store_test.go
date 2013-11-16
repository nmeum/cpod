package store

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	store, err := Load("testdata/testLoad.json")
	if err != nil {
		t.Fatal(err)
	}

	f := store.Feeds[0]
	if f.Latest != 42 {
		t.Fatalf("Expected %q - got %q", 42, f.Latest)
	}

	if f.Title != "Foo" {
		t.Fatalf("Expected %q - got %q", "Foo", f.Title)
	}

	if f.Url != "http://example.com/rss.xml" {
		t.Fatalf("Expected %q - get %q", "http://example.com/rss.xml", f.Url)
	}
}

func TestAdd(t *testing.T) {
	store := new(Store)
	store.Add("test", "http://example.com/test.xml")

	if store.Feeds[0].Title != "test" {
		t.Fatalf("Expected %q - got %q", "test", store.Feeds[0].Title)
	}

	if store.Feeds[0].Url != "http://example.com/test.xml" {
		t.Fatalf("Expected %q - got %q", "http://example.com/test.xml", store.Feeds[0].Url)
	}
}

func TestSave(t *testing.T) {
	store, err := Load("testdata/testSave.json")
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}

	store.Add("test", "http://example.com/test.xml")
	if err := store.Save(); err != nil {
		t.Fatal(err)
	}

	loaded, err := Load("testdata/testSave.json")
	if err != nil {
		t.Fatal(err)
	}

	if loaded.Feeds[0].Title != "test" {
		t.Fatalf("Expected %q - got %q", "test", loaded.Feeds[0].Title)
	}

	os.Remove("testdata/testSave.json")
}
