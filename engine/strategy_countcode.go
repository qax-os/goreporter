package engine

import (
	"strconv"

	"github.com/360EntSecGroup-Skylar/goreporter/linters/countcode"
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
func (s *StrategyCountCode) Compute(parameters StrategyParameter) (summaries Summaries) {
	summaries = NewSummaries()

	fileCount, codeLines, commentLines, totalLines := countcode.CountCode(parameters.ProjectPath, parameters.ExceptPackages)
	summaries.Summaries["FileCount"] = Summary{
		Name:        "FileCount",
		Description: strconv.Itoa(fileCount),
	}
	summaries.Summaries["CodeLines"] = Summary{
		Name:        "CodeLines",
		Description: strconv.Itoa(codeLines),
	}
	summaries.Summaries["CommentLines"] = Summary{
		Name:        "CommentLines",
		Description: strconv.Itoa(commentLines),
	}
	summaries.Summaries["TotalLines"] = Summary{
		Name:        "TotalLines",
		Description: strconv.Itoa(totalLines),
	}
	// todo:get all package count
	return
}

func (s *StrategyCountCode) Percentage(summaries Summaries) float64 {
	return 0.
}
