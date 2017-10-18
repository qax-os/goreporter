package spellcheck

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	"time"

	"github.com/360EntSecGroup-Skylar/goreporter/linters/spellcheck/misspell"
)

var (
	defaultWrite *template.Template
	defaultRead  *template.Template

	stdout     *log.Logger
	debug      *log.Logger
	spellCheck []string
)

const (
	// Note for gometalinter it must be "File:Line:Column: Msg"
	//  note space between ": Msg"
	defaultWriteTmpl = `{{ .Filename }}:{{ .Line }}:{{ .Column }}: corrected "{{ .Original }}" to "{{ .Corrected }}"`
	defaultReadTmpl  = `{{ .Filename }}:{{ .Line }}:{{ .Column }}: "{{ .Original }}" is a misspelling of "{{ .Corrected }}"`
	csvTmpl          = `{{ printf "%q" .Filename }},{{ .Line }},{{ .Column }},{{ .Original }},{{ .Corrected }}`
	csvHeader        = `file,line,column,typo,corrected`
	sqliteTmpl       = `INSERT INTO misspell VALUES({{ printf "%q" .Filename }},{{ .Line }},{{ .Column }},{{ printf "%q" .Original }},{{ printf "%q" .Corrected }});`
	sqliteHeader     = `PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE misspell(
	"file" TEXT, "line" INTEGER, "column" INTEGER, "typo" TEXT, "corrected" TEXT
);`
	sqliteFooter = "COMMIT;"
)

func worker(writeit bool, r *misspell.Replacer, mode string, files <-chan string, results chan<- int) {
	count := 0
	for filename := range files {
		orig, err := misspell.ReadTextFile(filename)
		if err != nil {
			log.Println(err)
			continue
		}
		if len(orig) == 0 {
			continue
		}

		debug.Printf("Processing %s", filename)

		updated, changes := r.Replace(orig)

		if len(changes) == 0 {
			continue
		}
		count += len(changes)
		for _, diff := range changes {
			// add in filename
			diff.Filename = filename

			// output can be done by doing multiple goroutines
			// and can clobber os.Stdout.
			//
			// the log package can be used simultaneously from multiple goroutines
			var output bytes.Buffer
			if writeit {
				defaultWrite.Execute(&output, diff)
			} else {
				defaultRead.Execute(&output, diff)
			}

			// goroutine-safe print to os.Stdout
			// stdout.Println(output.String())
			spellCheck = append(spellCheck, output.String())
		}

		if writeit {
			ioutil.WriteFile(filename, []byte(updated), 0)
		}
	}
	results <- count
}

