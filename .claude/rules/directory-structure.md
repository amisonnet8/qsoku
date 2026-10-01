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
                            reference/examplesへのリンク）。v0.1.0公開後に
                            「映え対応」で作り直した。冒頭の「開発初期」の注意書きはv0.1.0公開
                            後に外した（人間の判断）
◯ qsokufile                Makefileの主なターゲットを写した見本（`build`・`fmt`・`vet`・
                            `lint`・`test`・`check`・`e2e`・`docs-examples`・`race`・`trivy`・
                            `shellcheck`・`goreleaser-check`）。**このリポジトリ自身の開発には
                            使わない**（実開発は引き続き`make`。CLAUDE.md・testing.md）。
                            `test`は看板READMEのデモGIFが依存する`(cd //; go test ./...)`の
                            まま（Makefileの`unit`相当。名前がMakefile側のe2eターゲットと
                            ずれるため、そちらは`e2e`という名前にした）
◯ Makefile                ビルド・テストの入口（`make build`・`make check`など）
◯ go.mod                   `go.sum`はまだ無い（依存が無いため）
◯ .mtqg/                   mtqgの記録（journal.jsonl・SCHEMA.md・version・
                            .gitattributes・.gitignore）。`.local/`はコミットされない
                            （マシンごとのロック・一時ファイル。`.claude/rules/mtqg.md`）
◯ docs/
◯ ├── assets/               看板READMEが使う画像。`logo.svg`（手書きのSVG。
◯ │                          `@media (prefers-color-scheme: dark)`を埋め込み、ラスター画像は
◯ │                          置かない）と`demo.gif`（`vhs`で録画した本物のターミナル操作。
◯ │                          このリポジトリ自身の`qsokufile`を`internal/cli/`から実行し、
◯ │                          `//`でルートに戻って`go test ./...`が走る様子。リポジトリには
◯ │                          残さない一回限りのセットアップ）
◯ ├── design/              設計判断と理由の記録（日本語）。README.md
◯ │   ├── README.md         このディレクトリの位置づけと索引
◯ │   ├── note.md           最初の設計メモ（原文のまま。qsokuとは何か、`//`の規則など）
◯ │   └── history.md        決めたことの時系列
◯ ├── reference/            仕様。英語版（正）と日本語版`*_ja.md`の2本立て。実行例は
◯ │                          `e2e/examples_test.go`が実測して書き込む
◯ │   └── README.md         置き場所の決まりと、実行例の印の書き方
◯ └── examples/             歩いて回る入門（README.md/README_ja.mdの前半。`.init`から
◯                            シェル連携・補完まで、実行例つきで手を動かして追う）と、
◯                            コピーして使える実例（同後半＋go/node/monorepoの
◯                            qsokufile）。旧`docs/tour/`を2026-09-28に統合した
◯                            （小さすぎたため。`docs/design/history.md`参照）。実例の
◯                            qsokufileのパースだけ`e2e/run_test.go`の
◯                            `TestDocsExamplesQsokufilesParse`が確かめる（実行はしない）
◯ cmd/qsoku/main.go        エントリポイント。引数を渡すだけ
◯ internal/cli/             実装本体。`.`で始まらない名前は実際に見つけて実行する
◯                            （`internal/qsokufile`・`internal/run`を呼ぶ）。管理用コマンド
◯                            （`.version`・`.init`・`.add`・`.rm`・`.list`・`.names`・
◯                            `.edit`・`.where`・`.shell`・`.help`／引数なし）を1コマンド
◯                            1ファイルで実装（mtqgの`internal/cli/`に倣う）。未知の`.foo`は
◯                            終了コード2
◯ internal/cli/shells/        `qsoku.bash`・`qsoku.zsh`・`qsoku.fish`・`qsoku.ps1`
◯                            （`go:embed`。mtqgの`internal/cli/completions/`に相当）。
◯                            居場所を持ち帰る`qsoku`関数と、そのシェルの補完を1本に
◯                            まとめて持つ。`qsoku.ps1`はpwsh（PowerShell 7+）向けで、
◯                            対応OSはWindows・Linux・macOS共通（2026-09-27）
◯ internal/qsokufile/        qsokufileの探索（`Find`）と解析（`Parse`）、両方をまとめた
◯                            `Load`、名前引き（`Lookup`）、`//`の置き換え（`Substitute`・
◯                            `SubstituteArg`）、書き込み（`SetEntry`・`RemoveEntry`。
◯                            対象の行だけ最小限に書き換え、コメント・並び順は崩さない）。
◯                            `internal/cli`・`internal/run`を知らない層
◯ internal/run/               `sh`を実際に起動するだけの層（`Execute`）。qsokufileが
◯                            何かは知らない。居場所の持ち帰り・終了コードの素通しを担う
◯                            （`internal/qsokufile`・`internal/cli`のどちらも知らない。
◯                            `.golangci.yaml`のdepguardで3層の依存の向きを強制）
◯ e2e/                      本物のバイナリと本物のシェル（bash・zsh・fish・pwsh）で動かす
◯                            テスト（`e2e_test.go`が`TestMain`でバイナリを1回ビルド。
◯                            Windowsでは`qsoku.exe`）
◯ ├── shell_test.go          シェル連携・補完
◯ ├── run_test.go            本物のバイナリを直接実行（終了コードの素通し・シグナル・
◯                            `//`置き換え・`QSOKU_CWD_FILE`の受け渡し）。
◯                            `docs/examples/*/qsokufile`とリポジトリ直下の`qsokufile`が
◯                            パースできることも、ここで（`qsoku .list`を走らせるだけで）
◯                            確かめる（コマンドの中身は実行しない）
◯ ├── examples_test.go       docs/reference/・README・docs/examples/の実行例を実測で確かめる
◯                            （`make docs-examples`で出力を文書へ書き込む。対象文書は
◯                            `documentPairs`。mtqgの`e2e/examples_test.go`の簡素化版）
◯ └── testdata/examples/     ↑の例が使うfixture（言語非依存。英日どちらの文書からも参照）
◯ .devcontainer/           devcontainer.json・postCreate.sh。mtqgのインストール
                            （v1タグ公開後は`@latest`。2026-09-29）とbash補完（`mtqg
                            completion bash`）、qsoku自身のインストール（ローカルソース
                            から`go install ./cmd/qsoku`。mtqg自身のdevcontainerと同じ
                            やり方。2026-09-29）とシェル連携（`~/.bashrc`に`qsoku .shell
                            bash`。リポジトリ直下のqsokufileはあくまで見本で実開発はmake
                            のまま——動作確認・dogfooding目的。2026-09-29）もここ。
                            gh・pwshはdevcontainer.jsonの公式feature（`github-cli`・
                            `powershell`、2026-10-01にaptの手動インストールから置換）
                            `mkdir -p ~/.cache/trivy`（2026-09-30）は、Bashサンドボックスの
                            書き込み許可が既存パスへのbindマウント式で、無いパスを新しく
                            掘れないため——サンドボックス化後の`make trivy`初回失敗の対策
