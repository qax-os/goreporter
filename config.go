package main

import (
	"html/template"
)

type Reporter struct {
	Project     string              `json:"project"`
	UnitTestx   UnitTest            `json:"unit_test"`
	Cyclox      map[string]Cycloi   `json:"cyclo"`
	SimpleTips  map[string][]string `json:"simple_tips"`
	CopyTips    [][]string          `json:"copy_tips"`
	ScanTips    map[string][]string `json:"scan_tips"`
	DependGraph string              `json:"depend_graph"`
	DeadCode    []string            `json:"dead_code"`
	NoTestPkg   []string            `json:"no_test_package"`
}

type UnitTest struct {
	AvgCover           string                 `json:"average_cover"`
	PackagesTestDetail map[string]PackageTest `json:"packages_test_detail"`
	PackagesRaceDetail map[string][]string    `json:"packages_race_detail"`
}

type PackageTest struct {
	IsPass   bool    `json:"is_pass"`
	Coverage string  `json:"coverage"`
	Time     float64 `json:"time"`
}

type Cycloi struct {
	Average string   `json:"average"`
	Result  []string `json:"result"`
}

type Test struct {
	Path   string  `json:path`
	Result int     `json:result`
	Time   float64 `json:time`
	Cover  float64 `json:cover`
}

type File struct {
	Color     string
	CycloVal  string
	CycloInfo string
}

type Copycode struct {
	Files string
	Path  []string
}

type Race struct {
	Pkg  string
	Len  string
	Leng string
	Info []string
}

type Simple struct {
	Path string `json:path`
	Info string `json:info`
}

type Scan struct {
	Path string `json:path`
	Info string `json:info`
}

type Deadcode struct {
	Path string `json:path`
	Info string `json:info`
}
type CycloInfo struct {
	Comp int
	Info string
}
type Cyclo struct {
	Pkg  string
	Size int
	Info []CycloInfo
}

type HtmlData struct {
	Project               string
	Score                 int
	Tests                 string
	TestSummaryCoverAvg   string
	Races                 []Race
	NoTests               string
	Simples               string
	SimpleLevel           int
	Deadcodes             string
	Copycodes             string
	Cyclos                string
	CycloSummarysCycloAvg string
	CycloSummarysCyclo50  string
	CycloSummarysCyclo15  string
	DepGraph              template.HTML
}

type ApolloMeta struct {
	Branch         string `json:"branch"`
	Project        string `json:"project"`
	CommitID       string `json:"commitid"`
	CommitUser     string `json:"commituser"`
	User           string `json:"user"`
	UnintTestCover string `json:"uninttestcover"`
	StaticCheck    string `json:"staticcheck"`
	CycloBig       string `json:"cyclobig"`
	Score          int    `json:"score"`
	Starttime      string `json:"starttime"`
	Endtime        string `json:"endtime"`
}

// type OptionMeta struct {
// 	Branch     string `json:"branch"`
// 	CommitID   string `json:"commitID"`
// 	CommitUser string `json:"commitUser"`
// 	User       string `json:"user"`
// }

// type Value struct {
// 	Filepath string     `json:"filepath"`
// 	Meta     ApolloMeta `json:"meta"`
// }
