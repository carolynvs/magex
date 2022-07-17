package pkg

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/carolynvs/magex/pkg/downloads"
	"github.com/carolynvs/magex/pkg/gopath"
	"github.com/carolynvs/magex/shx"
	"github.com/carolynvs/magex/xplat"
)

// EnsureMage checks if mage is installed, and installs it if needed.
//
// When version is specified, detect if a compatible version is installed, and
// if not, install that specific version of mage. Otherwise, install the most
// recent code from the main branch.
func EnsureMage(defaultVersion string) error {
	versionConstraint := makeDefaultVersionConstraint(defaultVersion)
	found, err := IsCommandAvailable("mage", "-version", versionConstraint)
	if err != nil {
		return err
	}

	if !found {
		return InstallMage(defaultVersion)
	}
	return nil
}

// EnsurePackage checks if the package is installed and installs it if needed.
// Optionally accepts the argument or flag to pass to the command to check the
// installed version, and a semver range to use to validate the installed
// version, such as ^1.2.3 or 2.x. When no version arguments are supplied, any
// installed version is acceptable.
//
// When defaultVersion is specified, and a version constraint is not, the default
// is used as the minimum version and sets the allowed major version. For example,
// a defaultVersion of 1.2.3 would result in a constraint of ^1.2.3.
// When no defaultVersion is specified, the latest version is installed.
func EnsurePackage(pkg string, defaultVersion string, versionArgs ...string) error {
	cmd := getCommandName(pkg)

	// Apply optional arguments: versionCmd, versionConstraint
	versionCmd := ""
	versionConstraint := ""
	if len(versionArgs) > 0 {
		versionCmd = versionArgs[0]
		if len(versionArgs) > 1 {
			versionConstraint = versionArgs[1]
		}
	}

	// Default the constraint to [defaultVersion - next major)
	if versionConstraint == "" {
		versionConstraint = makeDefaultVersionConstraint(defaultVersion)
	}

	found, err := IsCommandAvailable(cmd, versionCmd, versionConstraint)
	if err != nil {
		return err
	}

	if !found {
		return InstallPackage(pkg, defaultVersion)
	}
	return nil
}

// create a semver constraint of ^defaultVersion, otherwise use no constraint
func makeDefaultVersionConstraint(defaultVersion string) string {
	defaultVersion = strings.TrimPrefix(defaultVersion, "v")
	if v, err := semver.NewVersion(defaultVersion); err == nil {
		return fmt.Sprintf("^%s", v.String())
	}
	return ""
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
// recent version.
func InstallPackage(pkg string, version string) error {
	gopath.EnsureGopathBin()

	cmd := getCommandName(pkg)

	if version == "" {
		version = "latest"
	} else {
		if version != "latest" && !strings.HasPrefix(version, "v") {
			version = "v" + version
		}
	}

	fmt.Printf("Installing %s@%s\n", cmd, version)
	return shx.Command("go", "install", pkg+"@"+version).
		Env("GO111MODULE=on").In(os.TempDir()).RunE()
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
		return fmt.Errorf("could not create a temp directory to install mage: %w", err)
	}
	defer os.RemoveAll(tmp)

	repoUrl := "https://github.com/magefile/mage.git"
	err = shx.Command("git", "clone", tag, repoUrl).CollapseArgs().In(tmp).RunE()
	if err != nil {
		return fmt.Errorf("could not clone %s: %w", repoUrl, err)
	}

	repoPath := filepath.Join(tmp, "mage")
	if err := shx.Command("go", "run", "bootstrap.go").In(repoPath).RunE(); err != nil {
		return fmt.Errorf("could not build mage with version info: %w", err)
	}

	return nil
}

// IsCommandAvailable determines if a command can be called based on the current PATH.
func IsCommandAvailable(cmd string, versionCmd string, versionConstraint string) (bool, error) {
	cmd, err := exec.LookPath(cmd)
	if err != nil {
		return false, nil
	}

	return CheckCommandVersion(cmd, versionCmd, versionConstraint)
}

// CheckCommandVersion determines if the specified command is available and
// if specified, that the version command returned a version that matches the semver constraint.
// For example, 1.x or ~2.3. See https://github.com/Masterminds/semver for details
// on how to specify a version constrain.
func CheckCommandVersion(cmd string, versionCmd string, versionConstraint string) (bool, error) {
	// Get the installed version
	scrapedVersion, err := GetCommandVersion(cmd, versionCmd)
	if err != nil {
		return false, err
	}

	// Parse the version from the command output as a semantic version
	ver, err := semver.NewVersion(scrapedVersion)
	if err != nil {
		return true, nil
	}

	// Try to use the version constraint, ignore it if it's not a valid semver constraint
	// such as "" or "latest" or a tag/branch
	constraint, err := semver.NewConstraint(versionConstraint)
	if err != nil {
		return true, nil
	}

	return constraint.Check(ver), nil
}

// This is the same regex that masterminds/semver uses
const semVerRegex string = `v?([0-9]+)(\.[0-9]+)?(\.[0-9]+)?` +
	`(-([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?` +
	`(\+([0-9A-Za-z\-]+(\.[0-9A-Za-z\-]+)*))?`

// GetCommandVersion executes the specified command to get its version
// The result is the contents of standard output of calling the command, and
// probably includes additional text besides the version number.
func GetCommandVersion(cmd string, versionCmd string) (string, error) {
	prettyCmd := cmd
	if versionCmd != "" {
		prettyCmd = fmt.Sprintf("%s %s", cmd, versionCmd)
	}

	// Get the installed version
	versionOutput, err := shx.Command(cmd, versionCmd).CollapseArgs().OutputE()
	if err != nil {
		return "", fmt.Errorf("could not determine the installed version of %s with '%s': %w", cmd, prettyCmd, err)
	}

	cmdRegex := regexp.MustCompile(semVerRegex)
	matches := cmdRegex.FindStringSubmatch(versionOutput)
	if len(matches) == 0 {
		return "", fmt.Errorf("the output of %s did not include a 3-part semver value: %s", prettyCmd, versionOutput)
	}
	return matches[0], nil
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
