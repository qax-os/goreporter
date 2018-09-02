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

	"github.com/360EntSecGroup-Skylar/goreporter/utils"
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
		codeTestHtmlData.Summary.TotalTime, _ = strconv.ParseFloat(strconv.FormatFloat(totalTime, 'f', 1, 64), 64)
		codeTestHtmlData.Summary.CodeCover, _ = strconv.ParseFloat(strconv.FormatFloat(result.Percentage, 'f', 1, 64), 64)
		if (len(codeTestHtmlData.Content.Pkg) + len(codeTestHtmlData.Content.NoTest)) == 0 {
			codeTestHtmlData.Summary.PackageCover = 0
		} else {
			codeTestHtmlData.Summary.PackageCover, _ = strconv.ParseFloat(strconv.FormatFloat(100*float64(len(codeTestHtmlData.Content.Pkg))*1.0/float64(len(codeTestHtmlData.Content.Pkg)+len(codeTestHtmlData.Content.NoTest)), 'f', 1, 64), 64)
		}
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
	codeStyleHtmlData.Content.GoLint = codeLintHtmlData

	codeFmtHtmlData := converterCodeFmt(structData)
	codeStyleHtmlData.Summary.FilesNum = codeStyleHtmlData.Summary.FilesNum + codeFmtHtmlData.filesNum
	codeStyleHtmlData.Summary.IssuesNum = codeStyleHtmlData.Summary.IssuesNum + codeFmtHtmlData.issuesNum
	codeStyleHtmlData.Content.GoFmt = codeFmtHtmlData

	codeVetHtmlData := converterCodeVet(structData)
	codeStyleHtmlData.Summary.FilesNum = codeStyleHtmlData.Summary.FilesNum + codeVetHtmlData.filesNum
	codeStyleHtmlData.Summary.IssuesNum = codeStyleHtmlData.Summary.IssuesNum + codeVetHtmlData.issuesNum
	codeStyleHtmlData.Content.GoVet = codeVetHtmlData

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
	codeSmellHtmlData.Content.List = make([]CodeSmellItem, 0)
	codeSmellHtmlData.Content.Percentage["1-15"] = 0
	codeSmellHtmlData.Content.Percentage["15-50"] = 0
	codeSmellHtmlData.Content.Percentage["50+"] = 0

	sumComp, sumNum := 0, 0
	if result, ok := structData.Metrics["CycloTips"]; ok {
		filesMap := make(map[string]bool, 0)
		for pkgName, summary := range result.Summaries {
			var compNum, compSum int
			for i := 0; i < len(summary.Errors); i++ {
				cycloTip := strings.Split(summary.Errors[i].ErrorString, ":")
				if len(cycloTip) >= 3 {
					if summary.Errors[i].LineNumber < 15 {
						codeSmellHtmlData.Content.Percentage["1-15"]++
						smellItem := CodeSmellItem{
							Path:  strings.Join(cycloTip[0:], ":"),
							Cyclo: summary.Errors[i].LineNumber,
						}
						codeSmellHtmlData.Content.List = append(codeSmellHtmlData.Content.List, smellItem)
						filesMap[strings.Join(cycloTip[0:], ":")] = true
					} else if summary.Errors[i].LineNumber >= 15 {
						smellItem := CodeSmellItem{
							Path:  strings.Join(cycloTip[0:], ":"),
							Cyclo: summary.Errors[i].LineNumber,
						}
						codeSmellHtmlData.Content.List = append(codeSmellHtmlData.Content.List, smellItem)
						filesMap[strings.Join(cycloTip[0:], ":")] = true
						if summary.Errors[i].LineNumber < 50 {
							codeSmellHtmlData.Content.Percentage["15-50"]++
						} else {
							codeSmellHtmlData.Content.Percentage["50+"]++
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
			sumComp = sumComp + compSum
			sumNum = sumNum + compNum
		}

		sortCycloByComp(codeSmellHtmlData.Content.List, 0, len(codeSmellHtmlData.Content.List))
		if sumNum == 0 {
			codeSmellHtmlData.Summary.CycloAvg = 0
		} else {
			codeSmellHtmlData.Summary.CycloAvg = sumComp / sumNum
		}
		codeSmellHtmlData.Summary.CycloHigh = codeSmellHtmlData.Content.Percentage["15-50"]
		codeSmellHtmlData.Summary.CycloGrave = codeSmellHtmlData.Content.Percentage["50+"]
	}

	stringCodeSmellJson, err := jsoniter.Marshal(codeSmellHtmlData)
	if err != nil {
		glog.Errorln(err)
	}
	hd.CodeSmell = string(stringCodeSmellJson)
}

// sortCycloByComp implements the quick sorting algorithm sort list by complexity.
func sortCycloByComp(input []CodeSmellItem, l, u int) {
	if l < u {
		m := partition(input, l, u)
		sortCycloByComp(input, l, m-1)
		sortCycloByComp(input, m, u)
	}
}

func partition(input []CodeSmellItem, l, u int) int {
	var (
		pivot = input[l]
		left  = l
		right = l + 1
	)
	for ; right < u; right++ {
		if input[right].Cyclo >= pivot.Cyclo {
			left++
			input[left], input[right] = input[right], input[left]
		}
	}
	input[l], input[left] = input[left], input[l]
	return left + 1
}

// converterCount provides function that convert countcode data into the
// format required in the html template.It will extract from the structData
// need to convert the data.The result will be saved in the hd's attributes.
func (hd *HtmlData) converterCodeCount(structData Reporter) {
	var codeCountHtmlData CodeCount

	fileFuncsCount := make(map[string]int, 0)
	pkgFuncsCount := make(map[string]int, 0)
	if result, ok := structData.Metrics["CycloTips"]; ok {
		for _, summary := range result.Summaries {
			for i := 0; i < len(summary.Errors); i++ {
				sepFileIndex := strings.LastIndex(summary.Errors[i].ErrorString, ".go")
				if sepFileIndex >= 0 && (sepFileIndex+3) < len(summary.Errors[i].ErrorString) {
					fileFuncsCount[string(summary.Errors[i].ErrorString[0:sepFileIndex+3])]++
				}
				sepPkgIndex := strings.LastIndex(summary.Errors[i].ErrorString, string(filepath.Separator))
				if sepPkgIndex >= 0 && sepPkgIndex < len(summary.Errors[i].ErrorString) {
					pkgFuncsCount[string(summary.Errors[i].ErrorString[0:sepPkgIndex])]++
				}
			}
		}
	}
	pkgCommentCount := make(map[string]int, 0)
	pkgLineCount := make(map[string]int, 0)

	if result, ok := structData.Metrics["CountCodeTips"]; ok {
		for fileName, codeCount := range result.Summaries {
			counts := strings.Split(codeCount.Description, ";")
			if len(counts) == 4 {
				codeCountHtmlData.Content.File = append(codeCountHtmlData.Content.File, fileName)
				fileCommentCount, _ := strconv.Atoi(counts[2])
				codeCountHtmlData.Content.FileCommentCount = append(codeCountHtmlData.Content.FileCommentCount, fileCommentCount)

				codeCountHtmlData.Content.FileFunctionCount = append(codeCountHtmlData.Content.FileFunctionCount, fileFuncsCount[fileName])
				// Add into summary.
				codeCountHtmlData.Summary.FunctionCount = codeCountHtmlData.Summary.FunctionCount + fileFuncsCount[fileName]
				fileLineCount, _ := strconv.Atoi(counts[1])
				codeCountHtmlData.Content.FileLineCount = append(codeCountHtmlData.Content.FileLineCount, fileLineCount)

				sepPkgIndex := strings.LastIndex(fileName, string(filepath.Separator))
				if sepPkgIndex >= 0 && sepPkgIndex < len(fileName) {
					pkgCommentCount[string(fileName[0:sepPkgIndex])] = pkgCommentCount[string(fileName[0:sepPkgIndex])] + fileCommentCount
					pkgLineCount[string(fileName[0:sepPkgIndex])] = pkgLineCount[string(fileName[0:sepPkgIndex])] + fileLineCount
				}
			}
		}

		for pkgName, commentCount := range pkgCommentCount {
			codeCountHtmlData.Content.Pkg = append(codeCountHtmlData.Content.Pkg, pkgName)
			codeCountHtmlData.Content.PkgCommentCount = append(codeCountHtmlData.Content.PkgCommentCount, commentCount)
			codeCountHtmlData.Content.PkgFunctionCount = append(codeCountHtmlData.Content.PkgFunctionCount, pkgFuncsCount[pkgName])
			codeCountHtmlData.Content.PkgLineCount = append(codeCountHtmlData.Content.PkgLineCount, pkgLineCount[pkgName])

			codeCountHtmlData.Summary.LineCount = codeCountHtmlData.Summary.LineCount + pkgLineCount[pkgName]
			codeCountHtmlData.Summary.CommentCount = codeCountHtmlData.Summary.CommentCount + commentCount
			codeCountHtmlData.Summary.FunctionCount = codeCountHtmlData.Summary.FunctionCount + pkgFuncsCount[pkgName]
		}
		codeCountHtmlData.Summary.FileCount = len(result.Summaries)
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
		fileMap := make(map[string]bool, 0)
		mapItem2DetailIndex := make(map[string]int, 0)
		for _, summary := range result.Summaries {
			simpleCodeTips := summary.Errors
			for i := 0; i < len(simpleCodeTips); i++ {
				simpleCodeTip := strings.Split(simpleCodeTips[i].ErrorString, ":")
				if len(simpleCodeTip) == 4 {
					if fileIndex, ok := mapItem2DetailIndex[strings.Join(simpleCodeTip[0:2], ":")]; ok {
						simpleHtmlData.Detail[fileIndex].Content = append(simpleHtmlData.Detail[fileIndex].Content, strings.Join(simpleCodeTip[2:], ":"))
					} else {
						spellcode := Item{
							File: strings.Join(simpleCodeTip[0:2], ":"),
						}
						spellcode.Content = append(spellcode.Content, strings.Join(simpleCodeTip[2:], ":"))
						fileMap[spellcode.File] = true
						mapItem2DetailIndex[strings.Join(simpleCodeTip[0:2], ":")] = len(simpleHtmlData.Detail)
						simpleHtmlData.Detail = append(simpleHtmlData.Detail, spellcode)
					}
				} else if len(simpleCodeTip) == 5 {
					if fileIndex, ok := mapItem2DetailIndex[strings.Join(simpleCodeTip[0:3], ":")]; ok {
						simpleHtmlData.Detail[fileIndex].Content = append(simpleHtmlData.Detail[fileIndex].Content, strings.Join(simpleCodeTip[3:], ":"))
					} else {
						spellcode := Item{
							File: strings.Join(simpleCodeTip[0:3], ":"),
						}
						spellcode.Content = append(spellcode.Content, strings.Join(simpleCodeTip[3:], ":"))
						fileMap[spellcode.File] = true
						mapItem2DetailIndex[strings.Join(simpleCodeTip[0:3], ":")] = len(simpleHtmlData.Detail)
						simpleHtmlData.Detail = append(simpleHtmlData.Detail, spellcode)
					}
				}
			}
		}
		simpleHtmlData.filesNum = len(fileMap)
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
		mapItem2DetailIndex := make(map[string]int, 0)
		for _, summary := range result.Summaries {
			interfacerCodeTips := summary.Errors
			for i := 0; i < len(interfacerCodeTips); i++ {
				interfacerCodeTip := strings.Split(interfacerCodeTips[i].ErrorString, ":")
				if len(interfacerCodeTip) == 4 {
					if fileIndex, ok := mapItem2DetailIndex[strings.Join(interfacerCodeTip[0:2], ":")]; ok {
						interfacerHtmlData.Detail[fileIndex].Content = append(interfacerHtmlData.Detail[fileIndex].Content, strings.Join(interfacerCodeTip[2:], ":"))
					} else {
						spellcode := Item{
							File: strings.Join(interfacerCodeTip[0:2], ":"),
						}
						spellcode.Content = append(spellcode.Content, strings.Join(interfacerCodeTip[2:], ":"))
						filesMap[spellcode.File] = true
						mapItem2DetailIndex[strings.Join(interfacerCodeTip[0:2], ":")] = len(interfacerHtmlData.Detail)
						interfacerHtmlData.Detail = append(interfacerHtmlData.Detail, spellcode)
					}
				} else if len(interfacerCodeTip) == 5 {
					if fileIndex, ok := mapItem2DetailIndex[strings.Join(interfacerCodeTip[0:3], ":")]; ok {
						interfacerHtmlData.Detail[fileIndex].Content = append(interfacerHtmlData.Detail[fileIndex].Content, strings.Join(interfacerCodeTip[3:], ":"))
					} else {
						spellcode := Item{
							File: strings.Join(interfacerCodeTip[0:3], ":"),
						}
						spellcode.Content = append(spellcode.Content, strings.Join(interfacerCodeTip[3:], ":"))
						filesMap[spellcode.File] = true
						mapItem2DetailIndex[strings.Join(interfacerCodeTip[0:3], ":")] = len(interfacerHtmlData.Detail)
						interfacerHtmlData.Detail = append(interfacerHtmlData.Detail, spellcode)
					}
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
		mapItem2DetailIndex := make(map[string]int, 0)
		for _, deadCodeResult := range result.Summaries {
			deadCodeTips := deadCodeResult.Errors
			for i := 0; i < len(deadCodeTips); i++ {
				deadCodeTip := strings.Split(deadCodeTips[i].ErrorString, ":")
				if len(deadCodeTip) == 4 {
					if fileIndex, ok := mapItem2DetailIndex[strings.Join(deadCodeTip[0:2], ":")]; ok {
						deadHtmlData.Detail[fileIndex].Content = append(deadHtmlData.Detail[fileIndex].Content, strings.Join(deadCodeTip[2:], ":"))
					} else {
						spellcode := Item{
							File: strings.Join(deadCodeTip[0:2], ":"),
						}
						spellcode.Content = append(spellcode.Content, strings.Join(deadCodeTip[2:], ":"))
						filesMap[spellcode.File] = true
						mapItem2DetailIndex[strings.Join(deadCodeTip[0:2], ":")] = len(deadHtmlData.Detail)
						deadHtmlData.Detail = append(deadHtmlData.Detail, spellcode)
					}
				} else if len(deadCodeTip) == 5 {
					if fileIndex, ok := mapItem2DetailIndex[strings.Join(deadCodeTip[0:3], ":")]; ok {
						deadHtmlData.Detail[fileIndex].Content = append(deadHtmlData.Detail[fileIndex].Content, strings.Join(deadCodeTip[3:], ":"))
					} else {
						spellcode := Item{
							File: strings.Join(deadCodeTip[0:3], ":"),
						}
						spellcode.Content = append(spellcode.Content, strings.Join(deadCodeTip[3:], ":"))
						filesMap[spellcode.File] = true
						mapItem2DetailIndex[strings.Join(deadCodeTip[0:3], ":")] = len(deadHtmlData.Detail)
						deadHtmlData.Detail = append(deadHtmlData.Detail, spellcode)
					}
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
		filesMap := make(map[string]bool, 0)
		mapItem2DetailIndex := make(map[string]int, 0)
		for _, summary := range result.Summaries {
			spellCodeTips := summary.Errors
			for i := 0; i < len(spellCodeTips); i++ {
				spellCodeTip := strings.Split(spellCodeTips[i].ErrorString, ":")
				if len(spellCodeTip) == 4 {
					if fileIndex, ok := mapItem2DetailIndex[strings.Join(spellCodeTip[0:2], ":")]; ok {
						spellHtmlData.Detail[fileIndex].Content = append(spellHtmlData.Detail[fileIndex].Content, strings.Join(spellCodeTip[2:], ":"))
					} else {
						spellcode := Item{
							File: strings.Join(spellCodeTip[0:2], ":"),
						}
						spellcode.Content = append(spellcode.Content, strings.Join(spellCodeTip[2:], ":"))
						filesMap[spellcode.File] = true
						mapItem2DetailIndex[strings.Join(spellCodeTip[0:2], ":")] = len(spellHtmlData.Detail)
						spellHtmlData.Detail = append(spellHtmlData.Detail, spellcode)
					}
				} else if len(spellCodeTip) == 5 {
					if fileIndex, ok := mapItem2DetailIndex[strings.Join(spellCodeTip[0:3], ":")]; ok {
						spellHtmlData.Detail[fileIndex].Content = append(spellHtmlData.Detail[fileIndex].Content, strings.Join(spellCodeTip[3:], ":"))
					} else {
						spellcode := Item{
							File: strings.Join(spellCodeTip[0:3], ":"),
						}
						spellcode.Content = append(spellcode.Content, strings.Join(spellCodeTip[3:], ":"))
						filesMap[spellcode.File] = true
						mapItem2DetailIndex[strings.Join(spellCodeTip[0:3], ":")] = len(spellHtmlData.Detail)
						spellHtmlData.Detail = append(spellHtmlData.Detail, spellcode)
					}
				}
			}
		}
		spellHtmlData.filesNum = len(filesMap)
		spellHtmlData.issuesNum = len(spellHtmlData.Detail)
	}

	return spellHtmlData
}

// converterCodeLint provides function that convert spellcheck data into the
// format required in the html template.It will extract from the structData
// need to convert the data.The result will be saved in the hd's attributes.
func converterCodeLint(structData Reporter) (lintHtmlData StyleItem) {
	lintHtmlData.Label = `Correct commonly misspelled English words... quickly`
	if result, ok := structData.Metrics["GoLintTips"]; ok {
		fileMap := make(map[string]bool, 0)
		mapItem2DetailIndex := make(map[string]int, 0)
		for _, summary := range result.Summaries {
			spellCodeTips := summary.Errors
			for i := 0; i < len(spellCodeTips); i++ {
				lintCodeTip := strings.Split(spellCodeTips[i].ErrorString, ":")
				if len(lintCodeTip) == 4 {
					if fileIndex, ok := mapItem2DetailIndex[strings.Join(lintCodeTip[0:2], ":")]; ok {
						lintHtmlData.Detail[fileIndex].Content = append(lintHtmlData.Detail[fileIndex].Content, strings.Join(lintCodeTip[2:], ":"))
					} else {
						spellcode := Item{
							File: strings.Join(lintCodeTip[0:2], ":"),
						}
						spellcode.Content = append(spellcode.Content, strings.Join(lintCodeTip[2:], ":"))
						fileMap[spellcode.File] = true
						mapItem2DetailIndex[strings.Join(lintCodeTip[0:2], ":")] = len(lintHtmlData.Detail)
						lintHtmlData.Detail = append(lintHtmlData.Detail, spellcode)
					}
				} else if len(lintCodeTip) == 5 {
					if fileIndex, ok := mapItem2DetailIndex[strings.Join(lintCodeTip[0:3], ":")]; ok {
						lintHtmlData.Detail[fileIndex].Content = append(lintHtmlData.Detail[fileIndex].Content, strings.Join(lintCodeTip[3:], ":"))
					} else {
						spellcode := Item{
							File: strings.Join(lintCodeTip[0:3], ":"),
						}
						spellcode.Content = append(spellcode.Content, strings.Join(lintCodeTip[3:], ":"))
						fileMap[spellcode.File] = true
						mapItem2DetailIndex[strings.Join(lintCodeTip[0:3], ":")] = len(lintHtmlData.Detail)
						lintHtmlData.Detail = append(lintHtmlData.Detail, spellcode)
					}
				}
			}
		}
		lintHtmlData.filesNum = len(fileMap)
		lintHtmlData.issuesNum = len(lintHtmlData.Detail)
	}

	return lintHtmlData
}

// converterCodeFmt provides function that convert spellcheck data into the
// format required in the html template.It will extract from the structData
// need to convert the data.The result will be saved in the hd's attributes.
func converterCodeFmt(structData Reporter) (fmtHtmlData StyleItem) {
	fmtHtmlData.Label = `Correct commonly misspelled English words... quickly`
	if result, ok := structData.Metrics["GoFmtTips"]; ok {
		filesMap := make(map[string]bool, 0)
		mapItem2DetailIndex := make(map[string]int, 0)
		for _, summary := range result.Summaries {
			fmtCodeTips := summary.Errors
			for i := 0; i < len(fmtCodeTips); i++ {
				fmtCodeTip := strings.Split(fmtCodeTips[i].ErrorString, ":")
				if len(fmtCodeTip) == 3 {
					if fileIndex, ok := mapItem2DetailIndex[strings.Join(fmtCodeTip[0:1], ":")]; ok {
						fmtHtmlData.Detail[fileIndex].Content = append(fmtHtmlData.Detail[fileIndex].Content, strings.Join(fmtCodeTip[1:], ":"))
					} else {
						spellcode := Item{
							File: strings.Join(fmtCodeTip[0:1], ":"),
						}
						spellcode.Content = append(spellcode.Content, strings.Join(fmtCodeTip[1:], ":"))
						filesMap[spellcode.File] = true
						fmtHtmlData.Detail = append(fmtHtmlData.Detail, spellcode)
					}
				} else if len(fmtCodeTip) == 4 {
					if fileIndex, ok := mapItem2DetailIndex[strings.Join(fmtCodeTip[0:2], ":")]; ok {
						fmtHtmlData.Detail[fileIndex].Content = append(fmtHtmlData.Detail[fileIndex].Content, strings.Join(fmtCodeTip[2:], ":"))
					} else {
						spellcode := Item{
							File: strings.Join(fmtCodeTip[0:2], ":"),
						}
						spellcode.Content = append(spellcode.Content, strings.Join(fmtCodeTip[2:], ":"))
						filesMap[spellcode.File] = true
						mapItem2DetailIndex[strings.Join(fmtCodeTip[0:3], ":")] = len(fmtHtmlData.Detail)
						fmtHtmlData.Detail = append(fmtHtmlData.Detail, spellcode)
					}
				}
			}
		}
		fmtHtmlData.filesNum = len(filesMap)
		fmtHtmlData.issuesNum = len(fmtHtmlData.Detail)
	}

	return fmtHtmlData
}

// converterCodeVet provides function that convert spellcheck data into the
// format required in the html template.It will extract from the structData
// need to convert the data.The result will be saved in the hd's attributes.
func converterCodeVet(structData Reporter) (vetHtmlData StyleItem) {
	vetHtmlData.Label = `Correct commonly misspelled English words... quickly`
	if result, ok := structData.Metrics["GoVetTips"]; ok {
		fileMap := make(map[string]bool, 0)
		mapItem2DetailIndex := make(map[string]int, 0)
		for _, summary := range result.Summaries {
			vetCodeTips := summary.Errors
			for i := 0; i < len(vetCodeTips); i++ {
				vetCodeTip := strings.Split(vetCodeTips[i].ErrorString, ":")
				if len(vetCodeTip) == 4 {
					if fileIndex, ok := mapItem2DetailIndex[strings.Join(vetCodeTip[0:2], ":")]; ok {
						vetHtmlData.Detail[fileIndex].Content = append(vetHtmlData.Detail[fileIndex].Content, strings.Join(vetCodeTip[2:], ":"))
					} else {
						spellcode := Item{
							File: strings.Join(vetCodeTip[0:2], ":"),
						}
						spellcode.Content = append(spellcode.Content, strings.Join(vetCodeTip[2:], ":"))
						fileMap[spellcode.File] = true
						mapItem2DetailIndex[strings.Join(vetCodeTip[0:2], ":")] = len(vetHtmlData.Detail)
						vetHtmlData.Detail = append(vetHtmlData.Detail, spellcode)
					}
				} else if len(vetCodeTip) == 5 {
					if fileIndex, ok := mapItem2DetailIndex[strings.Join(vetCodeTip[0:3], ":")]; ok {
						vetHtmlData.Detail[fileIndex].Content = append(vetHtmlData.Detail[fileIndex].Content, strings.Join(vetCodeTip[3:], ":"))
					} else {
						spellcode := Item{
							File: strings.Join(vetCodeTip[0:3], ":"),
						}
						spellcode.Content = append(spellcode.Content, strings.Join(vetCodeTip[3:], ":"))
						fileMap[spellcode.File] = true
						mapItem2DetailIndex[strings.Join(vetCodeTip[0:3], ":")] = len(vetHtmlData.Detail)
						vetHtmlData.Detail = append(vetHtmlData.Detail, spellcode)
					}
				}

			}
		}
		vetHtmlData.filesNum = len(fileMap)
		vetHtmlData.issuesNum = len(vetHtmlData.Detail)
	}

	return vetHtmlData
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
	projectName := utils.ProjectName(projectPath)
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
