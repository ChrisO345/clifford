package spec

import (
	"fmt"
	"reflect"
)

// BuildHelp generates a help string for the provided target structure.
func BuildHelp(target any) (string, error) {
	helpText := ""

	helpText += ansiHelp("Usage:", ansiBold, ansiUnderline) + " "
	t := reflect.TypeOf(target)

	// Get the name of the cli from a meta tag.
	// TODO: Also support using the name of the struct by if no tag.
	tag_name, _ := t.FieldByName("Clifford")
	if tag_name.Tag.Get("name") == "" {
		return "", fmt.Errorf("invalid type: must pass pointer to struct")
	}
	name := tag_name.Tag.Get("name")
	if name == "" {
		return "", fmt.Errorf("invalid type: must pass pointer to struct with name tag")
	}

	helpText += ansiHelp(name, ansiBold)

	hasOptions := hasOptions(target)
	if hasOptions {
		helpText += ansiHelp(" [OPTIONS] ")
	}

	// TODO: Add positional arguments

	if hasOptions {
		helpText += "\n\n"
		helpText += ansiHelp("Options:", ansiBold, ansiUnderline) + "\n"
		helpText += buildOptions(target)
	}

	return helpText, nil
}

// hasOptions checks if the target struct has any options defined.
func hasOptions(_ any) bool {
	return true
}

// buildOptions generates a string representation of the options for the target struct.
func buildOptions(_ any) string {
	builder := ""
	return builder
}

// ansiHelp returns a string with ANSI escape codes for formatting.
func ansiHelp(text string, format ...ansiFormat) string {
	builder := ""
	if len(format) == 0 {
		return text
	}

	for _, f := range format {
		builder += string(f)
	}
	builder += text
	return builder + string(ansiReset)
}

// ansiFormat defines the format for ANSI escape codes.
type ansiFormat string

const (
	ansiReset     ansiFormat = "\033[0m"
	ansiBold      ansiFormat = "\033[1m"
	ansiDim       ansiFormat = "\033[2m"
	ansiItalic    ansiFormat = "\033[3m"
	ansiUnderline ansiFormat = "\033[4m"
	ansiRed       ansiFormat = "\033[31m"
	ansiGreen     ansiFormat = "\033[32m"
	ansiYellow    ansiFormat = "\033[33m"
	ansiBlue      ansiFormat = "\033[34m"
	ansiMagenta   ansiFormat = "\033[35m"
	ansiCyan      ansiFormat = "\033[36m"
	ansiWhite     ansiFormat = "\033[37m"
)
