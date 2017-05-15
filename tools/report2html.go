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

	"github.com/golang/glog"
	"github.com/wgliang/goreporter/engine"
)

// Json2Html will remake the reporter's data for the
// format we need for the template.
func Json2Html(jsonData []byte) (HtmlData, error) {
	var structData engine.Reporter
	var htmlData HtmlData

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
			test := Test{
				Path: pkgName,
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
				cycloInfo := CycloInfo{
					Comp: cycloTips[i].LineNumber,
					Info: cycloTips[i].ErrorString,
				}
				if cycloTips[i].LineNumber > 15 {
					htmlData.CycloBigThan15 = htmlData.CycloBigThan15 + 1
				}
				infos = append(infos, cycloInfo)
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

	// convert copy code result
	copyHtmlRes := make([]Simple, 0)
	if result, ok := structData.Metrics["CopyCodeTips"]; ok {
		for _, simpleResult := range result.Summaries {
			simpleTips := simpleResult.Errors

			for i := 0; i < len(simpleTips); i++ {
				copyCodeTip := strings.Split(simpleTips[i].ErrorString, ":")
				if len(copyCodeTip) == 4 {
					copycode := Simple{
						Path: strings.Join(copyCodeTip[0:3], ":"),
						Info: copyCodeTip[3],
					}
					copyHtmlRes = append(copyHtmlRes, copycode)
				}
			}
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
		for _, copyResult := range result.Summaries {
			simpleTips := copyResult.Errors

			for i := 0; i < len(simpleTips); i++ {
				deadCodeTip := strings.Split(simpleTips[i].ErrorString, ":")
				if len(deadCodeTip) == 4 {
					deadcode := Deadcode{
						Path: strings.Join(deadCodeTip[0:3], ":"),
						Info: deadCodeTip[3],
					}
					deadcodeHtmlRes = append(deadcodeHtmlRes, deadcode)
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
	htmlData.Date = structData.TimeStamp
	if len(importPackages) > 0 {
		htmlData.AveragePackageCover = float64(100 * (len(importPackages) - len(noTestPackages)) / len(importPackages))
	}
	return htmlData, nil
}

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
		glog.Infoln("create html report:", htmlpath)
		err = ioutil.WriteFile(htmlpath, out.Bytes(), 0666)
		if err != nil {
			glog.Errorln(err)
		}
	} else {
		htmlpath := projectName + "-" + timestamp + ".html"
		glog.Infoln("create html report:", htmlpath)
		err = ioutil.WriteFile(htmlpath, out.Bytes(), 0666)
		if err != nil {
			glog.Errorln(err)
		}
	}
}
