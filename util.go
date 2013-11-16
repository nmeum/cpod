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

func download(url string, dir string) (err error) {
	if err = os.MkdirAll(dir, 0755); err != nil && !os.IsExist(err) {
		return
	}

	file, err := os.Create(filepath.Join(dir, filepath.Base(url)))
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
