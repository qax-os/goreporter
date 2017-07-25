package foo

type Stringer interface {
	String() string
}

type st struct{}

func (s st) String() string {
	return "foo"
}

func Exported(s st) string { // WARN s can be Stringer
	return s.String()
}

func unexported(s st) string {
	return s.String()
}
