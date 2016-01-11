// Copyright (C) 2013-2016 Sören Tempel
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package util

import (
	"errors"
	"os"
	"os/signal"
	"os/user"
	"strings"
	"unicode"
)

// Lock creates a lockfile at the given path and creates a signal
// handler which removes the lockfile on interrupt or kill.
func Lock(path string) (err error) {
	_, err = os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0444)
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
// filename. It returns an error if the escaped string is empty.
func Escape(name string) (escaped string, err error) {
	mfunc := func(r rune) rune {
		switch {
		case unicode.IsLetter(r):
			return r
		case unicode.IsNumber(r):
			return r
		case unicode.IsPunct(r) && r != os.PathSeparator:
			return r
		case unicode.IsSpace(r):
			return ' '
		}

		return -1
	}

	escaped = strings.Map(mfunc, name)
	escaped = strings.TrimSpace(escaped)

	for len(escaped) > 0 && escaped[0] == '.' {
		if len(escaped) >= 1 {
			escaped = escaped[1:]
		} else {
			escaped = ""
		}
	}

	if len(escaped) <= 0 {
		err = errors.New("couldn't escape title")
	}

	return
}

// HomeDir returns the absolute path of the users home directory. It
// tries to determine the path using os/user first and if that doesn't
// work it returns the value of the HOME environment variable.
func HomeDir() string {
	var home string
	user, err := user.Current()
	if err == nil {
		home = user.HomeDir
	} else {
		home = os.Getenv("HOME")
	}

	return home
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
