# barely [![report](https://goreportcard.com/badge/github.com/reconquest/barely)](https://goreportcard.com/report/reconquest/barely)

Dead simple but yet extensible status bar for displaying interactive progress
for the shell-based tools, written in Go-lang.

![barely-example](https://cloud.githubusercontent.com/assets/674812/16452828/4c788a68-3e2c-11e6-8247-19e3db8f71fe.gif)

# Example

```go
package main

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/reconquest/barely"
	"github.com/reconquest/loreley"
)

func main() {
	format, err := loreley.CompileWithReset(
		` {bg 2}{fg 15}{bold} {.Mode} `+
			`{bg 253}{fg 0} `+
			`{if .Updated}{fg 70}{end}{.Done}{fg 0}/{.Total} `,
		nil,
	)
	if err != nil {
		panic(err)
	}

	var (
		bar = barely.NewStatusBar(format.Template)
		wg  = sync.WaitGroup{}

		status = &struct {
			Mode  string
			Total int
			Done  int64

			Updated int64
		}{
			Mode:  "PROCESSING",
			Total: 100,
		}
	)

	bar.SetStatus(status)
	bar.Render(os.Stderr)

	for i := 1; i <= status.Total; i++ {
		wg.Add(1)

		go func(i int) {
			sleep := time.Duration(rand.Intn(i)) * time.Millisecond * 300
			time.Sleep(sleep)

			atomic.AddInt64(&status.Done, 1)
			atomic.AddInt64(&status.Updated, 1)

			if i%10 == 0 {
				bar.Clear(os.Stderr)
				fmt.Printf("go-routine %d done (%s sleep)\n", i, sleep)
			}

			bar.Render(os.Stderr)

			wg.Done()

			<-time.After(time.Millisecond * 500)
			atomic.AddInt64(&status.Updated, -1)
			bar.Render(os.Stderr)
		}(i)
	}

	wg.Wait()

	bar.Clear(os.Stderr)
}
```
