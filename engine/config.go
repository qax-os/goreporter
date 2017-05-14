package engine

// Error contains the line number and the reason for
// an error output from a command
type Error struct {
	LineNumber  int    `json:"line_number"`
	ErrorString string `json:"error_string"`
}

// FileSummary contains the filename, location of the file
// on GitHub, and all of the errors related to the file
type Summary struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Errors      []Error `json:"errors"`
}

// Metric as template of report and will save all linters result
// data.But may have some diffreence in diffrent linter.
type Metric struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Summaries   map[string]Summary `json:"summaries"`
	Weight      float64            `json:"weight"`
	Percentage  float64            `json:"percentage"`
	Error       string             `json:"error"`
}

// Reporter is the top struct of GoReporter.
type Reporter struct {
	Project   string            `json:"project"`
	Score     int               `json:"score"`
	Grade     int               `json:"grade"`
	Metrics   map[string]Metric `json:"metrics"`
	Issues    int               `json:"issues"`
	TimeStamp string            `json:"time_stamp"`
}
