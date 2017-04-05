package unittest

import (
	"testing"
)

func Test_UnitTest(t *testing.T) {
	UnitTest("../aligncheck")
}

func Test_GoListWithImportPackages(t *testing.T) {
	GoListWithImportPackages("../copycheck")
}
