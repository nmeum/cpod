package util

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

// retry describes the amount of times a failed http get request is
// retried if the error is temporary or a timeout error.
const retry = 3

// Filename returns the fiilename of an URL. It removes the query
// parameters etc.
func Filename(uri string) (fn string, err error) {
	u, err := url.Parse(uri)
	if err != nil {
		return
	}

	fn = filepath.Base(u.Path)
	if fn == string(os.PathSeparator) || fn == "." {
		fn = "unnamed"
	}

	return
}

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
func GetFile(uri, target string) error {
	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	fn, err := Filename(uri)
	if err != nil {
		return err
	}

	partPath := filepath.Join(target, fmt.Sprintf("%s.part", fn))
	if _, err = os.Open(partPath); os.IsNotExist(err) {
		if err = newGet(uri, partPath); err != nil {
			return err
		}
	} else if err == nil {
		if err = resumeGet(uri, partPath); err != nil {
			return err
		}
	}

	return os.Rename(partPath, filepath.Join(target, fn))
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

	req.Header.Add("Range", fmt.Sprint("bytes=%d-", fi.Size()))
	resp, err := doReq(req)
	if err != nil {
		return err
	}

	reader := resp.Body
	defer reader.Close()

	if resp.Status != "206" {
		fmt.Println("Doesn't support Partial Content", resp.Status)
	}

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

// doReq does the same as net.client.Do but it retries sending the
// request if it failed.
func doReq(req *http.Request) (resp *http.Response, err error) {
	client := http.DefaultClient
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
