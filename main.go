package main

import (
	"flag"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/wgliang/goreporter/engine"
	"github.com/wgliang/goreporter/tools"
)

// receive parameters
var (
	// project path:Must Be Relative path
	project = flag.String("p", "", "path of project.")
	// save path of report
	report = flag.String("r", "", "path of report.")
	// except packages,multiple packages are separated by semicolons
	except = flag.String("e", "", "except packages.")
	// template
	tplpath = flag.String("t", "", "report html template path.")
	// report formate
	formate = flag.String("f", "", "project report formate(json/html).")
)

func main() {
	flag.Parse()
	if *project == "" {
		glog.Fatal("The project path is not specified")
	} else {
		_, err := os.Stat(*project)
		if err != nil {
			glog.Fatal("project path is invalid")
		}
	}

	var templateHtml string
	if *tplpath == "" {
		glog.Warningln("The template path is not specified,and will use the default template")
	} else {
		if !strings.HasSuffix(*report, ".html") {
			glog.Warningln("The template file is not a html template")
		}
		fileData, err := ioutil.ReadFile(*tplpath)
		if err != nil {
			glog.Fatal(err)
		} else {
			templateHtml = string(fileData)
		}
	}

	if *report == "" {
		glog.Warningln("The report path is not specified, and the current path is used by default")
	} else {
		_, err := os.Stat(*report)
		if err != nil {
			glog.Fatal("report path is invalid")
		}
	}

	if *except == "" {
		glog.Warningln("There are no packages that are excepted, review all items of the package")
	}

	startTime := strconv.FormatInt(time.Now().Unix(), 10)
	reporter := engine.NewReporter()
	reporter.Engine(*project, *except)
	jsonData := reporter.FormateReport2Json()

	if *formate == "json" {
		tools.SaveAsJson(jsonData, *project, *report, startTime)
	} else if *formate == "text" {
		tools.DisplayAsText(jsonData)
	} else {
		htmlData, err := tools.Json2Html(jsonData)
		if err != nil {
			glog.Errorln("Json2Html error")
			return
		}
		tools.SaveAsHtml(htmlData, *project, *report, startTime, templateHtml)
	}
}
