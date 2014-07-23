package main

import (
	"flag"
	"fmt"
	"github.com/nmeum/cpod/store"
	"github.com/nmeum/cpod/util"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	appName    = "cpod"
	appVersion = "1.6dev"
)

var (
	recent     = flag.Int("r", 0, "only download latest n episodes")
	version    = flag.Bool("v", false, "print version and exit")
	noUpdate   = flag.Bool("u", false, "don't update feeds and don't download new episodes")
	noDownload = flag.Bool("d", false, "don't download new episodes and skip them")
)

var (
	logger      = log.New(os.Stderr, appName+": ", 0)
	downloadDir = util.EnvDefault("CPOD_DOWNLOAD_DIR", "podcasts")
)

func main() {
	flag.Parse()
	if *version {
		logger.Fatal(appVersion)
	}

	cacheDir := filepath.Join(util.EnvDefault("XDG_CACHE_HOME", ".cache"), appName)
	storeDir := filepath.Join(util.EnvDefault("XDG_CONFIG_HOME", ".config"), appName)

	for _, dir := range []string{cacheDir, storeDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			logger.Fatal(err)
		}
	}

	storage, err := store.Load(filepath.Join(storeDir, "urls"))
	if err != nil && !os.IsNotExist(err) {
		logger.Fatal(err)
	}

	lockPath := filepath.Join(cacheDir, "lock")
	if err = util.Lock(lockPath); os.IsExist(err) {
		logger.Fatalf("database is locked, remove %q to force unlock\n", lockPath)
	} else if err != nil {
		logger.Fatal(err)
	}

	err = updateFeeds(storage)
	os.Remove(lockPath)
	if err != nil {
		logger.Fatal(err)
	}
}

func updateFeeds(storage *store.Store) error {
	feeds := storage.Fetch()
	for f := range feeds {
		name := util.Escape(f.Title)
		if len(name) <= 0 {
			name = f.Title
		}

		path := filepath.Join(downloadDir, name)
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}

		file, err := os.OpenFile(filepath.Join(path, ".latest"), os.O_RDWR+os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		defer file.Close()

		var timestamp int64
		fmt.Fscanf(file, "%d\n", &timestamp) /// XXX

		latest := time.Unix(timestamp, 0)
		for _, i := range f.Items {
			if i.Date.After(latest) {
				latest = i.Date
			}
		}

		if _, err := fmt.Fprintf(file, "%d\n", latest.Unix()); err != nil {
			return err
		}

		// TODO 1. Create the podcast dir and the .latest file for the podcast
		// TODO 2. Use the .latest file to find out which episodes are new
		// TODO 3. Download the new episodes
	}

	return nil
}
