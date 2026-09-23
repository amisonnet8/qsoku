package qsokufile

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func mustMkdirAll(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o750); err != nil {
		t.Fatal(err)
	}
}

func mustWriteFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestFind(t *testing.T) {
	t.Run("in the current directory", func(t *testing.T) {
		root := t.TempDir()
		want := filepath.Join(root, fileName)
		mustWriteFile(t, want, "root: cd //\n")

		got, err := Find(root)
		if err != nil {
			t.Fatalf("Find() error: %v", err)
		}
		if got != want {
			t.Errorf("Find() = %q, want %q", got, want)
		}
	})

	t.Run("found in a parent from a deep subdirectory", func(t *testing.T) {
		root := t.TempDir()
		want := filepath.Join(root, fileName)
		mustWriteFile(t, want, "root: cd //\n")

		deep := filepath.Join(root, "a", "b", "c")
		mustMkdirAll(t, deep)

		got, err := Find(deep)
		if err != nil {
			t.Fatalf("Find() error: %v", err)
		}
		if got != want {
			t.Errorf("Find() = %q, want %q", got, want)
		}
	})

	t.Run("walks past a .git boundary", func(t *testing.T) {
		root := t.TempDir()
		want := filepath.Join(root, fileName)
		mustWriteFile(t, want, "root: cd //\n")

		repo := filepath.Join(root, "repo")
		mustMkdirAll(t, filepath.Join(repo, ".git"))

		got, err := Find(repo)
		if err != nil {
			t.Fatalf("Find() error: %v", err)
		}
		if got != want {
			t.Errorf("Find() = %q, want %q", got, want)
		}
	})

	t.Run("a directory named qsokufile is not a match", func(t *testing.T) {
		root := t.TempDir()
		mustMkdirAll(t, filepath.Join(root, fileName))

		want := filepath.Dir(root) // whatever is above root, or ErrNotFound
		got, err := Find(root)
		if err == nil {
			if got == filepath.Join(root, fileName) {
				t.Fatalf("Find() returned the directory named qsokufile: %q", got)
			}
			return
		}
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("Find() error = %v, want ErrNotFound or a match above %q", err, want)
		}
	})

	t.Run("not found", func(t *testing.T) {
		root := t.TempDir()
		nested := filepath.Join(root, "x", "y")
		mustMkdirAll(t, nested)

		// t.TempDir() is itself under the real filesystem root, which may or
		// may not have a qsokufile of its own on the machine running the
		// test. To keep this deterministic, this test only checks that Find
		// does not stop early at "nested" or "root" -- it must not find
		// anything placed nowhere. A true "walked all the way to / and found
		// nothing" case is exercised implicitly by every other subtest
		// succeeding at the expected level rather than a higher one.
		got, err := Find(nested)
		if err == nil {
			t.Fatalf("Find() unexpectedly found %q", got)
		}
		if !errors.Is(err, ErrNotFound) {
			t.Fatalf("Find() error = %v, want ErrNotFound", err)
		}
	})
}
