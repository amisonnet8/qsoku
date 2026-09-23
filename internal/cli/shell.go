package cli

import (
	"embed"
	"fmt"
	"io"
)

// The scripts are fixed text: they list the management commands themselves
// (unlike mtqg's fully dynamic completion), since that list only changes
// with a new qsoku release, at which point these scripts -- embedded in the
// same binary -- change with it. Only the defined names, which can change
// on every Tab press without a new qsoku, are fetched dynamically via
// `qsoku .names` (docs/reference/cli.md "Shell completion").
//
//go:embed shells
var shellScripts embed.FS

var shellFiles = map[string]string{
	"bash": "shells/qsoku.bash",
	"zsh":  "shells/qsoku.zsh",
	"fish": "shells/qsoku.fish",
}

// runShell implements "qsoku .shell <shell>" (docs/reference/cli.md
// "Management commands"): prints that shell's integration and completion
// script, to eval.
func runShell(args []string, stdout, stderr io.Writer) int {
	if len(args) != 1 {
		_, _ = fmt.Fprintln(stderr, "usage: qsoku .shell <bash|zsh|fish>")
		return exitCommandLine
	}
	path, ok := shellFiles[args[0]]
	if !ok {
		_, _ = fmt.Fprintf(stderr, "qsoku: unknown shell %q (want bash, zsh, or fish)\n", args[0])
		return exitCommandLine
	}

	script, err := shellScripts.ReadFile(path)
	if err != nil {
		// The scripts are embedded at build time; this only fails if the
		// embed itself is missing, not on anything a user did.
		_, _ = fmt.Fprintf(stderr, "qsoku: %v\n", err)
		return exitNotFound
	}
	_, _ = stdout.Write(script)
	return exitOK
}
