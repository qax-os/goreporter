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

// UnitTest is a struct that contains AvgCover, PackagesTestDetail and
// PackagesRaceDetail. The type of AvgCover MUST string that represents
// the code coverage of the entire project. The type of PackagesTestDetail
// MUST map[string]PackageTest(contains pass-status,code-coverage and time).
// and it has all packages' detail information. PackagesRaceDetail contains
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

// Interfacer is a struct that contains Path and Info. The type of Path and Info MUST string.
// The Interfacer warns about the usage of types that are more specific than necessary.
type Interfacer struct {
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

// Cyclo is a struct that contains Pkg, Size and Info. The type of Info MUST []CycloInfo that represents
// detail information of all function.
type Cyclo struct {
	Pkg  string
	Size int
	Info []CycloInfo
}

// CycloInfo is a struct that contains Comp and Info. The type of Comp MUST int that represents
// the cyclo of one function.The CycloInfo represents one cyclo function information.
type CycloInfo struct {
	Comp int
	Info string
}

// Depth is a struct that contains Pkg, Size and Info. Info is an alias to CycloInfo
type Depth struct {
	Pkg  string
	Size int
	Info []DepthInfo
}

type DepthInfo CycloInfo

// CodeTest is a struct that contains Summary and Content. It represents the result data
// of the project unit test.
type CodeTest struct {
	Summary struct {
		CodeCover    float64 `json:"code_cover"`
		PackageCover float64 `json:"pkg_cover"`
		TotalTime    float64 `json:"total_time"`
	} `json:"summary"`
	Content struct {
		Pkg    []string  `json:"pkg"`
		Cover  []float64 `json:"cover"`
		Time   []float64 `json:"time"`
		NoTest []string  `json:"no_test"`
	} `json:"content"`
}

// Item is a struct that contains File and Content. It is the code-based infrastructure.
type Item struct {
	File    string   `json:"rep"`
	Content []string `json:"content"`
}

// CodeTest is a struct that contains Summary and Content. It represents the result data
// of the project unit test.
type StyleItem struct {
	Label  string `json:"label"`
	Score  int    `json:"score"`
	Detail []Item `json:"detail"`

	filesNum  int
	issuesNum int
}

// CodeTest is a struct that contains Summary and Content. It represents the result data
// of the project unit test.
type CodeStyle struct {
	Summary struct {
		IssuesNum int    `json:"issue_num"`
		FilesNum  int    `json:"file_num"`
		Quality   string `json:"quality"`
	} `json:"summary"`
	Content struct {
		GoFmt     StyleItem `json:"go_fmt"`
		GoVet     StyleItem `json:"go_vet"`
		GoLint    StyleItem `json:"go_lint"`
		MissSpell StyleItem `json:"miss_spell"`
	} `json:"content"`
}

// CodeTest is a struct that contains Summary and Content. It represents the result data
// of the project unit test.
type CopyItem struct {
	Label  string     `json:"label"`
	Score  int        `json:"score"`
	Detail [][]string `json:"detail"`

	filesNum  int
	issuesNum int
}

// CodeTest is a struct that contains Summary and Content. It represents the result data
// of the project unit test.
type CodeOptimization struct {
	Summary struct {
		IssuesNum int    `json:"issue_num"`
		FilesNum  int    `json:"file_num"`
		Quality   string `json:"quality"`
	} `json:"summary"`
	Content struct {
		DeadCode       StyleItem `json:"dead_code"`
		SimpleCode     StyleItem `json:"simple_code"`
		StaticCode     StyleItem `json:"static_code"`
		CopyCode       CopyItem  `json:"copy_code"`
		InterfacerCode StyleItem `json:"interfacer_code"`
	} `json:"content"`
}

// CodeCount is a struct that contains Summary and Content. t represents the code
// statistics of the data, you can understand the project from the perspective
// of the number of lines.
type CodeCount struct {
	Summary struct {
		LineCount     int `json:"line_count"`
		CommentCount  int `json:"comment_count"`
		FunctionCount int `json:"function_count"`
		FileCount     int `json:"file_count"`
	} `json:"summary"`
	Content struct {
		Pkg               []string `json:"pkg"`
		PkgLineCount      []int    `json:"pkg_line_count"`
		PkgCommentCount   []int    `json:"pkg_comment_count"`
		PkgFunctionCount  []int    `json:"pkg_function_count"`
		File              []string `json:"file"`
		FileLineCount     []int    `json:"file_line_count"`
		FileCommentCount  []int    `json:"file_comment_count"`
		FileFunctionCount []int    `json:"file_function_count"`
	} `json:"content"`
}

// CodeSmellItem is a struct that contains path and cyclo.
type CodeSmellItem struct {
	Path  string `json:"path"`
	Cyclo int    `json:"cyclo"`
}

// CodeSmell is a struct that contains Summary and Content. It represents the taste of
// the code data, you can understand the project from the perspective of complexity.
type CodeSmell struct {
	Summary struct {
		CycloAvg   int `json:"cyclo_avg"`
		CycloHigh  int `json:"cyclo_high"`
		CycloGrave int `json:"cyclo_grave"`
	} `json:"summary"`
	Content struct {
		Percentage map[string]int  `json:"percentage"`
		Pkg        []string        `json:"pkg"`
		Cyclo      []int           `json:"cyclo"`
		List       []CodeSmellItem `json:"list"`
	} `json:"content"`
}
