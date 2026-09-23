package cli

import (
	"fmt"
	"io"

	"github.com/amisonnet8/qsoku/internal/qsokufile"
)

const usage = `usage: qsoku <name> [args...]
       qsoku <command> [args...]

Commands:
  qsoku .init                  Create an empty qsokufile here
  qsoku .add <name> <command>  Add a name, or replace its command
  qsoku .rm <name>             Remove a name
  qsoku .list                  List the defined names and their commands
  qsoku .names                 List just the names
  qsoku .edit                  Edit the qsokufile with $EDITOR
  qsoku .where                 Show the qsokufile's directory
  qsoku .version               Show qsoku's version
  qsoku .shell <shell>         Print shell integration code
  qsoku .help                  Show this help
`

// runHelp implements "qsoku .help", and running qsoku with no arguments at
// all (docs/reference/cli.md "Management commands"). Unlike every other
// management command, it never fails: whatever arguments it was called
// with, and whether a qsokufile can be found at all, it always has
// something useful to show.
func runHelp(stdout, _ io.Writer) int {
	_, _ = io.WriteString(stdout, usage)

	f, err := qsokufile.Load(".")
	if err != nil {
		return exitOK
	}
	if len(f.Entries) == 0 {
		return exitOK
	}
	_, _ = fmt.Fprintf(stdout, "\nDefined in %s:\n", f.Path)
	for _, e := range f.Entries {
		_, _ = fmt.Fprintf(stdout, "  %s: %s\n", e.Name, e.Command)
	}
	return exitOK
}
