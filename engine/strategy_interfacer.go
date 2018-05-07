package engine

import (
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/goreporter/linters/interfacer"
	"github.com/360EntSecGroup-Skylar/goreporter/utils"
)

type StrategyInterfacer struct {
	Sync *Synchronizer `inject:""`
}

func (s *StrategyInterfacer) GetName() string {
	return "Interfacer"
}

func (s *StrategyInterfacer) GetDescription() string {
	return "Suggests interface types. In other words, it warns about the usage of types that are more specific than necessary."
}

func (s *StrategyInterfacer) GetWeight() float64 {
	return 0.05
}

// linterInterfacer is a function that scan the interface of all packages in the
// project helps you optimize the project architecture.It will extract from the
// linter need to convert the data.The result will be saved in the r's attributes.
func (s *StrategyInterfacer) Compute(parameters StrategyParameter) (summaries *Summaries) {
	summaries = NewSummaries()

	interfacers := interfacer.Interfacer(parameters.AllDirs)
	sumProcessNumber := int64(5)
	processUnit := utils.GetProcessUnit(sumProcessNumber, len(interfacers))
	for _, interfaceTip := range interfacers {
		interfaceTips := strings.Split(interfaceTip, ":")
		if len(interfaceTips) == 4 {
			packageName := utils.PackageNameFromGoPath(interfaceTips[0])
			line, _ := strconv.Atoi(interfaceTips[1])
			erroru := Error{
				LineNumber:  line,
				ErrorString: utils.AbsPath(interfaceTips[0]) + ":" + strings.Join(interfaceTips[1:], ":"),
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
	return
}

func (s *StrategyInterfacer) Percentage(summaries *Summaries) float64 {
	return 0.
}
