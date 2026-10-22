package cli

import (
	"runtime/debug"
	"testing"
)

// TestBuildVersion_ldflags checks that a version embedded at build time
// (-ldflags, by .goreleaser.yaml) takes priority over the build information
// (which the "go install" path still relies on -- see buildVersion).
func TestBuildVersion_ldflags(t *testing.T) {
	old := Version
	defer func() { Version = old }()

	Version = "v0.1.0"
	if got := buildVersion(); got != "v0.1.0" {
		t.Errorf("buildVersion() = %q, want %q", got, "v0.1.0")
	}
}

// TestVersionFrom gives each kind of build information. The commit is cut to
// 7 characters (git's short hash); a shorter one is not shown.
func TestVersionFrom(t *testing.T) {
	built := func(version string, settings ...debug.BuildSetting) func() (*debug.BuildInfo, bool) {
		return func() (*debug.BuildInfo, bool) {
			return &debug.BuildInfo{Main: debug.Module{Version: version}, Settings: settings}, true
		}
	}
	rev := func(v string) debug.BuildSetting { return debug.BuildSetting{Key: "vcs.revision", Value: v} }
	none := func() (*debug.BuildInfo, bool) { return nil, false }

	for _, c := range []struct {
		name string
		set  string
		read func() (*debug.BuildInfo, bool)
		want string
	}{
		{"set by the release build wins", "v1.2.3", built("v0.1.0"), "v1.2.3"},
		{"no build information", "", none, "dev"},
		{"module version (go install ...@v0.1.0)", "", built("v0.1.0", rev("0123456789abcdef")), "v0.1.0"},
		{"devel with a commit", "", built("(devel)", rev("0123456789abcdef")), "dev (0123456)"},
		{"no version with a commit", "", built("", rev("0123456789abcdef")), "dev (0123456)"},
		{"a commit of exactly 7", "", built("(devel)", rev("abcdefg")), "dev (abcdefg)"},
		{"a commit shorter than 7", "", built("(devel)", rev("abcdef")), "dev"},
		{"the commit after another setting", "", built("(devel)", debug.BuildSetting{Key: "vcs.time", Value: "2026-10-23T00:00:00Z"}, rev("0123456789abcdef")), "dev (0123456)"},
		{"no commit", "", built("(devel)", debug.BuildSetting{Key: "vcs.time", Value: "2026-10-23T00:00:00Z"}), "dev"},
	} {
		if got := versionFrom(c.set, c.read); got != c.want {
			t.Errorf("%s: versionFrom = %q, want %q", c.name, got, c.want)
		}
	}
}
