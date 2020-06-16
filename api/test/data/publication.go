package data

import (
	"github.com/andrewzulaybar/books/api/pkg/publication"
	"github.com/andrewzulaybar/books/api/pkg/work"
)

// GetPublications reads the test data in `publication.json` and returns it.
func GetPublications(p *publication.Service) publication.Publications {
	var buf struct {
		Publications publication.Publications `json:"publications"`
	}
	loadBuffer(&buf, "publication.json")

	works := LoadWorks(&p.WorkService)
	return mergeWorks(buf.Publications, works)
}

// LoadPublications reads the test data in `publication.json` and loads it into the database.
func LoadPublications(p *publication.Service) publication.Publications {
	var pubs publication.Publications
	for _, pub := range GetPublications(p) {
		s, pb := p.PostPublication(&pub)
		if s.Err() == nil {
			pubs = append(pubs, *pb)
		}
	}
	return pubs
}

func mergeWorks(publications publication.Publications, works work.Works) publication.Publications {
	for i := range publications {
		for _, w := range works {
			work := &publications[i].Work
			if work.ID == w.ID {
				*work = w
			}
		}
	}
	return publications
}
