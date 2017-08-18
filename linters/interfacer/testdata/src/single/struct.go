package single

type St struct{}

func (s *St) Close() {}

func (s *St) Basic(c Closer) {
	c.Close()
}

func (s *St) BasicWrong(rc ReadCloser) { // WARN rc can be Closer
	rc.Close()
}
