package tools

import (
	"html/template"
	"time"
)

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
	Average string
	Result  []string
}

type Test struct {
	Path   string
	Result int
	Time   float64
	Cover  float64
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
	Path string
	Info string
}

type Scan struct {
	Path string
	Info string
}

type Deadcode struct {
	Path string
	Info string
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
	Project             string
	Score               int
	Tests               string
	Date                string
	TestSummaryCoverAvg string
	AveragePackageCover float64
	SimpleIssues        int
	DeadcodeIssues      int
	CycloBigThan15      int
	Races               []Race
	NoTests             string
	Simples             string
	SimpleLevel         int
	Deadcodes           string
	Copycodes           string
	Cyclos              string
	DepGraph            template.HTML

	LastRefresh          time.Time `json:"last_refresh"`
	HumanizedLastRefresh string    `json:"humanized_last_refresh"`
}
