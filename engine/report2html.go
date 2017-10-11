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
	"bytes"
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
	"time"

	"github.com/golang/glog"
	"github.com/json-iterator/go"
)

var (
	issues int
)

// UnitTest is a struct that contains some features in a report of html.
//         GoReporter HTML Report Features
//
//    +----------------------------------------------------------------------+
//    |        Feature        |                 Information                  |
//    +=======================+==============================================+
//    | Project               | The path address of the item being detected  |
//    +-----------------------+----------------------------------------------+
//    | Score                 | The score of the tested project              |
//    +-----------------------+----------------------------------------------+
//    | CodeTest              | Unit test results                            |
//    +-----------------------+----------------------------------------------+
//    | IssuesNum             | Issues number of the project                 |
//    +-----------------------+----------------------------------------------+
//    | CodeCount             | Number of lines of code                      |
//    +-----------------------+----------------------------------------------+
//    | CodeStyle             | Code style check                             |
//    +-----------------------+----------------------------------------------+
//    | CodeOptimization      | Code optimization                            |
//    +-----------------------+----------------------------------------------+
//    | CodeSmell             | Code smell                                   |
//    +-----------------------+----------------------------------------------+
//    | DepGraph              | Depend graph of all packages in the project  |
//    +-----------------------+----------------------------------------------+
//    | Date                  | Date assessment of the project               |
//    +-----------------------+----------------------------------------------+
//    | LastRefresh           | Last refresh time of one project             |
//    +-----------------------+----------------------------------------------+
//    | HumanizedLastRefresh  | Humanized last refresh setting               |
//    +-----------------------+----------------------------------------------+
//
// And the HtmlData just as data for default html template. If you want to customize
// your own template file, please follow these data, or you can redefine it yourself.
type HtmlData struct {
	Project          string
	Score            int
	IssuesNum        int
	CodeTest         string
	CodeStyle        string
	CodeOptimization string
	CodeCount        string
	CodeSmell        string
	DepGraph         template.HTML

	Date                 string
	LastRefresh          time.Time `json:"last_refresh"`
	HumanizedLastRefresh string    `json:"humanized_last_refresh"`
}

// converterUnitTest provides function that convert unit test data into the
// format required in the html template.It will extract from the structData
// need to convert the data.The result will be saved in the hd's attributes.
func (hd *HtmlData) converterCodeTest(structData Reporter) {
	var codeTestHtmlData CodeTest
	if result, ok := structData.Metrics["UnitTestTips"]; ok {
		var totalTime float64
		for pkgName, testRes := range result.Summaries {
			var packageUnitTestResult PackageTest
			jsoniter.Unmarshal([]byte(testRes.Description), &packageUnitTestResult)
			srcLastIndex := strings.LastIndex(pkgName, hd.Project)
			if !packageUnitTestResult.IsPass {
				codeTestHtmlData.Content.NoTest = append(codeTestHtmlData.Content.NoTest, pkgName)
			} else if srcLastIndex < len(pkgName) && srcLastIndex >= 0 {
				codeTestHtmlData.Content.Pkg = append(codeTestHtmlData.Content.Pkg, pkgName[srcLastIndex:])
				codeTestHtmlData.Content.Time = append(codeTestHtmlData.Content.Time, packageUnitTestResult.Time)
				totalTime = totalTime + packageUnitTestResult.Time
				if len(packageUnitTestResult.Coverage) > 1 {
					cover, _ := strconv.ParseFloat(packageUnitTestResult.Coverage[:(len(packageUnitTestResult.Coverage)-1)], 64)
					codeTestHtmlData.Content.Cover = append(codeTestHtmlData.Content.Cover, cover)
				}
			}
		}
		codeTestHtmlData.Summary.TotalTime = totalTime
		codeTestHtmlData.Summary.CodeCover = result.Percentage
		codeTestHtmlData.Summary.PackageCover = float64(len(codeTestHtmlData.Content.Pkg)) * 1.0 / float64(len(codeTestHtmlData.Content.Pkg)+len(codeTestHtmlData.Content.NoTest))
	}

	stringCodeTestJson, err := jsoniter.Marshal(codeTestHtmlData)
	if err != nil {
		glog.Errorln(err)
	}
	hd.CodeTest = string(stringCodeTestJson)
}

