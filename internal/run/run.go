// Package run executes a qsokufile entry's command with sh
// (docs/reference/cli.md "Running a name"). It does not know what a
// qsokufile is; it only runs a command string.
package run

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

// trailer brings back the working directory qsoku finishes in
// (docs/reference/cli.md "Bringing the working directory back"). It is
// joined to the command with a newline, not ';', so a trailing '#' comment
// on the command can't swallow it.
//
// On Windows, `sh` is Git for Windows' MSYS `sh.exe`: its own `pwd` prints
// the MSYS path form (/c/Users/...), which the shell integration's
// Set-Location/cd cannot use. `pwd -W` prints the Windows form (C:/Users/...)
// instead; a plain POSIX sh (Linux, macOS) does not know that flag and
// exits 1 with a message on stderr, so it is written first with stderr
// discarded, falling back to plain `pwd` if it fails.
const trailer = "\n__status=$?; pwd > \"$QSOKU_CWD_FILE\"; exit $__status"
const trailerWindows = "\n__status=$?; { pwd -W 2>/dev/null || pwd; } > \"$QSOKU_CWD_FILE\"; exit $__status"

// Execute runs command with sh, from the current directory, with root
// available to it as QSOKU_ROOT. args become $1, $2, and so on inside
// command.
//
// QSOKU_CWD_FILE is read from qsoku's own environment (set there by the
// shell integration function), not passed in: when it is present, the
// trailer above is appended so the location qsoku finishes in gets written
// there, on success or failure alike. When it is absent, the trailer is
// left out entirely -- appending it regardless would redirect pwd to an
// empty filename, which sh reports as an error, and nothing would read the
// file anyway.
//
// Execute returns the command's own exit code (128+n if it was killed by
// signal n). A non-nil error means sh itself could not be started; the
// returned code is meaningless in that case.
func Execute(command, root string, args []string, stdin io.Reader, stdout, stderr io.Writer) (int, error) {
	script := command
	if _, ok := os.LookupEnv("QSOKU_CWD_FILE"); ok {
		if runtime.GOOS == "windows" {
			script += trailerWindows
		} else {
			script += trailer
		}
	}

	cmdArgs := append([]string{"-c", script, "qsoku"}, args...)
	cmd := exec.Command("sh", cmdArgs...) //nolint:gosec // "sh" is a fixed program name; running a qsokufile's own command text is qsoku's own job (see the repo-wide G204 note in .golangci.yaml)
	cmd.Env = append(os.Environ(), "QSOKU_ROOT="+root)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()
	if err == nil {
		return 0, nil
	}

	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		// sh itself could not be started (not found, permission, ...).
		return 0, err
	}
	if ws, ok := exitErr.Sys().(syscall.WaitStatus); ok && ws.Signaled() {
		return 128 + int(ws.Signal()), nil
	}
	return exitErr.ExitCode(), nil
}
