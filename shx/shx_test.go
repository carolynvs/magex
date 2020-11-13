package shx_test

import (
	"testing"

	"github.com/carolynvs/magex/shx"
	"github.com/stretchr/testify/assert"
)

func TestCollapseArgs(t *testing.T) {
	verbose := "" // -v was omitted
	args := shx.CollapseArgs("go", "test", verbose, "./...")
	assert.Equal(t, []string{"go", "test", "./..."}, args)
}
