package foo

type Closer interface {
	Close() error
}

type FooCloser interface {
	Foo()
	Close() error
}

func ShadowArg(fc FooCloser) { // WARN fc can be Closer
	fc.Close()
	for {
		fc := 3
		println(fc + 1)
	}
}
