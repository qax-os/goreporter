package simpler

import (
	"github.com/360EntSecGroup-Skylar/goreporter/linters/simpler/lint/lintutil"
)

func Simpler(projectPath map[string]string) []string {
	fs := lintutil.FlagSet("gosimple")
	gen := fs.Bool("generated", false, "Check generated code")
	paths := make([]string, len(projectPath))
	for _, v := range projectPath {
		paths = append(paths, v)
	}
	fs.Parse(paths)
	c := NewChecker()
	c.CheckGenerated = *gen
	return lintutil.ProcessFlagSet(c, fs)
}
