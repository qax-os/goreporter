package engine

import (
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/goreporter/linters/golint"
)

type StrategyLint struct {
	Sync *Synchronizer `inject:""`
}

func (s *StrategyLint) GetName() string {
	return "Lint"
}

func (s *StrategyLint) GetDescription() string {
	return "All golang code hints that can be optimized and give suggestions for changes."
}

func (s *StrategyLint) GetWeight() float64 {
	return 0.1
}

func (s *StrategyLint) Compute(parameters StrategyParameter) (summaries Summaries) {
	summaries = NewSummaries()
	slicePackagePaths := make([]string, 0)
	for _, packagePath := range parameters.AllDirs {
		slicePackagePaths = append(slicePackagePaths, packagePath)
	}
	lints := golint.GoLinter(slicePackagePaths)
	sumProcessNumber := int64(10)
	processUnit := GetProcessUnit(sumProcessNumber, len(lints))
	for _, lintTip := range lints {
		lintTips := strings.Split(lintTip, ":")
		if len(lintTips) == 4 {
			packageName := PackageNameFromGoPath(lintTips[0])
			line, _ := strconv.Atoi(lintTips[1])
			erroru := Error{
				LineNumber:  line,
				ErrorString: AbsPath(lintTips[0]) + ":" + strings.Join(lintTips[1:], ":"),
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

func (s *StrategyLint) Percentage(summaries Summaries) float64 {
	summaries.RLock()
	defer summaries.RUnlock()
	return CountPercentage(len(summaries.Summaries))
}
