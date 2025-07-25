package clifford

import (
	"os"
	"testing"

	"github.com/chriso345/gore/assert"
)

func TestParse_HelpFlag(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmd", "--help"}

	cli := struct {
		Clifford `name:"mytool"`
		Help
		Name struct {
			Value string
			Required
		}
	}{}

	// Temporarily override os.Exit
	calledExit := false
	osExit = func(code int) {
		calledExit = true
		panic("os.Exit called")
	}
	defer func() { osExit = os.Exit }()

	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, true, calledExit)
		}
	}()

	_ = Parse(&cli)
	t.Errorf("should have exited before this line")
}

func TestParse_VersionFlag(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmd", "--version"}

	cli := struct {
		Clifford `name:"mytool"`
		Version  `version:"1.2.3"`
	}{}

	calledExit := false
	osExit = func(code int) {
		calledExit = true
		panic("os.Exit called")
	}
	defer func() { osExit = os.Exit }()

	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, true, calledExit)
		}
	}()

	_ = Parse(&cli)
	t.Errorf("should have exited before this line")
}
