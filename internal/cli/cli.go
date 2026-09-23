// Package cli implements qsoku's command line.
package cli

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"

	"github.com/amisonnet8/qsoku/internal/qsokufile"
	"github.com/amisonnet8/qsoku/internal/run"
)

// Run executes qsoku's command line (the arguments after "qsoku" itself) and
// returns the process exit code.
func Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		return runHelp(stdout, stderr)
	}

	switch args[0] {
	case ".version":
		if len(args) != 1 {
			_, _ = fmt.Fprintln(stderr, "usage: qsoku .version")
			return exitCommandLine
		}
		_, _ = fmt.Fprintln(stdout, buildVersion())
		return exitOK
	case ".init":
		return runInit(args[1:], stdout, stderr)
	case ".add":
		return runAdd(args[1:], stderr)
	case ".rm":
		return runRm(args[1:], stderr)
	case ".list":
		return runList(args[1:], stdout, stderr)
	case ".names":
		return runNames(args[1:], stdout, stderr)
	case ".edit":
		return runEdit(args[1:], stdin, stdout, stderr)
	case ".where":
		return runWhere(args[1:], stdout, stderr)
	case ".help":
		return runHelp(stdout, stderr)
	case ".shell":
		// Shell integration and completion (Step 7): not implemented yet.
		_, _ = fmt.Fprintln(stderr, "qsoku: not implemented yet")
		return exitNotFound
	}

	if args[0] != "" && args[0][0] == '.' {
		_, _ = fmt.Fprintf(stderr, "qsoku: unknown command %q\n", args[0])
		return exitCommandLine
	}
	return runName(args[0], args[1:], stdin, stdout, stderr)
}

// runName runs the qsokufile entry called name (docs/reference/cli.md
// "Running a name").
func runName(name string, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	f, err := qsokufile.Load(".")
	if err != nil {
		return reportLoadError(err, stderr)
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

// reportFindError and reportLoadError print a qsokufile.Find/Load error the
// way every management command that needs an existing qsokufile does, and
// return the exit code to use.
func reportFindError(err error, stderr io.Writer) int {
	if errors.Is(err, qsokufile.ErrNotFound) {
		_, _ = fmt.Fprintln(stderr, "qsoku: no qsokufile found; run qsoku .init")
		return exitNotFound
	}
	_, _ = fmt.Fprintf(stderr, "qsoku: %v\n", err)
	return exitNotFound
}

func reportLoadError(err error, stderr io.Writer) int {
	return reportFindError(err, stderr)
}

// qsoku's own exit codes (docs/reference/cli.md "Exit codes"). Anything
// past this point is the command's own exit code, returned unchanged.
const (
	exitOK          = 0
	exitNotFound    = 1 // qsoku could not do what was asked
	exitCommandLine = 2 // the command line is wrong
)
