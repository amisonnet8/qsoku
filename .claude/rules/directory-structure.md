# ディレクトリ構成

**◯は今あるもの、△はまだ無いもの**（実装しながら作る）。名前だけでは中身が分からないものに説明を添える。

```
qsoku/
◯ CLAUDE.md               プロジェクトルール（参照先の案内、mtqgでの記録の使い分け）
◯ LICENSE                 MIT
◯ .gitignore
◯ .gitattributes          `* text=auto eol=lf`（Windowsでの改行コード変換による誤検知を防ぐ）
◯ .golangci.yaml          lintの設定
◯ trivy.yaml               依存の脆弱性・ライセンス検査の設定
◯ .goreleaser.yaml         バイナリ配布の設定（distribution.md）。動くのは人間が
◯                          v*タグをpushしたとき（.github/workflows/release.yml）だけ。
◯                          設定が壊れていないかは.github/workflows/ci.yml の
◯                          goreleaserジョブ（make goreleaser-check）が毎pushで確認
◯ README.md / README_ja.md 看板（ロゴ・バッジ・目次・特徴・デモ・実行例・インストール・
                            tour/reference/examplesへのリンク）。冒頭に開発初期の注意書きを
                            残す（Step 9で作成、v0.1.0公開後に「映え対応」で作り直し）
◯ qsokufile                このリポジトリ自身の近道（`build`・`test`、`//`でルートへ戻って
                            `go build`・`go test`を呼ぶ）。看板READMEのデモGIF録画用に新設
                            したが、実際に使える本物の近道として残す（ドッグフーディング）
