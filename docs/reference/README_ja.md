# qsoku 仕様

*[English](README.md) | **日本語***

`qsokufile`の書式・`//`の置き換え規則・コマンドラインの仕様など、qsokuが**何をするか**の正式な仕様をここに置く。

- **ここが仕様の正。** `docs/design/`は「なぜこうなっているか」の記録で、仕様と食い違う場合はここが優先する
- **英語版（`*.md`）が正、日本語版（`*_ja.md`）と2本立て。** 同じ変更の中で両方直すこと
- 見出しのすぐ下に相互リンクを置く（現在の言語を太字に）。英語版は`*[日本語](xxx_ja.md) | **English***`、日本語版は`*[English](xxx.md) | **日本語***`の形
- **実行例は手で書かない。** `$ qsoku ...`で始まるコードブロックの直前に
  `<!-- qsoku:example dir=<fixture> [cwd=<subdir>] [skip="理由"] -->`という印を置くと、
  `e2e/examples_test.go`（`make docs-examples`）が実際にビルドしたqsokuで実行し、
  その出力をブロックへ書き込む。fixtureは`e2e/testdata/examples/`（`none`は
  qsokufileの無い空ディレクトリ）。各`$ `行は`qsoku`で始まるものに限る（任意の
  シェルセッションではなく、qsoku自身の出力を見せるため）。`make test`が
  `TestDocExamples`で文書と実際の出力の食い違いを検知する。この仕組みは
  `docs/reference/`だけでなくトップの`README.md`・`docs/tour/`の実行例にも使う
  （`e2e/examples_test.go`の`documentPairs`が対象文書の一覧）
- **実装しながら育てる文書。** 実装と仕様がずれたらここを更新する。設計判断を変えるときは、まずここを更新してから着手する

## ファイル

| ファイル | 中身 |
|---|---|
| [qsokufile.md](qsokufile.md) / [qsokufile_ja.md](qsokufile_ja.md) | `qsokufile`の書式：置き場所と探し方、`名前: コマンド`の1行形式、コメント、名前に使える文字、`//`の置き換え規則 |
| [cli.md](cli.md) / [cli_ja.md](cli_ja.md) | `qsoku <名前>`コマンドライン：実行のしくみと環境変数、管理用コマンド（`.init`〜`.help`）、シェル連携、シェル補完、終了コード |

手を動かす入門は[docs/tour/](../tour/)、コピーして使える`qsokufile`は
[docs/examples/](../examples/)を参照。
