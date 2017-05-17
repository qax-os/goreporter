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
					var githubPath string
					srcLastIndex := strings.LastIndex(cycloTip[0], htmlData.Project)
					if srcLastIndex < len(cycloTip[0]) && srcLastIndex >= 0 {
						cycloTip[0] = cycloTip[0][srcLastIndex:]
						cycloTip[1] = "#L" + cycloTip[1]
						if len(htmlData.Project) < len(cycloTip[0]) {
							if strings.HasPrefix(htmlData.Project, "github.com") {
								cycloTip[0] = htmlData.Project + "/blob/master" + cycloTip[0][len(htmlData.Project):]
							}
							githubPath = cycloTip[0] + strings.Join(cycloTip[1:], ":")
						}
					}
					if githubPath == "" {
						githubPath = strings.Join(cycloTip[0:], ":")
					}

					cycloInfo := CycloInfo{
						Comp: cycloTips[i].LineNumber,
						Info: githubPath,
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
					var githubPath string
					srcLastIndex := strings.LastIndex(simpleCodeTip[0], htmlData.Project)
					if srcLastIndex < len(simpleCodeTip[0]) && srcLastIndex >= 0 {
						simpleCodeTip[0] = simpleCodeTip[0][srcLastIndex:]
						simpleCodeTip[1] = "#L" + simpleCodeTip[1]
						if len(htmlData.Project) < len(simpleCodeTip[0]) {
							if strings.HasPrefix(htmlData.Project, "github.com") {
								simpleCodeTip[0] = htmlData.Project + "/blob/master" + simpleCodeTip[0][len(htmlData.Project):]
							}
							githubPath = simpleCodeTip[0] + strings.Join(simpleCodeTip[1:3], ":")
						}
					}
					if githubPath == "" {
						githubPath = strings.Join(simpleCodeTip[0:], ":")
					}
					simpecode := Simple{
						Path: githubPath,
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

	// convert copy code result
	copyHtmlRes := make([]Copycode, 0)
	if result, ok := structData.Metrics["CopyCodeTips"]; ok {
		for _, copyResult := range result.Summaries {
			copyTips := copyResult.Errors
			var copyCodePathList []string
			for i := 0; i < len(copyTips); i++ {
				copyCodeTip := strings.Split(copyTips[i].ErrorString, ":")
				if len(copyCodeTip) == 2 {
					var githubPath string
					srcLastIndex := strings.LastIndex(copyCodeTip[0], htmlData.Project)
					if srcLastIndex < len(copyCodeTip[0]) && srcLastIndex >= 0 {
						copyCodeTip[0] = copyCodeTip[0][srcLastIndex:]
						copyCodeTip[1] = "#L" + copyCodeTip[1]
						if len(htmlData.Project) < len(copyCodeTip[0]) {
							if strings.HasPrefix(htmlData.Project, "github.com") {
								copyCodeTip[0] = htmlData.Project + "/blob/master" + copyCodeTip[0][len(htmlData.Project):]
							}
							githubPath = copyCodeTip[0] + copyCodeTip[1]
						}
					}
					if githubPath == "" {
						githubPath = strings.Join(copyCodeTip[0:], ":")
					}
					copyCodePathList = append(copyCodePathList, githubPath)
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
					var githubPath string
					srcLastIndex := strings.LastIndex(deadCodeTip[0], htmlData.Project)
					if srcLastIndex < len(deadCodeTip[0]) && srcLastIndex >= 0 {
						deadCodeTip[0] = deadCodeTip[0][srcLastIndex:]
						deadCodeTip[1] = "#L" + deadCodeTip[1]
						if len(htmlData.Project) < len(deadCodeTip[0]) {
							if strings.HasPrefix(htmlData.Project, "github.com") {
								deadCodeTip[0] = htmlData.Project + "/blob/master" + deadCodeTip[0][len(htmlData.Project):]
							}
							githubPath = deadCodeTip[0] + strings.Join(deadCodeTip[1:3], ":")
						}
					}
					if githubPath == "" {
						githubPath = strings.Join(deadCodeTip[0:], ":")
					}
					deadcode := Deadcode{
						Path: githubPath,
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

	if len(importPackages) > 0 {
		htmlData.AveragePackageCover = float64(100 * (len(importPackages) - len(noTestPackages)) / len(importPackages))
	} else {
		htmlData.AveragePackageCover = float64(100)
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
