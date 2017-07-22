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
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/360EntSecGroup-Skylar/goreporter/linters/copycheck"
	"github.com/360EntSecGroup-Skylar/goreporter/linters/countcode"
	"github.com/360EntSecGroup-Skylar/goreporter/linters/cyclo"
	"github.com/360EntSecGroup-Skylar/goreporter/linters/deadcode"
	"github.com/360EntSecGroup-Skylar/goreporter/linters/depend"
	"github.com/360EntSecGroup-Skylar/goreporter/linters/simplecode"
	"github.com/360EntSecGroup-Skylar/goreporter/linters/spellcheck"
	"github.com/360EntSecGroup-Skylar/goreporter/linters/unittest"
	"github.com/golang/glog"
)

// WaitGroupWrapper is a struct that as a waiter for all linetr-tasks.And it
// encapsulates sync.WaitGroup that can be call as a interface.
type WaitGroupWrapper struct {
	sync.WaitGroup
}

// Wrap implements a interface that run the function cd as a goroutine.And it
// encapsulates Add(1) and Done() operation.You can just think go cd() but not
// worry about synchronization and security issues.
func (w *WaitGroupWrapper) Wrap(cb func()) {
	w.Add(1)
	go func() {
		cb()
		w.Done()
	}()
}

// NewReporter will initialize a Reporter struct and return address of the struct
// which is safe for use.
func NewReporter() *Reporter {
	return &Reporter{
		Metrics: make(map[string]Metric, 0),
		syncRW:  new(sync.RWMutex),
	}
}

