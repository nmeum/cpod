package main

import (
	"flag"
	"github.com/nmeum/cpod/feed"
	"github.com/nmeum/cpod/opml"
	"github.com/nmeum/cpod/store"
	"github.com/nmeum/cpod/util"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	appName    = "cpod"
	appVersion = "1.1"
)

var (
	recent     = flag.Int("r", 0, "only download latest n episodes")
	version    = flag.Bool("v", false, "print version and exit")
	noUpdate   = flag.Bool("u", false, "don't update feeds and don't download new episodes")
	noDownload = flag.Bool("d", false, "don't download new episodes and skip them")
	opmlImport = flag.String("i", "", "import opml file at path")
	opmlExport = flag.String("e", "", "export opml file to path")
)

var (
	storage     *store.Store
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

	var err error
	for _, dir := range []string{cacheDir, storeDir} {
		if err = os.MkdirAll(dir, 0755); err != nil {
			logger.Fatal(err)
		}
	}

	storage, err = store.Load(filepath.Join(storeDir, "feeds.json"))
	if err != nil && !os.IsNotExist(err) {
		logger.Fatal(err)
	}

	lockPath := filepath.Join(cacheDir, "lock")
	if err = util.Lock(lockPath); err != nil && os.IsExist(err) {
		logger.Fatalf("database is locked, remove %q to force unlock\n", lockPath)
	} else if err != nil {
		logger.Fatal(err)
	}

	err = processInput()
	os.Remove(lockPath)
	if err != nil {
		logger.Fatal(err)
	}
}

func processInput() (err error) {
	if len(*opmlImport) > 0 {
		if err = importCmd(*opmlImport); err != nil {
			return
		}
	}

	if !*noUpdate {
		if err = updateCmd(); err != nil {
			return
		}
	}

	if len(*opmlExport) > 0 {
		if err = exportCmd(*opmlExport); err != nil {
			return
		}
	}

	return
}

func updateCmd() error {
	for _, p := range storage.Podcasts {
		feed, err := feed.Parse(p.URL)
		if err != nil {
			return err
		}

		p.Type = feed.Type
		items := feed.Items

		if *recent > 0 {
			items = items[0:*recent]
		}

		latest := time.Unix(p.Latest, 0)
		for _, item := range items {
			if item.Date.After(time.Unix(p.Latest, 0)) {
				p.Latest = item.Date.Unix()
			}

			if item.Date.After(latest) && len(item.Attachment) > 0 && !*noDownload {
				path, err := util.Get(item.Attachment, filepath.Join(downloadDir, p.Title))
				if err != nil {
					return err
				}

				name := util.Escape(item.Title)
				if len(name) > 1 {
					err = os.Rename(path, filepath.Join(filepath.Dir(path), name+filepath.Ext(path)))
					if err != nil {
						return err
					}
				}
			}

		}
	}

	return storage.Save()
}

func importCmd(path string) (err error) {
	file, err := opml.Load(path)
	if err != nil {
		return
	}

	for _, o := range file.Outlines {
		if !util.Subscribed(storage.Podcasts, o.URL) {
			storage.Add(o.Text, o.Type, o.URL)
		}
	}

	return storage.Save()
}

func exportCmd(path string) (err error) {
	export := opml.Create("Podcast subscriptions")
	for _, cast := range storage.Podcasts {
		export.Add(cast.Title, cast.Type, cast.URL)
	}

	if err = export.Save(*opmlExport); err != nil {
		return
	}

	return
}
