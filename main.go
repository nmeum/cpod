package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/nmeum/freddie"
	"github.com/nmeum/freddie/feed"
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
	downloadDir = envDefault("CPOD_DOWNLOAD_DIR", "podcasts")
)

func main() {
	flag.Parse()
	if *version {
		logger.Fatal(appVersion)
	}

	cacheDir := filepath.Join(envDefault("XDG_CACHE_HOME", ".cache"), appName)
	storeDir := filepath.Join(envDefault("XDG_CONFIG_HOME", ".config"), appName)

	for _, dir := range []string{cacheDir, storeDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			logger.Fatal(err)
		}
	}

	lockPath := filepath.Join(cacheDir, "lock")
	if err := lock(lockPath); os.IsExist(err) {
		logger.Fatalf("database is locked, remove %q to force unlock\n", lockPath)
	} else if err != nil {
		logger.Fatal(err)
	}

	update(readFeeds(filepath.Join(storeDir, "urls")))
	os.Remove(lockPath)
}

func update(podcasts <-chan feed.Feed) {
	var wg sync.WaitGroup
	var counter int

	for cast := range podcasts {
		wg.Add(1)
		counter++

		go func(p feed.Feed) {
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

func readFeeds(fp string) <-chan feed.Feed {
	file, err := os.Open(fp)
	if err != nil {
		logger.Fatal(err)
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	var feeds []string
	for scanner.Scan() {
		feeds = append(feeds, scanner.Text())
	}

	out := make(chan feed.Feed)
	go func(c chan feed.Feed, urls []string) {
		for _, url := range urls {
			feed, err := freddie.Parse(url)
			if err != nil {
				logger.Println(err)
			} else {
				c <- feed
			}
		}

		close(out)
	}(out, feeds)

	return out
}

func newItems(cast feed.Feed) (items []feed.Item, err error) {
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
		if len(item.Attachment) <= 0 || item.Date.Before(unread) {
			break
		}

		items = append(items, item)
	}

	return
}

func getItem(cast feed.Feed, item feed.Item) error {
	url := strings.TrimSpace(item.Attachment)

	fp := filepath.Join(downloadDir, cast.Title, path.Base(url))
	if err := os.MkdirAll(filepath.Dir(fp), 0755); err != nil {
		return err
	}

	if err := get(url, fp, *retry); err != nil {
		return err
	}

	name, err := escape(item.Title)
	if err == nil {
		os.Rename(fp, filepath.Join(filepath.Dir(fp), name+filepath.Ext(fp)))
	}

	return nil
}

func readMarker(name string) (marker time.Time, err error) {
	name, err = escape(name)
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
	name, err := escape(name)
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
