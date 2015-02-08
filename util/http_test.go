package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
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

func TestGet(t *testing.T) {
	expected := "Success\n"
	th := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, expected)
	}

	ts := httptest.NewServer(http.HandlerFunc(th))
	defer ts.Close()

	resp, err := Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	result := string(data)
	if result != expected {
		t.Fatalf("Expected %q - got %q", expected, result)
	}
}
