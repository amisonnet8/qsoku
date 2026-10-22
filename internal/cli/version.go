package cli

import "runtime/debug"

// Version is the version set at release build time, handed over by package
// main (its own version, which .goreleaser.yaml sets with
// -ldflags "-X main.version={{ .Tag }}"). Without it, the version comes from
// the build information (distribution.md).
var Version string

// buildVersion reports qsoku's own version: the version set at build time
// (a GoReleaser release build), else the module version from the build
// information (a tag for "go install ...@v0.1.0", a pseudo-version made of
// the commit time and hash for "go build"), else dev with the commit if it
// is known ("go install"-ed builds have no -ldflags to embed a version
// into, so the build information is the fallback for that path, not the
// primary source overall; distribution.md).
func buildVersion() string {
	return versionFrom(Version, debug.ReadBuildInfo)
}

// versionFrom is buildVersion with the build information passed in, so the
// tests can give it each kind (the test binary's own has no tag).
func versionFrom(set string, read func() (*debug.BuildInfo, bool)) string {
	if set != "" {
		return set
	}
	info, ok := read()
	if !ok {
		return "dev"
	}
	if v := info.Main.Version; v != "" && v != "(devel)" {
		return v
	}
	for _, s := range info.Settings {
		if s.Key == "vcs.revision" && len(s.Value) >= 7 {
			return "dev (" + s.Value[:7] + ")"
		}
	}
	return "dev"
}
