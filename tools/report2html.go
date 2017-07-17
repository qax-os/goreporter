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

package tools

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/goreporter/engine"
	"github.com/golang/glog"
)

// Json2Html will rebuild the reporter's json data into html data in template.
// It will parse json data and organize the data structure.
func Json2Html(jsonData []byte) (HtmlData, error) {
	var (
		structData engine.Reporter
		htmlData   HtmlData
		issues     int
	)
	issues = 0

	if jsonData == nil {
		return htmlData, errors.New("json is null")
	}
	json.Unmarshal(jsonData, &structData)

	htmlData.Project = structData.Project
	var score float64
	for _, metric := range structData.Metrics {
		score = score + metric.Percentage*metric.Weight
	}
	htmlData.Score = int(score)
	// convert test result
	testHtmlRes := make([]Test, 0)
	if result, ok := structData.Metrics["UnitTestTips"]; ok {
		for pkgName, testRes := range result.Summaries {
			var packageUnitTestResult PackageTest
			json.Unmarshal([]byte(testRes.Description), &packageUnitTestResult)
			srcLastIndex := strings.LastIndex(pkgName, htmlData.Project)
			if srcLastIndex < len(pkgName) && srcLastIndex >= 0 {
				test := Test{
					Path: pkgName[srcLastIndex:],
					Time: packageUnitTestResult.Time,
				}
				if len(packageUnitTestResult.Coverage) > 1 {
					test.Cover, _ = strconv.ParseFloat(packageUnitTestResult.Coverage[:(len(packageUnitTestResult.Coverage)-1)], 64)
				}
				if packageUnitTestResult.IsPass {
					test.Result = 1
				}
				testHtmlRes = append(testHtmlRes, test)
			}
		}
		htmlData.TestSummaryCoverAvg = fmt.Sprintf("%0.1f", result.Percentage)
	}

	stringTestJson, err := json.Marshal(testHtmlRes)
	if err != nil {
		glog.Errorln(err)
	}
	htmlData.Tests = string(stringTestJson)

	// convert cyclo result
	cycloHtmlRes := make([]Cyclo, 0)
	if result, ok := structData.Metrics["CycloTips"]; ok {
		for pkgName, summary := range result.Summaries {
			cycloTips := summary.Errors
			cyclo := Cyclo{
				Pkg:  pkgName,
				Size: len(cycloTips),
			}
			var infos []CycloInfo
			for i := 0; i < len(cycloTips); i++ {
				cycloTip := strings.Split(cycloTips[i].ErrorString, ":")
				if len(cycloTip) == 3 {
					cycloInfo := CycloInfo{
						Comp: cycloTips[i].LineNumber,
						Info: strings.Join(cycloTip[0:], ":"),
					}
					if cycloTips[i].LineNumber > 15 {
						htmlData.CycloBigThan15 = htmlData.CycloBigThan15 + 1
						issues = issues + 1
					}
					infos = append(infos, cycloInfo)
				}
			}
			cyclo.Info = infos
			cycloHtmlRes = append(cycloHtmlRes, cyclo)
		}
	}

	stringCycloJson, err := json.Marshal(cycloHtmlRes)
	if err != nil {
		glog.Errorln(err)
	}
	htmlData.Cyclos = string(stringCycloJson)

	// convert simple code result
	simpleHtmlRes := make([]Simple, 0)
	if result, ok := structData.Metrics["SimpleTips"]; ok {
		for _, summary := range result.Summaries {
			simpleCodeTips := summary.Errors

			for i := 0; i < len(simpleCodeTips); i++ {
				simpleCodeTip := strings.Split(simpleCodeTips[i].ErrorString, ":")
				if len(simpleCodeTip) == 4 {
					simpecode := Simple{
						Path: strings.Join(simpleCodeTip[0:3], ":"),
						Info: simpleCodeTip[3],
					}
					simpleHtmlRes = append(simpleHtmlRes, simpecode)
					issues = issues + 1
				}
			}
		}
	}

	stringSimpleJson, err := json.Marshal(simpleHtmlRes)
	if err != nil {
		glog.Errorln(err)
	}
	htmlData.Simples = string(stringSimpleJson)
	htmlData.SimpleIssues = len(simpleHtmlRes)

	// convert spell code result
	spellHtmlRes := make([]Spell, 0)
	if result, ok := structData.Metrics["SpellCheckTips"]; ok {
		for _, summary := range result.Summaries {
			spellCodeTips := summary.Errors

			for i := 0; i < len(spellCodeTips); i++ {
				spellCodeTip := strings.Split(spellCodeTips[i].ErrorString, ":")
				if len(spellCodeTip) == 4 {
					spellcode := Spell{
						Path: strings.Join(spellCodeTip[0:3], ":"),
						Info: spellCodeTip[3],
					}
					spellHtmlRes = append(spellHtmlRes, spellcode)
				}
			}
		}
	}

	stringSpellJson, err := json.Marshal(spellHtmlRes)
	if err != nil {
		glog.Errorln(err)
	}
	htmlData.Spells = string(stringSpellJson)

	// convert copy code result
	copyHtmlRes := make([]Copycode, 0)
	if result, ok := structData.Metrics["CopyCodeTips"]; ok {
		for _, copyResult := range result.Summaries {
			copyTips := copyResult.Errors
			var copyCodePathList []string
			for i := 0; i < len(copyTips); i++ {
				copyCodeTip := strings.Split(copyTips[i].ErrorString, ":")
				if len(copyCodeTip) == 2 {
					copyCodePathList = append(copyCodePathList, strings.Join(copyCodeTip[0:], ":"))
				}
			}
			copycode := Copycode{
				Files: strconv.Itoa(len(copyTips)),
				Path:  copyCodePathList,
			}
			copyHtmlRes = append(copyHtmlRes, copycode)
			issues = issues + 1
		}
	}

	stringCopyJson, err := json.Marshal(copyHtmlRes)
	if err != nil {
		glog.Errorln(err)
	}
	htmlData.Copycodes = string(stringCopyJson)

	// convert simple code result
	deadcodeHtmlRes := make([]Deadcode, 0)
	if result, ok := structData.Metrics["DeadCodeTips"]; ok {
		for _, deadCodeResult := range result.Summaries {
			deadCodeTips := deadCodeResult.Errors

			for i := 0; i < len(deadCodeTips); i++ {
				deadCodeTip := strings.Split(deadCodeTips[i].ErrorString, ":")
				if len(deadCodeTip) == 4 {
					deadcode := Deadcode{
						Path: strings.Join(deadCodeTip[0:3], ":"),
						Info: deadCodeTip[3],
					}
					deadcodeHtmlRes = append(deadcodeHtmlRes, deadcode)
					issues = issues + 1
				}
			}
		}
	}

	stringDeadCodeJson, err := json.Marshal(deadcodeHtmlRes)
	if err != nil {
		glog.Errorln(err)
	}
	htmlData.Deadcodes = string(stringDeadCodeJson)
	htmlData.DeadcodeIssues = len(deadcodeHtmlRes)

	// convert countline result
	if result, ok := structData.Metrics["CountCodeTips"]; ok {
		htmlData.FileCount, _ = strconv.Atoi(result.Summaries["FileCount"].Description)
		htmlData.CodeLines, _ = strconv.Atoi(result.Summaries["CodeLines"].Description)
		htmlData.CommentLines, _ = strconv.Atoi(result.Summaries["CommentLines"].Description)
		htmlData.TotalLines, _ = strconv.Atoi(result.Summaries["TotalLines"].Description)
	}

	// convert depend graph
	htmlData.DepGraph = template.HTML(structData.Metrics["DependGraphTips"].Summaries["graph"].Description)
	noTestPackages := make([]string, 0)
	importPackages := structData.Metrics["ImportPackagesTips"].Summaries
	unitTestPackages := structData.Metrics["UnitTestTips"].Summaries
	for packageName, _ := range importPackages {
		if _, ok := unitTestPackages[packageName]; !ok {
			noTestPackages = append(noTestPackages, packageName)
		}
	}
	stringNoTestJson, err := json.Marshal(noTestPackages)
	if err != nil {
		glog.Errorln(err)
	}
	htmlData.NoTests = string(stringNoTestJson)
	htmlData.Issues = issues
	htmlData.Date = structData.TimeStamp

	if len(importPackages) > 0 && len(noTestPackages) == 0 {
		htmlData.AveragePackageCover = float64(100)
	} else if len(importPackages) > 0 {
		htmlData.AveragePackageCover = float64(100 * (len(importPackages) - len(noTestPackages)) / len(importPackages))
	} else {
		htmlData.AveragePackageCover = float64(0)
	}
	return htmlData, nil
}

