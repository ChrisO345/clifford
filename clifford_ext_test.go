package clifford_test

import (
	"os"
	"testing"

	"github.com/chriso345/clifford"
	"github.com/chriso345/gore/assert"
)

func TestParse_ShortAndLongFlags(t *testing.T) {
	// Save original args and restore later
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmd", "--name", "Alice", "-a", "30"}

	cli := struct {
		clifford.Clifford `name:"mytool"`

		Name struct {
			Value             string
			clifford.Clifford `short:"n" long:"name" desc:"User name"`
		}

		Age struct {
			Value string
			clifford.ShortTag
			clifford.LongTag
			clifford.Desc `desc:"Age of user"`
		}
	}{}

	err := clifford.Parse(&cli)
	assert.Nil(t, err)
	assert.Equal(t, "Alice", cli.Name.Value)
	assert.Equal(t, "30", cli.Age.Value)
}

func TestParse_PositionalArgs(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmd", "Alice", "30"}

	cli := struct {
		clifford.Clifford `name:"mytool"`

		Name struct {
			Value string
			clifford.Required
		}
		Age struct {
			Value string
		}
	}{}

	err := clifford.Parse(&cli)
	assert.Nil(t, err)
	assert.Equal(t, "Alice", cli.Name.Value)
	assert.Equal(t, "30", cli.Age.Value)
}

func TestParse_DebugPositionalArgs(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"junk", "junk", "--", "Alice", "30"}

	cli := struct {
		clifford.Clifford `name:"mytool"`

		Name struct {
			Value string
			clifford.Required
		}
		Age struct {
			Value string
		}
	}{}

	err := clifford.Parse(&cli)
	assert.Nil(t, err)
	assert.Equal(t, cli.Name.Value, "Alice")
	assert.Equal(t, cli.Age.Value, "30")
}

func TestParse_MissingRequired(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmd", "--age", "30"}

	cli := struct {
		clifford.Clifford `name:"mytool"`

		Name struct {
			Value string
			clifford.Required
		}
		Age struct {
			Value string
			clifford.LongTag
		}
	}{}

	err := clifford.Parse(&cli)
	assert.NotNil(t, err)
	assert.StringContains(t, err.Error(), "missing required argument: Name")
}
