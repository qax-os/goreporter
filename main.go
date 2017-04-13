package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
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
	log.SetPrefix("[Apollo]")
	if *project == "" {
		log.Fatal("The project path is not specified")
	} else {
		_, err := os.Stat(*project)
		if err != nil {
			log.Fatal("project path is invalid")
		}
	}

	if *tplpath == "" {
		log.Println("The template path is not specified,and will use the default template")
	} else {
		if !strings.HasSuffix(*report, ".html") {
			log.Println("The template file is not a html template")
		}
		fileData, err := ioutil.ReadFile(*tplpath)
		if err != nil {
			log.Fatal(err)
		} else {
			tpl = string(fileData)
		}
	}

	if *report == "" {
		log.Println("The report path is not specified, and the current path is used by default")
	} else {
		_, err := os.Stat(*report)
		if err != nil {
			log.Fatal("report path is invalid")
		}
	}

	if *except == "" {
		log.Println("There are no packages that are excepted, review all items of the package")
	}

	startTime := strconv.FormatInt(time.Now().Unix(), 10)
	reporter := NewReporter()
	reporter.Engine(*project, *except)
	htmlData, err := reporter.Json2Html()
	if err != nil {
		log.Println("Json2Html error")
		return
	}
	if *formate == "json" {
		reporter.SaveAsJson(*project, *report, startTime)
	} else {
		reporter.SaveAsHtml(htmlData, *project, *report, startTime)
	}
}
