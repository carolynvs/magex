package pkg

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/carolynvs/magex/pkg/downloads"
	"github.com/carolynvs/magex/pkg/gopath"
	"github.com/carolynvs/magex/shx"
	"github.com/carolynvs/magex/xplat"
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
	cmd := getCommandName(pkg)

	found, err := IsCommandAvailable(cmd, version, versionArgs...)
	if err != nil {
		return err
	}

	if !found {
		return InstallPackage(pkg, version)
	}
	return nil
}

func getCommandName(pkg string) string {
	name := path.Base(pkg)
	if ok, _ := regexp.MatchString(`v[\d]+`, name); ok {
		return getCommandName(path.Dir(pkg))
	}
	return name
}

// InstallPackage installs the latest version of a package.
//
// When version is specified, install that version. Otherwise install the most
// recent code from the default branch.
func InstallPackage(pkg string, version string) error {
	gopath.EnsureGopathBin()

	cmd := getCommandName(pkg)

	// Optionally install a specific version of the package
	moduleVersion := ""
	if version != "" {
		if strings.HasPrefix(version, "v") {
			moduleVersion = "@" + version
		} else {
			moduleVersion = "@v" + version
		}
	}

	fmt.Printf("Installing %s%s\n", cmd, moduleVersion)
	err := shx.Command("go", "install", pkg+moduleVersion).
		Env("GO111MODULE=on").In(os.TempDir()).RunE()
	if err != nil {
		return err
	}

	log.Printf("Checking if %s is accessible from the PATH", cmd)
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
	var tag string
	if version != "" {
		tag = "-b" + version
	}

	tmp, err := ioutil.TempDir("", "magefile")
	if err != nil {
		return errors.Wrap(err, "could not create a temp directory to install mage")
	}
	defer os.RemoveAll(tmp)

	repoUrl := "https://github.com/magefile/mage.git"
	err = shx.Command("git", "clone", tag, repoUrl).CollapseArgs().In(tmp).RunE()
	if err != nil {
		return errors.Wrapf(err, "could not clone %s", repoUrl)
	}

	repoPath := filepath.Join(tmp, "mage")
	err = shx.Command("go", "run", "bootstrap.go").In(repoPath).RunE()
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

// DownloadToGopathBin downloads an executable file to GOPATH/bin.
// src can include the following template values:
//   - {{.GOOS}}
//   - {{.GOARCH}}
//   - {{.EXT}}
//   - {{.VERSION}}
func DownloadToGopathBin(srcTemplate string, name string, version string) error {
	opts := downloads.DownloadOptions{
		UrlTemplate: srcTemplate,
		Name:        name,
		Version:     version,
		Ext:         xplat.FileExt(),
	}
	return downloads.DownloadToGopathBin(opts)
}
