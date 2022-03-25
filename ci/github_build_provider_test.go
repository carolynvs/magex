package ci

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGitHubBuildProvider_SetEnv(t *testing.T) {
	tmp, err := ioutil.TempFile("", "magex")
	require.NoError(t, err)
	defer os.Remove(tmp.Name())

	os.Setenv(GitHubVariablesEnvVar, tmp.Name())
	defer os.Unsetenv(GitHubVariablesEnvVar)

	p := GitHubBuildProvider{}
	err = p.SetEnv("FOO", "1")
	require.NoError(t, err)
	err = p.SetEnv("BAR", "A")
	require.NoError(t, err)

	contents, err := ioutil.ReadFile(tmp.Name())
	require.NoError(t, err)
	assert.Contains(t, string(contents), "FOO=1\nBAR=A\n")
}

func TestGitHubBuildProvider_PrependPath(t *testing.T) {
	tmp, err := ioutil.TempFile("", "magex")
	require.NoError(t, err)
	defer os.Remove(tmp.Name())

	os.Setenv(GitHubPathEnvVar, tmp.Name())
	defer os.Unsetenv(GitHubPathEnvVar)

	p := GitHubBuildProvider{}
	err = p.PrependPath("/usr/bin")
	require.NoError(t, err)
	err = p.PrependPath("/home/me/bin")
	require.NoError(t, err)

	contents, err := ioutil.ReadFile(tmp.Name())
	require.NoError(t, err)
	assert.Contains(t, string(contents), "/usr/bin\n/home/me/bin\n")
}
