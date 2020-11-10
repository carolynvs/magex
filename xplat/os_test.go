package xplat

import (
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInPath(t *testing.T) {
	sep := string(PathListSeparator())

	mkPath := func(segments ...string) string {
		return strings.Join(segments, sep)
	}

	testcases := []struct {
		name  string
		path  string
		value string
		want  bool
	}{
		{"missing", mkPath("/test/bin"), "/test", false},
		{"incorrect case", mkPath("/Test"), "/test", false},
		{"exact", mkPath("/test"), "/test", true},
		{"trailing", mkPath("/test"), "/test/", true},
		{"trailing in path", mkPath("/test/", "/tmp"), "/test", true},
		{"embedded", mkPath("/bin", "/test", "/tmp"), "/test", true},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

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

	if runtime.GOOS == "windows" && !IsMingw() {
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
