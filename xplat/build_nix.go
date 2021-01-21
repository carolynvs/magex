// +build !windows

package xplat

import (
	"go/build"
	"os"
)

// GOPATH returns the current GOPATH that is safe to use on any OS, including
// when run through Git Bash (mingw).
func GOPATH() string {
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		return gopath
	}
	return build.Default.GOPATH
}
