package shx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandBuilder_Command(t *testing.T) {
	b := CommandBuilder{
		StopOnError: true,
		Env:         []string{"a=1"},
		Dir:         "tmp",
	}

	cmd := b.Command("go", "build")
	assert.True(t, cmd.StopOnError, "incorrect StopOnError")
	assert.Contains(t, cmd.Cmd.Env, "a=1", "incorrect Env")
	assert.Equal(t, "tmp", cmd.Cmd.Dir, "incorrect Dir")
}
