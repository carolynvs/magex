package pkg

import (
	"go/build"
	"os"
)

// GOPATH returns the current GOPATH.
func GOPATH() string {
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		return gopath
	}
	return build.Default.GOPATH
}
