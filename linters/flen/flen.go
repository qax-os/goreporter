package flen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"strconv"
	"strings"
)

type funcLen struct {
	Name           string
	Size           int // len of function
	Filepath       string
	Lbrace, Rbrace int
	Type           funcType
}

// Options provides the user of this package to configre it
// for according to his/her needs.
type Options struct {
	IncludeTests bool
	BucketSize   int
}

type funcType int

const (
	// Sentinel defines max length of func this pkg can handle
	Sentinel                     = 1000000
	defaultBucketSize            = 5
	defaultIncludeTests          = false
	streakChar                   = "∎"
	lenLowerLimit                = 0
	lenUpperLimit                = Sentinel
	implemented         funcType = iota
	implementedAtRuntime
)

var (
	goroot  = os.Getenv("GOROOT")
	pkgpath string
	gopath  = os.Getenv("GOPATH")
	opts    *Options
)

func rangeAsked(ll, ul int) bool { return ll != 0 || ul != Sentinel }

func FuncLen(packageName string) (lengAll []int, data [][]string) {
	if packageName == "" {
		return
	}

	flenOptions := &Options{
		IncludeTests: false,
		BucketSize:   5,
	}
	flens, _, err := GenerateFuncLens(packageName, flenOptions)
	if err != nil {
		fmt.Println(err)
		return
	}

	lengAll = flens.computeHistogram()

	for i, f := range flens {
		data = append(data, []string{strconv.Itoa(i), f.Name, f.Filepath, strconv.Itoa(f.Lbrace), strconv.Itoa(f.Size)})
	}
	return
}

// FuncLens is the main object which flen pkg exposes to client.
// All the operation are done on this.
type FuncLens []funcLen

// Sort interface is implemented for FuncLens type.
// Sorting is needed for computing percentiles.
func (flens *FuncLens) Len() int { return len(*flens) }
func (flens *FuncLens) Less(i, j int) bool {
	switch {
	case (*flens)[i].Size < (*flens)[j].Size:
		return true
	case (*flens)[i].Size > (*flens)[j].Size:
		return false
	case strings.Compare((*flens)[i].Name, (*flens)[j].Name) == -1:
		return true
	case strings.Compare((*flens)[i].Name, (*flens)[j].Name) == 1:
		return false
	}
	return false
}
func (flens *FuncLens) Swap(i, j int) { (*flens)[i], (*flens)[j] = (*flens)[j], (*flens)[i] }

// createHistogram computes and returns slice of histogram data points.
func (flens *FuncLens) computeHistogram() []int {
	var hg []int
	var x int
	if len(*flens) == 0 {
		return nil
	}
	// find max func len
	var maxFlen int = -Sentinel
	for _, flen := range *flens {
		if flen.Size > maxFlen {
			maxFlen = flen.Size
		}
	}
	hglen := maxFlen / opts.BucketSize
	hg = make([]int, hglen+1)
	for _, v := range *flens {
		if v.Size > 0 {
			x = v.Size % opts.BucketSize
			if x == 0 {
				x = v.Size/opts.BucketSize - 1
			} else {
				x = v.Size / opts.BucketSize
			}
			hg[x]++
		}
	}
	return hg
}

// GenerateFuncLens generates FuncLens for the given package. If options.InclTests is true,
// functions in tests are also evaluated. For ease in readibility of func lens in table,
// result is sorted.
func GenerateFuncLens(pkg string, options *Options) (FuncLens, string, error) {
	opts = options
	if opts == nil {
		opts = &Options{
			IncludeTests: defaultIncludeTests,
			BucketSize:   defaultBucketSize,
		}
	}
	pkgpath, err := getPkgPath(pkg)
	if err != nil {
		return nil, pkgpath, err
	}

	fset := token.NewFileSet()
	pkgs, ferr := parser.ParseDir(fset, pkgpath, func(f os.FileInfo) bool {
		if opts.IncludeTests {
			return true
		}
		return !strings.HasSuffix(f.Name(), "_test.go")
	}, parser.AllErrors)
	if ferr != nil {
		panic(ferr)
	}
	flens := make(FuncLens, 0)
	for _, v := range pkgs {
		for filepath, astf := range v.Files {
			for _, decl := range astf.Decls {
				ast.Inspect(decl, func(node ast.Node) bool {
					var (
						funcname string
						diff     int
						lb, rb   token.Pos
						rln, lln int
						ftype    funcType
					)

					if x, ok := node.(*ast.FuncDecl); ok {
						ftype = implemented
						funcname = x.Name.Name
						if x.Body == nil {
							ftype = implementedAtRuntime // externally implemented
						} else {
							lb = x.Body.Lbrace
							rb = x.Body.Rbrace
							if !lb.IsValid() || !rb.IsValid() {
								return false
							}
							rln = fset.Position(rb).Line
							lln = fset.Position(lb).Line
							diff = rln - lln - 1
							if diff == -1 {
								diff = 1 // single line func
							}
						}
						flens = append(flens, funcLen{
							Name:     funcname,
							Size:     diff,
							Filepath: filepath,
							Lbrace:   lln,
							Rbrace:   rln,
							Type:     ftype,
						})
					}
					return false

				})
			}
		}
	}
	sort.Sort(&flens)
	return flens, pkgpath, nil
}

// getPkgPath tries to get path of pkg. Path is platform dependent.
// First pkg is checked in GOPATH, then in GOROOT, then err.
func getPkgPath(pkgname string) (string, *os.PathError) {
	var ppath string
	if gopath != "" {
		for _, godir := range strings.Split(gopath, string(os.PathListSeparator)) {
			ppath = strings.Join([]string{godir, "src", pkgname}, string(os.PathSeparator))
			_, err := os.Stat(ppath)
			if err != nil {
				continue
			}
			return ppath, nil
		}
	}
	ppath = strings.Join([]string{goroot, "src", pkgname}, string(os.PathSeparator))
	_, err := os.Stat(ppath)
	if err != nil {
		return "", err.(*os.PathError)
	}
	return ppath, nil
}
