//go:build mage
// +build mage

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/carolynvs/magex/pkg"
	"github.com/carolynvs/magex/shx"
	"github.com/carolynvs/magex/xplat"
	"github.com/magefile/mage/mg"
)

var Default = Test

func Test() error {
	tmpHome, err := ioutil.TempDir("", "magex")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpHome)

	fmt.Println("Running tests on", xplat.DetectShell())
	var v string
	if mg.Verbose() {
		v = "-v"
	}
	return shx.Command("go", "test", v, "./...").CollapseArgs().
		Env("GOPATH=" + tmpHome).RunV()
}

func EnsureMage() error {
	return pkg.EnsureMage("v1.11.0")
}
