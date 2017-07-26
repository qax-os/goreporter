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

package tools

import (
	"html/template"
	"time"
)

// UnitTest is a struct that contains AvgCover, PackagesTestDetail and
// PackagesRaceDetail. The type of AvgCover MUST string that represents
// the code coverage of the entire project. The type of PackagesTestDetail
// MUST map[string]PackageTest(contains pass-status,code-coverage and time).
// and it has all packages' detail infomation. PackagesRaceDetail contains
// all packages' race cases.
//
// And the UnitTest contains all packages' result.
type UnitTest struct {
	AvgCover           string                 `json:"average_cover"`
	PackagesTestDetail map[string]PackageTest `json:"packages_test_detail"`
	PackagesRaceDetail map[string][]string    `json:"packages_race_detail"`
}

// PackageTest is a struct that contains IsPass, Coverage and Time. The
// type of Time MUST float64.
type PackageTest struct {
	IsPass   bool    `json:"is_pass"`
	Coverage string  `json:"coverage"`
	Time     float64 `json:"time"`
}

// Cycloi is a struct that contains Average and Result. The Average is
// one package's cyclo coverage. And Result is the detail cyclo of the package's
// all function.
type Cycloi struct {
	Average string
	Result  []string
}

// Test is a struct that contains Path, Result, Time and Cover. The type of
// Time and Cover MUST float64. And it is just for one package's display.
type Test struct {
	Path   string
	Result int
	Time   float64
	Cover  float64
}

// File is a struct that contains Color, CycloVal and CycloInfo. And it is just
// for one file's display. The CycloInfo contains all cyclo detail information.
type File struct {
	Color     string
	CycloVal  string
	CycloInfo string
}

// Copycode is a struct that contains Files and Path. The type of Path MUST []string
// that contains more than one file path. The Copycode represents some copyed code
// information.
type Copycode struct {
	Files string
	Path  []string
}

// Race is a struct that contains Pkg, Len, Leng and Info. The type of Info MUST
// []string that represents more than one race case. Len is the number of cases.
type Race struct {
	Pkg  string
	Len  string
	Leng string
	Info []string
}

// Simple is a struct that contains Path and Info. The type of Path and Info MUST string.
// The Simple represents one can be simpled code case.
type Simple struct {
	Path string
	Info string
}

// Spell is a struct that contains Path and Info. The type of Path and Info MUST string.
// The Spell represents one word is misspelled.
type Spell struct {
	Path string
	Info string
}

// Scan is a struct that contains Path and Info. The type of Path and Info MUST string.
// The Scan represents one defect case.
type Scan struct {
	Path string
	Info string
}

// Deadcode is a struct that contains Path and Info. The type of Path and Info MUST string.
// The Deadcode represents one dead code.
type Deadcode struct {
	Path string
	Info string
}

// CycloInfo is a struct that contains Comp and Info. The type of Comp MUST int that represents
// the cyclo of one function.The CycloInfo represents one cyclo function information.
type CycloInfo struct {
	Comp int
	Info string
}

// Cyclo is a struct that contains Pkg, Size and Info. The type of Info MUST []CycloInfo that represents
// detail information of all function.
type Cyclo struct {
	Pkg  string
	Size int
	Info []CycloInfo
}

// Depth is a struct that contains Pkg, Size and Info. Info is an alias to CycloInfo
type Depth struct {
	Pkg  string
	Size int
	Info []DepthInfo
}

type DepthInfo CycloInfo

// UnitTest is a struct that contains some features in a report of html.
//         GoReporter HTML Report Features
//
//    +----------------------------------------------------------------------+
//    |        Feature        |                 Information                  |
//    +=======================+==============================================+
//    | Project               | The path address of the item being detected  |
//    +-----------------------+----------------------------------------------+
//    | Score                 | The score of the tested project              |
//    +-----------------------+----------------------------------------------+
//    | Tests                 | Unit test results                            |
//    +-----------------------+----------------------------------------------+
//    | Date                  | Date assessment of the project               |
//    +-----------------------+----------------------------------------------+
//    | Issues                | Issues number of the project                 |
//    +-----------------------+----------------------------------------------+
//    | FileCount             | Go file number of the peoject                |
//    +-----------------------+----------------------------------------------+
//    | CodeLines             | Number of lines of code                      |
//    +-----------------------+----------------------------------------------+
//    | CommentLines          | Number of lines of Comment                   |
//    +-----------------------+----------------------------------------------+
//    | TestSummaryCoverAvg   | Code cover average of all unit test          |
//    +-----------------------+----------------------------------------------+
//    | AveragePackageCover   | Package cover average of all packages        |
//    +-----------------------+----------------------------------------------+
//    | SimpleIssues          | Simpled issues number                        |
//    +-----------------------+----------------------------------------------+
//    | DeadcodeIssues        | Dead code issues number                      |
//    +-----------------------+----------------------------------------------+
//    | CycloBigThan15        | Cyclo more than 15 number                    |
//    +-----------------------+----------------------------------------------+
//    | Races                 | Race result of all packages                  |
//    +-----------------------+----------------------------------------------+
//    | NoTests               | No unit test packages information            |
//    +-----------------------+----------------------------------------------+
//    | Simples               | Simpled cases of all packages information    |
//    +-----------------------+----------------------------------------------+
//    | SimpleLevel           | Simple level of path                         |
//    +-----------------------+----------------------------------------------+
//    | Deadcodes             | Dead code cases information                  |
//    +-----------------------+----------------------------------------------+
//    | Copycodes             | Copy code cases information                  |
//    +-----------------------+----------------------------------------------+
//    | Cyclos                | Cyclo of funtion cases information           |
//    +-----------------------+----------------------------------------------+
//    | DepGraph              | Depend graph of all packages in the project  |
//    +-----------------------+----------------------------------------------+
//    | LastRefresh           | Last refresh time of one project             |
//    +-----------------------+----------------------------------------------+
//    | HumanizedLastRefresh  | Humanized last refresh setting               |
//    +-----------------------+----------------------------------------------+
//
// And the HtmlData just as data for default html template. If you want to customize
// your own template file, please follow these data, or you can redefine it yourself.
type HtmlData struct {
	Project             string
	Score               int
	Tests               string
	Date                string
	Issues              int
	FileCount           int
	CodeLines           int
	CommentLines        int
	TotalLines          int
	TestSummaryCoverAvg string
	AveragePackageCover float64
	SimpleIssues        int
	DeadcodeIssues      int
	CycloBigThan15      int
	Races               []Race
	NoTests             string
	Simples             string
	Spells              string
	SimpleLevel         int
	Deadcodes           string
	Copycodes           string
	Cyclos              string
	Depths              string
	DepGraph            template.HTML

	LastRefresh          time.Time `json:"last_refresh"`
	HumanizedLastRefresh string    `json:"humanized_last_refresh"`
}
