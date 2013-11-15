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
	update      = flag.Bool("u", false, "update all feeds")
	opmlImport  = flag.String("i", "", "import opml file")
	downloadNew = flag.Bool("d", false, "download new episodes")
)

var (
	storage     *store.Store
	downloadDir string
)

func main() {
	storeDir := os.Getenv("XDG_DATA_HOME")
	if len(storeDir) <= 0 {
		storeDir = filepath.Join(os.Getenv("HOME"), ".local", "share")
	}

	downloadDir = os.Getenv("CPOD_DOWNLOAD_DIR")
	if len(downloadDir) <= 0 {
		downloadDir = filepath.Join(os.Getenv("HOME"), "Podcasts")
	}

	storeDir = filepath.Join(storeDir, "cpod")
	if err := os.Mkdir(storeDir, 0755); err != nil && !os.IsExist(err) {
		panic(err)
	}

	var err error
	storage, err = store.Load(filepath.Join(storeDir, "feeds.json"))
	if err != nil {
		panic(err)
	}

	defer storage.Close()
	flag.Parse()

	if err := processFlags(); err != nil {
		panic(err)
	}

	storage.Save()
}

func processFlags() (err error) {
	if *update {
		return updateCmd()
	}

	if *downloadNew {
		return downloadCmd()
	}

	if len(*opmlImport) > 0 {
		return importCmd(*opmlImport)
	}

	return errors.New("No operation specified")
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
