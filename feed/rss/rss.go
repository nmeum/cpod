package rss

type Feed struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Items []Item `xml:"item"`
}

type Item struct {
	PubDate string `xml:"pubDate"`
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	Enclosure Enclosure `xml:"enclosure"`
}

type Enclosure struct {
	Type string `xml:"type,attr"`
	Url  string `xml:"url,attr"`
}
