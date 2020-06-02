package data

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"runtime"

	"github.com/andrewzulaybar/books/api/pkg/work"
)

func GetWorks() work.Works {
	bytes := readFile("work.json")

	var works struct {
		Works work.Works `json:"works"`
	}
	if err := json.Unmarshal(bytes, &works); err != nil {
		panic(err)
	}
	return works.Works
}

func LoadWorks(w *work.Service) work.Works {
	var works work.Works
	for _, work := range GetWorks() {
		s, wk := w.PostWork(&work)
		if s.Err() == nil {
			works = append(works, *wk)
		}
	}
	return works
}

func readFile(fileName string) []byte {
	_, file, _, _ := runtime.Caller(0)
	path := path.Join(path.Dir(file), fileName)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bytes
}
