package engine

import (
	"path/filepath"
	"strconv"
	"sync"

	"github.com/golang/glog"
	"github.com/json-iterator/go"

	"github.com/360EntSecGroup-Skylar/goreporter/linters/unittest"
	"github.com/360EntSecGroup-Skylar/goreporter/utils"
)

type StrategyUnitTest struct {
	Sync       *Synchronizer `inject:""`
	sumCover   float64
	countCover int
}

func (s *StrategyUnitTest) GetName() string {
	return "UnitTest"
}

func (s *StrategyUnitTest) GetDescription() string {
	return "Run all valid TEST in your golang package.And will measure from both coverage and time-consuming."
}

func (s *StrategyUnitTest) GetWeight() float64 {
	return 0.3
}

func (s *StrategyUnitTest) Compute(parameters StrategyParameter) (summaries *Summaries) {
	summaries = NewSummaries()

	sumProcessNumber := int64(30)
	processUnit := utils.GetProcessUnit(sumProcessNumber, len(parameters.UnitTestDirs))

	var pkg sync.WaitGroup

	for pkgName, pkgPath := range parameters.UnitTestDirs {
		pkg.Add(1)
		go func(pkgName, pkgPath string) {
			unitTestRes, _ := unittest.UnitTest("." + string(filepath.Separator) + pkgPath)
			var packageTest PackageTest
			if len(unitTestRes) >= 5 {
				if unitTestRes[0] == "ok" {
					packageTest.IsPass = true
				} else {
					packageTest.IsPass = false
				}
				timeLen := len(unitTestRes[2])
				if timeLen > 1 {
					t, err := strconv.ParseFloat(unitTestRes[2][:(timeLen-1)], 64)
					if err == nil {
						packageTest.Time = t
					} else {
						glog.Errorln(err)
					}
				}
				packageTest.Coverage = unitTestRes[4]

				coverLen := len(unitTestRes[4])
				if coverLen > 1 {
					coverFloat, _ := strconv.ParseFloat(unitTestRes[4][:(coverLen-1)], 64)
					s.sumCover = s.sumCover + coverFloat
				}
				s.countCover++
			} else {
				packageTest.Coverage = "0%"
				s.countCover++
			}
			jsonStringPackageTest, err := jsoniter.Marshal(packageTest)
			if err != nil {
				glog.Errorln(err)
			}
			summaries.Lock()
			summaries.Summaries[pkgName] = Summary{
				Name:        pkgName,
				Description: string(jsonStringPackageTest),
			}
			summaries.Unlock()
			pkg.Done()
		}(pkgName, pkgPath)
		if sumProcessNumber > 0 {
			s.Sync.LintersProcessChans <- processUnit
			sumProcessNumber = sumProcessNumber - processUnit
		}
	}

	pkg.Wait()

	return
}

func (s *StrategyUnitTest) Percentage(summaries *Summaries) float64 {
	if s.countCover == 0 {
		return 0.0
	} else {
		return s.sumCover / float64(s.countCover)
	}
}
