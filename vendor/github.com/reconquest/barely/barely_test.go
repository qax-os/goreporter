package barely

import (
	"bytes"
	"strings"
	"testing"

	"github.com/reconquest/loreley"
	"github.com/stretchr/testify/assert"
)

func TestStatusBar_Render_RendersTemplateIntoWriter(t *testing.T) {
	assertStatusBarRendering(
		t,
		`{.Blah} test`,
		"foo test\r",
		struct {
			Blah string
		}{
			Blah: "foo",
		},
	)
}

func TestStatusBar_Render_RendersTemplateIntoWriterWithColors(t *testing.T) {
	test := assert.New(t)

	expected, err := loreley.CompileAndExecuteToString("{fg 1}foo test\r", nil, nil)
	test.NoError(err)

	assertStatusBarRendering(
		t,
		`{fg 1}{.Blah} test`,
		expected,
		struct {
			Blah string
		}{
			Blah: "foo",
		},
	)
}

func TestStatusBar_Clear_ErasesPreviousOutputWithSpaces(t *testing.T) {
	test := assert.New(t)

	bar := assertStatusBarRendering(
		t,
		`{.Blah} test`,
		"foo test\r",
		struct {
			Blah string
		}{
			Blah: "foo",
		},
	)

	buffer := bytes.Buffer{}

	bar.Clear(&buffer)

	test.Equal(strings.Repeat(` `, len(`foo test`))+"\r", buffer.String())
}

func TestStatusBar_Clear_ErasesPreviousOutputWithSpacesButSkipColors(t *testing.T) {
	test := assert.New(t)

	expected, err := loreley.CompileAndExecuteToString("{fg 1}foo test\r", nil, nil)
	test.NoError(err)

	bar := assertStatusBarRendering(
		t,
		`{fg 1}{.Blah} test`,
		expected,
		struct {
			Blah string
		}{
			Blah: "foo",
		},
	)

	buffer := bytes.Buffer{}

	bar.Clear(&buffer)

	test.Equal(strings.Repeat(` `, len(`foo test`))+"\r", buffer.String())
}

func assertStatusBarRendering(
	t *testing.T,
	formatting string,
	expected string,
	data interface{},
) *StatusBar {
	test := assert.New(t)

	buffer := bytes.Buffer{}

	style, err := loreley.Compile(formatting, nil)
	test.NoError(err)

	bar := NewStatusBar(style.Template)

	err = bar.Render(&buffer)
	test.NoError(err)

	test.Equal(``, buffer.String())

	bar.SetStatus(data)

	err = bar.Render(&buffer)
	test.NoError(err)

	test.Equal(expected, buffer.String())

	return bar
}
