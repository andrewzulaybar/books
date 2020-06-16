package data

import (
	"github.com/andrewzulaybar/books/api/pkg/author"
	"github.com/andrewzulaybar/books/api/pkg/work"
)

// GetWorks reads the test data in `work.json` and returns it.
func GetWorks(ws *work.Service) work.Works {
	var buffer struct {
		Works work.Works `json:"works"`
	}
	loadBuffer(&buffer, "work.json")

	authors := LoadAuthors(&ws.AuthorService)
	return mergeAuthors(buffer.Works, authors)
}

// LoadWorks reads the test data in `work.json` and loads it into the database.
func LoadWorks(ws *work.Service) work.Works {
	var works work.Works
	for _, work := range GetWorks(ws) {
		s, w := ws.PostWork(&work)
		if s.Err() == nil {
			s, a := ws.AuthorService.GetAuthor(w.Author.ID)
			if s.Err() == nil {
				w.Author = *a
				works = append(works, *w)
			}
		}
	}
	return works
}

func mergeAuthors(works work.Works, authors author.Authors) work.Works {
	for i := range works {
		author := &works[i].Author
		for j := range authors {
			if author.ID == authors[j].ID {
				*author = authors[j]
			}
		}
	}
	return works
}
