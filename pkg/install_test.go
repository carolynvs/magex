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

	url := "https://dl.k8s.io/release/{{.VERSION}}/bin/{{.GOOS}}/{{.GOARCH}}/kubectl{{.EXT}}"
	err = DownloadToGopathBin(url, "kubectl", "v1.23.0")
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

func TestEnsurePackage_FreshInstall(t *testing.T) {
	os.Setenv(mg.VerboseEnv, "true")
	defer os.Unsetenv(mg.VerboseEnv)

	testcases := []struct {
		name              string
		versionConstraint string
		defaultVersion    string
		wantVersion       string
	}{
		{name: "with prefix", versionConstraint: "v2.0.x", defaultVersion: "v2.0.2", wantVersion: "v2.0.2"},
		{name: "without prefix", versionConstraint: "2.0.x", defaultVersion: "2.0.2", wantVersion: "v2.0.2"},
		{name: "no version", versionConstraint: "", wantVersion: "v2.0.3"},
		{name: "latest version", versionConstraint: "latest", wantVersion: "v2.0.3"},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err, cleanup := gopath.UseTempGopath()
			require.NoError(t, err, "Failed to set up a temporary GOPATH")
			defer cleanup()

			// Verify it's not currently installed
			hasCmd, err := IsCommandAvailable("testpkg", "", "")
			require.False(t, hasCmd)

			err = EnsurePackage("github.com/carolynvs/testpkg/v2", tc.defaultVersion, "--version", tc.versionConstraint)
			require.NoError(t, err)

			installedVersion, err := GetCommandVersion("testpkg", "")
			require.Equal(t, tc.wantVersion, installedVersion, "incorrect version was resolved")
		})
	}
}

func TestEnsurePackage_Upgrade(t *testing.T) {
	os.Setenv(mg.VerboseEnv, "true")
	defer os.Unsetenv(mg.VerboseEnv)

	testcases := []struct {
		name              string
		versionConstraint string
		defaultVersion    string
		wantVersion       string
	}{
		{name: "constraint allows existing version", versionConstraint: "v2.0.x", defaultVersion: "v2.0.2", wantVersion: "v2.0.2"},
		{name: "constraint requires higher version", versionConstraint: "^2.0.3", defaultVersion: "2.0.3", wantVersion: "v2.0.3"},
		{name: "no version constraint", versionConstraint: "", defaultVersion: "", wantVersion: "v2.0.2"},
		{name: "latest version", defaultVersion: "latest", wantVersion: "v2.0.2"},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err, cleanup := gopath.UseTempGopath()
			require.NoError(t, err, "Failed to set up a temporary GOPATH")
			defer cleanup()

			// Install an old version
			err = InstallPackage("github.com/carolynvs/testpkg/v2", "v2.0.2")
			require.NoError(t, err)

			// Ensure it's installed with a higher default version
			err = EnsurePackage("github.com/carolynvs/testpkg/v2", tc.defaultVersion, "--version", tc.versionConstraint)
			require.NoError(t, err)

			installedVersion, err := GetCommandVersion("testpkg", "")
			require.Equal(t, tc.wantVersion, installedVersion, "incorrect version was resolved")
		})
	}
}
