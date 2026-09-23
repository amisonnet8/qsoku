// Package cli implements qsoku's command line.
package cli

import (
	"fmt"
	"io"
)

// Run executes qsoku's command line (the arguments after "qsoku" itself) and
// returns the process exit code.
func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 1 && args[0] == ".version" {
		_, _ = fmt.Fprintln(stdout, buildVersion())
		return 0
	}
	_, _ = fmt.Fprintln(stderr, "qsoku: not implemented yet")
	return 1
}
