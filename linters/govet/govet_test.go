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
	"reflect"
	"testing"
)

var wantVetResult = []string{
	`../../engine/config.go:204: struct field FilesNum repeats json tag "content" also at ../../engine/config.go:203`,
	`../../engine/config.go:205: struct field Quality repeats json tag "content" also at ../../engine/config.go:203`,
	`../../engine/reporter.go:128: call of strategy.Percentage copies lock value: engine.Summaries`,
	`../../engine/strategy_copycheck.go:69: Percentage passes lock by value: engine.Summaries`,
	`../../engine/strategy_countcode.go:51: Percentage passes lock by value: engine.Summaries`,
	`../../engine/strategy_cyclo.go:76: Percentage passes lock by value: engine.Summaries`,
	`../../engine/strategy_deadcode.go:66: Percentage passes lock by value: engine.Summaries`,
	`../../engine/strategy_dependgraph.go:38: Percentage passes lock by value: engine.Summaries`,
	`../../engine/strategy_depth.go:74: Percentage passes lock by value: engine.Summaries`,
	`../../engine/strategy_importpackages.go:38: Percentage passes lock by value: engine.Summaries`,
	`../../engine/strategy_interfacer.go:66: Percentage passes lock by value: engine.Summaries`,
	`../../engine/strategy_lint.go:64: return copies lock value: engine.Summaries`,
	`../../engine/strategy_lint.go:67: Percentage passes lock by value: engine.Summaries`,
	`../../engine/strategy_simplecode.go:61: return copies lock value: engine.Summaries`,
	`../../engine/strategy_simplecode.go:64: Percentage passes lock by value: engine.Summaries`,
	`../../engine/strategy_spellcheck.go:64: Percentage passes lock by value: engine.Summaries`,
	`../../engine/strategy_unittest.go:94: Percentage passes lock by value: engine.Summaries`,
}

func Test_GoVet(t *testing.T) {
	res, err := GoVet([]string{"../../engine"})
	if err != nil {
		t.Error("go vet failed.")
	} else {
		if !reflect.DeepEqual(res, wantVetResult) {
			t.Errorf("want %v, but got %v", wantVetResult, res)
		}
	}
}
