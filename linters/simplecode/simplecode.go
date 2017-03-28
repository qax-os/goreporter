package simplecode

import (
	// "fmt"
	// "os"

	"github.com/wgliang/goreporter/linters/simplecode/lint/lintutil"
	"github.com/wgliang/goreporter/linters/simplecode/simple"
)

func SimpleCode(prohectPath string) []string {
	fs := lintutil.FlagSet("gosimple")
	gen := fs.Bool("generated", false, "Check generated code")
	fs.Parse([]string{prohectPath})
	c := simple.NewChecker()
	c.CheckGenerated = *gen

	return lintutil.ProcessFlagSet(c, fs)
}
