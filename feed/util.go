package feed

import (
	"github.com/nmeum/cpod/feed/atom"
)

func findLink(links []atom.Link) atom.Link {
	var score int
	var match atom.Link

	for _, link := range links {
		switch {
		case link.Rel == "alternate" && link.Type == "text/html":
			return link
		case score < 3 && link.Type == "text/html":
			score = 3
			match = link
		case score < 2 && link.Rel == "self":
			score = 2
			match = link
		case score < 1 && link.Rel == "":
			score = 1
			match = link
		case &match == nil:
			match = link
		}
	}

	return match
}

func findAttachment(links []atom.Link) atom.Link {
	for _, link := range links {
		if link.Rel == "enclosure" {
			return link
		}
	}

	return atom.Link{}
}
