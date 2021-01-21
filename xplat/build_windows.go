// +build windows

package xplat

import (
	"go/build"
	"os"
	"path/filepath"
	"strings"
)

// GOPATH returns the current GOPATH that is safe to use on any OS, including
// when run through Git Bash (mingw).
func GOPATH() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	if IsMSys2() {
		// Remove volume separator
		gopath = strings.ReplaceAll(gopath, ":", "")
		// Convert to unix path separator
		gopath = filepath.ToSlash(gopath)
		return gopath
	}

	return gopath
}
