package display

import (
	"strings"
	"testing"

	"github.com/chriso345/gore/assert"
)

func TestAnsiHelp_NoFormat(t *testing.T) {
	input := "Hello, World!"
	output := ansiHelp(input)
	assert.Equal(t, output, input)
}

func TestAnsiHelp_SingleFormat(t *testing.T) {
	input := "Hello"
	expectedPrefix := string(ansiBold)
	expectedSuffix := string(ansiReset)

	output := ansiHelp(input, ansiBold)

	if !strings.HasPrefix(output, expectedPrefix) {
		t.Errorf("output does not start with expected ANSI code %q: got %q", expectedPrefix, output)
	}
	if !strings.HasSuffix(output, expectedSuffix) {
		t.Errorf("output does not end with ANSI reset code %q: got %q", expectedSuffix, output)
	}
	assert.StringContains(t, output, input)
}

func TestAnsiHelp_MultipleFormats(t *testing.T) {
	input := "Test"
	output := ansiHelp(input, ansiBold, ansiUnderline)

	// Expect output to start with concatenation of ansiBold + ansiUnderline
	expectedPrefix := string(ansiBold) + string(ansiUnderline)
	if !strings.HasPrefix(output, expectedPrefix) {
		t.Errorf("output does not start with expected ANSI codes %q: got %q", expectedPrefix, output)
	}
	if !strings.HasSuffix(output, string(ansiReset)) {
		t.Errorf("output does not end with ANSI reset code: got %q", output)
	}
	assert.StringContains(t, output, input)
}
