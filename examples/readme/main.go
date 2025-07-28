package main

import (
	"fmt"
	"log"

	"github.com/chriso345/clifford"
)

func main() {
	target := struct {
		clifford.Clifford `name:"mytool"`   // Set the name of the CLI tool
		clifford.Version  `version:"1.2.3"` // Enable automatic version flag
		clifford.Help                       // Enable automatic help flags

		Name struct {
			Value             string
			clifford.Clifford `short:"n" long:"name" desc:"User name"`
		}
		Age struct {
			Value             string
			clifford.ShortTag // auto generates -a
			clifford.LongTag  // auto generates --age
			clifford.Desc     `desc:"Age of the user"`
		}
	}{}

	err := clifford.Parse(&target)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Name: %s\n", target.Name.Value)
	fmt.Printf("Age: %s\n", target.Age.Value)
}
