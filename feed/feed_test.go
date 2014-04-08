package feed

import (
	"github.com/nmeum/cpod/feed/atom"
	"github.com/nmeum/cpod/feed/rss"
	"testing"
	"time"
)

var testItem = Item{
	Title:      "Some Title",
	Link:       "http://example.org/posts/some_post.html",
	Date:       time.Unix(1393528968, 0),
	Attachment: "http://example.org/posts/some_post.ogg",
}

type testpair struct {
	URL  string
	Type string
}

func TestParse(t *testing.T) {
	tests := []testpair{
		{"http://cyber.law.harvard.edu/rss/examples/sampleRss092.xml", "rss"},
		{"http://cyber.law.harvard.edu/rss/examples/rss2sample.xml", "rss"},
		{"http://www.heise.de/developer/rss/news-atom.xml", "atom"},
		{"http://blog.case.edu/news/feed.atom", "atom"},
	}

	for _, test := range tests {
		feed, err := Parse(test.URL)
		if err != nil {
			t.Fatal(err)
		}

		if feed.Type != test.Type {
			t.Fatalf("Expected %q - got %q", test.Type, feed.Type)
		}
	}
}

func TestConvertRss(t *testing.T) {
	rssFeed := rss.Feed{
		Title: "Some Title",
		Link:  "http://example.org",
	}

	rssFeed.Items = make([]rss.Item, 1)
	rssFeed.Items[0] = rss.Item{
		PubDate:   "Thu, 27 Feb 2014 20:22:48 +0100",
		Title:     "Some Title",
		Link:      "http://example.org/posts/some_post.html",
		Enclosure: rss.Enclosure{"audio/ogg", "http://example.org/posts/some_post.ogg"},
	}

	feed, err := convertRss(rssFeed)
	if err != nil {
		t.Fatal(err)
	}

	if feed.Title != "Some Title" {
		t.Fatalf("Expected %q - got %q", "Some Title", feed.Title)
	}

	if feed.Type != "rss" {
		t.Fatalf("Expected %q - got %q", "rss", feed.Type)
	}

	if feed.Link != "http://example.org" {
		t.Fatalf("Expected %q - got %q", "http://example.org", feed.Link)
	}

	if feed.Items[0] != testItem {
		t.Fatalf("Expected %q - got %q", testItem, feed.Items[0])
	}
}

func TestConvertAtom(t *testing.T) {
	atomFeed := atom.Feed{
		Title: "Some Title",
	}

	atomFeed.Links = make([]atom.Link, 1)
	atomFeed.Links[0] = atom.Link{Href: "http://example.org"}

	atomFeed.Entries = make([]atom.Entry, 1)
	atomFeed.Entries[0] = atom.Entry{
		Published: "Thu, 27 Feb 2014 20:22:48 +0100",
		Title:     "Some Title",
	}

	atomFeed.Entries[0].Links = make([]atom.Link, 2)
	atomFeed.Entries[0].Links[0] = atom.Link{Type: "text/html", Href: "http://example.org/posts/some_post.html"}
	atomFeed.Entries[0].Links[1] = atom.Link{
		Type: "audio/ogg",
		Href: "http://example.org/posts/some_post.ogg",
		Rel:  "enclosure",
	}

	feed, err := convertAtom(atomFeed)
	if err != nil {
		t.Fatal(err)
	}

	if feed.Title != "Some Title" {
		t.Fatalf("Expected %q - got %q", "Some Title", feed.Title)
	}

	if feed.Type != "atom" {
		t.Fatalf("Expected %q - got %q", "atom", feed.Type)
	}

	if feed.Link != "http://example.org" {
		t.Fatalf("Expected %q - got %q", "http://example.org", feed.Link)
	}

	if feed.Items[0] != testItem {
		t.Fatalf("Expected %q - got %q", testItem, feed.Items[0])
	}
}

func TestParseDate(t *testing.T) {
	testFormat := "Thu, 27 Feb 2014 18:46:18 +0100"
	var timestamp int64 = 1393523178

	date, err := parseDate(testFormat)
	if err != nil {
		t.Fatal(err)
	}

	if date.Unix() != timestamp {
		t.Fatalf("Expected %q - got %q", timestamp, date.Unix())
	}
}
