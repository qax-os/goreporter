package main

type Reporter struct {
	UnitTestx  UnitTest            `json:unit_test`
	Cyclox     map[string]Cyclo    `json:cyclo`
	SimpleTips map[string][]string `json:simple_tips`
	CopyTips   [][]string          `json:copy_tips`
	ScanTips   map[string][]string `json:scan_tips`
}

type UnitTest struct {
	AvgCover           string                 `json:average_cover`
	PackagesTestDetail map[string]PackageTest `json:packages_test_detail`
	PackagesRaceDetail map[string][]string    `json:packages_race_detail`
}

type PackageTest struct {
	IsPass   bool    `json:is_pass`
	Coverage string  `json:coverage`
	Time     float64 `json:time`
}

type Cyclo struct {
	Average string   `json:average`
	Result  []string `json:cyclo_res`
}
