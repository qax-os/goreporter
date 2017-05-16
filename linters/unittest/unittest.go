package unittest

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
)

func UnitTest(packagePath string) (packageUnitTestResults map[string][]string, packageTestRaceResults map[string][]string) {
	packageUnitTestResults = make(map[string][]string, 0)
	packageTestRaceResults = make(map[string][]string, 0)

	packageName := PackageAbsPath(packagePath)
	if "" == packageName {
		packageName = packagePath
	}

	out, err := GoTestWithCoverAndRace(packagePath)
	if err != nil {
		if !strings.Contains(out, "==================") {
			glog.Infoln("[UnitTest] package->:", packageName, " ... ", err)
		} else {
			glog.Infoln("[UnitTest] package->:", packageName, " ... pass")
		}
	} else {
		glog.Infoln("[UnitTest] package->:", packageName, " ... pass")
	}

	if out == "" || !strings.Contains(out, "ok") {
		packageUnitTestResults[packageName] = []string{}
	} else if err != nil {
		lindex := strings.LastIndex(out, "coverage:")
		res := strings.Split(out[lindex:], "\n")
		info := strings.Fields(res[2])
		cov := strings.Fields(res[0])

		if len(info) >= 3 && len(cov) >= 2 {
			rest := info[0] + " " + info[1] + " " + info[2] + " " + cov[0] + " " + cov[1]
			packageUnitTestResults[packageName] = strings.Fields(rest)

			for in, val := range strings.Split(out, "==================") {
				if (in+1)%2 == 0 {
					packageTestRaceResults[packageName] = append(packageTestRaceResults[packageName], val)
				}
			}
		}
	} else {
		test := strings.Fields(out)
		packageUnitTestResults[packageName] = test
	}

	return packageUnitTestResults, packageTestRaceResults
}

// run go test -cover
func GoTestWithCoverAndRace(packagePath string) (packageUnitResult string, err error) {
	cmd := exec.Command("go", "test", packagePath, "-cover", "-race")
	var out bytes.Buffer
	cmd.Stdout = &out
	// cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		if !strings.Contains(out.String(), "==================") {
			return "", err
		}
	}

	return out.String(), err
}

// run go list -cover
func GoListWithImportPackages(packagePath string) (importPackages []string) {
	importPackages = make([]string, 0)
	cmd := exec.Command("go", "list", "-f", `'{{ join .Imports " " }}'`, packagePath)
	var out bytes.Buffer
	cmd.Stdout = &out
	// cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		glog.Warningln(err)
		return importPackages
	}
	packagesString := out.String()
	packagesString = strings.Replace(packagesString, `'`, "", -1)
	packages := strings.Fields(packagesString)

	var out2 bytes.Buffer
	cmd = exec.Command("go", "list", "std")
	cmd.Stdout = &out2
	// cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		glog.Warningln(err)
		return importPackages
	}
	stdPackages := strings.Split(out2.String(), "\n")
	mapStdPackages := make(map[string]string, 0)
	for i := 0; i < len(stdPackages); i++ {
		mapStdPackages[stdPackages[i]] = stdPackages[i]
	}
	// remove std package
	for i := 0; i < len(packages); i++ {
		if strings.Contains(packages[i], string(filepath.Separator)) && !strings.Contains(packages[i], "vendor") {
			if _, ok := mapStdPackages[packages[i]]; !ok {
				importPackages = append(importPackages, strings.Replace(packages[i], "'", "", -1))
			}
		}
	}

	return importPackages
}

func PackageAbsPath(path string) (packagePath string) {
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
