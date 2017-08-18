package simpler

import (
	"testing"

	"github.com/360EntSecGroup-Skylar/goreporter/linters/simpler/lint/testutil"
)

func TestAll(t *testing.T) {
	testutil.TestAll(t, NewChecker(), "")
}
