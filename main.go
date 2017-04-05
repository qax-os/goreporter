package main

import (
	"flag"
	"io/ioutil"
	"log"
	"runtime"
	"strconv"
	"time"
)

// receive parameters
var (
	// project path:Must Be Relative path
	project = flag.String("p", "", "path of project.")
	// save path of report
	report = flag.String("d", "", "path of report.")
	// except packages,multiple packages are separated by semicolons
	except = flag.String("e", "", "except packages.")
	// template
	tplpath = flag.String("t", "", "project meta information.")
	// report formate
	formate = flag.String("r", "", "project report formate(json/html).")
)

func init() {
	if runtime.GOOS == `windows` {
		system = `\`
	} else {
		system = `/`
	}
}

func main() {
	flag.Parse()
	log.SetPrefix("[Apollo]")
	if *project == "" {
		log.Fatal("The project path is not specified")
	}

	if *tplpath == "" {
		log.Fatal("The template path is not specified")
	} else {
		fileData, err := ioutil.ReadFile(*tplpath)
		if err != nil {
			log.Fatal(err)
		} else {
			tpl = string(fileData)
		}
	}

	if *report == "" {
		log.Println("The report path is not specified, and the current path is used by default")
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
