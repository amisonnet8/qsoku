package cli

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// fakeEditor writes a shell script that appends a marker line to whatever
// file it's given as $1, and returns its absolute path in forward-slash
// form. runEdit splices $EDITOR into the sh script text unquoted (so a
// multi-word value like "code --wait" still splits into arguments --
// docs/reference/cli.md ".edit"); on Windows, sh (Git for Windows' MSYS
// sh.exe) would otherwise parse t.TempDir()'s native backslashes as
// escapes and mangle the path. A real Windows $EDITOR pointing at an
// absolute path needs the same forward-slash form for the same reason --
// this fixture matches that, rather than working around a test-only
// artifact.
func fakeEditor(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "fake-editor.sh")
	script := "#!/bin/sh\necho edited-by-fake-editor >> \"$1\"\n"
	if err := os.WriteFile(path, []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}
	return filepath.ToSlash(path)
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

	// $1 is the path itself, so a directory with a space in its name still
	// reaches the editor as one argument.
	t.Run("a path with a space reaches the editor whole", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "a b")
		if err := os.Mkdir(dir, 0o700); err != nil {
			t.Fatal(err)
		}
		writeQsokufile(t, dir, "root: cd //\n")
		t.Chdir(dir)
		t.Setenv("EDITOR", fakeEditor(t))

		var stdout, stderr bytes.Buffer
		if got := Run([]string{".edit"}, nil, &stdout, &stderr); got != 0 {
			t.Fatalf("Run(.edit) = %d, want 0; stderr: %s", got, stderr.String())
		}
		want := "root: cd //\nedited-by-fake-editor\n"
		if got := readTestFile(t, filepath.Join(dir, "qsokufile")); got != want {
			t.Errorf("file = %q, want %q", got, want)
		}
	})

	// PATH holds only a fake nano and sh: a real editor is never started,
	// and one other than nano is not found.
	t.Run("no EDITOR falls back to nano", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("Git for Windows' sh.exe needs the DLLs beside it, so it cannot be linked into a PATH of its own")
		}
		sh, err := exec.LookPath("sh")
		if err != nil {
			t.Fatal(err)
		}
		bin := t.TempDir()
		if err := os.Symlink(sh, filepath.Join(bin, "sh")); err != nil {
			t.Fatal(err)
		}
		script := "#!" + sh + "\necho edited-by-fake-nano >> \"$1\"\n"
		if err := os.WriteFile(filepath.Join(bin, "nano"), []byte(script), 0o700); err != nil {
			t.Fatal(err)
		}
		dir := t.TempDir()
		writeQsokufile(t, dir, "root: cd //\n")
		t.Chdir(dir)
		t.Setenv("EDITOR", "")
		t.Setenv("PATH", bin)

		var stdout, stderr bytes.Buffer
		if got := Run([]string{".edit"}, nil, &stdout, &stderr); got != 0 {
			t.Fatalf("Run(.edit) = %d, want 0; stderr: %s", got, stderr.String())
		}
		want := "root: cd //\nedited-by-fake-nano\n"
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
