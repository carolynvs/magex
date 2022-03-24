package ci

func ExampleAzureBuildProvider_SetEnv() {
	p := AzureBuildProvider{}
	p.SetEnv("FOO", "1")
	p.SetEnv("BAR", "A")

	// Output: ##vso[task.setvariable variable=FOO]1
	// ##vso[task.setvariable variable=BAR]A
}

func ExampleAzureBuildProvider_PrependPath() {
	p := AzureBuildProvider{}
	p.PrependPath("/usr/bin")
	p.PrependPath("/home/me/bin")

	// Output: ##vso[task.prependpath]/usr/bin
	// ##vso[task.prependpath]/home/me/bin
}
