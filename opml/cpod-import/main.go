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
)

func usage() {
	fmt.Fprintf(os.Stderr, "USAGE: cpod-import OPMLFILE...\n")
	os.Exit(1)
}

func load(files []string) (out []opml.Outline, err error) {
	for _, fp := range files {
		var file *os.File
		file, err = os.Open(fp)
		if err != nil {
			return
		}
		defer file.Close()

		var op *opml.OPML
		op, err = opml.Load(file)
		if err != nil {
			return
		}

		for _, o := range op.Outlines {
			out = append(out, o)
		}
	}

	return
}

func main() {
	if len(os.Args) <= 1 {
		usage()
	}

	outlines, err := load(os.Args[1:])
	if err != nil {
		panic(err)
	}

	for _, outline := range outlines {
		fmt.Println(outline.URL)
	}
}
