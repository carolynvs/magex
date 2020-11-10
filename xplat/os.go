package xplat

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// InPath determines if the path is in the PATH environment variable.
func InPath(value string) bool {
	pathSep := string(PathSeparator())
	pathListSep := string(PathListSeparator())
	value = strings.TrimRight(value, pathSep)

	path := os.Getenv("PATH")
	paths := strings.Split(path, pathListSep)
	for _, p := range paths {
		p = strings.TrimRight(p, pathSep)
		if p == value {
			return true
		}
	}

	return false
}

// EnsureInPath adds the specified path to the beginning of the PATH environment
// variable when it is already in PATH.
func EnsureInPath(value string) {
	if !InPath(value) {
		PrependPath(value)
	}
}

// PrependPath adds the specified path to the beginning of the PATH environment
// variable.
func PrependPath(value string) {
	path := os.Getenv("PATH")
	sep := string(PathListSeparator())

	path = fmt.Sprintf("%s%s%s", value, sep, path)
	os.Setenv("PATH", path)
	log.Printf("Added %s to $PATH\n", value)
}
