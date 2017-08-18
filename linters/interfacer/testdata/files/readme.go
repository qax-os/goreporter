package foo

import (
	"io/ioutil"
	"os"
)

func ProcessInput(f *os.File) error { // WARN f can be io.Reader
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	return processBytes(b)
}

func processBytes(b []byte) error { return nil }
