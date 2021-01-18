package shx_test

import (
	"log"

	"github.com/carolynvs/magex/shx"
)

func ExampleRunS() {
	err := shx.RunS("go", "run", "echo.go", "hello world")
	if err != nil {
		log.Fatal(err)
	}

	// Output:
}

func ExampleRunE() {
	err := shx.RunE("go", "run")
	if err == nil {
		log.Fatal("expected the command to fail")
	}

	// Output:
}

func ExampleOutputS() {
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
	output, err := shx.OutputE("go", "run", "echo.go", "hello world")
	if err != nil {
		log.Fatal(err)
	}

	if output != "hello world\n" {
		log.Fatal("expected to capture the output of the command")
	}

	// Nothing is printed because the command passed
	// Output:
}
