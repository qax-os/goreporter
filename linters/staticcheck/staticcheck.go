package staticcheck

import (
	"github.com/360EntSecGroup-Skylar/goreporter/linters/simpler/lint/lintutil"
)

func StaticCheck(projectPath []string) []string {
	fs := lintutil.FlagSet("staticcheck")
	gen := fs.Bool("generated", false, "Check generated code")
	fs.Parse(projectPath)
	c := NewChecker()
	c.CheckGenerated = *gen
	return lintutil.ProcessFlagSet(c, fs)
}
