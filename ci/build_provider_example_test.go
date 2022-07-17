package ci_test

import (
	"fmt"
	"testing"

	"github.com/carolynvs/magex/ci"
)

func TestExampleDetectBuildProvider(t *testing.T) {
	ExampleDetectBuildProvider()
}

func ExampleDetectBuildProvider() {
	// Figure out if you are on a build provider that is supported
	p, detected := ci.DetectBuildProvider()
	if !detected {
		fmt.Println("no build provider was detected, using a noop implementation")
	}

	// Set the LOG_LEVEL environment variable
	p.SetEnv("LOG_LEVEL", "3")

	// Add the gopath bin directory to the beginning of the PATH environment variable
	p.PrependPath("/go/bin")
}