// converterUnitTest provides function that convert unit test data into the
// format required in the html template.It will extract from the structData
// need to convert the data.The result will be saved in the hd's attributes.
func (hd *HtmlData) converterCodeStyle(structData Reporter) {
	var codeStyleHtmlData CodeStyle
	codeSpellHtmlData := converterCodeSpell(structData)
	codeStyleHtmlData.Summary.FilesNum = codeStyleHtmlData.Summary.FilesNum + codeSpellHtmlData.filesNum
	codeStyleHtmlData.Summary.IssuesNum = codeStyleHtmlData.Summary.IssuesNum + codeSpellHtmlData.issuesNum
	codeStyleHtmlData.Content.MissSpell = codeSpellHtmlData

	codeLintHtmlData := converterCodeLint(structData)
	codeStyleHtmlData.Summary.FilesNum = codeStyleHtmlData.Summary.FilesNum + codeLintHtmlData.filesNum
	codeStyleHtmlData.Summary.IssuesNum = codeStyleHtmlData.Summary.IssuesNum + codeLintHtmlData.issuesNum
	codeStyleHtmlData.Content.MissSpell = codeLintHtmlData

	codeFmtHtmlData := converterCodeFmt(structData)
	codeStyleHtmlData.Summary.FilesNum = codeStyleHtmlData.Summary.FilesNum + codeFmtHtmlData.filesNum
	codeStyleHtmlData.Summary.IssuesNum = codeStyleHtmlData.Summary.IssuesNum + codeFmtHtmlData.issuesNum
	codeStyleHtmlData.Content.MissSpell = codeFmtHtmlData

	codeVetHtmlData := converterCodeVet(structData)
	codeStyleHtmlData.Summary.FilesNum = codeStyleHtmlData.Summary.FilesNum + codeVetHtmlData.filesNum
	codeStyleHtmlData.Summary.IssuesNum = codeStyleHtmlData.Summary.IssuesNum + codeVetHtmlData.issuesNum
	codeStyleHtmlData.Content.MissSpell = codeVetHtmlData

	stringCodeStyleJson, err := jsoniter.Marshal(codeStyleHtmlData)
	if err != nil {
		glog.Errorln(err)
	}
	hd.CodeStyle = string(stringCodeStyleJson)
}

// converterUnitTest provides function that convert unit test data into the
// format required in the html template.It will extract from the structData
// need to convert the data.The result will be saved in the hd's attributes.
func (hd *HtmlData) converterCodeOptimization(structData Reporter) {
	var codeOptimizationHtmlData CodeOptimization
	codeSimpleHtmlData := converterCodeSimple(structData)
	codeOptimizationHtmlData.Summary.FilesNum = codeOptimizationHtmlData.Summary.FilesNum + codeSimpleHtmlData.filesNum
	codeOptimizationHtmlData.Summary.IssuesNum = codeOptimizationHtmlData.Summary.IssuesNum + codeSimpleHtmlData.issuesNum
	codeOptimizationHtmlData.Content.SimpleCode = codeSimpleHtmlData

	codeDeadHtmlData := converterCodeDead(structData)
	codeOptimizationHtmlData.Summary.FilesNum = codeOptimizationHtmlData.Summary.FilesNum + codeDeadHtmlData.filesNum
	codeOptimizationHtmlData.Summary.IssuesNum = codeOptimizationHtmlData.Summary.IssuesNum + codeDeadHtmlData.issuesNum
	codeOptimizationHtmlData.Content.DeadCode = codeDeadHtmlData

	copyCodeHtmlData := converterCopyCode(structData)
	codeOptimizationHtmlData.Summary.FilesNum = codeOptimizationHtmlData.Summary.FilesNum + copyCodeHtmlData.filesNum
	codeOptimizationHtmlData.Summary.IssuesNum = codeOptimizationHtmlData.Summary.IssuesNum + copyCodeHtmlData.issuesNum
	codeOptimizationHtmlData.Content.CopyCode = copyCodeHtmlData

	codeInterfacerHtmlData := converterCodeInterfacer(structData)
	codeOptimizationHtmlData.Summary.FilesNum = codeOptimizationHtmlData.Summary.FilesNum + codeInterfacerHtmlData.filesNum
	codeOptimizationHtmlData.Summary.IssuesNum = codeOptimizationHtmlData.Summary.IssuesNum + codeInterfacerHtmlData.issuesNum
	codeOptimizationHtmlData.Content.InterfacerCode = codeInterfacerHtmlData

	stringCodeOptimizationJson, err := jsoniter.Marshal(codeOptimizationHtmlData)
	if err != nil {
		glog.Errorln(err)
	}
	hd.CodeOptimization = string(stringCodeOptimizationJson)
}

