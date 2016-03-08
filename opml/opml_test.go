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
	"encoding/xml"
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	outline := Outline{
		XMLName: xml.Name{Local: "outline"},
		Text:    "Chaosradio",
		Type:    "rss",
		URL:     "http://chaosradio.ccc.de/chaosradio-latest.rss",
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

	if len(o.Body.Outlines) != 1 {
		t.Fatalf("Expected %d - got %d", 1, len(o.Body.Outlines))
	}

	if o.Body.Outlines[0] != outline {
		t.Fatalf("Expected %q - got %q", outline, o.Body.Outlines[0])
	}

	if o.Version != Version {
		t.Fatalf("Expected %q - got %q", "2.0", o.Version)
	}

	if o.Head.Title != "Subscriptions" {
		t.Fatalf("Expected %q - got %q", "Subscriptions", o.Head.Title)
	}

	if o.Head.Created != "Wed, 15 May 2013 19:30:58 +0200" {
		t.Fatalf("Expected %q - got %q", "Wed, 15 May 2013 19:30:58 +0200", o.Head.Created)
	}
}
