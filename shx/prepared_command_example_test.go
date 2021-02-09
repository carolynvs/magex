package shx_test

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/carolynvs/magex/shx"
)

func ExamplePreparedCommand_RunV() {
	err := shx.Command("go", "run", "echo.go", "hello world").RunV()
	if err != nil {
		log.Fatal(err)
	}
	// Output: hello world
}

func ExamplePreparedCommand_Args() {
	cmd := shx.Command("go", "run", "echo.go")
	// Append arguments to the command
	err := cmd.Args("hello", "world").RunV()
	if err != nil {
		log.Fatal(err)
	}

	// Output: hello world
}

func ExamplePreparedCommand_In() {
	tmp, err := ioutil.TempDir("", "mage")
	if err != nil {
		log.Fatal(err)
	}

	contents := `package main

import "fmt"

func main() {
	fmt.Println("hello world")
}
`
	err = ioutil.WriteFile(filepath.Join(tmp, "test_main.go"), []byte(contents), 0644)
	if err != nil {
		log.Fatal(err)
	}

	// Run `go run test_main.go` in /tmp
	err = shx.Command("go", "run", "test_main.go").In(tmp).RunV()
	if err != nil {
		log.Fatal(err)
	}
	// Output: hello world
}

func ExamplePreparedCommand_RunS() {
	err := shx.Command("go", "run", "echo.go", "hello world").RunS()
	if err != nil {
		log.Fatal(err)
	}
	// Output:
}

func ExamplePreparedCommand_CollapseArgs() {
	err := shx.Command("go", "run", "echo.go", "hello", "", "world").CollapseArgs().RunV()
	if err != nil {
		log.Fatal(err)
	}

	// Output: hello world
}
