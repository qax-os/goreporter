package simple

import (
	"testing"

	"github.com/360EntSecGroup-Skylar/goreporter/linters/simplecode/lint/testutil"
)

func TestAll(t *testing.T) {
	testutil.TestAll(t, Funcs, "../../")
}
