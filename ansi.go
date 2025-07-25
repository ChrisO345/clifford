package clifford

import "strings"

// ansiHelp formats the text with ANSI escape codes for styling.
func ansiHelp(text string, format ...ansiFormat) string {
	if len(format) == 0 {
		return text
	}
	var builder strings.Builder
	for _, f := range format {
		builder.WriteString(string(f))
	}
	builder.WriteString(text)
	builder.WriteString(string(ansiReset))
	return builder.String()
}

// ansiFormat defines the ANSI escape codes for text formatting.
type ansiFormat string

const (
	ansiReset     ansiFormat = "\033[0m"
	ansiBold      ansiFormat = "\033[1m"
	ansiUnderline ansiFormat = "\033[4m"
)
