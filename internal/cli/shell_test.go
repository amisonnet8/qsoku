package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun_shell(t *testing.T) {
	tests := []struct {
		shell string
		want  string // a substring that should be in the script
	}{
		{"bash", "qsoku()"},
		{"zsh", "qsoku()"},
		{"fish", "function qsoku"},
	}
	for _, tt := range tests {
		t.Run(tt.shell, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			got := Run([]string{".shell", tt.shell}, nil, &stdout, &stderr)
			if got != 0 {
				t.Fatalf("Run(.shell %s) = %d, want 0; stderr: %s", tt.shell, got, stderr.String())
			}
			if !strings.Contains(stdout.String(), tt.want) {
				t.Errorf("script for %s does not contain %q:\n%s", tt.shell, tt.want, stdout.String())
			}
			if !strings.Contains(stdout.String(), "QSOKU_CWD_FILE") {
				t.Errorf("script for %s does not mention QSOKU_CWD_FILE", tt.shell)
			}
		})
	}

	t.Run("unknown shell", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		got := Run([]string{".shell", "powershell"}, nil, &stdout, &stderr)
		if got != 2 {
			t.Errorf("Run(.shell powershell) = %d, want 2", got)
		}
		if stderr.Len() == 0 {
			t.Error("Run(.shell powershell) wrote nothing to stderr")
		}
	})

	t.Run("wrong number of arguments", func(t *testing.T) {
		var stdout, stderr bytes.Buffer
		if got := Run([]string{".shell"}, nil, &stdout, &stderr); got != 2 {
			t.Errorf("Run(.shell) = %d, want 2", got)
		}
		stdout.Reset()
		stderr.Reset()
		if got := Run([]string{".shell", "bash", "extra"}, nil, &stdout, &stderr); got != 2 {
			t.Errorf("Run(.shell bash extra) = %d, want 2", got)
		}
	})
}
