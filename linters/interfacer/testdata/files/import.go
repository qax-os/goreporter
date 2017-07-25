package foo

import (
	"io"
	"os"
)

func Empty() {
}

func Basic(c io.Closer) {
	c.Close()
}

func BasicWrong(rc io.ReadCloser) { // WARN rc can be io.Closer
	rc.Close()
}

type St struct{}

func (s *St) Basic(c io.Closer) {
	c.Close()
}

func (s *St) BasicWrong(rc io.ReadCloser) { // WARN rc can be io.Closer
	rc.Close()
}

type Namer interface {
	Name() string
}

func WalkFuncImpl(path string, info os.FileInfo, err error) error {
	info.Name()
	return nil
}

func WalkFuncImplWrong(path string, info os.FileInfo, err error) { // WARN info can be Namer
	info.Name()
}

type MyWalkFunc func(path string, info os.FileInfo, err error) error
