package tools

import (
	"os"
	"encoding/json"

	"github.com/fatih/color"
	"github.com/wgliang/goreporter/engine"
)

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
	errorInfoTpl = `  %s at line %d`
)

// DisplayAsText will display the json data to console
func DisplayAsText(jsonData []byte) {
	var structData engine.Reporter
	json.Unmarshal(jsonData, &structData)

	color.Magenta(
		headerTpl, 
		structData.Project, 
		structData.Score, 
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

	if structData.Issues > 0 {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}