package staticcheck

import (
	"testing"
)

func TestStaticCheck(t *testing.T) {
	StaticCheck([]string{"../../engine"})
}
