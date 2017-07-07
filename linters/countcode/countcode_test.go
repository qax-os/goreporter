package countcode

import (
	"testing"
)

func Test_CountCode(t *testing.T) {
	CountCode("../../../goreporter", "vendor/wgliang,linters,tools")
}
