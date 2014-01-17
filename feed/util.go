package feed

import (
	"github.com/nmeum/cpod/feed/atom"
	"github.com/nmeum/cpod/feed/rss"
)

func convertRss(r rss.Feed) (f Feed) {
	f.Title = r.Title
	f.Link = r.Link

	for _, i := range r.Items {
		item := Item{
			Title:      i.Title,
			Link:       i.Link,
			Date:       i.PubDate,
			Attachment: i.Enclosure.Url,
		}

		f.Items = append(f.Items, item)
	}

	return
}

func convertAtom(a atom.Feed) (f Feed) {
	f.Title = a.Title
	f.Link = findLink(a.Links).Href

	for _, e := range a.Entries {
		item := Item{
			Title:      e.Title,
			Link:       findLink(e.Links).Href,
			Date:       e.Published,
			Attachment: findAttachment(e.Links).Href,
		}

		f.Items = append(f.Items, item)
	}

	return
}

func findLink(links []atom.Link) (l atom.Link) {
	for _, link := range links {
		if link.Type == "text/html" {
			l = link
		}
	}

	return
}

func findAttachment(links []atom.Link) (l atom.Link) {
	for _, link := range links {
		if link.Rel == "enclosure" {
			l = link
		}
	}

	return
}
