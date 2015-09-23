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
	"github.com/nmeum/cpod/store"
	"github.com/nmeum/cpod/util"
	"os"
	"path/filepath"
)

func usage() {
	fmt.Fprintf(os.Stderr, "USAGE: cpod-import FILE...\n")
	os.Exit(1)
}

func load(files []string) (out []opml.Outline, err error) {
	for _, file := range files {
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

	storeDir := filepath.Join(util.EnvDefault("XDG_CONFIG_HOME", ".config"), "cpod")
	if err := os.MkdirAll(storeDir, 0755); err != nil {
		panic(err)
	}

	store, err := store.Load(filepath.Join(storeDir, "urls"))
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	outlines, err := load(os.Args)
	if err != nil {
		panic(err)
	}

	for _, outline := range outlines {
		url := outline.URL
		if !store.Contains(url) {
			store.Add(url)
		}
	}

	if err = store.Save(); err != nil {
		panic(err)
	}
}
