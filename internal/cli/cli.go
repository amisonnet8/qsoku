// Package cli implements qsoku's command line.
package cli

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/amisonnet8/qsoku/internal/qsokufile"
	"github.com/amisonnet8/qsoku/internal/run"
)

// Run executes qsoku's command line (the arguments after "qsoku" itself) and
// returns the process exit code.
func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) == 1 && args[0] == ".version" {
		_, _ = fmt.Fprintln(stdout, buildVersion())
		return 0
	}
	// Management commands (Step 6) and no-args .help are not implemented
	// yet.
	if len(args) == 0 || strings.HasPrefix(args[0], ".") {
		_, _ = fmt.Fprintln(stderr, "qsoku: not implemented yet")
		return 1
	}
	return runName(args[0], args[1:], stdin, stdout, stderr)
}

// runName runs the qsokufile entry called name (docs/reference/cli.md
// "Running a name").
func runName(name string, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	f, err := qsokufile.Load(".")
	if err != nil {
		if errors.Is(err, qsokufile.ErrNotFound) {
			_, _ = fmt.Fprintln(stderr, "qsoku: no qsokufile found; run qsoku .init")
			return exitNotFound
		}
		_, _ = fmt.Fprintf(stderr, "qsoku: %v\n", err)
		return exitNotFound
	}

	entry, ok := f.Lookup(name)
	if !ok {
		_, _ = fmt.Fprintf(stderr, "qsoku: %q is not defined in %s\n", name, f.Path)
		return exitCommandLine
	}

	root := filepath.Dir(f.Path)
	command := qsokufile.Substitute(entry.Command)
	subArgs := make([]string, len(args))
	for i, a := range args {
		subArgs[i] = qsokufile.SubstituteArg(a, root)
	}

	code, err := run.Execute(command, root, subArgs, stdin, stdout, stderr)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "qsoku: %v\n", err)
		return exitNotFound
	}
	return code
}

// qsoku's own exit codes (docs/reference/cli.md "Exit codes"). Anything
// past this point is the command's own exit code, returned unchanged.
const (
	exitNotFound    = 1 // qsoku could not do what was asked
	exitCommandLine = 2 // the command line is wrong
)
