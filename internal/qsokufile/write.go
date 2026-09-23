package qsokufile

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

// NotDefinedError is returned by RemoveEntry when name has no entry to
// remove.
type NotDefinedError struct {
	Name string
}

func (e *NotDefinedError) Error() string {
	return fmt.Sprintf("%q is not defined", e.Name)
}

// InvalidNameError is returned by SetEntry when name is not a valid
// qsokufile name (docs/reference/qsokufile.md "Names").
type InvalidNameError struct {
	Name string
}

func (e *InvalidNameError) Error() string {
	return fmt.Sprintf("invalid name %q", e.Name)
}

// SetEntry adds "name: command" to the qsokufile at path. If name already
// has an entry, that one line is replaced in place (same position); every
// other line -- comments, blank lines, the rest of the order -- is left
// untouched. The file's own permissions are kept.
func SetEntry(path, name, command string) error {
	if !validName(name) {
		return &InvalidNameError{Name: name}
	}

	raw, mode, err := readFile(path)
	if err != nil {
		return err
	}
	f, err := Parse(bytes.NewReader(raw))
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}

	lines := strings.Split(string(raw), "\n")
	newLine := name + ": " + command
	if entry, ok := f.Lookup(name); ok {
		lines[entry.Line-1] = newLine
	} else {
		lines = appendLine(lines, newLine)
	}

	return writeFile(path, lines, mode)
}

// RemoveEntry removes name's line from the qsokufile at path, leaving every
// other line untouched. It returns a *NotDefinedError if name has no entry.
func RemoveEntry(path, name string) error {
	raw, mode, err := readFile(path)
	if err != nil {
		return err
	}
	f, err := Parse(bytes.NewReader(raw))
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	entry, ok := f.Lookup(name)
	if !ok {
		return &NotDefinedError{Name: name}
	}

	lines := strings.Split(string(raw), "\n")
	lines = append(lines[:entry.Line-1], lines[entry.Line:]...)

	return writeFile(path, lines, mode)
}

func readFile(path string) (data []byte, mode os.FileMode, err error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, 0, err
	}
	data, err = os.ReadFile(path) //nolint:gosec // path comes from Find, which only ever joins a fixed name onto directories walked from a caller-given start; reading/writing the qsokufile is qsoku's own job, like G204 elsewhere for running one
	if err != nil {
		return nil, 0, err
	}
	return data, info.Mode().Perm(), nil
}

// appendLine adds line to the end of lines, keeping exactly one trailing
// empty element (so the joined text ends with a single newline) regardless
// of whether the original content already ended with one.
func appendLine(lines []string, line string) []string {
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines[len(lines)-1] = line
		return append(lines, "")
	}
	return append(lines, line, "")
}

func writeFile(path string, lines []string, mode os.FileMode) error {
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")), mode)
}
