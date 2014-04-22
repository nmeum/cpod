package util

import (
	"github.com/nmeum/cpod/store"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestGet(t *testing.T) {
	path, err := Get("http://paste42.de/6915.txt", os.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(path)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	contents := string(data)
	if contents != "Foobar" {
		t.Fatalf("Expected %q - got %q", "Foobar", contents)
	}
}

func TestLock1(t *testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "testLock")
	if err != nil {
		t.Fatal(err)
	}

	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		t.Fatal(err)
	}

	lockPath := filepath.Join(os.TempDir(), fi.Name())
	if err := Lock(lockPath); !os.IsExist(err) {
		t.Fail()
	}

}

func TestLock2(t *testing.T) {
	lockPath := filepath.Join(os.TempDir(), "lockTest")
	if err := Lock(lockPath); err != nil {
		t.Fatal(err)
	}

	os.Remove(lockPath)
}

func TestEscape(t *testing.T) {
	type testpair struct {
		unescaped string
		escaped   string
	}

	tests := []testpair{
		{"$$foo /", "foo"},
		{"Foo bar, baz!", "Foo-bar-baz"},
		{"LNP007: Foobar!", "LNP007-Foobar"},
		{"$:(=== >$-%)", ""},
	}

	for _, test := range tests {
		e := Escape(test.unescaped)
		if e != test.escaped {
			t.Fatalf("Expected %q - got %q", test.escaped, e)
		}
	}
}

func TestSubscribed(t *testing.T) {
	casts := []*store.Podcast{
		{0, "Foocast", "rss", "http://example.com/foocast.rss"},
		{0, "Barcast", "rss", "http://barcast.org/feed.xml"},
	}

	if !Subscribed(casts, "http://example.com/foocast.rss") {
		t.Fail()
	}

	if !Subscribed(casts, "http://barcast.org/feed.xml") {
		t.Fail()
	}

	if Subscribed(casts, "None sense podcast") {
		t.Fail()
	}
}

func TestEnvDefault1(t *testing.T) {
	if err := os.Setenv("TESTDIR", "/foo"); err != nil {
		t.Fatal(err)
	}

	dir := EnvDefault("TESTDIR", "")
	if dir != "/foo" {
		t.Fatalf("Expected %q - got %q", "/foo", dir)
	}
}

func TestEnvDefault2(t *testing.T) {
	dir := EnvDefault("TESTDIR2", "bar")
	if dir != filepath.Join(os.Getenv("HOME"), "bar") {
		t.Fatalf("Expected %q - got %q", filepath.Join(os.Getenv("HOME"), "bar"), dir)
	}
}