◯ Makefile                ビルド・テストの入口（`make build`・`make check`など）
◯ go.mod                   `go.sum`はまだ無い（依存が無いため）
◯ docs/
◯ ├── assets/               看板READMEが使う画像。`logo.svg`（手書きのSVG。
◯ │                          `@media (prefers-color-scheme: dark)`を埋め込み、ラスター画像は
◯ │                          置かない）と`demo.gif`（`vhs`で録画した本物のターミナル操作。
◯ │                          このリポジトリ自身の`qsokufile`を`internal/cli/`から実行し、
◯ │                          `//`でルートに戻って`go test ./...`が走る様子。録画の手順は
◯ │                          `docs/design/history.md`参照、リポジトリには残さない一回限りの
◯ │                          セットアップ）
◯ ├── design/              設計判断と理由の記録（日本語）。README.md
◯ │   ├── README.md         このディレクトリの位置づけと索引
◯ │   ├── note.md           最初の設計メモ（原文のまま。qsokuとは何か、`//`の規則など）
◯ │   └── history.md        決めたことの時系列
◯ ├── reference/            仕様。英語版（正）と日本語版`*_ja.md`の2本立て。実行例は
◯ │                          `e2e/examples_test.go`が実測して書き込む
◯ │   └── README.md         置き場所の決まりと、実行例の印の書き方
◯ ├── tour/                 歩いて回る入門（README.md/README_ja.md）。`.init`から
◯ │                          シェル連携・補完まで、実行例つきで手を動かして追う（Step 9）
◯ └── examples/             実例（README.md/README_ja.md＋go/node/monorepoの
◯                            qsokufile）。コピーして使える。パースだけ`e2e/run_test.go`の
◯                            `TestDocsExamplesQsokufilesParse`が確かめる（実行はしない）
◯ cmd/qsoku/main.go        エントリポイント。引数を渡すだけ
◯ internal/cli/             実装本体。`.`で始まらない名前は実際に見つけて実行する
◯                            （`internal/qsokufile`・`internal/run`を呼ぶ）。管理用コマンド
◯                            （`.version`・`.init`・`.add`・`.rm`・`.list`・`.names`・
◯                            `.edit`・`.where`・`.shell`・`.help`／引数なし）を1コマンド
◯                            1ファイルで実装（mtqgの`internal/cli/`に倣う）。未知の`.foo`は
◯                            終了コード2
◯ internal/cli/shells/        `qsoku.bash`・`qsoku.zsh`・`qsoku.fish`（`go:embed`。
◯                            mtqgの`internal/cli/completions/`に相当）。居場所を持ち帰る
◯                            `qsoku`関数と、そのシェルの補完を1本にまとめて持つ
◯ internal/qsokufile/        qsokufileの探索（`Find`）と解析（`Parse`）、両方をまとめた
◯                            `Load`、名前引き（`Lookup`）、`//`の置き換え（`Substitute`・
◯                            `SubstituteArg`）、書き込み（`SetEntry`・`RemoveEntry`。
◯                            対象の行だけ最小限に書き換え、コメント・並び順は崩さない）。
◯                            `internal/cli`・`internal/run`を知らない層
◯ internal/run/               `sh`を実際に起動するだけの層（`Execute`）。qsokufileが
◯                            何かは知らない。居場所の持ち帰り・終了コードの素通しを担う
◯                            （`internal/qsokufile`・`internal/cli`のどちらも知らない。
◯                            `.golangci.yaml`のdepguardで3層の依存の向きを強制）
◯ e2e/                      本物のバイナリと本物のシェル（bash・zsh・fish）で動かすテスト
◯                            （`e2e_test.go`が`TestMain`でバイナリを1回ビルド）。土台は
◯                            予定より前倒しでStep 7に作った（`.mtqg`のhistory参照）。
◯ ├── shell_test.go          シェル連携・補完（Step 7）
◯ ├── run_test.go            本物のバイナリを直接実行（終了コードの素通し・シグナル・
◯                            `//`置き換え・`QSOKU_CWD_FILE`の受け渡し。Step 8）
◯ ├── examples_test.go       docs/reference/・README・docs/tour/の実行例を実測で確かめる
◯                            （`make docs-examples`で出力を文書へ書き込む。対象文書は
◯                            `documentPairs`。mtqgの`e2e/examples_test.go`の簡素化版。
◯                            Step 8で作り、Step 9でREADME・tour/にも対象を広げた）
◯ └── testdata/examples/     ↑の例が使うfixture（言語非依存。英日どちらの文書からも参照）
◯ .devcontainer/           devcontainer.json・postCreate.sh
◯ .github/workflows/
◯ ├── ci.yml                 CI（Linux・macOSのマトリクス。check・race・shellcheck・
◯ │                            trivy・goreleaser（make goreleaser-checkの安全網）の5系統）
◯ └── release.yml            v*タグのpushだけで動く。goreleaser-actionで実際に
◯                            ビルド・GitHub Releaseの公開まで行う
◯ .claude/
◯ ├── settings.json         権限（deny/ask）とビルドフックの設定。人間が管理する
◯ ├── rules/                 このファイルを含む、育てていくルール
◯ └── hooks/                 build.sh：`.go`・`go.mod`・`go.sum`編集後に`make build`する
◯                            （mtqgの`.claude/hooks/build.sh`と同じ）
```

## `docs/`の4つの違い（迷いやすいので明記）

| ディレクトリ | 何のためのものか | 読者 |
|---|---|---|
| `design/` | **なぜこうなっているか**（判断の理由、経緯）。仕様と食い違えば`reference/`が正しい | 開発するAI・人間 |
| `reference/` | **今、何をするか**（仕様そのもの）。正式な契約 | 利用者・実装者 |
| `tour/` | qsokuを初めて触る人が、手を動かしながら一通り追える入門 | 利用者（読み物） |
| `examples/` | 動く実例の置き場（コピーして使える`qsokufile`など） | 利用者（コピー元） |

`tour/`・`examples/`は、mtqg本体の`docs/tour/`・`docs/examples/`と同じ役割分担を踏襲している（実装完了後に作る、英日2本立てで作る、という決まりも同じ。Step 9で作った）。

## 配置の判断基準

- **`internal/`**：Goの仕組みとして、リポジトリの外からimportできない。外との約束は、配布物（`go install`で入るバイナリ）とデータ形式（`qsokufile`の書式）だけ
- **配布物（ビルド済みバイナリ）はコミットしない**（`.gitignore`にルート直下限定で`/qsoku`・`/qsoku.exe`・`/dist/`）
- **シェル連携・補完のスクリプト**は`internal/cli/shells/`に埋め込み（`go:embed`）で配る（mtqgの`internal/cli/completions/`と同じやり方）。対応シェルはbash・zsh・fish（PowerShellは対象外）
- **看板としてのREADME・`tour/`・`examples/`は実装完了後（Step 9）に作った。** 未完成の間の注意書きは看板の冒頭に残す
