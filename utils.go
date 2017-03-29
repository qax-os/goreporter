package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func DirList(path string, suffix, expect string) (dirs map[string]string, err error) {
	dirs = make(map[string]string, 0)

	err = filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, suffix) {
			dir := path[0:strings.LastIndex(path, system)]
			if ExpectPkg(expect, dir) {
				return nil
			}
			dirs[PackageAbsPath(dir)] = dir
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return dirs, nil
}

func ExpectPkg(expect, pkg string) bool {
	if expect == "" {
		return false
	}
	expects := strings.Split(expect, ",")
	for _, va := range expects {
		if strings.Contains(pkg, va) {
			return true
		}
	}
	return false
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

func PackageAbsPathExceptSuffix(path string) (packagePath string) {
	if strings.LastIndex(path, system) <= 0 {
		path, _ = os.Getwd()
	}
	path = path[0:strings.LastIndex(path, system)]
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

func projectName(projectPath string) (project string) {
	absPath, err := filepath.Abs(projectPath)
	if err != nil {
		fmt.Println(err)
	}
	projectPathIndex := strings.Index(absPath, "src")
	if -1 != projectPathIndex {
		project = absPath[(projectPathIndex + 4):]
	}

	return project
}

func absPath(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Println(err)
		return path
	}
	return absPath
}
