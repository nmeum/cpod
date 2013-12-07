package atom

type Feed struct {
	Title   string  `xml:"title"`
	Links   []Link  `xml:"link"`
	Entries []Entry `xml:"entry"`
}

type Entry struct {
	Published string `xml:"published"`
	Title     string `xml:"title"`
	Links     []Link `xml:"link"`
}

type Link struct {
	Type string `xml:"type,attr"`
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
}
