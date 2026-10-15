# qsokuの使い方（このリポジトリでは make）

qsokuは、ビルド・テスト・lintなどの決まった作業を、リポジトリ直下の`qsokufile`に名前を付けて置き、`qsoku <名前>`で呼ぶツール。**このリポジトリはqsoku自身なので、自分の開発には`make`を使う**（`Makefile`が正）。開発中のqsokuが壊れても、ビルド・テストの足場は壊れない。

**ビルド・テスト・lintは、`go`・`golangci-lint`などを直接打たず、`make <名前>`で実行する。** Claude Codeは、Bashツールでは`bcuz`を通して呼ぶ：`bcuz "理由" make <名前>`（`.claude/rules/bcuz-usage.md`）。

## このリポジトリの名前（`Makefile`）

| 名前 | 内容 | いつ使うか |
|---|---|---|
| `make build` | `go build ./...`と`go build ./cmd/qsoku`（直下に`./qsoku`を作る） | 編集後のビルド確認（Goファイルの編集後は、`.claude/hooks/build.sh`が自動で呼ぶ） |
| `make fmt` | `golangci-lint fmt` | 整形を書き直したいとき |
| `make vet` | `go vet`（e2eのタグつきも） | |
| `make lint` | `golangci-lint run`（Goより古いlintで落ちたら、更新のしかたを案内する） | |
| `make unit` | `go test ./...` | |
| `make check` | `vet`・`lint`・`unit`をまとめて実行 | **作業の区切りで必ず通す**（`.claude/rules/testing.md`） |
| `make test` | e2e（本物のバイナリと本物のシェル） | シェル連携・実行に関わる変更 |
| `make docs-examples` | `docs/reference/`などの実行例を、実測の出力で書き直す | 実行例のある文書や、出力を変えたとき |
| `make race` | `CGO_ENABLED=1 go test -race` | 並行性・排他に関わる変更 |
| `make trivy` | 既知の脆弱性・ライセンスの検査 | 依存を足したとき |
| `make shellcheck` | 追跡中の`*.sh`・`*.bash`をShellCheckにかける | シェルスクリプトを足したとき |
| `make hook` | `.claude/hooks/require-bcuz.test.sh`（bcuz を強制するフックの判定の表。`jq`が要る） | `require-bcuz.sh`を変えたとき |
| `make goreleaser-check` | `.goreleaser.yaml`の検査と、公開しない試しのビルド | 配布の設定を変えたとき（`.claude/rules/distribution.md`） |

- **`windows-latest`のCIには`make`が無い。** `ci.yml`のWindowsのステップは、`make check`・`make test`の中身を直接呼ぶ。`Makefile`の`check`・`test`を変えたら、`ci.yml`のWindowsのステップも同じ変更の中で直す
- 名前を足したら、この表と、下の見本の`qsokufile`と、`.claude/rules/testing.md`の該当箇所も同じ変更の中で直す

## リポジトリ直下の`qsokufile`（見本）

- `Makefile`の主な名前を写した**見本**。このリポジトリの開発には使わない
- 名前は`Makefile`とほぼ同じ。ただし`test`は`go test ./...`（`Makefile`の`unit`）で、e2eは`e2e`という名前（看板READMEのデモGIFが、この`test`に依存するため）
- `.devcontainer/postCreate.sh`が、作業中のソースから入れたqsokuで、このqsokufileのシェル連携を`~/.bashrc`に入れる（動作確認・dogfooding のため）
- `e2e/run_test.go`が、この`qsokufile`を解析できることだけを確かめる（中身は実行しない）

## qsokuそのものの使い方

qsokuの使い方（呼び方、`qsokufile`の書き方、`//`、探し方、管理用コマンド、シェル連携、OSごとの注意）は、このリポジトリの仕様`docs/reference/`が正（英語が正、日本語版`*_ja.md`）。ルールのファイルには写さない（二重に持たない）。
