package engine

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/goreporter/linters/govet"
	"github.com/360EntSecGroup-Skylar/goreporter/utils"
)

type StrategyGoVet struct {
	Sync *Synchronizer `inject:""`
}

func (s *StrategyGoVet) GetName() string {
	return "GoVet"
}

func (s *StrategyGoVet) GetDescription() string {
	return "go vet examines Go source code and reports suspicious constructs, such as Printf calls whose arguments do not align with the format string."
}

func (s *StrategyGoVet) GetWeight() float64 {
	return 0.1
}

func (s *StrategyGoVet) Compute(parameters StrategyParameter) (summaries *Summaries) {
	summaries = NewSummaries()
	slicePackagePaths := make([]string, 0)
	for _, packagePath := range parameters.AllDirs {
		slicePackagePaths = append(slicePackagePaths, packagePath)
	}
	lints, err := govet.GoVet(slicePackagePaths)
	if err != nil {
		fmt.Println(err)
	}
	sumProcessNumber := int64(10)
	processUnit := utils.GetProcessUnit(sumProcessNumber, len(lints))
	for _, lintTip := range lints {
		lintTips := strings.Split(lintTip, ":")
		if len(lintTips) == 4 {
			packageName := utils.PackageNameFromGoPath(lintTips[0])
			line, _ := strconv.Atoi(lintTips[1])
			erroru := Error{
				LineNumber:  line,
				ErrorString: utils.AbsPath(lintTips[0]) + ":" + strings.Join(lintTips[1:], ":"),
			}
			summaries.Lock()
			if summarie, ok := summaries.Summaries[packageName]; ok {
				summarie.Errors = append(summarie.Errors, erroru)
				summaries.Summaries[packageName] = summarie
			} else {
				summarie := Summary{
					Name:   packageName,
					Errors: make([]Error, 0),
				}
				summarie.Errors = append(summarie.Errors, erroru)
				summaries.Summaries[packageName] = summarie
			}
			summaries.Unlock()
		}
		if sumProcessNumber > 0 {
			s.Sync.LintersProcessChans <- processUnit
			sumProcessNumber = sumProcessNumber - processUnit
		}
	}

	return summaries
}

func (s *StrategyGoVet) Percentage(summaries *Summaries) float64 {
	summaries.RLock()
	defer summaries.RUnlock()
	return utils.CountPercentage(len(summaries.Summaries))
}
