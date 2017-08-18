package foo

type EmptyIface interface{}

type UninterestingMethods interface {
	Foo() error
	bar() int
}

type InterestingUnexported interface {
	Foo(f string) error
	bar(f string) int
}

type St struct{}

func (s St) Foo(f string) {}

func (s St) nonExported() {}

func Bar(s St) {
	s.Foo("")
}

type NonInterestingFunc func() error
