package data

import (
	"github.com/andrewzulaybar/books/api/pkg/author"
	"github.com/andrewzulaybar/books/api/pkg/location"
)

// GetAuthors reads the test data in `author.json` and returns it.
func GetAuthors(as *author.Service) author.Authors {
	var buffer struct {
		Authors author.Authors `json:"authors"`
	}
	loadBuffer(&buffer, "author.json")

	locations := LoadLocations(&as.LocationService)
	return mergeLocations(buffer.Authors, locations)
}

// LoadAuthors reads the test data in `author.json` and loads it into the database.
func LoadAuthors(as *author.Service) author.Authors {
	var authors author.Authors
	for _, au := range GetAuthors(as) {
		s, a := as.PostAuthor(&au)
		if s.Err() == nil {
			s, l := as.LocationService.GetLocation(a.PlaceOfBirth.ID)
			if s.Err() == nil {
				a.PlaceOfBirth = *l
				authors = append(authors, *a)
			}
		}
	}
	return authors
}

func mergeLocations(authors author.Authors, locations location.Locations) author.Authors {
	for i := range authors {
		pob := &authors[i].PlaceOfBirth
		for j := range locations {
			if pob.ID == locations[j].ID {
				*pob = locations[j]
			}
		}
	}
	return authors
}
