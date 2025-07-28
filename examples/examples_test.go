package examples_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestExamplesDoNotPanic(t *testing.T) {
	examplesRoot := "."

	entries, err := os.ReadDir(examplesRoot)
	if err != nil {
		t.Fatalf("Failed to read examples directory: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		examplePath := filepath.Join(examplesRoot, entry.Name())
		mainGo := filepath.Join(examplePath, "main.go")

		if _, err := os.Stat(mainGo); err != nil {
			continue
		}

		cmd := getCmdArgs(examplePath, mainGo)

		t.Run(entry.Name(), func(t *testing.T) {
			cmd := exec.Command("go", cmd...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				t.Errorf("Example %q failed: %v", entry.Name(), err)
			}
		})
	}
}

func getCmdArgs(examplepath string, mainGo string) []string {
	// Reads from a file containing command line arguments found at filepath/args.txt.
	argsFile := filepath.Join(examplepath, "args.txt")
	if _, err := os.Stat(argsFile); os.IsNotExist(err) {
		return nil // No args file, return empty slice
	}

	data, err := os.ReadFile(argsFile)
	if err != nil {
		return nil // Error reading file, return empty slice
	}

	strData := strings.Trim(string(data), "\n")
	if strings.ContainsRune(strData, '\n') {
		fmt.Println("Warning: args.txt contains multiple lines, only the first line will be used.")
		strData = strings.Split(strData, "\n")[0] // Use only the first line
	}
	args := strings.Fields(strData)

	baseArgs := []string{"run", mainGo, "--"}
	return append(baseArgs, args...)
}
