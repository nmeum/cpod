package rss

import (
	"encoding/xml"
	"io/ioutil"
	"testing"
)

var rssFeed Feed

func TestFeed(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/feed.rss")
	if err != nil {
		t.Fatal(err)
	}

	if err := xml.Unmarshal(data, &rssFeed); err != nil {
		t.Fatal(err)
	}

	if rssFeed.Title != "Some Title" {
		t.Fatalf("Expected %q - got %q", "Some Title", rssFeed.Title)
	}

	if rssFeed.Link != "http://example.org" {
		t.Fatalf("Expected %q - got %q", "http://example.org", rssFeed.Link)
	}
}

func TestItem(t *testing.T) {
	item := rssFeed.Items[0]

	if item.PubDate != "Tue, 20 May 2003 08:56:02 GMT" {
		t.Fatalf("Expected %q - got %q", "Tue, 20 May 2003 08:56:02 GMT", item.PubDate)
	}

	if item.Title != "Test Post" {
		t.Fatalf("Expected %q - got %q", "Test Post", item.Title)
	}

	if item.Link != "http://example.org/posts/test.html" {
		t.Fatalf("Expected %q - got %q", "http://example.org/posts/test.html", item.Link)
	}
}

func TestEnclosure(t *testing.T) {
	enclosure := rssFeed.Items[0].Enclosure

	if enclosure.Type != "audio/ogg" {
		t.Fatalf("Expected %q - got %q", "audio/ogg", enclosure.Type)
	}

	if enclosure.URL != "http://example.org/posts/test.ogg" {
		t.Fatalf("Expected %q - got %q", "http://example.org/posts/test.ogg", enclosure.URL)
	}
}
