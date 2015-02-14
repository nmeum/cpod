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
		var op *opml.Opml

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
	if len(os.Args) <= 0 {
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