// converterCyclo provides function that convert cyclo data into the
// format required in the html template.It will extract from the structData
// need to convert the data.The result will be saved in the hd's attributes.
func (hd *HtmlData) converterCodeSmell(structData Reporter) {
	var codeSmellHtmlData CodeSmell
	codeSmellHtmlData.Content.Percentage = make(map[string]int, 0)
	codeSmellHtmlData.Content.List = make(map[string]int, 0)
	codeSmellHtmlData.Content.Percentage["1-15"] = 0
	codeSmellHtmlData.Content.Percentage["15-50"] = 0
	codeSmellHtmlData.Content.Percentage["50-"] = 0

	if result, ok := structData.Metrics["CycloTips"]; ok {
		filesMap := make(map[string]bool, 0)
		for pkgName, summary := range result.Summaries {
			var compNum, compSum int
			for i := 0; i < len(summary.Errors); i++ {
				cycloTip := strings.Split(summary.Errors[i].ErrorString, ":")
				if len(cycloTip) >= 3 {
					if summary.Errors[i].LineNumber < 15 {
						codeSmellHtmlData.Content.Percentage["1-15"]++
					} else if summary.Errors[i].LineNumber >= 15 {
						codeSmellHtmlData.Content.List[strings.Join(cycloTip[0:], ":")] = summary.Errors[i].LineNumber
						filesMap[strings.Join(cycloTip[0:], ":")] = true
						if summary.Errors[i].LineNumber < 50 {
							codeSmellHtmlData.Content.Percentage["15-50"]++
						} else {
							codeSmellHtmlData.Content.Percentage["50-"]++
						}
					}
					compNum++
					compSum = compSum + summary.Errors[i].LineNumber
				}
			}

			if compNum > 0 {
				codeSmellHtmlData.Content.Pkg = append(codeSmellHtmlData.Content.Pkg, pkgName)
				codeSmellHtmlData.Content.Cyclo = append(codeSmellHtmlData.Content.Cyclo, compSum/compNum)
			}

		}
		codeSmellHtmlData.Summary.FilesNum = len(filesMap)
		codeSmellHtmlData.Summary.IssuesNum = len(codeSmellHtmlData.Content.List)
	}

	stringCodeSmellJson, err := jsoniter.Marshal(codeSmellHtmlData)
	if err != nil {
		glog.Errorln(err)
	}
	hd.CodeSmell = string(stringCodeSmellJson)
}

// converterCount provides function that convert countcode data into the
// format required in the html template.It will extract from the structData
// need to convert the data.The result will be saved in the hd's attributes.
func (hd *HtmlData) converterCodeCount(structData Reporter) {
	var codeCountHtmlData CodeCount
	if result, ok := structData.Metrics["CountCodeTips"]; ok {
		codeCountHtmlData.Summary.FileCount, _ = strconv.Atoi(result.Summaries["FileCount"].Description)
		codeCountHtmlData.Summary.LineCount, _ = strconv.Atoi(result.Summaries["CodeLines"].Description)
		codeCountHtmlData.Summary.CommentCount, _ = strconv.Atoi(result.Summaries["CommentLines"].Description)
	}
	stringCodeCountJson, err := jsoniter.Marshal(codeCountHtmlData)
	if err != nil {
		glog.Errorln(err)
	}
	hd.CodeCount = string(stringCodeCountJson)
}

// converterCyclo provides function that convert cyclo data into the
// format required in the html template.It will extract from the structData
// need to convert the data.The result will be saved in the hd's attributes.
// func converterFunctionDepth(structData Reporter) {
// 	if result, ok := structData.Metrics["DepthTips"]; ok {
// 		for pkgName, summary := range result.Summaries {
// 			depthTips := summary.Errors
// 			depth := Depth{
// 				Pkg:  pkgName,
// 				Size: len(depthTips),
// 			}
// 			var infos []DepthInfo
// 			for i := 0; i < len(depthTips); i++ {
// 				depthTip := strings.Split(depthTips[i].ErrorString, ":")
// 				if len(depthTip) >= 3 {
// 					depthInfo := DepthInfo{
// 						Comp: depthTips[i].LineNumber,
// 						Info: strings.Join(depthTip[0:], ":"),
// 					}
// 					if depthTips[i].LineNumber > 3 {
// 						hd.CycloBigThan15 = hd.CycloBigThan15 + 1
// 						issues = issues + 1
// 					}
// 					infos = append(infos, depthInfo)
// 				}
// 			}
// 			depth.Info = infos
// 			depthHtmlRes = append(depthHtmlRes, depth)
// 		}
// 	}

