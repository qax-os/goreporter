package simpler

import (
	"github.com/360EntSecGroup-Skylar/goreporter/linters/simpler/lint/lintutil"
)

func Simpler(projectPath []string) []string {
	fs := lintutil.FlagSet("gosimple")
	gen := fs.Bool("generated", false, "Check generated code")
	fs.Parse(projectPath)
	c := NewChecker()
	c.CheckGenerated = *gen
	return lintutil.ProcessFlagSet(c, fs)
}
