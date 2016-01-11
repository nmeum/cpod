// Copyright (C) 2013-2016 Sören Tempel
//
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

	if err := os.Remove(lockPath); err != nil {
		t.Fatal(err)
	}
}

func TestLock2(t *testing.T) {
	lockPath := filepath.Join(os.TempDir(), "lockTest")
	if err := Lock(lockPath); os.IsExist(err) {
		t.Fatal(err)
	}

	if err := os.Remove(lockPath); err != nil {
		t.Fatal(err)
	}
}

func TestEscape1(t *testing.T) {
	type testpair struct {
		unescaped string
		escaped   string
	}

	tests := []testpair{
		{"$$foo /", "foo"},
		{"Foo, bar, baz!", "Foo, bar, baz!"},
		{"LNP007: Foobar!", "LNP007: Foobar!"},
		{"B$:(=== >$-%)/A/R", "B:( -%)AR"},
		{"foobar  ", "foobar"},
		{"../foo..", "foo.."},
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

func TestEscape2(t *testing.T) {
	if _, err := Escape(".."); err == nil {
		t.Fail()
	}
}
