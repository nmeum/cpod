package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func download(url, target, name string) (err error) {
	if err = os.MkdirAll(target, 0755); err != nil && !os.IsExist(err) {
		return
	}

	path := filepath.Join(target, name+filepath.Ext(url))
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

func isPodcast(url string) (b bool) {
	for _, cast := range storage.Podcasts {
		if cast.URL == url {
			return true
		}
	}

	return
}

func envDefault(key, fallback string) (d string) {
	d = os.Getenv(key)
	if len(d) <= 0 {
		d = filepath.Join(os.Getenv("HOME"), fallback)
	}

	return

}

func abort(err error) {
	logger.Fatalln("ERROR:", err.Error())
}
