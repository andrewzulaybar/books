package data

import (
	"github.com/andrewzulaybar/books/api/pkg/work"
)

// GetWorks reads the test data in `work.json` and returns it.
func GetWorks(_ *work.Service) work.Works {
	var buf struct {
		Works work.Works `json:"works"`
	}
	loadBuffer(&buf, "work.json")
	return buf.Works
}

// LoadWorks reads the test data in `work.json` and loads it into the database.
func LoadWorks(w *work.Service) work.Works {
	var works work.Works
	for _, work := range GetWorks(w) {
		s, wk := w.PostWork(&work)
		if s.Err() == nil {
			works = append(works, *wk)
		}
	}
	return works
}
