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

// Lock creates a lockfile at the given path and creates a signal
// handler which removes the lockfile on interrupt or kill.
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

// Escape escapes the given data to make sure it is safe to use it as a
// filename. It also replaces spaces and other seperation characters
// with the '-' character. It returns an error if the escaped string is
// empty.
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

// EnvDefault returns the value of the given environment variable
// key if it is not empty. If it is empty it returns the fallback
// as an absolute path joined with the users home.
func EnvDefault(key, fallback string) string {
	dir := os.Getenv(key)
	if len(dir) > 0 {
		return dir
	}

	var home string
	user, err := user.Current()
	if err == nil {
		home = user.HomeDir
	} else {
		home = os.Getenv("HOME")
	}

	return filepath.Join(home, fallback)
}

// Username returns the username of the current user. It tries to
// determine the username using os/user first and if that doesn't
// work it returns the value of the USER environment variable.
func Username() string {
	var name string
	user, err := user.Current()
	if err == nil {
		name = user.Username
	} else {
		name = os.Getenv("USER")
	}

	return name
}
