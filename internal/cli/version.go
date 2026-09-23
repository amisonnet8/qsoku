package cli

import "runtime/debug"

// version is set at release build time (-ldflags
// "-X github.com/amisonnet8/qsoku/internal/cli.version={{.Tag}}", by
// .goreleaser.yaml). Without it, the version comes from the build
// information (distribution.md).
var version string

// buildVersion reports qsoku's own version: the version set at build time
// (a GoReleaser release build), else the module version from the build
// information (a tag for "go install ...@v0.1.0", a pseudo-version made of
// the commit time and hash for "go build"), else dev with the commit if it
// is known ("go install"-ed builds have no -ldflags to embed a version
// into, so the build information is the fallback for that path, not the
// primary source overall; distribution.md).
func buildVersion() string {
	if version != "" {
		return version
	}
	info, ok := debug.ReadBuildInfo()
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
