// +build !windows

package xplat

import "os"

// PathListSeparator determines the PATH separator that is safe to use on any OS,
// including when run through Git Bash (mingw).
func PathListSeparator() rune {
	return os.PathListSeparator
}

// PathSeparator determines the correct path separator based on the operating
// system _and_ shell.
//
// For example, on Windows with Git Bash (mingw) the path is formatted
// linux-style.
func PathSeparator() rune {
	return os.PathSeparator
}

// FileExt returns the default file extension based on the operating system.
func FileExt() string {
	return ""
}
