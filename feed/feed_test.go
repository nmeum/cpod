package feed

import (
	"testing"
	"github.com/nmeum/cpod/feed/rss"
)

var testItem = Item{
	Title: "Some Title",
	Link: "http://example.org/posts/some_post.html",
	Date: "Tue, 10 Jun 2003 04:00:00 GMT",
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
		Link: "http://example.org",
	}

	rssFeed.Items = make([]rss.Item, 1)
	rssFeed.Items[0] = rss.Item{
		PubDate: "Tue, 10 Jun 2003 04:00:00 GMT",
		Title: "Some Title",
		Link: "http://example.org/posts/some_post.html",
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
