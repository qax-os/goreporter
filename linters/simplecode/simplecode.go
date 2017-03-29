package simplecode

import (
	// "fmt"
	// "os"

	"github.com/wgliang/goreporter/linters/simplecode/simple"
	"github.com/wgliang/goreporter/linters/staticscan/lint/lintutil"
)

func SimpleCode(prohectPath string) []string {
	fs := lintutil.FlagSet("gosimple")
	gen := fs.Bool("generated", false, "Check generated code")
	fs.Parse([]string{prohectPath})
	c := simple.NewChecker()
	c.CheckGenerated = *gen

	return lintutil.ProcessFlagSet(c, fs)
}
