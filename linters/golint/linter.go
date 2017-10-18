// Copyright (c) 2013 The Go Authors. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.

// golint lints the Go source files named on its command line.
package golint

import (
	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	minConfidence = flag.Float64("min_confidence", 0.8, "minimum confidence of a problem to print it")
	setExitStatus = flag.Bool("set_exit_status", false, "set exit status to 1 if any issues are found")
	suggestions   int
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\tgolint [flags] # runs on package in current directory\n")
	fmt.Fprintf(os.Stderr, "\tgolint [flags] [packages]\n")
	fmt.Fprintf(os.Stderr, "\tgolint [flags] [directories] # where a '/...' suffix includes all sub-directories\n")
	fmt.Fprintf(os.Stderr, "\tgolint [flags] [files] # all must belong to a single package\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func GoLinter(projectPath []string) (results []string) {
	flag.Usage = usage
	flag.Parse()

	// dirsRun, filesRun, and pkgsRun indicate whether golint is applied to
	// directory, file or package targets. The distinction affects which
	// checks are run. It is no valid to mix target types.
	var dirsRun, filesRun, pkgsRun int
	var args []string
	for _, arg := range projectPath {
		if strings.HasSuffix(arg, "/...") && isDir(arg[:len(arg)-len("/...")]) {
			dirsRun = 1
			args = append(args, allPackagesInFS(arg)...)
		} else if isDir(arg) {
			dirsRun = 1
			args = append(args, arg)
		} else if exists(arg) {
			filesRun = 1
			args = append(args, arg)
		} else {
			pkgsRun = 1
			args = append(args, arg)
		}
	}

	if dirsRun+filesRun+pkgsRun != 1 {
		usage()
		os.Exit(2)
	}
	switch {
	case dirsRun == 1:
		for _, dir := range args {
			results = append(results, lintDir(dir)...)
		}
	case filesRun == 1:
		return lintFiles(args...)
	case pkgsRun == 1:
		for _, pkg := range importPaths(args) {
			results = append(results, lintPackage(pkg)...)
		}
	}

	return
}

func isDir(filename string) bool {
	fi, err := os.Stat(filename)
	return err == nil && fi.IsDir()
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func lintFiles(filenames ...string) (results []string) {
	files := make(map[string][]byte)
	for _, filename := range filenames {
		src, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		files[filename] = src
	}

	l := new(Linter)
	ps, err := l.LintFiles(files)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	for _, p := range ps {
		if p.Confidence >= *minConfidence {
			results = append(results, fmt.Sprintf("%v: %s", p.Position, p.Text))
			// suggestions++
		}
	}
	return results
}

func lintDir(dirname string) []string {
	pkg, err := build.ImportDir(dirname, 0)
	return lintImportedPackage(pkg, err)
}

func lintPackage(pkgname string) []string {
	pkg, err := build.Import(pkgname, ".", 0)
	return lintImportedPackage(pkg, err)
}

func lintImportedPackage(pkg *build.Package, err error) (results []string) {
	if err != nil {
		if _, nogo := err.(*build.NoGoError); nogo {
			// Don't complain if the failure is due to no Go source files.
			return
		}
		fmt.Fprintln(os.Stderr, err)
		return
	}

	var files []string
	files = append(files, pkg.GoFiles...)
	files = append(files, pkg.CgoFiles...)
	files = append(files, pkg.TestGoFiles...)
	if pkg.Dir != "." {
		for i, f := range files {
			files[i] = filepath.Join(pkg.Dir, f)
		}
	}
	// TODO(dsymonds): Do foo_test too (pkg.XTestGoFiles)

	return lintFiles(files...)
}
