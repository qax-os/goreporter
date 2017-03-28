package unittest

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

func UnitTest(packagePath string) (packageUnitTestResults map[string][]string, packageTestRaceResults map[string][]string) {
	packageUnitTestResults = make(map[string][]string, 0)
	packageTestRaceResults = make(map[string][]string, 0)

	fmt.Println("Testing package " + packagePath + " ...")

	packageName := PackageAbsPath(packagePath)
	if "" == packageName {
		packageName = packagePath
	}
	out, err := GoTestWithCoverAndRace(packagePath)
	if err != nil {
		if !strings.Contains(out, "==================") {
			fmt.Println(err)
		}
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
			fmt.Println(err)
			return "", err
		}
	}
	fmt.Println(err)
	return out.String(), err
}

func PackageAbsPath(path string) (packagePath string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Println(err)
	}
	packagePathIndex := strings.Index(absPath, "src")
	if -1 != packagePathIndex {
		packagePath = absPath[(packagePathIndex + 4):]
	}

	return packagePath
}
