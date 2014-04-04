package main

import (
	"flag"
	"fmt"
	"github.com/nmeum/cpod/feed"
	"github.com/nmeum/cpod/opml"
	"github.com/nmeum/cpod/store"
	"os"
	"path/filepath"
	"time"
)

const (
	appName    = "cpod"
	appVersion = "1.1"
)

var (
	recent     = flag.Int("r", 0, "download latest n episodes")
	cleanup    = flag.Bool("c", false, "remove old episodes")
	version    = flag.Bool("v", false, "print version and exit")
	noUpdate   = flag.Bool("u", false, "don't update feeds")
	noDownload = flag.Bool("d", false, "don't download new episodes")
	opmlImport = flag.String("i", "", "import opml file")
	opmlExport = flag.String("e", "", "export opml file")
)

var (
	storage     *store.Store
	downloadDir string
)

func main() {
	var err error
	downloadDir = envDefault("CPOD_DOWNLOAD_DIR", "podcasts")

	storeDir := filepath.Join(envDefault("XDG_DATA_HOME", ".local/share"), appName)
	if err = os.Mkdir(storeDir, 0755); err != nil && !os.IsExist(err) {
		return
	}

	storage, err = store.Load(filepath.Join(storeDir, "feeds.json"))
	if err != nil && !os.IsNotExist(err) {
		abort(err)
	}

	flag.Parse()
	if err := processInput(); err != nil {
		abort(err)
	}

	if err := storage.Save(); err != nil {
		abort(err)
	}
}

func processInput() (err error) {
	if *version {
		fmt.Fprintf(os.Stderr, "%s %s\n", appName, appVersion)
		os.Exit(2)
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

	if *cleanup {
		if err = cleanupCmd(); err != nil {
			return
		}
	}

	return
}

func updateCmd() error {
	for _, p := range storage.Podcasts {
		feed, err := feed.Parse(p.Url)
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
		if !isPodcast(o.Text) {
			storage.Add(o.Text, o.Type, o.XmlUrl)
		}
	}

	return
}

func exportCmd(path string) (err error) {
	export := opml.Create("Podcast subscriptions")
	for _, cast := range storage.Podcasts {
		export.Add(cast.Title, cast.Type, cast.Url)
	}

	if err = export.Save(*opmlExport); err != nil {
		return
	}

	return
}

func cleanupCmd() (err error) {
	dir, err := os.Open(downloadDir)
	if err != nil {
		return
	}

	defer dir.Close()
	files, err := dir.Readdir(-1)
	if err != nil {
		return
	}

	for _, file := range files {
		if !file.IsDir() || !isPodcast(file.Name()) {
			continue
		}

		path := filepath.Join(downloadDir, file.Name())
		if err = cleanupDir(path); err != nil {
			return
		}
	}

	return
}