// Engine is a important function of goreporter, it will run all linters and rebuild
// metrics data in a golang prohject. And all linters' result will be as one metric
// data for Reporter.
func (r *Reporter) Engine(projectPath string, exceptPackages string, lintersProcessChans chan int64, lintersFinishedSignal chan string, start time.Time) {
	glog.Infoln("start code quality assessment...")
	wg := &WaitGroupWrapper{}
	lintersFunction := make(map[string]func(), 9)
	// All directory that has _test.go files will be add into.
	dirsUnitTest, err := DirList(projectPath, "_test.go", exceptPackages)
	if err != nil {
		glog.Errorln(err)
	}
	r.syncRW.Lock()
	r.Project = PackageAbsPath(projectPath)
	r.syncRW.Unlock()

	// linterFunction:unitTestF,Run all valid TEST in your golang package.And will measure
	// from both coverage and time-consuming
	lintersFunction["unitTestF"] = func() {
		glog.Infoln("running unit test...")

		metricUnitTest := Metric{
			Name:        "UnitTest",
			Description: "Run all valid TEST in your golang package.And will measure from both coverage and time-consuming.",
			Weight:      0.4,
		}

		packagesTestDetail := struct {
			Values map[string]Summary
			mux    *sync.RWMutex
		}{make(map[string]Summary, 0), new(sync.RWMutex)}

		sumCover := 0.0
		countCover := 0
		sumProcessNumber := int64(30)
		processUnit := getProcessUnit(sumProcessNumber, len(dirsUnitTest))
		var pkg sync.WaitGroup
		for pkgName, pkgPath := range dirsUnitTest {
			pkg.Add(1)
			go func(pkgName, pkgPath string) {
				unitTestRes, _ := unittest.UnitTest("./" + pkgPath)
				var packageTest PackageTest
				if len(unitTestRes) >= 5 {
					if unitTestRes[0] == "ok" {
						packageTest.IsPass = true
					} else {
						packageTest.IsPass = false
					}
					timeLen := len(unitTestRes[2])
					if timeLen > 1 {
						time, err := strconv.ParseFloat(unitTestRes[2][:(timeLen-1)], 64)
						if err == nil {
							packageTest.Time = time
						} else {
							glog.Errorln(err)
						}
					}
					packageTest.Coverage = unitTestRes[4]

					coverLen := len(unitTestRes[4])
					if coverLen > 1 {
						coverFloat, _ := strconv.ParseFloat(unitTestRes[4][:(coverLen-1)], 64)
						sumCover = sumCover + coverFloat
						countCover = countCover + 1
					} else {
						countCover = countCover + 1
					}
				} else {
					packageTest.Coverage = "0%"
					countCover = countCover + 1
				}
				jsonStringPackageTest, err := json.Marshal(packageTest)
				if err != nil {
					glog.Errorln(err)
				}
				summarie := Summary{
					Name:        pkgName,
					Description: string(jsonStringPackageTest),
				}
				packagesTestDetail.mux.Lock()
				packagesTestDetail.Values[pkgName] = summarie
				packagesTestDetail.mux.Unlock()

				pkg.Done()
			}(pkgName, pkgPath)
			if sumProcessNumber > 0 {
				lintersProcessChans <- processUnit
				sumProcessNumber = sumProcessNumber - processUnit
			}
		}

		pkg.Wait()

		packagesTestDetail.mux.Lock()
		metricUnitTest.Summaries = packagesTestDetail.Values
		packagesTestDetail.mux.Unlock()
		if countCover == 0 {
			metricUnitTest.Percentage = 0
		} else {
			metricUnitTest.Percentage = sumCover / float64(countCover)
		}

		r.syncRW.Lock()
		r.Metrics["UnitTestTips"] = metricUnitTest
		r.syncRW.Unlock()
		lintersProcessChans <- sumProcessNumber
		lintersFinishedSignal <- fmt.Sprintf("Linter:UnitTest over,time consuming %vs", time.Now().Sub(start).Seconds())
		glog.Infoln("unit test over!")
	}
	// All directory that has .go files will be add into.
	dirsAll, err := DirList(projectPath, ".go", exceptPackages)
	if err != nil {
		glog.Errorln(err)
	}

	// linterFunnction:cycloF,Computing all [.go] file's cyclo,and as an important
	// indicator of the quality of the code.
	lintersFunction["cycloF"] = func() {
		glog.Infoln("computing cyclo...")

		metricCyclo := Metric{
			Name:        "Cyclo",
			Description: "Computing all [.go] file's cyclo,and as an important indicator of the quality of the code.",
			Weight:      0.2,
		}

		summaries := make(map[string]Summary, 0)
		sumAverageCyclo := 0.0
		sumProcessNumber := int64(10)
		processUnit := getProcessUnit(sumProcessNumber, len(dirsAll))
		var compBigThan15 int
		for pkgName, pkgPath := range dirsAll {
			summary := Summary{
				Name: pkgName,
			}
			summary.Errors = make([]Error, 0)
			errors := make([]Error, 0)
			cyclo, avg := cyclo.Cyclo(pkgPath)
			avgfloat, _ := strconv.ParseFloat(avg, 64)
			sumAverageCyclo = sumAverageCyclo + avgfloat
			for _, val := range cyclo {
				cyclovalues := strings.Split(val, " ")
				if len(cyclovalues) == 4 {
					comp, _ := strconv.Atoi(cyclovalues[0])
					erroru := Error{
						LineNumber:  comp,
						ErrorString: AbsPath(cyclovalues[3]),
					}
					if comp >= 15 {
						compBigThan15 = compBigThan15 + 1
					}
					errors = append(errors, erroru)
				}
			}
			summary.Errors = errors
			summary.Description = avg
			summaries[pkgName] = summary
			if sumProcessNumber > 0 {
				lintersProcessChans <- processUnit
				sumProcessNumber = sumProcessNumber - processUnit
			}
		}

		metricCyclo.Summaries = summaries
		metricCyclo.Percentage = countPercentage(compBigThan15 + int(sumAverageCyclo/float64(len(dirsAll))) - 1)
		r.syncRW.Lock()
		r.Issues = r.Issues + len(summaries)
		r.Metrics["CycloTips"] = metricCyclo
		r.syncRW.Unlock()
		lintersProcessChans <- sumProcessNumber
		lintersFinishedSignal <- fmt.Sprintf("Linter:Cyclo over,time consuming %vs", time.Now().Sub(start).Seconds())
		glog.Infoln("comput cyclo done!")
	}
	// linterfunction:simpleCodeF,all golang code hints that can be optimized
	// and give suggestions for changes.
	lintersFunction["simpleCodeF"] = func() {
		glog.Infoln("simpling code...")

		metricSimple := Metric{
			Name:        "Simple",
			Description: "All golang code hints that can be optimized and give suggestions for changes.",
			Weight:      0.1,
		}
		summaries := make(map[string]Summary, 0)

		simples := simplecode.Simple(dirsAll)
		sumProcessNumber := int64(10)
		processUnit := getProcessUnit(sumProcessNumber, len(simples))
		for _, simpleTip := range simples {
			simpleTips := strings.Split(simpleTip, ":")
			if len(simpleTips) == 4 {
				packageName := packageNameFromGoPath(simpleTips[0])
				line, _ := strconv.Atoi(simpleTips[1])
				erroru := Error{
					LineNumber:  line,
					ErrorString: AbsPath(simpleTips[0]) + ":" + strings.Join(simpleTips[1:], ":"),
				}
				if summarie, ok := summaries[packageName]; ok {
					summarie.Errors = append(summarie.Errors, erroru)
					summaries[packageName] = summarie
				} else {
					summarie := Summary{
						Name:   packageName,
						Errors: make([]Error, 0),
					}
					summarie.Errors = append(summarie.Errors, erroru)
					summaries[packageName] = summarie
				}

			}
			if sumProcessNumber > 0 {
				lintersProcessChans <- processUnit
				sumProcessNumber = sumProcessNumber - processUnit
			}
		}
		metricSimple.Summaries = summaries
		metricSimple.Percentage = countPercentage(len(summaries))
		r.syncRW.Lock()
		r.Issues = r.Issues + len(summaries)
		r.Metrics["SimpleTips"] = metricSimple
		r.syncRW.Unlock()
		lintersProcessChans <- sumProcessNumber
		lintersFinishedSignal <- fmt.Sprintf("Linter:Simple over,time consuming %vs", time.Now().Sub(start).Seconds())
		glog.Infoln("simple code done!")
	}

	// linterFunction:copycode,query all duplicate code in the project and give
	// duplicate code locations and rows.
	lintersFunction["copyCheckF"] = func() {
		glog.Infoln("checking copy code...")
		metricCopyCode := Metric{
			Name:        "CopyCode",
			Description: "Query all duplicate code in the project and give duplicate code locations and rows.",
			Weight:      0.1,
		}

		summaries := make(map[string]Summary, 0)
		copyCodeList := copycheck.CopyCheck(projectPath, "_test.go")
		sumProcessNumber := int64(10)
		processUnit := getProcessUnit(sumProcessNumber, len(copyCodeList))
		for i := 0; i < len(copyCodeList); i++ {
			summary := Summary{
				Errors: make([]Error, 0),
			}
			for j := 0; j < len(copyCodeList[i]); j++ {
				var line int
				values := strings.Split(copyCodeList[i][j], ":")
				if len(values) > 1 {
					lines := strings.Split(strings.TrimSpace(values[1]), ",")
					if len(lines) == 2 {
						lineright, _ := strconv.Atoi(lines[1])
						lineleft, _ := strconv.Atoi(lines[0])
						if lineright-lineleft >= 0 {
							line = lineright - lineleft + 1
						}
					}
					values[0] = AbsPath(values[0])
				}

				summary.Errors = append(summary.Errors, Error{LineNumber: line, ErrorString: strings.Join(values, ":")})
			}
			summary.Name = strconv.Itoa(len(summary.Errors))
			summaries[string(i)] = summary
			if sumProcessNumber > 0 {
				lintersProcessChans <- processUnit
				sumProcessNumber = sumProcessNumber - processUnit
			}
		}

		metricCopyCode.Summaries = summaries
		metricCopyCode.Percentage = countPercentage(len(summaries))
		r.syncRW.Lock()
		r.Issues = r.Issues + len(summaries)
		r.Metrics["CopyCodeTips"] = metricCopyCode
		r.syncRW.Unlock()
		lintersProcessChans <- sumProcessNumber
		lintersFinishedSignal <- fmt.Sprintf("Linter:CopyCode over,time consuming %vs", time.Now().Sub(start).Seconds())
		glog.Infoln("checked copy code!")
	}
	// linterFunction:deadCodeF,all useless code, or never obsolete obsolete code.
	lintersFunction["deadCodeF"] = func() {
		glog.Infoln("checking dead code...")

		metricDeadCode := Metric{
			Name:        "DeadCode",
			Description: "All useless code, or never obsolete obsolete code.",
			Weight:      0.1,
		}
		summaries := make(map[string]Summary, 0)

		deadcode := deadcode.DeadCode(projectPath)
		sumProcessNumber := int64(10)
		processUnit := getProcessUnit(sumProcessNumber, len(deadcode))
		for _, simpleTip := range deadcode {
			deadCodeTips := strings.Split(simpleTip, ":")
			if len(deadCodeTips) == 4 {
				packageName := packageNameFromGoPath(deadCodeTips[0])
				line, _ := strconv.Atoi(deadCodeTips[1])
				erroru := Error{
					LineNumber:  line,
					ErrorString: AbsPath(deadCodeTips[0]) + ":" + strings.Join(deadCodeTips[1:], ":"),
				}
				if summarie, ok := summaries[packageName]; ok {
					summarie.Errors = append(summarie.Errors, erroru)
					summaries[packageName] = summarie
				} else {
					summarie := Summary{
						Name:   PackageAbsPathExceptSuffix(deadCodeTips[0]),
						Errors: make([]Error, 0),
					}
					summarie.Errors = append(summarie.Errors, erroru)
					summaries[packageName] = summarie
				}
			}
			if sumProcessNumber > 0 {
				lintersProcessChans <- processUnit
				sumProcessNumber = sumProcessNumber - processUnit
			}
		}
		metricDeadCode.Summaries = summaries
		metricDeadCode.Percentage = countPercentage(len(summaries))
		r.syncRW.Lock()
		r.Issues = r.Issues + len(summaries)
		r.Metrics["DeadCodeTips"] = metricDeadCode
		r.syncRW.Unlock()
		lintersProcessChans <- sumProcessNumber
		lintersFinishedSignal <- fmt.Sprintf("Linter:SwadCode over,time consuming %vs", time.Now().Sub(start).Seconds())
		glog.Infoln("check dead code done.")
	}
	// linterFunction:spellCheckF,check the project variables, functions,
	// etc. naming spelling is wrong.
	lintersFunction["spellCheckF"] = func() {
		glog.Infoln("checking spell error...")

		metricSpellTips := Metric{
			Name:        "SpellCheck",
			Description: "Check the project variables, functions, etc. naming spelling is wrong.",
			Weight:      0.1,
		}
		summaries := make(map[string]Summary, 0)

		spelltips := spellcheck.SpellCheck(projectPath, exceptPackages)
		sumProcessNumber := int64(10)
		processUnit := getProcessUnit(sumProcessNumber, len(spelltips))
		for _, simpleTip := range spelltips {
			simpleTips := strings.Split(simpleTip, ":")
			if len(simpleTips) == 4 {
				packageName := packageNameFromGoPath(simpleTips[0])
				line, _ := strconv.Atoi(simpleTips[1])
				erroru := Error{
					LineNumber:  line,
					ErrorString: AbsPath(simpleTips[0]) + ":" + strings.Join(simpleTips[1:], ":"),
				}
				if summarie, ok := summaries[packageName]; ok {
					summarie.Errors = append(summarie.Errors, erroru)
					summaries[packageName] = summarie
				} else {
					summarie := Summary{
						Name:   PackageAbsPathExceptSuffix(simpleTips[0]),
						Errors: make([]Error, 0),
					}
					summarie.Errors = append(summarie.Errors, erroru)
					summaries[packageName] = summarie
				}
			}
			if sumProcessNumber > 0 {
				lintersProcessChans <- processUnit
				sumProcessNumber = sumProcessNumber - processUnit
			}
		}
		metricSpellTips.Summaries = summaries
		metricSpellTips.Percentage = countPercentage(len(summaries))
		r.syncRW.Lock()
		r.Issues = r.Issues + len(summaries)
		r.Metrics["SpellCheckTips"] = metricSpellTips
		r.syncRW.Unlock()
		lintersProcessChans <- sumProcessNumber
		lintersFinishedSignal <- fmt.Sprintf("Linter:Spell over,time consuming %vs", time.Now().Sub(start).Seconds())
		glog.Infoln("checked spell error")
	}
	// linterFunction:importPackagesF,The project contains all the package lists.
	lintersFunction["ImportPackagesF"] = func() {
		glog.Infoln("getting import packages...")
		metricImportPackageTips := Metric{
			Name:        "ImportPackages",
			Description: "Check the project variables, functions, etc. naming spelling is wrong.",
			Weight:      0,
			Summaries:   make(map[string]Summary, 0),
		}
		summaries := make(map[string]Summary, 0)
		importPkgs := unittest.GoListWithImportPackages(projectPath)
		for i := 0; i < len(importPkgs); i++ {
			summaries[importPkgs[i]] = Summary{Name: importPkgs[i]}
		}
		metricImportPackageTips.Summaries = summaries
		metricImportPackageTips.Percentage = countPercentage(len(summaries))
		r.syncRW.Lock()
		r.Metrics["ImportPackagesTips"] = metricImportPackageTips
		r.syncRW.Unlock()
		lintersProcessChans <- int64(5)
		lintersFinishedSignal <- fmt.Sprintf("Linter:ImportPackages over,time consuming %vs", time.Now().Sub(start).Seconds())
		glog.Infoln("import packages done.")
	}

	// linterFunction:countCodeF,Count go files and go code lines of project.
	lintersFunction["CountCodeF"] = func() {
		glog.Infoln("countting code...")
		metricCountCodeTips := Metric{
			Name:        "CountCode",
			Description: "Count lines and files of go project.",
			Weight:      0,
			Summaries:   make(map[string]Summary, 0),
		}
		summaries := make(map[string]Summary, 0)
		fileCount, codeLines, commentLines, totalLines := countcode.CountCode(projectPath, exceptPackages)
		summaries["FileCount"] = Summary{
			Name:        "FileCount",
			Description: strconv.Itoa(fileCount),
		}
		summaries["CodeLines"] = Summary{
			Name:        "CodeLines",
			Description: strconv.Itoa(codeLines),
		}
		summaries["CommentLines"] = Summary{
			Name:        "CommentLines",
			Description: strconv.Itoa(commentLines),
		}
		summaries["TotalLines"] = Summary{
			Name:        "TotalLines",
			Description: strconv.Itoa(totalLines),
		}
		metricCountCodeTips.Summaries = summaries
		metricCountCodeTips.Percentage = 0
		r.syncRW.Lock()
		r.Metrics["CountCodeTips"] = metricCountCodeTips
		r.syncRW.Unlock()
		lintersProcessChans <- int64(5)
		lintersFinishedSignal <- fmt.Sprintf("Linter:CountCode over,time consuming %vs", time.Now().Sub(start).Seconds())
		glog.Infoln("count code done.")
	}

	// linterFunction:dependGraphF,The dependency graph for all packages in the
	// project helps you optimize the project architecture.
	lintersFunction["dependGraphF"] = func() {
		glog.Infoln("creating depend graph...")
		metricDependGraphTips := Metric{
			Name:        "DependGraph",
			Description: "The dependency graph for all packages in the project helps you optimize the project architecture.",
			Weight:      0,
		}
		summaries := make(map[string]Summary, 0)

		graph := depend.Depend(projectPath, exceptPackages)
		summaries["graph"] = Summary{
			Name:        "graph",
			Description: graph,
		}
		metricDependGraphTips.Summaries = summaries
		metricDependGraphTips.Percentage = countPercentage(len(summaries))
		r.syncRW.Lock()
		r.Issues = r.Issues + len(summaries)
		r.Metrics["DependGraphTips"] = metricDependGraphTips
		r.syncRW.Unlock()
		lintersProcessChans <- int64(10)
		lintersFinishedSignal <- fmt.Sprintf("Linter:DependGraph over,time consuming %vs", time.Now().Sub(start).Seconds())
		glog.Infoln("created depend graph")
	}
	r.TimeStamp = time.Now().Format("2006-01-02 15:04:05")
	// run all linters.
	for _, funcRun := range lintersFunction {
		wg.Wrap(funcRun)
	}

	wg.Wait()
	// ensure peocessbar quit.
	lintersProcessChans <- 100
	glog.Infoln("finished code quality assessment...")
}

// FormateReport2Json will marshal struct Reporter into json and
// return a []byte data.
func (r *Reporter) FormateReport2Json() []byte {
	report, err := json.Marshal(r)
	if err != nil {
		glog.Errorln("json err:", err)
	}

	return report
}

// countPercentage will count all linters' percentage.And rule is
//
//    +--------------------------------------------------+
//    |   issues    |               score                |
//    +==================================================+
//    | 5           | 100-issues*2                            |
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
func countPercentage(issues int) float64 {
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

func getProcessUnit(sumProcessNumber int64, number int) int64 {
	if number == 0 {
		return sumProcessNumber
	} else if sumProcessNumber/int64(number) <= 0 {
		return int64(1)
	} else {
		return sumProcessNumber / int64(number)
	}
}