// 	stringDepthJson, err := jsoniter.Marshal(depthHtmlRes)
// 	if err != nil {
// 		glog.Errorln(err)
// 	}
// 	hd.Depths = string(stringDepthJson)
// }

// converterSimple provides function that convert simplecode data into the
// format required in the html template.It will extract from the structData
// need to convert the data.The result will be saved in the hd's attributes.
func converterCodeSimple(structData Reporter) (simpleHtmlData StyleItem) {
	simpleHtmlData.Label = `gosimple is a linter for Go source code that specialises on simplifying code.`
	if result, ok := structData.Metrics["SimpleTips"]; ok {
		filesMap := make(map[string]bool, 0)
		for _, summary := range result.Summaries {
			simpleCodeTips := summary.Errors
			for i := 0; i < len(simpleCodeTips); i++ {
				simpleCodeTip := strings.Split(simpleCodeTips[i].ErrorString, ":")
				if len(simpleCodeTip) == 4 {
					simpecode := Item{
						File:    strings.Join(simpleCodeTip[0:3], ":"),
						Content: simpleCodeTip[3],
					}
					filesMap[simpecode.File] = true
					simpleHtmlData.Detail = append(simpleHtmlData.Detail, simpecode)
				} else if len(simpleCodeTip) == 5 {
					simpecode := Item{
						File:    strings.Join(simpleCodeTip[0:4], ":"),
						Content: simpleCodeTip[4],
					}
					filesMap[simpecode.File] = true
					simpleHtmlData.Detail = append(simpleHtmlData.Detail, simpecode)
				}
			}
		}
		simpleHtmlData.filesNum = len(filesMap)
		simpleHtmlData.issuesNum = len(simpleHtmlData.Detail)
	}

	return simpleHtmlData
}

// converterInterfacer provides function that convert interfacer data into the
// format required in the html template.It will extract from the structData
// need to convert the data.The result will be saved in the hd's attributes.
func converterCodeInterfacer(structData Reporter) (interfacerHtmlData StyleItem) {
	interfacerHtmlData.Label = `A linter that suggests interface types. In other words, it warns about the usage of types that are more specific than necessary.`
	if result, ok := structData.Metrics["InterfacerTips"]; ok {
		filesMap := make(map[string]bool, 0)
		for _, summary := range result.Summaries {
			interfacerCodeTips := summary.Errors
			for i := 0; i < len(interfacerCodeTips); i++ {
				interfacerCodeTip := strings.Split(interfacerCodeTips[i].ErrorString, ":")
				if len(interfacerCodeTip) == 4 {
					interfacer := Item{
						File:    strings.Join(interfacerCodeTip[0:3], ":"),
						Content: interfacerCodeTip[3],
					}
					filesMap[interfacer.File] = true
					interfacerHtmlData.Detail = append(interfacerHtmlData.Detail, interfacer)
				} else if len(interfacerCodeTip) == 5 {
					interfacer := Item{
						File:    strings.Join(interfacerCodeTip[0:4], ":"),
						Content: interfacerCodeTip[4],
					}
					filesMap[interfacer.File] = true
					interfacerHtmlData.Detail = append(interfacerHtmlData.Detail, interfacer)
				}
			}
		}
	}

	return interfacerHtmlData
}

