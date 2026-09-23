package run

import (
	"bytes"
	"os"
	"path/filepath"
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

func TestExecute_cwdHandoff(t *testing.T) {
	target := t.TempDir()
	start := t.TempDir()

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
		got := readCwdFile(t, "cd "+target, 0)
		if got != target {
			t.Errorf("cwd = %q, want %q", got, target)
		}
	})

	// A qsokufile command never contains a literal top-level "exit": qsoku
	// appends its own, and running one early would skip the trailer
	// entirely (exit terminates the shell immediately, before anything
	// appended after it can run). A command "fails" the way a real
	// qsokufile entry would -- the natural exit status of its own last
	// command, here a nested sh -c that itself exits 3.
	t.Run("no parens, failure: location is still brought back", func(t *testing.T) {
		got := readCwdFile(t, "cd "+target+`; sh -c "exit 3"`, 3)
		if got != target {
			t.Errorf("cwd = %q, want %q", got, target)
		}
	})

	t.Run("parens: a subshell cd does not affect the outer location", func(t *testing.T) {
		got := readCwdFile(t, "(cd "+target+`; sh -c "exit 5")`, 5)
		if got != start {
			t.Errorf("cwd = %q, want %q (unchanged)", got, start)
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
