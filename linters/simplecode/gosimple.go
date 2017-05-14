package simplecode // import "github.com/wgliang/goreporter/linters/simplecode"

import (
	"github.com/wgliang/goreporter/linters/simplecode/lint/lintutil"
	"github.com/wgliang/goreporter/linters/simplecode/simple"
)

func Simple(path map[string]string) []string {
	var res []string
	for _, p := range path {
		res = append(res, lintutil.ProcessArgs("gosimple", simple.Funcs, []string{p})...)
	}
	return res
}
