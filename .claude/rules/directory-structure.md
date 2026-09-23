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
◯ README.md / README_ja.md 未完成の間の注意書き（「開発中でまだ使えません」だけ）
                            看板としてのREADME（売り文句・使い方）は実装完了後に作る
◯ Makefile                ビルド・テストの入口（`make build`・`make check`など）
◯ go.mod                   `go.sum`はまだ無い（依存が無いため）
◯ docs/
◯ ├── design/              設計判断と理由の記録（日本語）。README.md
◯ │   ├── README.md         このディレクトリの位置づけと索引
◯ │   ├── note.md           最初の設計メモ（原文のまま。qsokuとは何か、`//`の規則など）
◯ │   └── history.md        決めたことの時系列
◯ ├── reference/            仕様。英語版（正）と日本語版`*_ja.md`の2本立て
◯ │   └── README.md         今は置き場所の決まりだけ。実装の前にここへ書く
△ ├── tour/                 歩いて回る入門。手を動かしながらqsokuの一通りの使い方を
△ │                          追える読み物（README.mdの「使い方」の詳しい版）
△ └── examples/             実例。実際のqsokufileと、そのまま動く使い方の例
△                            tour/・examples/は**実装完了後**に作る。英語版＋`*_ja.md`
◯ cmd/qsoku/main.go        エントリポイント。引数を渡すだけ
◯ internal/cli/             実装本体。今は`Run`（`.version`だけ実装、ほかは仮に
◯                            `not implemented yet`で終了コード1）と`buildVersion`のみ
◯ internal/qsokufile/        qsokufileの探索（`Find`）と解析（`Parse`）、両方をまとめた
◯                            `Load`、`//`の置き換え（`Substitute`・`SubstituteArg`）。
◯                            `internal/cli`を知らない層（`.golangci.yaml`のdepguardで
◯                            強制）。層への分け方は引き続きStep 5〜6で足していく
△ e2e/                      本物のバイナリと本物のシェル（bash・zsh・fish）で動かすテスト
◯ .devcontainer/           devcontainer.json・postCreate.sh
◯ .github/workflows/        CI（Linux・macOSのマトリクス）
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

`tour/`・`examples/`は、mtqg本体の`docs/tour/`・`docs/examples/`と同じ役割分担を踏襲している（実装完了後に作る、英日2本立てで作る、という決まりも同じ）。

## 配置の判断基準

- **`internal/`**：Goの仕組みとして、リポジトリの外からimportできない。外との約束は、配布物（`go install`で入るバイナリ）とデータ形式（`qsokufile`の書式）だけ
- **配布物（ビルド済みバイナリ）はコミットしない**（`.gitignore`にルート直下限定で`/qsoku`・`/qsoku.exe`・`/dist/`）
- **シェル連携・補完のスクリプト**は、埋め込み（`go:embed`）で配る想定（mtqgの`internal/cli/completions/`と同じやり方が参考になる）。対応シェルはbash・zsh・fish（PowerShellは対象外）
- **看板としてのREADMEは最後に作る**（今あるのは注意書きだけの版）。`tour/`・`examples/`も同じタイミング（実装完了後）
