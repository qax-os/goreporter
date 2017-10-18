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

package govet

import (
	"bytes"
	"os/exec"
	"strings"
)

// GoVet if a function that will run command go tool vet,return all result of
// warnning issues.
func GoVet(packagePath []string) (goVetData []string, err error) {
	cmd := exec.Command("go", append([]string{"tool", "vet"}, packagePath...)...)
	var out, outerr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &outerr
	err = cmd.Run()
	if err != nil {
		goVetData = strings.Split(strings.TrimSuffix(outerr.String(), "\n"), "\n")
		if len(goVetData) > 0 {
			return goVetData, nil
		} else {
			return goVetData, err
		}
	}
	goVetData = strings.Split(strings.TrimSuffix(out.String(), "\n"), "\n")
	return goVetData, nil
}
