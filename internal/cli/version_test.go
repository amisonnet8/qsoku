package cli

import "testing"

// TestBuildVersion_ldflags checks that a version embedded at build time
// (-ldflags, by .goreleaser.yaml) takes priority over the build information
// (which the "go install" path still relies on -- see buildVersion).
func TestBuildVersion_ldflags(t *testing.T) {
	old := version
	defer func() { version = old }()

	version = "v0.1.0"
	if got := buildVersion(); got != "v0.1.0" {
		t.Errorf("buildVersion() = %q, want %q", got, "v0.1.0")
	}
}
