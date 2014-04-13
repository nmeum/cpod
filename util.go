package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func download(url, target, name string) (err error) {
	if err = os.MkdirAll(target, 0755); err != nil && !os.IsExist(err) {
		return
	}

	filename := escape(name) + filepath.Ext(url)
	file, err := os.Create(filepath.Join(target, filename))
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
		case r == '.' || r == '_':
			return r
		case r == ' ' || r == '\t':
			return '-'
		}

		return -1
	}

	escaped := strings.Map(mfunc, name)
	for strings.HasPrefix(escaped, "-") && len(escaped) > 1 {
		escaped = escaped[1:]
	}

	return escaped
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
	logger.Fatalf("%s: %s", appName, err.Error())
}
