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

package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/nmeum/cpod/util"
	"github.com/nmeum/go-feedparser"
	"html"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	appName    = "cpod"
	appVersion = "2.0dev"
)

var (
	limit   = flag.Int("p", 5, "number of maximal parallel downloads")
	recent  = flag.Int("r", 0, "number of most recent episodes to download")
	version = flag.Bool("v", false, "display version number and exit")
)

var (
	logger    = log.New(os.Stderr, fmt.Sprintf("%s: ", appName), 0)
	targetDir = filepath.Join(util.HomeDir(), "podcasts")
)

func main() {
	flag.Parse()
	if *version {
		logger.Fatal(appVersion)
	} else if flag.NArg() >= 1 {
		targetDir = flag.Arg(0)
	}

	lockPath := filepath.Join(os.TempDir(), fmt.Sprintf("%s-%s", appName, util.Username()))
	if err := util.Lock(lockPath); os.IsExist(err) {
		logger.Fatalf("database is locked, remove %q to force unlock\n", lockPath)
	} else if err != nil {
		logger.Fatal(err)
	}

	update()
	if err := os.Remove(lockPath); err != nil {
		logger.Fatal(err)
	}
}

func update() {
	var wg sync.WaitGroup
	var counter int

	feeds := make(chan feedparser.Feed)
	go fetchFeeds(feeds)

	for cast := range feeds {
		wg.Add(1)
		counter++

		go func(feed feedparser.Feed) {
			defer func() {
				wg.Done()
				counter--
			}()

			title, err := util.Escape(html.UnescapeString(feed.Title))
			if err != nil {
				return
			} else {
				feed.Title = title
			}

			items, err := newItems(feed)
			if err != nil {
				logger.Println(err)
				return
			}

			for _, item := range items {
				wg.Add(1)
				go func(i feedparser.Item) {
					if err := getItem(feed, item); err != nil {
						logger.Println(err)
					}
					wg.Done()
				}(item)
			}
		}(cast)

		for *limit > 0 && counter >= *limit {
			time.Sleep(3 * time.Second)
		}
	}

	wg.Wait()
}

func fetchFeeds(ch chan<- feedparser.Feed) {
	file, err := os.Open(filepath.Join(targetDir, "urls.txt"))
	if err != nil {
		logger.Fatal(err)
	}
	defer file.Close()

	urlChan := make(chan string)
	go func(r io.Reader) {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			urlChan <- scanner.Text()
		}

		if err := scanner.Err(); err != nil {
			logger.Fatal(err)
		}

		close(urlChan)
	}(file)

	var wg sync.WaitGroup
	for url := range urlChan {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			resp, err := util.Get(u)
			if err != nil {
				logger.Println(err)
				return
			}

			reader := resp.Body
			defer reader.Close()

			feed, err := feedparser.Parse(reader)
			if err != nil {
				logger.Println(err)
			} else {
				ch <- feed
			}
		}(url)
	}

	wg.Wait()
	close(ch)
}

func newItems(cast feedparser.Feed) (items []feedparser.Item, err error) {
	var latest time.Time
	latestFi, err := findLatest(filepath.Join(targetDir, cast.Title))
	if err == nil {
		latest = latestFi.ModTime()
	} else if os.IsNotExist(err) {
		err = nil
	} else {
		return
	}

	if *recent > 0 && len(cast.Items) >= *recent {
		cast.Items = cast.Items[0:*recent]
	}

	for _, item := range cast.Items {
		if !item.PubDate.After(latest) {
			break
		}

		if len(item.Attachment) > 0 {
			items = append(items, item)
		}
	}

	return
}

func getItem(cast feedparser.Feed, item feedparser.Item) error {
	target := filepath.Join(targetDir, cast.Title)
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return err
	}

	fp, err := util.GetFile(item.Attachment, target)
	if err != nil {
		return err
	}

	name, err := util.Escape(html.UnescapeString(item.Title))
	if err == nil {
		newfp := filepath.Join(target, name+filepath.Ext(fp))
		if err = os.Rename(fp, newfp); err == nil {
			fp = newfp
		} else {
			return err
		}
	}

	return os.Chtimes(fp, item.PubDate, item.PubDate)
}

func findLatest(fp string) (fi os.FileInfo, err error) {
	dir, err := os.Open(fp)
	if err != nil {
		return
	}
	defer dir.Close()

	files, err := dir.Readdir(0)
	if err != nil {
		return
	}

	var latest *os.FileInfo
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".part" {
			continue
		}

		t := (*latest).ModTime()
		if latest == nil || file.ModTime().Before(t) {
			latest = &file
		}
	}

	if latest == nil {
		err = os.ErrNotExist
	} else {
		fi = *latest
	}

	return
}
