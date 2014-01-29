package feed

import (
	"testing"
)

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