// converterCopy provides function that convert copycode data into the
// format required in the html template.It will extract from the structData
// need to convert the data.The result will be saved in the hd's attributes.
func converterCopyCode(structData Reporter) (copyHtmlData CopyItem) {
	copyHtmlData.Label = `Find code clones. So far it can find clones only in the Go source files. The method uses suffix tree for serialized ASTs. It ignores values of AST nodes.`
	if result, ok := structData.Metrics["CopyCodeTips"]; ok {
		filesMap := make(map[string]bool, 0)
		for _, copyResult := range result.Summaries {
			copyTips := copyResult.Errors
			var copyCodePathList []string
			for i := 0; i < len(copyTips); i++ {
				copyCodeTip := strings.Split(copyTips[i].ErrorString, ":")
				if len(copyCodeTip) >= 2 {
					newPath := strings.Join(copyCodeTip[0:], ":")
					filesMap[newPath] = true
					copyCodePathList = append(copyCodePathList, newPath)
				}
			}
			copyHtmlData.Detail = append(copyHtmlData.Detail, copyCodePathList)
		}
		copyHtmlData.filesNum = len(filesMap)
		copyHtmlData.issuesNum = len(copyHtmlData.Detail)
	}

	return copyHtmlData
}

// converterDead provides function that convert deadcode data into the
// format required in the html template.It will extract from the structData
// need to convert the data.The result will be saved in the hd's attributes.
func converterCodeDead(structData Reporter) (deadHtmlData StyleItem) {
	deadHtmlData.Label = `Unused code.`
	if result, ok := structData.Metrics["DeadCodeTips"]; ok {
		filesMap := make(map[string]bool, 0)
		for _, deadCodeResult := range result.Summaries {
			deadCodeTips := deadCodeResult.Errors
			for i := 0; i < len(deadCodeTips); i++ {
				deadCodeTip := strings.Split(deadCodeTips[i].ErrorString, ":")
				if len(deadCodeTip) == 4 {
					deadcode := Item{
						File:    strings.Join(deadCodeTip[0:3], ":"),
						Content: deadCodeTip[3],
					}
					filesMap[deadcode.File] = true
					deadHtmlData.Detail = append(deadHtmlData.Detail, deadcode)
				} else if len(deadCodeTip) == 5 {
					deadcode := Item{
						File:    strings.Join(deadCodeTip[0:4], ":"),
						Content: deadCodeTip[4],
					}
					filesMap[deadcode.File] = true
					deadHtmlData.Detail = append(deadHtmlData.Detail, deadcode)
				}
			}
		}
	}

	return deadHtmlData
}

// converterSpell provides function that convert spellcheck data into the
// format required in the html template.It will extract from the structData
// need to convert the data.The result will be saved in the hd's attributes.
func converterCodeSpell(structData Reporter) (spellHtmlData StyleItem) {
	spellHtmlData.Label = `Correct commonly misspelled English words... quickly`
	if result, ok := structData.Metrics["SpellCheckTips"]; ok {
		fileMap := make(map[string]bool, 0)
		for _, summary := range result.Summaries {
			spellCodeTips := summary.Errors
			for i := 0; i < len(spellCodeTips); i++ {
				spellCodeTip := strings.Split(spellCodeTips[i].ErrorString, ":")
				if len(spellCodeTip) == 4 {
					spellcode := Item{
						File:    strings.Join(spellCodeTip[0:3], ":"),
						Content: spellCodeTip[3],
					}
					fileMap[spellcode.File] = true
					spellHtmlData.Detail = append(spellHtmlData.Detail, spellcode)
				} else if len(spellCodeTip) == 5 {
					spellcode := Item{
						File:    strings.Join(spellCodeTip[0:4], ":"),
						Content: spellCodeTip[4],
					}
					fileMap[spellcode.File] = true
					spellHtmlData.Detail = append(spellHtmlData.Detail, spellcode)
				}
			}
		}
		spellHtmlData.filesNum = len(fileMap)
		spellHtmlData.issuesNum = len(spellHtmlData.Detail)
	}

	return spellHtmlData
}

// converterCodeLint provides function that convert spellcheck data into the
// format required in the html template.It will extract from the structData
// need to convert the data.The result will be saved in the hd's attributes.
func converterCodeLint(structData Reporter) (spellHtmlData StyleItem) {
	spellHtmlData.Label = `Correct commonly misspelled English words... quickly`
	if result, ok := structData.Metrics["GoLintTips"]; ok {
		fileMap := make(map[string]bool, 0)
		for _, summary := range result.Summaries {
			spellCodeTips := summary.Errors
			for i := 0; i < len(spellCodeTips); i++ {
				spellCodeTip := strings.Split(spellCodeTips[i].ErrorString, ":")
				if len(spellCodeTip) == 4 {
					spellcode := Item{
						File:    strings.Join(spellCodeTip[0:3], ":"),
						Content: spellCodeTip[3],
					}
					fileMap[spellcode.File] = true
					spellHtmlData.Detail = append(spellHtmlData.Detail, spellcode)
				} else if len(spellCodeTip) == 5 {
					spellcode := Item{
						File:    strings.Join(spellCodeTip[0:4], ":"),
						Content: spellCodeTip[4],
					}
					fileMap[spellcode.File] = true
					spellHtmlData.Detail = append(spellHtmlData.Detail, spellcode)
				}
			}
		}
		spellHtmlData.filesNum = len(fileMap)
		spellHtmlData.issuesNum = len(spellHtmlData.Detail)
	}

	return spellHtmlData
}

