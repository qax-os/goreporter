package engine

import (
	"testing"
)

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

func Test_DirList_NoPath(t *testing.T) {
	DirList("../../nopath", ".go", "_test.go")
}

func Test_ExceptPkg_NoPath(t *testing.T) {
	ExceptPkg("execpt", "github.com/except/src")
}

func Test_PackageAbsPath_NoPath(t *testing.T) {
	PackageAbsPath("../nopath")
}

func Test_PackageAbsPathExceptSuffix_NoPath(t *testing.T) {
	PackageAbsPathExceptSuffix("../engine/nopath.go")
}

func Test_ProjectName_NoPath(t *testing.T) {
	ProjectName("../nopath")
}

func Test_AbsPath_NoPath(t *testing.T) {
	AbsPath("../nopath")
}
