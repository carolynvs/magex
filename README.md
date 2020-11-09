# Magefile Extensions

![test](https://github.com/carolynvs/magex/workflows/test/badge.svg)

This library provides helper methods to use with [mage](https://magefile.org).

Below is a sample of the type of helpers available. Full examples and
documentation is on [godoc](godoc.org/github.com/carolynvs/magex).

```go
// +build mage

package main

import (
	"github.com/carolynvs/magex"
	"github.com/carolynvs/magex/shx"
)

// Install packr2 v2.8.0 if it's not available, and ensure it's in PATH.
func Packr2() error {
    if magex.IsCommandAvailable("packr2") {
        err := magex.EnsurePackage("github.com/gobuffalo/packr/v2/packr2/cmd", "v2.8.0", "version")
        if err != nil {
            log.Fatal("could not install packr2")
        }
    }
}

// Install mage if it's not available, and ensure it's in PATH. We don't care which version
func Mage() error {
    return magex.EnsureMage("")
}

// Run a docker registry in a container. Do not print stdout and only print
// stderr when the command fails even when -v is set.
//
// Useful for commands that you only care about when it fails, keeping unhelpful
// output out of your logs.
func StartRegistry() error {
    return shx.RunE("docker", "run", "-d", "-p", "5000:5000", "--name", "registry", "registry:2")
}
```
