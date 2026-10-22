package qsokufile

import "testing"

func TestSubstitute(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    string
	}{
		// docs/reference/qsokufile.md "// : the qsokufile's location" table.
		{"start of command", "cd //src", `cd "$QSOKU_ROOT"/src`},
		{"after ( and after ;", "(cd //; make build)", `(cd "$QSOKU_ROOT"/; make build)`},
		{"after |", "foo | cd //", `foo | cd "$QSOKU_ROOT"/`},
		{"after &", "foo && cd //", `foo && cd "$QSOKU_ROOT"/`},
		{"after : is not a word start", "curl http://example.com", "curl http://example.com"},
		{"inside double quotes", `grep "//" main.go`, `grep "//" main.go`},
		{"inside single quotes", "echo '//' is a comment", "echo '//' is a comment"},
		// The two above start the quoted text with //, which is not at a word
		// start once the quote is read; these have // at a word start inside
		// the quotes, so only the quote state keeps them from being replaced.
		{"a word start inside single quotes", "echo 'a //b'", "echo 'a //b'"},
		{"a word start inside double quotes", `echo "a //b"`, `echo "a //b"`},
		{"right after ;", "x;//y", `x;"$QSOKU_ROOT"/y`},
		{"right after (", "(//bin/x)", `("$QSOKU_ROOT"/bin/x)`},
		{
			"a trailing comment is never reached",
			`(cd //; make build) # see //docs`,
			`(cd "$QSOKU_ROOT"/; make build) # see //docs`,
		},

		// Additional cases.
		{
			"an escaped double quote does not close the string",
			`echo "x\" //y" z`,
			`echo "x\" //y" z`,
		},
		{"an escaped quote does not start a quoted string", `echo \'foo //bar`, `echo \'foo "$QSOKU_ROOT"/bar`},
		{"a third slash is left alone", "cd ///x", `cd "$QSOKU_ROOT"//x`},
		{"not at a word start mid-word", "a//b", "a//b"},
		{"a lone slash is not substituted", "cd /x", "cd /x"},
		{"empty command", "", ""},
		{"comment as the whole command", "# //not substituted", "# //not substituted"},
		{"an escaped hash does not start a comment", `echo \#foo //bar`, `echo \#foo "$QSOKU_ROOT"/bar`},
		{"tab counts as whitespace", "cd\t//x", "cd\t\"$QSOKU_ROOT\"/x"},
		// A quoted or escaped part starts the word: // right after it is mid-word.
		{"mid-word after an escaped character", `echo a\x//y`, `echo a\x//y`},
		{"after an empty single-quoted string", "echo ''//x", "echo ''//x"},
		{"after an empty double-quoted string", `echo ""//x`, `echo ""//x`},
		{"a slash at the end", "cd /", "cd /"},
		{"four slashes: only the first two", "cd ////x", `cd "$QSOKU_ROOT"///x`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Substitute(tt.command); got != tt.want {
				t.Errorf("Substitute(%q) = %q, want %q", tt.command, got, tt.want)
			}
		})
	}
}

func TestSubstituteArg(t *testing.T) {
	const root = "/repo"

	tests := []struct {
		name string
		arg  string
		want string
	}{
		{"prefixed path", "//src", "/repo/src"},
		{"bare //", "//", "/repo/"},
		{"no prefix", "src", "src"},
		{"single slash is not a prefix", "/src", "/src"},
		{"empty argument", "", ""},
		{"a third slash is left alone", "///x", "/repo//x"},
		{"// not at the start", "a//b", "a//b"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SubstituteArg(tt.arg, root); got != tt.want {
				t.Errorf("SubstituteArg(%q, %q) = %q, want %q", tt.arg, root, got, tt.want)
			}
		})
	}
}
