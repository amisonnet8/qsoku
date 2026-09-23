package cli

import (
	"errors"
	"fmt"
	"io"

	"github.com/amisonnet8/qsoku/internal/qsokufile"
)

// runAdd implements "qsoku .add <name> <command>" (docs/reference/cli.md
// "Management commands").
func runAdd(args []string, stderr io.Writer) int {
	if len(args) != 2 {
		_, _ = fmt.Fprintln(stderr, "usage: qsoku .add <name> <command>")
		return exitCommandLine
	}
	name, command := args[0], args[1]

	path, err := qsokufile.Find(".")
	if err != nil {
		return reportFindError(err, stderr)
	}

	if err := qsokufile.SetEntry(path, name, command); err != nil {
		var invalid *qsokufile.InvalidNameError
		if errors.As(err, &invalid) {
			_, _ = fmt.Fprintf(stderr, "qsoku: %v\n", err)
			return exitCommandLine
		}
		_, _ = fmt.Fprintf(stderr, "qsoku: %v\n", err)
		return exitNotFound
	}
	return exitOK
}
