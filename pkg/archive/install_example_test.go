package archive_test

import (
	"log"

	"github.com/carolynvs/magex/pkg/archive"
	"github.com/carolynvs/magex/pkg/downloads"
	"github.com/carolynvs/magex/pkg/gopath"
)

func ExampleDownloadToGopathBin() {
	opts := archive.DownloadArchiveOptions{
		DownloadOptions: downloads.DownloadOptions{
			UrlTemplate: "https://get.helm.sh/helm-{{.VERSION}}-{{.GOOS}}-{{.GOARCH}}{{.EXT}}",
			Name:        "helm",
			Version:     "v3.5.3",
		},
		ArchiveExtensions: map[string]string{
			"darwin":  ".tar.gz",
			"linux":   ".tar.gz",
			"windows": ".zip",
		},
		TargetFileTemplate: "{{.GOOS}}-{{.GOARCH}}/helm{{.EXT}}",
	}
	err := archive.DownloadToGopathBin(opts)
	if err != nil {
		log.Fatal("could not download helm")
	}

	// Add GOPATH/bin to PATH if necessary so that we can immediately
	// use the installed tool
	gopath.EnsureGopathBin()
}
