package engine

import (
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/goreporter/linters/spellcheck"
	"github.com/360EntSecGroup-Skylar/goreporter/utils"
)

type StrategySpellCheck struct {
	Sync *Synchronizer `inject:""`
}

func (s *StrategySpellCheck) GetName() string {
	return "SpellCheck"
}

func (s *StrategySpellCheck) GetDescription() string {
	return "Check the project variables, functions, etc. naming spelling is wrong."
}

func (s *StrategySpellCheck) GetWeight() float64 {
	return 0.05
}

func (s *StrategySpellCheck) Compute(parameters StrategyParameter) (summaries *Summaries) {
	summaries = NewSummaries()

	spelltips := spellcheck.SpellCheck(parameters.ProjectPath, parameters.ExceptPackages)
	sumProcessNumber := int64(10)
	processUnit := utils.GetProcessUnit(sumProcessNumber, len(spelltips))

	for _, simpleTip := range spelltips {
		simpleTips := strings.Split(simpleTip, ":")
		if len(simpleTips) == 4 {
			packageName := utils.PackageNameFromGoPath(simpleTips[0])
			line, _ := strconv.Atoi(simpleTips[1])
			erroru := Error{
				LineNumber:  line,
				ErrorString: utils.AbsPath(simpleTips[0]) + ":" + strings.Join(simpleTips[1:], ":"),
			}
			summaries.Lock()
			if summarie, ok := summaries.Summaries[packageName]; ok {
				summarie.Errors = append(summarie.Errors, erroru)
				summaries.Summaries[packageName] = summarie
			} else {
				summarie := Summary{
					Name:   utils.PackageAbsPathExceptSuffix(simpleTips[0]),
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
	return
}

func (s *StrategySpellCheck) Percentage(summaries *Summaries) float64 {
	summaries.RLock()
	defer summaries.RUnlock()
	return utils.CountPercentage(len(summaries.Summaries))
}
