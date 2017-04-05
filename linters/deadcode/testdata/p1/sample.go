package main

// main is used
func main() {
	f(x)
	return
}

// x is used
var x int

// unused is unused
var unused int

// f is used
func f(x int) {
}

// g is unused
func g(x int) {
}

// H is exported
func H(x int) {
}

// init is used
func init() {
}

var _ int

func h(x int) {
	if x > 0 {
		h(x - 1)
	}
}
