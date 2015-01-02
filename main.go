package main

import (
	"errors"
	"strings"
	"path"
	"flag"
	"fmt"
	"github.com/nmeum/cpod/feed"
	"github.com/nmeum/cpod/store"
	"github.com/nmeum/cpod/util"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	appName    = "cpod"
	appVersion = "1.8dev"
)

var (
	limit   = flag.Int("p", 5, "number of maximal parallel downloads")
	retry   = flag.Int("t", 3, "number of times a failed download is retried")
	recent  = flag.Int("r", 0, "number of most recent episodes to download")
	version = flag.Bool("v", false, "display version number and exit")
)

var (
	logger      = log.New(os.Stderr, fmt.Sprintf("%s: ", appName), 0)
	downloadDir = util.EnvDefault("CPOD_DOWNLOAD_DIR", "podcasts")
)

type episode struct {
	item feed.Item
	cast feed.Feed
}

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
	if err != nil {
		logger.Fatal(err)
	}

	lockPath := filepath.Join(cacheDir, "lock")
	if err = util.Lock(lockPath); os.IsExist(err) {
		logger.Fatalf("database is locked, remove %q to force unlock\n", lockPath)
	} else if err != nil {
		logger.Fatal(err)
	}

	update(storage)
	os.Remove(lockPath)
}

func update(storage *store.Store) {
	podcasts := storage.Fetch()
	episodes := newEpisodes(podcasts)

	var wg sync.WaitGroup
	var counter int

	for e := range episodes {
		wg.Add(1)
		counter++

		go func(item episode, count int) {
			if err := getEpisode(item); err != nil {
				logger.Println(err)
			}

			wg.Done()
			count--
		}(e, counter)

		for *limit > 0 && counter >= *limit {
			time.Sleep(3 * time.Second)
		}
	}

	wg.Wait()
}

func newEpisodes(podcasts <-chan feed.Feed) <-chan episode {
	out := make(chan episode)
	go func(pcasts <-chan feed.Feed) {
		for p := range pcasts {
			unread, err := readMarker(p.Title)
			if err != nil && !os.IsNotExist(err) {
				logger.Println(err)
				continue
			}

			items := p.Items
			if *recent > 0 && len(items) >= *recent {
				items = items[0:*recent]
			}

			for _, i := range items {
				if len(i.Attachment) <= 0 || i.Date.Before(unread) {
					break
				}

				out <- episode{i, p}
			}

			// TODO only do this if the download was succesfull
			if err := writeMarker(p.Title, items[0].Date); err != nil {
				logger.Println(err)
			}
		}

		close(out)
	}(podcasts)

	return out
}

func getEpisode(e episode) error {
	cast := util.Escape(e.cast.Title)
	if len(cast) <= 0 {
		return errors.New(fmt.Sprintf("couldn't escape title %q", e.cast.Title))
	}

	url := strings.TrimSpace(e.item.Attachment)
	fp := filepath.Join(downloadDir, cast, path.Base(url))
	if err := os.MkdirAll(filepath.Dir(fp), 0755); err != nil {
		return err
	}

	if err := util.Get(url, fp, *retry); err != nil {
		return err
	}

	name := util.Escape(e.item.Title)
	if len(name) > 0 {
		os.Rename(fp, filepath.Join(filepath.Dir(fp), name+filepath.Ext(fp)))
	}

	return nil
}

func readMarker(name string) (marker time.Time, err error) {
	file, err := os.Open(filepath.Join(downloadDir, name, ".latest"))
	if err != nil {
		return
	}

	defer file.Close()
	var timestamp int64

	if _, err = fmt.Fscanf(file, "%d\n", &timestamp); err != nil {
		return
	}

	marker = time.Unix(timestamp, 0)
	return
}

func writeMarker(name string, latest time.Time) error {
	path := filepath.Join(downloadDir, name, ".latest")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()
	if _, err := fmt.Fprintf(file, "%d\n", latest.Unix()); err != nil {
		return err
	}

	return nil
}
