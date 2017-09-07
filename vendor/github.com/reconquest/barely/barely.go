package barely

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"sync"
	"text/template"
)

var (
	escapeSequenceRegexp = regexp.MustCompile("\x1b[^m]+m")
)

// StatusBar is the implementation of configurable status bar.
type StatusBar struct {
	sync.Locker

	format *template.Template
	last   string

	status interface{}
}

// NewStatusBar returns new StatusBar object, initialized with given template.
//
// Template will be used later on Render calls.
func NewStatusBar(format *template.Template) *StatusBar {
	return &StatusBar{
		format: format,
	}
}

// Lock locks StatusBar object if locker object was set with SetLock method
// to prevent multi-threading race conditions.
//
// StatusBar will be locked in the Set and Render methods.
func (bar *StatusBar) Lock() {
	if bar.Locker != nil {
		bar.Locker.Lock()
	}
}

// Unlock unlocks previously locked StatusBar object.
func (bar *StatusBar) Unlock() {
	if bar.Locker != nil {
		bar.Locker.Unlock()
	}
}

// SetLock sets locker object, that will be used for Lock and Unlock methods.
func (bar *StatusBar) SetLock(lock sync.Locker) {
	bar.Locker = lock
}

// SetStatus sets data which will be used in the template execution, which is
// previously set through NewStatusBar function.
func (bar *StatusBar) SetStatus(data interface{}) {
	bar.Lock()
	defer bar.Unlock()

	bar.status = data
}

// Render renders specified template and writes it to the specified writer.
//
// Also, render result will be remembered and will be used to generate clear
// sequence which can be obtained from Clear method call.
func (bar *StatusBar) Render(writer io.Writer) error {
	bar.Lock()
	defer bar.Unlock()

	buffer := &bytes.Buffer{}

	if bar.status == nil {
		return nil
	}

	err := bar.format.Execute(buffer, bar.status)
	if err != nil {
		return fmt.Errorf(
			`error during rendering status bar: %s`,
			err,
		)
	}

	fmt.Fprintf(buffer, "\r")

	bar.last = escapeSequenceRegexp.ReplaceAllLiteralString(
		buffer.String(),
		``,
	)

	_, err = io.Copy(writer, buffer)
	if err != nil {
		return fmt.Errorf(
			`can't write status bar: %s`,
			err,
		)
	}

	return nil
}

// Clear writes clear sequence in the specified writer, which is represented by
// terminal erase line sequence.
func (bar *StatusBar) Clear(writer io.Writer) {
	bar.Lock()
	defer bar.Unlock()

	if bar.last != "" {
		fmt.Fprint(writer, "\x1b[2K")
	}

	bar.last = ""
}
