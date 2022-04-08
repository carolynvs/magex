package downloads

import (
	"github.com/carolynvs/magex/pkg/gopath"
	"github.com/carolynvs/magex/xplat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadToGopathBin(t *testing.T) {
	err, cleanup := gopath.UseTempGopath()
	require.NoError(t, err, "Failed to set up a temporary GOPATH")
	defer cleanup()

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("echo ok"))
	}))
	defer svr.Close()

	opts := DownloadOptions{
		UrlTemplate: svr.URL,
		Name:        "mybin",
	}
	err = DownloadToGopathBin(opts)
	require.NoError(t, err)
	assert.FileExists(t, filepath.Join(gopath.GetGopathBin(), "mybin"+xplat.FileExt()))
}

func TestDownload(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
		}))
		defer svr.Close()

		opts := DownloadOptions{
			UrlTemplate: svr.URL,
			Name:        "mybin",
		}
		err := Download("bin", opts)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "404 Not Found")
	})

	t.Run("found", func(t *testing.T) {
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("echo ok"))
		}))
		defer svr.Close()

		dest, err := ioutil.TempDir("", "magex")
		require.NoError(t, err)
		defer os.RemoveAll(dest)

		opts := DownloadOptions{
			UrlTemplate: svr.URL,
			Name:        "mybin",
		}
		err = Download(dest, opts)
		require.NoError(t, err)
		assert.FileExists(t, filepath.Join(dest, "mybin"+xplat.FileExt()))
	})
}
