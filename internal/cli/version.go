package cli

import "runtime/debug"

// buildVersion reports qsoku's own version: the module version from the
// build information (a tag for "go install ...@v0.1.0", a pseudo-version
// made of the commit time and hash for "go build"), else dev with the
// commit if it is known ("go install"-ed builds have no -ldflags to embed
// a version into, so the build information is the primary source, not a
// fallback; distribution.md).
func buildVersion() string {
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
