// Copyright (C) 2013-2015 Sören Tempel
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
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	appName    = "cpod"
	appVersion = "1.9"
)

var (
	limit   = flag.Int("p", 5, "number of maximal parallel downloads")
	recent  = flag.Int("r", 0, "number of most recent episodes to download")
	version = flag.Bool("v", false, "display version number and exit")
)

var (
	logger    = log.New(os.Stderr, fmt.Sprintf("%s: ", appName), 0)
	targetDir = util.EnvDefault("CPOD_DOWNLOAD_DIR", "podcasts")
)

func main() {
	flag.Parse()
	if *version {
		logger.Fatal(appVersion)
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

			items, err := newItems(feed)
			if err != nil {
				logger.Println(err)
				return
			}

			for i := len(items) - 1; i >= 0; i-- {
				item := items[i]
				if err := getItem(feed, item); err != nil {
					logger.Println(err)
					break
				}

				if err := writeMarker(feed.Title, item.PubDate); err != nil {
					logger.Println(err)
					break
				}
			}
		}(cast)

		for *limit > 0 && counter >= *limit {
			time.Sleep(3 * time.Second)
		}
	}

	wg.Wait()
}

func fetchFeeds(och chan<- feedparser.Feed) {
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

	for url := range urlChan {
		resp, err := util.Get(url)
		if err != nil {
			logger.Println(err)
			continue
		}

		reader := resp.Body
		defer reader.Close()

		feed, err := feedparser.Parse(reader)
		if err != nil {
			logger.Println(err)
		} else {
			och <- feed
		}
	}

	close(och)
}

func newItems(cast feedparser.Feed) (items []feedparser.Item, err error) {
	unread, err := readMarker(cast.Title)
	if os.IsNotExist(err) {
		err = nil
	} else if err != nil {
		return
	}

	if *recent > 0 && len(cast.Items) >= *recent {
		cast.Items = cast.Items[0:*recent]
	}

	for _, item := range cast.Items {
		if !item.PubDate.After(unread) {
			break
		}

		if len(item.Attachment) > 0 {
			items = append(items, item)
		}
	}

	return
}

func getItem(cast feedparser.Feed, item feedparser.Item) error {
	title, err := util.Escape(cast.Title)
	if err != nil {
		return err
	}

	target := filepath.Join(targetDir, title)
	if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
		return err
	}

	fp, err := util.GetFile(item.Attachment, target)
	if err != nil {
		return err
	}

	name, err := util.Escape(item.Title)
	if err == nil {
		newfp := filepath.Join(target, name+filepath.Ext(fp))
		if err = os.Rename(fp, newfp); err != nil {
			return err
		}
	}

	return nil
}

func readMarker(name string) (marker time.Time, err error) {
	name, err = util.Escape(name)
	if err != nil {
		return
	}

	file, err := os.Open(filepath.Join(targetDir, name, ".latest"))
	if err != nil {
		return
	}

	defer file.Close()
	var timestamp int64

	if _, err = fmt.Fscanf(file, "%d\n", &timestamp); err != nil {
		return
	}

	marker = time.Unix(timestamp, 0)
	return
}

func writeMarker(name string, latest time.Time) error {
	name, err := util.Escape(name)
	if err != nil {
		return err
	}

	path := filepath.Join(targetDir, name, ".latest")
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()
	if _, err := fmt.Fprintf(file, "%d\n", latest.Unix()); err != nil {
		return err
	}

	return nil
}
