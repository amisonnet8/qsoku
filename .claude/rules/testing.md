# テスト方針

## 開発環境の前提

開発者の手元環境は**Linux（devcontainer）のみ**。Windowsでの動作確認は、qsokuが`sh`実行を前提とする以上、優先度は低い（Git BashやWSLでの動作は、要望が出てから考える）。

**CI（`.github/workflows/`）の対象OSはLinux・macOS**（Windowsは対象外。人間の判断、2026-09-23、`docs/reference/`仕様執筆時に確認）。

| 対象 | 進め方 |
| :--- | :--- |
| `qsokufile`の解析・`//`の置き換え・管理用コマンド | Linux・Goだけで単体テストを書く。字句解析（引用符の状態を追う部分）は、URL・引用符の中・括弧の中・`;`の後を表で押さえる |
| シェル連携・補完（bash・zsh・fish） | 3つとも devcontainer に入れてある（`postCreate.sh`）。**本物のシェルで動かして確かめる**（スクリプトを模したものに置き換えない）。`e2e/shell_test.go`（`go test -tags e2e`、`make test`）が実施 |
| 本物のバイナリを直接実行（終了コード・シグナル・`//`置き換え・`QSOKU_CWD_FILE`） | `e2e/run_test.go`（`make test`）。`internal/cli`の単体テストは`Run`をプロセス内で呼ぶだけなので、ビルドした実バイナリでしか見えない部分をここで押さえる |
| `docs/reference/`の実行例 | 手で書かない。`$ qsoku ...`の直前に`<!-- qsoku:example dir=<fixture> ... -->`を置くと、`e2e/examples_test.go`が実際にビルドしたqsokuで実行し、`make docs-examples`が出力を文書へ書き込む（`docs/reference/README.md`参照）。比較だけなら`make test`が実施 |
| 依存の脆弱性・ライセンス | `trivy`（`trivy.yaml`） |
| シェルスクリプト | `shellcheck`（追跡中の`*.sh`・`*.bash`のみ） |

## 実装後の動作確認

- ビルド確認に加えて、実際に`qsoku`コマンドを打って確かめること
- `qsokufile`のコマンドを実行するテストは、**本物の`sh`**で動かす（sh実行を模したものに置き換えない）

## mtqg固有の検証項目（qsokuに置き換えたもの）

- **常に`sh`で実行すること**：`()`がサブシェルとして扱われ、外の居場所を変えないこと。`cd`を含むコマンドが、成功でも失敗でも実行後の居場所を正しく持ち帰ること
- **`//`の置き換え**：単語の先頭で引用符の外にあるときだけ置き換わり、URL（`http://`）・引用符の中（シングル・ダブル）では置き換わらないこと
- **居場所の受け渡し**：`QSOKU_CWD_FILE`に書き出す1行が、コマンド自身の標準出力（`make build`のログなど）と混ざらないこと。qsoku自身のエラー（名前がない、`sh`が起動できない）では一時ファイルに何も書かれず、シェル関数側で`cd`が起きないこと
- **終了コード**：コマンドを実行した後は`sh`の終了コードをそのまま返すこと（qsoku自身の決まりで上書きしない）。コマンドを実行する前のqsoku自身のエラーは0/1/2（`docs/reference/cli.md`「Exit codes」）を表テストで押さえる
- **シェル補完**：登録されている名前を補完する。bash・zsh・fishで確認する（mtqgの`e2e/completion_test.go`のやり方が参考になる）。**fishは同じ候補にならない**：`.`で始まる候補（管理用コマンド）は、打っている語自体が`.`で始まるまでfishが隠すため（ドットファイルの扱いと同じ規則）。`qsoku <TAB>`は定義済みの名前だけ、`qsoku .<TAB>`で管理用コマンドが出ることを別々に確かめる（`docs/reference/cli.md`「Shell completion」に明記）

## Trivy・ShellCheckの運用

- `go get`で依存を足したら、必ずTrivyを通す。依存は原則Go標準ライブラリと`golang.org/x/`だけ（`CLAUDE.md`）
- シェルスクリプトのコメント行を、行頭が小文字の`# shellcheck`で始まる文にしない（ShellCheckがインラインディレクティブとして誤解釈し、SC1072/SC1073でパースが止まる）

## `-race`の運用

- devcontainerは`CGO_ENABLED=0`。`-race`が要るときはcgoを有効にして走らせる（`Makefile`ができたら`make race`として分ける）
