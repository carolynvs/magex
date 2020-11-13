package shx_test

import (
	"fmt"
	"log"

	"github.com/carolynvs/magex/shx"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

func ExampleRunS() {
	err := shx.RunS("bash", "-c", "echo hello world")
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	// Output:
}

func ExampleRunE() {
	err := shx.RunE("bash", "-c", "oops")
	if err == nil {
		fmt.Println("error was expected")
	}
	// Output:
}

func ExampleOutputS() {
	output, err := shx.OutputS("bash", "-c", "echo hello world")
	if err != nil {
		log.Fatal(err)
	}
	if output != "hello world" {
		log.Fatal(`expected to capture "hello world"`)
	}
	// Output:
}

func ExampleOutputE() {
	// Get the output of `docker --version`, only logging stderr when the command fails.
	versionOutput, err := shx.OutputE("docker", "--version")
	if err != nil {
		log.Println("could not get the docker version")
	}

	fmt.Println(versionOutput)
}

func ExampleCollapseArgs() {
	// Only pass -v to go test when the target was called with -v
	// mage -v test -> go test -v ./...
	// mage test -> go test ./...
	v := ""
	if mg.Verbose() {
		v = "-v"
	}
	sh.RunV("go", "test", v, "./...")
}
