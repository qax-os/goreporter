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
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/goreporter/engine"
	"github.com/golang/glog"
	"github.com/json-iterator/go"
)

var (
	issues int
)

// Json2Html will rebuild the reporter's json data into html data in template.
// It will parse json data and organize the data structure.
func Json2Html(jsonData []byte) (HtmlData, error) {
	var (
		structData engine.Reporter
		htmlData   HtmlData
	)
	issues = 0

	if jsonData == nil {
		return htmlData, errors.New("json is null")
	}
	jsoniter.Unmarshal(jsonData, &structData)

	htmlData.Project = structData.Project
	var score float64
	for _, metric := range structData.Metrics {
		score = score + metric.Percentage*metric.Weight
	}
	htmlData.Score = int(score)
	// convert all linter's data.
	htmlData.converterUnitTest(structData)
	htmlData.converterCopy(structData)
	htmlData.converterCyclo(structData)
	htmlData.converterDepth(structData)
	htmlData.converterInterfacer(structData)
	htmlData.converterSimple(structData)
	htmlData.converterSpell(structData)
	htmlData.converterCount(structData)
	htmlData.converterDead(structData)
	htmlData.converterDependGraph(structData)

	noTestPackages := make([]string, 0)
	importPackages := structData.Metrics["ImportPackagesTips"].Summaries
	unitTestPackages := structData.Metrics["UnitTestTips"].Summaries
	for packageName, _ := range importPackages {
		if _, ok := unitTestPackages[packageName]; !ok {
			noTestPackages = append(noTestPackages, packageName)
		}
	}
	stringNoTestJson, err := jsoniter.Marshal(noTestPackages)
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

// converterUnitTest provides function that convert unit test data into the
// format required in the html template.It will extract from the structData
// need to convert the data.The result will be saved in the hd's attributes.
func (hd *HtmlData) converterUnitTest(structData engine.Reporter) {
	testHtmlRes := make([]Test, 0)
	if result, ok := structData.Metrics["UnitTestTips"]; ok {
		for pkgName, testRes := range result.Summaries {
			var packageUnitTestResult PackageTest
			jsoniter.Unmarshal([]byte(testRes.Description), &packageUnitTestResult)
			srcLastIndex := strings.LastIndex(pkgName, hd.Project)
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
		hd.TestSummaryCoverAvg = fmt.Sprintf("%0.1f", result.Percentage)
	}

	stringTestJson, err := jsoniter.Marshal(testHtmlRes)
	if err != nil {
		glog.Errorln(err)
	}
	hd.Tests = string(stringTestJson)
}

// converterCyclo provides function that convert cyclo data into the
// format required in the html template.It will extract from the structData
// need to convert the data.The result will be saved in the hd's attributes.
func (hd *HtmlData) converterCyclo(structData engine.Reporter) {
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
				if len(cycloTip) >= 3 {
					cycloInfo := CycloInfo{
						Comp: cycloTips[i].LineNumber,
						Info: strings.Join(cycloTip[0:], ":"),
					}
					if cycloTips[i].LineNumber > 15 {
						hd.CycloBigThan15 = hd.CycloBigThan15 + 1
						issues = issues + 1
					}
					infos = append(infos, cycloInfo)
				}
			}
			cyclo.Info = infos
			cycloHtmlRes = append(cycloHtmlRes, cyclo)
		}
	}

	stringCycloJson, err := jsoniter.Marshal(cycloHtmlRes)
	if err != nil {
		glog.Errorln(err)
	}
	hd.Cyclos = string(stringCycloJson)
}

// converterCyclo provides function that convert cyclo data into the
// format required in the html template.It will extract from the structData
// need to convert the data.The result will be saved in the hd's attributes.
func (hd *HtmlData) converterDepth(structData engine.Reporter) {
	// convert depth result
	depthHtmlRes := make([]Depth, 0)
	if result, ok := structData.Metrics["DepthTips"]; ok {
		for pkgName, summary := range result.Summaries {
			depthTips := summary.Errors
			depth := Depth{
				Pkg:  pkgName,
				Size: len(depthTips),
			}
			var infos []DepthInfo
			for i := 0; i < len(depthTips); i++ {
				depthTip := strings.Split(depthTips[i].ErrorString, ":")
				if len(depthTip) >= 3 {
					depthInfo := DepthInfo{
						Comp: depthTips[i].LineNumber,
						Info: strings.Join(depthTip[0:], ":"),
					}
					if depthTips[i].LineNumber > 3 {
						hd.CycloBigThan15 = hd.CycloBigThan15 + 1
						issues = issues + 1
					}
					infos = append(infos, depthInfo)
				}
			}
			depth.Info = infos
			depthHtmlRes = append(depthHtmlRes, depth)
		}
	}

	stringDepthJson, err := jsoniter.Marshal(depthHtmlRes)
	if err != nil {
		glog.Errorln(err)
	}
	hd.Depths = string(stringDepthJson)
}

