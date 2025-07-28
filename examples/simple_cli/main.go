package main

import (
	"fmt"
	"os"

	"github.com/chriso345/clifford"
)

func main() {
	target := struct {
		clifford.Clifford `name:"example_cli"` // This is the name of the cli command
		clifford.Version  `version:"1.0.0"`
		clifford.Help

		Name struct {
			Value string
			clifford.ShortTag
		}

		// Age struct {
		// 	Value string
		// 	clifford.Required
		// }
	}{}

	// os.Args = []string{"cmd", "Alice", "30"}
	fmt.Println(os.Args)

	err := clifford.Parse(&target)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("Parsed values: Name=%s, Age=%s\n", target.Name.Value, target.Age.Value)

	if target.Name.Value != "" {
		println("Hello, " + target.Name.Value + "!")
	} else {
		println("Hello, World!")
	}
}
