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
	"encoding/json"

	"github.com/360EntSecGroup-Skylar/goreporter/engine"
	"github.com/fatih/color"
)

// Text display description and logo.
const (
	headerTpl = `
		                                                                                                        
                                                                            _/                          
     _/_/_/    _/_/    _/  _/_/    _/_/    _/_/_/      _/_/    _/  _/_/  _/_/_/_/    _/_/    _/  _/_/   
  _/    _/  _/    _/  _/_/      _/_/_/_/  _/    _/  _/    _/  _/_/        _/      _/_/_/_/  _/_/        
 _/    _/  _/    _/  _/        _/        _/    _/  _/    _/  _/          _/      _/        _/           
  _/_/_/    _/_/    _/          _/_/_/  _/_/_/      _/_/    _/            _/_/    _/_/_/  _/            
     _/                                _/                                                               
_/_/                                  _/                                                                
		
	Project: %s 
	Score: %d
	Grade: %d
	Time: %s
	Issues: %d

	`
	metricsHeaderTpl = `>> %s Linter %s find:`
	summaryHeaderTpl = ` %s: %s`
	errorInfoTpl     = `  %s at line %d`
)

// DisplayAsText will display the json data to console.In your CI, all
// the tips will be given a variety of color tips, very beautiful.
func DisplayAsText(jsonData []byte) {
	var structData engine.Reporter
	json.Unmarshal(jsonData, &structData)

	var score float64
	for _, metric := range structData.Metrics {
		score = score + metric.Percentage*metric.Weight
	}

	color.Magenta(
		headerTpl,
		structData.Project,
		int(score),
		structData.Grade,
		structData.TimeStamp,
		structData.Issues,
	)
	for _, metric := range structData.Metrics {
		if metric.Name == "DependGraph" || 0 == len(metric.Summaries) {
			continue
		}
		color.Cyan(metricsHeaderTpl, metric.Name, metric.Description)
		for _, summary := range metric.Summaries {
			color.Blue(summaryHeaderTpl, summary.Name, summary.Description)
			for _, errorInfo := range summary.Errors {
				color.Red(errorInfoTpl, errorInfo.ErrorString, errorInfo.LineNumber)
			}
		}
	}
}
