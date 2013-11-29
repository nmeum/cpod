package main

import (
	"errors"
	"flag"
	"github.com/nmeum/cpod/opml"
	"github.com/nmeum/cpod/rss"
	"github.com/nmeum/cpod/store"
	"os"
	"path/filepath"
)

var (
	recent     = flag.Int("r", 0, "download latest n episodes")
	cleanup    = flag.Bool("c", false, "remove old episodes")
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
		handleError(err)
	}

	flag.Parse()
	if err := processInput(); err != nil {
		handleError(err)
	}

	if err := storage.Save(); err != nil {
		handleError(err)
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
		xml, err := rss.Parse(f.Url)
		if err != nil {
			return err
		}

		var latest int64
		items := xml.Channel.Items

		if *recent > 0 {
			items = items[0:*recent]
		}

		for _, i := range items {
			if len(i.Enclosure.Url) <= 0 {
				continue
			}

			date, err := parseDate(i.PubDate)
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
				if err = download(i.Enclosure.Url, dir, i.Title); err != nil {
					return err
				}
			}
		}

		storage.Feeds[n].Latest = latest
	}

	return nil
}

func importCmd(path string) (err error) {
	file, err := opml.Load(path)
	if err != nil {
		return
	}

	for _, o := range file.Body.Outlines {
		storage.Add(o.Text, o.XmlUrl)
	}

	return
}

func exportCmd(path string) (err error) {
	return errors.New("Not implemented yet!")
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