func SpellCheck(projectPath, except string) []string {
	spellCheck = make([]string, 0)
	t := time.Now()
	var (
		workers   = 0
		writeit   = false
		quietFlag = false
		outFlag   = "stdout"
		format    = ""
		ignores   = ""
		locale    = ""
		mode      = "auto"
		debugFlag = false
		exitError = false
	)

	if debugFlag {
		debug = log.New(os.Stderr, "DEBUG ", 0)
	} else {
		debug = log.New(ioutil.Discard, "", 0)
	}

	r := misspell.Replacer{
		Replacements: misspell.DictMain,
		Debug:        debugFlag,
	}
	//
	// Figure out regional variations
	//
	switch strings.ToUpper(locale) {
	case "":
		// nothing
	case "US":
		r.AddRuleList(misspell.DictAmerican)
	case "UK", "GB":
		r.AddRuleList(misspell.DictBritish)
	case "NZ", "AU", "CA":
		log.Fatalf("Help wanted.  https://github.com/client9/misspell/issues/6")
	default:
		log.Fatalf("Unknow locale: %q", locale)
	}

	//
	// Stuff to ignore
	//
	_ = ignores
	// if len(ignores) > 0 {
	// 	r.RemoveRule(strings.Split(ignores, ","))
	// }

	//
	// Source input mode
	//
	switch mode {
	case "auto":
	case "go":
	case "text":
	default:
		log.Fatalf("Mode must be one of auto=guess, go=golang source, text=plain or markdown-like text")
	}

	//
	// Custom output
	//
	switch {
	case format == "csv":
		tmpl := template.Must(template.New("csv").Parse(csvTmpl))
		defaultWrite = tmpl
		defaultRead = tmpl
		stdout.Println(csvHeader)
	case format == "sqlite" || format == "sqlite3":
		tmpl := template.Must(template.New("sqlite3").Parse(sqliteTmpl))
		defaultWrite = tmpl
		defaultRead = tmpl
		stdout.Println(sqliteHeader)
	case len(format) > 0:
		t, err := template.New("custom").Parse(format)
		if err != nil {
			log.Fatalf("Unable to compile log format: %s", err)
		}
		defaultWrite = t
		defaultRead = t
	default: // format == ""
		defaultWrite = template.Must(template.New("defaultWrite").Parse(defaultWriteTmpl))
		defaultRead = template.Must(template.New("defaultRead").Parse(defaultReadTmpl))
	}

	// we cant't just write to os.Stdout directly since we have multiple goroutine
	// all writing at the same time causing broken output.  Log is routine safe.
	// we see it so it doesn't use a prefix or include a time stamp.
	switch {
	case quietFlag || outFlag == "/dev/null":
		stdout = log.New(ioutil.Discard, "", 0)
	case outFlag == "/dev/stderr" || outFlag == "stderr":
		stdout = log.New(os.Stderr, "", 0)
	case outFlag == "/dev/stdout" || outFlag == "stdout":
		stdout = log.New(os.Stdout, "", 0)
	case outFlag == "" || outFlag == "-":
		stdout = log.New(os.Stdout, "", 0)
	default:
		fo, err := os.Create(outFlag)
		if err != nil {
			log.Fatalf("unable to create outfile %q: %s", outFlag, err)
		}
		defer fo.Close()
		stdout = log.New(fo, "", 0)
	}

	//
	// Number of Workers / CPU to use
	//
	if workers < 0 {
		log.Fatalf("-j must >= 0")
	}
	if workers == 0 {
		workers = runtime.NumCPU()
	}
	if debugFlag {
		workers = 1
	}

	//
	// Done with Flags.
	//  Compile the Replacer and process files
	//
	r.Compile()

	args := []string{projectPath}
	debug.Printf("initialization complete in %v", time.Since(t))

	// stdin/stdout
	if len(args) == 0 {
		// if we are working with pipes/stdin/stdout
		// there is no concurrency, so we can directly
		// send data to the writers
		var fileout io.Writer
		var errout io.Writer
		switch writeit {
		case true:
			// if we ARE writing the corrected stream
			// the the corrected stream goes to stdout
			// and the misspelling errors goes to stderr
			// so we can do something like this:
			// curl something | misspell -w | gzip > afile.gz
			fileout = os.Stdout
			errout = os.Stderr
		case false:
			// if we are not writing out the corrected stream
			// then work just like files.  Misspelling errors
			// are sent to stdout
			fileout = ioutil.Discard
			errout = os.Stdout
		}
		count := 0
		next := func(diff misspell.Diff) {
			count++

			// don't even evaluate the output templates
			if quietFlag {
				return
			}
			diff.Filename = "stdin"
			if writeit {
				defaultWrite.Execute(errout, diff)
			} else {
				defaultRead.Execute(errout, diff)
			}
			errout.Write([]byte{'\n'})

		}
		err := r.ReplaceReader(os.Stdin, fileout, next)
		if err != nil {
			return spellCheck
		}
		switch format {
		case "sqlite", "sqlite3":
			fileout.Write([]byte(sqliteFooter))
		}
		if count != 0 && exitError {
			// error
			return spellCheck
		}
		return spellCheck
	}

	c := make(chan string, 64)
	results := make(chan int, workers)

	for i := 0; i < workers; i++ {
		go worker(writeit, &r, mode, c, results)
	}
	excepts := setExcept(except)
	for _, filename := range args {
		filepath.Walk(filename, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() && strings.HasSuffix(path, ".go") && !checkExcept(path, excepts) {
				c <- path
			}
			return nil
		})
	}
	close(c)

	return spellCheck
}

func checkExcept(path string, excepts []string) bool {
	if path == "" || path == " " {
		return false
	}
	for _, val := range excepts {
		if val != "" && val != " " {
			return strings.Contains(path, val)
		}
	}
	return false
}

func setExcept(except string) (excepts []string) {
	excepts = append(excepts, "vendor")
	if except == "" || except == " " {
		return excepts
	}

	return append(excepts, strings.Split(except, ",")...)
}
