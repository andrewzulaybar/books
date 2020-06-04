package data

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"runtime"
)

func loadBuffer(buf interface{}, fileName string) {
	bytes := readFile(fileName)
	if err := json.Unmarshal(bytes, &buf); err != nil {
		panic(err)
	}
}

func readFile(fileName string) []byte {
	_, file, _, _ := runtime.Caller(0)
	path := path.Join(path.Dir(file), "json", fileName)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bytes
}
