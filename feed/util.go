package feed

import (
	"github.com/nmeum/cpod/feed/atom"
)

func findLink(links []atom.Link) (l atom.Link) {
	var score int

	for _, link := range links {
		if link.Rel == "alternate" && link.Type == "text/html" {
			return link
		}

		if score < 3 && link.Type == "text/html" {
			l = link
		}

		if score < 2 && link.Rel == "self" {
			l = link
		}

		if score < 1 && link.Rel == "" {
			l = link
		}

		if score <= 0 {
			l = link
		}
	}

	return
}

func findAttachment(links []atom.Link) (l atom.Link) {
	for _, link := range links {
		if link.Rel == "enclosure" {
			return link
		}
	}

	return
}
