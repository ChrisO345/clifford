package clifford

import "fmt"

func BuildVersion(target any) (string, error) {
	if !isStructPtr(target) {
		return "", fmt.Errorf("invalid type: must pass pointer to struct")
	}

	t := getStructType(target)
	for i := range t.NumField() {
		field := t.Field(i)
		if field.Name == "Version" {
			if val := field.Tag.Get("version"); val != "" {
				return val, nil
			}
		}
	}

	// FIXME: If no version tag is found, try to infer a version
	return "No version specified", nil
}
