package clifford

import (
	"runtime/debug"
)

const root = "github.com/chriso345/clifford"

// ModuelVersion returns the version of the clifford module, if available.
func ModuleVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Path == root && info.Main.Version != "" {
			return info.Main.Version
		}
		for _, dep := range info.Deps {
			if dep.Path == root {
				return dep.Version
			}
		}
	}

	return "unknown"
}
