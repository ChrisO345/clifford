package spec

import (
	"fmt"
	"testing"

	"github.com/chriso345/gore/assert"
)

func TestHelp(t *testing.T) {
	target := struct {
		Clifford `name:"clifford"`

		Name struct {
			Value    string
			ShortTag `short:"n"`
		}
	}{}

	a, err := BuildHelp(target)
	assert.Nil(t, err)
	fmt.Println(a)
}
