package shx

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type MoveOption int

const (
	MoveDefault MoveOption = iota
	// MoveNoOverwrite does not overwrite existing files in the destination
	MoveNoOverwrite
	MoveRecursive
)

// Move a file or directory with the specified set of MoveOption.
// The source may use globbing, which is resolved with filepath.Glob.
func Move(src string, dest string, opts ...MoveOption) error {
	items, err := filepath.Glob(src)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		return fmt.Errorf("no such file or directory '%s'", src)
	}

	var combinedOpts MoveOption
	for _, opt := range opts {
		combinedOpts |= opt
	}

	// Check if the destination exists, e.g. if we are moving to /tmp/foo, /tmp should already exist
	if _, err := os.Stat(filepath.Dir(dest)); err != nil {
		return err
	}

	for _, item := range items {
		err := moveFileOrDirectory(item, dest, combinedOpts)
		if err != nil {
			return err
		}
	}

	return nil
}

func moveFileOrDirectory(src string, dest string, opts MoveOption) error {
	// If the destination is a directory that exists,
	// move into the directory.
	destInfo, err := os.Stat(dest)
	if err == nil && destInfo.IsDir() {
		dest = filepath.Join(dest, filepath.Base(src))
	}

	return move(src, dest, opts)
}

func move(src string, dest string, opts MoveOption) error {
	destExists := true
	destInfo, err := os.Stat(dest)
	if err != nil {
		if os.IsNotExist(err) {
			destExists = false
		} else {
			return err
		}
	}

	overwrite := opts&MoveNoOverwrite != MoveNoOverwrite
	if destExists {
		if overwrite {
			// Do not try to rename a file to an existing directory (mimics mv behavior)
			if destInfo.IsDir() {
				srcInfo, err := os.Stat(src)
				if err != nil {
					return err
				}
				if !srcInfo.IsDir() {
					return fmt.Errorf("rename %s to %s: not a directory", src, dest)
				}
			}

			os.RemoveAll(dest)
		} else {
			// Do not overwrite, skip
			log.Printf("%s not overwritten\n", dest)
			return nil
		}
	}

	log.Printf("%s -> %s\n", src, dest)
	return os.Rename(src, dest)
}
