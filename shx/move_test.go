package shx

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func resetTestdata(t *testing.T) {
	err := exec.Command("git", "checkout", "testdata").Run()
	require.NoError(t, err, "error resetting the testdata directory")
}

func TestMove(t *testing.T) {
	t.Run("recursively move directory into empty dest dir", func(t *testing.T) {
		defer resetTestdata(t)

		// Make the temp directory on the same physical drive, os.Rename doesn't work across drives and /tmp may be on another drive
		tmp, err := os.MkdirTemp("testdata", "magex")
		require.NoError(t, err, "could not create temp directory for test")
		defer os.RemoveAll(tmp)

		err = Move("testdata/copy/a", tmp, MoveRecursive)
		require.NoError(t, err, "Move into empty directory failed")

		assert.DirExists(t, filepath.Join(tmp, "a"))
		assertFile(t, filepath.Join(tmp, "a/a1.txt"))
		assertFile(t, filepath.Join(tmp, "a/a2.txt"))
		assert.DirExists(t, filepath.Join(tmp, "a/ab"))
		assertFile(t, filepath.Join(tmp, "a/ab/ab1.txt"))
		assertFile(t, filepath.Join(tmp, "a/ab/ab2.txt"))
	})

	t.Run("recursively move directory into populated dest dir", func(t *testing.T) {
		defer resetTestdata(t)

		// Make the temp directory on the same physical drive, os.Rename doesn't work across drives and /tmp may be on another drive
		tmp, err := os.MkdirTemp("testdata", "magex")
		require.NoError(t, err, "could not create temp directory for test")
		defer os.RemoveAll(tmp)

		require.NoError(t, os.MkdirAll(filepath.Join(tmp, "a"), 0700))
		require.NoError(t, os.WriteFile(filepath.Join(tmp, "a/a1.txt"), []byte("a lot of extra data that should be overwritten"), 0600))

		err = Move("testdata/copy/a", tmp, MoveRecursive)
		require.NoError(t, err, "Move into directory with same directory name")

		assert.DirExists(t, filepath.Join(tmp, "a"))
		assertFile(t, filepath.Join(tmp, "a/a1.txt"))
		assertFile(t, filepath.Join(tmp, "a/a2.txt"))
		assert.DirExists(t, filepath.Join(tmp, "a/ab"))
		assertFile(t, filepath.Join(tmp, "a/ab/ab1.txt"))
		assertFile(t, filepath.Join(tmp, "a/ab/ab2.txt"))
	})

	t.Run("move glob", func(t *testing.T) {
		defer resetTestdata(t)

		// Make the temp directory on the same physical drive, os.Rename doesn't work across drives and /tmp may be on another drive
		tmp, err := os.MkdirTemp("testdata", "magex")
		require.NoError(t, err, "could not create temp directory for test")
		defer os.RemoveAll(tmp)

		err = Move("testdata/copy/a/*.txt", tmp)
		require.NoError(t, err, "Move into empty directory failed")

		assertFile(t, filepath.Join(tmp, "a1.txt"))
		assertFile(t, filepath.Join(tmp, "a2.txt"))
		assert.NoDirExists(t, filepath.Join(tmp, "a"))
		assert.NoDirExists(t, filepath.Join(tmp, "ab"))
	})

	t.Run("missing parent dir", func(t *testing.T) {
		defer resetTestdata(t)

		// Make the temp directory on the same physical drive, os.Rename doesn't work across drives and /tmp may be on another drive
		tmp, err := os.MkdirTemp("testdata", "magex")
		require.NoError(t, err, "could not create temp directory for test")
		defer os.RemoveAll(tmp)

		err = Move("testdata/copy/a", filepath.Join(tmp, "missing-parent/dir"))
		require.Error(t, err)
	})

	t.Run("missing src", func(t *testing.T) {
		defer resetTestdata(t)

		// Make the temp directory on the same physical drive, os.Rename doesn't work across drives and /tmp may be on another drive
		tmp, err := os.MkdirTemp("testdata", "magex")
		require.NoError(t, err, "could not create temp directory for test")
		defer os.RemoveAll(tmp)

		err = Move("testdata/missing-src", tmp)
		require.Error(t, err)
		require.Contains(t, err.Error(), "no such file or directory")
	})

	t.Run("recursively move directory to new name", func(t *testing.T) {
		defer resetTestdata(t)

		// Make the temp directory on the same physical drive, os.Rename doesn't work across drives and /tmp may be on another drive
		tmp, err := os.MkdirTemp("testdata", "magex")
		require.NoError(t, err, "could not create temp directory for test")
		defer os.RemoveAll(tmp)

		dest := filepath.Join(tmp, "dest")
		err = Move("testdata/copy/a", dest, MoveRecursive)
		require.NoError(t, err, "Move into empty directory failed")

		assert.DirExists(t, dest)
		assertFile(t, filepath.Join(dest, "a1.txt"))
		assertFile(t, filepath.Join(dest, "a2.txt"))
		assert.DirExists(t, filepath.Join(dest, "ab"))
		assertFile(t, filepath.Join(dest, "ab/ab1.txt"))
		assertFile(t, filepath.Join(dest, "ab/ab2.txt"))
	})

	t.Run("recursively merge dest dir", func(t *testing.T) {
		defer resetTestdata(t)

		// Make the temp directory on the same physical drive, os.Rename doesn't work across drives and /tmp may be on another drive
		tmp, err := os.MkdirTemp("testdata", "magex")
		require.NoError(t, err, "could not create temp directory for test")
		defer os.RemoveAll(tmp)

		err = Move("testdata/copy/partial-dest", tmp, MoveRecursive)
		require.NoError(t, err, "Move partial destination failed")

		err = Move("testdata/copy/a", tmp, MoveRecursive)
		require.NoError(t, err, "Merge into non-empty destination failed")

		assert.DirExists(t, filepath.Join(tmp, "a"))
		assertFile(t, filepath.Join(tmp, "a/a1.txt"))
		assertFile(t, filepath.Join(tmp, "a/a2.txt"))
		assert.DirExists(t, filepath.Join(tmp, "a/ab"))
		assertFile(t, filepath.Join(tmp, "a/ab/ab1.txt"))
		assertFile(t, filepath.Join(tmp, "a/ab/ab2.txt"))
	})

	t.Run("move file into empty directory", func(t *testing.T) {
		defer resetTestdata(t)

		// Make the temp directory on the same physical drive, os.Rename doesn't work across drives and /tmp may be on another drive
		tmp, err := os.MkdirTemp("testdata", "magex")
		require.NoError(t, err, "could not create temp directory for test")
		defer os.RemoveAll(tmp)

		err = Move("testdata/copy/a/a1.txt", tmp)
		require.NoError(t, err, "Move file failed")

		assertFile(t, filepath.Join(tmp, "a1.txt"))
	})

	t.Run("overwrite directory should fail", func(t *testing.T) {
		defer resetTestdata(t)

		// Make the temp directory on the same physical drive, os.Rename doesn't work across drives and /tmp may be on another drive
		tmp, err := os.MkdirTemp("testdata", "magex")
		require.NoError(t, err, "could not create temp directory for test")
		defer os.RemoveAll(tmp)

		err = Move("testdata/copy/directory-conflict/a", tmp, MoveRecursive)
		require.NoError(t, err, "Setup failed")

		err = Move("testdata/copy/a/*", filepath.Join(tmp, "a"))
		require.Error(t, err, "Overwrite directory should have failed")
	})
}

func TestMove_MoveNoOverwrite(t *testing.T) {
	testcases := []struct {
		name         string
		opts         MoveOption
		wantContents string
	}{
		{name: "overwrite", opts: MoveDefault, wantContents: "a2.txt"},
		{name: "no overwrite", opts: MoveNoOverwrite, wantContents: "a1.txt"},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			defer resetTestdata(t)

			// Make the temp directory on the same physical drive, os.Rename doesn't work across drives and /tmp may be on another drive
			tmp, err := os.MkdirTemp("testdata", "magex")
			require.NoError(t, err, "could not create temp directory for test")
			defer os.RemoveAll(tmp)

			err = Move("testdata/copy/a/a1.txt", tmp)
			require.NoError(t, err, "Move a1.txt failed")

			err = Move("testdata/copy/a/a2.txt", filepath.Join(tmp, "a1.txt"), tc.opts)
			require.NoError(t, err, "Overwrite failed")

			gotContents, err := ioutil.ReadFile(filepath.Join(tmp, "a1.txt"))
			require.NoError(t, err, "could not read file")
			assert.Equal(t, tc.wantContents, string(gotContents), "invalid contents, want: %s, got: %s", tc.wantContents, gotContents)
		})
	}
}
