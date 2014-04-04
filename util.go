package main

import (
	"fmt"
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

func cleanupDir(path string) (err error) {
	dir, err := os.Open(path)
	if err != nil {
		return
	}

	defer dir.Close()
	files, err := dir.Readdir(-1)
	if err != nil {
		return
	}

	if len(files) <= 1 {
		return
	}

	latest := latestFile(files)
	for _, file := range files {
		if file.Name() == latest.Name() {
			continue
		}

		path := filepath.Join(path, file.Name())
		if err = os.Remove(path); err != nil {
			return
		}
	}

	return
}

func latestFile(files []os.FileInfo) (f os.FileInfo) {
	f = files[0]
	for _, file := range files {
		if file.ModTime().After(f.ModTime()) {
			f = file
		}
	}

	return
}

func isPodcast(title string) (b bool) {
	for _, cast := range storage.Podcasts {
		if cast.Title == title {
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

	return d
}

func abort(err error) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
	os.Exit(2)
}
