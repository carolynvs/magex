package pkg_test

import (
	"log"
	"testing"

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
	// Install packr2@v2.8.0 using the command `packr2 version` to detect if the
	// correct version is installed.
	err := pkg.EnsurePackage("github.com/gobuffalo/packr/v2/packr2", "v2.8.3", "version")
	if err != nil {
		log.Fatal("could not install packr2")
	}
}

func TestExampleEnsurePackage_WithVersionConstraint(t *testing.T) {
	ExampleEnsurePackage_WithVersionConstraint()
}

func ExampleEnsurePackage_WithVersionConstraint() {
	// Install packr2@v2.8.0 using the command `packr2 version` to detect if
	// any v2 version is installed
	err := pkg.EnsurePackage("github.com/gobuffalo/packr/v2/packr2", "v2.8.3", "version", "2.x")
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
