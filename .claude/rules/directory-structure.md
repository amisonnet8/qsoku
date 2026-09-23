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
◯                            （`e2e_test.go`が`TestMain`でバイナリを1回ビルド）。予定より
◯                            前倒しでStep 7に作った（`.mtqg`のhistory参照）。docsの例の
◯                            確認の仕組みはまだ無い（Step 8）
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
- **シェル連携・補完のスクリプト**は`internal/cli/shells/`に埋め込み（`go:embed`）で配る（mtqgの`internal/cli/completions/`と同じやり方）。対応シェルはbash・zsh・fish（PowerShellは対象外）
- **看板としてのREADMEは最後に作る**（今あるのは注意書きだけの版）。`tour/`・`examples/`も同じタイミング（実装完了後）
