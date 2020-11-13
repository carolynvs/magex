package xplat

import (
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetermineShell(t *testing.T) {
	os.Setenv("HOME", "/home/me")
	defer os.Unsetenv("HOME")

	testcases := []struct {
		name      string
		env       map[string]string
		goos      string // skip if the current goos doesn't match
		shell     string // skip if the current shell doesn't match
		wantShell string
	}{
		{name: "shell set - linux", env: map[string]string{"SHELL": "bash"},
			goos: "linux", shell: "bash", wantShell: "bash"},
		{name: "shell set - macos", env: map[string]string{"SHELL": "bash"},
			goos: "darwin", shell: "bash", wantShell: "bash"},
		{name: "git bash", env: nil,
			goos: "windows", shell: "mingw64", wantShell: "mingw64"},
		{name: "powershell", env: nil,
			goos: "windows", shell: "powershell", wantShell: "powershell"},
		{name: "default - windows", env: nil,
			goos: "windows", shell: "cmd", wantShell: "cmd"},
		{name: "default - linux", env: map[string]string{"SHELL": ""},
			goos: "linux", shell: "bash", wantShell: "posix"},
		{name: "default - macos", env: map[string]string{"SHELL": ""},
			goos: "darwin", shell: "bash", wantShell: "posix"},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.goos != runtime.GOOS {
				t.Skip("Skipping because GOOS doesn't match")
			}

			testShell := os.Getenv("TEST_SHELL")
			if tc.shell != testShell {
				t.Skipf("Skipping because TEST_SHELL isn't %s", tc.shell)
			}

			for k, v := range tc.env {
				os.Setenv(k, v)
				//goland:noinspection GoDeferInLoop
				defer func() {
					os.Unsetenv(k)
				}()
			}

			gotShell := DetectShell()
			assert.Equal(t, tc.wantShell, gotShell)
		})

	}
}
