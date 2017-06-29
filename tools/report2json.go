// Copyright 2017 The GoReporter Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tools

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
	"github.com/wgliang/goreporter/engine"
)

// SaveAsJson is a function that save data as a json report.And will receive
// jsonData, projectPath, savePath and timestamp parameters.
func SaveAsJson(jsonData []byte, projectPath, savePath, timestamp string) {
	saveAbsPath := engine.AbsPath(savePath)
	projectName := engine.ProjectName(projectPath)
	if saveAbsPath != "" && saveAbsPath != savePath {
		jsonpath := strings.Replace(saveAbsPath+string(filepath.Separator)+projectName+"-"+timestamp+".json", string(filepath.Separator)+string(filepath.Separator), string(filepath.Separator), -1)
		err := ioutil.WriteFile(jsonpath, jsonData, 0666)
		if err != nil {
			glog.Errorln(err)
		} else {
			glog.Info("Json report saved in:", jsonpath)
		}
	} else {
		jsonpath := projectName + "-" + timestamp + ".json"
		err := ioutil.WriteFile(jsonpath, jsonData, 0666)
		if err != nil {
			glog.Errorln(err)
		} else {
			glog.Info("Json report saved in:", jsonpath)
		}

	}
}
