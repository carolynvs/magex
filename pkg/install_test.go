package pkg

import (
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/carolynvs/magex/xplat"
	"github.com/magefile/mage/mg"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDownloadToGopathBin(t *testing.T) {
	url := "https://storage.googleapis.com/kubernetes-release/release/{{.VERSION}}/bin/{{.GOOS}}/{{.GOARCH}}/kubectl{{.EXT}}"
	err := DownloadToGopathBin(url, "kubectl", "v1.19.0")
	require.NoError(t, err)

	dest := filepath.Join(GOPATH(), "bin/kubectl"+xplat.FileExt())
	_, err = os.Stat(dest)
	require.NoError(t, err)

	os.Remove(dest)
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

func TestEnsurePackage_MajorVersion(t *testing.T) {
	os.Setenv(mg.VerboseEnv, "true")
	err, cleanup := UseTempGopath(t)
	defer cleanup()
	require.NoError(t, err, "Failed to set up a temporary GOPATH")

	hasCmd, err := IsCommandAvailable("testpkg", "")
	require.False(t, hasCmd)
	err = EnsurePackage("github.com/carolynvs/testpkg/v2", "v2.0.1", "--version")
	require.NoError(t, err)
}

func UseTempGopath(t *testing.T) (error, func()) {
	oldpath := os.Getenv("PATH")
	tmp, err := ioutil.TempDir("", "magex")
	if err != nil {
		return errors.Wrap(err, "Failed to create a temp directory"), func() {}
	}

	cleanup := func() {
		os.RemoveAll(tmp)
		defer os.Setenv("PATH", oldpath)
		defer os.Setenv("GOPATH", build.Default.GOPATH)
	}

	// Remove actual GOPATH/bin from PATH so the test doesn't accidentally pass because the package was installed before the test was run
	gopathBin := filepath.Join(build.Default.GOPATH, "bin")
	os.Setenv("PATH", strings.ReplaceAll(oldpath, gopathBin, ""))

	// Use temp dir for GOPATH
	os.Setenv("GOPATH", tmp)

	return EnsureGopathBin(), cleanup
}