◯ .mcp.json                mtqgのMCPサーバー（`mtqg mcp`）の起動設定。`mtqg init --agent
                            claude-code`が作った（2026-09-25）
◯ .github/workflows/
◯ ├── ci.yml                 CI（check・race・shellcheck・trivy・goreleaser
◯ │                            （make goreleaser-checkの安全網）の5系統。
◯ │                            checkはLinux・macOS・Windowsの3OS、raceは
◯ │                            Linux・macOS、shellcheck・trivy・goreleaserは
◯ │                            Linuxのみ（Windowsは2026-09-27に対象に加えた）
◯ └── release.yml            v*タグのpushだけで動く。goreleaser-actionで実際に
◯                            ビルド・GitHub Releaseの公開まで行う
◯ .claude/
◯ ├── settings.json         権限（deny/ask）・サンドボックス（ネットワーク許可ドメイン・
◯                            書き込み許可パス。2026-09-30、mtqg本体のsettings.jsonを参考に
◯                            変更。CLAUDE.md「権限・自動化について」参照）・フックの設定。
◯                            人間が管理する部分（permissions・sandbox・build.shの
◯                            PostToolUse）に加え、`mtqg init --agent claude-code`
◯                            （2026-09-25）が`env`（記録者の自動設定）・`SessionStart`／
◯                            `Stop`フック（`mtqg hook claude-code`。`mtqg context`の自動
◯                            読み込みと、記録漏れがあれば終了時に促す）を足した
◯ ├── rules/                 このファイルを含む、育てていくルール
◯ └── hooks/                 build.sh：`.go`・`go.mod`・`go.sum`編集後に`make build`する
◯                            （mtqgの`.claude/hooks/build.sh`と同じ）
```

## `docs/`の4つの違い（迷いやすいので明記）

| ディレクトリ | 何のためのものか | 読者 |
|---|---|---|
| `design/` | **なぜこうなっているか**（判断の理由、経緯）。仕様と食い違えば`reference/`が正しい | 開発するAI・人間 |
| `reference/` | **今、何をするか**（仕様そのもの）。正式な契約 | 利用者・実装者 |
| `examples/` | qsokuを初めて触る人が手を動かしながら一通り追える入門と、動く実例の置き場（コピーして使える`qsokufile`など） | 利用者（読み物・コピー元） |

`examples/`は、mtqg本体の`docs/tour/`・`docs/examples/`の役割を1つにまとめたもの（qsoku側は入門を`tour/`として独立させていたが、小さすぎたため2026-09-28に`examples/`へ統合した。英日2本立てで作る、という決まりは変えていない）。

## 配置の判断基準

- **`internal/`**：Goの仕組みとして、リポジトリの外からimportできない。外との約束は、配布物（`go install`で入るバイナリ）とデータ形式（`qsokufile`の書式）だけ
- **配布物（ビルド済みバイナリ）はコミットしない**（`.gitignore`にルート直下限定で`/qsoku`・`/qsoku.exe`・`/dist/`）
- **シェル連携・補完のスクリプト**は`internal/cli/shells/`に埋め込み（`go:embed`）で配る（mtqgの`internal/cli/completions/`と同じやり方）。対応シェルはbash・zsh・fish・pwsh（pwshは2026-09-27に追加。以前はPowerShell対象外としていたが、方針を改めた）
- **看板としてのREADME・`examples/`は実装完了後に作った。**
