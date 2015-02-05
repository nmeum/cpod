package util

import (
	"errors"
	"html"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"strings"
	"unicode"
)

func Lock(path string) (err error) {
	_, err = os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0644)
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

func Escape(name string) (escaped string, err error) {
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

	escaped = strings.Map(mfunc, html.UnescapeString(name))
	for strings.Contains(escaped, "--") {
		escaped = strings.Replace(escaped, "--", "-", -1)
	}

	escaped = strings.TrimPrefix(escaped, "-")
	escaped = strings.TrimSuffix(escaped, "-")

	if len(escaped) <= 0 {
		err = errors.New("couldn't escape title")
	}

	return
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

func Username() string {
	var name string
	user, err := user.Current()
	if err == nil && len(user.Username) > 0 {
		name = user.Username
	} else {
		name = os.Getenv("USER")
	}

	return name
}
