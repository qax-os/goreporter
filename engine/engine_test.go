package engine

import (
	"testing"
)

func Test_Engine(t *testing.T) {
	report := NewReporter("")
	report.Engine("../../logcool", "")
}
