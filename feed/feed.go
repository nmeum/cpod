package feed

import (
	"time"
)

type Feed struct {
	Title string
	Type  string
	Link  string
	Items []Item
}

type Item struct {
	Title      string
	Link       string
	Date       time.Time
	Attachment string
}

func ParseDate(date string) (t time.Time, err error) {
	formats := []string{
		time.RFC1123Z, time.RFC1123, time.RFC822Z,
		time.RFC822, time.ANSIC, time.RFC3339,
		time.RFC850, time.RubyDate, time.UnixDate,
		"2 January 2006 15:04:05 -0700", "2 January 2006 15:04:05 MST",
		"2 Jan 2006 15:04:05 -0700", "2 Jan 2006 15:04:05 MST",
		"Mon, 2 Jan 2006 15:04:05 -0700", "Mon, 2 Jan 2006 15:04:05 MST",
		"2006-01-02T15:04:05", "2006-01-02T15:04:05Z",
	}

	for _, format := range formats {
		t, err = time.Parse(format, date)
		if err == nil {
			return
		}
	}

	return
}
