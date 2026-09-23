package qsokufile

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []Entry
		wantErr string // substring expected in the error, "" if no error expected
	}{
		{
			name: "note.md example",
			input: "root: cd //\n" +
				"cd: cd $1\n" +
				"build: (cd //; make build)\n" +
				"test: (cd //; go test \"$@\")\n" +
				"\n" +
				"# Full clean rebuild.\n" +
				"rebuild: (cd //; rm -rf dist; make build)\n",
			want: []Entry{
				{Name: "root", Command: "cd //", Line: 1},
				{Name: "cd", Command: "cd $1", Line: 2},
				{Name: "build", Command: "(cd //; make build)", Line: 3},
				{Name: "test", Command: "(cd //; go test \"$@\")", Line: 4},
				{Name: "rebuild", Command: "(cd //; rm -rf dist; make build)", Line: 7},
			},
		},
		{
			name:  "blank lines are ignored",
			input: "\n   \n\t\nroot: cd //\n",
			want:  []Entry{{Name: "root", Command: "cd //", Line: 4}},
		},
		{
			name:  "comment with leading whitespace",
			input: "  # a comment\nroot: cd //\n",
			want:  []Entry{{Name: "root", Command: "cd //", Line: 2}},
		},
		{
			name:  "trailing comment on a command is not special",
			input: "greet: echo hi # not a qsoku comment\n",
			want:  []Entry{{Name: "greet", Command: "echo hi # not a qsoku comment", Line: 1}},
		},
		{
			name:  "leading whitespace before the whole line is trimmed",
			input: "  root:   cd //  \n",
			want:  []Entry{{Name: "root", Command: "cd //", Line: 1}},
		},
		{
			name:    "no colon",
			input:   "not a valid line\n",
			wantErr: "no ':'",
		},
		{
			name:    "name with a space is invalid",
			input:   "root : cd //\n",
			wantErr: `invalid name "root "`,
		},
		{
			name:    "name with an invalid character",
			input:   "ro ot: cd //\n",
			wantErr: "invalid name",
		},
		{
			name:    "name starting with a dot is invalid",
			input:   ".init: echo no\n",
			wantErr: "invalid name",
		},
		{
			name:  "name starting with a dash is allowed",
			input: "-x: echo ok\n",
			want:  []Entry{{Name: "-x", Command: "echo ok", Line: 1}},
		},
		{
			name:  "name with digits, dots, underscores and dashes",
			input: "build.v2_1-rc: echo ok\n",
			want:  []Entry{{Name: "build.v2_1-rc", Command: "echo ok", Line: 1}},
		},
		{
			name:    "empty name",
			input:   ": echo no\n",
			wantErr: `invalid name ""`,
		},
		{
			name:    "duplicate name",
			input:   "root: cd //\nroot: cd /tmp\n",
			wantErr: `"root" already defined at line 1`,
		},
		{
			name:  "CRLF line endings",
			input: "root: cd //\r\ncd: cd $1\r\n",
			want: []Entry{
				{Name: "root", Command: "cd //", Line: 1},
				{Name: "cd", Command: "cd $1", Line: 2},
			},
		},
		{
			name:  "empty file",
			input: "",
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := Parse(strings.NewReader(tt.input))

			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("Parse() = nil error, want error containing %q", tt.wantErr)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("Parse() error = %q, want it to contain %q", err.Error(), tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Parse() unexpected error: %v", err)
			}
			if len(f.Entries) != len(tt.want) {
				t.Fatalf("Parse() = %d entries, want %d: %+v", len(f.Entries), len(tt.want), f.Entries)
			}
			for i, got := range f.Entries {
				if got != tt.want[i] {
					t.Errorf("entry %d = %+v, want %+v", i, got, tt.want[i])
				}
			}
		})
	}
}
