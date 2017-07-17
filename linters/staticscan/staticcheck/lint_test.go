package staticcheck

import (
	"testing"

	"github.com/360EntSecGroup-Skylar/goreporter/linters/staticscan/lint/testutil"
)

func TestAll(t *testing.T) {
	c := NewChecker()
	testutil.TestAll(t, c, "")
}
