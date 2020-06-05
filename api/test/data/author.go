package data

import (
	"github.com/andrewzulaybar/books/api/pkg/author"
)

// GetAuthors reads the test data in `author.json` and returns it.
func GetAuthors(_ *author.Service) author.Authors {
	var buf struct {
		Authors author.Authors `json:"authors"`
	}
	loadBuffer(&buf, "author.json")
	return buf.Authors
}

// LoadAuthors reads the test data in `author.json` and loads it into the database.
func LoadAuthors(as *author.Service) author.Authors {
	var authors author.Authors
	for _, au := range GetAuthors(as) {
		s, a := as.PostAuthor(&au)
		if s.Err() == nil {
			authors = append(authors, *a)
		}
	}
	return authors
}
