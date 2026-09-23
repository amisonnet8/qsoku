package cli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/amisonnet8/qsoku/internal/qsokufile"
	"github.com/amisonnet8/qsoku/internal/run"
)

// runEdit implements "qsoku .edit" (docs/reference/cli.md "Management
// commands"): opens the qsokufile in use with $EDITOR, falling back to
// nano. It reuses internal/run.Execute -- the same subprocess machinery
// running a name uses -- so an editor gets the same cwd handoff and exit
// code passthrough as everything else, rather than a second copy of that
// logic.
func runEdit(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) != 0 {
		_, _ = fmt.Fprintln(stderr, "usage: qsoku .edit")
		return exitCommandLine
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		if _, err := exec.LookPath("nano"); err != nil {
			_, _ = fmt.Fprintln(stderr, "qsoku: $EDITOR is not set, and nano is not on PATH")
			return exitNotFound
		}
		editor = "nano"
	}

	path, err := qsokufile.Find(".")
	if err != nil {
		return reportFindError(err, stderr)
	}

	// $EDITOR may itself contain arguments (e.g. "code --wait"), so it is
	// run through sh rather than as a single program name; the path is
	// passed as $1 (not baked into the command text) so it is never
	// re-parsed by sh even if it contains spaces or shell metacharacters.
	root := filepath.Dir(path)
	command := editor + ` "$1"`
	code, err := run.Execute(command, root, []string{path}, stdin, stdout, stderr)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "qsoku: %v\n", err)
		return exitNotFound
	}
	return code
}
