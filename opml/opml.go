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

// Package opml implements a parser for OPML files.
// See also: http://dev.opml.org/spec2.html
package opml

import (
	"encoding/xml"
	"golang.org/x/net/html/charset"
	"io"
)

// OPML version supported by this library.
const Version = "2.0"

// Root element of an OPML document.
type OPML struct {
	// XML name.
	XMLName xml.Name `xml:"opml"`

	// OPML standard version implemented by this file.
	Version string `xml:"version,attr"`

	// OPML head element of this document.
	Head Head `xml:"head"`

	// OPML body element of this document.
	Body Body `xml:"body"`
}

// Head element containg metadata.
type Head struct {
	// XML name.
	XMLName xml.Name `xml:"head"`

	// Title of the OPML document.
	Title string `xml:"title"`

	// Time the document was created.
	Created string `xml:"dateCreated"`
}

// Body element containg outlines.
type Body struct {
	// XML name.
	XMLName xml.Name `xml:"body"`

	// Array of outlines, each represents a subscription.
	Outlines []Outline `xml:"outline"`
}

// Outline represents an arbitrary OPML outline.
type Outline struct {
	// XML name.
	XMLName xml.Name `xml:"outline"`

	// Text attribute, might contain HTML markup.
	Text string `xml:"text,attr"`

	// Type of file found at the outline URL.
	Type string `xml:"type,attr"`

	// Arbitrary outline URL.
	URL string `xml:"xmlUrl,attr"`
}

// Parse parses an existing OPML using the given reader.
func Load(r io.Reader) (o *OPML, err error) {
	decoder := xml.NewDecoder(r)
	decoder.CharsetReader = charset.NewReaderLabel
	if err = decoder.Decode(&o); err != nil {
		return
	}

	return
}
