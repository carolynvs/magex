// +build !windows

package xplat

import "go/build"

// GOPATH returns the current GOPATH that is safe to use on any OS, including
// when run through Git Bash (mingw).
func GOPATH() string {
	return build.Default.GOPATH
}
