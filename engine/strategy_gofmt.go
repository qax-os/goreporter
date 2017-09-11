package engine

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/goreporter/linters/gofmt"
)

type StrategyGoFmt struct {
	Sync *Synchronizer `inject:""`
}

func (s *StrategyGoFmt) GetName() string {
	return "GoFmt"
}

func (s *StrategyGoFmt) GetDescription() string {
	return `go fmt formats Go programs. We run gofmt -s on your code, where -s is for the "simplify" command.`
}

func (s *StrategyGoFmt) GetWeight() float64 {
	return 0.1
}

func (s *StrategyGoFmt) Compute(parameters StrategyParameter) (summaries Summaries) {
	summaries = NewSummaries()
	slicePackagePaths := make([]string, 0)
	for _, packagePath := range parameters.AllDirs {
		slicePackagePaths = append(slicePackagePaths, packagePath)
	}
	lints, err := gofmt.GoFmt(slicePackagePaths)
	if err != nil {
		fmt.Println(err)
	}
	sumProcessNumber := int64(10)
	processUnit := GetProcessUnit(sumProcessNumber, len(lints))
	for _, lintTip := range lints {
		packageName := PackageNameFromGoPath(lintTip)
		erroru := Error{
			LineNumber:  1,
			ErrorString: AbsPath(lintTip) + ":warning: file is not gofmted with -s (gofmt)",
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

		if sumProcessNumber > 0 {
			s.Sync.LintersProcessChans <- processUnit
			sumProcessNumber = sumProcessNumber - processUnit
		}
	}

	return summaries
}

func (s *StrategyGoFmt) Percentage(summaries Summaries) float64 {
	summaries.RLock()
	defer summaries.RUnlock()
	return CountPercentage(len(summaries.Summaries))
}
