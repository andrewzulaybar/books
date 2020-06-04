package data

import (
	"encoding/json"

	"github.com/andrewzulaybar/books/api/pkg/publication"
)

// GetPublications reads the test data in `publication.json` and returns it.
func GetPublications(_ *publication.Service) publication.Publications {
	bytes := readFile("publication.json")

	var pubs struct {
		Publications publication.Publications `json:"publications"`
	}
	if err := json.Unmarshal(bytes, &pubs); err != nil {
		panic(err)
	}
	return pubs.Publications
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
