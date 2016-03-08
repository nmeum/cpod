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
	"fmt"
	"github.com/nmeum/cpod/opml"
	"os"
	"sync"
)

func usage() {
	fmt.Fprintf(os.Stderr, "USAGE: cpod-import OPMLFILE...\n")
	os.Exit(1)
}

func loadFiles(ch chan<- opml.Outline, files []string) {
	var wg sync.WaitGroup
	wg.Add(len(files))

	for _, path := range files {
		go func(fp string) {
			defer wg.Done()
			file, err := os.Open(fp)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				return
			}
			defer file.Close()

			op, err := opml.Load(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				return
			}

			for _, out := range op.Body.Outlines {
				ch <- out
			}
		}(path)
	}

	wg.Wait()
	close(ch)
}

func main() {
	if len(os.Args) <= 1 {
		usage()
	}

	ch := make(chan opml.Outline)
	go loadFiles(ch, os.Args[1:])

	for out := range ch {
		fmt.Println(out.URL)
	}
}
