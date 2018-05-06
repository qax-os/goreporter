package engine

import (
	"fmt"
	"strconv"

	"github.com/360EntSecGroup-Skylar/goreporter/linters/countcode"
	"github.com/360EntSecGroup-Skylar/goreporter/utils"
)

type StrategyCountCode struct {
	Sync *Synchronizer `inject:""`
}

func (s *StrategyCountCode) GetName() string {
	return "CountCode"
}

func (s *StrategyCountCode) GetDescription() string {
	return "Count lines and files of go project."
}

func (s *StrategyCountCode) GetWeight() float64 {
	return 0.
}

// linterCount is a function that counts go files and go code lines of
// project.It will extract from the linter need to convert the data.
// The result will be saved in the r's attributes.
func (s *StrategyCountCode) Compute(parameters StrategyParameter) (summaries *Summaries) {
	summaries = NewSummaries()

	codeCounts := countcode.CountCode(parameters.ProjectPath, parameters.ExceptPackages)
	for packageName, codeCount := range codeCounts {
		if len(codeCount) == 4 {
			absFilePath := utils.AbsPath(packageName)
			summaries.Summaries[absFilePath] = Summary{
				Name:        absFilePath,
				Description: fmt.Sprintf("%s;%s;%s;%s", strconv.Itoa(codeCount[0]), strconv.Itoa(codeCount[1]), strconv.Itoa(codeCount[2]), strconv.Itoa(codeCount[3])),
			}
		}
	}
	return
}

func (s *StrategyCountCode) Percentage(summaries *Summaries) float64 {
	return 0.
}
