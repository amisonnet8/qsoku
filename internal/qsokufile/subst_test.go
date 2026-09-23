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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SubstituteArg(tt.arg, root); got != tt.want {
				t.Errorf("SubstituteArg(%q, %q) = %q, want %q", tt.arg, root, got, tt.want)
			}
		})
	}
}
