//go:build e2e

package e2e

// The examples in the documents documentPairs lists (docs/reference/,
// README.md, docs/tour/, and their _ja counterparts) are run against the
// built qsoku binary, and what they show has to be what qsoku actually
// prints. Every example is a code block that starts with "$ qsoku ...",
// and the line before it says which fixture to run it against:
//
//	<!-- qsoku:example dir=basic -->
//
// dir names a directory of e2e/testdata/examples (or "none" for an empty
// directory with no qsokufile). cwd=<subdir> runs the commands from inside
// that subdirectory of the copy instead of its root, and skip="<why>" skips
// the example instead of running it. Every "$ " line must start with
// "qsoku " (or be exactly "qsoku"): these examples show qsoku's own output,
// not an arbitrary shell session.
//
// go test -tags e2e ./e2e -run '^TestDocExamples$' -update (make
// docs-examples) writes what qsoku prints into the documents. Nobody writes
// an example's output by hand.

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "write what qsoku prints into the documents that hold examples (make docs-examples)")

// documentPairs are the documents that hold examples, as paths relative to
// the repository root: the English document, then its _ja counterpart.
var documentPairs = [][2]string{
	{"docs/reference/cli.md", "docs/reference/cli_ja.md"},
	{"docs/reference/qsokufile.md", "docs/reference/qsokufile_ja.md"},
	{"README.md", "README_ja.md"},
	{"docs/tour/README.md", "docs/tour/README_ja.md"},
}

// documents flattens documentPairs into the list TestDocExamples walks.
var documents = flattenPairs(documentPairs)

func flattenPairs(pairs [][2]string) []string {
	all := make([]string, 0, len(pairs)*2)
	for _, p := range pairs {
		all = append(all, p[0], p[1])
	}
	return all
}

// placeholderPath is what every copy's absolute path is replaced with in an
// example's output, so the document does not depend on where the test ran.
const placeholderPath = "/home/you/project"

var markPattern = regexp.MustCompile(`^<!-- qsoku:example (.*) -->$`)

// example is a code block of a document that is run.
type example struct {
	line int // the line of the opening fence, from 1
	opts exampleOptions
	cmds []exampleCommand
}

type exampleOptions struct {
	dir  string
	cwd  string
	skip string
}

// exampleCommand is a "$ " line and the output the document shows after it,
// which is lines[outStart:outEnd] of the document (trailing blank lines
// that end it are left out).
type exampleCommand struct {
	text     string
	outStart int
	outEnd   int
}

// document is a file that has been read and split into examples.
type document struct {
	name     string
	lines    []string
	examples []*example
	problems []string
}

func loadDocument(t *testing.T, name string) *document {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("..", name)) //nolint:gosec // name is one of the fixed "documents" this file lists, not user input
	if err != nil {
		t.Fatal(err)
	}
	d := &document{name: name, lines: strings.Split(string(data), "\n")}
	d.parse()
	return d
}

func (d *document) problem(line int, format string, args ...any) {
	d.problems = append(d.problems, fmt.Sprintf("%s:%d: %s", d.name, line, fmt.Sprintf(format, args...)))
}

func (d *document) parse() {
	marked := map[int]bool{} // the marks that an example took
	for i := 0; i < len(d.lines); i++ {
		if !strings.HasPrefix(d.lines[i], "```") {
			continue
		}
		end := i + 1
		for end < len(d.lines) && !strings.HasPrefix(d.lines[end], "```") {
			end++
		}
		if end == len(d.lines) {
			d.problem(i+1, "a code block that is never closed")
			return
		}
		if end > i+1 && strings.HasPrefix(d.lines[i+1], "$ ") {
			d.example(i, end, marked)
		}
		i = end
	}
	for i, l := range d.lines {
		if strings.HasPrefix(l, "<!-- qsoku:example") && !marked[i] {
			d.problem(i+1, "an example mark that is not right before a code block that starts with $")
		}
	}
}

