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
	"log"
	"sync"
	"testing"

	"github.com/360EntSecGroup-Skylar/goreporter/engine/processbar"
	"github.com/facebookgo/inject"
	"github.com/golang/glog"
)

func Test_Engine(t *testing.T) {
	synchronizer := &Synchronizer{
		LintersProcessChans:   make(chan int64, 20),
		LintersFinishedSignal: make(chan string, 10),
	}
	syncRW := &sync.RWMutex{}
	waitGW := &WaitGroupWrapper{}

	reporter := NewReporter("../../../wgliang/logcool", "foo", "foo", "baz")
	strategyCopyCheck := &StrategyCopyCheck{}
	strategyCountCode := &StrategyCountCode{}
	strategyCyclo := &StrategyCyclo{}
	strategyDeadCode := &StrategyDeadCode{}
	strategyDependGraph := &StrategyDependGraph{}
	strategyDepth := &StrategyDepth{}
	strategyImportPackages := &StrategyImportPackages{}
	strategyInterfacer := &StrategyInterfacer{}
	strategySimpleCode := &StrategySimpleCode{}
	strategySpellCheck := &StrategySpellCheck{}
	strategyUnitTest := &StrategyUnitTest{}

	if err := inject.Populate(
		reporter,
		synchronizer,
		strategyCopyCheck,
		strategyCountCode,
		strategyCyclo,
		strategyDeadCode,
		strategyDependGraph,
		strategyDepth,
		strategyImportPackages,
		strategyInterfacer,
		strategySimpleCode,
		strategySpellCheck,
		strategyUnitTest,
		syncRW,
		waitGW,
	); err != nil {
		log.Fatal(err)
	}

	reporter.AddLinters(strategyCopyCheck, strategyCountCode, strategyCyclo, strategyDeadCode, strategyDependGraph,
		strategyDepth, strategyImportPackages, strategyInterfacer, strategySimpleCode, strategySpellCheck, strategyUnitTest)

	go processbar.LinterProcessBar(synchronizer.LintersProcessChans, synchronizer.LintersFinishedSignal)

	if err := reporter.Report(); err != nil {
		glog.Errorln(err)
	}
}
