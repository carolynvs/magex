package pkg

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

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
	EnsureGopathBin()

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

	log.Printf("Installing %s%s\n", cmd, moduleVersion)
	err := shx.Command("go", "get", "-u", pkg+moduleVersion).
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
	err := InstallPackage("github.com/magefile/mage", version)
	if err != nil {
		return err
	}

	src := filepath.Join(GOPATH(), "src/github.com/magefile/mage")
	err = shx.Command("go", "run", "bootstrap.go").In(src).RunE()
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

// GetGopathBin returns GOPATH/bin.
func GetGopathBin() string {
	return filepath.Join(GOPATH(), "bin")
}

// EnsureGopathBin ensures that GOPATH/bin exists and is in PATH.
// Detects if this is an Azure CI build and exports the updated PATH.
func EnsureGopathBin() error {
	gopathBin := GetGopathBin()
	err := os.MkdirAll(gopathBin, 0755)
	if err != nil {
		errors.Wrapf(err, "could not create GOPATH/bin at %s", gopathBin)
	}
	xplat.EnsureInPath(GetGopathBin())
	return nil
}

// DownloadToGopathBin downloads an executable file to GOPATH/bin.
// src can include the following template values:
//   - {{.GOOS}}
//   - {{.GOARCH}}
//   - {{.EXT}}
//   - {{.VERSION}}
func DownloadToGopathBin(srcTemplate string, name string, version string) error {
	src, err := renderUrlTemplate(srcTemplate, version)
	if err != nil {
		return err
	}
	log.Printf("Downloading %s to $GOPATH/bin\n", src)

	err = EnsureGopathBin()
	if err != nil {
		return err
	}

	// Download to a temp file
	f, err := ioutil.TempFile("", path.Base(src))
	if err != nil {
		return errors.Wrap(err, "could not create temp file")
	}
	defer f.Close()

	// Make it executable
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

	// Move it to GOPATH/bin
	dest := filepath.Join(GetGopathBin(), name+xplat.FileExt())
	err = os.Rename(f.Name(), dest)
	return errors.Wrapf(err, "error moving %s to %s", src, dest)
}

func renderUrlTemplate(srcTemplate string, version string) (string, error) {
	tmpl, err := template.New("url").Parse(srcTemplate)
	if err != nil {
		return "", errors.Wrapf(err, "error parsing %s as a Go template", srcTemplate)
	}

	srcData := struct {
		GOOS    string
		GOARCH  string
		EXT     string
		VERSION string
	}{
		GOOS:    runtime.GOOS,
		GOARCH:  runtime.GOARCH,
		EXT:     xplat.FileExt(),
		VERSION: version,
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, srcData)
	if err != nil {
		return "", errors.Wrapf(err, "error rendering %s as a Go template with data: %#v", srcTemplate, srcData)
	}

	return buf.String(), nil
}
