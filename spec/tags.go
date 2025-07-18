package spec

import (
	"fmt"
	"reflect"
)

func Parse(target any) error {
	// Get cli args
	// args := os.Args[1:]

	t := reflect.TypeOf(target)

	for i := range t.NumField() {
		field := t.Field(i)
		fmt.Printf("Field: %v\n", field)
		for j := range field.Type.NumField() {
			subField := field.Type.Field(j)
			fmt.Printf("  SubField: %v\n", subField)

			// Check for tags
			if tag, ok := subField.Tag.Lookup("short"); ok {
				fmt.Printf("    Short Tag: %s\n", tag)
			}
			if tag, ok := subField.Tag.Lookup("desc"); ok {
				fmt.Printf("    Description: %s\n", tag)
			}
			if tag, ok := subField.Tag.Lookup("required"); ok {
				fmt.Printf("    Required: %s\n", tag)
			}
		}
	}

	return nil
}