// converterCodeFmt provides function that convert spellcheck data into the
// format required in the html template.It will extract from the structData
// need to convert the data.The result will be saved in the hd's attributes.
func converterCodeFmt(structData Reporter) (spellHtmlData StyleItem) {
	spellHtmlData.Label = `Correct commonly misspelled English words... quickly`
	if result, ok := structData.Metrics["GoFmtTips"]; ok {
		fileMap := make(map[string]bool, 0)
		for _, summary := range result.Summaries {
			spellCodeTips := summary.Errors
			for i := 0; i < len(spellCodeTips); i++ {
				spellCodeTip := strings.Split(spellCodeTips[i].ErrorString, ":")
				if len(spellCodeTip) == 3 {
					spellcode := Item{
						File:    strings.Join(spellCodeTip[0:1], ":"),
						Content: strings.Join(spellCodeTip[1:3], ":"),
					}
					fileMap[spellcode.File] = true
					spellHtmlData.Detail = append(spellHtmlData.Detail, spellcode)
				} else if len(spellCodeTip) == 4 {
					spellcode := Item{
						File:    strings.Join(spellCodeTip[0:2], ":"),
						Content: strings.Join(spellCodeTip[2:4], ":"),
					}
					fileMap[spellcode.File] = true
					spellHtmlData.Detail = append(spellHtmlData.Detail, spellcode)
				}
			}
		}
		spellHtmlData.filesNum = len(fileMap)
		spellHtmlData.issuesNum = len(spellHtmlData.Detail)
	}

	return spellHtmlData
}

// converterCodeVet provides function that convert spellcheck data into the
// format required in the html template.It will extract from the structData
// need to convert the data.The result will be saved in the hd's attributes.
func converterCodeVet(structData Reporter) (spellHtmlData StyleItem) {
	spellHtmlData.Label = `Correct commonly misspelled English words... quickly`
	if result, ok := structData.Metrics["GoVetTips"]; ok {
		fileMap := make(map[string]bool, 0)
		for _, summary := range result.Summaries {
			spellCodeTips := summary.Errors
			for i := 0; i < len(spellCodeTips); i++ {
				spellCodeTip := strings.Split(spellCodeTips[i].ErrorString, ":")
				if len(spellCodeTip) == 4 {
					spellcode := Item{
						File:    strings.Join(spellCodeTip[0:3], ":"),
						Content: spellCodeTip[3],
					}
					fileMap[spellcode.File] = true
					spellHtmlData.Detail = append(spellHtmlData.Detail, spellcode)
				} else if len(spellCodeTip) == 5 {
					spellcode := Item{
						File:    strings.Join(spellCodeTip[0:4], ":"),
						Content: spellCodeTip[4],
					}
					fileMap[spellcode.File] = true
					spellHtmlData.Detail = append(spellHtmlData.Detail, spellcode)
				}
			}
		}
		spellHtmlData.filesNum = len(fileMap)
		spellHtmlData.issuesNum = len(spellHtmlData.Detail)
	}

	return spellHtmlData
}

// converterDependGraph provides function that convert depend graph data into the
// format required in the html template.It will extract from the structData
// need to convert the data.The result will be saved in the hd's attributes.
func (hd *HtmlData) converterDependGraph(structData Reporter) {
	hd.DepGraph = template.HTML(structData.Metrics["DependGraphTips"].Summaries["graph"].Description)
}

// SaveAsHtml is a function that save HtmlData as a html report.And will receive
// htmlData, projectPath, savePath and tpl parameters.
func SaveAsHtml(htmlData HtmlData, projectPath, savePath, timestamp, tpl string) {
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
	projectName := ProjectName(projectPath)
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
