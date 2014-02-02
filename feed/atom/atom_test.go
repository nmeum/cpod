package atom

import (
	"encoding/xml"
	"io/ioutil"
	"testing"
)

var atomFeed Feed

func TestFeed(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/feed.atom")
	if err != nil {
		t.Fatal(err)
	}

	if err := xml.Unmarshal(data, &atomFeed); err != nil {
		t.Fatal(err)
	}

	if atomFeed.Title != "Some Title" {
		t.Fatalf("Expected %q - got %q", "Some Title", atomFeed.Title)
	}

	testLink := Link{
		Type: "text/html",
		Href: "http://example.org/feed.atom",
		Rel:  "self",
	}

	if atomFeed.Links[0] != testLink {
		t.Fatalf("Expected %q - got %q", testLink, atomFeed.Links[0])
	}
}

func TestEntry(t *testing.T) {
	entry := atomFeed.Entries[0]

	if entry.Published != "2013-10-11T23:56:00Z" {
		t.Fatalf("Expected %q - got %q", "2013-10-11T23:56:00Z", entry.Published)
	}

	if entry.Title != "Test Post" {
		t.Fatalf("Expected %q - got %q", "Test Post", entry.Title)
	}

	testLink := Link{
		Type: "text/html",
		Href: "http://example.org/posts/test.html",
		Rel:  "alternate",
	}

	if entry.Links[0] != testLink {
		t.Fatalf("Expected %q - got %q", testLink, entry.Links[0])
	}
}
