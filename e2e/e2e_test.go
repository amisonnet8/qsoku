//go:build e2e

// Package e2e runs the real qsoku binary against real shells. It is built
// with the tag e2e and run by make test; it is not part of make check.
//
// The tests in internal/cli call Run directly and check the details. These
// check that the pieces are joined: the binary that is built, real bash,
// zsh and fish, and the terminal that is not one.
package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// binary is the qsoku that TestMain built.
var binary string

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "qsoku-e2e-")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	binary = filepath.Join(dir, "qsoku")

	build := exec.Command("go", "build", "-o", binary, "github.com/amisonnet8/qsoku/cmd/qsoku") //nolint:gosec // fixed arguments building the real qsoku binary under test, not user input
	build.Stdout, build.Stderr = os.Stderr, os.Stderr
	if err := build.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "cannot build qsoku:", err)
		_ = os.RemoveAll(dir)
		os.Exit(1)
	}

	code := m.Run()
	_ = os.RemoveAll(dir)
	os.Exit(code)
}
