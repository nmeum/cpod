package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

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

	os.Remove(lockPath)
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
		{"B$:(=== >$-%)/A/R", "B-A-R"},
	}

	for _, test := range tests {
		e, err := Escape(test.unescaped)
		if err != nil {
			t.Fatal(err)
		}

		if e != test.escaped {
			t.Fatalf("Expected %q - got %q", test.escaped, e)
		}
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
