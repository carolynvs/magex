// +build !windows

package pkg

import (
	"go/build"
	"os"
	"path"
)

func JoinPath(elem ...string) string {
	return path.Join(elem...)
}

func GOPATH() string {
	return build.Default.GOPATH
}

func listPathSeparator(path string) string {
	return string(os.PathListSeparator)
}
