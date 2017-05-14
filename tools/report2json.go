package tools

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
	"github.com/wgliang/goreporter/engine"
)

// SaveAsJson will save as a json file with Reporter struct.
func SaveAsJson(jsonData []byte, projectPath, savePath, timestamp string) {
	savePath = engine.AbsPath(savePath)
	projectName := engine.ProjectName(projectPath)
	if savePath != "" {
		jsonpath := strings.Replace(savePath+string(filepath.Separator)+projectName+"-"+timestamp+".json", string(filepath.Separator)+string(filepath.Separator), string(filepath.Separator), -1)
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
