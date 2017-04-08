package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	// "github.com/wgliang/goreporter/linters/aligncheck"
	"github.com/wgliang/goreporter/linters/copycheck"
	"github.com/wgliang/goreporter/linters/cyclo"
	"github.com/wgliang/goreporter/linters/deadcode"
	"github.com/wgliang/goreporter/linters/depend"
	// "github.com/wgliang/goreporter/linters/errorcheck"
	"github.com/wgliang/goreporter/linters/simplecode"
	"github.com/wgliang/goreporter/linters/staticscan"
	// "github.com/wgliang/goreporter/linters/structcheck"
	"github.com/wgliang/goreporter/linters/spellcheck"
	"github.com/wgliang/goreporter/linters/unittest"
	// "github.com/wgliang/goreporter/linters/varcheck"
)

var (
	tpl string
)

func NewReporter() *Reporter {
	return &Reporter{}
}

func (r *Reporter) Engine(projectPath string, exceptPackages string) {
	log.Println("start code quality assessment...")

	dirsUnitTest, err := DirList(projectPath, "_test.go", exceptPackages)
	if err != nil {
		log.Println(err)
	}
	r.Project = projectName(projectPath)
	var wg sync.WaitGroup

	// run linter:unit test
	wg.Add(1)
	go func() {
		log.Println("running unit test...")
		packagesTestDetail := struct {
			Values map[string]PackageTest
			mux    *sync.RWMutex
		}{make(map[string]PackageTest, 0), new(sync.RWMutex)}
		packagesRaceDetail := struct {
			Values map[string][]string
			mux    *sync.RWMutex
		}{make(map[string][]string, 0), new(sync.RWMutex)}

		sumCover := 0.0
		countCover := 0
		var pkg sync.WaitGroup
		for pkgName, pkgPath := range dirsUnitTest {
			pkg.Add(1)
			go func(pkgName, pkgPath string) {
				unitTestRes, unitRaceRes := unittest.UnitTest(pkgPath)
				var packageTest PackageTest
				if len(unitTestRes) >= 1 {
					testres := unitTestRes[pkgName]
					if len(testres) > 5 {
						if testres[0] == "ok" {
							packageTest.IsPass = true
						} else {
							packageTest.IsPass = false
						}
						timeLen := len(testres[2])
						if timeLen > 1 {
							time, err := strconv.ParseFloat(testres[2][:(timeLen-1)], 64)
							if err == nil {
								packageTest.Time = time
							} else {
								log.Println(err)
							}
						}
						packageTest.Coverage = testres[4]

						coverLen := len(testres[4])
						if coverLen > 1 {
							coverFloat, _ := strconv.ParseFloat(testres[4][:(coverLen-1)], 64)
							sumCover = sumCover + coverFloat
							countCover = countCover + 1
						} else {
							countCover = countCover + 1
						}
					} else {
						packageTest.Coverage = "0%"
						countCover = countCover + 1
					}
				} else {
					packageTest.Coverage = "0%"
					countCover = countCover + 1
				}
				packagesTestDetail.mux.Lock()
				packagesTestDetail.Values[pkgName] = packageTest
				packagesTestDetail.mux.Unlock()

				if len(unitRaceRes[pkgName]) > 0 {
					packagesRaceDetail.mux.Lock()
					packagesRaceDetail.Values[pkgName] = unitRaceRes[pkgName]
					packagesRaceDetail.mux.Unlock()
				}
				pkg.Done()
			}(pkgName, pkgPath)
		}

		pkg.Wait()
		packagesTestDetail.mux.Lock()
		r.UnitTestx.PackagesTestDetail = packagesTestDetail.Values
		packagesTestDetail.mux.Unlock()
		r.UnitTestx.AvgCover = fmt.Sprintf("%.1f", sumCover/float64(countCover)) + "%"
		packagesRaceDetail.mux.Lock()
		r.UnitTestx.PackagesRaceDetail = packagesRaceDetail.Values
		packagesRaceDetail.mux.Unlock()

		wg.Done()
		log.Println("unit test over!")
	}()

	log.Println("computing cyclo...")
	dirsAll, err := DirList(projectPath, ".go", exceptPackages)
	if err != nil {
		log.Println(err)
	}
	wg.Add(1)
	go func() {
		cycloRes := make(map[string]Cycloi, 0)
		for pkgName, pkgPath := range dirsAll {
			cyclo, avg := cyclo.Cyclo(pkgPath)
			cycloRes[pkgName] = Cycloi{
				Average: avg,
				Result:  cyclo,
			}
		}
		r.Cyclox = cycloRes
		wg.Done()
		log.Println("cyclo over!")
	}()

	log.Println("simpling code...")
	wg.Add(1)
	go func() {
		simples := simplecode.SimpleCode(projectPath)
		simpleTips := make(map[string][]string, 0)
		for _, tips := range simples {
			index := strings.Index(tips, ":")
			simpleTips[PackageAbsPathExceptSuffix(tips[0:index])] = append(simpleTips[PackageAbsPathExceptSuffix(tips[0:index])], tips)
		}
		r.SimpleTips = simpleTips
		wg.Done()
	}()
	log.Println("simpled code!")

	log.Println("checking copy code...")
	wg.Add(1)
	go func() {
		x := copycheck.CopyCheck(projectPath, "_test.go")
		r.CopyTips = x
		wg.Done()
		log.Println("checked copy code!")
	}()

	log.Println("running staticscan...")
	wg.Add(1)
	go func() {
		staticscan.StaticScan(projectPath)
		scanTips := make(map[string][]string, 0)
		tips := staticscan.StaticScan(projectPath)
		for _, tip := range tips {
			index := strings.Index(tip, ":")
			scanTips[PackageAbsPathExceptSuffix(tip[0:index])] = append(scanTips[PackageAbsPathExceptSuffix(tip[0:index])], tip)
		}
		r.ScanTips = scanTips
		wg.Done()
		log.Println("staticscan over!")
	}()

	log.Println("creating depend graph...")
	wg.Add(1)
	go func() {
		r.DependGraph = depend.Depend(projectPath, exceptPackages)
		wg.Done()
		log.Println("created depend graph")
	}()

	log.Println("checking dead code...")
	wg.Add(1)
	go func() {
		r.DeadCode = deadcode.DeadCode(projectPath)
		wg.Done()
		log.Println("checked dead code")
	}()

	log.Println("checking spell error...")
	wg.Add(1)
	go func() {
		r.SpellError = spellcheck.SpellCheck(projectPath, exceptPackages)
		wg.Done()
		log.Println("checked spell error")
	}()

	log.Println("getting import packages...")
	var importPkgs []string
	wg.Add(1)
	go func() {
		importPkgs = unittest.GoListWithImportPackages(projectPath)
		wg.Done()
	}()

	wg.Wait()

	// get all no unit test packages
	noTestPackage := make([]string, 0)
	for i := 0; i < len(importPkgs); i++ {
		if _, ok := r.UnitTestx.PackagesTestDetail[importPkgs[i]]; !ok {
			noTestPackage = append(noTestPackage, importPkgs[i])
		}
	}
	r.NoTestPkg = noTestPackage

	log.Println("finished code quality assessment...")
}

