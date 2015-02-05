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

func usage() {
	fmt.Fprintf(os.Stderr, "USAGE: cpod-export FILE\n")
	os.Exit(1)
}

func warn(err error) {
	fmt.Fprintf(os.Stderr, err.Error()+"\n")
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
	opmlFile := opml.Create("Podcast subscriptions")

	for cast := range storage.Fetch() {
		wg.Add(1)
		go func(p store.Podcast) {
			feed := p.Feed
			if p.Error != nil {
				warn(p.Error)
				return
			}

			opmlFile.Add(feed.Title, feed.Type, p.URL)
			wg.Done()
		}(cast)
	}

	wg.Wait()
	if err = opmlFile.Save(os.Args[1]); err != nil {
		panic(err)
	}
}
