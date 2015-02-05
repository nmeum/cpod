package main

import (
	"fmt"
	"github.com/nmeum/cpod/opml"
	"github.com/nmeum/cpod/store"
	"github.com/nmeum/cpod/util"
	"github.com/nmeum/go-feedparser"
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

	for _, url := range storage.URLs {
		wg.Add(1)
		go func(u string) {
			resp, err := util.Get(url)
			if err != nil {
				warn(err)
				return
			}

			reader := resp.Body
			defer resp.Body.Close()

			cast, err := feedparser.Parse(reader, store.Parsers)
			if err != nil {
				warn(err)
				return
			}

			opmlFile.Add(cast.Title, cast.Type, u)
			wg.Done()
		}(url)
	}

	wg.Wait()
	if err = opmlFile.Save(os.Args[1]); err != nil {
		panic(err)
	}
}
