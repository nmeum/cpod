package opml

import (
	"testing"
)

func TestLoad(t *testing.T) {
	opml, err := Load("testdata/testLoad.opml")
	if err != nil {
		t.Fatal(err)
	}

	if opml.Head.Title != "Subscriptions" {
		t.Fatalf("Expected %q - got %q", "Subscriptions", opml.Head.Title)
	}

	if opml.Head.Created != "Wed, 15 May 2013 19:30:58 +0200" {
		t.Fatalf("Expected %q - got %q", "Wed, 15 May 2013 19:30:58 +0200", opml.Head.Created)
	}

	if opml.Body.Outlines[0].Text != "Chaosradio" {
		t.Fatalf("Expected %q - got %q", "Chaosradio", opml.Body.Outlines[0].Text)
	}
}
