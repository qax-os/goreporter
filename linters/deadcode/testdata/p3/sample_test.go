package p

import (
	"testing"
)

func TestX(t *testing.T) {
	if x != 42 {
		t.Fatalf("x != 42")
	}
}
