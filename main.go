package main

import (
	"flag"
	"fmt"
	"github.com/nmeum/cpod/store"
	"github.com/nmeum/cpod/util"
	"github.com/nmeum/go-feedparser"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
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
	if err := util.Lock(lockPath); os.IsExist(err) {
		logger.Fatalf("database is locked, remove %q to force unlock\n", lockPath)
	} else if err != nil {
		logger.Fatal(err)
	}

	update(storage)
	os.Remove(lockPath)
}

func update(storage *store.Store) {
	var wg sync.WaitGroup
	var counter int

	for cast := range storage.Fetch() {
		wg.Add(1)
		counter++

		go func(p feedparser.Feed) {
			items, err := newItems(p)
			if err != nil {
				logger.Println(err)
				return
			}

			for _, i := range items {
				if err := getItem(p, i); err != nil {
					logger.Println(err)
					continue
				}

				if err := writeMarker(p.Title, i.Date); err != nil {
					logger.Println(err)
				}
			}

			wg.Done()
			counter--
		}(cast)

		for *limit > 0 && counter >= *limit {
			time.Sleep(3 * time.Second)
		}
	}

	wg.Wait()
}

func newItems(cast feedparser.Feed) (items []feedparser.Item, err error) {
	unread, err := readMarker(cast.Title)
	if os.IsNotExist(err) {
		err = nil
	} else if err != nil {
		return
	}

	if *recent > 0 && len(cast.Items) >= *recent {
		cast.Items = cast.Items[0:*recent]
	}

	for _, item := range cast.Items {
		if len(item.Attachment) <= 0 || item.Date.Before(unread) ||
			item.Date.Equal(unread) {
			break
		}

		items = append(items, item)
	}

	return
}

func getItem(cast feedparser.Feed, item feedparser.Item) error {
	title, err := util.Escape(cast.Title)
	if err != nil {
		return err
	}

	url := strings.TrimSpace(item.Attachment)
	fp := filepath.Join(downloadDir, title, path.Base(url))
	if err := os.MkdirAll(filepath.Dir(fp), 0755); err != nil {
		return err
	}

	if err := util.Get(url, fp, *retry); err != nil {
		return err
	}

	name, err := util.Escape(item.Title)
	if err == nil {
		os.Rename(fp, filepath.Join(filepath.Dir(fp), name+filepath.Ext(fp)))
	}

	return nil
}

func readMarker(name string) (marker time.Time, err error) {
	name, err = util.Escape(name)
	if err != nil {
		return
	}

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
	name, err := util.Escape(name)
	if err != nil {
		return err
	}

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
