package downloads

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDownloadToGopathBin(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(404)
		}))
		defer svr.Close()

		t.Run("not found", func(t *testing.T) {
			opts := DownloadOptions{
				UrlTemplate: svr.URL,
			}
			err := DownloadToGopathBin(opts)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "404 Not Found")
		})
	})
}
