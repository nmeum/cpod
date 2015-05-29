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
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const (
	// Number of times a failed HTTP request is retried.
	retry = 3

	// Number of maximal allowed redirects.
	maxRedirects = 10

	// Time limit for HTTP requests.
	timeout = 30 * time.Second

	// HTTP User-Agent.
	useragent = "cpod"
)

// Get performs a HTTP GET request, just like http.get, however, it has
// a few handy extra features: I adds a User-Agent header and it retries
// a failed get request if the error was a temporary one.
func Get(uri string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return
	}

	return doReq(req)
}

// GetFile downloads the file from the given uri and stores it in the
// specified target directory. If a download was interrupted previously
// GetFile is able to resume it.
func GetFile(uri, target string) (fp string, err error) {
	if err = os.MkdirAll(target, 0755); err != nil {
		return
	}

	fn, err := filename(uri)
	if err != nil {
		return
	}

	fp = filepath.Join(target, fn)
	partPath := fmt.Sprintf("%s.part", fp)
	if _, err = os.Open(partPath); os.IsNotExist(err) {
		if err = newGet(uri, partPath); err != nil {
			return
		}
	} else {
		if err = resumeGet(uri, partPath); err != nil {
			return
		}
	}

	if err = os.Rename(partPath, fp); err != nil {
		return
	}

	return
}

// resumeGet resumes an canceled download started by the newGet
// function.
func resumeGet(uri, target string) error {
	fi, err := os.Stat(target)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Range", fmt.Sprintf("bytes=%d-", fi.Size()))
	req.Header.Add("If-Range", fi.ModTime().Format(time.RFC1123))

	resp, err := doReq(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusPartialContent {
		return newGet(uri, target)
	}

	file, err := os.OpenFile(target, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	defer file.Close()
	if _, err = io.Copy(file, resp.Body); err != nil {
		return err
	}

	return nil
}

// newGet starts a new file download, if the download wasn't completed
// it can be resumed later on using the resumeGet function.
func newGet(uri, target string) error {
	resp, err := Get(uri)
	if err != nil {
		return err
	}

	reader := resp.Body
	defer reader.Close()

	file, err := os.Create(target)
	if err != nil {
		return err
	}

	defer file.Close()
	if _, err = io.Copy(file, reader); err != nil {
		return err
	}

	return err
}

// filename returns the fiilename of an URL. Basically it just uses
// path.Base to determine the filename but it also removes queries.
// Furthermore it also guarantees that the filename is not empty by
// setting it to "unnamed" if it couldn't determine a proper filename.
func filename(uri string) (fn string, err error) {
	u, err := url.Parse(uri)
	if err != nil {
		return
	}

	fn = strings.TrimSpace(path.Base(u.Path))
	if len(fn) <= 0 || fn == "/" || fn == "." {
		fn = "unnamed"
	}

	return
}

// doReq does the same as net.client.Do but it retries sending the request if a
// temporary error on layer 4 is encountered. Furthermore, it also ensure that
// headers remain the same after a redirect and it adds a User-Agent header.
func doReq(req *http.Request) (resp *http.Response, err error) {
	req.Header.Add("User-Agent", useragent)
	client := headerClient(req.Header)

	for i := 1; i <= retry; i++ {
		resp, err = client.Do(req)
		if nerr, ok := err.(net.Error); ok && (nerr.Temporary() || nerr.Timeout()) {
			time.Sleep(time.Duration(i*3) * time.Second)
			continue
		}

		break
	}

	return
}

// headerClient returns a client witch a custom CheckRedirect function which
// ensures that the given headers will be readded after a redirect. Furthermore
// it also set a time limit for requests.
func headerClient(headers http.Header) *http.Client {
	redirectFunc := func(req *http.Request, via []*http.Request) error {
		if len(via) >= maxRedirects {
			return errors.New("too many redirects")
		}

		req.Header = headers
		return nil
	}

	return &http.Client{CheckRedirect: redirectFunc, Timeout: timeout}
}
