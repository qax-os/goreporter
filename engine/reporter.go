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
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/golang/glog"
	"github.com/json-iterator/go"

	"github.com/360EntSecGroup-Skylar/goreporter/utils"
)

type Synchronizer struct {
	SyncRW                *sync.RWMutex     `inject:""`
	WaitGW                *WaitGroupWrapper `inject:""`
	LintersProcessChans   chan int64        `json:"-"`
	LintersFinishedSignal chan string       `json:"-"`
}

// Reporter is the top struct of GoReporter.
type Reporter struct {
	Project   string            `json:"project"`
	Score     int               `json:"score"`
	Grade     int               `json:"grade"`
	Metrics   map[string]Metric `json:"metrics"`
	Issues    int               `json:"issues"`
	TimeStamp string            `json:"time_stamp"`
	Linters   []StrategyLinter
	Sync      *Synchronizer `inject:"" json:"-"`

	ProjectPath    string `json:"-"`
	ReportPath     string `json:"-"`
	HtmlTemplate   string `json:"-"`
	ReportFormat   string `json:"-"`
	ExceptPackages string `json:"-"`

	StartTime time.Time
}

// WaitGroupWrapper is a struct that as a waiter for all linetr-tasks.And it
// encapsulates Sync.WaitGroup that can be call as a interface.
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

type StrategyParameter struct {
	AllDirs, UnitTestDirs       map[string]string
	ProjectPath, ExceptPackages string
}

// Report is a important function of goreporter, it will run all linters and rebuild
// metrics data in a golang project. And all linters' result will be as one metric
// data for Reporter.
func (r *Reporter) Report() error {
	defer r.Close()
	glog.Infoln("start code quality assessment...")

	r.Project = utils.PackageAbsPath(r.ProjectPath)

	// All directory that has _test.go files will be add into.
	dirsUnitTest, err := utils.DirList(r.ProjectPath, "_test.go", r.ExceptPackages)
	if err != nil {
		return err
	}

	// All directory that has .go files will be add into.
	dirsAll, err := utils.DirList(r.ProjectPath, ".go", r.ExceptPackages)
	if err != nil {
		return err
	}

	params := StrategyParameter{
		AllDirs:      dirsAll,
		UnitTestDirs: dirsUnitTest,
		ProjectPath:  r.ProjectPath,
	}

	for _, linter := range r.Linters {
		r.compute(linter, params)
	}

	r.TimeStamp = time.Now().Format("2006-01-02-15-04-05")

	// ensure peocessbar quit.
	r.Sync.LintersProcessChans <- 100
	glog.Infoln("finished code quality assessment...")
	return nil
}

func (r *Reporter) compute(strategy StrategyLinter, params StrategyParameter) {
	glog.Infof("running %s...", strategy.GetName())

	summaries := strategy.Compute(params)

	r.Metrics[strategy.GetName()+"Tips"] = Metric{
		Name:        strategy.GetName(),
		Description: strategy.GetDescription(),
		Weight:      strategy.GetWeight(),
		Summaries:   summaries.Summaries,
		Percentage:  strategy.Percentage(summaries),
	}

	r.Sync.LintersFinishedSignal <- fmt.Sprintf("Linter:%s over,time consuming %vs", strategy.GetName(), time.Since(r.StartTime).Seconds())
	glog.Infof("%s over!", strategy.GetName())
}

func (r *Reporter) Render() (err error) {
	switch r.ReportFormat {
	case "json":
		err = r.toJson()
	case "text":
		err = r.toText()
	default:
		glog.Infof(fmt.Sprintf("Generating HTML report,time consuming %vs", time.Since(r.StartTime).Seconds()))
		err = r.toHtml()
		if err != nil {
			glog.Infoln("Json2Html error:", err)
			return
		}
	}
	return
}

// toJson will marshal struct Reporter into json and
// return a []byte data.
func (r *Reporter) toJson() (err error) {
	glog.Infof(fmt.Sprintf("Generating json report,time consuming %vs", time.Since(r.StartTime).Seconds()))
	jsonReport, err := jsoniter.Marshal(r)
	if err != nil {
		return
	}

	saveAbsPath := utils.AbsPath(r.ReportPath)
	projectName := utils.ProjectName(r.ProjectPath)

	jsonpath := projectName + "-" + r.TimeStamp + ".json"
	if saveAbsPath != "" && saveAbsPath != r.ReportPath {
		jsonpath = strings.Replace(saveAbsPath+string(filepath.Separator)+projectName+"-"+r.TimeStamp+".json", string(filepath.Separator)+string(filepath.Separator), string(filepath.Separator), -1)
	}
	if err = ioutil.WriteFile(jsonpath, jsonReport, 0666); err != nil {
		return
	}
	glog.Info("Json report saved in:", jsonpath)
	return
}

