// +build mage

package main

import (
	"fmt"

	"github.com/carolynvs/magex/xplat"
	"github.com/magefile/mage/sh"
)

var Default = Test

func Test() error {
	fmt.Println("Running tests on", xplat.DetectShell())
	return sh.RunV("go", "test", "-v", "./...")
}
