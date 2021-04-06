package downloads

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"text/template"

	"github.com/carolynvs/magex/pkg/gopath"
	"github.com/carolynvs/magex/xplat"
	"github.com/pkg/errors"
)

// PostDownloadHook is the handler called after downloading a file, which returns the absolute path to the binary.
type PostDownloadHook func(archivePath string) (string, error)

// DownloadOptions
type DownloadOptions struct {
	// UrlTemplate is the Go template for the URL to download. Required.
	// Available Template Variables:
	//   - {{.GOOS}}
	//   - {{.GOARCH}}
	//   - {{.EXT}}
	//   - {{.VERSION}}
	UrlTemplate string

	// Name of the binary, excluding OS specific file extension. Required.
	Name string

	// Version to replace {{.VERSION}} in the URL template. Optional depending on whether or not the version is in the UrlTemplate.
	Version string

	// Ext to replace {{.EXT}} in the URL template. Optional, defaults to xplat.FileExt().
	Ext string

	// OsReplacement maps from a GOOS to the os keyword used for the download. Optional, defaults to empty.
	OsReplacement map[string]string

	// ArchReplacement maps from a GOARCH to the arch keyword used for the download. Optional, defaults to empty.
	ArchReplacement map[string]string

	// Hook to call after downloading the file.
	Hook PostDownloadHook
}

// DownloadToGopathBin takes a Go templated URL and expands template variables
// - srcTemplate is the URL
// - version is the version to substitute into the template
// - ext is the file extension to substitute into the template
//
// Template Variables:
// - {{.GOOS}}
// - {{.GOARCH}}
// - {{.EXT}}
// - {{.VERSION}}
func DownloadToGopathBin(opts DownloadOptions) error {
	src, err := RenderTemplate(opts.UrlTemplate, opts)
	if err != nil {
		return err
	}
	log.Printf("Downloading %s...", src)

	err = gopath.EnsureGopathBin()
	if err != nil {
		return err
	}

	// Download to a temp file
	tmpDir, err := ioutil.TempDir("", "magex")
	if err != nil {
		return errors.Wrap(err, "could not create temporary directory")
	}
	tmpFile := filepath.Join(tmpDir, filepath.Base(src))

	r, err := http.Get(src)
	if err != nil {
		return errors.Wrapf(err, "could not resolve %s", src)
	}
	defer r.Body.Close()

	f, err := os.OpenFile(tmpFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		return errors.Wrapf(err, "could not open %s", tmpFile)
	}
	defer f.Close()

	// Download to the temp file
	_, err = io.Copy(f, r.Body)
	if err != nil {
		errors.Wrapf(err, "error downloading %s", src)
	}
	f.Close()

	// Call a hook to allow for extracting or modifying the downloaded file
	var tmpBin = tmpFile
	if opts.Hook != nil {
		tmpBin, err = opts.Hook(tmpFile)
		if err != nil {
			return err
		}
	}

	// Make the binary executable
	err = os.Chmod(tmpBin, 0755)
	if err != nil {
		return errors.Wrapf(err, "could not make %s executable", tmpBin)
	}

	// Move it to GOPATH/bin
	dest := filepath.Join(gopath.GetGopathBin(), opts.Name+xplat.FileExt())
	err = os.Rename(tmpBin, dest)
	return errors.Wrapf(err, "error moving %s to %s", src, dest)
}

// RenderTemplate takes a Go templated string and expands template variables
// Available Template Variables:
// - {{.GOOS}}
// - {{.GOARCH}}
// - {{.EXT}}
// - {{.VERSION}}
func RenderTemplate(tmplContents string, opts DownloadOptions) (string, error) {
	tmpl, err := template.New("url").Parse(tmplContents)
	if err != nil {
		return "", errors.Wrapf(err, "error parsing %s as a Go template", opts.UrlTemplate)
	}

	srcData := struct {
		GOOS    string
		GOARCH  string
		EXT     string
		VERSION string
	}{
		GOOS:    runtime.GOOS,
		GOARCH:  runtime.GOARCH,
		EXT:     opts.Ext,
		VERSION: opts.Version,
	}

	if overrideGoos, ok := opts.OsReplacement[runtime.GOOS]; ok {
		srcData.GOOS = overrideGoos
	}

	if overrideGoarch, ok := opts.ArchReplacement[runtime.GOARCH]; ok {
		srcData.GOARCH = overrideGoarch
	}

	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, srcData)
	if err != nil {
		return "", errors.Wrapf(err, "error rendering %s as a Go template with data: %#v", opts.UrlTemplate, srcData)
	}

	return buf.String(), nil
}
