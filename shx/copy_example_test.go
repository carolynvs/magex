package shx_test

import (
	"github.com/carolynvs/magex/shx"
)

func ExampleCopy() {
	// Copy a file from the current directory into TEMP
	shx.Copy("a.txt", "/tmp")

	// Copy matching files in the current directory into TEMP
	shx.Copy("*.txt", "/tmp")

	// Overwrite a file
	shx.Copy("/tmp/a.txt", "/tmp/b.txt")

	// Copy the contents of a directory into TEMP
	// Do not overwrite existing files
	shx.Copy("a/*", "/tmp", shx.CopyNoOverwrite)

	// Recursively copy a directory into TEMP
	shx.Copy("a", "/tmp", shx.CopyRecursive)
}
