package util

import (
	"testing"
)

type testpair struct {
	inputData string
	expected  string
}

func TestFilename(t *testing.T) {
	testpairs := []testpair{
		{"http://example.com/foo/bar/foo/bar/foo.mp3", "foo.mp3"},
		{"http://example.com/bar.opus?foo=bar&bar=foo", "bar.opus"},
	}

	for _, p := range testpairs {
		f, err := Filename(p.inputData)
		if err != nil {
			t.Fatal(err)
		}

		if f != p.expected {
			t.Fatalf("Expected %q - got %q", p.expected, f)
		}
	}
}
