package clifford

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var osExit = os.Exit // Mockable for testing

// Parse parses command-line arguments into the provided target struct.
//
// The target must be a pointer to a struct where each field represents either
// a CLI argument or a group of options. Each sub-struct should contain a `Value`
// field to hold the parsed value, and may be annotated using `clifford` tags or
// helper types like `ShortTag`, `LongTag`, `Required`, and `Desc`.
//
// If the first argument passed to the CLI is `-h` or `--help`, Parse will
// automatically call BuildHelp and exit the program.
//
// Usage:
//
//	target := struct {
//		clifford.Clifford `name:"mytool"`
//
//		Name struct {
//			Value    string
//			clifford.Clifford `short:"n" long:"name" desc:"User name"`
//		}
//
//		Age struct {
//			Value    string
//			clifford.ShortTag // Auto-generates: -a
//			clifford.LongTag  // Auto-generates: --age
//			clifford.Desc     `desc:"Age of the user"`
//		}
//	}{}
//
//	err := clifford.Parse(&target)
//	if err != nil {
//		log.Fatal(err)
//	}
func Parse(target any) error {
	if !isStructPtr(target) {
		return fmt.Errorf("invalid type: must pass pointer to struct")
	}

	args := os.Args[1:]

	// For test mode: drop args before "--"
	if i := indexOf(args, "--"); i >= 0 {
		args = args[i+1:]
	}

	// Pre-process args into a map for fast lookup
	argMap := map[string]string{}
	argFlags := map[string]bool{}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "--") || strings.HasPrefix(arg, "-") {
			argFlags[arg] = true
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				argMap[arg] = args[i+1]
				i++ // skip the value
			}
		}
	}

	// Handle help/version
	if metaEnabled("Help", target) {
		if argFlags["-h"] || argFlags["--help"] {
			help, err := BuildHelp(target)
			if err != nil {
				return err
			}
			fmt.Println(help)
			osExit(0)
		}
	}
	if metaEnabled("Version", target) {
		if argFlags["--version"] {
			version, err := BuildVersion(target)
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

		// Skip internal meta fields
		if field.Type.Name() == "Clifford" {
			continue
		}
		if field.Type.Kind() != reflect.Struct {
			continue
		}

		subVal := v.Field(i)
		subType := field.Type
		tags := getTagsFromEmbedded(subType, field.Name)

		var value string
		found := false

		longFlag := "--" + tags["long"]
		shortFlag := "-" + tags["short"]

		// Check for long or short flags in map
		if tags["long"] != "" {
			if val, ok := argMap[longFlag]; ok {
				value = val
				found = true
			}
		}
		if !found && tags["short"] != "" {
			if val, ok := argMap[shortFlag]; ok {
				value = val
				found = true
			}
		}

		// Handle boolean flags (no value)
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

		// Handle positional arguments
		if !found && tags["short"] == "" && tags["long"] == "" {
			if positionalIndex < len(args) && !strings.HasPrefix(args[positionalIndex], "-") {
				value = args[positionalIndex]
				positionalIndex++
				found = true
			}
		}

		// Required check
		if !found && tags["required"] == "true" {
			// FIXME: Replace this with a custom error handling system as this is an external user error
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
