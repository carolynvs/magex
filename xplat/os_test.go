package xplat

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInPath(t *testing.T) {
	pathSep := string(PathSeparator())
	listSep := string(PathListSeparator())

	mkPath := func(pathParts ...[]string) string {
		paths := make([]string, 0, len(pathParts))
		for _, segments := range pathParts {
			p := filepath.Join(segments...)
			paths = append(paths, p)
		}
		return strings.Join(paths, listSep)
	}

	testcases := []struct {
		name  string
		path  string
		value string
		want  bool
	}{
		{"missing", mkPath([]string{"test", "bin"}), "test", false},
		{"incorrect case", "Test", "test", false},
		{"exact", "test", "test", true},
		{"trailing", "test", "test" + pathSep, true},
		{"trailing in path", mkPath([]string{"test" + pathSep}, []string{"tmp"}), "test", true},
		{"embedded", mkPath([]string{"bin"}, []string{"test"}, []string{"tmp"}), "test", true},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			origPath := os.Getenv("PATH")
			os.Setenv("PATH", tc.path)
			defer os.Setenv("PATH", origPath)

			got := InPath(tc.value)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestPrependPath(t *testing.T) {
	origPath := os.Getenv("PATH")
	defer os.Setenv("PATH", origPath)

	if runtime.GOOS == "windows" && !IsMSys2() {
		os.Setenv("PATH", `C:\Temp`)
		PrependPath(`C:\test`)
		gotPath := os.Getenv("PATH")
		assert.Equal(t, `C:\test;C:\Temp`, gotPath)
	} else {
		os.Setenv("PATH", "/tmp")
		PrependPath("/test")
		gotPath := os.Getenv("PATH")
		assert.Equal(t, "/test:/tmp", gotPath)
	}
}
