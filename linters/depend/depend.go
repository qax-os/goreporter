package depend

import (
	"bytes"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/golang/glog"
)

var (
	pkgs   map[string]*build.Package
	ids    map[string]int
	nextId int

	ignored = map[string]bool{
		"C": true,
	}
	ignoredPrefixes []string
	// ignore packages in the Go standard library
	ignoreStdlib   = true
	delveGoroot    = false
	ignorePrefixes = ""
	// a comma-separated list of packages to ignore
	ignorePackages = ""
	// a comma-separated list of build tags to consider satisified during the build
	tagList    = ""
	horizontal = true
	// include test packages
	includeTests = false

	buildTags    []string
	buildContext = build.Default

	vendors []string
)

func Depend(path, expect string) string {
	ignorePackages = expect
	vendors = getVendorlist(path)
	// add root vendor
	vendors = append(vendors, "vendor")
	pkgs = make(map[string]*build.Package)
	ids = make(map[string]int)

	args := []string{path}

	if len(args) != 1 {
		glog.Errorln("need one package name to process")
		return ""
	}

	if ignorePrefixes != "" {
		if runtime.GOOS == `windows` {
			ignorePrefixes = strings.Replace(ignorePrefixes, "/", `\`, -1)
		}
		ignoredPrefixes = strings.Split(ignorePrefixes, ",")
	}
	if ignorePackages != "" {
		for _, p := range strings.Split(ignorePackages, ",") {
			ignored[p] = true
		}
	}
	if tagList != "" {
		buildTags = strings.Split(tagList, ",")
	}
	buildContext.BuildTags = buildTags

	cwd, err := os.Getwd()
	if err != nil {
		glog.Errorf("failed to get cwd: %s", err)
		return ""
	}
	if err := processPackage(cwd, strings.Replace(args[0], `\`, "/", -1), path); err != nil {
		glog.Errorln(err)
		return ""
	}

	graph := "digraph godep {"
	// fmt.Println("digraph godep {")
	if horizontal {
		// fmt.Println(`rankdir="LR"`)
		graph += `rankdir="LR"`
	}
	for pkgName, pkg := range pkgs {
		pkgId := getId(pkgName)

		if isIgnored(pkg) {
			continue
		}

		var color string
		if pkg.Goroot {
			color = "palegreen"
		} else if len(pkg.CgoFiles) > 0 {
			color = "darkgoldenrod1"
		} else {
			color = "paleturquoise"
		}
		graph += fmt.Sprintf("%d [label=\"%s\" style=\"filled\" color=\"%s\"];\n", pkgId, pkgName, color)
		// fmt.Printf("%d [label=\"%s\" style=\"filled\" color=\"%s\"];\n", pkgId, pkgName, color)

		// Don't render imports from packages in Goroot
		if pkg.Goroot && !delveGoroot {
			continue
		}

		for _, imp := range getImports(pkg) {
			impPkg := pkgs[imp]
			if impPkg == nil || isIgnored(impPkg) {
				continue
			}

			impId := getId(imp)
			graph += fmt.Sprintf("%d -> %d;\n", pkgId, impId)
			// fmt.Printf("%d -> %d;\n", pkgId, impId)
		}
	}
	graph += `}`

	err = ioutil.WriteFile("graph.gv", []byte(graph), 0666)
	if err != nil {
		glog.Errorln(err)
	}

	// convert file formate
	cmdsvg := exec.Command("dot", "-Tsvg", "-o", "pkgdep.svg", "graph.gv")
	var outsvg bytes.Buffer
	cmdsvg.Stdout = &outsvg
	cmdsvg.Stderr = os.Stderr
	err = cmdsvg.Run()
	if err != nil {
		glog.Errorln(err)
	}

	svg, err := ioutil.ReadFile("pkgdep.svg")
	if err != nil {
		glog.Errorln(err)
	}

	err = os.Remove("pkgdep.svg")
	if err != nil {
		glog.Errorln(err)
	}

	err = os.Remove("graph.gv")
	if err != nil {
		glog.Errorln(err)
	}

	return string(svg)
}

func processPackage(root string, pkgName, path string) error {
	if ignored[pkgName] {
		return nil
	}
	var err error
	if !build.IsLocalImport(pkgName) {
		root, err = filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("failed to convert path to absolute path %s ", err)
		}
	}

	pkg, err := buildContext.Import(pkgName, root, 0)
	if err != nil {
		flag := false
		for i := 0; i < len(vendors); i++ {
			pkg, err = buildContext.Import(vendors[i]+string(filepath.Separator)+pkgName, root, 0)
			if err == nil {
				flag = true
				break
			}
		}
		if !flag {
			return fmt.Errorf("failed to import %s: %s", pkgName, err)
		}
	}

	if isIgnored(pkg) {
		return nil
	}

	pkgs[pkg.ImportPath] = pkg

	// Don't worry about dependencies for stdlib packages
	if pkg.Goroot && !delveGoroot {
		return nil
	}

	for _, imp := range getImports(pkg) {
		if _, ok := pkgs[imp]; !ok {
			if err := processPackage(root, imp, path); err != nil {
				return err
			}
		}
	}
	return nil
}

func getImports(pkg *build.Package) []string {
	allImports := pkg.Imports
	if includeTests {
		allImports = append(allImports, pkg.TestImports...)
		allImports = append(allImports, pkg.XTestImports...)
	}
	var imports []string
	found := make(map[string]struct{})
	for _, imp := range allImports {
		if imp == pkg.ImportPath {
			// Don't draw a self-reference when foo_test depends on foo.
			continue
		}
		if _, ok := found[imp]; ok {
			continue
		}
		found[imp] = struct{}{}
		imports = append(imports, imp)
	}
	return imports
}

func getId(name string) int {
	id, ok := ids[name]
	if !ok {
		id = nextId
		nextId++
		ids[name] = id
	}
	return id
}

func hasPrefixes(s string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

func hasPackage(pkgpath string) bool {
	for k, _ := range ignored {
		if strings.Contains(pkgpath, k) {
			return true
		}
	}
	return false
}

func isIgnored(pkg *build.Package) bool {
	return ignored[pkg.ImportPath] || (pkg.Goroot && ignoreStdlib) || hasPrefixes(pkg.ImportPath, ignoredPrefixes) || hasPackage(pkg.ImportPath)
}

func debug(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
}

func debugf(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s, args...)
}

func getVendorlist(path string) []string {
	vendors := make([]string, 0)
	err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if !f.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, "vendor") && path != "" {
			vendors = append(vendors, PackageAbsPath(path))
		}
		return nil
	})
	if err != nil {
		glog.Errorf("filepath.Walk() returned %v\n", err)
	}
	return vendors
}

func PackageAbsPath(path string) (packagePath string) {
	_, err := os.Stat(path)
	if err != nil {
		glog.Errorln("package path is invalid")
		return ""
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		glog.Errorln(err)
	}
	packagePathIndex := strings.Index(absPath, "src")
	if -1 != packagePathIndex {
		packagePath = absPath[(packagePathIndex + 4):]
	}

	return packagePath
}