// example reads the block between the fences d.lines[open] and d.lines[end].
func (d *document) example(open, end int, marked map[int]bool) {
	ex := &example{line: open + 1}
	var m []string
	if open > 0 {
		m = markPattern.FindStringSubmatch(d.lines[open-1])
	}
	if m == nil {
		d.problem(ex.line, "an example (a code block that starts with $) without a line <!-- qsoku:example dir=... --> before it")
		return
	}
	marked[open-1] = true
	opts, err := parseOptions(m[1])
	if err != nil {
		d.problem(ex.line, "%v", err)
		return
	}
	ex.opts = opts

	for k := open + 1; k < end; k++ {
		if text, ok := strings.CutPrefix(d.lines[k], "$ "); ok {
			ex.cmds = append(ex.cmds, exampleCommand{text: text, outStart: k + 1})
		}
	}
	for i := range ex.cmds {
		stop := end
		if i+1 < len(ex.cmds) {
			stop = ex.cmds[i+1].outStart - 1
		}
		for stop > ex.cmds[i].outStart && strings.TrimSpace(d.lines[stop-1]) == "" {
			stop--
		}
		ex.cmds[i].outEnd = stop
	}
	if ex.opts.skip == "" {
		for _, c := range ex.cmds {
			if err := checkInvocation(c.text); err != nil {
				d.problem(ex.line, "%v", err)
			}
		}
	}
	d.examples = append(d.examples, ex)
}

// checkInvocation requires that a "$ " line runs qsoku itself, since these
// examples show qsoku's own output, not an arbitrary shell session.
func checkInvocation(text string) error {
	if text != "qsoku" && !strings.HasPrefix(text, "qsoku ") {
		return fmt.Errorf("a command of an example starts with qsoku: %q", text)
	}
	return nil
}

func parseOptions(mark string) (exampleOptions, error) {
	words, err := splitWords(mark)
	if err != nil {
		return exampleOptions{}, err
	}
	var o exampleOptions
	for _, w := range words {
		key, val, hasValue := strings.Cut(w, "=")
		switch {
		case key == "dir" && hasValue:
			o.dir = val
		case key == "cwd" && hasValue:
			o.cwd = val
		case key == "skip" && val != "":
			o.skip = val
		default:
			return exampleOptions{}, fmt.Errorf("cannot read %q in the mark of the example", w)
		}
	}
	if o.skip == "" && o.dir == "" {
		return exampleOptions{}, errors.New("an example needs dir=... (or skip=\"why\")")
	}
	return o, nil
}

// splitWords splits a line at spaces, the way a shell does for the simple
// cases: 'text' and "text" are one word, or a part of one, and nothing is
// an escape.
func splitWords(s string) ([]string, error) {
	var words []string
	var cur strings.Builder
	inWord := false
	var quote rune
	for _, r := range s {
		switch {
		case quote != 0:
			if r == quote {
				quote = 0
			} else {
				cur.WriteRune(r)
			}
		case r == '"' || r == '\'':
			quote, inWord = r, true
		case r == ' ' || r == '\t':
			if inWord {
				words = append(words, cur.String())
				cur.Reset()
				inWord = false
			}
		default:
			cur.WriteRune(r)
			inWord = true
		}
	}
	if quote != 0 {
		return nil, fmt.Errorf("a quote is not closed in %q", s)
	}
	if inWord {
		words = append(words, cur.String())
	}
	return words, nil
}

// TestDocExamplesAreMarkedAndMatch checks the documents without running
// anything: that no example is left without a mark, and that the English
// document and its _ja counterpart show the same examples (qsoku's output
// is English regardless, so both must show identical commands and output).
func TestDocExamplesAreMarkedAndMatch(t *testing.T) {
	docs := map[string]*document{}
	for _, name := range documents {
		d := loadDocument(t, name)
		for _, p := range d.problems {
			t.Error(p)
		}
		docs[name] = d
	}

	for _, pair := range documentPairs {
		en, ja := docs[pair[0]], docs[pair[1]]
		if len(en.examples) != len(ja.examples) {
			t.Fatalf("%s has %d examples and %s has %d: a change to one has to be made to the other", en.name, len(en.examples), ja.name, len(ja.examples))
		}
		for i := range en.examples {
			a, b := en.examples[i], ja.examples[i]
			if len(a.cmds) != len(b.cmds) {
				t.Errorf("example %d has %d commands in %s (line %d) and %d in %s (line %d)", i+1, len(a.cmds), en.name, a.line, len(b.cmds), ja.name, b.line)
				continue
			}
			for j := range a.cmds {
				if a.cmds[j].text != b.cmds[j].text {
					t.Errorf("example %d, command %d differs: %s:%d has %q, %s:%d has %q", i+1, j+1, en.name, a.line, a.cmds[j].text, ja.name, b.line, b.cmds[j].text)
				}
			}
		}
	}
}

