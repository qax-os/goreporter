package foo

type FooCloser interface {
	Foo()
	Close() error
}

type Barer interface {
	Bar(fc FooCloser) int
}

type St struct{}

func (s St) Bar(fc FooCloser) int {
	return 2
}

func Foo(s St) { // WARN s can be Barer
	_ = s.Bar(nil)
}

func Bar(fc FooCloser) int {
	fc.Close()
	return 3
}
