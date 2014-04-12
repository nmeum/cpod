package main

import (
	"flag"
	"github.com/nmeum/cpod/feed"
	"github.com/nmeum/cpod/opml"
	"github.com/nmeum/cpod/store"
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
	logger      = log.New(os.Stderr, "", 0)
	downloadDir = envDefault("CPOD_DOWNLOAD_DIR", "podcasts")
)

func main() {
	storeDir := filepath.Join(envDefault("XDG_CONFIG_HOME", ".config"), appName)

	err := os.MkdirAll(storeDir, 0755)
	if err != nil && !os.IsExist(err) {
		abort(err)
	}

	storage, err = store.Load(filepath.Join(storeDir, "feeds.json"))
	if err != nil && !os.IsNotExist(err) {
		abort(err)
	}

	flag.Parse()
	if err = processInput(); err != nil {
		abort(err)
	}

	if err = storage.Save(); err != nil {
		abort(err)
	}
}

func processInput() (err error) {
	if *version {
		logger.Fatalln(appName, appVersion)
	}

	if len(*opmlImport) > 0 {
		if err = importCmd(*opmlImport); err != nil {
			return
		}
	}

	if len(*opmlExport) > 0 {
		if err = exportCmd(*opmlExport); err != nil {
			return
		}
	}

	if !*noUpdate {
		if err = updateCmd(); err != nil {
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
				err := download(item.Attachment, filepath.Join(downloadDir, p.Title), item.Title)
				if err != nil {
					return err
				}
			}

		}
	}

	return nil
}

func importCmd(path string) (err error) {
	file, err := opml.Load(path)
	if err != nil {
		return
	}

	for _, o := range file.Outlines {
		if !isPodcast(o.XMLURL) {
			storage.Add(o.Text, o.Type, o.XMLURL)
		}
	}

	return
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