// TestDocExamples runs every example against the built binary and compares
// what it prints with the document.
func TestDocExamples(t *testing.T) {
	root := resolvedTempDir(t)
	r := &exampleRunner{root: root}
	for _, name := range documents {
		d := loadDocument(t, name)
		if len(d.problems) > 0 {
			t.Fatalf("%s has %d problems; TestDocExamplesAreMarkedAndMatch lists them", name, len(d.problems))
		}
		var edits []edit
		for _, ex := range d.examples {
			t.Run(fmt.Sprintf("%s:%d", name, ex.line), func(t *testing.T) {
				if ex.opts.skip != "" {
					t.Skipf("skipped: %s", ex.opts.skip)
				}
				got := r.outputs(t, ex)
				for i, c := range ex.cmds {
					shown := strings.Join(d.lines[c.outStart:c.outEnd], "\n")
					if shown == got[i] {
						continue
					}
					edits = append(edits, edit{c.outStart, c.outEnd, lines(got[i])})
					if !*update {
						t.Errorf("`$ %s` prints something else than the document shows (make docs-examples writes what it prints into the document):\n--- document\n%s\n--- qsoku\n%s", c.text, shown, got[i])
					}
				}
			})
		}
		if *update && len(edits) > 0 {
			d.apply(t, edits)
		}
	}
}

func lines(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}

// edit replaces lines[start:end] of a document.
type edit struct {
	start, end int
	with       []string
}

// apply writes the edits into the document. The rest of the file is not touched.
func (d *document) apply(t *testing.T, edits []edit) {
	t.Helper()
	sort.Slice(edits, func(i, j int) bool { return edits[i].start > edits[j].start })
	out := slices.Clone(d.lines)
	for _, e := range edits {
		out = slices.Concat(out[:e.start], e.with, out[e.end:])
	}
	path := filepath.Join("..", d.name)
	if err := os.WriteFile(path, []byte(strings.Join(out, "\n")), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Logf("wrote %d examples into %s", len(edits), path)
}

// exampleRunner runs examples in fresh copies of the fixtures of
// e2e/testdata/examples.
type exampleRunner struct {
	root   string
	copies int
}

// outputs runs the commands of an example, in a fresh copy of its fixture,
// and returns what each printed (standard output, then standard error),
// with the copy's absolute path replaced by placeholderPath.
func (r *exampleRunner) outputs(t *testing.T, ex *example) []string {
	t.Helper()
	dir := r.copyOf(t, ex.opts.dir)
	runDir := dir
	if ex.opts.cwd != "" {
		runDir = filepath.Join(dir, ex.opts.cwd)
	}
	var got []string
	for _, c := range ex.cmds {
		out := r.run(t, runDir, c.text)
		out = strings.ReplaceAll(out, dir, placeholderPath)
		got = append(got, strings.TrimRight(out, "\n"))
	}
	return got
}

// run runs one "$ " line through sh, with the built qsoku on PATH.
func (r *exampleRunner) run(t *testing.T, dir, text string) string {
	t.Helper()
	home := t.TempDir()
	cmd := exec.Command("sh", "-c", text) //nolint:gosec // text is a "$ qsoku ..." line from docs/reference/, checked by checkInvocation to start with qsoku, run against a disposable fixture copy
	cmd.Dir = dir
	cmd.Env = []string{
		"PATH=" + filepath.Dir(binary) + string(os.PathListSeparator) + os.Getenv("PATH"),
		"HOME=" + home,
	}
	var out, errOut strings.Builder
	cmd.Stdout, cmd.Stderr = &out, &errOut
	_ = cmd.Run() // the example's own exit code is not shown in the document
	return out.String() + errOut.String()
}

// copyOf makes a fresh copy of a fixture's directory (or an empty directory
// for "none"), and returns its path. Every example gets its own copy, so
// what one writes (an .add, say) is not there for the next.
func (r *exampleRunner) copyOf(t *testing.T, fixture string) string {
	t.Helper()
	r.copies++
	dst := filepath.Join(r.root, fmt.Sprintf("run-%d", r.copies), "project")
	if err := os.MkdirAll(dst, 0o750); err != nil {
		t.Fatal(err)
	}
	if fixture == "none" {
		return dst
	}
	src := filepath.Join("testdata", "examples", fixture)
	if _, err := os.Stat(src); err != nil {
		t.Fatalf("fixture %q: %v", fixture, err)
	}
	err := filepath.WalkDir(src, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, p)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		if d.IsDir() {
			return os.MkdirAll(filepath.Join(dst, rel), info.Mode().Perm()|0o700)
		}
		data, err := os.ReadFile(p) //nolint:gosec // p is walked from e2e/testdata/examples/<fixture>, a fixed test fixture tree, not user input
		if err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(dst, rel), data, info.Mode().Perm()|0o600) //nolint:gosec // dst is this test's own t.TempDir()-derived copy, and rel comes from the same fixed fixture walk above
	})
	if err != nil {
		t.Fatal(err)
	}
	return dst
}
