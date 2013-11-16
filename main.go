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
	update      = flag.Bool("u", false, "update all feeds")
	opmlImport  = flag.String("i", "", "import opml file")
	downloadNew = flag.Bool("d", false, "download new episodes")
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

	if *downloadNew {
		if err = downloadCmd(); err != nil {
			return
		}
	}

	return
}

func updateCmd() (err error) {
	for _, f := range storage.Feeds {
		xml, err := rss.Parse(f.Url)
		if err != nil {
			return err
		}

		date, err := parseDate(xml.Channel.Items[0].PubDate)
		if err != nil {
			return err
		}

		f.Latest = date.Unix()
	}

	return
}

func downloadCmd() error {
	for _, f := range storage.Feeds {
		xml, err := rss.Parse(f.Url)
		if err != nil {
			return err
		}

		for _, i := range xml.Channel.Items {
			if len(i.Enclosure.Url) <= 0 {
				continue
			}

			date, err := parseDate(i.PubDate)
			if err != nil {
				return err
			}

			if f.Latest < date.Unix() {
				dir := filepath.Join(downloadDir, f.Title)
				if err = download(i.Enclosure.Url, dir); err != nil {
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

	for _, o := range file.Body.Outline {
		storage.Add(o.Text, o.XmlUrl)
	}

	return
}
