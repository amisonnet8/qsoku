package qsokufile

import (
	"fmt"
	"os"
)

// Load finds the qsokufile in use starting from startDir (see Find) and
// parses it (see Parse). An error reading or parsing the file (but not one
// from Find itself, which already names the directories it tried) is
// wrapped with the file's path, so a *ParseError reads the way
// docs/reference/qsokufile.md "Format" describes it: naming the file and
// the line.
func Load(startDir string) (*File, error) {
	path, err := Find(startDir)
	if err != nil {
		return nil, err
	}

	r, err := os.Open(path) //nolint:gosec // path comes from Find, which only ever joins a fixed name onto directories walked from startDir; reading the qsokufile is qsoku's own job, like G204 below it for running one
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	defer func() { _ = r.Close() }()

	f, err := Parse(r)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	f.Path = path
	return f, nil
}
