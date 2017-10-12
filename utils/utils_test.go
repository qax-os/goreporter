// Copyright 2017 The GoReporter Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"testing"
)

func Test_DirList(t *testing.T) {
	DirList("../goreporter", ".go", "_test.go")
}

func Test_ExceptPkg(t *testing.T) {
	ExceptPkg("github.com/except/src")
}

func Test_PackageAbsPath(t *testing.T) {
	PackageAbsPath("./engine")
}

func Test_PackageAbsPathExceptSuffix(t *testing.T) {
	PackageAbsPathExceptSuffix("./engine/engine.go")
}

func Test_ProjectName(t *testing.T) {
	ProjectName("./engine")
}

func Test_AbsPath(t *testing.T) {
	AbsPath("./engine")
}

func Test_DirList_NoPath(t *testing.T) {
	DirList("../../nopath", ".go", "_test.go")
}

func Test_ExceptPkg_NoPath(t *testing.T) {
	ExceptPkg("github.com/except/src")
}

func Test_PackageAbsPath_NoPath(t *testing.T) {
	PackageAbsPath("../nopath")
}

func Test_PackageAbsPathExceptSuffix_NoPath(t *testing.T) {
	PackageAbsPathExceptSuffix("./engine/nopath.go")
}

func Test_ProjectName_NoPath(t *testing.T) {
	ProjectName("../nopath")
}

func Test_AbsPath_NoPath(t *testing.T) {
	AbsPath("../nopath")
}
