package pkg

import (
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/carolynvs/magex/xplat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDownloadToGopathBin(t *testing.T) {
	url := "https://storage.googleapis.com/kubernetes-release/release/{{.VERSION}}/bin/{{.GOOS}}/{{.GOARCH}}/kubectl{{.EXT}}"
	err := DownloadToGopathBin(url, "kubectl", "v1.19.0")
	require.NoError(t, err)

	dest := filepath.Join(xplat.GOPATH(), "bin/kubectl")
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
	tmp, err := ioutil.TempDir("", "magex")
	require.NoError(t, err)
	defer os.RemoveAll(tmp)

	os.Setenv("GOPATH", tmp)
	defer os.Setenv("GOPATH", build.Default.GOPATH)

	path := os.Getenv("PATH")
	os.Setenv("PATH", strings.ReplaceAll(path, build.Default.GOPATH, ""))
	defer os.Setenv("PATH", path)

	hasCmd, err := IsCommandAvailable("yq", "")
	require.False(t, hasCmd)
	err = EnsurePackage("github.com/mikefarah/yq/v4", "4.4.1", "--version")
	require.NoError(t, err)
}

func TestInstallPackage_MajorVersion(t *testing.T) {
	tmp, err := ioutil.TempDir("", "magex")
	require.NoError(t, err)
	defer os.RemoveAll(tmp)

	os.Setenv("GOPATH", tmp)

	err = InstallPackage("github.com/mikefarah/yq/v4", "4.4.1")
	if err != nil {
		log.Fatal("could not install yq")
	}
}
