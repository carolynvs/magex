package ci

// BuildProvider is a common interface to interact with a CI build provider
// such as GitHub Actions, or Azure DevOps.
type BuildProvider interface {
	// SetEnv exports an environment variable. Changes from this command become
	// available in subsequent steps in the CI pipeline. You must call os.SetEnv
	// if you want to use the environment variable in the current process.
	SetEnv(name string, value string) error

	// PrependPath adds the specified path to the beginning of the PATH
	// environment variable. Changes from this command become available in
	// subsequent steps in the CI pipeline. You must call os.SetEnv if you want
	// to use the PATH environment variable in the current process.
	PrependPath(value string) error

	// IsDetected determines if this build provider was detected and is available
	// to use.
	IsDetected() bool
}

// DetectBuildProvider determines the current build provider that the code is
// executing upon, returning a NoopBuildProvider and false when nothing is
// detected. By default, only build providers implemented in this package are
// included in the search. Specify additional providers with the providers
// argument.
func DetectBuildProvider(providers ...BuildProvider) (BuildProvider, bool) {
	providers = append(providers, AzureBuildProvider{}, GitHubBuildProvider{})
	for _, provider := range providers {
		if provider.IsDetected() {
			return provider, true
		}
	}

	return NoopBuildProvider{}, false
}
