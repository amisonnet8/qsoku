//go:build e2e

package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"testing"
)

// newShellRepo creates a qsokufile with three entries that exercise the
// working-directory handoff (docs/reference/cli.md "Bringing the working
// directory back"): a plain cd, a plain cd followed by a failure, and the
// same inside parentheses. It returns the qsokufile's directory and the
// "sub" directory the entries cd into.
func newShellRepo(t *testing.T) (dir, sub string) {
	t.Helper()
	// Canonicalized up front (filepath.EvalSymlinks), not just t.TempDir()'s
	// raw value: on macOS this resolves the /tmp -> /private/tmp symlink,
	// and on Windows it normalizes a short 8.3-style path segment (e.g.
	// "RUNNER~1") to its real long name -- both forms a real shell's own
	// "pwd", after cd'ing here, reports back regardless of what t.TempDir()
	// itself happened to return.
	dir = resolveSymlinks(t, t.TempDir())
	sub = filepath.Join(dir, "sub")
	if err := os.Mkdir(sub, 0o750); err != nil {
		t.Fatal(err)
	}
	content := "root: cd //sub\n" +
		"failcd: cd //sub; sh -c \"exit 3\"\n" +
		"parens: (cd //sub; sh -c \"exit 5\")\n"
	if err := os.WriteFile(filepath.Join(dir, "qsokufile"), []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return dir, sub
}

// shellPath finds a shell, or skips the test.
func shellPath(t *testing.T, name string) string {
	t.Helper()
	path, err := exec.LookPath(name)
	if err != nil {
		t.Skipf("%s is not installed", name)
	}
	return path
}

// shellEnv is the environment of a shell that runs qsoku from the PATH.
func shellEnv() []string {
	path := filepath.Dir(binary) + string(os.PathListSeparator) + os.Getenv("PATH")
	return append(os.Environ(), "PATH="+path)
}

// runShellScript runs a program with a script on its standard input, in
// dir, and returns what it printed on standard output.
func runShellScript(t *testing.T, dir, script, program string, args ...string) string {
	t.Helper()
	cmd := exec.Command(program, args...) //nolint:gosec // program is a shell binary found via exec.LookPath in this same test file, not user input
	cmd.Dir = dir
	cmd.Env = shellEnv()
	cmd.Stdin = strings.NewReader(script)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s: %v\n%s", program, err, out)
	}
	return string(out)
}

// cwdCase is one of the three qsokufile entries in newShellRepo.
type cwdCase struct {
	name       string
	wantStatus int
	wantMoved  bool // whether the shell should end up in "sub"
}

var cwdCases = []cwdCase{
	{"root", 0, true},
	{"failcd", 3, true},
	{"parens", 5, false},
}

// parseCwdOutput reads len(cwdCases) pairs of "STATUS:n" then a pwd line.
func parseCwdOutput(t *testing.T, out string) []struct {
	status int
	pwd    string
} {
	t.Helper()
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) != len(cwdCases)*2 {
		t.Fatalf("expected %d lines, got %d:\n%s", len(cwdCases)*2, len(lines), out)
	}
	results := make([]struct {
		status int
		pwd    string
	}, len(cwdCases))
	for i := range cwdCases {
		statusLine, pwdLine := lines[i*2], lines[i*2+1]
		n, ok := strings.CutPrefix(statusLine, "STATUS:")
		if !ok {
			t.Fatalf("line %q does not start with STATUS:", statusLine)
		}
		status, err := strconv.Atoi(n)
		if err != nil {
			t.Fatalf("bad status line %q: %v", statusLine, err)
		}
		results[i].status = status
		results[i].pwd = pwdLine
	}
	return results
}

func checkCwdResults(t *testing.T, dir, sub string, results []struct {
	status int
	pwd    string
}) {
	t.Helper()
	for i, c := range cwdCases {
		want := dir
		if c.wantMoved {
			want = sub
		}
		if results[i].status != c.wantStatus {
			t.Errorf("%s: status = %d, want %d", c.name, results[i].status, c.wantStatus)
		}
		if results[i].pwd != want {
			t.Errorf("%s: pwd = %q, want %q", c.name, results[i].pwd, want)
		}
	}
}

