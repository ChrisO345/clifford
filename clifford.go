package clifford

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

var osExit = os.Exit // For testing purposes

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

	// If running inside a testcase, we only care about args after the first "--"
	if indexOf(args, "--") >= 0 {
		args = args[indexOf(args, "--")+1:]
	}

	// Help flag check
	if metaEnabled("Help", target) || metaEnabled("Version", target) {
		for _, arg := range args {
			switch arg {
			case "-h", "--help":
				help, err := BuildHelp(target)
				if err != nil {
					return err
				}
				fmt.Println(help)
				osExit(0)
			case "--version":
				version, err := BuildVersion(target)
				if err != nil {
					return err
				}
				fmt.Println(version)
				osExit(0)
			}
		}
	}

	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("must pass pointer to struct")
	}
	v = v.Elem()
	t := v.Type()

	argIndex := 0 // for positional args

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip meta fields
		if field.Type.Name() == "Clifford" {
			continue
		}

		if field.Type.Kind() != reflect.Struct {
			continue
		}

		subVal := v.Field(i)
		subType := field.Type

		tags := getTagsFromEmbedded(subType, field.Name)

		// CLI Flag Parsing
		var value string
		found := false

		for idx := range args {
			arg := args[idx]

			// Handle --long
			if tags["long"] != "" && strings.HasPrefix(arg, "--") {
				name := strings.TrimPrefix(arg, "--")
				if name == tags["long"] && idx+1 < len(args) {
					value = args[idx+1]
					found = true
					break
				}
			}

			// Handle -s
			if tags["short"] != "" && strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") {
				name := strings.TrimPrefix(arg, "-")
				if name == tags["short"] && idx+1 < len(args) {
					value = args[idx+1]
					found = true
					break
				}
			}
		}

		// Positional arg
		if !found && tags["short"] == "" && tags["long"] == "" {
			if !(argIndex < len(args) && strings.HasPrefix(args[argIndex], "-")) && argIndex < len(args) {
				value = args[argIndex]
				argIndex++
				found = true
			}
		}

		// Check required
		if !found && tags["required"] == "true" {
			// FIXME: Replace this with a custom error handling system as this is an external user error
			return fmt.Errorf("missing required argument: %s", field.Name)
		}

		// Set to Value field
		if found {
			for j := 0; j < subType.NumField(); j++ {
				subField := subType.Field(j)
				if subField.Name == "Value" {
					valField := subVal.FieldByName("Value")
					if valField.IsValid() && valField.CanSet() {
						valField.SetString(value) // supports string only for now
					}
				}
			}
		}
	}

	return nil
}
