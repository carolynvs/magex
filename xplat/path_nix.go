// +build !windows

package xplat

import "path"

// FilePathJoin elements accounting for the operating system _and_ shell.
// For example, on Windows with MingW the path is formatted in the linux-style.
func FilePathJoin(elem ...string) string {
	return path.Join(elem...)
}
