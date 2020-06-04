package data

import (
	"io/ioutil"
	"path"
	"runtime"
)

func readFile(fileName string) []byte {
	_, file, _, _ := runtime.Caller(0)
	path := path.Join(path.Dir(file), fileName)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bytes
}