// SaveAsHtml is a function that save HtmlData as a html report.And will receive
// htmlData, projectPath, savePath and tpl parameters.
func SaveAsHtml(htmlData HtmlData, projectPath, savePath, timestamp, tpl string) {
	if tpl == "" {
		tpl = defaultTpl
	}

	t, err := template.New("goreporter").Parse(tpl)
	if err != nil {
		glog.Errorln(err)
	}

	var out bytes.Buffer
	err = t.Execute(&out, htmlData)
	if err != nil {
		glog.Errorln(err)
	}
	projectName := engine.ProjectName(projectPath)
	if savePath != "" {
		htmlpath := strings.Replace(savePath+string(filepath.Separator)+projectName+"-"+timestamp+".html", string(filepath.Separator)+string(filepath.Separator), string(filepath.Separator), -1)
		err = ioutil.WriteFile(htmlpath, out.Bytes(), 0666)
		if err != nil {
			glog.Errorln(err)
		} else {
			glog.Info("Html report was saved in:", htmlpath)
		}
	} else {
		htmlpath := projectName + "-" + timestamp + ".html"
		err = ioutil.WriteFile(htmlpath, out.Bytes(), 0666)
		if err != nil {
			glog.Errorln(err)
		} else {
			glog.Info("Html report was saved in:", htmlpath)
		}
	}
}
