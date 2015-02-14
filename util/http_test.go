package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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
		{"http://example.com/", "unnamed"},
		{"http://example.com", "unnamed"},
	}

	for _, p := range testpairs {
		f, err := filename(p.inputData)
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

func TestGetFile1(t *testing.T) {
	expected := "Hello World!\n"
	testFile := filepath.Join("testdata", "hello.txt")

	th := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, testFile)
	}

	ts := httptest.NewServer(http.HandlerFunc(th))
	defer ts.Close()

	fp, err := GetFile(ts.URL, os.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(fp)

	data, err := ioutil.ReadFile(fp)
	if err != nil {
		t.Fatal(err)
	}

	result := string(data)
	if result != expected {
		t.Fatalf("Expected %q - got %q", expected, result)
	}
}
