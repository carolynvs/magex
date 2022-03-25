package ci

import (
	"fmt"
	"os"
	"strconv"
)

var _ BuildProvider = GitHubBuildProvider{}

const (
	// GitHubCIEnvVar is the environment variable used to detect the
	// GitHubBuildProvider.
	GitHubCIEnvVar = "GITHUB_ACTIONS"

	// GitHubVariablesEnvVar is a GitHub environment variable that contains the
	// path to a file where you can set variable assignments.
	GitHubVariablesEnvVar = "GITHUB_ENV"

	// GitHubPathEnvVar is a GitHub environment variable that contains the path
	// to a file where you can prepend PATH values.
	GitHubPathEnvVar = "GITHUB_PATH"
)

// GitHubBuildProvider supports GitHub Actions.
type GitHubBuildProvider struct{}

// SetEnv exports an environment variable. Changes from this command become
// available in subsequent steps in the CI pipeline. You must call os.SetEnv
// if you want to use the environment variable in the current process.
func (p GitHubBuildProvider) SetEnv(name string, value string) error {
	assignment := fmt.Sprintf("%s=%s", name, value)
	return p.appendFile(GitHubVariablesEnvVar, assignment)
}

// PrependPath adds the specified path to the beginning of the PATH
// environment variable. Changes from this command become available in
// subsequent steps in the CI pipeline. You must call os.SetEnv if you want
// to use the PATH environment variable in the current process.
func (p GitHubBuildProvider) PrependPath(value string) error {
	return p.appendFile(GitHubPathEnvVar, value)
}

// IsDetected determines if this build provider was detected and is available
// to use.
func (p GitHubBuildProvider) IsDetected() bool {
	detected, _ := strconv.ParseBool(os.Getenv(GitHubCIEnvVar))
	return detected
}

func (p GitHubBuildProvider) appendFile(envVar string, line string) error {
	path := os.Getenv(envVar)

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0660)
	if err != nil {
		return fmt.Errorf("could not open the file referenced by %s: %w", envVar, err)
	}

	_, err = fmt.Fprintln(f, line)
	if err != nil {
		return fmt.Errorf("could not write to the file referenced by %s: %w", envVar, err)
	}

	return nil
}
