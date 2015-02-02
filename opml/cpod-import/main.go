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
	fmt.Fprintf(os.Stderr, "USAGE: cpod-import [path]\n")
	os.Exit(1)
}

func main() {
	if len(os.Args) <= 0 {
		usage()
	}

	opmlFile, err := opml.Load(os.Args[1])
	if os.IsNotExist(err) {
		usage()
	} else if err != nil {
		panic(err)
	}

	storeDir := filepath.Join(util.EnvDefault("XDG_CONFIG_HOME", ".config"), "cpod")
	if err := os.MkdirAll(storeDir, 0755); err != nil {
		panic(err)
	}

	store, err := store.Load(filepath.Join(storeDir, "urls"))
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}

	for _, outline := range opmlFile.Outlines {
		url := outline.URL
		if !store.Contains(url) {
			store.Add(url)
		}
	}

	if err = store.Save(); err != nil {
		panic(err)
	}
}
