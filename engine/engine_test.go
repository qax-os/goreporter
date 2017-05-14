package engine

import (
	"testing"
)

func Test_Engine(t *testing.T) {
	report := NewReporter("")
	report.Engine("../../logcool", "")
}

func Test_DirList(t *testing.T) {
	DirList("../../goreporter", ".go", "_test.go")
}

func Test_ExceptPkg(t *testing.T) {
	ExceptPkg("execpt", "github.com/except/src")
}

func Test_PackageAbsPath(t *testing.T) {
	PackageAbsPath("../engine")
}

func Test_PackageAbsPathExceptSuffix(t *testing.T) {
	PackageAbsPathExceptSuffix("../engine/engine.go")
}

func Test_ProjectName(t *testing.T) {
	ProjectName("../engine")
}

func Test_AbsPath(t *testing.T) {
	AbsPath("../engine")
}
