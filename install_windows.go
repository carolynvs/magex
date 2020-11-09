// +build windows

package magex

import (
	"go/build"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func listPathSeparator(path string) string {
	if isMingw() {
		return ":"
	}

	return string(os.PathListSeparator)
}

func isMingw() bool {
	path := os.Getenv("PATH")
	return strings.Contains(path, "/mingw")
}

// GOPATH returns the current gopath that is safe to use on any OS, including
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

func JoinPath(elem ...string) string {
	if isMingw() {
		return path.Join(elem...)
	}

	return filepath.Join(elem...)
}
