package clifford

import "fmt"

func BuildVersion(target any) (string, error) {
	if !isStructPtr(target) {
		return "", fmt.Errorf("invalid type: must pass pointer to struct")
	}

	t := getStructType(target)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Name == "Version" {
			if val := field.Tag.Get("version"); val != "" {
				return val, nil
			}
		}
	}

	return "No version specified", nil
}
