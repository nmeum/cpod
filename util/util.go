package util

import (
	"github.com/nmeum/cpod/store"
	"io"
	"net/http"
	"os"
	"os/signal"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

func Get(url, target string) (fp string, err error) {
	if err = os.MkdirAll(target, 0755); err != nil {
		return
	}

	fp = filepath.Join(target, strings.TrimSpace(path.Base(url)))
	file, err := os.Create(fp)
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

func Lock(path string) (err error) {
	_, err = os.OpenFile(path, os.O_CREATE+os.O_EXCL, 0666)
	if err != nil {
		return
	}

	// Setup unlock handler
	ch := make(chan os.Signal, 1)
	go func() {
		<-ch // Block until signal is received
		os.Remove(path)
		os.Exit(2)
	}()
	signal.Notify(ch, os.Interrupt, os.Kill)

	return
}

func Escape(name string) string {
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

func Subscribed(casts []*store.Podcast, url string) bool {
	for _, cast := range casts {
		if cast.URL == url {
			return true
		}
	}

	return false
}

func EnvDefault(key, fallback string) string {
	dir := os.Getenv(key)
	if len(dir) <= 0 {
		dir = filepath.Join(home(), fallback)
	}

	return dir
}

func home() string {
	user, err := user.Current()
	if err == nil && len(user.HomeDir) > 0 {
		return user.HomeDir
	}

	return os.Getenv("HOME")
}
