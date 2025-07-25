package clifford

import (
	"fmt"
	"reflect"
	"strings"
)

// BuildHelp generates and returns a formatted help message for a CLI tool
// defined by the given struct pointer.
//
// The `target` must be a pointer to a struct that embeds a `Clifford` field
// with a `name` tag. This tag specifies the CLI tool's name and is displayed
// in the usage header.
//
// The function inspects the struct to determine CLI arguments and options,
// including those marked as required. It outputs a help string that includes:
//   - The usage line with the command name and expected arguments
//   - A section for required arguments (based on `Required` tags)
//   - A section for optional flags (based on `short` or `long` tags)
//
// If no `name` tag is found on any embedded `Clifford` field, the function
// returns an error.
//
// Example:
//
//	target := struct {
//		clifford.Clifford `name:"mytool"`
//
//		Filename struct {
//			Value    string
//			clifford.Required
//			clifford.Desc `desc:"Input file path"`
//		}
//
//		Verbose struct {
//			Value    bool
//			clifford.Clifford `short:"v" long:"verbose" desc:"Enable verbose output"`
//		}
//	}{}
//
//	helpText, err := clifford.BuildHelp(&target)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(helpText)
func BuildHelp(target any) (string, error) {
	if !isStructPtr(target) {
		return "", fmt.Errorf("invalid type: must pass pointer to struct")
	}

	t := getStructType(target)

	// Find struct tag with `name`
	name := ""
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Tag.Get("name") != "" {
			name = field.Tag.Get("name")
			break
		}
	}
	if name == "" {
		return "", fmt.Errorf("struct must embed `Clifford` with `name` tag")
	}

	var builder strings.Builder
	builder.WriteString(ansiHelp("Usage:", ansiBold, ansiUnderline) + " ")
	builder.WriteString(ansiHelp(name, ansiBold))

	// Collect required args
	requiredArgs := getRequiredArgs(target)
	for _, arg := range requiredArgs {
		builder.WriteString(fmt.Sprintf(" [%s]", strings.ToUpper(arg)))
	}

	if hasOptions(target) {
		builder.WriteString(" [OPTIONS]")
	}
	builder.WriteString("\n")

	if len(requiredArgs) > 0 {
		builder.WriteString("\n" + ansiHelp("Arguments:", ansiBold, ansiUnderline) + "\n")
		builder.WriteString(argsHelp(target))
	}

	if hasOptions(target) {
		builder.WriteString("\n" + ansiHelp("Options:", ansiBold, ansiUnderline) + "\n")
		builder.WriteString(optionsHelp(target))
	}

	return builder.String(), nil
}

// === HELPERS ===

// argsHelp generates help text for positional arguments in the target struct.
func argsHelp(target any) string {
	t := getStructType(target)

	var lines []string
	maxLen := 0

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Name() == "Clifford" || field.Type.Name() == "Version" || field.Type.Name() == "Help" {
			continue
		}

		if field.Type.Kind() != reflect.Struct {
			continue
		}

		tags := getTagsFromEmbedded(field.Type, field.Name)
		if tags["short"] != "" || tags["long"] != "" {
			continue
		}

		argName := field.Name
		desc := tags["desc"]
		line := fmt.Sprintf("  [%s]", strings.ToUpper(argName))
		if len(line) > maxLen {
			maxLen = len(line)
		}
		lines = append(lines, fmt.Sprintf("%s||%s", line, desc))
	}

	// Format with aligned colons
	var builder strings.Builder
	for _, line := range lines {
		parts := strings.SplitN(line, "||", 2)
		padding := strings.Repeat(" ", maxLen-len(parts[0])+1)
		builder.WriteString(fmt.Sprintf("%s%s %s\n", parts[0], padding, parts[1]))
	}
	return builder.String()
}

// optionsHelp generates help text for options in the target struct.
func optionsHelp(target any) string {
	t := getStructType(target)

	var lines []string
	maxLen := 0

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Name() == "Clifford" {
			if field.Tag.Get("version") != "" {
				curr := "  --version||Show version information"
				lines = append(lines, curr)
				if 11 > maxLen {
					maxLen = 11
				}
			}
			if field.Tag.Get("help") != "" {
				curr := "  --help||Show this help message"
				lines = append(lines, curr)
				if 8 > maxLen {
					maxLen = 8
				}
			}
			continue
		}

		if field.Type.Name() == "Version" {
			curr := "  --version||Show version information"
			lines = append(lines, curr)
			if 11 > maxLen {
				maxLen = 11
			}
			continue
		}

		if field.Type.Name() == "Help" {
			curr := "  --help||Show this help message"
			lines = append(lines, curr)
			if 8 > maxLen {
				maxLen = 8
			}
			continue
		}

		if field.Type.Kind() != reflect.Struct {
			continue
		}

		tags := getTagsFromEmbedded(field.Type, field.Name)
		if tags["short"] == "" && tags["long"] == "" {
			continue
		}

		short := tags["short"]
		long := tags["long"]
		desc := tags["desc"]
		typeHint := fmt.Sprintf("[%s]", strings.ToUpper(field.Name))

		var flag string
		switch {
		case short != "" && long != "":
			flag = fmt.Sprintf("  -%s, --%s %s", short, long, typeHint)
		case short != "":
			flag = fmt.Sprintf("  -%s %s", short, typeHint)
		case long != "":
			flag = fmt.Sprintf("  --%s %s", long, typeHint)
		}

		if len(flag) > maxLen {
			maxLen = len(flag)
		}
		lines = append(lines, fmt.Sprintf("%s||%s", flag, desc))
	}

	// Format with aligned colons
	var builder strings.Builder
	for _, line := range lines {
		parts := strings.SplitN(line, "||", 2)
		padding := strings.Repeat(" ", maxLen-len(parts[0]))
		builder.WriteString(fmt.Sprintf("%s%s  %s\n", parts[0], padding, parts[1]))
	}
	return builder.String()
}

// getRequiredArgs returns a list of required argument names from the target struct.
func getRequiredArgs(target any) []string {
	t := getStructType(target)

	var args []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() != reflect.Struct {
			continue
		}

		tags := getTagsFromEmbedded(field.Type, field.Name)
		if tags["short"] != "" || tags["long"] != "" {
			continue
		}

		if _, ok := tags["required"]; ok {
			args = append(args, field.Name)
		}
	}
	return args
}

// hasOptions checks if the target struct has any options defined with short or long flags.
func hasOptions(target any) bool {
	t := getStructType(target)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() != reflect.Struct {
			continue
		}

		if field.Type.Name() == "Version" || field.Tag.Get("version") != "" {
			return true
		}
		if field.Type.Name() == "Help" || field.Tag.Get("help") != "" {
			return true
		}

		tags := getTagsFromEmbedded(field.Type, field.Name)
		if tags["short"] != "" || tags["long"] != "" {
			return true
		}
	}
	return false
}
