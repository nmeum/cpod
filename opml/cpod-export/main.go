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

package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"github.com/nmeum/cpod/opml"
	"github.com/nmeum/cpod/util"
	"github.com/nmeum/go-feedparser"
	"os"
	"sync"
	"time"
)

// OPML document title
const title = "Podcast subscriptions"

func usage() {
	fmt.Fprintf(os.Stderr, "USAGE: cpod-export URLFILE...\n")
	os.Exit(1)
}

func encodeOutline(url string, enc *xml.Encoder) {
	resp, err := util.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}
	defer resp.Body.Close()

	feed, err := feedparser.Parse(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	enc.Encode(opml.Outline{
		Text: feed.Title,
		Type: feed.Type,
		URL:  url,
	})
}

func encodeOPML(enc *xml.Encoder, urls []string) {
	rootElem := xml.StartElement{Name: xml.Name{Local: "opml"}}
	rootElem.Attr = []xml.Attr{
		{Name: xml.Name{Local: "version"}, Value: opml.Version},
	}

	bodyElem := xml.StartElement{Name: xml.Name{Local: "body"}}
	head := opml.Head{
		Title:   title,
		Created: time.Now().Format(time.RFC1123Z),
	}

	enc.EncodeToken(rootElem)
	enc.Encode(head)
	enc.EncodeToken(bodyElem)

	enc.Flush()
	var wg sync.WaitGroup

	wg.Add(len(urls))
	for _, url := range urls {
		go func(u string) {
			encodeOutline(u, enc)
			wg.Done()
		}(url)
	}

	wg.Wait()
	enc.EncodeToken(bodyElem.End())
	enc.EncodeToken(rootElem.End())
	enc.Flush()
}

func main() {
	if len(os.Args) < 1 {
		usage()
	}

	var urls []string
	for _, fp := range os.Args[1:] {
		file, err := os.Open(fp)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			urls = append(urls, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}

	encodeOPML(xml.NewEncoder(os.Stdout), urls)
}
