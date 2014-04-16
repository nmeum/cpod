package main

import (
	"flag"
	"github.com/nmeum/cpod/feed"
	"github.com/nmeum/cpod/opml"
	"github.com/nmeum/cpod/store"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	downloadDir = envDefault("CPOD_DOWNLOAD_DIR", "podcasts")
)

func main() {
	storeDir := filepath.Join(envDefault("XDG_CONFIG_HOME", ".config"), appName)

	err := os.MkdirAll(storeDir, 0755)
	if err != nil && !os.IsExist(err) {
		logger.Fatal(err)
	}

	storage, err = store.Load(filepath.Join(storeDir, "feeds.json"))
	if err != nil && !os.IsNotExist(err) {
		logger.Fatal(err)
	}

	flag.Parse()
	if err = processInput(); err != nil {
		logger.Fatal(err)
	}

	if err = storage.Save(); err != nil {
		logger.Fatal(err)
	}
}

func processInput() (err error) {
	if *version {
		logger.Fatal(appVersion)
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
				path, err := download(item.Attachment, filepath.Join(downloadDir, p.Title))
				if err != nil {
					return err
				}

				name := escape(item.Title)
				if len(name) > 1 {
					err = os.Rename(path, filepath.Join(filepath.Dir(path), name+filepath.Ext(path)))
					if err != nil {
						return err
					}
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

func download(url, target string) (path string, err error) {
	if err = os.MkdirAll(target, 0755); err != nil && !os.IsExist(err) {
		return
	}

	path = filepath.Join(target, strings.TrimSpace(filepath.Base(url)))
	file, err := os.Create(path)
	if err != nil {
		return
	}

	defer file.Close()
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	if _, err = io.Copy(file, resp.Body); err != nil {
		return
	}

	return
}

func escape(name string) string {
	mfunc := func(r rune) rune {
		switch {
		case r >= '0' && r <= '9':
			return r
		case r >= 'A' && r <= 'Z':
			return r
		case r >= 'a' && r <= 'z':
			return r
		case r == '.' || r == ':':
			return '-'
		case r == ' ' || r == '_':
			return '-'
		}

		return -1
	}

	escaped := strings.Map(mfunc, name)
	for strings.Contains(escaped, "--") {
		escaped = strings.Replace(escaped, "--", "-", -1)
	}

	escaped = strings.TrimPrefix(escaped, "-")
	escaped = strings.TrimSuffix(escaped, "-")

	return escaped
}

func isPodcast(url string) bool {
	for _, cast := range storage.Podcasts {
		if cast.URL == url {
			return true
		}
	}

	return false
}

func envDefault(key, fallback string) string {
	dir := os.Getenv(key)
	if len(dir) <= 0 {
		dir = filepath.Join(os.Getenv("HOME"), fallback)
	}

	return dir
}
