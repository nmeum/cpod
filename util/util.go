package util

import (
	"io"
	"net/http"
	"os"
	"os/signal"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"unicode"
)

func Get(url, target string) (fp string, err error) {
	fp = filepath.Join(target, strings.TrimSpace(path.Base(url)))
	if err = os.MkdirAll(filepath.Dir(fp), 0755); err != nil {
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	file, err := os.Create(fp)
	if err != nil {
		return
	}

	defer file.Close()
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
		case unicode.IsLetter(r):
			return r
		case unicode.IsNumber(r):
			return r
		case unicode.IsSpace(r):
			return '-'
		case unicode.IsPunct(r):
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

func EnvDefault(key, fallback string) string {
	dir := os.Getenv(key)
	if len(dir) > 0 {
		return dir
	}

	var home string
	user, err := user.Current()
	if err == nil && len(user.HomeDir) > 0 {
		home = user.HomeDir
	} else {
		home = os.Getenv("HOME")
	}

	return filepath.Join(home, fallback)
}
