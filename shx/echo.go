// +build ignore

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	// Read from stdin
	if len(os.Args) == 2 && os.Args[1] == "-" {
		msg, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(string(msg))
	} else {
		fmt.Println(strings.Join(os.Args[1:], " "))
	}
}
