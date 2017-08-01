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
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestLinterProcessBar(t *testing.T) {
	lintersProcessChans := make(chan int64, 20)
	lintersFinishedSignal := make(chan string, 10)
	go LinterProcessBar(lintersProcessChans, lintersFinishedSignal)
	wg := sync.WaitGroup{}
	for i := 1; i <= 100; i++ {
		wg.Add(1)
		time.Sleep(100 * time.Microsecond)
		go func(i int) {
			lintersProcessChans <- 1
			if i%10 == 0 {
				lintersFinishedSignal <- fmt.Sprintf("go-routine %d done\n", i)
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
	close(lintersFinishedSignal)
	close(lintersProcessChans)
	return
}
