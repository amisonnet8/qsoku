package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRun_init(t *testing.T) {
	t.Run("creates an empty qsokufile", func(t *testing.T) {
		dir := t.TempDir()
		t.Chdir(dir)

		var stdout, stderr bytes.Buffer
		got := Run([]string{".init"}, nil, &stdout, &stderr)
		if got != 0 {
			t.Fatalf("Run(.init) = %d, want 0; stderr: %s", got, stderr.String())
		}

		data, err := os.ReadFile(filepath.Join(dir, "qsokufile")) //nolint:gosec // dir is t.TempDir()
		if err != nil {
			t.Fatalf("qsokufile was not created: %v", err)
		}
		if len(data) != 0 {
			t.Errorf("qsokufile content = %q, want empty", data)
		}
	})

	t.Run("errors if one already exists here", func(t *testing.T) {
		dir := t.TempDir()
		writeQsokufile(t, dir, "root: cd //\n")
		t.Chdir(dir)

		var stdout, stderr bytes.Buffer
		got := Run([]string{".init"}, nil, &stdout, &stderr)
		if got != 1 {
			t.Errorf("Run(.init) = %d, want 1", got)
		}
		if stderr.Len() == 0 {
			t.Error("Run(.init) wrote nothing to stderr")
		}
	})

	t.Run("a qsokufile in a parent does not block it", func(t *testing.T) {
		dir := t.TempDir()
		writeQsokufile(t, dir, "root: cd //\n")
		sub := filepath.Join(dir, "sub")
		if err := os.Mkdir(sub, 0o750); err != nil {
			t.Fatal(err)
		}
		t.Chdir(sub)

		var stdout, stderr bytes.Buffer
		if got := Run([]string{".init"}, nil, &stdout, &stderr); got != 0 {
			t.Fatalf("Run(.init) = %d, want 0; stderr: %s", got, stderr.String())
		}
		if _, err := os.Stat(filepath.Join(sub, "qsokufile")); err != nil {
			t.Errorf("qsokufile was not created in %s: %v", sub, err)
		}
	})

	t.Run("rejects extra arguments", func(t *testing.T) {
		t.Chdir(t.TempDir())

		var stdout, stderr bytes.Buffer
		if got := Run([]string{".init", "extra"}, nil, &stdout, &stderr); got != 2 {
			t.Errorf("Run(.init extra) = %d, want 2", got)
		}
	})
}
