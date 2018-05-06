package engine

import (
	"github.com/360EntSecGroup-Skylar/goreporter/linters/depend"
)

type StrategyDependGraph struct {
	Sync *Synchronizer `inject:""`
}

func (s *StrategyDependGraph) GetName() string {
	return "DependGraph"
}

func (s *StrategyDependGraph) GetDescription() string {
	return "The dependency graph for all packages in the project helps you optimize the project architecture."
}

func (s *StrategyDependGraph) GetWeight() float64 {
	return 0.
}

// linterDependGraph is a function that builds the dependency graph of all packages in the
// project helps you optimize the project architecture.It will extract from the linter need
// to convert the data.The result will be saved in the r's attributes.
func (s *StrategyDependGraph) Compute(parameters StrategyParameter) (summaries *Summaries) {
	summaries = NewSummaries()

	graph := depend.Depend(parameters.ProjectPath, parameters.ExceptPackages)
	summaries.Summaries["graph"] = Summary{
		Name:        s.GetName(),
		Description: graph,
	}

	return
}

func (s *StrategyDependGraph) Percentage(summaries *Summaries) float64 {
	return 0.
}
