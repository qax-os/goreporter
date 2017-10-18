package golint

import (
	"testing"
)

func TestGoLinter(t *testing.T) {
	GoLinter([]string{"../../../goreporter"})
}
