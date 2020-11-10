// +build !windows

package pkg

import (
	"go/build"
	"os"
	"path"
)

// JoinPath elements accounting for the operating system _and_ shell.
// For example, on Windows with MingW the path is formatted in the linux-style.
func JoinPath(elem ...string) string {
	return path.Join(elem...)
}

// GOPATH returns the current GOPATH that is safe to use on any OS, including
// when run through Git Bash (mingw).
func GOPATH() string {
	return build.Default.GOPATH
}

// listPathSeparator determines the PATH separator that is safe to use on any OS,
// including when run through Git Bash (mingw).
func listPathSeparator(path string) string {
	return string(os.PathListSeparator)
}
