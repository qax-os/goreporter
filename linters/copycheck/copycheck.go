package copycheck

import (
	"bufio"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/360EntSecGroup-Skylar/goreporter/linters/copycheck/job"
	"github.com/360EntSecGroup-Skylar/goreporter/linters/copycheck/output"
	"github.com/360EntSecGroup-Skylar/goreporter/linters/copycheck/syntax"
	"github.com/golang/glog"
)

const defaultThreshold = 15

var (
	//flag.Bool("vendor", false, "check files in vendor directory")
	vendor = false
	//flag.Bool("verbose", false, "explain what is being done")
	verbose = false
	//flag.Int("threshold", defaultThreshold, "minimum token sequence as a clone")
	threshold = 50
	//flag.Bool("files", false, "files names from stdin")
	files = false

	// flag.Bool("html", false, "html output")
	html = false
	//flag.Bool("plumbing", false, "plumbing output for consumption by scripts or tools")
	plumbing = false
)

const (
	vendorDirPrefix = "vendor" + string(filepath.Separator)
	vendorDirInPath = string(filepath.Separator) + "vendor" + string(filepath.Separator)
)

func CopyCheck(projectPath string, except string) (result [][]string) {
	flag.Parse()
	if html && plumbing {
		glog.Errorln("you can have either plumbing or HTML output")
		return result
	}
	paths := []string{projectPath}
	if verbose {
		glog.Errorln("Building suffix tree")
	}
	schan := job.Parse(filesFeed(paths, except))
	t, data, done := job.BuildTree(schan)
	<-done

	// finish stream
	t.Update(&syntax.Node{Type: -1})

	if verbose {
		glog.Errorln("Searching for clones")
	}
	mchan := t.FindDuplOver(threshold)
	duplChan := make(chan syntax.Match)
	go func() {
		for m := range mchan {
			match := syntax.FindSyntaxUnits(*data, m, threshold)
			if len(match.Frags) > 0 {
				duplChan <- match
			}
		}
		close(duplChan)
	}()
	return printDupls(duplChan)
}

func filesFeed(paths []string, except string) chan string {
	if files {
		fchan := make(chan string)
		go func() {
			s := bufio.NewScanner(os.Stdin)
			for s.Scan() {
				f := s.Text()
				f = strings.TrimPrefix(f, "./")
				fchan <- f
			}
			close(fchan)
		}()
		return fchan
	}
	return crawlPaths(paths, except)
}

func crawlPaths(paths []string, except string) chan string {
	fchan := make(chan string)
	go func() {
		for _, path := range paths {
			info, err := os.Lstat(path)
			if err != nil {
				glog.Errorln(err)
				break
			}
			if !info.IsDir() {
				fchan <- path
				continue
			}
			filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
				if checkExcept(path, except) {
					return nil
				}
				if !vendor && (strings.HasPrefix(path, vendorDirPrefix) ||
					strings.Contains(path, vendorDirInPath)) {
					return nil
				}
				if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") && !strings.HasSuffix(info.Name(), "_test.go") {
					fchan <- path
				}
				return nil
			})
		}
		close(fchan)
	}()
	return fchan
}

func checkExcept(path, except string) bool {
	if except == "" || except == " " {
		return false
	}
	excepts := strings.Split(except, ",")
	for _, val := range excepts {
		if val != "" && val != " " {
			if strings.Contains(path, val) {
				return true
			}
		}
	}
	return false
}

type LocalFileReader struct{}

func (LocalFileReader) ReadFile(node *syntax.Node) ([]byte, error) {
	return ioutil.ReadFile(node.Filename)
}

func printDupls(duplChan <-chan syntax.Match) (copys [][]string) {
	groups := make(map[string][][]*syntax.Node)
	for dupl := range duplChan {
		groups[dupl.Hash] = append(groups[dupl.Hash], dupl.Frags...)
	}
	keys := make([]string, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	p := getPrinter()
	for _, k := range keys {
		uniq := unique(groups[k])
		if len(uniq) > 1 {
			// p.Print(uniq)
			copys = append(copys, p.SPrint(uniq))
		}
	}
	// p.Finish()
	return copys
}

func getPrinter() output.Printer {
	fr := LocalFileReader{}
	if html {
		return output.NewHtmlPrinter(os.Stdout, fr)
	} else if plumbing {
		return output.NewPlumbingPrinter(os.Stdout, fr)
	}
	return output.NewTextPrinter(os.Stdout, fr)
}

func unique(group [][]*syntax.Node) [][]*syntax.Node {
	fileMap := make(map[string]map[int]struct{})

	var newGroup [][]*syntax.Node
	for _, seq := range group {
		node := seq[0]
		file, ok := fileMap[node.Filename]
		if !ok {
			file = make(map[int]struct{})
			fileMap[node.Filename] = file
		}
		if _, ok := file[node.Pos]; !ok {
			file[node.Pos] = struct{}{}
			newGroup = append(newGroup, seq)
		}
	}
	return newGroup
}
