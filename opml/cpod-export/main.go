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
	"sync"
)

// OPML document title
const title = "Podcast subscriptions"

func usage() {
	fmt.Fprintf(os.Stderr, "USAGE: cpod-export FILE\n")
	os.Exit(1)
}

func warn(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err.Error())
}

func main() {
	if len(os.Args) <= 0 {
		usage()
	}

	storeDir := filepath.Join(util.EnvDefault("XDG_CONFIG_HOME", ".config"), "cpod")
	storage, err := store.Load(filepath.Join(storeDir, "urls"))
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	opmlFile := opml.Create(title)

	for cast := range storage.Fetch() {
		wg.Add(1)
		go func(p store.Podcast) {
			defer wg.Done()
			if p.Error != nil {
				warn(p.Error)
				return
			}

			feed := p.Feed
			opmlFile.Add(feed.Title, feed.Type, p.URL)

			return
		}(cast)
	}

	wg.Wait()
	if err = opmlFile.Save(os.Args[1]); err != nil {
		panic(err)
	}
}
