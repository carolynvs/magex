package pkg

import (
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/carolynvs/magex/shx"
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
	err := shx.InDir(os.TempDir(), func() error {
		return shx.RunE("go", "get", "-u", pkg+moduleVersion)
	})
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

	src := JoinPath(GOPATH(), "src/github.com/magefile/mage")
	err = shx.InDir(src, func() error {
		return shx.RunE("go", "run", "bootstrap.go")
	})
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

// EnsureGopathBinInPath checks if GOPATH/bin is in PATH and adds it if necessary.
func EnsureGopathBinInPath() {
	path := os.Getenv("PATH")
	bin := JoinPath(GOPATH(), "bin")
	if !strings.Contains(path, bin) {
		log.Printf("Adding %s to $PATH\n", bin)
		sep := listPathSeparator(path)
		path += sep + bin
		os.Setenv("PATH", path)
	}
}
