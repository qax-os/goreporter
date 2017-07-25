package foo

func Empty() {
}

type Closer interface {
	Close()
}

type ReadCloser interface {
	Closer
	Read()
}

func Basic(c Closer) {
	c.Close()
}

func BasicInteresting(rc ReadCloser) {
	rc.Read()
	rc.Close()
}

func BasicWrong(rc ReadCloser) { // WARN rc can be Closer
	rc.Close()
}

type St struct{}

func (s *St) Basic(c Closer) {
	c.Close()
}

func (s *St) BasicWrong(rc ReadCloser) { // WARN rc can be Closer
	rc.Close()
}
