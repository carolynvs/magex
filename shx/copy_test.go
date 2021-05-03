package shx

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("recursively copy directory into empty dest dir", func(t *testing.T) {
		tmp, err := ioutil.TempDir("", "magex")
		require.NoError(t, err, "could not create temp directory for test")
		defer os.RemoveAll(tmp)

		err = Copy("testdata/copy/a", tmp, CopyRecursive)
		require.NoError(t, err, "Copy into empty directory failed")

		assert.DirExists(t, filepath.Join(tmp, "a"))
		assertFile(t, filepath.Join(tmp, "a/a1.txt"))
		assertFile(t, filepath.Join(tmp, "a/a2.txt"))
		assert.DirExists(t, filepath.Join(tmp, "a/ab"))
		assertFile(t, filepath.Join(tmp, "a/ab/ab1.txt"))
		assertFile(t, filepath.Join(tmp, "a/ab/ab2.txt"))
	})

	t.Run("recursively copy directory into populated dest dir", func(t *testing.T) {
		tmp, err := ioutil.TempDir("", "magex")
		require.NoError(t, err, "could not create temp directory for test")
		defer os.RemoveAll(tmp)

		require.NoError(t, os.MkdirAll(filepath.Join(tmp, "a"), 0755))

		err = Copy("testdata/copy/a", tmp, CopyRecursive)
		require.NoError(t, err, "Copy into directory with same directory name")

		assert.DirExists(t, filepath.Join(tmp, "a"))
		assertFile(t, filepath.Join(tmp, "a/a1.txt"))
		assertFile(t, filepath.Join(tmp, "a/a2.txt"))
		assert.DirExists(t, filepath.Join(tmp, "a/ab"))
		assertFile(t, filepath.Join(tmp, "a/ab/ab1.txt"))
		assertFile(t, filepath.Join(tmp, "a/ab/ab2.txt"))
	})

	t.Run("copy glob", func(t *testing.T) {
		tmp, err := ioutil.TempDir("", "magex")
		require.NoError(t, err, "could not create temp directory for test")
		defer os.RemoveAll(tmp)

		err = Copy("testdata/copy/a/*.txt", tmp)
		require.NoError(t, err, "Copy into empty directory failed")

		assertFile(t, filepath.Join(tmp, "a1.txt"))
		assertFile(t, filepath.Join(tmp, "a2.txt"))
		assert.NoDirExists(t, filepath.Join(tmp, "a"))
		assert.NoDirExists(t, filepath.Join(tmp, "ab"))
	})

	t.Run("missing parent dir", func(t *testing.T) {
		tmp, err := ioutil.TempDir("", "magex")
		require.NoError(t, err, "could not create temp directory for test")
		defer os.RemoveAll(tmp)

		err = Copy("testdata/copy/a", filepath.Join(tmp, "missing-parent/dir"))
		require.Error(t, err)
	})

	t.Run("missing src", func(t *testing.T) {
		tmp, err := ioutil.TempDir("", "magex")
		require.NoError(t, err, "could not create temp directory for test")
		defer os.RemoveAll(tmp)

		err = Copy("testdata/missing-src", tmp)
		require.Error(t, err)
		require.Contains(t, err.Error(), "no such file or directory")
	})

	t.Run("recursively copy directory to new name", func(t *testing.T) {
		tmp, err := ioutil.TempDir("", "magex")
		require.NoError(t, err, "could not create temp directory for test")
		defer os.RemoveAll(tmp)

		dest := filepath.Join(tmp, "dest")
		err = Copy("testdata/copy/a", dest, CopyRecursive)
		require.NoError(t, err, "Copy into empty directory failed")

		assert.DirExists(t, dest)
		assertFile(t, filepath.Join(dest, "a1.txt"))
		assertFile(t, filepath.Join(dest, "a2.txt"))
		assert.DirExists(t, filepath.Join(dest, "ab"))
		assertFile(t, filepath.Join(dest, "ab/ab1.txt"))
		assertFile(t, filepath.Join(dest, "ab/ab2.txt"))
	})

	t.Run("recursively merge dest dir", func(t *testing.T) {
		tmp, err := ioutil.TempDir("", "magex")
		require.NoError(t, err, "could not create temp directory for test")
		defer os.RemoveAll(tmp)

		err = Copy("testdata/copy/partial-dest", tmp, CopyRecursive)
		require.NoError(t, err, "Copy partial destination failed")

		err = Copy("testdata/copy/a", tmp, CopyRecursive)
		require.NoError(t, err, "Merge into non-empty destination failed")

		assert.DirExists(t, filepath.Join(tmp, "a"))
		assertFile(t, filepath.Join(tmp, "a/a1.txt"))
		assertFile(t, filepath.Join(tmp, "a/a2.txt"))
		assert.DirExists(t, filepath.Join(tmp, "a/ab"))
		assertFile(t, filepath.Join(tmp, "a/ab/ab1.txt"))
		assertFile(t, filepath.Join(tmp, "a/ab/ab2.txt"))
	})

	t.Run("copy file into empty directory", func(t *testing.T) {
		tmp, err := ioutil.TempDir("", "magex")
		require.NoError(t, err, "could not create temp directory for test")
		defer os.RemoveAll(tmp)

		err = Copy("testdata/copy/a/a1.txt", tmp)
		require.NoError(t, err, "Copy file failed")

		assertFile(t, filepath.Join(tmp, "a1.txt"))
	})

	t.Run("overwrite directory should fail", func(t *testing.T) {
		tmp, err := ioutil.TempDir("", "magex")
		require.NoError(t, err, "could not create temp directory for test")
		defer os.RemoveAll(tmp)

		err = Copy("testdata/copy/directory-conflict/a", tmp, CopyRecursive)
		require.NoError(t, err, "Setup failed")

		err = Copy("testdata/copy/a/*", filepath.Join(tmp, "a"))
		require.Error(t, err, "Overwrite directory should have failed")
	})
}

