package cli

import (
	"fmt"
	"io"

	"github.com/amisonnet8/qsoku/internal/qsokufile"
)

// runList implements "qsoku .list" (docs/reference/cli.md "Management
// commands"): the defined names and their commands, one per line, in file
// order.
func runList(args []string, stdout, stderr io.Writer) int {
	if len(args) != 0 {
		_, _ = fmt.Fprintln(stderr, "usage: qsoku .list")
		return exitCommandLine
	}

	f, err := qsokufile.Load(".")
	if err != nil {
		return reportLoadError(err, stderr)
	}
	for _, e := range f.Entries {
		_, _ = fmt.Fprintf(stdout, "%s: %s\n", e.Name, e.Command)
	}
	return exitOK
}

// runNames implements "qsoku .names" (docs/reference/cli.md "Management
// commands"): just the names, one per line. Unlike every other management
// command, it never fails -- a missing or broken qsokufile silently means
// no names, not an error, so a shell completion script that calls it on
// every Tab press never sees an error interrupt the line being typed.
func runNames(args []string, stdout, stderr io.Writer) int {
	if len(args) != 0 {
		_, _ = fmt.Fprintln(stderr, "usage: qsoku .names")
		return exitCommandLine
	}

	f, err := qsokufile.Load(".")
	if err != nil {
		return exitOK
	}
	for _, e := range f.Entries {
		_, _ = fmt.Fprintln(stdout, e.Name)
	}
	return exitOK
}
