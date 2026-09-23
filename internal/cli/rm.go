package cli

import (
	"errors"
	"fmt"
	"io"

	"github.com/amisonnet8/qsoku/internal/qsokufile"
)

// runRm implements "qsoku .rm <name>" (docs/reference/cli.md "Management
// commands").
func runRm(args []string, stderr io.Writer) int {
	if len(args) != 1 {
		_, _ = fmt.Fprintln(stderr, "usage: qsoku .rm <name>")
		return exitCommandLine
	}
	name := args[0]

	path, err := qsokufile.Find(".")
	if err != nil {
		return reportFindError(err, stderr)
	}

	if err := qsokufile.RemoveEntry(path, name); err != nil {
		var notDefined *qsokufile.NotDefinedError
		if errors.As(err, &notDefined) {
			_, _ = fmt.Fprintf(stderr, "qsoku: %v in %s\n", err, path)
			return exitCommandLine
		}
		_, _ = fmt.Fprintf(stderr, "qsoku: %v\n", err)
		return exitNotFound
	}
	return exitOK
}