func TestCopy_CopyNoRecursive(t *testing.T) {
	tmp, err := ioutil.TempDir("", "magex")
	require.NoError(t, err, "could not create temp directory for test")
	defer os.RemoveAll(tmp)

	err = Copy("testdata/copy/a", tmp)
	require.NoError(t, err, "Copy into empty directory failed")

	assert.DirExists(t, filepath.Join(tmp, "a"))
	assert.NoFileExists(t, filepath.Join(tmp, "a/a1.txt"))
	assert.NoFileExists(t, filepath.Join(tmp, "a/a2.txt"))
	assert.NoDirExists(t, filepath.Join(tmp, "a/ab"))
}

func TestCopy_CopyNoOverwrite(t *testing.T) {
	testcases := []struct {
		name         string
		opts         CopyOption
		wantContents string
	}{
		{name: "overwrite", opts: CopyDefault, wantContents: "a2.txt"},
		{name: "no overwrite", opts: CopyNoOverwrite, wantContents: "a1.txt"},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tmp, err := ioutil.TempDir("", "magex")
			require.NoError(t, err, "could not create temp directory for test")
			defer os.RemoveAll(tmp)

			err = Copy("testdata/copy/a/a1.txt", tmp)
			require.NoError(t, err, "Copy a1.txt failed")

			err = Copy("testdata/copy/a/a2.txt", filepath.Join(tmp, "a1.txt"), tc.opts)
			require.NoError(t, err, "Overwrite failed")

			gotContents, err := ioutil.ReadFile(filepath.Join(tmp, "a1.txt"))
			require.NoError(t, err, "could not read file")
			assert.Equal(t, tc.wantContents, string(gotContents), "invalid contents, want: %s, got: %s", tc.wantContents, gotContents)
		})
	}
}

func assertFile(t *testing.T, f string) {
	gotContents, err := ioutil.ReadFile(f)
	require.NoErrorf(t, err, "could not read file %s", f)

	wantContents := filepath.Base(f)
	assert.Equal(t, wantContents, string(gotContents), "invalid contents for %s, want: %s, got: %s", f, wantContents, gotContents)
}
