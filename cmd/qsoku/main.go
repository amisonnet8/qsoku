// Command qsoku runs a repository's command shortcuts, defined in a
// qsokufile. It only hands its arguments to internal/cli.
package main

import (
	"os"

	"github.com/amisonnet8/qsoku/internal/cli"
)

// version is set by the release build (-ldflags "-X main.version={{ .Tag }}",
// .goreleaser.yaml); empty for go install, which falls back to the build
// information (.claude/rules/distribution.md).
var version string

func main() {
	cli.Version = version
	os.Exit(cli.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
