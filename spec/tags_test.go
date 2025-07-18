package spec

import (
	"testing"

	"github.com/chriso345/gore/assert"
)

func TestParse_Basic(t *testing.T) {
	target := struct {
		Name struct {
			Value    string
			ShortTag `short:"n"`
		}
	}{}

	err := Parse(target)
	assert.Nil(t, err)
}
