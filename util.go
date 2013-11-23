package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func parseDate(date string) (t time.Time, err error) {
	formats := []string{
		time.RFC1123Z, time.RFC1123, time.RFC822Z,
		time.RFC822, time.ANSIC, time.RFC3339,
		time.RFC850, time.RubyDate, time.UnixDate,
	}

	for _, format := range formats {
		t, err = time.Parse(format, date)
		if err == nil {
			return
		} else {
			err = nil
		}
	}

	return
}

func download(url string, target string, name string) (err error) {
	if err = os.MkdirAll(target, 0755); err != nil && !os.IsExist(err) {
		return
	}

	path := filepath.Join(target, name, filepath.Ext(url))
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
