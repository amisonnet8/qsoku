package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// runInit implements "qsoku .init" (docs/reference/cli.md "Management
// commands"): creates an empty qsokufile in the current directory only, not
// walking up to a parent (unlike everything else, which does).
func runInit(args []string, stdout, stderr io.Writer) int {
	if len(args) != 0 {
		_, _ = fmt.Fprintln(stderr, "usage: qsoku .init")
		return exitCommandLine
	}

	cwd, err := os.Getwd()
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "qsoku: %v\n", err)
		return exitNotFound
	}
	path := filepath.Join(cwd, "qsokufile")

	if _, err := os.Stat(path); err == nil {
		_, _ = fmt.Fprintf(stderr, "qsoku: %s already exists\n", path)
		return exitNotFound
	} else if !os.IsNotExist(err) {
		_, _ = fmt.Fprintf(stderr, "qsoku: %v\n", err)
		return exitNotFound
	}

	if err := os.WriteFile(path, nil, 0o600); err != nil {
		_, _ = fmt.Fprintf(stderr, "qsoku: %v\n", err)
		return exitNotFound
	}
	_, _ = fmt.Fprintf(stdout, "Created %s\n", path)
	return exitOK
}
