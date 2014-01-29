package opml

import (
	"os"
	"testing"
)

func TestCreate(t *testing.T) {
	o := Create("Test subscriptions")

	if o.Title != "Test subscriptions" {
		t.Fatalf("Expected %q - got %q", "Test subscriptions", o.Title)
	}

	if o.Version != "2.0" {
		t.Fatalf("Expected %q - got %q", "2.0", o.Version)
	}
}

func TestLoad(t *testing.T) {
	outline := Outline{
		Text:   "Chaosradio",
		Type:   "rss",
		XmlUrl: "http://chaosradio.ccc.de/chaosradio-latest.rss",
	}

	o, err := Load("testdata/testLoad.opml")
	if err != nil {
		t.Fatal(err)
	}

	if o.Outlines[0] != outline {
		t.Fatalf("Expected %q - got %q", outline, o.Outlines[0])
	}

	if o.Version != "2.0" {
		t.Fatalf("Expected %q - got %q", "2.0", o.Version)
	}

	if o.Title != "Subscriptions" {
		t.Fatalf("Expected %q - got %q", "Subscriptions", o.Title)
	}

	if o.Created != "Wed, 15 May 2013 19:30:58 +0200" {
		t.Fatalf("Expected %q - got %q", "Wed, 15 May 2013 19:30:58 +0200", o.Created)
	}
}

func TestAdd(t *testing.T) {
	testOutline := Outline{
		Text:   "Testcast",
		Type:   "atom",
		XmlUrl: "http://testcast.com/atom-feed.xml",
	}

	o := new(Opml)
	o.Add("Testcast", "atom", "http://testcast.com/atom-feed.xml")

	if o.Outlines[0] != testOutline {
		t.Fatalf("Expected %q - got %q", testOutline, o.Outlines[0])
	}
}

func TestSave(t *testing.T) {
	o := Create("Podcasts")
	o.Add("Somecast", "rss", "http://somecast.io/feed.rss")

	if err := o.Save("testdata/testSave.json"); err != nil {
		t.Fatal(err)
	}

	loaded, err := Load("testdata/testSave.json")
	if err != nil {
		t.Fatal(err)
	}

	if loaded.Title != "Podcasts" {
		t.Fatal(err)
	}

	os.Remove("testdata/testSave.json")
}
