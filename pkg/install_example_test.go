package pkg_test

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/carolynvs/magex/pkg"
	"github.com/carolynvs/magex/pkg/gopath"
)

func TestExampleEnsureMage(t *testing.T) {
	ExampleEnsureMage()
}

func ExampleEnsureMage() {
	// Leave the version parameter blank to only check if it is installed, and
	// if not install the latest version.
	err := pkg.EnsureMage("")
	if err != nil {
		log.Fatal("could not install mage")
	}
}

func TestExampleEnsurePackage(t *testing.T) {
	ExampleEnsurePackage()
}

func ExampleEnsurePackage() {
	// Install packr2@v2.8.3 using the command `packr2 version` to detect if the
	// correct version is installed.
	err := pkg.EnsurePackage("github.com/gobuffalo/packr/v2/packr2", "v2.8.3", "version")
	if err != nil {
		log.Fatal("could not install packr2")
	}
}

func TestExampleEnsurePackageWith_LatestVersion(t *testing.T) {
	ExampleEnsurePackageWith_LatestVersion()
}

func ExampleEnsurePackageWith_LatestVersion() {
	// Install packr2@latest into bin/ using the command `packr2 version` to detect if the
	// correct version is installed.
	err := pkg.EnsurePackageWith(pkg.EnsurePackageOptions{
		Name:           "github.com/gobuffalo/packr/v2/packr2",
		VersionCommand: "version",
		Destination:    "bin",
	})
	if err != nil {
		log.Fatal("could not install packr2")
	}
}

func TestExampleEnsurePackageWith_DefaultVersion(t *testing.T) {
	ExampleEnsurePackageWith_DefaultVersion()
}

func ExampleEnsurePackageWith_DefaultVersion() {
	// Install packr2@v2.8.3 into bin/ using the command `packr2 version` to detect if the
	// correct version is installed.
	err := pkg.EnsurePackageWith(pkg.EnsurePackageOptions{
		Name:           "github.com/gobuffalo/packr/v2/packr2",
		DefaultVersion: "v2.8.3",
		VersionCommand: "version",
		Destination:    "bin",
	})
	if err != nil {
		log.Fatal("could not install packr2")
	}
}

func TestExampleEnsurePackage_VersionConstraint(t *testing.T) {
	ExampleEnsurePackage_VersionConstraint()
}

func ExampleEnsurePackage_VersionConstraint() {
	// Install packr2@v2.8.3 using the command `packr2 version` to detect if
	// any v2 version is installed
	err := pkg.EnsurePackage("github.com/gobuffalo/packr/v2/packr2", "v2.8.3", "version", "2.x")
	if err != nil {
		log.Fatal("could not install packr2")
	}
}

func TestExampleEnsurePackageWith_VersionConstraint(t *testing.T) {
	ExampleEnsurePackageWith_VersionConstraint()
}

func ExampleEnsurePackageWith_VersionConstraint() {
	// Install packr2@v2.8.3 into bin/ using the command `packr2 version` to detect if
	// any v2 version is installed
	err := pkg.EnsurePackageWith(pkg.EnsurePackageOptions{
		Name:           "github.com/gobuffalo/packr/v2/packr2",
		DefaultVersion: "v2.8.3",
		VersionCommand: "version",
		Destination:    "bin",
		AllowedVersion: "2.x",
	})
	if err != nil {
		log.Fatal("could not install packr2")
	}
}

func TestExampleInstallPackage(t *testing.T) {
	ExampleInstallPackage()
}

func ExampleInstallPackage() {
	// Install packr2@v2.8.3
	err := pkg.InstallPackage("github.com/gobuffalo/packr/v2/packr2", "v2.8.3")
	if err != nil {
		log.Fatal("could not install packr2")
	}
}

func TestExampleInstallPackageWith(t *testing.T) {
	tmpBin, err := os.MkdirTemp("", "magex")
	require.NoError(t, err)
	defer os.RemoveAll(tmpBin)

	ExampleInstallPackageWith()
}

func ExampleInstallPackageWith() {
	// Install packr2@v2.8.3 into the bin/ directory
	opts := pkg.InstallPackageOptions{
		Name:        "github.com/gobuffalo/packr/v2/packr2",
		Destination: "bin",
		Version:     "v2.8.3",
	}
	err := pkg.InstallPackageWith(opts)
	if err != nil {
		log.Fatal("could not install packr2@v2.8.3 into bin")
	}
}

func TestExampleDownloadToGopathBin(t *testing.T) {
	ExampleDownloadToGopathBin()
}

func ExampleDownloadToGopathBin() {
	url := "https://storage.googleapis.com/kubernetes-release/release/{{.VERSION}}/bin/{{.GOOS}}/{{.GOARCH}}/kubectl{{.EXT}}"
	err := pkg.DownloadToGopathBin(url, "kubectl", "v1.19.0")
	if err != nil {
		log.Fatal("could not download kubectl")
	}

	// Add GOPATH/bin to PATH if necessary so that we can immediately
	// use the installed tool
	gopath.EnsureGopathBin()
}
