//go:build e2e

package e2e

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// runBinary runs the built qsoku binary in dir, with extraEnv appended to a
// clean PATH-only environment, and returns its standard output, standard
// error and exit code.
func runBinary(t *testing.T, dir string, extraEnv []string, args ...string) (stdout, stderr string, code int) {
	t.Helper()
	cmd := exec.Command(binary, args...) //nolint:gosec // binary is the qsoku this package's TestMain built, not user input
	cmd.Dir = dir
	cmd.Env = append([]string{"PATH=" + os.Getenv("PATH")}, extraEnv...)
	var out, errOut strings.Builder
	cmd.Stdout, cmd.Stderr = &out, &errOut
	err := cmd.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			t.Fatalf("cannot run qsoku: %v", err)
		}
		code = exitErr.ExitCode()
	}
	return out.String(), errOut.String(), code
}

// resolvedTempDir returns t.TempDir(), resolved through any symlinks
// (macOS puts /tmp under /private/var), so comparisons against pwd's output
// match exactly.
func resolvedTempDir(t *testing.T) string {
	t.Helper()
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func writeQsokufile(t *testing.T, dir, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, "qsokufile"), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

// TestRunPassesThroughExitCode checks that once sh starts running the
// entry's command, qsoku's own exit codes stop applying (docs/reference/cli.md
// "Exit codes").
func TestRunPassesThroughExitCode(t *testing.T) {
	dir := resolvedTempDir(t)
	writeQsokufile(t, dir, "seven: sh -c \"exit 7\"\n")

	_, _, code := runBinary(t, dir, nil, "seven")
	if code != 7 {
		t.Errorf("code = %d, want 7", code)
	}
}

// TestRunSignaledCommand checks that a command killed by a signal reports
// 128+n, the same convention mtqg itself uses.
func TestRunSignaledCommand(t *testing.T) {
	dir := resolvedTempDir(t)
	writeQsokufile(t, dir, "term: kill -TERM $$\n")

	_, _, code := runBinary(t, dir, nil, "term")
	if code != 128+15 { // SIGTERM
		t.Errorf("code = %d, want %d", code, 128+15)
	}
}

// TestRunNoQsokufile checks the exit code and message when no qsokufile is
// found anywhere above the current directory.
func TestRunNoQsokufile(t *testing.T) {
	dir := resolvedTempDir(t)

	_, stderr, code := runBinary(t, dir, nil, "anything")
	if code != 1 {
		t.Errorf("code = %d, want 1", code)
	}
	if !strings.Contains(stderr, "qsoku .init") {
		t.Errorf("stderr = %q, want it to suggest qsoku .init", stderr)
	}
}

// TestRunParseError checks that a qsokufile line with no ':' is reported
// with its line number and stops qsoku before it runs anything.
func TestRunParseError(t *testing.T) {
	dir := resolvedTempDir(t)
	writeQsokufile(t, dir, "ok: echo hi\nno colon here\n")

	_, stderr, code := runBinary(t, dir, nil, "ok")
	if code != 1 {
		t.Errorf("code = %d, want 1", code)
	}
	if !strings.Contains(stderr, "line 2:") {
		t.Errorf("stderr = %q, want it to mention line 2", stderr)
	}
}

// TestRunUndefinedNameAndUnknownManagementCommand checks the command-line
// error case (exit 2), for both an undefined qsokufile name and an unknown
// ".foo" management command.
func TestRunUndefinedNameAndUnknownManagementCommand(t *testing.T) {
	dir := resolvedTempDir(t)
	writeQsokufile(t, dir, "ok: echo hi\n")

	for _, name := range []string{"nope", ".foo"} {
		_, _, code := runBinary(t, dir, nil, name)
		if code != 2 {
			t.Errorf("qsoku %s: code = %d, want 2", name, code)
		}
	}
}

// TestRunWrongArgumentCount checks that a management command given the
// wrong number of arguments is a command-line error (exit 2).
func TestRunWrongArgumentCount(t *testing.T) {
	dir := resolvedTempDir(t)
	writeQsokufile(t, dir, "ok: echo hi\n")

	_, _, code := runBinary(t, dir, nil, ".where", "extra")
	if code != 2 {
		t.Errorf("code = %d, want 2", code)
	}
}

// TestRunArgumentSubstitution checks that "//" is substituted in an
// argument typed by the user, the same way it is in the entry's own command
// (docs/reference/qsokufile.md "Arguments typed by the user").
func TestRunArgumentSubstitution(t *testing.T) {
	dir := resolvedTempDir(t)
	if err := os.Mkdir(filepath.Join(dir, "src"), 0o750); err != nil {
		t.Fatal(err)
	}
	writeQsokufile(t, dir, "show: echo \"$1\"\n")

	stdout, _, code := runBinary(t, dir, nil, "show", "//src")
	if code != 0 {
		t.Fatalf("code = %d, want 0 (stdout=%q)", code, stdout)
	}
	want := filepath.Join(dir, "src") + "\n"
	if stdout != want {
		t.Errorf("stdout = %q, want %q", stdout, want)
	}
}

// TestRunCwdHandoffDoesNotMixWithStdout checks that the pwd line qsoku
// appends for the working-directory handoff goes only to QSOKU_CWD_FILE,
// never into the command's own standard output.
func TestRunCwdHandoffDoesNotMixWithStdout(t *testing.T) {
	dir := resolvedTempDir(t)
	if err := os.Mkdir(filepath.Join(dir, "sub"), 0o750); err != nil {
		t.Fatal(err)
	}
	writeQsokufile(t, dir, "root: echo hello; cd //sub\n")

	cwdFile := filepath.Join(t.TempDir(), "cwd")
	stdout, _, code := runBinary(t, dir, []string{"QSOKU_CWD_FILE=" + cwdFile}, "root")
	if code != 0 {
		t.Fatalf("code = %d, want 0", code)
	}
	if stdout != "hello\n" {
		t.Errorf("stdout = %q, want %q", stdout, "hello\n")
	}
	got, err := os.ReadFile(cwdFile) //nolint:gosec // cwdFile is this test's own path, built from t.TempDir() a few lines above
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(dir, "sub") + "\n"
	if string(got) != want {
		t.Errorf("cwd file = %q, want %q", got, want)
	}
}

// TestRunCwdHandoffNotWrittenOnQsokuOwnError checks that qsoku's own errors
// (before the command runs) leave QSOKU_CWD_FILE untouched, so the shell
// function knows not to cd.
func TestRunCwdHandoffNotWrittenOnQsokuOwnError(t *testing.T) {
	dir := resolvedTempDir(t)
	writeQsokufile(t, dir, "ok: echo hi\n")

	cwdFile := filepath.Join(t.TempDir(), "cwd")
	_, _, code := runBinary(t, dir, []string{"QSOKU_CWD_FILE=" + cwdFile}, "nope")
	if code != 2 {
		t.Fatalf("code = %d, want 2", code)
	}
	if _, err := os.Stat(cwdFile); !os.IsNotExist(err) {
		t.Errorf("QSOKU_CWD_FILE was written on qsoku's own error: err = %v", err)
	}
}

// TestDocsExamplesQsokufilesParse checks that every qsokufile under
// docs/examples/ (the ones the README and tour link to, meant to be copied
// into a real repository) actually parses: qsoku .list must succeed run
// from that directory. It does not run any of their entries (some call
// real tools like npm or go that this repository does not depend on).
func TestDocsExamplesQsokufilesParse(t *testing.T) {
	dirs, err := filepath.Glob(filepath.Join("..", "docs", "examples", "*", "qsokufile"))
	if err != nil {
		t.Fatal(err)
	}
	if len(dirs) == 0 {
		t.Fatal("no docs/examples/*/qsokufile found")
	}
	for _, qf := range dirs {
		dir := filepath.Dir(qf)
		t.Run(filepath.Base(dir), func(t *testing.T) {
			_, stderr, code := runBinary(t, dir, nil, ".list")
			if code != 0 {
				t.Errorf(".list in %s: code = %d, stderr = %q", dir, code, stderr)
			}
		})
	}
}

// TestRootQsokufileParses checks that this repository's own qsokufile (a
// worked sample, not used for qsoku's own development — see
// docs/design/history.md) parses: qsoku .list must succeed run from the
// repository root. It does not run any of the entries.
func TestRootQsokufileParses(t *testing.T) {
	_, stderr, code := runBinary(t, "..", nil, ".list")
	if code != 0 {
		t.Errorf(".list in repository root: code = %d, stderr = %q", code, stderr)
	}
}
