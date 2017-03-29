package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func Json2Html(jsonData []byte, path string) (HtmlData, error) {
	var structData Reporter
	var htmlData HtmlData
	if jsonData == nil {
		return htmlData, errors.New("json is null")
	}
	json.Unmarshal(jsonData, &structData)

	htmlData.Project = structData.Project

	// convert test result
	testHtmlRes := make([]Test, 0)
	for pkgName, testRes := range structData.UnitTestx.PackagesTestDetail {
		test := Test{
			Path: pkgName,
			Time: testRes.Time,
		}
		if len(testRes.Coverage) > 1 {
			test.Cover, _ = strconv.ParseFloat(testRes.Coverage[:(len(testRes.Coverage)-1)], 64)
		}
		if testRes.IsPass {
			test.Result = 1
		}
		testHtmlRes = append(testHtmlRes, test)
	}

	stringTestJson, err := json.Marshal(testHtmlRes)
	if err != nil {
		fmt.Println(err)
	}
	htmlData.Tests = string(stringTestJson)

	// convert cyclo result
	cycloHtmlRes := make([]Cyclo, 0)
	for pkgName, cycloRes := range structData.Cyclox {
		cyclo := Cyclo{
			Pkg:  pkgName,
			Size: len(cycloRes.Result),
		}
		cycloInfos := make([]CycloInfo, 0)
		for _, value := range cycloRes.Result {
			values := strings.Fields(value)
			if len(values) == 4 {
				com, _ := strconv.Atoi(values[0])
				cycloInfo := CycloInfo{
					Comp: com,
					Info: values[3],
				}
				cycloInfos = append(cycloInfos, cycloInfo)
			}
		}
		cyclo.Info = cycloInfos
		cycloHtmlRes = append(cycloHtmlRes, cyclo)
	}

	stringCycloJson, err := json.Marshal(cycloHtmlRes)
	if err != nil {
		fmt.Println(err)
	}
	htmlData.Cyclos = string(stringCycloJson)

	// convert simple code result
	simpleHtmlRes := make([]Simple, 0)
	for _, simpleInfo := range structData.SimpleTips {
		for _, value := range simpleInfo {
			pathIndex := strings.Index(value, ":")
			pathIndexh := strings.Index(value[(pathIndex+1):], ":")
			if pathIndex > 0 && len(value) > (pathIndex+pathIndexh+2) {
				simple := Simple{
					Path: absPath(value[0:(pathIndex + pathIndexh + 1)]),
					Info: value[(pathIndex + pathIndexh + 2):],
				}
				simpleHtmlRes = append(simpleHtmlRes, simple)
			}
		}
	}

	stringSimpleJson, err := json.Marshal(simpleHtmlRes)
	if err != nil {
		fmt.Println(err)
	}
	htmlData.Cyclos = string(stringSimpleJson)

	// convert scan code result
	scanHtmlRes := make([]Scan, 0)
	for _, scanInfo := range structData.ScanTips {
		for _, value := range scanInfo {
			pathIndex := strings.Index(value, ":")
			pathIndexh := strings.Index(value[(pathIndex+1):], ":")
			pathIndexi := strings.Index(value[(pathIndex+pathIndexh+2):], ":")
			if pathIndex > 0 && len(value) > (pathIndex+pathIndexh+2) {
				scan := Scan{
					Path: absPath(value[0:(pathIndex + pathIndexh + 1)]),
					Info: value[(pathIndex + pathIndexh + pathIndexi + 3):],
				}
				scanHtmlRes = append(scanHtmlRes, scan)
			}
		}

	}

	scanSimpleJson, err := json.Marshal(scanHtmlRes)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(scanSimpleJson))
	// scan := string(scanSimpleJson)

	// convert copy code result
	copyHtmlRes := make([]Copycode, 0)
	for _, copys := range structData.CopyTips {
		copyFiles := Copycode{
			Files: strconv.Itoa(len(copys)),
			Path:  copys,
		}
		copyHtmlRes = append(copyHtmlRes, copyFiles)
	}

	stringCopyJson, err := json.Marshal(copyHtmlRes)
	if err != nil {
		fmt.Println(err)
	}
	htmlData.Copycodes = string(stringCopyJson)
	return htmlData, nil
}
