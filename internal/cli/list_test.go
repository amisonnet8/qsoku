package cli

import (
	"bytes"
	"testing"
)

func TestRun_list(t *testing.T) {
	dir := t.TempDir()
	writeQsokufile(t, dir, "root: cd //\nbuild: make build\n")
	t.Chdir(dir)

	var stdout, stderr bytes.Buffer
	got := Run([]string{".list"}, nil, &stdout, &stderr)
	if got != 0 {
		t.Fatalf("Run(.list) = %d, want 0; stderr: %s", got, stderr.String())
	}
	want := "root: cd //\nbuild: make build\n"
	if stdout.String() != want {
		t.Errorf("stdout = %q, want %q", stdout.String(), want)
	}
}

func TestRun_names(t *testing.T) {
	t.Run("lists just the names", func(t *testing.T) {
		dir := t.TempDir()
		writeQsokufile(t, dir, "root: cd //\nbuild: make build\n")
		t.Chdir(dir)

		var stdout, stderr bytes.Buffer
		got := Run([]string{".names"}, nil, &stdout, &stderr)
		if got != 0 {
			t.Fatalf("Run(.names) = %d, want 0; stderr: %s", got, stderr.String())
		}
		want := "root\nbuild\n"
		if stdout.String() != want {
			t.Errorf("stdout = %q, want %q", stdout.String(), want)
		}
	})

	t.Run("never fails: no qsokufile", func(t *testing.T) {
		t.Chdir(t.TempDir())

		var stdout, stderr bytes.Buffer
		got := Run([]string{".names"}, nil, &stdout, &stderr)
		if got != 0 {
			t.Errorf("Run(.names) = %d, want 0", got)
		}
		if stderr.Len() != 0 {
			t.Errorf("stderr = %q, want empty", stderr.String())
		}
		if stdout.Len() != 0 {
			t.Errorf("stdout = %q, want empty", stdout.String())
		}
	})

	t.Run("never fails: a broken qsokufile", func(t *testing.T) {
		dir := t.TempDir()
		writeQsokufile(t, dir, "this is not valid\n")
		t.Chdir(dir)

		var stdout, stderr bytes.Buffer
		got := Run([]string{".names"}, nil, &stdout, &stderr)
		if got != 0 {
			t.Errorf("Run(.names) = %d, want 0", got)
		}
		if stderr.Len() != 0 {
			t.Errorf("stderr = %q, want empty", stderr.String())
		}
	})

	t.Run("still validates its own argument count", func(t *testing.T) {
		t.Chdir(t.TempDir())

		var stdout, stderr bytes.Buffer
		got := Run([]string{".names", "extra"}, nil, &stdout, &stderr)
		if got != 2 {
			t.Errorf("Run(.names extra) = %d, want 2", got)
		}
	})
}
