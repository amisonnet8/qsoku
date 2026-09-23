package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestRun_version(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if got := Run([]string{".version"}, nil, &stdout, &stderr); got != 0 {
		t.Errorf("Run(.version) = %d, want 0", got)
	}
	if stdout.Len() == 0 {
		t.Error("Run(.version) wrote nothing to stdout")
	}
}

func TestRun_shellIsNotImplementedYet(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if got := Run([]string{".shell", "bash"}, nil, &stdout, &stderr); got != 1 {
		t.Errorf("Run(.shell bash) = %d, want 1", got)
	}
}

func TestRun_unknownDotCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if got := Run([]string{".foo"}, nil, &stdout, &stderr); got != 2 {
		t.Errorf("Run(.foo) = %d, want 2", got)
	}
	if stderr.Len() == 0 {
		t.Error("Run(.foo) wrote nothing to stderr")
	}
}

func TestRun_name(t *testing.T) {
	dir := t.TempDir()
	writeQsokufile(t, dir, "hello: echo hi\nfail: exit 7\n")
	t.Chdir(dir)

	t.Run("runs the entry's command", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := Run([]string{"hello"}, nil, &stdout, &stderr)
		if got != 0 {
			t.Fatalf("Run(hello) = %d, want 0; stderr: %s", got, stderr.String())
		}
		if want := "hi\n"; stdout.String() != want {
			t.Errorf("stdout = %q, want %q", stdout.String(), want)
		}
	})

	t.Run("passes the command's own exit code through", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := Run([]string{"fail"}, nil, &stdout, &stderr)
		if got != 7 {
			t.Errorf("Run(fail) = %d, want 7", got)
		}
	})

	t.Run("a name not in the qsokufile is a command-line error", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := Run([]string{"missing"}, nil, &stdout, &stderr)
		if got != 2 {
			t.Errorf("Run(missing) = %d, want 2", got)
		}
		if stderr.Len() == 0 {
			t.Error("Run(missing) wrote nothing to stderr")
		}
	})
}

func TestRun_noQsokufile(t *testing.T) {
	t.Chdir(t.TempDir())

	var stdout, stderr bytes.Buffer
	got := Run([]string{"anything"}, nil, &stdout, &stderr)
	if got != 1 {
		t.Errorf("Run(anything) = %d, want 1", got)
	}
	if stderr.Len() == 0 {
		t.Error("Run(anything) wrote nothing to stderr")
	}
}

func TestBuildVersion(t *testing.T) {
	if v := buildVersion(); v == "" {
		t.Error("buildVersion() returned an empty string")
	}
}

func writeQsokufile(t *testing.T, dir, content string) {
	t.Helper()
	path := filepath.Join(dir, "qsokufile")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path) //nolint:gosec // path is one writeQsokufile just built under t.TempDir()
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}
