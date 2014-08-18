package main

import (
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
	appVersion = "1.6dev"
)

var (
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
	cast string
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
	if err != nil && !os.IsNotExist(err) {
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
	for e := range episodes {
		wg.Add(1)
		go func(item episode) {
			defer wg.Done()
			if err := getEpisode(item); err != nil {
				logger.Println(err)
			}
		}(e)
	}

	wg.Wait()
}

func newEpisodes(podcasts <-chan feed.Feed) <-chan episode {
	out := make(chan episode)
	go func() {
		for p := range podcasts {
			name := util.Escape(p.Title)
			if len(name) <= 0 {
				logger.Println("Skipping %q, couldn't escape name", p.Title)
				continue
			}

			unread, err := unreadMarker(name, p)
			if err != nil {
				logger.Println(err)
				continue
			}

			items := p.Items
			if *recent > 0 && len(items) >= *recent {
				items = items[0:*recent]
			}

			for _, i := range items {
				if len(i.Attachment) > 0 && i.Date.After(unread) {
					out <- episode{i, name}
				}
			}
		}

		close(out)
	}()

	return out
}

func getEpisode(e episode) error {
	var path string
	var err error

	for i := 1; i <= *retry; i++ {
		path, err = util.Get(e.item.Attachment, filepath.Join(downloadDir, e.cast))
		if err == nil {
			break
		}
	}

	// Last error returned by the loop above
	if err != nil {
		return err
	}

	name := util.Escape(e.item.Title)
	if len(name) > 0 {
		// Fallback to provided file name if rename fails
		os.Rename(path, filepath.Join(filepath.Dir(path), name+filepath.Ext(path)))
	}

	return nil
}

func unreadMarker(name string, cast feed.Feed) (marker time.Time, err error) {
	path := filepath.Join(downloadDir, name)
	if err = os.MkdirAll(path, 0755); err != nil {
		return
	}

	file, err := os.OpenFile(filepath.Join(path, ".latest"), os.O_RDWR+os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer file.Close()

	var timestamp int64
	fmt.Fscanf(file, "%d\n", &timestamp) // XXX
	marker = time.Unix(timestamp, 0)
	latest := marker

	for _, i := range cast.Items {
		if i.Date.After(latest) {
			latest = i.Date
		}
	}

	_, err = fmt.Fprintf(file, "%d\n", latest.Unix())
	return
}