func (r *Reporter) toText() (err error) {
	glog.Infof(fmt.Sprintf("Generating text report,time consuming %vs", time.Since(r.StartTime).Seconds()))
	color.Magenta(
		headerTpl,
		r.Project,
		int(r.GetFinalScore()),
		r.Grade,
		r.TimeStamp,
		r.Issues,
	)
	for _, metric := range r.Metrics {
		if metric.Name == "DependGraph" || 0 == len(metric.Summaries) {
			continue
		}
		color.Cyan(metricsHeaderTpl, metric.Name, metric.Description)
		for _, summary := range metric.Summaries {
			color.Blue(summaryHeaderTpl, summary.Name, summary.Description)
			for _, errorInfo := range summary.Errors {
				color.Red(errorInfoTpl, errorInfo.ErrorString, errorInfo.LineNumber)
			}
		}
	}
	return
}

// toHtml will rebuild the reporter's json data into html data in template.
// It will parse json data and organize the data structure.
func (r *Reporter) toHtml() (err error) {
	glog.Infof(fmt.Sprintf("Generating json report,time consuming %vs", time.Since(r.StartTime).Seconds()))
	jsonReport, err := jsoniter.Marshal(r)
	if err != nil {
		return
	}
	if jsonReport == nil {
		return errors.New("json is null")
	}

	var htmlData HtmlData
	issues := 0

	htmlData.Project = r.Project
	htmlData.Score = int(r.GetFinalScore())
	// convert all linter's data.
	htmlData.converterCodeTest(*r)
	htmlData.converterCodeSmell(*r)
	htmlData.converterCodeOptimization(*r)
	htmlData.converterCodeStyle(*r)
	htmlData.converterCodeCount(*r)
	htmlData.converterDependGraph(*r)

	noTestPackages := make([]string, 0)
	importPackages := r.Metrics["ImportPackagesTips"].Summaries
	unitTestPackages := r.Metrics["UnitTestTips"].Summaries
	for packageName, _ := range importPackages {
		if _, ok := unitTestPackages[packageName]; !ok {
			noTestPackages = append(noTestPackages, packageName)
		}
	}
	htmlData.IssuesNum = issues
	htmlData.Date = r.TimeStamp

	SaveAsHtml(htmlData, r.ProjectPath, r.ReportPath, r.TimeStamp, r.HtmlTemplate)

	return
}

func (r *Reporter) GetFinalScore() (score float64) {
	for _, metric := range r.Metrics {
		score = score + metric.Percentage*metric.Weight
	}
	return
}

func NewReporter(projectPath, reportPath, reportFormat, htmlTemplate string) *Reporter {
	return &Reporter{
		StartTime:    time.Now(),
		Metrics:      make(map[string]Metric, 0),
		ProjectPath:  projectPath,
		ReportPath:   reportPath,
		ReportFormat: reportFormat,
		HtmlTemplate: htmlTemplate,
	}
}

func (r *Reporter) AddLinters(strategies ...StrategyLinter) {
	r.Linters = append(r.Linters, strategies...)
}

func (r *Reporter) Close() {
	close(r.Sync.LintersFinishedSignal)
	close(r.Sync.LintersProcessChans)
}

// Error contains the line number and the reason for
// an error output from a command
type Error struct {
	LineNumber  int    `json:"line_number"`
	ErrorString string `json:"error_string"`
}

// FileSummary contains the filename, location of the file
// on GitHub, and all of the errors related to the file
type Summary struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Errors      []Error `json:"errors"`
	SumCover    float64
	CountCover  int
	Avg         float64
}

type Summaries struct {
	Summaries map[string]Summary
	sync.RWMutex
}

func NewSummaries() *Summaries {
	return &Summaries{Summaries: make(map[string]Summary, 0)}
}

// Metric as template of report and will save all linters result
// data.But may have some difference in different linter.
type Metric struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Summaries   map[string]Summary `json:"summaries"`
	Weight      float64            `json:"weight"`
	Percentage  float64            `json:"percentage"`
	Error       string             `json:"error"`
}
