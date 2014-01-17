package main

import (
	"flag"
	"fmt"
	"github.com/nmeum/cpod/feed"
	"github.com/nmeum/cpod/opml"
	"github.com/nmeum/cpod/store"
	"os"
	"path/filepath"
)

const (
	appName    = "cpod"
	appVersion = "0.0"
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
	var storeDir string

	storeDir, downloadDir = getDirs()
	if err := os.Mkdir(storeDir, 0755); err != nil && !os.IsExist(err) {
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

func getDirs() (s string, d string) {
	s = os.Getenv("XDG_DATA_HOME")
	if len(s) <= 0 {
		s = filepath.Join(os.Getenv("HOME"), ".local", "share")
	}

	s = filepath.Join(s, "cpod")
	d = os.Getenv("CPOD_DOWNLOAD_DIR")
	if len(d) <= 0 {
		d = filepath.Join(os.Getenv("HOME"), "podcasts")
	}

	return
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
	for n, f := range storage.Feeds {
		xml, err := feed.Parse(f.Url)
		if err != nil {
			return err
		}

		var latest int64
		items := xml.Items

		if *recent > 0 {
			items = items[0:*recent]
		}

		for _, i := range items {
			if len(i.Attachment) <= 0 {
				continue
			}

			date, err := parseDate(i.Date)
			if err != nil {
				return err
			}

			if latest == 0 {
				latest = date.Unix()
				if *noDownload {
					break
				}
			}

			if f.Latest < date.Unix() && !*noDownload {
				dir := filepath.Join(downloadDir, f.Title)
				if err = download(i.Attachment, dir, i.Title); err != nil {
					return err
				}
			}
		}

		storage.Feeds[n].Type = xml.Type
		storage.Feeds[n].Latest = latest
	}

	return nil
}

func importCmd(path string) (err error) {
	file, err := opml.Load(path)
	if err != nil {
		return
	}

	for _, o := range file.Outlines {
		storage.Add(o.Text, "", o.XmlUrl)
	}

	return
}

func exportCmd(path string) (err error) {
	export := opml.Create("Podcast subscriptions")
	for _, feed := range storage.Feeds {
		export.Add(feed.Title, feed.Type, feed.Url)
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
		if !file.IsDir() {
			continue
		}

		path := filepath.Join(downloadDir, file.Name())
		if err = cleanupDir(path); err != nil {
			return
		}
	}

	return
}
