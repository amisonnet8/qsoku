package cli

import (
	"bytes"
	"path/filepath"
	"testing"
)

func TestRun_rm(t *testing.T) {
	t.Run("removes the named line", func(t *testing.T) {
		dir := t.TempDir()
		writeQsokufile(t, dir, "root: cd //\nbuild: make\n")
		t.Chdir(dir)

		var stdout, stderr bytes.Buffer
		got := Run([]string{".rm", "build"}, nil, &stdout, &stderr)
		if got != 0 {
			t.Fatalf("Run(.rm) = %d, want 0; stderr: %s", got, stderr.String())
		}

		want := "root: cd //\n"
		if got := readTestFile(t, filepath.Join(dir, "qsokufile")); got != want {
			t.Errorf("file = %q, want %q", got, want)
		}
	})

	t.Run("name not defined", func(t *testing.T) {
		dir := t.TempDir()
		writeQsokufile(t, dir, "root: cd //\n")
		t.Chdir(dir)

		var stdout, stderr bytes.Buffer
		got := Run([]string{".rm", "missing"}, nil, &stdout, &stderr)
		if got != 2 {
			t.Errorf("Run(.rm missing) = %d, want 2", got)
		}
	})

	t.Run("no qsokufile found", func(t *testing.T) {
		t.Chdir(t.TempDir())

		var stdout, stderr bytes.Buffer
		if got := Run([]string{".rm", "x"}, nil, &stdout, &stderr); got != 1 {
			t.Errorf("Run(.rm x) = %d, want 1", got)
		}
	})

	t.Run("wrong number of arguments", func(t *testing.T) {
		dir := t.TempDir()
		writeQsokufile(t, dir, "root: cd //\n")
		t.Chdir(dir)

		var stdout, stderr bytes.Buffer
		if got := Run([]string{".rm"}, nil, &stdout, &stderr); got != 2 {
			t.Errorf("Run(.rm) = %d, want 2", got)
		}
	})
}
