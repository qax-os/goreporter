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
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
	"go/build"
	"fmt"
)

var (
	excepts []string
)

// DirList is a function that traverse the file directory containing the
// specified file format according to the specified rule.
func DirList(projectPath string, suffix, except string) (dirs map[string]string, err error) {
	var relativePath string = ""
	dirs = make(map[string]string, 0)
	_, err = os.Stat(projectPath)
	if err != nil {
		glog.Errorln("dir path is invalid")
	}
	if build.IsLocalImport(projectPath) {
		toPos := strings.LastIndex(projectPath, string(filepath.Separator))
		relativePath = projectPath[0:toPos+1]
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
					dir = fmt.Sprintf("%s%s", relativePath, dir)
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

// FileList is a function that traverse the file is the specified file format
// according to the specified rule.
func FileList(projectPath string, suffix, except string) (files []string, err error) {
	_, err = os.Stat(projectPath)
	if err != nil {
		glog.Errorln("project path is invalid")
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

			if ExceptPkg(subPath) {
				return nil
			}
			files = append(files, subPath)
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
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

	packagePathIndex := strings.LastIndex(absPath, string(filepath.Separator)+"src"+string(filepath.Separator))
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
		project = absPath[(projectPathIndex + 1):]
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

// PackageNameFromGoPath is a function that will get package's name from GOPATH.
func PackageNameFromGoPath(path string) string {
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

// CountPercentage will count all linters' percentage.And rule is
//
//    +--------------------------------------------------+
//    |   issues    |               score                |
//    +==================================================+
//    | 5           | 100-issues*2                       |
//    +--------------------------------------------------+
//    | [5,10)      | 100 - 10 - (issues-5)*4            |
//    +--------------------------------------------------+
//    | [10,20)     | 100 - 10 - 20 - (issues-10)*5      |
//    +--------------------------------------------------+
//    | [20,40)     | 100 - 10 - 20 - 50 - (issues-20)*1 |
//    +--------------------------------------------------+
//    | [40,*)      | 0                                  |
//    +--------------------------------------------------+
//
// It will return a float64 type score.
func CountPercentage(issues int) float64 {
	if issues < 5 {
		return float64(100 - issues*2)
	} else if issues < 10 {
		return float64(100 - 10 - (issues-5)*4)
	} else if issues < 20 {
		return float64(100 - 10 - 20 - (issues-10)*5)
	} else if issues < 40 {
		return float64(100 - 10 - 20 - 50 - (issues-20)*1)
	} else {
		return 0.0
	}
}

// GetProcessUnit provides function that will get sumProcessNumber of linter's
// weight and the number of current linter's case.It will return 1 if
// sumProcessNumber/int64(number) <= 0 or  sumProcessNumber / int64(number).
// Just for communication.
func GetProcessUnit(sumProcessNumber int64, number int) int64 {
	if number == 0 {
		return sumProcessNumber
	} else if sumProcessNumber/int64(number) <= 0 {
		return int64(1)
	} else {
		return sumProcessNumber / int64(number)
	}
}
