// +build windows

package xplat

import (
	"path"
	"path/filepath"
)

// FilePathJoin elements accounting for the operating system _and_ shell.
//
// For example, on Windows with Git Bash (mingw) the path is formatted
// linux-style.
func FilePathJoin(elem ...string) string {
	if IsMSys2() {
		return path.Join(elem...)
	}

	return filepath.Join(elem...)
}