func quote(s string) string { return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'" }

func quoteFish(s string) string {
	return "'" + strings.ReplaceAll(strings.ReplaceAll(s, `\`, `\\`), "'", `\'`) + "'"
}

// posixCwdScript builds a script usable by both bash and zsh (their
// function syntax is compatible): eval the integration, then for each case
// cd back to dir, run the entry, and print its status and the resulting pwd.
func posixCwdScript(shell, dir string) string {
	var b strings.Builder
	if shell == "zsh" {
		// zsh's own completion system (compinit) is not loaded in this bare
		// -f session; qsoku.zsh's trailing `compdef` call would otherwise
		// fail with "command not found" (mtqg's e2e stubs it the same way).
		b.WriteString("compdef() { :; }\n")
	}
	b.WriteString("eval \"$(command qsoku .shell " + shell + ")\"\n")
	for _, c := range cwdCases {
		b.WriteString("cd " + quote(dir) + "\n")
		b.WriteString("qsoku " + c.name + "\n")
		b.WriteString("echo \"STATUS:$?\"\n")
		b.WriteString("pwd\n")
	}
	return b.String()
}

// toMsysPath converts a native Windows path (whatever separator or drive-
// letter case it happens to use) to the form Git for Windows' MSYS bash/zsh
// report back from their own "pwd" builtin (/c/Users/... rather than
// C:/Users/... or C:\Users\...): MSYS mounts each drive under a lowercase
// single-letter root of its own POSIX-style namespace, regardless of what
// was passed to "cd" to get there. A no-op everywhere else, where "pwd"
// already reports back exactly what was cd'd to.
func toMsysPath(p string) string {
	if runtime.GOOS != "windows" {
		return p
	}
	p = filepath.ToSlash(p)
	if len(p) >= 2 && p[1] == ':' {
		return "/" + strings.ToLower(p[:1]) + p[2:]
	}
	return p
}

func TestShellIntegrationBash(t *testing.T) {
	bash := shellPath(t, "bash")
	dir, sub := newShellRepo(t)

	out := runShellScript(t, dir, posixCwdScript("bash", dir), bash, "--norc", "--noprofile")
	checkCwdResults(t, toMsysPath(dir), toMsysPath(sub), parseCwdOutput(t, out))
}

func TestShellIntegrationZsh(t *testing.T) {
	zsh := shellPath(t, "zsh")
	dir, sub := newShellRepo(t)

	out := runShellScript(t, dir, posixCwdScript("zsh", dir), zsh, "-f")
	checkCwdResults(t, toMsysPath(dir), toMsysPath(sub), parseCwdOutput(t, out))
}

func TestShellIntegrationFish(t *testing.T) {
	fish := shellPath(t, "fish")
	dir, sub := newShellRepo(t)

	var b strings.Builder
	b.WriteString("command qsoku .shell fish | source\n")
	for _, c := range cwdCases {
		b.WriteString("cd " + quoteFish(dir) + "\n")
		b.WriteString("qsoku " + c.name + "\n")
		b.WriteString("echo \"STATUS:$status\"\n")
		b.WriteString("pwd\n")
	}
	out := runShellScript(t, dir, b.String(), fish, "--no-config")
	checkCwdResults(t, dir, sub, parseCwdOutput(t, out))
}

// wantCandidates is what completion should offer after "qsoku ": the three
// entries from newShellRepo, and the fixed management commands.
var wantCandidates = []string{
	"root", "failcd", "parens",
	".init", ".add", ".rm", ".list", ".names", ".edit", ".where", ".version", ".shell", ".help",
}

func checkCandidates(t *testing.T, got []string) {
	t.Helper()
	got = slices.Clone(got)
	slices.Sort(got)
	want := slices.Clone(wantCandidates)
	slices.Sort(want)
	if !slices.Equal(got, want) {
		t.Errorf("candidates:\n got  %q\n want %q", got, want)
	}
}

func TestCompletionInBash(t *testing.T) {
	bash := shellPath(t, "bash")
	dir, _ := newShellRepo(t)

	script := "eval \"$(command qsoku .shell bash)\"\n" +
		"COMP_WORDS=(qsoku '')\n" +
		"COMP_CWORD=1\n" +
		"_qsoku\n" +
		"printf '%s\\n' \"${COMPREPLY[@]}\"\n"
	out := runShellScript(t, dir, script, bash, "--norc", "--noprofile")
	checkCandidates(t, strings.Split(strings.TrimRight(out, "\n"), "\n"))
}

func TestCompletionInZsh(t *testing.T) {
	zsh := shellPath(t, "zsh")
	dir, _ := newShellRepo(t)

	// _qsoku can be called directly once defined. compdef is stubbed for the
	// same reason as in TestShellIntegrationZsh: no compinit in this bare
	// session.
	script := "compdef() { :; }\n" +
		"eval \"$(command qsoku .shell zsh)\"\n" +
		"words=(qsoku '')\n" +
		"CURRENT=2\n" +
		"reply=()\n" +
		"compadd() { shift; reply+=(\"$@\"); }\n" +
		"_qsoku\n" +
		"printf '%s\\n' \"${reply[@]}\"\n"
	out := runShellScript(t, dir, script, zsh, "-f")
	checkCandidates(t, strings.Split(strings.TrimRight(out, "\n"), "\n"))
}

// fish hides completion candidates that start with '.' until the word being
// completed itself starts with '.' -- the same rule it uses for dotfiles in
// path completion, applied here too since qsoku's management commands are
// spelled that way (docs/reference/cli.md "Shell completion"). So "qsoku "
// only offers the qsokufile's own names, and "qsoku ." is needed to see the
// management commands -- unlike bash and zsh, which offer both at once.
func TestCompletionInFish(t *testing.T) {
	fish := shellPath(t, "fish")
	dir, _ := newShellRepo(t)

	complete := func(line string) []string {
		script := "command qsoku .shell fish | source\n" +
			"complete -C " + quoteFish(line) + "\n"
		out := runShellScript(t, dir, script, fish, "--no-config")
		var got []string
		for _, l := range strings.Split(strings.TrimRight(out, "\n"), "\n") {
			if l == "" {
				continue
			}
			value, _, _ := strings.Cut(l, "\t")
			got = append(got, value)
		}
		return got
	}

	t.Run("names", func(t *testing.T) {
		want := []string{"root", "failcd", "parens"}
		got := complete("qsoku ")
		slices.Sort(got)
		slices.Sort(want)
		if !slices.Equal(got, want) {
			t.Errorf("candidates for \"qsoku \":\n got  %q\n want %q", got, want)
		}
	})

	t.Run("management commands", func(t *testing.T) {
		want := []string{".init", ".add", ".rm", ".list", ".names", ".edit", ".where", ".version", ".shell", ".help"}
		got := complete("qsoku .")
		slices.Sort(got)
		slices.Sort(want)
		if !slices.Equal(got, want) {
			t.Errorf("candidates for \"qsoku .\":\n got  %q\n want %q", got, want)
		}
	})
}

// quotePwsh quotes a string as a PowerShell single-quoted string literal
// (only ' itself needs doubling; unlike sh, PowerShell does no other
// expansion inside single quotes).
func quotePwsh(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

// ansiEscape matches a terminal control sequence (CSI, e.g. "\x1b[?1h").
// pwsh emits these around running a native command (qsoku, here) even when
// its own standard output is a pipe rather than a real terminal -- observed
// on both Linux and macOS CI runners, not just interactively. cleanPwshOutput
// strips them and drops any line that turns out to be empty once they are
// gone (a sequence that landed on a line of its own), so the real output
// lines this package's other parsing helpers expect line up exactly as if
// pwsh had never touched the terminal at all.
var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;?]*[a-zA-Z]`)

// cleanPwshOutput also strips '\r': on Windows, pwsh's Write-Output ends
// lines with the platform's own newline (CRLF), which this package's other
// parsing helpers -- built around bash/zsh/fish's plain '\n' -- do not
// expect.
func cleanPwshOutput(s string) string {
	s = ansiEscape.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, "\r", "")
	lines := strings.Split(s, "\n")
	kept := lines[:0]
	for _, l := range lines {
		if l != "" {
			kept = append(kept, l)
		}
	}
	return strings.Join(kept, "\n")
}

// resolveSymlinks resolves a path (typically t.TempDir()'s raw value) to
// its canonical form: on macOS this follows /tmp -> /private/tmp, and on
// Windows it normalizes a short 8.3-style path segment (e.g. "RUNNER~1") to
// its real long name. newShellRepo uses this for the directory it returns,
// since a real shell reports back one of these canonical forms regardless
// of what was passed to "cd" to get there -- notably pwsh, whose
// Set-Location/current-directory tracking does not preserve a "logical",
// possibly-symlinked or short-named path the way a POSIX shell's own $PWD
// does.
func resolveSymlinks(t *testing.T, dir string) string {
	t.Helper()
	resolved, err := filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatal(err)
	}
	return resolved
}

func TestShellIntegrationPwsh(t *testing.T) {
	pwsh := shellPath(t, "pwsh")
	dir, sub := newShellRepo(t)

	var b strings.Builder
	b.WriteString("Invoke-Expression (& qsoku .shell pwsh | Out-String)\n")
	for _, c := range cwdCases {
		b.WriteString("Set-Location -LiteralPath " + quotePwsh(dir) + "\n")
		b.WriteString("qsoku " + c.name + "\n")
		b.WriteString("Write-Output (\"STATUS:\" + $LASTEXITCODE)\n")
		b.WriteString("Write-Output (Get-Location).Path\n")
	}
	out := runShellScript(t, dir, b.String(), pwsh, "-NoProfile", "-NoLogo", "-Command", "-")
	checkCwdResults(t, dir, sub, parseCwdOutput(t, cleanPwshOutput(out)))
}

func TestCompletionInPwsh(t *testing.T) {
	pwsh := shellPath(t, "pwsh")
	dir, _ := newShellRepo(t)

	script := "Invoke-Expression (& qsoku .shell pwsh | Out-String)\n" +
		"$r = TabExpansion2 -inputScript 'qsoku ' -cursorColumn 6\n" +
		"$r.CompletionMatches | ForEach-Object { $_.CompletionText }\n"
	out := runShellScript(t, dir, script, pwsh, "-NoProfile", "-NoLogo", "-Command", "-")
	out = cleanPwshOutput(out)
	checkCandidates(t, strings.Split(strings.TrimRight(out, "\n"), "\n"))
}
