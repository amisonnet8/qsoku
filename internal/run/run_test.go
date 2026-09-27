package run

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestExecute_exitCode(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    int
	}{
		{"success", "exit 0", 0},
		{"failure", "exit 7", 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			code, err := Execute(tt.command, t.TempDir(), nil, nil, &stdout, &stderr)
			if err != nil {
				t.Fatalf("Execute() error: %v", err)
			}
			if code != tt.want {
				t.Errorf("Execute() code = %d, want %d", code, tt.want)
			}
		})
	}
}

func TestExecute_stdout(t *testing.T) {
	var stdout, stderr bytes.Buffer
	_, err := Execute("echo hello", t.TempDir(), nil, nil, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
	if got := stdout.String(); got != "hello\n" {
		t.Errorf("stdout = %q, want %q", got, "hello\n")
	}
}

func TestExecute_root(t *testing.T) {
	root := t.TempDir()
	var stdout, stderr bytes.Buffer
	_, err := Execute("echo $QSOKU_ROOT", root, nil, nil, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
	if got := strings.TrimSpace(stdout.String()); got != root {
		t.Errorf("QSOKU_ROOT = %q, want %q", got, root)
	}
}

// canonicalPath resolves p to the form the OS itself reports for it once a
// real cwd lookup is involved (filepath.EvalSymlinks resolves symlinks on
// Linux/macOS -- notably macOS's /tmp -> /private/tmp -- and, on Windows,
// also normalizes a short 8.3-style path segment such as "RUNNER~1" to its
// real long name, which "pwd -W" reports rather than preserving the form
// t.TempDir() happened to return). Comparing against this instead of the
// raw t.TempDir() value directly is what lets the expectations below match
// what actually comes back on every platform.
func canonicalPath(t *testing.T, p string) string {
	t.Helper()
	resolved, err := filepath.EvalSymlinks(p)
	if err != nil {
		t.Fatal(err)
	}
	return resolved
}

func TestExecute_cwdHandoff(t *testing.T) {
	target := canonicalPath(t, t.TempDir())
	start := canonicalPath(t, t.TempDir())
	// Windows: these paths are spliced into the sh script text below
	// unquoted, and sh (Git for Windows' MSYS sh.exe) parses a backslash
	// outside quotes as an escape, stripping it -- t.TempDir()'s native
	// backslash form would arrive at "cd" mangled (C:\Users\... becomes
	// C:Users...). The trailer this package itself appends on Windows
	// (trailerWindows) reports the location via "pwd -W" in the same
	// forward-slash form regardless, so both the command text and the
	// expected values below are normalized with it. This is a no-op on
	// Linux/macOS, where the path already uses '/'.
	targetSh := filepath.ToSlash(target)
	startSh := filepath.ToSlash(start)

	// Windows paths are case-insensitive; whether "pwd -W" happens to
	// report the drive letter in the same case Go's own t.TempDir() used is
	// not something this package controls, so the comparison below ignores
	// case there. Elsewhere, exact case is expected as before.
	pathsEqual := func(a, b string) bool {
		if runtime.GOOS == "windows" {
			return strings.EqualFold(a, b)
		}
		return a == b
	}

	readCwdFile := func(t *testing.T, command string, wantCode int) string {
		t.Helper()
		cwdFile := filepath.Join(t.TempDir(), "cwd")
		t.Setenv("QSOKU_CWD_FILE", cwdFile)
		t.Chdir(start)

		var stdout, stderr bytes.Buffer
		code, err := Execute(command, "/root-unused", nil, nil, &stdout, &stderr)
		if err != nil {
			t.Fatalf("Execute() error: %v", err)
		}
		if code != wantCode {
			t.Fatalf("Execute() code = %d, want %d", code, wantCode)
		}
		if stderr.Len() != 0 {
			t.Fatalf("stderr = %q, want empty", stderr.String())
		}
		data, err := os.ReadFile(cwdFile) //nolint:gosec // cwdFile is a path this test built itself under t.TempDir()
		if err != nil {
			t.Fatalf("reading cwd file: %v", err)
		}
		return strings.TrimSpace(string(data))
	}

	t.Run("no parens, success: location is brought back", func(t *testing.T) {
		got := readCwdFile(t, "cd "+targetSh, 0)
		if !pathsEqual(got, targetSh) {
			t.Errorf("cwd = %q, want %q", got, targetSh)
		}
	})

	// A qsokufile command never contains a literal top-level "exit": qsoku
	// appends its own, and running one early would skip the trailer
	// entirely (exit terminates the shell immediately, before anything
	// appended after it can run). A command "fails" the way a real
	// qsokufile entry would -- the natural exit status of its own last
	// command, here a nested sh -c that itself exits 3.
	t.Run("no parens, failure: location is still brought back", func(t *testing.T) {
		got := readCwdFile(t, "cd "+targetSh+`; sh -c "exit 3"`, 3)
		if !pathsEqual(got, targetSh) {
			t.Errorf("cwd = %q, want %q", got, targetSh)
		}
	})

	t.Run("parens: a subshell cd does not affect the outer location", func(t *testing.T) {
		got := readCwdFile(t, "(cd "+targetSh+`; sh -c "exit 5")`, 5)
		if !pathsEqual(got, startSh) {
			t.Errorf("cwd = %q, want %q (unchanged)", got, startSh)
		}
	})
}

func TestExecute_noCwdFile(t *testing.T) {
	// Deliberately not calling t.Setenv("QSOKU_CWD_FILE", ...): it must be
	// absent from this test's own environment for this case to be
	// meaningful.
	if _, ok := os.LookupEnv("QSOKU_CWD_FILE"); ok {
		t.Fatal("QSOKU_CWD_FILE is unexpectedly set in the test environment")
	}

	var stdout, stderr bytes.Buffer
	code, err := Execute("exit 0", t.TempDir(), nil, nil, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
	if code != 0 {
		t.Errorf("code = %d, want 0", code)
	}
	if stderr.Len() != 0 {
		t.Errorf("stderr = %q, want empty (no redirect-to-empty-filename error)", stderr.String())
	}
}

func TestExecute_signal(t *testing.T) {
	if runtime.GOOS == "windows" {
		// sh on Windows is Git for Windows' MSYS sh.exe; whether "kill -TERM
		// $$" against it still surfaces as 128+SIGTERM through Go's
		// syscall.WaitStatus (which on Windows has no real signal concept)
		// has not been confirmed on real Windows CI yet (see todo
		// a251592e2c). Skip rather than assert something unverified.
		t.Skip("signal delivery to sh on Windows is not confirmed yet")
	}
	var stdout, stderr bytes.Buffer
	code, err := Execute("kill -TERM $$; sleep 1", t.TempDir(), nil, nil, &stdout, &stderr)
	if err != nil {
		t.Fatalf("Execute() error: %v", err)
	}
	const sigterm = 15
	if want := 128 + sigterm; code != want {
		t.Errorf("code = %d, want %d (128+SIGTERM)", code, want)
	}
}

func TestExecute_shNotFound(t *testing.T) {
	emptyPath := t.TempDir()
	t.Setenv("PATH", emptyPath)

	var stdout, stderr bytes.Buffer
	_, err := Execute("exit 0", t.TempDir(), nil, nil, &stdout, &stderr)
	if err == nil {
		t.Fatal("Execute() error = nil, want an error (sh not found)")
	}
}
