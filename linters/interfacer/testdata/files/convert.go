package foo

type mint int

func (m mint) Close() error {
	return nil
}

type mint2 mint

func ConvertNamed(m mint) {
	m.Close()
	_ = mint2(m)
}

func ConvertBasic(m mint) {
	m.Close()
	println(int(m))
}

type Closer interface {
	Close() error
}

func ConvertIface(m mint) { // WARN m can be Closer
	m.Close()
	_ = Closer(m)
}
