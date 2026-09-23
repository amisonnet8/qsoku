# qsoku

*[日本語](README_ja.md) | **English***

> **Work in progress. Not ready to use yet.**
> qsoku is under active development. There is no released version yet
> (`@latest` installs the newest commit, not a tagged release), and the data
> format and commands may change without notice. Please do not use it in
> your projects for now.

**qsoku** writes a repository's own command shortcuts into a `qsokufile`, so
anyone (or anything) in that repository can run `qsoku <name>` from
anywhere inside it.

- Unlike a shell `alias`, a `qsokufile` is committed with the repository, so
  the whole team gets the same shortcuts, and they never leak into your
  other projects.
- Unlike `make`, qsoku has no dependency graph, no `.PHONY`, and no
  tab-sensitive recipes — just a flat list of names, and `cd` works fine
  inside a shortcut.
- An AI agent reading the repository sees the same `qsokufile` a human
  does — one list of names either can run.
- Only three rules are qsoku's own: how it finds the `qsokufile`, what `//`
  means, and how it brings the working directory back. Everything else is
  ordinary `sh`.

## Example

<!-- qsoku:example dir=basic -->
```
$ qsoku hello world
hello, world
```

```
# qsokufile
hello: echo "hello, $1"
```

## Install

```sh
go install github.com/amisonnet8/qsoku/cmd/qsoku@latest
```

Then add one line to your shell's startup file (see
[cli.md](docs/reference/cli.md#shell-integration) for bash/zsh/fish
specifics):

```sh
eval "$(qsoku .shell bash)"
```

This also wires up shell completion for the names in your `qsokufile`.

## Learn more

- [docs/tour/](docs/tour/) — a hands-on walk through qsoku, from an empty
  `qsokufile` to shell integration and completion
- [docs/reference/](docs/reference/) — the formal specification of the
  `qsokufile` format and the `qsoku` command line
- [docs/examples/](docs/examples/) — `qsokufile`s you can copy into your own
  repository

## License

MIT (see [LICENSE](LICENSE)).
