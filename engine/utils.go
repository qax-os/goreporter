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

package engine

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
)

var (
	excepts []string
)

// DirList is a function that traverse the file directory containing the
// specified file format according to the specified rule.
func DirList(projectPath string, suffix, except string) (dirs map[string]string, err error) {
	dirs = make(map[string]string, 0)
	_, err = os.Stat(projectPath)
	if err != nil {
		glog.Errorln("dir path is invalid")
	}
	exceptsFilter(except)
	err = filepath.Walk(projectPath, func(subPath string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		if strings.HasSuffix(subPath, suffix) {
			sepIdx := strings.LastIndex(subPath, string(filepath.Separator))
			var dir string
			if sepIdx == -1 {
				dir = "."
			} else {
				if len(subPath) > sepIdx {
					dir = subPath[0:sepIdx]
				}
			}
			if ExceptPkg(dir) {
				return nil
			}
			dirs[PackageAbsPath(dir)] = dir
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return dirs, nil
}

// ExceptPkg is a function that will determine whether the package is an exception.
func ExceptPkg(pkg string) bool {
	for _, va := range excepts {
		if strings.Contains(pkg, va) {
			return true
		}
	}
	return false
}

// PackageTest is an intermediate variables.
type PackageTest struct {
	IsPass   bool    `json:"is_pass"`
	Coverage string  `json:"coverage"`
	Time     float64 `json:"time"`
}

// PackageAbsPath will gets the absolute path of the specified
// package from GOPATH's [src].
func PackageAbsPath(path string) (packagePath string) {
	_, err := os.Stat(path)
	if err != nil {
		glog.Errorln("package path is invalid")
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		glog.Fatal(err)
	}

	packagePathIndex := strings.LastIndex(absPath, "/src/")
	if -1 != packagePathIndex {
		packagePath = absPath[(packagePathIndex + 5):]
	}

	return packagePath
}

// PackageAbsPath will gets the absolute directory path of
// the specified file from GOPATH's [src].
func PackageAbsPathExceptSuffix(path string) (packagePath string) {
	if strings.LastIndex(path, string(filepath.Separator)) <= 0 {
		path, _ = os.Getwd()
	}
	path = path[0:strings.LastIndex(path, string(filepath.Separator))]
	absPath, err := filepath.Abs(path)
	if err != nil {
		glog.Errorln(err)
	}
	packagePathIndex := strings.Index(absPath, "src")
	if -1 != packagePathIndex && (packagePathIndex+4) < len(absPath) {
		packagePath = absPath[(packagePathIndex + 4):]
	}

	return packagePath
}

// ProjectName is a function that gets project's name.
func ProjectName(projectPath string) (project string) {
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		glog.Errorln(err)
	}
	projectPathIndex := strings.LastIndex(absPath, string(filepath.Separator))
	if -1 != projectPathIndex {
		project = absPath[(projectPathIndex + 1):len(absPath)]
	}

	return project
}

// AbsPath is a function that will get absolute path of file.
func AbsPath(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		glog.Errorln(err)
		return path
	}
	return absPath
}

// packageNameFromGoPath is a function that will get package's name from GOPATH.
func packageNameFromGoPath(path string) string {
	names := strings.Split(path, string(filepath.Separator))
	if len(names) >= 2 {
		return names[len(names)-2]
	}
	return "null"
}

// exceptsFilter provides function that filte except string check it's
// value is not a null string.
func exceptsFilter(except string) {
	temp := strings.Split(except, ",")
	temp = append(temp, "vendor")
	for i, _ := range temp {
		if temp[i] != "" {
			excepts = append(excepts, temp[i])
		}
	}
}
