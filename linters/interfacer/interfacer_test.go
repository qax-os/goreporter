// Copyright (c) 2015, Daniel Martí <mvdan@mvdan.cc>
// See LICENSE for licensing information

package interfacer

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/kisielk/gotool"
)

const testdata = "testdata"

var (
	issuesRe = regexp.MustCompile(`^WARN (.*)\n?$`)
	singleRe = regexp.MustCompile(`([^ ]*) can be ([^ ]*)(,|$)`)
)

func goFiles(t *testing.T, p string) []string {
	if strings.HasSuffix(p, ".go") {
		return []string{p}
	}
	dirs := gotool.ImportPaths([]string{p})
	var paths []string
	for _, dir := range dirs {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			t.Fatal(err)
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			paths = append(paths, filepath.Join(dir, file.Name()))
		}
	}
	return paths
}

type identVisitor struct {
	fset   *token.FileSet
	idents map[string]token.Pos
}

func identKey(line int, name string) string {
	return fmt.Sprintf("%d %s", line, name)
}

func (v *identVisitor) Visit(n ast.Node) ast.Visitor {
	switch x := n.(type) {
	case *ast.Ident:
		line := v.fset.Position(x.Pos()).Line
		v.idents[identKey(line, x.Name)] = x.Pos()
	}
	return v
}

func identPositions(fset *token.FileSet, f *ast.File) map[string]token.Pos {
	v := &identVisitor{
		fset:   fset,
		idents: make(map[string]token.Pos),
	}
	ast.Walk(v, f)
	return v.idents
}

func wantedIssues(t *testing.T, p string) []string {
	fset := token.NewFileSet()
	lines := make([]string, 0)
	for _, path := range goFiles(t, p) {
		src, err := os.Open(path)
		if err != nil {
			t.Fatal(err)
		}
		f, err := parser.ParseFile(fset, path, src, parser.ParseComments)
		src.Close()
		if err != nil {
			t.Fatal(err)
		}
		identPos := identPositions(fset, f)
		for _, group := range f.Comments {
			cm := issuesRe.FindStringSubmatch(group.Text())
			if cm == nil {
				continue
			}
			for _, m := range singleRe.FindAllStringSubmatch(cm[1], -1) {
				vname, tname := m[1], m[2]
				line := fset.Position(group.Pos()).Line
				pos := fset.Position(identPos[identKey(line, vname)])
				lines = append(lines, fmt.Sprintf("%s: %s can be %s",
					pos, vname, tname))
			}
		}
	}
	return lines
}

func doTest(t *testing.T, p string) {
	t.Run(p, func(t *testing.T) {
		lines := wantedIssues(t, p)
		doTestLines(t, p, lines, p)
	})
}

func doTestLines(t *testing.T, name string, want []string, args ...string) {
	got, err := CheckArgs(args)
	if err != nil {
		t.Fatalf("Did not want error in %s:\n%v", name, err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Output mismatch in %s:\nwant:\n%s\ngot:\n%s",
			name, strings.Join(want, "\n"), strings.Join(got, "\n"))
	}
}

func doTestString(t *testing.T, name, want string, args ...string) {
	switch len(args) {
	case 0:
		args = []string{name}
	case 1:
		if args[0] == "" {
			args = nil
		}
	}
	issues, err := CheckArgs(args)
	if err != nil {
		t.Fatalf("Did not want error in %s:\n%v", name, err)
	}
	got := strings.Join(issues, "\n")
	if want != got {
		t.Fatalf("Output mismatch in %s:\nExpected:\n%s\nGot:\n%s",
			name, want, got)
	}
}

func inputPaths(t *testing.T, glob string) []string {
	all, err := filepath.Glob(glob)
	if err != nil {
		t.Fatal(err)
	}
	return all
}

func chdirUndo(t *testing.T, d string) func() {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(d); err != nil {
		t.Fatal(err)
	}
	return func() {
		if err := os.Chdir(wd); err != nil {
			t.Fatal(err)
		}
	}
}

func runFileTests(t *testing.T, paths ...string) {
	defer chdirUndo(t, "files")()
	if len(paths) == 0 {
		paths = inputPaths(t, "*")
	}
	for _, p := range paths {
		doTest(t, p)
	}
}

func runLocalTests(t *testing.T, paths ...string) {
	defer chdirUndo(t, "local")()
	if len(paths) > 0 {
		for _, p := range paths {
			doTest(t, p)
		}
		return
	}
	for _, p := range inputPaths(t, "*") {
		paths = append(paths, "./"+p+"/...")
	}
	for _, p := range paths {
		doTest(t, p)
	}
	// non-recursive
	doTest(t, "./single")
	doTestString(t, "no-args", "", "")
}

func runNonlocalTests(t *testing.T, paths ...string) {
	// std
	doTestString(t, "std-pkg", "", "sync/atomic")
	defer chdirUndo(t, "src")()
	if len(paths) > 0 {
		for _, p := range paths {
			doTest(t, p)
		}
		return
	}
	paths = inputPaths(t, "*")
	for _, p := range paths {
		doTest(t, p+"/...")
	}
	// local recursive
	doTest(t, "./nested/...")
	// non-recursive
	doTest(t, "single")
	// make sure we don't miss a package's imports
	doTestString(t, "grab-import", "grab-import/use.go:27:15: s can be grab-import/def/nested.Fooer")
	defer chdirUndo(t, "nested/pkg")()
	// relative paths
	doTestString(t, "rel-path", "simple.go:12:17: rc can be Closer", "./...")
}

func TestMain(m *testing.M) {
	if err := os.Chdir(testdata); err != nil {
		panic(err)
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	build.Default.GOPATH = wd
	gotool.DefaultContext.BuildContext.GOPATH = wd
	os.Exit(m.Run())
}

func TestIssues(t *testing.T) {
	runFileTests(t)
	runLocalTests(t)
	runNonlocalTests(t)
}

func TestExtraArg(t *testing.T) {
	_, err := CheckArgs([]string{"single", "--", "foo", "bar"})
	got := err.Error()
	want := "unwanted extra args: [foo bar]"
	if got != want {
		t.Fatalf("Error mismatch:\nExpected:\n%s\nGot:\n%s", want, got)
	}
}
