package depth

// godepth calculates maximum depth of go methods in Go source code.
//
// This work was mainly inspired by github.com/fzipp/gocyclo
//
// Usage:
//      godepth [<flag> ...] <Go file or directory> ...
//
// Flags:
//      -over N   show functions with depth > N only and
//                return exit code 1 if the output is non-empty
//      -top N    show the top N most complex functions only
//      -avg      show the average depth
//
// The output fields for each line are:
// <depth> <package> <function> <file:row:column>

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
)

const usageDoc = `Calculate maximum depth of Go functions.
Usage:
        godepth [flags...] <Go file or directory> ...

Flags:
        -over N        show functions with depth > N only and
                       return exit code 1 if the set is non-empty
        -top N         show the top N most complex functions only
        -avg           show the average depth over all functions,
                       not depending on whether -over or -top are set

The output fields for each line are:
<depth> <package> <function> <file:row:column>
`

func usage() {
	fmt.Fprintf(os.Stderr, usageDoc)
	os.Exit(2)
}

var (
	over = 0
	top  = 10
	avg  = false
)

func Depth(packagePath string) ([]string, string) {
	args := []string{packagePath}
	if len(args) == 0 {
		usage()
	}

	stats := analyze(args)
	sort.Sort(byDepth(stats))

	packageAvg := "0"
	if avg {
		packageAvg = getAverage(stats)
	}

	result := make([]string, 0)
	for _, stat := range stats {
		result = append(result, stat.String())
	}

	return result, packageAvg
}

func analyze(paths []string) []stat {
	stats := []stat{}
	for _, path := range paths {
		if isDir(path) {
			stats = analyzeDir(path, stats)
		} else {
			stats = analyzeFile(path, stats)
		}
	}
	return stats
}

func isDir(filename string) bool {
	fi, err := os.Stat(filename)
	return err == nil && fi.IsDir()
}

func analyzeDir(dirname string, stats []stat) []stat {
	files, _ := filepath.Glob(filepath.Join(dirname, "*.go"))
	for _, file := range files {
		stats = analyzeFile(file, stats)
	}
	return stats
}

func analyzeFile(fname string, stats []stat) []stat {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fname, nil, 0)
	if err != nil {
		exitError(err)
	}
	return buildStats(f, fset, stats)
}

func exitError(err error) {
	fmt.Fprintln(os.Stderr, err)
	// os.Exit(1)
}

func getAverage(stats []stat) string {
	return fmt.Sprintf("%.2f", average(stats))
}

func average(stats []stat) float64 {
	total := 0
	for _, s := range stats {
		total += s.Depth
	}
	return float64(total) / float64(len(stats))
}

type stat struct {
	PkgName  string
	FuncName string
	Depth    int
	Pos      token.Position
}

func (s stat) String() string {
	return fmt.Sprintf("%d %s %s %s", s.Depth, s.PkgName, s.FuncName, s.Pos)
}

type byDepth []stat

func (s byDepth) Len() int      { return len(s) }
func (s byDepth) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byDepth) Less(i, j int) bool {
	return s[i].Depth >= s[j].Depth
}

func buildStats(f *ast.File, fset *token.FileSet, stats []stat) []stat {
	for _, decl := range f.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			stats = append(stats, stat{
				PkgName:  f.Name.Name,
				FuncName: funcName(fn),
				Depth:    getdepth(fn),
				Pos:      fset.Position(fn.Pos()),
			})
		}
	}
	return stats
}

// funcName returns the name representation of a function or method:
// "(Type).Name" for methods or simply "Name" for functions.
func funcName(fn *ast.FuncDecl) string {
	if fn.Recv != nil {
		if fn.Recv.NumFields() > 0 {
			typ := fn.Recv.List[0].Type
			return fmt.Sprintf("(%s).%s", recvString(typ), fn.Name)
		}
	}
	return fn.Name.Name
}

// recvString returns a string representation of recv of the
// form "T", "*T", or "BADRECV" (if not a proper receiver type).
func recvString(recv ast.Expr) string {
	switch t := recv.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + recvString(t.X)
	}
	return "BADRECV"
}

func max(s []int) (m int) {
	for _, value := range s {
		if value > m {
			m = value
		}
	}
	return
}

// getdepth calculates the depth of a function
func getdepth(fn *ast.FuncDecl) int {
	allDepth := []int{}
	if fn.Body == nil {
		return 0
	}
	for _, lvl := range fn.Body.List {
		v := maxDepthVisitor{}
		ast.Walk(&v, lvl)
		allDepth = append(allDepth, max(v.NodeDepth))
	}
	return max(allDepth)
}

type maxDepthVisitor struct {
	Depth     int
	NodeDepth []int
	Lbrace    token.Pos
	Rbrace    token.Pos
}

// Visit implements the ast.Visitor interface.
// Basically, it counts the number of consecutive brackets
// Each time Visit is called, we store the current depth in a slice
// When it encounters a sibling eg: if {} followed by if{},
// it saves the current depth and decreases the depth counter
func (v *maxDepthVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.BlockStmt:
		if v.Rbrace == 0 && v.Lbrace == 0 {
			v.Lbrace = n.Lbrace
			v.Rbrace = n.Rbrace
		}

		if n.Lbrace > v.Lbrace && n.Rbrace > v.Rbrace {
			v.Depth--
		}

		v.Lbrace = n.Lbrace
		v.Rbrace = n.Rbrace
		v.Depth++
		v.NodeDepth = append(v.NodeDepth, v.Depth)
	}

	return v
}
