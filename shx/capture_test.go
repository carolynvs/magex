package shx

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecordStderr(t *testing.T) {
	const msg = "printed to os.Stderr"
	orig := *os.Stderr

	stderr := RecordStderr()

	// Verify that stderr has been redirected
	assert.Equal(t, "/dev/stderr", orig.Name())
	assert.NotEqual(t, *os.Stderr, orig)

	fmt.Fprint(os.Stderr, msg)
	got := stderr.Output()
	assert.Equal(t, msg, got)

	assert.Equal(t, *os.Stderr, orig, "Stdout was not restored")
}

func TestRecordStdout(t *testing.T) {
	const msg = "printed to os.Stdout"
	orig := *os.Stdout

	stderr := RecordStdout()

	// Verify that stdout has been redirected
	assert.Equal(t, "/dev/stdout", orig.Name())
	assert.NotEqual(t, *os.Stdout, orig)

	fmt.Fprint(os.Stdout, msg)
	got := stderr.Output()
	assert.Equal(t, msg, got)

	assert.Equal(t, *os.Stdout, orig, "Stdout was not restored")
}
