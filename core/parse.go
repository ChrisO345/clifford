package core

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/chriso345/clifford/display"
	"github.com/chriso345/clifford/internal/common"
)

var osExit = os.Exit // Mockable for testing

func Parse(target any) error {
	if !common.IsStructPtr(target) {
		return fmt.Errorf("invalid type: must pass pointer to struct")
	}

	args := os.Args[1:]

	// For test mode: drop args before "--"
	if i := common.ArgsIndexOf(args, "--"); i >= 0 {
		args = args[i+1:]
	}

	// Pre-process args into maps for fast lookup
	argMap := map[string]string{}
	argFlags := map[string]bool{}
	used := map[int]bool{}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "--") || strings.HasPrefix(arg, "-") {
			argFlags[arg] = true
			used[i] = true
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				argMap[arg] = args[i+1]
				used[i+1] = true
				i++ // skip the value
			}
		}
	}

	// Extract positional arguments (non-flag args)
	var positionals []string
	for i, arg := range args {
		if !used[i] {
			positionals = append(positionals, arg)
		}
	}

	// Handle --help
	if common.MetaArgEnabled("Help", target) {
		if argFlags["-h"] || argFlags["--help"] {
			help, err := display.BuildHelp(target, argFlags["--help"])
			if err != nil {
				return err
			}
			fmt.Println(help)
			osExit(0)
		}
	}

	// Handle --version
	if common.MetaArgEnabled("Version", target) {
		if argFlags["--version"] {
			version, err := display.BuildVersion(target)
			if err != nil {
				return err
			}
			fmt.Println(version)
			osExit(0)
		}
	}

	v := reflect.ValueOf(target).Elem()
	t := v.Type()
	positionalIndex := 0

	for i := range t.NumField() {
		field := t.Field(i)

		// Skip meta fields like Clifford, Version, Help
		if field.Type.Name() == "Clifford" || field.Type.Name() == "Version" || field.Type.Name() == "Help" {
			continue
		}
		if field.Type.Kind() != reflect.Struct {
			continue
		}

		subVal := v.Field(i)
		subType := field.Type
		tags := common.GetTagsFromEmbedded(subType, field.Name)

		var value string
		found := false

		longFlag := "--" + tags["long"]
		shortFlag := "-" + tags["short"]

		// Check long flag
		if tags["long"] != "" {
			if val, ok := argMap[longFlag]; ok {
				value = val
				found = true
			}
		}
		// Check short flag
		if !found && tags["short"] != "" {
			if val, ok := argMap[shortFlag]; ok {
				value = val
				found = true
			}
		}
		// Handle boolean flags (without values)
		if !found && tags["long"] != "" {
			if _, ok := argFlags[longFlag]; ok {
				value = "true"
				found = true
			}
		}
		if !found && tags["short"] != "" {
			if _, ok := argFlags[shortFlag]; ok {
				value = "true"
				found = true
			}
		}

		// Handle positional arguments (no short or long tag)
		if !found && tags["short"] == "" && tags["long"] == "" {
			if positionalIndex < len(positionals) {
				value = positionals[positionalIndex]
				positionalIndex++
				found = true
			}
		}

		// Required check
		if !found && tags["required"] == "true" {
			return fmt.Errorf("missing required argument: %s", field.Name)
		}

		// Set the value to the `Value` field
		if found {
			valField := subVal.FieldByName("Value")
			if !valField.IsValid() || !valField.CanSet() {
				continue
			}

			switch valField.Kind() {
			case reflect.String:
				valField.SetString(value)
			case reflect.Int:
				if intVal, err := strconv.Atoi(value); err == nil {
					valField.SetInt(int64(intVal))
				}
			case reflect.Float64:
				if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
					valField.SetFloat(floatVal)
				}
			case reflect.Bool:
				if boolVal, err := strconv.ParseBool(value); err == nil {
					valField.SetBool(boolVal)
				}
			default:
				return fmt.Errorf("unsupported type for field %s: %s", field.Name, valField.Kind())
			}
		}
	}

	return nil
}
