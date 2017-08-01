package engine

import (
	"github.com/360EntSecGroup-Skylar/goreporter/linters/interfacer"
	"strings"
	"strconv"
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
	return 0.06
}

// linterInterfacer is a function that scan the interface of all packages in the
// project helps you optimize the project architecture.It will extract from the
// linter need to convert the data.The result will be saved in the r's attributes.
func (s *StrategyInterfacer) Compute(parameters StrategyParameter) (summaries Summaries) {
	summaries = NewSummaries()

	interfacers := interfacer.Interfacer(parameters.AllDirs)
	sumProcessNumber := int64(5)
	processUnit := GetProcessUnit(sumProcessNumber, len(interfacers))
	for _, interfaceTip := range interfacers {
		interfaceTips := strings.Split(interfaceTip, ":")
		if len(interfaceTips) == 4 {
			packageName := PackageNameFromGoPath(interfaceTips[0])
			line, _ := strconv.Atoi(interfaceTips[1])
			erroru := Error{
				LineNumber:  line,
				ErrorString: AbsPath(interfaceTips[0]) + ":" + strings.Join(interfaceTips[1:], ":"),
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
	return
}

func (s *StrategyInterfacer) Percentage(summaries Summaries) float64 {
	return 0.
}


