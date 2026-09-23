# qsoku command reference

*[日本語](cli_ja.md) | **English***

```
qsoku <name> [args...]
```

Every argument after `qsoku` is looked up as a name. **qsoku has no
command-line options of its own** (no `-h`, `--help`, `--version`, and so
on) — a name may start with `-` (see [qsokufile.md](qsokufile.md#names)),
and treating it as an option would make that name unreachable. All of
qsoku's own functionality is under names starting with `.` (see
[Management commands](#management-commands)), which a `qsokufile` entry can
never define (see [qsokufile.md](qsokufile.md#names)).

For the file format itself, see [qsokufile.md](qsokufile.md).

## Running a name

1. qsoku finds the `qsokufile` in use (see
   [qsokufile.md](qsokufile.md#location-and-lookup)) and looks up `<name>`.
   If it is not defined, qsoku reports an error and exits (code 2 — see
   [Exit codes](#exit-codes)).
2. qsoku builds the command's text, substituting `//` per
   [qsokufile.md](qsokufile.md#-the-qsokufiles-location), and runs it with
   `sh`, from the current directory, as:

   ```sh
   sh -c '<command>
   __status=$?; pwd > "$QSOKU_CWD_FILE"; exit $__status' qsoku "$@"
   ```

   — but only when `QSOKU_CWD_FILE` is set in qsoku's own environment (see
   [Bringing the working directory back](#bringing-the-working-directory-back)
   below); otherwise it is just `sh -c '<command>' qsoku "$@"`, with no
   second line.

   - The `pwd` line, when present, is joined with a **newline**, not `;` —
     a trailing `# …` comment on the command would otherwise swallow it.
   - `"$@"` carries `[args...]` through to the command as `$1`, `$2`, and so
     on, exactly as `sh` itself would expand them.
   - Standard input, standard output and standard error are passed straight
     through; qsoku does not read or alter them.
   - A `qsokufile` entry should not itself call `exit`: since the second
     line is appended to the *same* script, an `exit` inside the entry's own
     command ends the whole script right there, before qsoku's own line
     ever runs. Let the entry's last command's own exit status decide
     success or failure instead (as `make` recipes and ordinary shell
     scripts do), the same way `cd //src && make build` in the example below
     does.
3. `sh` is located via `PATH`, like any other subprocess.

### Example

<!-- qsoku:example dir=basic -->
```
$ qsoku hello world
hello, world
$ qsoku nope
qsoku: "nope" is not defined in /home/you/project/qsokufile
```

Looked up from a subdirectory, the same `qsokufile` (found by walking up) still applies:

<!-- qsoku:example dir=basic cwd=src -->
```
$ qsoku hello everyone
hello, everyone
```

### Environment passed to the command

| Variable | Value |
|---|---|
| `QSOKU_ROOT` | The absolute path of the directory holding the `qsokufile` in use (what `//` expands to) |
| `QSOKU_CWD_FILE` | Set by the shell integration function (see [Shell integration](#shell-integration)) to a temporary file's path. When set, the command above appends the `pwd` line. When unset (for example, qsoku run directly without the shell function), the `pwd` line is left out entirely — appending it regardless would redirect `pwd`'s output to an empty filename, which `sh` reports as an error, for no benefit: nothing would read the file anyway |

### Bringing the working directory back

Because the command runs in a separate process, a `cd` inside it does not
move the caller's shell. The shell integration function reads the location
`qsoku` finishes in from `QSOKU_CWD_FILE` and `cd`s there itself:

- The location is written to `QSOKU_CWD_FILE` **whether the command
  succeeds or fails**, and whether the process substitution itself succeeded — this
  is not conditional. Examples:
  - `cd //src` failing (no such directory) leaves `pwd` unchanged, so
    nothing moves when it's read back.
  - `(cd //; make build)` moves inside a subshell, so `pwd` is unaffected
    either way (see [qsokufile.md](qsokufile.md#-the-qsokufiles-location)).
  - `cd //src && make build` written *without* parentheses, with a build
    failure, brings back whatever directory that reached — the entry was
    written without parentheses, so moving there is taken to be intended
    even on failure.
- The location written is the **logical** working directory (`pwd`, not
  `pwd -P`): qsoku follows a symlinked path the same way an interactive
  shell's own `cd` does, rather than resolving it away.
- When qsoku's own errors happen **before** the command runs (see
  [Exit codes](#exit-codes)) — the name isn't found, the `qsokufile` can't be
  found or parsed, `sh` can't be started — nothing is written to
  `QSOKU_CWD_FILE`. The shell function treats an empty read as "don't move".

## Management commands

Management commands start with `.`, a prefix a `qsokufile` name can never
have (see [qsokufile.md](qsokufile.md#names)), so they never collide with a
user's own names.

| Command | Does |
|---|---|
| `qsoku .init` | Creates an empty `qsokufile` in the current directory. Error if one already exists there (a `qsokufile` in a parent directory does not block this) |
| `qsoku .add <name> <command>` | Adds `name: command` to the `qsokufile` in use. If `<name>` is already defined, **replaces that line in place** (same position; comments and the rest of the order are untouched) instead of erroring |
| `qsoku .rm <name>` | Removes that one line. Error if `<name>` is not defined |
| `qsoku .list` | Lists the defined names and their commands, one per line, in file order |
| `qsoku .names` | Lists just the names, one per line — nothing else, ever (see [Shell completion](#shell-completion)) |
| `qsoku .edit` | Opens the `qsokufile` in use with `$EDITOR`; if unset, falls back to `nano`; if `nano` is not on `PATH` either, errors (exit 1) instead of guessing further |
| `qsoku .where` | Prints the absolute path of the directory holding the `qsokufile` in use (what `//` and `QSOKU_ROOT` expand to) |
| `qsoku .version` | Prints qsoku's own version, from `runtime/debug.ReadBuildInfo` (see [distribution.md](../../.claude/rules/distribution.md) design note — this is why `go install` alone, without `-ldflags`, still reports a meaningful version) |
| `qsoku .shell <shell>` | Prints shell code for the given shell (`bash`, `zsh`, `fish`) that defines the `qsoku` function described below, to `eval`. Unknown `<shell>` is a command-line error (exit 2) |
| `qsoku .help` | Prints usage and the list of defined names. Running `qsoku` with **no arguments at all** prints the same thing (the `git`/`git --help` convention, without an option to spell) |

- `.add` and `.rm` write to the `qsokufile` found by walking up from the
  current directory, same as running a name. If none is found, they error
  instead of creating one — `.init` is suggested.
- `.add`'s rewrite is minimal: only the one line changes (added at the end,
  or replaced in place). Hand-written comments and the order of other
  entries are never touched.
- `.edit` and `.where` work even when the `qsokufile` currently fails to
  parse (a bad line, a duplicate name) — you need them to be able to fix it.
- `.names`, specifically, **never fails**: given no `qsokufile`, an
  unreadable one, or a `qsokufile` with a parse error, it prints nothing and
  exits 0, with nothing on standard error. This is so that pressing Tab in
  the middle of typing a command never interrupts the line with an error
  (see [Shell completion](#shell-completion)). Every other management
  command reports such problems normally.

### Example

<!-- qsoku:example dir=none -->
```
$ qsoku .init
Created /home/you/project/qsokufile
$ qsoku .add build 'go build ./...'
$ qsoku .add test 'go test ./...'
$ qsoku .list
build: go build ./...
test: go test ./...
$ qsoku .names
build
test
$ qsoku .where
/home/you/project
$ qsoku .rm test
$ qsoku .list
build: go build ./...
```

## Shell integration

`.bashrc` (or the equivalent for your shell) needs one line:

```sh
eval "$(qsoku .shell bash)"
```

| Shell | Typical file |
|---|---|
| bash | `~/.bashrc` |
| zsh | `~/.zshrc` |
| fish | `~/.config/fish/config.fish` (using `qsoku .shell fish \| source`, fish's own idiom) |

This defines a shell function named `qsoku` (not a single letter — short
enough already, and `command qsoku` inside it reaches the real binary
without recursing):

```sh
qsoku() {
  local f; f=$(mktemp)
  QSOKU_CWD_FILE=$f command qsoku "$@"
  local status=$?
  local dir; dir=$(cat "$f"); rm -f "$f"
  [ -n "$dir" ] && [ "$dir" != "$PWD" ] && cd "$dir"
  return $status
}
```

(fish's equivalent function is spelled differently but does the same three
things: run the real binary with `QSOKU_CWD_FILE` set, read it back, `cd` if
it changed.)

The same `.shell` output also defines that shell's completion for qsoku's
defined names (see [Shell completion](#shell-completion)) — one `eval`
covers both.

PowerShell is not a target: `qsokufile` commands always run under `sh` (see
[qsokufile.md](qsokufile.md)), which on Windows means Git Bash or WSL, both
of which already have bash.

## Shell completion

Pressing Tab after `qsoku ` completes on the names defined in the
`qsokufile` in use. This is **dynamic**: the completion script defined by
`.shell` calls `qsoku .names` on every Tab press rather than embedding a
fixed list, so completion always matches the `qsokufile` actually in use in
the current directory. `.names` is built for exactly this (see
[Management commands](#management-commands)): it never errors and never
writes to standard error, so an outdated, missing, or currently-broken
`qsokufile` never interrupts typing.

Management commands (`.init`, `.add`, and so on) are also completed, since
they are valid names to type after `qsoku`. They are the one part of
completion that **is** a fixed list, written directly into each shell's
script rather than fetched from `qsoku`: unlike the defined names, that list
only changes with a new qsoku release, and the script is embedded in, and
shipped with, that same release, so there is no version to fall out of sync
with.

- **fish**: candidates starting with `.` are hidden until the word being
  typed itself starts with `.` — the same rule fish applies to dotfiles in
  path completion, applied here too since it goes by the leading character,
  not by whether a candidate is a file. So `qsoku <TAB>` on fish offers only
  the defined names; `qsoku .<TAB>` is needed to see the management
  commands. bash and zsh have no such rule and offer both at once.

## Exit codes

| Code | Meaning |
|---|---|
| 0 | Success |
| 1 | qsoku could not do what was asked: no `qsokufile` found, the `qsokufile` fails to parse (bad line, duplicate name), `sh` could not be started, `$EDITOR`/`nano` not found for `.edit`, and so on |
| 2 | The command line is wrong: the given name is not defined in the `qsokufile` (management commands included — an unknown `.foo`), or a management command got the wrong number of arguments |
| *(other)* | Once `sh` starts running the entry's command, qsoku's own exit code stops applying: the command's exit code (0-255) is returned unchanged, including 128+*n* if it was killed by signal *n*. qsoku never substitutes its own value once the command has actually run |

This mirrors mtqg's own convention (0 / 1 / 2) rather than inventing a
separate scheme (see [CLAUDE.md](../../CLAUDE.md) "独自の決まりを増やさない").
Unlike `make`, which returns a fixed `2` for a failed recipe, qsoku is a thin
wrapper: `qsoku test` in CI sees exactly the same exit code `go test` itself
would have produced.
