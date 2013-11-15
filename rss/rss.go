package rss

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

type Feed struct {
	Channel struct {
		Title string `xml:"title"`
		Link  string `xml:"link"`

		Items []struct {
			PubDate string `xml:"pubDate"`
			Title   string `xml:"title"`
			Link    string `xml:"link"`

			Enclosure struct {
				Url string `xml:"url,attr"`
			} `xml:"enclosure"`
		} `xml:"item"`
	} `xml:"channel"`
}

func Parse(url string) (f *Feed, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if err = xml.Unmarshal(body, &f); err != nil {
		return
	}

	return
}
