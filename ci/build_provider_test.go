package ci

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectBuildProvider(t *testing.T) {
	// Unset any variables that were set by OUR ci system :-)
	os.Unsetenv(AzureCIEnvVar)
	os.Unsetenv(GitHubCIEnvVar)

	t.Run("azure", func(t *testing.T) {
		os.Setenv(AzureCIEnvVar, "false")
		defer os.Unsetenv(AzureCIEnvVar)

		_, detected := DetectBuildProvider()
		require.False(t, detected)

		os.Setenv(AzureCIEnvVar, "true")

		p, detected := DetectBuildProvider()
		require.True(t, detected)
		assert.IsType(t, AzureBuildProvider{}, p)
	})

	t.Run("github", func(t *testing.T) {
		os.Setenv(GitHubCIEnvVar, "false")
		defer os.Unsetenv(GitHubCIEnvVar)

		_, detected := DetectBuildProvider()
		require.False(t, detected)

		os.Setenv(GitHubCIEnvVar, "true")

		p, detected := DetectBuildProvider()
		require.True(t, detected)
		assert.IsType(t, GitHubBuildProvider{}, p)
	})
}
