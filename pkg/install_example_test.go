package pkg_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/carolynvs/magex/xplat"
	"github.com/stretchr/testify/require"

	"github.com/carolynvs/magex/pkg"
)

func ExampleEnsureMage() {
	// Leave the version parameter blank to only check if it is installed, and
	// if not install the latest version.
	err := pkg.EnsureMage("")
	if err != nil {
		log.Fatal("could not install mage")
	}
}

func ExampleEnsurePackage() {
	// Install packr2@v2.8.0 using the command `packr2 version` to detect if the
	// correct version is installed.
	err := pkg.EnsurePackage("github.com/gobuffalo/packr/v2/packr2", "v2.8.0", "version")
	if err != nil {
		log.Fatal("could not install packr2")
	}
}

func ExampleInstallPackage() {
	// Install packr2@v2.8.0
	err := pkg.InstallPackage("github.com/gobuffalo/packr/v2/packr2", "v2.8.0")
	if err != nil {
		log.Fatal("could not install packr2")
	}
}

func ExampleDownloadToGopathBin() {
	url := "https://storage.googleapis.com/kubernetes-release/release/{{.VERSION}}/bin/{{.GOOS}}/{{.GOARCH}}/kubectl{{.EXT}}"
	err := pkg.DownloadToGopathBin(url, "kubectl", "v1.19.0")
	if err != nil {
		log.Fatal("could not download kubectl")
	}

	// Add GOPATH/bin to PATH if necessary so that we can immediately
	// use the installed tool
	pkg.EnsureGopathBin()
}

func TestDownloadToGopathBin(t *testing.T) {
	url := "https://storage.googleapis.com/kubernetes-release/release/{{.VERSION}}/bin/{{.GOOS}}/{{.GOARCH}}/kubectl{{.EXT}}"
	err := pkg.DownloadToGopathBin(url, "kubectl", "v1.19.0")
	require.NoError(t, err)

	dest := filepath.Join(xplat.GOPATH(), "bin/kubectl")
	_, err = os.Stat(dest)
	require.NoError(t, err)

	os.Remove(dest)
}
