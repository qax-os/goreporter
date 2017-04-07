package staticscan

import (
	"github.com/wgliang/goreporter/linters/staticscan/lint/lintutil"
	"github.com/wgliang/goreporter/linters/staticscan/staticcheck"
)

func StaticScan(projectPath string) []string {
	fs := lintutil.FlagSet("staticcheck")
	gen := fs.Bool("generated", false, "Check generated code")
	fs.Parse([]string{projectPath})
	c := staticcheck.NewChecker()
	c.CheckGenerated = *gen
	return lintutil.ProcessFlagSet(c, fs)
}
