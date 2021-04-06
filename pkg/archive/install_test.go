package archive

import (
	"os"
	"os/exec"
	"runtime"
	"testing"

	"github.com/carolynvs/magex/pkg/downloads"
	"github.com/carolynvs/magex/pkg/gopath"
	"github.com/carolynvs/magex/xplat"
	"github.com/magefile/mage/mg"
	"github.com/stretchr/testify/require"
)

func TestDownloadArchiveToGopathBin(t *testing.T) {
	os.Setenv(mg.VerboseEnv, "true")
	err, cleanup := gopath.UseTempGopath()
	require.NoError(t, err, "Failed to set up a temporary GOPATH")
	defer cleanup()

	// gh cli unfortunately uses a different archive schema depending on the OS
	tmpl := "gh_{{.VERSION}}_{{.GOOS}}_{{.GOARCH}}/bin/gh{{.EXT}}"
	if runtime.GOOS == "windows" {
		tmpl = "bin/gh.exe"
	}

	opts := DownloadArchiveOptions{
		DownloadOptions: downloads.DownloadOptions{
			UrlTemplate: "https://github.com/cli/cli/releases/download/v{{.VERSION}}/gh_{{.VERSION}}_{{.GOOS}}_{{.GOARCH}}{{.EXT}}",
			Name:        "gh",
			Version:     "1.8.1",
			OsReplacement: map[string]string{
				"darwin": "macOS",
			},
		},
		ArchiveExtensions: map[string]string{
			"linux":   ".tar.gz",
			"darwin":  ".tar.gz",
			"windows": ".zip",
		},
		TargetFileTemplate: tmpl,
	}

	err = DownloadToGopathBin(opts)
	require.NoError(t, err)

	_, err = exec.LookPath("gh" + xplat.FileExt())
	require.NoError(t, err)
}
