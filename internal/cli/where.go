package cli

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/amisonnet8/qsoku/internal/qsokufile"
)

// runWhere implements "qsoku .where" (docs/reference/cli.md "Management
// commands"): the absolute path of the directory holding the qsokufile in
// use (what "//" and QSOKU_ROOT expand to). Like .edit, it only needs Find,
// not a successful Parse, so it still works on a qsokufile that currently
// fails to parse.
func runWhere(args []string, stdout, stderr io.Writer) int {
	if len(args) != 0 {
		_, _ = fmt.Fprintln(stderr, "usage: qsoku .where")
		return exitCommandLine
	}

	path, err := qsokufile.Find(".")
	if err != nil {
		return reportFindError(err, stderr)
	}
	_, _ = fmt.Fprintln(stdout, filepath.Dir(path))
	return exitOK
}
