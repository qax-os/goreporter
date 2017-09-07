package loreley

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	Colorize = ColorizeAlways
}

func TestCompile_CompilesEmptyStringToEmptyStyle(t *testing.T) {
	assertExecutedTemplate(t, ``, ``, nil)
}

func TestCompile_CompilesStringToString(t *testing.T) {
	assertExecutedTemplate(t, `root beer`, `root beer`, nil)
}

func TestCompile_CompilesAsGoTemplate(t *testing.T) {
	assertExecutedTemplate(
		t,
		`{if .sweet}bubblegum{end}`,
		`bubblegum`,
		map[string]interface{}{"sweet": true},
	)
}

func TestStyle_Execute_PutsEscapeSequenceToChangeBgColor(t *testing.T) {
	assertExecutedTemplate(t, `{bg 2}finn`, "\x1b[48;5;2mfinn", nil)
}

func TestStyle_Execute_PutsEscapeSequenceToChangeFgColor(t *testing.T) {
	assertExecutedTemplate(t, `{fg 2}finn`, "\x1b[38;5;2mfinn", nil)
}

func TestStyle_Execute_ResetsBackgroundColorToDefault(t *testing.T) {
	assertExecutedTemplate(t, `{nobg}finn`, "\x1b[49mfinn", nil)
}

func TestStyle_Execute_ResetsForegroundColorToDefault(t *testing.T) {
	assertExecutedTemplate(t, `{nofg}finn`, "\x1b[39mfinn", nil)
}

func TestStyle_Execute_PutsReverseCodeSequence(t *testing.T) {
	assertExecutedTemplate(t, `{reverse}jake`, "\x1b[7mjake", nil)
}

func TestStyle_Execute_PutsNoReverseCodeSequence(t *testing.T) {
	assertExecutedTemplate(t, `{noreverse}jake`, "\x1b[27mjake", nil)
}

func TestStyle_Execute_PutsBoldCodeSequence(t *testing.T) {
	assertExecutedTemplate(t, `{bold}jake`, "\x1b[1mjake", nil)
}

func TestStyle_Execute_PutsNoBoldCodeSequence(t *testing.T) {
	assertExecutedTemplate(t, `{nobold}finn`, "\x1b[22mfinn", nil)
}

func TestStyle_Execute_PutsResetCodeSequence(t *testing.T) {
	assertExecutedTemplate(t, `{reset}finn`, "\x1b[0mfinn", nil)
}

func TestStyle_Execute_PutsBgWithFg(t *testing.T) {
	assertExecutedTemplate(
		t,
		`{fg 1}{bg 2}finn`,
		"\x1b[38;5;1m\x1b[48;5;2mfinn", nil,
	)
}

func TestStyle_Execute_PutsTransitionStringFromOneBgToAnother(t *testing.T) {
	assertExecutedTemplate(
		t,
		`{fg 6}{bg 2}finn{from "" 4}jake`,
		"\x1b[38;5;6m\x1b[48;5;2mfinn"+
			"\x1b[38;5;2m\x1b[48;5;4m\x1b[38;5;6mjake", nil,
	)
}

func TestStyle_Execute_PutsTransitionStringFromOneBgToAnotherInverted(
	t *testing.T,
) {
	assertExecutedTemplate(
		t,
		`{fg 6}{bg 2}finn{to 4 ""}jake`,
		"\x1b[38;5;6m\x1b[48;5;2mfinn"+
			"\x1b[38;5;4m\x1b[48;5;4m\x1b[38;5;6mjake", nil,
	)
}

func TestCompileWithReset_AddsResetToEndOfStyle(t *testing.T) {
	test := assert.New(t)

	style, err := CompileWithReset(`123`, nil)
	test.Nil(err)

	actual, err := style.ExecuteToString(nil)
	test.Nil(err)
	test.Equal("123\x1b[0m", actual)
}

func TestTrimStyles_RemovesAllEscapeCodesFromString(t *testing.T) {
	assert.New(t).Equal(
		`finn`,
		TrimStyles(
			"\x1b[38;5;1m\x1b[48;5;2mfinn",
		),
	)
}

func assertExecutedTemplate(
	t *testing.T,
	template string,
	expected string,
	data map[string]interface{},
) {
	test := assert.New(t)

	style, err := Compile(template, nil)
	test.Nil(err)

	actual, err := style.ExecuteToString(data)
	test.Nil(err)
	test.Equal(expected, actual)
}