func (r *Reporter) formateReport2Json() []byte {
	report, err := json.Marshal(r)
	if err != nil {
		log.Println("json err:", err)
	}

	return report
}

func (r *Reporter) SaveAsHtml(htmlData HtmlData, projectPath, savePath, timestamp string) {
	t, err := template.New("skylar-apollo").Parse(tpl)
	if err != nil {
		log.Println(err)
	}

	var out bytes.Buffer
	err = t.Execute(&out, htmlData)
	if err != nil {
		log.Println(err)
	}
	projectName := projectName(projectPath)
	if savePath != "" {
		htmlpath := strings.Replace(savePath+string(filepath.Separator)+projectName+"-"+timestamp+".html", string(filepath.Separator)+string(filepath.Separator), string(filepath.Separator), -1)
		log.Println(htmlpath)
		err = ioutil.WriteFile(htmlpath, out.Bytes(), 0666)
		if err != nil {
			log.Println(err)
		}
	} else {
		htmlpath := projectName + "-" + timestamp + ".html"
		log.Println(htmlpath)
		err = ioutil.WriteFile(htmlpath, out.Bytes(), 0666)
		if err != nil {
			log.Println(err)
		}
	}
}

func (r *Reporter) Grade() int {
	score := 0.0
	tscore := float64(40)
	if len(r.UnitTestx.AvgCover) > 1 {
		cover, err := strconv.ParseFloat(r.UnitTestx.AvgCover[:(len(r.UnitTestx.AvgCover)-1)], 64)
		if err != nil {
			cover = 0
		}
		score = score + tscore*cover/100.0
	}

	countCopy := len(r.CopyTips)
	if countCopy < 10 {
		score = score + float64(10-1*countCopy)
	}

	countScan := 0
	for _, pkg := range r.ScanTips {
		countScan = countScan + len(pkg)
	}
	if countScan < 10 {
		score = score + float64(10-1*countScan)
	}

	countSimple := 0
	for _, pkg := range r.SimpleTips {
		countSimple = countSimple + len(pkg)
	}
	if countSimple < 10 {
		score = score + float64(10-1*countSimple)
	}

	sscore := 10.0
	sscore = sscore - float64(len(r.DeadCode)/5)
	if sscore < 0 {
		sscore = 0
	}
	score = score + sscore

	sum15 := 0
	sum50 := 0
	countcyclo := 0
	sum := 0
	pscore := 0
	for _, val := range r.Cyclox {
		for _, v := range val.Result {
			var num int
			in := strings.Index(v, " ")
			if in > 0 {
				countcyclo++
				num, _ = strconv.Atoi(v[0:in])
				if num >= 50 {
					sum50++
					sum15++
				} else if num >= 15 {
					sum15++
				} else {
					sum += num
				}
			}
		}
	}

	if (countcyclo - sum50 - sum15) > 0 {
		pscore = 20 * ((15 * 1.0 * (countcyclo - sum50 - sum15)) - sum) / (15 * (countcyclo - sum50 - sum15))
	} else {
		pscore = 0
	}

	pscore = pscore - sum50/5 - sum15/10
	if pscore < 0 {
		pscore = 0
	}
	score = score + float64(pscore)

	return int(score)
}

func (r *Reporter) SaveAsJson(projectPath, savePath, timestamp string) {
	jsonData := r.formateReport2Json()
	savePath = absPath(savePath)
	projectName := projectName(projectPath)
	if savePath != "" {
		jsonpath := strings.Replace(savePath+string(filepath.Separator)+projectName+"-"+timestamp+".json", string(filepath.Separator)+string(filepath.Separator), string(filepath.Separator), -1)
		err := ioutil.WriteFile(jsonpath, jsonData, 0666)
		if err != nil {
			log.Println(err)
		}
	} else {
		jsonpath := projectName + "-" + timestamp + ".json"
		err := ioutil.WriteFile(jsonpath, jsonData, 0666)
		if err != nil {
			log.Println(err)
		}
	}
}
