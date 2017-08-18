package single

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

func BasicWrong(rc ReadCloser) { // WARN rc can be Closer
	rc.Close()
}

func OtherWrong(s St) { // WARN s can be Closer
	s.Close()
}
