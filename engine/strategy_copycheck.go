package engine

import (
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/goreporter/linters/copycheck"
	"github.com/360EntSecGroup-Skylar/goreporter/utils"
)

type StrategyCopyCheck struct {
	Sync *Synchronizer `inject:""`
}

func (s *StrategyCopyCheck) GetName() string {
	return "CopyCheck"
}

func (s *StrategyCopyCheck) GetDescription() string {
	return "Query all duplicate code in the project and give duplicate code locations and rows."
}

func (s *StrategyCopyCheck) GetWeight() float64 {
	return 0.05
}

// linterCopy provides a function that scans all duplicate code in the project and give
// duplicate code locations and rows.It will extract from the linter need to convert the
// data.The result will be saved in the r's attributes.
func (s *StrategyCopyCheck) Compute(parameters StrategyParameter) (summaries *Summaries) {
	summaries = NewSummaries()

	copyCodeList := copycheck.CopyCheck(parameters.ProjectPath, parameters.ExceptPackages+",_test.go")
	sumProcessNumber := int64(7)
	processUnit := utils.GetProcessUnit(sumProcessNumber, len(copyCodeList))

	for i := 0; i < len(copyCodeList); i++ {
		errorSlice := make([]Error, 0)
		for j := 0; j < len(copyCodeList[i]); j++ {
			line := 0
			values := strings.Split(copyCodeList[i][j], ":")
			if len(values) > 1 {
				lines := strings.Split(strings.TrimSpace(values[1]), ",")
				if len(lines) == 2 {
					lineright, _ := strconv.Atoi(lines[1])
					lineleft, _ := strconv.Atoi(lines[0])
					if lineright-lineleft >= 0 {
						line = lineright - lineleft + 1
					}
				}
				values[0] = utils.AbsPath(values[0])
			}

			errorSlice = append(errorSlice, Error{LineNumber: line, ErrorString: strings.Join(values, ":")})
		}
		summaries.Lock()
		summaries.Summaries[string(i)] = Summary{
			Name:   strconv.Itoa(len(errorSlice)),
			Errors: errorSlice,
		}
		summaries.Unlock()
		if sumProcessNumber > 0 {
			s.Sync.LintersProcessChans <- processUnit
			sumProcessNumber = sumProcessNumber - processUnit
		}
	}
	return
}

func (s *StrategyCopyCheck) Percentage(summaries *Summaries) float64 {
	summaries.RLock()
	defer summaries.RUnlock()
	return utils.CountPercentage(len(summaries.Summaries))
}
