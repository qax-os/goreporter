package engine

type StrategyLinter interface {
	Compute(parameters StrategyParameter) *Summaries
	Percentage(summaries *Summaries) float64
	GetName() string
	GetDescription() string
	GetWeight() float64
}
