package simple

import (
	"testing"

	"github.com/wgliang/goreporter/linters/staticscan/lint/testutil"
)

func TestAll(t *testing.T) {
	testutil.TestAll(t, NewChecker(), "")
}
