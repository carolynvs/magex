// +build windows

package xplat

import (
	"go/build"
	"path/filepath"
	"strings"
)

// GOPATH returns the current GOPATH that is safe to use on any OS, including
// when run through Git Bash (mingw).
func GOPATH() string {
	if IsMSys2() {
		gopath := build.Default.GOPATH
		// Remove volume separator
		gopath = strings.ReplaceAll(gopath, ":", "")
		// Convert to unix path separator
		gopath = filepath.ToSlash(gopath)
		return gopath
	}

	return build.Default.GOPATH
}
