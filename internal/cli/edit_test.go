package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// fakeEditor writes a shell script that appends a marker line to whatever
// file it's given as $1, and returns its absolute path.
func fakeEditor(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "fake-editor.sh")
	script := "#!/bin/sh\necho edited-by-fake-editor >> \"$1\"\n"
	if err := os.WriteFile(path, []byte(script), 0o700); err != nil { //nolint:gosec // an executable test fixture needs the execute bit
		t.Fatal(err)
	}
	return path
}

func TestRun_edit(t *testing.T) {
	// The subtest name deliberately avoids a literal '$': t.TempDir()
	// embeds the (sanitized) test name verbatim in a real directory path,
	// and that path is then spliced into sh script text unescaped (the
	// same way a real $EDITOR value is, docs/reference/cli.md ".edit") --
	// a literal "$EDITOR" in the name would itself get shell-expanded.
	t.Run("runs the editor on the qsokufile in use", func(t *testing.T) {
		dir := t.TempDir()
		writeQsokufile(t, dir, "root: cd //\n")
		t.Chdir(dir)
		t.Setenv("EDITOR", fakeEditor(t))

		var stdout, stderr bytes.Buffer
		got := Run([]string{".edit"}, nil, &stdout, &stderr)
		if got != 0 {
			t.Fatalf("Run(.edit) = %d, want 0; stderr: %s", got, stderr.String())
		}

		want := "root: cd //\nedited-by-fake-editor\n"
		if got := readTestFile(t, filepath.Join(dir, "qsokufile")); got != want {
			t.Errorf("file = %q, want %q", got, want)
		}
	})

	t.Run("no EDITOR and no nano on PATH", func(t *testing.T) {
		dir := t.TempDir()
		writeQsokufile(t, dir, "root: cd //\n")
		t.Chdir(dir)
		t.Setenv("EDITOR", "")
		t.Setenv("PATH", t.TempDir())

		var stdout, stderr bytes.Buffer
		got := Run([]string{".edit"}, nil, &stdout, &stderr)
		if got != 1 {
			t.Errorf("Run(.edit) = %d, want 1", got)
		}
		if stderr.Len() == 0 {
			t.Error("Run(.edit) wrote nothing to stderr")
		}
	})

	t.Run("no qsokufile found", func(t *testing.T) {
		t.Chdir(t.TempDir())
		t.Setenv("EDITOR", fakeEditor(t))

		var stdout, stderr bytes.Buffer
		if got := Run([]string{".edit"}, nil, &stdout, &stderr); got != 1 {
			t.Errorf("Run(.edit) = %d, want 1", got)
		}
	})

	t.Run("rejects extra arguments", func(t *testing.T) {
		t.Chdir(t.TempDir())

		var stdout, stderr bytes.Buffer
		if got := Run([]string{".edit", "extra"}, nil, &stdout, &stderr); got != 2 {
			t.Errorf("Run(.edit extra) = %d, want 2", got)
		}
	})
}
