package data

import (
	"github.com/andrewzulaybar/books/api/pkg/publication"
)

// GetPublications reads the test data in `publication.json` and returns it.
func GetPublications(_ *publication.Service) publication.Publications {
	var buf struct {
		Publications publication.Publications `json:"publications"`
	}
	loadBuffer(&buf, "publication.json")
	return buf.Publications
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
