// +build windows

package xplat

import (
	"os"
)

// PathSeparator determines the correct path separator based on the operating
// system _and_ shell.
//
// For example, on Windows with Git Bash (mingw) the path is formatted
// linux-style.
func PathSeparator() rune {
	if IsMSys2() {
		return '/'
	}

	return os.PathSeparator
}

// PathListSeparator determines the PATH separator that is safe to use on any OS,
// including when run through Git Bash (mingw).
func PathListSeparator() rune {
	if IsMSys2() {
		return ':'
	}

	return os.PathListSeparator
}

// FileExt returns the default file extension based on the operating system.
func FileExt() string {
	return ".exe"
}
