package staticcheck

import (
	"testing"

	"github.com/wgliang/goreporter/linters/staticscan/lint/testutil"
)

func TestAll(t *testing.T) {
	c := NewChecker()
	testutil.TestAll(t, c, "")
}
