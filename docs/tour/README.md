# Tour

*[日本語](README_ja.md) | **English***

A hands-on walk through qsoku, from an empty repository to shell
integration and completion. For the formal specification, see
[docs/reference/](../reference/); for ready-to-copy examples, see
[docs/examples/](../examples/).

## Start from nothing

A repository with no `qsokufile` yet:

<!-- qsoku:example dir=none -->
```
$ qsoku .init
Created /home/you/project/qsokufile
```

`qsoku .init` creates an empty `qsokufile` in the current directory. It is
just a text file — open it in any editor, or use `qsoku .add` (next).

## Add a shortcut

<!-- qsoku:example dir=none -->
```
$ qsoku .init
Created /home/you/project/qsokufile
$ qsoku .add build 'echo "building the project"'
$ qsoku .add hello 'echo "hello, $1"'
$ qsoku .list
build: echo "building the project"
hello: echo "hello, $1"
```

Each line of a `qsokufile` is `name: command`. `qsoku .add` appends one
(or replaces it in place, if the name is already defined) without touching
anything else you've written by hand — comments and the order of the other
lines are left alone.

## Run it

<!-- qsoku:example dir=tour -->
```
$ qsoku build
building the project
$ qsoku run everyone
running with arg: everyone
```

Every argument after the name is passed straight through to the command, as
`$1`, `$2`, and so on — exactly as `sh` itself would expand them. The
command always runs with `sh`, no matter which shell you typed `qsoku`
from, so a `qsokufile` behaves the same for the whole team.

## `//` is the qsokufile's own location

<!-- qsoku:example dir=tour -->
```
$ qsoku test
testing from /home/you/project
```

The `test` entry is `(cd //; echo "testing from $QSOKU_ROOT")`: `//` expands
to the directory holding the `qsokufile` in use, so a shortcut can always
find its way back to the repository root, no matter where you ran it from.
The parentheses keep that `cd` inside a subshell — see below.

## The same `qsokufile` from any subdirectory

<!-- qsoku:example dir=tour cwd=src -->
```
$ qsoku build
building the project
```

qsoku walks up from the current directory until it finds a `qsokufile` —
the same one applies everywhere underneath it, the way `.git/` does for
git.

## Bringing the working directory back

A `qsokufile` entry runs in a separate process, so a plain `cd` inside it
would not normally move your own shell. The shell integration below (`qsoku
.shell`) fixes this: after running a name, the wrapper function reads back
where qsoku finished and `cd`s there itself, on success or failure alike.
Wrap a `cd` in parentheses, as `test` does above, if you want it to stay
inside the shortcut instead.

See [cli.md](../reference/cli.md#bringing-the-working-directory-back) for
the exact rules (what happens on a failed `cd`, symlinks, and so on).

## Shell integration and completion

Add one line to your shell's startup file:

```sh
eval "$(qsoku .shell bash)"   # ~/.bashrc (zsh: ~/.zshrc)
```

```fish
qsoku .shell fish | source    # ~/.config/fish/config.fish
```

This defines a `qsoku` shell function that runs the real binary, brings the
working directory back, and also turns on Tab completion for every name in
the `qsokufile` in use (plus qsoku's own management commands). See
[cli.md](../reference/cli.md#shell-integration) for what the function
actually does, and
[cli.md](../reference/cli.md#shell-completion) for how completion behaves
in each shell (fish, in particular, hides the management commands until you
type a leading `.`).

## Fixing a mistake, or finding the file

<!-- qsoku:example dir=tour -->
```
$ qsoku .where
/home/you/project
$ qsoku .rm run
$ qsoku .list
build: echo "building the project"
test: (cd //; echo "testing from $QSOKU_ROOT")
```

`qsoku .where` prints the directory holding the `qsokufile` in use (what
`//` expands to); `qsoku .edit` opens it with `$EDITOR`. Both work even when
the file currently fails to parse — you need them to fix it.

## Commit it

A `qsokufile` is meant to be committed, just like a `Makefile`. Once it's in
the repository, every teammate — and every AI agent working in that
repository — sees the same list of shortcuts and can run them the same way.

## Where to go next

- [docs/reference/](../reference/) for the full specification (name
  characters, comments, exit codes, and so on)
- [docs/examples/](../examples/) for `qsokufile`s you can copy into your own
  project
