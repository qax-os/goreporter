// Copyright 2017 The GoReporter Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package engine

import "sync"

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
// data.But may have some difference in different linter.
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

	syncRW *sync.RWMutex
}
