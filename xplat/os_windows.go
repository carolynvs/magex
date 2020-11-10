// +build windows

package xplat

import (
	"os"
	"strings"
)

// PathSeparator determines the correct path separator based on the operating
// system _and_ shell.
//
// For example, on Windows with Git Bash (mingw) the path is formatted
// linux-style.
func PathSeparator() rune {
	if IsMingw() {
		return '/'
	}

	return os.PathSeparator
}

// PathListSeparator determines the PATH separator that is safe to use on any OS,
// including when run through Git Bash (mingw).
func PathListSeparator() rune {
	if IsMingw() {
		return ':'
	}

	return os.PathListSeparator
}

// IsMingw determines if the process is executing on Git Bash (MingW).
func IsMingw() bool {
	path := os.Getenv("PATH")
	return strings.Contains(path, "/mingw")
}
