package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun_help(t *testing.T) {
	t.Run("no args and .help print the same thing", func(t *testing.T) {
		t.Chdir(t.TempDir())

		var out1, out2, stderr bytes.Buffer
		code1 := Run(nil, nil, &out1, &stderr)
		code2 := Run([]string{".help"}, nil, &out2, &stderr)
		if code1 != 0 || code2 != 0 {
			t.Fatalf("Run() = %d, %d, want 0, 0", code1, code2)
		}
		if out1.String() != out2.String() {
			t.Errorf("no-args output differs from .help output:\n%q\n%q", out1.String(), out2.String())
		}
		if !strings.Contains(out1.String(), "usage: qsoku") {
			t.Errorf("help output = %q, want it to contain usage text", out1.String())
		}
	})

	t.Run("lists the defined names when a qsokufile is found", func(t *testing.T) {
		dir := t.TempDir()
		writeQsokufile(t, dir, "root: cd //\n")
		t.Chdir(dir)

		var stdout, stderr bytes.Buffer
		Run([]string{".help"}, nil, &stdout, &stderr)
		if !strings.Contains(stdout.String(), "root: cd //") {
			t.Errorf("help output = %q, want it to list the defined entry", stdout.String())
		}
	})

	t.Run("never fails, even with a broken qsokufile", func(t *testing.T) {
		dir := t.TempDir()
		writeQsokufile(t, dir, "not valid\n")
		t.Chdir(dir)

		var stdout, stderr bytes.Buffer
		got := Run([]string{".help"}, nil, &stdout, &stderr)
		if got != 0 {
			t.Errorf("Run(.help) = %d, want 0", got)
		}
	})

	t.Run("never fails on extra arguments", func(t *testing.T) {
		t.Chdir(t.TempDir())

		var stdout, stderr bytes.Buffer
		got := Run([]string{".help", "extra", "args"}, nil, &stdout, &stderr)
		if got != 0 {
			t.Errorf("Run(.help extra args) = %d, want 0", got)
		}
	})
}