// converterSimple provides function that convert simplecode data into the
// format required in the html template.It will extract from the structData
// need to convert the data.The result will be saved in the hd's attributes.
func (hd *HtmlData) converterSimple(structData engine.Reporter) {
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
				} else if len(simpleCodeTip) == 5 {
					simpecode := Simple{
						Path: strings.Join(simpleCodeTip[0:4], ":"),
						Info: simpleCodeTip[4],
					}
					simpleHtmlRes = append(simpleHtmlRes, simpecode)
					issues = issues + 1
				}
			}
		}
	}

	stringSimpleJson, err := jsoniter.Marshal(simpleHtmlRes)
	if err != nil {
		glog.Errorln(err)
	}
	hd.Simples = string(stringSimpleJson)
	hd.SimpleIssues = len(simpleHtmlRes)
}

// converterInterfacer provides function that convert interfacer data into the
// format required in the html template.It will extract from the structData
// need to convert the data.The result will be saved in the hd's attributes.
func (hd *HtmlData) converterInterfacer(structData engine.Reporter) {
	interfacerHtmlRes := make([]Interfacer, 0)
	if result, ok := structData.Metrics["InterfacerTips"]; ok {
		for _, summary := range result.Summaries {
			interfacerCodeTips := summary.Errors

			for i := 0; i < len(interfacerCodeTips); i++ {
				interfacerCodeTip := strings.Split(interfacerCodeTips[i].ErrorString, ":")
				if len(interfacerCodeTip) == 4 {
					interfacer := Interfacer{
						Path: strings.Join(interfacerCodeTip[0:3], ":"),
						Info: interfacerCodeTip[3],
					}
					interfacerHtmlRes = append(interfacerHtmlRes, interfacer)
					issues = issues + 1
				} else if len(interfacerCodeTip) == 5 {
					interfacer := Interfacer{
						Path: strings.Join(interfacerCodeTip[0:4], ":"),
						Info: interfacerCodeTip[4],
					}
					interfacerHtmlRes = append(interfacerHtmlRes, interfacer)
					issues = issues + 1
				}
			}
		}
	}

	stringInterfacerJson, err := jsoniter.Marshal(interfacerHtmlRes)
	if err != nil {
		glog.Errorln(err)
	}
	hd.Interfacers = string(stringInterfacerJson)
}

// converterSpell provides function that convert spellcheck data into the
// format required in the html template.It will extract from the structData
// need to convert the data.The result will be saved in the hd's attributes.
func (hd *HtmlData) converterSpell(structData engine.Reporter) {
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
				} else if len(spellCodeTip) == 5 {
					spellcode := Spell{
						Path: strings.Join(spellCodeTip[0:4], ":"),
						Info: spellCodeTip[4],
					}
					spellHtmlRes = append(spellHtmlRes, spellcode)
				}
			}
		}
	}

	stringSpellJson, err := jsoniter.Marshal(spellHtmlRes)
	if err != nil {
		glog.Errorln(err)
	}
	hd.Spells = string(stringSpellJson)
}

