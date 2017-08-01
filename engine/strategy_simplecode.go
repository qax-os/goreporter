package engine

import (
	"strconv"
	"strings"
	"github.com/360EntSecGroup-Skylar/goreporter/linters/simplecode"
)

type StrategySimpleCode struct {
	Sync *Synchronizer `inject:""`
}

func (s *StrategySimpleCode) GetName() string {
	return "Simple"
}

func (s *StrategySimpleCode) GetDescription() string {
	return "All golang code hints that can be optimized and give suggestions for changes."
}

func (s *StrategySimpleCode) GetWeight() float64 {
	return 0.1
}

func (s *StrategySimpleCode) Compute(parameters StrategyParameter) (summaries Summaries) {
	summaries = NewSummaries()

	simples := simplecode.Simple(parameters.AllDirs)
	sumProcessNumber := int64(10)
	processUnit := GetProcessUnit(sumProcessNumber, len(simples))
	for _, simpleTip := range simples {
		simpleTips := strings.Split(simpleTip, ":")
		if len(simpleTips) == 4 {
			packageName := PackageNameFromGoPath(simpleTips[0])
			line, _ := strconv.Atoi(simpleTips[1])
			erroru := Error {
				LineNumber:  line,
				ErrorString: AbsPath(simpleTips[0]) + ":" + strings.Join(simpleTips[1:], ":"),
			}
			if summarie, ok := summaries[packageName]; ok {
				summarie.Errors = append(summarie.Errors, erroru)
				summaries[packageName] = summarie
			} else {
				summarie := Summary{
					Name:   packageName,
					Errors: make([]Error, 0),
				}
				summarie.Errors = append(summarie.Errors, erroru)
				summaries[packageName] = summarie
			}
		}
		if sumProcessNumber > 0 {
			s.Sync.LintersProcessChans <- processUnit
			sumProcessNumber = sumProcessNumber - processUnit
		}
	}

	return summaries
}

func (s *StrategySimpleCode) Percentage(summaries Summaries) float64 {
	return CountPercentage(len(summaries))
}
