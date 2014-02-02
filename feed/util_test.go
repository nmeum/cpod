package feed

import (
	"github.com/nmeum/cpod/feed/atom"
	"testing"
)

func TestFindLink(t *testing.T) {
	links := []atom.Link{
		{"text/html", "http://example.com/my_link", ""},
		{"audio/ogg", "http://example.com/my_foo", "enclosure"},
	}

	link := findLink(links)
	if link != links[0] {
		t.Fatalf("Expected %q - got %q", link, links[0])
	}
}

func TestFindAttachment(t *testing.T) {
	links := []atom.Link{
		{"text/html", "http://example.org/foo", "alternate"},
		{"image/png", "http://example.org/bar", "enclosure"},
		{"text/html", "http://example.org/baz", ""},
	}

	link := findAttachment(links)
	if link != links[1] {
		t.Fatalf("Expected %q - got %q", link, links[1])
	}
}
