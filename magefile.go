// +build mage

package main

import (
	"github.com/magefile/mage/sh"
)

var Default = Test

func Test() error {
	return sh.RunV("go", "test", "./...")
}
