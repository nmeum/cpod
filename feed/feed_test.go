package feed

import (
	"github.com/nmeum/cpod/feed/atom"
	"github.com/nmeum/cpod/feed/rss"
	"testing"
)

var testItem = Item{
	Title:      "Some Title",
	Link:       "http://example.org/posts/some_post.html",
	Date:       "Tue, 10 Jun 2003 04:00:00 GMT",
	Attachment: "http://example.org/posts/some_post.ogg",
}

type testpair struct {
	Url  string
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
		feed, err := Parse(test.Url)
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
		PubDate:   "Tue, 10 Jun 2003 04:00:00 GMT",
		Title:     "Some Title",
		Link:      "http://example.org/posts/some_post.html",
		Enclosure: rss.Enclosure{"audio/ogg", "http://example.org/posts/some_post.ogg"},
	}

	feed := convertRss(rssFeed)
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
		Published: "Tue, 10 Jun 2003 04:00:00 GMT",
		Title:     "Some Title",
	}

	atomFeed.Entries[0].Links = make([]atom.Link, 2)
	atomFeed.Entries[0].Links[0] = atom.Link{Type: "text/html", Href: "http://example.org/posts/some_post.html"}
	atomFeed.Entries[0].Links[1] = atom.Link{
		Type: "audio/ogg",
		Href: "http://example.org/posts/some_post.ogg",
		Rel:  "enclosure",
	}

	feed := convertAtom(atomFeed)
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
