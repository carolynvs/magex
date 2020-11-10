package shx

import (
	"bytes"
	"log"
	"os"
	"strings"

	"github.com/magefile/mage/sh"
	"github.com/pkg/errors"
)

// Runs a command silently, without logging stdout/stderr.
func RunS(cmd string, args ...string) error {
	_, err := sh.Exec(nil, nil, nil, cmd, args...)
	return err
}

// Capture stdout from a command, without logging stdout/stderr.
func OutputS(cmd string, args ...string) (string, error) {
	buf := &bytes.Buffer{}
	_, err := sh.Exec(nil, buf, nil, cmd, args...)
	return strings.TrimSuffix(buf.String(), "\n"), err
}

// Runs a command, only logging stderr when the command fails.
func RunE(cmd string, args ...string) error {
	buf := &bytes.Buffer{}
	_, err := sh.Exec(nil, nil, buf, cmd, args...)
	if err != nil {
		log.Println(strings.TrimSuffix(buf.String(), "\n"))
	}
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

// InDir executes a command in the specified directory.
//
// TODO: Consider making these extensions fluent, e.g. RunE().WithDir().Execute()
func InDir(dir string, cmd func() error) error {
	pwd, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err, "could not determine current working directory")
	}

	dir = os.ExpandEnv(dir)
	err = os.Chdir(dir)
	if err != nil {
		return errors.Wrapf(err, "could not change directory to %s", dir)
	}
	defer os.Chdir(pwd)

	return cmd()
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
