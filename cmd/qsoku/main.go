// Command qsoku runs a repository's command shortcuts, defined in a
// qsokufile. It only hands its arguments to internal/cli.
package main

import (
	"os"

	"github.com/amisonnet8/qsoku/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}
