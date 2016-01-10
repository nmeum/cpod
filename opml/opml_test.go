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

package opml

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreate(t *testing.T) {
	o := Create("Test subscriptions")
	if o.Title != "Test subscriptions" {
		t.Fatalf("Expected %q - got %q", "Test subscriptions", o.Title)
	}

	if o.Version != version {
		t.Fatalf("Expected %q - got %q", version, o.Version)
	}
}

func TestLoad(t *testing.T) {
	outline := Outline{
		Text: "Chaosradio",
		Type: "rss",
		URL:  "http://chaosradio.ccc.de/chaosradio-latest.rss",
	}

	file, err := os.Open("testdata/testLoad.opml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	o, err := Load(file)
	if err != nil {
		t.Fatal(err)
	}

	if len(o.Outlines) != 1 {
		t.Fatalf("Expected %d - got %d", 1, len(o.Outlines))
	}

	if o.Outlines[0] != outline {
		t.Fatalf("Expected %q - got %q", outline, o.Outlines[0])
	}

	if o.Version != version {
		t.Fatalf("Expected %q - got %q", "2.0", o.Version)
	}

	if o.Title != "Subscriptions" {
		t.Fatalf("Expected %q - got %q", "Subscriptions", o.Title)
	}

	if o.Created != "Wed, 15 May 2013 19:30:58 +0200" {
		t.Fatalf("Expected %q - got %q", "Wed, 15 May 2013 19:30:58 +0200", o.Created)
	}
}

func TestAdd(t *testing.T) {
	testOutline := Outline{
		Text: "Testcast",
		Type: "atom",
		URL:  "http://testcast.com/atom-feed.xml",
	}

	o := new(OPML)
	o.Add("Testcast", "atom", "http://testcast.com/atom-feed.xml")

	if o.Outlines[0] != testOutline {
		t.Fatalf("Expected %q - got %q", testOutline, o.Outlines[0])
	}

	if len(o.Outlines) != 1 {
		t.Fatalf("Expected %d - got %d", 1, len(o.Outlines))
	}
}

func TestSave(t *testing.T) {
	o := Create("Podcasts")
	o.Add("Somecast", "rss", "http://somecast.io/feed.rss")

	testPath := filepath.Join(os.TempDir(), "testSave.opml")
	if err := o.Save(testPath); err != nil {
		t.Fatal(err)
	}
	defer os.Remove(testPath)

	file, err := os.Open(testPath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	loaded, err := Load(file)
	if err != nil {
		t.Fatal(err)
	}

	if loaded.Title != "Podcasts" {
		t.Fatal(err)
	}
}
