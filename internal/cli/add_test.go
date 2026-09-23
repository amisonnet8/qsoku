package cli

import (
	"bytes"
	"path/filepath"
	"testing"
)

func TestRun_add(t *testing.T) {
	t.Run("appends a new name", func(t *testing.T) {
		dir := t.TempDir()
		writeQsokufile(t, dir, "root: cd //\n")
		t.Chdir(dir)

		var stdout, stderr bytes.Buffer
		got := Run([]string{".add", "build", "make build"}, nil, &stdout, &stderr)
		if got != 0 {
			t.Fatalf("Run(.add) = %d, want 0; stderr: %s", got, stderr.String())
		}

		want := "root: cd //\nbuild: make build\n"
		if got := readTestFile(t, filepath.Join(dir, "qsokufile")); got != want {
			t.Errorf("file = %q, want %q", got, want)
		}
	})

	t.Run("replaces an existing name in place", func(t *testing.T) {
		dir := t.TempDir()
		writeQsokufile(t, dir, "# comment\nbuild: old\nroot: cd //\n")
		t.Chdir(dir)

		var stdout, stderr bytes.Buffer
		got := Run([]string{".add", "build", "new"}, nil, &stdout, &stderr)
		if got != 0 {
			t.Fatalf("Run(.add) = %d, want 0; stderr: %s", got, stderr.String())
		}

		want := "# comment\nbuild: new\nroot: cd //\n"
		if got := readTestFile(t, filepath.Join(dir, "qsokufile")); got != want {
			t.Errorf("file = %q, want %q", got, want)
		}
	})

	t.Run("rejects an invalid name", func(t *testing.T) {
		dir := t.TempDir()
		writeQsokufile(t, dir, "root: cd //\n")
		t.Chdir(dir)

		var stdout, stderr bytes.Buffer
		got := Run([]string{".add", ".init", "echo no"}, nil, &stdout, &stderr)
		if got != 2 {
			t.Errorf("Run(.add .init ...) = %d, want 2", got)
		}
	})

	t.Run("no qsokufile found", func(t *testing.T) {
		t.Chdir(t.TempDir())

		var stdout, stderr bytes.Buffer
		if got := Run([]string{".add", "x", "y"}, nil, &stdout, &stderr); got != 1 {
			t.Errorf("Run(.add) = %d, want 1", got)
		}
	})

	t.Run("wrong number of arguments", func(t *testing.T) {
		dir := t.TempDir()
		writeQsokufile(t, dir, "root: cd //\n")
		t.Chdir(dir)

		var stdout, stderr bytes.Buffer
		if got := Run([]string{".add", "onlyname"}, nil, &stdout, &stderr); got != 2 {
			t.Errorf("Run(.add onlyname) = %d, want 2", got)
		}
	})
}
