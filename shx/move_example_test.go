package shx_test

import "github.com/carolynvs/magex/shx"

func ExampleMove() {
	// Move a file from the current directory into TEMP
	shx.Move("a.txt", "/tmp")

	// Move matching files in the current directory into TEMP
	shx.Move("*.txt", "/tmp")

	// Overwrite a file
	shx.Move("/tmp/a.txt", "/tmp/b.txt")

	// Move the contents of a directory into TEMP
	// Do not overwrite existing files
	shx.Move("a/*", "/tmp", shx.MoveNoOverwrite)
}
