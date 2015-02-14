package util

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

const (
	// Number of times failed HTTP request is retried.
	retry = 3

	// Number of maximal allowed redirects.
	maxRedirects = 10

	// HTTP User-Agent.
	useragent = "cpod"
)

// Get works like http.get but also retries the get a few times if it
// failed.
func Get(uri string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return
	}

	return doReq(req)
}

// GetFile downloads the file located at the given uri and saves it in
// the directory specified as target. It also supports continuous
// downloads.
func GetFile(uri, target string) (fp string, err error) {
	if err = os.MkdirAll(target, 0755); err != nil {
		return
	}

	fn, err := filename(uri)
	if err != nil {
		return
	}

	partPath := filepath.Join(target, fmt.Sprintf("%s.part", fn))
	if _, err = os.Open(partPath); os.IsNotExist(err) {
		if err = newGet(uri, partPath); err != nil {
			return
		}
	} else {
		if err = resumeGet(uri, partPath); err != nil {
			return
		}
	}

	fp = filepath.Join(target, fn)
	if err = os.Rename(partPath, fp); err != nil {
		return
	}

	return
}

// ResumeGet resumes a download which was already started.
func resumeGet(uri, target string) error {
	file, err := os.OpenFile(target, os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Range", fmt.Sprintf("bytes=%d-", fi.Size()))
	resp, err := doReq(req)
	if err != nil {
		return err
	}

	reader := resp.Body
	defer reader.Close()

	if _, err = io.Copy(file, reader); err != nil {
		return err
	}

	return nil
}

// newGet starts a new file download.
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

// filename returns the fiilename of an URL. It removes the query
// parameters etc.
func filename(uri string) (fn string, err error) {
	u, err := url.Parse(uri)
	if err != nil {
		return
	}

	fn = filepath.Base(u.Path)
	if fn == "/" || fn == "." {
		fn = "unnamed"
	}

	return
}

// doReq does the same as net.client.Do but it retries sending the
// request if it failed. Furthermore, it also ensure that headers
// remain the same after a redirect and it adds a User-Agent header.
func doReq(req *http.Request) (resp *http.Response, err error) {
	req.Header.Add("User-Agent", useragent)
	client := headerClient(req.Header)

	for i := 1; i <= retry; i++ {
		resp, err = client.Do(req)
		if nerr, ok := err.(net.Error); ok && (nerr.Temporary() || nerr.Timeout()) {
			time.Sleep((time.Duration)(i*3) * time.Second)
		} else {
			break
		}
	}

	return
}

// headerClient returns a client witch a custom CheckRedirect function
// which ensures that the given headers will be readded after a redirect.
func headerClient(headers http.Header) *http.Client {
	redirectFunc := func(req *http.Request, via []*http.Request) error {
		if len(via) >= maxRedirects {
			return errors.New("too many redirects")
		}

		req.Header = headers
		return nil
	}

	return &http.Client{CheckRedirect: redirectFunc}
}
