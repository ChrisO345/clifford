package common

import (
	"reflect"
	"strings"
)

// GetTagsFromEmbedded retrieves tags from embedded structs in the target struct.
func GetTagsFromEmbedded(t reflect.Type, fieldName string) map[string]string {
	tags := make(map[string]string)

	for i := range t.NumField() {
		field := t.Field(i)
		if !field.Anonymous {
			continue
		}

		switch field.Type.Name() {
		case "ShortTag":
			tags["short"] = strings.ToLower(string(fieldName[0]))
		case "LongTag":
			tags["long"] = strings.ToLower(fieldName)
		case "Required":
			tags["required"] = "true"
		case "Desc":
			if val := field.Tag.Get("desc"); val != "" {
				tags["desc"] = val
			}
		default:
			for _, key := range []string{"short", "long", "desc", "required"} {
				if val := field.Tag.Get(key); val != "" {
					tags[key] = val
				}
			}
		}
	}

	return tags
}

// ArgsIndexOf returns the index of the first occurrence of s in args, or -1 if not found.
func ArgsIndexOf(args []string, s string) int {
	for i, arg := range args {
		if arg == s {
			return i
		}
	}
	return -1
}

// IsStructPtr checks if the provided value is a pointer to a struct.
func IsStructPtr(v any) bool {
	t := reflect.TypeOf(v)
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}

// GetStructType returns the reflect.Type of the underlying struct pointer.
func GetStructType(v any) reflect.Type {
	return reflect.TypeOf(v).Elem()
}

// MetaArgEnabled returns true if the root struct has a `Clifford` field with tag or name matching s
// or if the field name itself matches s.
func MetaArgEnabled(s string, target any) bool {
	t := GetStructType(target)

	for i := range t.NumField() {
		field := t.Field(i)

		// Must be a Clifford field at root
		if field.Type.Name() != "Clifford" {
			if field.Type.Name() == s {
				return true
			}
			continue
		}

		// Check if tag contains the key s
		if tagValue := field.Tag.Get(s); tagValue != "" {
			return true
		}

		// Check if the field name itself matches s
		if field.Name == s {
			return true
		}
	}

	return false
}
