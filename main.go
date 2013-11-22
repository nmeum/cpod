package main

import (
	"flag"
	"github.com/nmeum/cpod/opml"
	"github.com/nmeum/cpod/rss"
	"github.com/nmeum/cpod/store"
	"os"
	"path/filepath"
)

var (
	update     = flag.Bool("u", false, "update all feeds")
	noDownload = flag.Bool("d", false, "don't download new episodes during update")
	opmlImport = flag.String("i", "", "import opml file")
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
		panic(err)
	}

	flag.Parse()
	if err := processFlags(); err != nil {
		panic(err)
	}

	if err := storage.Save(); err != nil {
		panic(err)
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
		d = filepath.Join(os.Getenv("HOME"), "Podcasts")
	}

	return
}

func processFlags() (err error) {
	if len(*opmlImport) > 0 {
		if err = importCmd(*opmlImport); err != nil {
			return
		}
	}

	if *update {
		if err = updateCmd(); err != nil {
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
		for _, i := range xml.Channel.Items {
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
				if err = download(i.Enclosure.Url, dir); err != nil {
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

	for _, o := range file.Body.Outline {
		storage.Add(o.Text, o.XmlUrl)
	}

	return
}
