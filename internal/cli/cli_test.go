package cli

import (
	"bytes"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantCode int
	}{
		{"version", []string{".version"}, 0},
		{"no args", nil, 1},
		{"unimplemented name", []string{"build"}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			if got := Run(tt.args, &stdout, &stderr); got != tt.wantCode {
				t.Errorf("Run(%v) = %d, want %d", tt.args, got, tt.wantCode)
			}
		})
	}
}

func TestBuildVersion(t *testing.T) {
	if v := buildVersion(); v == "" {
		t.Error("buildVersion() returned an empty string")
	}
}
