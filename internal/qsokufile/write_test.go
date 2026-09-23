package qsokufile

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func writeTestFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), fileName)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path) //nolint:gosec // path is one writeTestFile just built under t.TempDir()
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestSetEntry(t *testing.T) {
	tests := []struct {
		name    string
		initial string
		setName string
		command string
		want    string
		wantErr string
	}{
		{
			name:    "append to an empty file",
			initial: "",
			setName: "root",
			command: "cd //",
			want:    "root: cd //\n",
		},
		{
			name:    "append after existing entries, keeping comments and order",
			initial: "# top comment\nroot: cd //\n\nbuild: (cd //; make build)\n",
			setName: "test",
			command: "(cd //; go test ./...)",
			want:    "# top comment\nroot: cd //\n\nbuild: (cd //; make build)\ntest: (cd //; go test ./...)\n",
		},
		{
			name:    "append when the file has no trailing newline",
			initial: "root: cd //",
			setName: "build",
			command: "make",
			want:    "root: cd //\nbuild: make\n",
		},
		{
			name:    "replace an existing entry in place",
			initial: "# comment\nroot: cd //\nbuild: old command\ntest: (cd //; go test ./...)\n",
			setName: "build",
			command: "new command",
			want:    "# comment\nroot: cd //\nbuild: new command\ntest: (cd //; go test ./...)\n",
		},
		{
			name:    "invalid name is rejected without touching the file",
			initial: "root: cd //\n",
			setName: ".init",
			command: "echo no",
			wantErr: `invalid name ".init"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := writeTestFile(t, tt.initial)

			err := SetEntry(path, tt.setName, tt.command)

			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("SetEntry() error = %v, want %q", err, tt.wantErr)
				}
				if got := readTestFile(t, path); got != tt.initial {
					t.Errorf("file was modified on error: got %q, want unchanged %q", got, tt.initial)
				}
				return
			}
			if err != nil {
				t.Fatalf("SetEntry() unexpected error: %v", err)
			}
			if got := readTestFile(t, path); got != tt.want {
				t.Errorf("file = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRemoveEntry(t *testing.T) {
	tests := []struct {
		name    string
		initial string
		rmName  string
		want    string
		wantErr bool
	}{
		{
			name:    "removes only the named line",
			initial: "# comment\nroot: cd //\nbuild: make\ntest: go test\n",
			rmName:  "build",
			want:    "# comment\nroot: cd //\ntest: go test\n",
		},
		{
			name:    "the only entry",
			initial: "root: cd //\n",
			rmName:  "root",
			want:    "",
		},
		{
			name:    "name not defined",
			initial: "root: cd //\n",
			rmName:  "missing",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := writeTestFile(t, tt.initial)

			err := RemoveEntry(path, tt.rmName)

			if tt.wantErr {
				var notDefined *NotDefinedError
				if !errors.As(err, &notDefined) {
					t.Fatalf("RemoveEntry() error = %v, want *NotDefinedError", err)
				}
				if got := readTestFile(t, path); got != tt.initial {
					t.Errorf("file was modified on error: got %q, want unchanged %q", got, tt.initial)
				}
				return
			}
			if err != nil {
				t.Fatalf("RemoveEntry() unexpected error: %v", err)
			}
			if got := readTestFile(t, path); got != tt.want {
				t.Errorf("file = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSetEntry_preservesPermissions(t *testing.T) {
	path := writeTestFile(t, "root: cd //\n")
	if err := os.Chmod(path, 0o640); err != nil { //nolint:gosec // deliberately testing that a non-default, group-readable mode on a t.TempDir() fixture survives SetEntry unchanged
		t.Fatal(err)
	}

	if err := SetEntry(path, "build", "make"); err != nil {
		t.Fatalf("SetEntry() error: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o640 {
		t.Errorf("mode = %v, want %v", got, os.FileMode(0o640))
	}
}
