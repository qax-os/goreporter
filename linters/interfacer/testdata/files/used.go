package foo

type Closer interface {
	Close() error
}

type Reader interface {
	Read(p []byte) (n int, err error)
}

type ReadCloser interface {
	Reader
	Closer
}

type St struct{}

func (s St) Read(p []byte) (int, error) {
	return 0, nil
}
func (s St) Close() error {
	return nil
}
func (s St) Other() {}

func FooCloser(c Closer) {
	c.Close()
}

func FooSt(s St) {
	s.Other()
}

func Bar(s St) {
	s.Close()
	FooSt(s)
}

func BarWrong(s St) { // WARN s can be Closer
	s.Close()
	FooCloser(s)
}

func extra(n int, cs ...Closer) {}

func ArgExtraWrong(s1 St) { // WARN s1 can be Closer
	var s2 St
	s1.Close()
	s2.Close()
	extra(3, s1, s2)
}

func Assigned(s St) {
	s.Close()
	var s2 St
	s2 = s
	_ = s2
}

func Declared(s St) {
	s.Close()
	var s2 St = s
	_ = s2
}

func AssignedIface(s St) { // WARN s can be Closer
	s.Close()
	var c Closer
	c = s
	_ = c
}

func AssignedIfaceDiff(s St) { // WARN s can be ReadCloser
	s.Close()
	var r Reader
	r = s
	_ = r
}

func doRead(r Reader) {
	b := make([]byte, 10)
	r.Read(b)
}

func ArgIfaceDiff(s St) { // WARN s can be ReadCloser
	s.Close()
	doRead(s)
}
