// Package qsokufile finds and parses a qsokufile: the file that holds a
// repository's command shortcuts (docs/reference/qsokufile.md).
package qsokufile

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Entry is one name: command line of a qsokufile.
type Entry struct {
	Name    string
	Command string
	Line    int // 1-indexed
}

// File is a parsed qsokufile: its entries, in the order they appear.
type File struct {
	Path    string
	Entries []Entry
}

// ParseError reports a malformed line: a line with no ':', an invalid name,
// or a name already defined earlier in the file.
type ParseError struct {
	Line int
	Msg  string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("line %d: %s", e.Line, e.Msg)
}

// isNameByte reports whether b is one of the ASCII characters a qsokufile
// name may contain: letters, digits, '.', '_', '-'.
func isNameByte(b byte) bool {
	switch {
	case b >= 'a' && b <= 'z', b >= 'A' && b <= 'Z', b >= '0' && b <= '9':
		return true
	case b == '.' || b == '_' || b == '-':
		return true
	default:
		return false
	}
}

func validName(name string) bool {
	if name == "" || name[0] == '.' {
		return false
	}
	for i := 0; i < len(name); i++ {
		if !isNameByte(name[i]) {
			return false
		}
	}
	return true
}

// Parse reads a qsokufile's content and returns its entries. It stops at the
// first malformed line (a *ParseError): a line with no ':', an invalid name,
// or a name defined more than once.
func Parse(r io.Reader) (*File, error) {
	f := &File{}
	seen := make(map[string]int) // name -> first defining line

	scanner := bufio.NewScanner(r)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		i := strings.IndexByte(line, ':')
		if i < 0 {
			return nil, &ParseError{Line: lineNo, Msg: "no ':' found (expected \"name: command\")"}
		}
		name := line[:i]
		command := strings.TrimLeft(line[i+1:], " \t")

		if !validName(name) {
			return nil, &ParseError{Line: lineNo, Msg: fmt.Sprintf("invalid name %q", name)}
		}
		if first, ok := seen[name]; ok {
			return nil, &ParseError{Line: lineNo, Msg: fmt.Sprintf("name %q already defined at line %d", name, first)}
		}
		seen[name] = lineNo

		f.Entries = append(f.Entries, Entry{Name: name, Command: command, Line: lineNo})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return f, nil
}
