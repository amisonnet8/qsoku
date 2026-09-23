package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRun_where(t *testing.T) {
	t.Run("prints the qsokufile's directory", func(t *testing.T) {
		dir := t.TempDir()
		writeQsokufile(t, dir, "root: cd //\n")
		sub := filepath.Join(dir, "sub")
		if err := os.Mkdir(sub, 0o750); err != nil {
			t.Fatal(err)
		}
		t.Chdir(sub)

		var stdout, stderr bytes.Buffer
		got := Run([]string{".where"}, nil, &stdout, &stderr)
		if got != 0 {
			t.Fatalf("Run(.where) = %d, want 0; stderr: %s", got, stderr.String())
		}
		if want := dir + "\n"; stdout.String() != want {
			t.Errorf("stdout = %q, want %q", stdout.String(), want)
		}
	})

	t.Run("works even with a qsokufile that fails to parse", func(t *testing.T) {
		dir := t.TempDir()
		writeQsokufile(t, dir, "not valid\n")
		t.Chdir(dir)

		var stdout, stderr bytes.Buffer
		got := Run([]string{".where"}, nil, &stdout, &stderr)
		if got != 0 {
			t.Fatalf("Run(.where) = %d, want 0; stderr: %s", got, stderr.String())
		}
		if want := dir + "\n"; stdout.String() != want {
			t.Errorf("stdout = %q, want %q", stdout.String(), want)
		}
	})

	t.Run("no qsokufile found", func(t *testing.T) {
		t.Chdir(t.TempDir())

		var stdout, stderr bytes.Buffer
		if got := Run([]string{".where"}, nil, &stdout, &stderr); got != 1 {
			t.Errorf("Run(.where) = %d, want 1", got)
		}
	})
}
