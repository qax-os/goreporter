package engine

import (
	"math"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/goreporter/linters/cyclo"
	"github.com/360EntSecGroup-Skylar/goreporter/utils"
)

type StrategyCyclo struct {
	Sync            *Synchronizer `inject:""`
	compBigThan15   int
	sumAverageCyclo float64
	allDirs         map[string]string
}

func (s *StrategyCyclo) GetName() string {
	return "Cyclo"
}

func (s *StrategyCyclo) GetDescription() string {
	return "Computing all [.go] file's cyclo,and as an important indicator of the quality of the code."
}

func (s *StrategyCyclo) GetWeight() float64 {
	return 0.2
}

func (s *StrategyCyclo) Compute(parameters StrategyParameter) (summaries *Summaries) {
	summaries = NewSummaries()

	s.allDirs = parameters.AllDirs

	sumProcessNumber := int64(10)
	processUnit := utils.GetProcessUnit(sumProcessNumber, len(s.allDirs))

	for pkgName, pkgPath := range s.allDirs {
		errSlice := make([]Error, 0)

		cyclos, avg := cyclo.Cyclo(pkgPath, parameters.ExceptPackages)
		average, _ := strconv.ParseFloat(avg, 64)
		if math.IsNaN(average) == false {
			s.sumAverageCyclo = s.sumAverageCyclo + average
		} else {
			average = 0
		}

		for _, val := range cyclos {
			cyclovalues := strings.Split(val, " ")
			if len(cyclovalues) == 4 {
				comp, _ := strconv.Atoi(cyclovalues[0])
				erroru := Error{
					LineNumber:  comp,
					ErrorString: utils.AbsPath(cyclovalues[3]),
				}
				if comp >= 15 {
					s.compBigThan15 = s.compBigThan15 + 1
				}
				errSlice = append(errSlice, erroru)
			}
		}
		summaries.Lock()
		summaries.Summaries[pkgName] = Summary{
			Name:   pkgName,
			Errors: errSlice,
			Avg:    average,
		}
		summaries.Unlock()
		if sumProcessNumber > 0 {
			s.Sync.LintersProcessChans <- processUnit
			sumProcessNumber = sumProcessNumber - processUnit
		}
	}
	return
}

func (s *StrategyCyclo) Percentage(summaries *Summaries) float64 {
	return utils.CountPercentage(s.compBigThan15 + int(s.sumAverageCyclo/float64(len(s.allDirs))) - 1)
}
