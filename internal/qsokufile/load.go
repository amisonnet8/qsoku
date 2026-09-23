package qsokufile

import "os"

// Load finds the qsokufile in use starting from startDir (see Find) and
// parses it (see Parse).
func Load(startDir string) (*File, error) {
	path, err := Find(startDir)
	if err != nil {
		return nil, err
	}

	r, err := os.Open(path) //nolint:gosec // path comes from Find, which only ever joins a fixed name onto directories walked from startDir; reading the qsokufile is qsoku's own job, like G204 below it for running one
	if err != nil {
		return nil, err
	}
	defer func() { _ = r.Close() }()

	f, err := Parse(r)
	if err != nil {
		return nil, err
	}
	f.Path = path
	return f, nil
}
