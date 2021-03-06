// Copyright (C) 2013-2015 Sören Tempel
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

// Package opml implements a parser for OPML files.
// See also: http://dev.opml.org/spec2.html
package opml

import (
	"encoding/xml"
	"golang.org/x/net/html/charset"
	"os"
	"time"
)

// OPML version supported by this library.
const version = "2.0"

// OPML represent an OPML document.
type OPML struct {
	// XML name.
	XMLName xml.Name `xml:"opml"`

	// OPML standard version implemented by this file.
	Version string `xml:"version,attr"`

	// Title of the OPML document.
	Title string `xml:"head>title"`

	// Time the document was created.
	Created string `xml:"head>dateCreated"`

	// Array of outlines, each represents a subscription.
	Outlines []Outline `xml:"body>outline"`
}

// Outline represents an arbitrary OPML outline.
type Outline struct {
	// Text attribute, might contain HTML markup.
	Text string `xml:"text,attr"`

	// Type of file found at the outline URL.
	Type string `xml:"type,attr"`

	// Arbitrary outline URL.
	URL string `xml:"xmlUrl,attr"`
}

// Create returns a new OPML document with the given title. However,
// this is just syntax sugar. A file is only written after a call Save,
// it's the callers responsibility to do so if desired.
func Create(title string) *OPML {
	return &OPML{
		Version: version,
		Title:   title,
		Created: time.Now().Format(time.RFC1123Z),
	}
}

// Load reads an existing OPML document located at the given path.
func Load(path string) (o *OPML, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	if err = decoder.Decode(&o); err != nil {
		return
	}

	return
}

// Add appends a new outline to the OPML document, even if the outline
// is already a part of the document.
func (o *OPML) Add(text, ftype, url string) {
	outline := Outline{
		Text: text,
		Type: ftype,
		URL:  url,
	}

	o.Outlines = append(o.Outlines, outline)
}

// Save writes an indented version of the OPML document to the given
// file path.
func (o *OPML) Save(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()
	data, err := xml.MarshalIndent(o, "", "\t")
	if err != nil {
		return err
	}

	if _, err = file.WriteString(xml.Header); err != nil {
		return err
	}

	if _, err = file.Write(data); err != nil {
		return err
	}

	return err
}
