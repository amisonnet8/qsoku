package qsokufile

import "strings"

// rootVar is the literal shell text substituted for a "//" that stands for
// the qsokufile's own location, inside a qsokufile command
// (docs/reference/qsokufile.md "// : the qsokufile's location"). Its actual
// value comes from the QSOKU_ROOT environment variable at run time.
const rootVar = `"$QSOKU_ROOT"/`

// quote state for Substitute's lexer.
const (
	normal = iota
	single
	double
)

// Substitute rewrites every "//" in a qsokufile command that stands for the
// qsokufile's own location to rootVar: outside quotes, at the start of a
// word (the start of the command, or right after whitespace, ';', '&', '|',
// or '('). Recognizing quotes only needs to track single quotes, double
// quotes, and a backslash escaping the next character.
//
// Once an unquoted '#' is reached at the start of a word, the rest of the
// command is a shell comment (the same start-of-word rule sh itself uses)
// and is copied through unchanged: sh never looks at it, so scanning it for
// "//" would rewrite text that has no effect either way.
func Substitute(command string) string {
	var b strings.Builder
	state := normal
	escape := false
	atWordStart := true

	for i := 0; i < len(command); i++ {
		c := command[i]

		if escape {
			b.WriteByte(c)
			escape = false
			atWordStart = false
			continue
		}

		switch state {
		case single:
			b.WriteByte(c)
			if c == '\'' {
				state = normal
			}
			continue
		case double:
			b.WriteByte(c)
			switch c {
			case '\\':
				escape = true
			case '"':
				state = normal
			}
			continue
		}

		// state == normal
		switch {
		case c == '\\':
			b.WriteByte(c)
			escape = true
			atWordStart = false
		case c == '\'':
			b.WriteByte(c)
			state = single
			atWordStart = false
		case c == '"':
			b.WriteByte(c)
			state = double
			atWordStart = false
		case c == '#' && atWordStart:
			b.WriteString(command[i:])
			return b.String()
		case c == '/' && atWordStart && i+1 < len(command) && command[i+1] == '/':
			b.WriteString(rootVar)
			i++ // the second '/' is part of this match too
			atWordStart = false
		case c == ' ' || c == '\t' || c == ';' || c == '&' || c == '|' || c == '(':
			b.WriteByte(c)
			atWordStart = true
		default:
			b.WriteByte(c)
			atWordStart = false
		}
	}
	return b.String()
}

// SubstituteArg rewrites a command-line argument that stands for the
// qsokufile's own location: an argument starting with "//" (already a
// single word, with quotes removed by the user's own shell before qsoku saw
// it) has that "//" replaced with root itself.
func SubstituteArg(arg, root string) string {
	if !strings.HasPrefix(arg, "//") {
		return arg
	}
	return root + arg[1:]
}
