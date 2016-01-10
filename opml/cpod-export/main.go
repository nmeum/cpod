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
)

// OPML document title
const title = "Podcast subscriptions"

func usage() {
	fmt.Fprintf(os.Stderr, "USAGE: cpod-export URLFILE...\n")
	os.Exit(1)
}

func createOpml(urls []string) *opml.OPML {
	var wg sync.WaitGroup
	opmlFile := opml.Create(title)

	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			resp, err := util.Get(u)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				return
			}
			defer resp.Body.Close()

			feed, err := feedparser.Parse(resp.Body)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			}

			opmlFile.Add(feed.Title, feed.Type, u)
		}(url)
	}

	wg.Wait()
	return opmlFile
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

	data, err := xml.MarshalIndent(createOpml(urls), "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}
