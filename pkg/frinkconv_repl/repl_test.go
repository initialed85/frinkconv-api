package frinkconv_repl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestREPL(t *testing.T) {
	var r *REPL
	var output float64
	var err error

	r, err = New()
	assert.Nil(t, err)
	defer r.Close()

	output, err = r.Convert(10, "apples", "oranges")
	assert.Equal(
		t,
		0.0,
		output,
	)
	assert.NotNil(t, err)
	assert.Equal(
		t,
		"Warning: undefined symbol \"apples\".\nUnknown symbol \"oranges\"\nWarning: undefined symbol \"apples\".\nWarning: undefined symbol \"oranges\".\nUnconvertable expression:\n  10 apples (undefined symbol) -> oranges (undefined symbol)",
		err.Error())

	output, err = r.Convert(10, "feet", "pound feet")
	assert.Equal(
		t,
		0.0,
		output,
	)
	assert.NotNil(t, err)
	assert.Equal(
		t,
		"Conformance error\n   Left side is: 381/125 m (length)\n  Right side is: 17281869297/125000000000 m kg (unknown unit type)\n     Suggestion: multiply left side by mass\n\n For help, type: units[mass]\n                   to list known units with these dimensions.",
		err.Error())

	output, err = r.Convert(10, "feet", "inches")
	assert.Equal(
		t,
		120.0,
		output,
	)
	assert.Nil(t, err)

	output, err = r.Convert(120, "inches", "feet")
	assert.Equal(
		t,
		10.0,
		output,
	)
	assert.Nil(t, err)

	output, err = r.Convert(120, "feet", "metres")
	assert.Equal(
		t,
		36.576,
		output,
	)
	assert.Nil(t, err)

	r.Close()

	output, err = r.Convert(120, "inches", "feet")
	assert.NotNil(t, err)
}
