package qsokufile

import (
	"errors"
	"os"
	"path/filepath"
)

// fileName is the exact name qsoku looks for (docs/reference/qsokufile.md
// "Location and lookup"): always lowercase, regardless of the platform.
const fileName = "qsokufile"

// ErrNotFound is returned by Find when no qsokufile exists from startDir up
// to the filesystem root.
var ErrNotFound = errors.New("no qsokufile found")

// Find walks up from startDir, through each parent directory, and returns
// the path of the first qsokufile it finds. It does not stop at a git
// repository boundary, and does not look above the filesystem root.
//
// A directory that happens to be named "qsokufile" is not a match: Find
// keeps walking up past it, the same as if nothing were there.
//
// The returned path always uses '/' (filepath.ToSlash), even on Windows:
// this is the one place a qsokufile's path is produced, and every path a
// caller builds from it -- printed in a message, or, via filepath.Dir,
// turned into QSOKU_ROOT -- should read the same way "pwd -W" already does
// (docs/reference/cli.md "Bringing the working directory back") rather than
// mixing native backslashes into an otherwise forward-slash world. The walk
// itself still uses native paths throughout; only the final result is
// converted.
func Find(startDir string) (string, error) {
	dir, err := filepath.Abs(startDir)
	if err != nil {
		return "", err
	}

	for {
		candidate := filepath.Join(dir, fileName)
		info, err := os.Stat(candidate)
		switch {
		case err == nil && !info.IsDir():
			return filepath.ToSlash(candidate), nil
		case err == nil, os.IsNotExist(err):
			// Either a directory named "qsokufile" (not a match) or nothing
			// there at all: keep walking up.
		default:
			return "", err
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", ErrNotFound
		}
		dir = parent
	}
}
