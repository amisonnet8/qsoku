# qsoku specification

*[日本語](README_ja.md) | **English***

The formal specification of **what qsoku does** goes here: the
`qsokufile` format, the `//` substitution rule, the command line, and so
on.

- **This is the specification of record.** `docs/design/` records *why*
  things are this way; where it disagrees with this directory, this
  directory is right.
- **The English document (`*.md`) is the source of truth, paired with a
  Japanese one (`*_ja.md`).** Fix both in the same change.
- Put the cross-language link right below the heading (the current
  language in bold): `*[日本語](xxx_ja.md) | **English***` in the English
  document, `*[English](xxx.md) | **日本語***` in the Japanese one.
- **Executable examples are never written by hand.** Put a mark
  `<!-- qsoku:example dir=<fixture> [cwd=<subdir>] [skip="why"] -->`
  right before a code block that starts with `$ qsoku ...`, and
  `e2e/examples_test.go` (`make docs-examples`) runs it against the
  actually-built qsoku binary and writes its output into that block. A
  fixture is a directory under `e2e/testdata/examples/` (`none` is an
  empty directory with no `qsokufile`). Every `$ ` line must start with
  `qsoku` — these examples show qsoku's own output, not an arbitrary
  shell session. `make test` (`TestDocExamples`) catches a document that
  has drifted from what qsoku actually prints. This mechanism is not
  limited to `docs/reference/`: the top-level `README.md` and
  `docs/tour/` use it too (see `e2e/examples_test.go`'s `documentPairs`
  for the full list of documents it covers).
- **This is a living document, grown alongside the implementation.**
  Update it when the implementation and the spec drift apart; when
  changing a design decision, update this first, before touching code.

## Files

| File | Contents |
|---|---|
| [qsokufile.md](qsokufile.md) / [qsokufile_ja.md](qsokufile_ja.md) | The `qsokufile` format: location and lookup, the `name: command` line format, comments, the characters allowed in a name, the `//` substitution rule |
| [cli.md](cli.md) / [cli_ja.md](cli_ja.md) | The `qsoku <name>` command line: how running a name works and its environment variables, the management commands (`.init` through `.help`), shell integration, shell completion, exit codes |

See also [docs/tour/](../tour/) for a hands-on walk-through, and
[docs/examples/](../examples/) for `qsokufile`s you can copy.
