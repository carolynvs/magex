// +build windows

package pkg

import (
	"go/build"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// listPathSeparator determines the PATH separator that is safe to use on any OS,
// including when run through Git Bash (mingw).
func listPathSeparator(path string) string {
	if isMingw() {
		return ":"
	}

	return string(os.PathListSeparator)
}

// isMingw determines if the current process is running within Git Bash (mingw).
func isMingw() bool {
	path := os.Getenv("PATH")
	return strings.Contains(path, "/mingw")
}

// GOPATH returns the current GOPATH that is safe to use on any OS, including
// when run through Git Bash (mingw).
func GOPATH() string {
	if isMingw() {
		gopath := build.Default.GOPATH
		// Remove volume separator
		gopath = strings.ReplaceAll(gopath, ":", "")
		// Convert to unix path separator
		gopath = filepath.ToSlash(gopath)
		return gopath
	}

	return build.Default.GOPATH
}

// JoinPath elements accounting for the operating system _and_ shell.
// For example, on Windows with MingW the path is formatted in the linux-style.
func JoinPath(elem ...string) string {
	if isMingw() {
		return path.Join(elem...)
	}

	return filepath.Join(elem...)
}
