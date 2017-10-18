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

package gofmt

import (
	"reflect"
	"testing"
)

var wantFmtResult = []string{
	"../../engine/processbar/processbar.go",
}

func Test_GoFmt(t *testing.T) {
	res, err := GoFmt([]string{"../../engine"})
	if err != nil {
		t.Error("go vet failed.")
	} else {
		if !reflect.DeepEqual(res, wantFmtResult) {
			t.Errorf("want %v, but got %v", wantFmtResult, res)
		}
	}
}
