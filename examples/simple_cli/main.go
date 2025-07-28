package main

import (
	"github.com/chriso345/clifford"
)

func main() {
	target := struct {
		clifford.Clifford `name:"example_cli"` // This is the name of the cli command
		clifford.Version  `version:"1.0.0"`
		clifford.Help

		Name struct {
			Value string
		}
	}{}

	err := clifford.Parse(&target)
	if err != nil {
		panic(err)
	}

	if target.Name.Value != "" {
		println("Hello, " + target.Name.Value + "!")
	} else {
		println("Hello, World!")
	}
}
