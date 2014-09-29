package feed

import (
	"fmt"
	"sort"
	"testing"
	"time"
)

var testItems = []Item{
	{"Number 2", "http://example.com/three.html", time.Unix(0, 0), ""},
	{"Number 0", "http://example.com/first.html", time.Now(), ""},
	{"Number 1", "http://example.com/second.html", time.Unix(1412004199, 0), ""},
}

func TestByDate(t *testing.T) {
	sort.Sort(ByDate(testItems))
	for n, i := range testItems {
		if i.Title != fmt.Sprintf("Number %d", n) {
			t.Fatalf("Expected %q - got %q", fmt.Sprintf("Number %d", n), i.Title)
		}
	}
}
