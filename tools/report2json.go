package tools

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
	"github.com/wgliang/goreporter/engine"
)

// SaveAsJson will save as a json file with Reporter struct.
func SaveAsJson(jsonData []byte, projectPath, savePath, timestamp string) {
	saveAbsPath := engine.AbsPath(savePath)
	fmt.Println(savePath)
	fmt.Println(saveAbsPath)
	projectName := engine.ProjectName(projectPath)
	if saveAbsPath != "" && saveAbsPath != savePath {
		jsonpath := strings.Replace(saveAbsPath+string(filepath.Separator)+projectName+"-"+timestamp+".json", string(filepath.Separator)+string(filepath.Separator), string(filepath.Separator), -1)
		err := ioutil.WriteFile(jsonpath, jsonData, 0666)
		if err != nil {
			glog.Errorln(err)
		}
	} else {
		jsonpath := projectName + "-" + timestamp + ".json"
		err := ioutil.WriteFile(jsonpath, jsonData, 0666)
		if err != nil {
			glog.Errorln(err)
		}
	}
}
