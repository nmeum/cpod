// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

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
		{"http://example.org/foobar ", "foobar"},
		{"http://example.com/", "unnamed"},
		{"http://example.com", "unnamed"},
		{"", "unnamed"},
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
