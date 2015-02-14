package main

import (
	"flag"
	"fmt"
	"github.com/nmeum/cpod/store"
	"github.com/nmeum/cpod/util"
	"github.com/nmeum/go-feedparser"
	"log"
	"os"
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

	storeDir := filepath.Join(util.EnvDefault("XDG_CONFIG_HOME", ".config"), appName)
	lockPath := filepath.Join(os.TempDir(), fmt.Sprintf("%s-%s", appName, util.Username()))

	storage, err := store.Load(filepath.Join(storeDir, "urls"))
	if err != nil {
		logger.Fatal(err)
	}

	if err := util.Lock(lockPath); os.IsExist(err) {
		logger.Fatalf("database is locked, remove %q to force unlock\n", lockPath)
	} else if err != nil {
		logger.Fatal(err)
	}

	update(storage)
	if err := os.Remove(lockPath); err != nil {
		logger.Fatal(err)
	}
}

func update(storage *store.Store) {
	var wg sync.WaitGroup
	var counter int

	for cast := range storage.Fetch() {
		wg.Add(1)
		counter++

		go func(p store.Podcast) {
			defer func() {
				wg.Done()
				counter--
			}()

			feed := p.Feed
			if p.Error != nil {
				logger.Println(p.Error)
				return
			}

			items, err := newItems(feed)
			if err != nil {
				logger.Println(err)
				return
			}

			for _, i := range items {
				if err := getItem(feed, i); err != nil {
					logger.Println(err)
					continue
				}

				if err := writeMarker(feed.Title, i.Date); err != nil {
					logger.Println(err)
					continue
				}
			}

			return
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

	target := filepath.Join(downloadDir, title)
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return err
	}

	if err := util.GetFile(url, target); err != nil {
		return err
	}

	name, err := util.Escape(item.Title)
	if err == nil {
		fn, err := util.Filename(url)
		if err != nil {
			return err
		}

		newPath := filepath.Join(target, name+filepath.Ext(fn))
		if err = os.Rename(filepath.Join(target, fn), newPath); err != nil {
			return err
		}
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
