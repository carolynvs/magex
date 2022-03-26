package pkg

import (
	"os"
	"os/exec"
	"testing"

	"github.com/carolynvs/magex/pkg/gopath"
	"github.com/carolynvs/magex/xplat"
	"github.com/magefile/mage/mg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDownloadToGopathBin(t *testing.T) {
	err, cleanup := gopath.UseTempGopath()
	require.NoError(t, err, "Failed to set up a temporary GOPATH")
	defer cleanup()

	url := "https://storage.googleapis.com/kubernetes-release/release/{{.VERSION}}/bin/{{.GOOS}}/{{.GOARCH}}/kubectl{{.EXT}}"
	err = DownloadToGopathBin(url, "kubectl", "v1.19.0")
	require.NoError(t, err)

	_, err = exec.LookPath("kubectl" + xplat.FileExt())
	require.NoError(t, err)
}

func TestGetCommandName(t *testing.T) {
	t.Run("v suffix without version", func(t *testing.T) {
		cmd := getCommandName("github.com/foo/verynotsemver")
		assert.Equal(t, "verynotsemver", cmd)
	})

	t.Run("semver suffix", func(t *testing.T) {
		cmd := getCommandName("github.com/foo/bar/v2")
		assert.Equal(t, "bar", cmd)
	})

	t.Run("command suffix", func(t *testing.T) {
		cmd := getCommandName("github.com/foo/bar/cmd/baz")
		assert.Equal(t, "baz", cmd)
	})
}

func TestEnsurePackage(t *testing.T) {
	os.Setenv(mg.VerboseEnv, "true")
	defer os.Unsetenv(mg.VerboseEnv)

	testcases := []struct {
		name    string
		version string
	}{
		{name: "with prefix", version: "v2.0.1"},
		{name: "without prefix", version: "2.0.1"},
		{name: "no version", version: ""},
		{name: "latest version", version: "latest"},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err, cleanup := gopath.UseTempGopath()
			require.NoError(t, err, "Failed to set up a temporary GOPATH")
			defer cleanup()

			hasCmd, err := IsCommandAvailable("testpkg", "")
			require.False(t, hasCmd)
			err = EnsurePackage("github.com/carolynvs/testpkg/v2", tc.version, "--version")
			require.NoError(t, err)
		})
	}
}
