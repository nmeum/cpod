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

// Number of times a download is retried.
const retry = 3

func Filename(uri string) (fn string, err error) {
	u, err := url.Parse(uri)
	if err != nil {
		return
	}

	fn = filepath.Base(u.Path)
	return
}

func Get(uri string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return
	}

	return doReq(req)
}

func GetFile(uri, target string) error {
	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	fn, err := Filename(uri)
	if err != nil {
		return err
	}

	partPath := filepath.Join(target, fmt.Sprintf("%s.part", fn))
	_, err = os.Open(partPath)
	if os.IsNotExist(err) {
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

func resumeGet(uri, target string) error {
	file, err := os.OpenFile(target, os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

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

	if _, err = io.Copy(file, reader); err != nil {
		return err
	}

	return nil
}

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
