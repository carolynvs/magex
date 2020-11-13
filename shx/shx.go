package shx

import (
	"bytes"
	"log"
	"strings"

	"github.com/magefile/mage/sh"
)

// Runs a command silently, without logging stdout/stderr.
func RunS(cmd string, args ...string) error {
	c := sh.Command(cmd, args...).Silent()

	_, _, err := c.Run()
	return err
}

// Capture stdout from a command, without logging stdout/stderr.
func OutputS(cmd string, args ...string) (string, error) {
	return sh.Command(cmd, args...).Silent().Output()
}

// Runs a command, only logging stderr when the command fails.
func RunE(cmd string, args ...string) error {
	_, _, err := sh.Command(cmd, args...).Stdout(nil).Run()
	return err
}

// Capture stdout from a command, only logging stderr when the command fails.
func OutputE(cmd string, args ...string) (string, error) {
	buf := &bytes.Buffer{}
	_, err := sh.Exec(nil, nil, buf, cmd, args...)
	if err != nil {
		log.Println(strings.TrimSuffix(buf.String(), "\n"))
	}
	return strings.TrimSuffix(buf.String(), "\n"), err
}

// CollapseArgs removes empty arguments from the argument list.
//
// This is helpful when sometimes a flag should be specified and
// sometimes it shouldn't.
func CollapseArgs(args ...string) []string {
	result := make([]string, 0, len(args))
	for _, arg := range args {
		if arg != "" {
			result = append(result, arg)
		}
	}
	return result
}
