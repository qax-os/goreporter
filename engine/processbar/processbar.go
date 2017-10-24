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

package processbar

import (
	"log"
	"os"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/reconquest/barely"
	"github.com/reconquest/loreley"
)

// LinterProcessBar provides a function to display the current progress of
// generating your project, you can predict the time to generate your report.
// And it will communicate with engine through chans(lintersProcessChans and
// lintersFinishedSignal).
func LinterProcessBar(lintersProcessChans chan int64, lintersFinishedSignal chan string) {
	// Set process bar corlor,for a better experience.
	format, err := loreley.CompileWithReset(
		` {bg 2}{fg 15}{bold} {.Mode} `+
			`{bg 253}{fg 0} `+
			`{if .Updated}{fg 70}{end}{.Done}{fg 0}/{.Total} `,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	// Initialize some of the parameters, set some of the properties of our progress bar.
	var (
		bar = barely.NewStatusBar(format.Template)

		status = &struct {
			Mode  string
			Total int
			Done  int64

			Updated int64
		}{
			Mode:  "GOREPORTER",
			Total: 100,
		}
	)

	bar.SetStatus(status)
	processCount := int64(100)

	if runtime.GOOS == "windows" {
	PROCESSNORENDER:
		for {
			select {
			case process := <-lintersProcessChans:
				atomic.AddInt64(&status.Done, process)
				atomic.AddInt64(&status.Updated, process)
				processCount = processCount - process
				if processCount <= int64(0) {
					break PROCESSNORENDER
				}
			case signal := <-lintersFinishedSignal:
				log.Println(signal)
			}
		}
	} else {
		bar.Render(os.Stderr)
	PROCESSRENDER:
		for {
			select {
			case process := <-lintersProcessChans:
				atomic.AddInt64(&status.Done, process)
				atomic.AddInt64(&status.Updated, process)
				bar.Render(os.Stderr)
				processCount = processCount - process
				if processCount <= int64(0) {
					break PROCESSRENDER
				}
			case signal := <-lintersFinishedSignal:
				log.Println(signal)
				bar.Render(os.Stderr)
			case <-time.After(1 * time.Second):
				bar.Render(os.Stderr)
			}
		}
	}

	bar.Clear(os.Stderr)
	return
}
