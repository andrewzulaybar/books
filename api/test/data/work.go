package data

import (
	"encoding/json"

	"github.com/andrewzulaybar/books/api/pkg/work"
)

// GetWorks reads the test data in `work.json` and returns it.
func GetWorks(_ *work.Service) work.Works {
	bytes := readFile("work.json")

	var works struct {
		Works work.Works `json:"works"`
	}
	if err := json.Unmarshal(bytes, &works); err != nil {
		panic(err)
	}
	return works.Works
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
