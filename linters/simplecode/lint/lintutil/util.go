// Copyright (c) 2013 The Go Authors. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.

// Package lintutil provides helpers for writing linter command lines.
package lintutil // import "github.com/360EntSecGroup-Skylar/goreporter/linters/simplecode/lint/lintutil"

import (
	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/360EntSecGroup-Skylar/goreporter/linters/simplecode/lint"

	"github.com/360EntSecGroup-Skylar/goreporter/linters/simplecode/gotool"
)

var (
	excepts = []string{"vendor"}
)

func usage(name string, flags *flag.FlagSet) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", name)
		fmt.Fprintf(os.Stderr, "\t%s [flags] # runs on package in current directory\n", name)
		fmt.Fprintf(os.Stderr, "\t%s [flags] packages\n", name)
		fmt.Fprintf(os.Stderr, "\t%s [flags] directory\n", name)
		fmt.Fprintf(os.Stderr, "\t%s [flags] files... # must be a single package\n", name)
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flags.PrintDefaults()
	}
}

type runner struct {
	funcs         []lint.Func
	minConfidence float64
	tags          []string
	SimpleResult  []string

	unclean bool
}

func (runner runner) resolveRelative(importPaths []string) (goFiles bool, err error) {
	if len(importPaths) == 0 {
		return false, nil
	}
	if strings.HasSuffix(importPaths[0], ".go") {
		// User is specifying a package in terms of .go files, don't resolve
		return true, nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return false, err
	}
	ctx := build.Default
	ctx.BuildTags = runner.tags
	for i, path := range importPaths {
		githubIndex := strings.LastIndex(path, "github.com")
		if githubIndex < len(path) && githubIndex >= 0 {
			path = path[githubIndex:]
		}
		bpkg, err := ctx.Import(path, wd, build.FindOnly)
		if err != nil {
			return false, fmt.Errorf("can't load package %q: %v", path, err)
		}
		importPaths[i] = bpkg.ImportPath
	}
	return false, nil
}

func ProcessArgs(except, name string, funcs []lint.Func, args []string) []string {
	excepts = append(excepts, strings.Split(except, ",")...)
	flags := &flag.FlagSet{}
	flags.Usage = usage(name, flags)
	var minConfidence = flags.Float64("min_confidence", 0.8, "minimum confidence of a problem to print it")
	var tags = flags.String("tags", "", "List of `build tags`")
	flags.Parse(args)

	runner := &runner{
		funcs:         funcs,
		minConfidence: *minConfidence,
		tags:          strings.Fields(*tags),
	}
	paths := gotool.ImportPaths(flags.Args())
	goFiles, err := runner.resolveRelative(paths)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		runner.unclean = true
	}
	if goFiles {
		runner.lintFiles(paths...)
	} else {
		for _, path := range paths {
			runner.lintPackage(path)
		}
	}
	return runner.SimpleResult
}

func (runner *runner) lintFiles(filenames ...string) {
	files := make(map[string][]byte)
	for _, filename := range filenames {
		src, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			runner.unclean = true
			continue
		}
		files[filename] = src
	}

	l := &lint.Linter{
		Funcs: runner.funcs,
	}
	ps, err := l.LintFiles(files)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		runner.unclean = true
		return
	}
	if len(ps) > 0 {
		runner.unclean = true
	}
	for _, p := range ps {
		if p.Confidence >= runner.minConfidence {
			runner.SimpleResult = append(runner.SimpleResult, fmt.Sprintf("%v: %s", p.Position, p.Text))
		}
	}
}

func (runner *runner) lintPackage(pkgname string) {
	ctx := build.Default
	ctx.BuildTags = runner.tags
	pkg, err := ctx.Import(pkgname, ".", 0)
	runner.lintImportedPackage(pkg, err)
}

func (runner *runner) lintImportedPackage(pkg *build.Package, err error) {
	if err != nil {
		if _, nogo := err.(*build.NoGoError); nogo {
			// Don't complain if the failure is due to no Go source files.
			return
		}
		fmt.Fprintln(os.Stderr, err)
		runner.unclean = true
		return
	}

	var files []string
	xtest := pkg.XTestGoFiles
	files = append(files, filterFiles(pkg.GoFiles, pkg.Dir)...)
	files = append(files, filterFiles(pkg.CgoFiles, pkg.Dir)...)
	files = append(files, filterFiles(pkg.TestGoFiles, pkg.Dir)...)
	if pkg.Dir != "." {
		for i, f := range xtest {
			xtest[i] = filepath.Join(pkg.Dir, f)
		}
	}
	runner.lintFiles(xtest...)
	runner.lintFiles(files...)
}

func filterFiles(files []string, pkgDir string) (filtedFiles []string) {
	if pkgDir != "." {
		for i, f := range files {
			files[i] = filepath.Join(pkgDir, f)
			if !exceptPkg(files[i]) {
				filtedFiles = append(filtedFiles, files[i])
			}
		}
	} else {
		return files
	}
	return filtedFiles
}

// exceptPkg is a function that will determine whether the package is an exception.
func exceptPkg(pkg string) bool {
	if len(excepts) == 0 {
		return false
	}
	for _, va := range excepts {
		if strings.Contains(pkg, va) {
			return true
		}
	}
	return false
}
