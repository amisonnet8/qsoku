# Examples

*[日本語](README_ja.md) | **English***

Working `qsokufile`s you can copy into your own repository and adjust. See
[docs/tour/](../tour/) for a walk-through of qsoku itself, and
[docs/reference/](../reference/) for the formal specification.

| Example | Shows |
|---|---|
| [go/qsokufile](go/qsokufile) | The usual Go shortcuts (`build`, `test`, `race`, `lint`, `run`), and `//` to always run `go mod tidy` at the repository root regardless of where you are |
| [node/qsokufile](node/qsokufile) | The usual npm shortcuts, and `//` to `rm -rf` build output from the repository root safely |
| [monorepo/qsokufile](monorepo/qsokufile) | A multi-package repository: an unparenthesized entry that deliberately moves you into a package directory, alongside parenthesized entries that build/test each package without moving you, and one entry calling another with `qsoku` |
