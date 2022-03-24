package ci

import (
	"fmt"
	"os"
	"strconv"
)

var _ BuildProvider = AzureBuildProvider{}

// AzureCIEnvVar is the environment variable used to detect the AzureBuildProvider.
const AzureCIEnvVar = "TF_BUILD"

// AzureBuildProvider supports Azure DevOps Pipelines.
type AzureBuildProvider struct{}

// SetEnv exports an environment variable. Changes from this command become
// available in subsequent steps in the CI pipeline. You must call os.SetEnv
// if you want to use the environment variable in the current process.
func (AzureBuildProvider) SetEnv(name string, value string) error {
	_, err := fmt.Printf("##vso[task.setvariable variable=%s]%s\n", name, value)
	return err
}

// PrependPath adds the specified path to the beginning of the PATH
// environment variable. Changes from this command become available in
// subsequent steps in the CI pipeline. You must call os.SetEnv if you want
// to use the PATH environment variable in the current process.
func (AzureBuildProvider) PrependPath(value string) error {
	_, err := fmt.Printf("##vso[task.prependpath]%s\n", value)
	return err
}

// IsDetected determines if this build provider was detected and is available
// to use.
func (AzureBuildProvider) IsDetected() bool {
	detected, _ := strconv.ParseBool(os.Getenv(AzureCIEnvVar))
	return detected
}
