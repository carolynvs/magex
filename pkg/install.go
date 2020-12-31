package pkg

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/carolynvs/magex/shx"
	"github.com/carolynvs/magex/xplat"
	"github.com/magefile/mage/sh"
	"github.com/pkg/errors"
)

// EnsureMage checks if mage is installed, and installs it if needed.
//
// When version is specified, detect if the specified version is installed, and
// if not, install that specific version of mage. Otherwise install the most
// recent code from the main branch.
func EnsureMage(version string) error {
	found, err := IsCommandAvailable("mage", version, "-version")
	if err != nil {
		return err
	}

	if !found {
		return InstallMage(version)
	}
	return nil
}

// EnsurePackage checks if the package is installed and installs it if needed.
//
// When version is specified, detect if the specified version is installed, and
// if not, install the package at that version. Otherwise install the most
// recent code from the main branch.
func EnsurePackage(pkg string, version string, versionArgs ...string) error {
	cmd := path.Base(pkg)

	found, err := IsCommandAvailable(cmd, version, versionArgs...)
	if err != nil {
		return err
	}

	if !found {
		return InstallPackage(pkg, version)
	}
	return nil
}

// InstallPackage installs the latest version of a package.
//
// When version is specified, install that version. Otherwise install the most
// recent code from the default branch.
func InstallPackage(pkg string, version string) error {
	cmd := path.Base(pkg)

	// Optionally install a specific version of the package
	moduleVersion := ""
	if version != "" {
		moduleVersion = "@" + version
	}

	log.Printf("Installing %s%s\n", cmd, moduleVersion)
	_, _, err := sh.Command("go", "get", "-u", pkg+moduleVersion).
		Stdout(nil).In(os.TempDir()).Run()
	if err != nil {
		return err
	}

	EnsureGopathBinInPath()

	// Check that it worked
	if found, _ := IsCommandAvailable(cmd, ""); !found {
		return errors.Errorf("Could not install %s. Please install it manually", pkg)
	}
	return nil
}

// InstallMage mage into GOPATH and add GOPATH/bin to PATH if necessary.
//
// When version is specified, install that version. Otherwise install the most
// recent code from the default branch.
func InstallMage(version string) error {
	err := InstallPackage("github.com/magefile/mage", version)
	if err != nil {
		return err
	}

	src := xplat.FilePathJoin(xplat.GOPATH(), "src/github.com/magefile/mage")
	_, _, err = sh.Command("go", "run", "bootstrap.go").
		Stdout(nil).In(src).Run()
	return errors.Wrap(err, "could not build mage with version info")
}

// IsCommandAvailable determines if a command can be called based on the current PATH.
func IsCommandAvailable(cmd string, version string, versionArgs ...string) (bool, error) {
	_, err := exec.LookPath(cmd)
	if err != nil {
		return false, nil
	}

	// If no version is specified, report that it is installed
	if version == "" {
		return true, nil
	}

	// Get the installed version
	versionOutput, err := shx.OutputE(cmd, versionArgs...)
	if err != nil {
		versionCmd := strings.Join(append([]string{cmd}, versionArgs...), " ")
		return false, errors.Wrapf(err, "could not determine the installed version of %s with '%s'", cmd, versionCmd)
	}

	versionFound := strings.Contains(versionOutput, version)
	return versionFound, nil
}

func getGopathBin() string {
	return xplat.FilePathJoin(xplat.GOPATH(), "bin")
}

// EnsureGopathBinInPath checks if GOPATH/bin is in PATH and adds it if necessary.
func EnsureGopathBinInPath() {
	xplat.EnsureInPath(getGopathBin())
}

func DownloadToGopathBin(src string, name string) error {
	log.Printf("Downloading %s to $GOPATH/bin\n", src)
	f, err := ioutil.TempFile("", path.Base(src))
	if err != nil {
		return errors.Wrap(err, "could not create temp file")
	}
	defer f.Close()

	err = os.Chmod(f.Name(), 0755)
	if err != nil {
		return errors.Wrapf(err, "could not make %s executable", f.Name())
	}

	r, err := http.Get(src)
	if err != nil {
		return errors.Wrapf(err, "could not resolve %s", src)
	}
	defer r.Body.Close()

	_, err = io.Copy(f, r.Body)
	if err != nil {
		errors.Wrapf(err, "error downloading %s", src)
	}
	f.Close()

	dest := filepath.Join(getGopathBin(), name)
	err = os.Rename(f.Name(), dest)
	return errors.Wrapf(err, "error moving %s to %s", src, dest)
}
