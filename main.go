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

// GoReporter is a Golang tool that does static analysis, unit testing, code
// review and generate code quality report.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/goreporter/engine"
	"github.com/360EntSecGroup-Skylar/goreporter/tools"
)

// Received parameters, you can control some features using:
//
// -p:Specify the relative path of your project(Must Be Relative path),
//    by default, the current path is used
// -r:Specifies the save path for the generated report,
//    by default, the current path is used
// -e:Ignored detection of packages and multiple packages separated by commas.
// -t:Customize the path of the report template, not necessarily using the
//    default report template
// -f:Set the format to generate reports, support text, html and json,not
//    necessarily using the default formate-html.

var (
	projectPath     = flag.String("p", "", "path of project.")
	reportPath      = flag.String("r", "", "path of report.")
	exceptPackages  = flag.String("e", "", "except packages.")
	templatePath    = flag.String("t", "", "report html template path.")
	formateOfReport = flag.String("f", "", "project report formate(text/json/html).")
)

func main() {
	flag.Parse()
	if *projectPath == "" {
		log.Fatal("The project path is not specified")
	} else {
		_, err := os.Stat(*projectPath)
		if err != nil {
			log.Fatal("project path is invalid")
		}
	}

	var templateHtml string
	if *templatePath == "" {
		log.Println("The template path is not specified,and will use the default template")
	} else {
		if !strings.HasSuffix(*templatePath, ".html") {
			log.Println("The template file is not a html template")
		}
		fileData, err := ioutil.ReadFile(*templatePath)
		if err != nil {
			log.Fatal(err)
		} else {
			templateHtml = string(fileData)
		}
	}

	if *reportPath == "" {
		log.Println("The report path is not specified, and the current path is used by default")
	} else {
		_, err := os.Stat(*reportPath)
		if err != nil {
			log.Fatal("report path is invalid:", err)
		}
	}

	if *exceptPackages == "" {
		log.Println("There are no packages that are excepted, review all items of the package")
	}
	// Displaying linters process bar.
	lintersProcessChans := make(chan int64, 20)
	lintersFinishedSignal := make(chan string, 10)
	go tools.LinterProcessBar(lintersProcessChans, lintersFinishedSignal)
	start := time.Now()
	startTime := strconv.FormatInt(start.Unix(), 10)
	reporter := engine.NewReporter()
	reporter.Engine(*projectPath, *exceptPackages, lintersProcessChans, lintersFinishedSignal, start)
	jsonData := reporter.FormateReport2Json()

	if *formateOfReport == "json" {
		log.Println(fmt.Sprintf("Generating json report,time consuming %vs", time.Now().Sub(start).Seconds()))
		tools.SaveAsJson(jsonData, *projectPath, *reportPath, startTime)
	} else if *formateOfReport == "text" {
		tools.DisplayAsText(jsonData)
	} else {
		log.Println(fmt.Sprintf("Generating HTML report,time consuming %vs", time.Now().Sub(start).Seconds()))
		htmlData, err := tools.Json2Html(jsonData)
		if err != nil {
			log.Println("Json2Html error:", err)
			return
		}
		tools.SaveAsHtml(htmlData, *projectPath, *reportPath, startTime, templateHtml)
	}
	log.Println(fmt.Sprintf("GoReporter Finished,time consuming %vs", time.Now().Sub(start).Seconds()))
}
