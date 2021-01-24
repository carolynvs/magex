package shx_test

import (
	"log"

	"github.com/carolynvs/magex/shx"
)

func ExampleRun() {
	// Only write to stdout when mage -v is set
	err := shx.RunS("go", "run", "echo.go", "hello world")
	if err != nil {
		log.Fatal(err)
	}

	// Output:
}

func ExampleRunS() {
	// Do not write to stdout even when mage -v is set
	err := shx.RunS("go", "run", "echo.go", "hello world")
	if err != nil {
		log.Fatal(err)
	}

	// Output:
}

func ExampleRunE() {
	// Only write to stderr when the command fails
	err := shx.RunE("go", "run")
	if err == nil {
		log.Fatal("expected the command to fail")
	}

	// Output:
}

func ExampleRunV() {
	// Always print the output
	err := shx.RunV("go", "run", "echo.go", "hello world")
	if err != nil {
		log.Fatal(err)
	}

	// Output: hello world
}

func ExampleOutput() {
	// The output is printed only when mage -v is set
	output, err := shx.Output("go", "run", "echo.go", "hello world")
	if err != nil {
		log.Fatal(err)
	}

	if output != "hello world" {
		log.Fatal("expected to capture the output of the command")
	}

	// Output:
}

func ExampleOutputV() {
	// The output is printed every time
	output, err := shx.OutputV("go", "run", "echo.go", "hello world")
	if err != nil {
		log.Fatal(err)
	}

	if output != "hello world" {
		log.Fatal("expected to capture the output of the command")
	}

	// Output: hello world
}

func ExampleOutputS() {
	// Never write to stdout/stderr, just capture the output in a variable
	output, err := shx.OutputS("go", "run", "echo.go", "hello world")
	if err != nil {
		log.Fatal(err)
	}

	if output != "hello world" {
		log.Fatal(`expected to capture the output of the command`)
	}

	// Output:
}

func ExampleOutputE() {
	// Nothing should print when the command succeeds printed because the command passed
	output, err := shx.OutputE("go", "run", "echo.go", "hello world")
	if err != nil {
		log.Fatal(err)
	}

	if output != "hello world" {
		log.Fatal("expected to capture the output of the command")
	}

	// Output:
}
