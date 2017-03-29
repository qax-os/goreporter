package main

import (
	"fmt"
)

func main() {
	report := NewReporter()
	jsonData := report.Engine("../goreporter", "")
	if jsonData == nil {
		fmt.Println("Engine error")
		return
	}
	jsonHtmlString, err := Json2Html(jsonData, "../goreporter")
	if err != nil {
		fmt.Println("Json2Html error")
		return
	}
	fmt.Println(jsonHtmlString)
}
