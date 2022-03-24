package ci

var _ BuildProvider = NoopBuildProvider{}

// NoopBuildProvider is a build provider that does nothing.
type NoopBuildProvider struct{}

// SetEnv does nothing.
func (n NoopBuildProvider) SetEnv(string, string) error { return nil }

// PrependPath does nothing.
func (n NoopBuildProvider) PrependPath(string) error { return nil }

// IsDetected always returns false.
func (n NoopBuildProvider) IsDetected() bool { return false }
