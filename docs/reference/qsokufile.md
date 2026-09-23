# qsokufile

*[日本語](qsokufile_ja.md) | **English***

`qsokufile` is the single file that holds a repository's command shortcuts.
This document specifies its format. For the `qsoku` command line itself, see
[cli.md](cli.md).

## Location and lookup

- The file must be named exactly `qsokufile`, lowercase.
- qsoku looks for it starting at the current directory and walks up through
  each parent directory, using the **first** `qsokufile` it finds. It does not
  stop at a git repository boundary, and does not look above the filesystem
  root.
- On a case-insensitive filesystem (macOS, Windows), a file named `Qsokufile`
  or `QSOKUFILE` would also be found. Always create and edit it as
  `qsokufile` (`qsoku .init` always writes that exact name) so the name stays
  consistent across machines.
- If no `qsokufile` is found, qsoku exits with an error (exit code 1) and
  suggests `qsoku .init`.

## Format

- One entry per line: `name: command`.
- The line is trimmed of leading and trailing whitespace first. The first `:`
  on the line separates the name from the command; everything after it, with
  leading whitespace trimmed, is the command. The command keeps its own
  internal spacing.
- A trailing `\r` (CRLF line ending) is stripped before parsing, so a file
  edited on Windows still reads correctly.
- Blank lines (empty, or whitespace only) are ignored.
- **Comment lines**: a line is a comment, and ignored, when its first
  non-whitespace character is `#`. Leading whitespace before the `#` is
  allowed. A `#` that appears anywhere else on a line (for example at the end
  of a command) is **not** special to qsoku — it is part of the command text
  handed to `sh`, which interprets it as a shell comment itself.
- A non-blank, non-comment line that has no `:` is an error: qsoku reports the
  file and line number and exits (code 1) without running anything.
- There is no line continuation; each entry is exactly one line.

### Example: a bad line

<!-- qsoku:example dir=broken -->
```
$ qsoku ok
qsoku: /home/you/project/qsokufile: line 2: no ':' found (expected "name: command")
```

## Names

- Allowed characters: ASCII letters and digits, `.`, `_`, `-`.
- A name must not **start** with `.` — that is reserved for qsoku's own
  management commands (`qsoku .init`, `qsoku .add`, and so on; see
  [cli.md](cli.md)).
- A name **may** start with `-`. qsoku has no command-line options of its
  own — every argument after `qsoku` is looked up as a name — so there is no
  ambiguity to avoid.
- Names are case-sensitive.
- A name must appear at most once in a `qsokufile`. A second line defining a
  name already defined earlier is an error (both line numbers are reported;
  exit code 1). `qsoku .add` never creates this situation: adding a name that
  already exists replaces that line instead (see [cli.md](cli.md)).

### Example: a duplicate name

<!-- qsoku:example dir=dup -->
```
$ qsoku build
qsoku: /home/you/project/qsokufile: line 2: name "build" already defined at line 1
```

## `//`: the qsokufile's location

`//` stands for the directory that holds the `qsokufile` in use (also called
`QSOKU_ROOT`; see [cli.md](cli.md)). A real root path, `/`, is written as `/`
and is never affected by this rule.

| Written as | Substituted? | Why |
|---|---|---|
| `cd //src` | Yes | Start of the command, so start of a word |
| `(cd //; make build)` | Yes | Start of a word: right after `(` and after `;` |
| `foo \| cd //` | Yes | Start of a word: right after `\|` |
| `foo && cd //` | Yes | Start of a word: right after `&` |
| `curl http://example.com` | No | After `:`, not at the start of a word |
| `grep "//" main.go` | No | Inside double quotes |
| `echo '//' is a comment` | No | Inside single quotes |
| `build: (cd //; make build) # see //docs` | No (the `# see //docs` part) | That part of the line is a comment (see Format above), never reached |

- **Where it substitutes**: outside quotes, at the start of a word — the very
  start of the command, or right after whitespace, `;`, `&`, `|`, or `(`.
- **How**: qsoku passes the qsokufile's directory to `sh` as the environment
  variable `QSOKU_ROOT`, and rewrites each matching `//` to `"$QSOKU_ROOT"/`
  before running the command. This keeps the substitution correct even when
  the path contains spaces.
- **Escaping**: to use a literal `//` (for example inside a URL or a
  comment), put it in single or double quotes, or write it inside a
  comment — nothing new is invented for this. To use the qsokufile's location
  *inside* quotes, write `$QSOKU_ROOT` directly.
- **Arguments typed by the user**: `qsoku cd //src` passes `//src` as an
  argument to the `cd` entry's command (after the user's own shell has
  already removed any quotes). qsoku substitutes it the same way, so `$1`
  inside the qsokufile's command receives `/path/to/qsokufile/dir/src`.
- Implementation note: recognizing quotes needs only to track whether the
  lexer is inside single quotes, double quotes, or after a backslash — full
  shell grammar is not required.

### Substituting, and not

Run against a `qsokufile` with `show: echo "$1"`, `url: echo http://example.com`,
`grepq: grep "//" main.go` (a file containing the line `// see docs`) and
`echoq: echo '//' is a comment`:

<!-- qsoku:example dir=slashes -->
```
$ qsoku url
http://example.com
$ qsoku grepq
// see docs
$ qsoku echoq
// is a comment
$ qsoku show //src
/home/you/project/src
```

`url` and `echoq` print their `//` back unchanged (not at the start of a
word, and inside single quotes, respectively). `grepq`'s `grep "//" main.go`
still finds the literal `//` in `main.go` — proof it wasn't substituted
inside the double quotes. `show //src`, an argument typed by the user, is
substituted to the fixture's own directory.

## Example

```
root: cd //
cd: cd $1
build: (cd //; make build)
test: (cd //; go test "$@")

# Full clean rebuild.
rebuild: (cd //; rm -rf dist; make build)
```