// converterCopy provides function that convert copycode data into the
// format required in the html template.It will extract from the structData
// need to convert the data.The result will be saved in the hd's attributes.
func (hd *HtmlData) converterCopy(structData engine.Reporter) {
	copyHtmlRes := make([]Copycode, 0)
	if result, ok := structData.Metrics["CopyCodeTips"]; ok {
		for _, copyResult := range result.Summaries {
			copyTips := copyResult.Errors
			var copyCodePathList []string
			for i := 0; i < len(copyTips); i++ {
				copyCodeTip := strings.Split(copyTips[i].ErrorString, ":")
				if len(copyCodeTip) >= 2 {
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

	stringCopyJson, err := jsoniter.Marshal(copyHtmlRes)
	if err != nil {
		glog.Errorln(err)
	}
	hd.Copycodes = string(stringCopyJson)
}

// converterDead provides function that convert deadcode data into the
// format required in the html template.It will extract from the structData
// need to convert the data.The result will be saved in the hd's attributes.
func (hd *HtmlData) converterDead(structData engine.Reporter) {
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
				} else if len(deadCodeTip) == 5 {
					deadcode := Deadcode{
						Path: strings.Join(deadCodeTip[0:4], ":"),
						Info: deadCodeTip[4],
					}
					deadcodeHtmlRes = append(deadcodeHtmlRes, deadcode)
					issues = issues + 1
				}
			}
		}
	}

	stringDeadCodeJson, err := jsoniter.Marshal(deadcodeHtmlRes)
	if err != nil {
		glog.Errorln(err)
	}
	hd.Deadcodes = string(stringDeadCodeJson)
	hd.DeadcodeIssues = len(deadcodeHtmlRes)
}

// converterCount provides function that convert countcode data into the
// format required in the html template.It will extract from the structData
// need to convert the data.The result will be saved in the hd's attributes.
func (hd *HtmlData) converterCount(structData engine.Reporter) {
	if result, ok := structData.Metrics["CountCodeTips"]; ok {
		hd.FileCount, _ = strconv.Atoi(result.Summaries["FileCount"].Description)
		hd.CodeLines, _ = strconv.Atoi(result.Summaries["CodeLines"].Description)
		hd.CommentLines, _ = strconv.Atoi(result.Summaries["CommentLines"].Description)
		hd.TotalLines, _ = strconv.Atoi(result.Summaries["TotalLines"].Description)
	}
}

// converterDependGraph provides function that convert depend graph data into the
// format required in the html template.It will extract from the structData
// need to convert the data.The result will be saved in the hd's attributes.
func (hd *HtmlData) converterDependGraph(structData engine.Reporter) {
	hd.DepGraph = template.HTML(structData.Metrics["DependGraphTips"].Summaries["graph"].Description)
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

	var (
		out      bytes.Buffer
		htmlpath string
	)
	err = t.Execute(&out, htmlData)
	if err != nil {
		glog.Errorln(err)
	}
	projectName := engine.ProjectName(projectPath)
	if savePath != "" {
		htmlpath = strings.Replace(savePath+string(filepath.Separator)+projectName+"-"+timestamp+".html", string(filepath.Separator)+string(filepath.Separator), string(filepath.Separator), -1)
		err = ioutil.WriteFile(htmlpath, out.Bytes(), 0666)
		if err != nil {
			glog.Errorln(err)
		} else {
			glog.Info("Html report was saved in:", htmlpath)
		}
		absPath, err := filepath.Abs(htmlpath)
		if err != nil {
			log.Println(err)
		} else {
			displayReport(absPath)
		}

	} else {
		htmlpath = projectName + "-" + timestamp + ".html"
		err = ioutil.WriteFile(htmlpath, out.Bytes(), 0666)
		if err != nil {
			glog.Errorln(err)
		} else {
			glog.Info("Html report was saved in:", htmlpath)
		}
		absPath, err := filepath.Abs("." + string(filepath.Separator) + htmlpath)
		if err != nil {
			log.Println(err)
		} else {
			displayReport(absPath)
		}
	}
}

// displayReport function can be open system default browser automatic.
func displayReport(filePath string) {
	fileURL := fmt.Sprintf("file://%v", filePath)
	log.Println("To display report", fileURL, "in browser")
	var err error
	switch runtime.GOOS {
	case "linux":
		err = callSystemCmd("xdg-open", fileURL)
	case "darwin":
		err = callSystemCmd("open", fileURL)
	case "windows":
		r := strings.NewReplacer("&", "^&")
		err = callSystemCmd("cmd", "/c", "start", r.Replace(fileURL))
	default:
		err = fmt.Errorf("Unsupported platform,please view report file.")
	}
	if err != nil {
		log.Println(err)
	}
}

// callSystemCmd call system command opens a new browser window pointing to url.
func callSystemCmd(prog string, args ...string) error {
	cmd := exec.Command(prog, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
