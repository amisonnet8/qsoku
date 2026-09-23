# ディレクトリ構成

## 今あるもの

```
qsoku/
├── CLAUDE.md               ← プロジェクトルール（参照先の案内）
├── LICENSE                 （MIT）
├── .gitignore
├── .gitattributes          （`* text=auto eol=lf`）
├── .golangci.yaml           ← lintの設定
├── trivy.yaml               ← 脆弱性・ライセンス検査の設定
├── README.md / README_ja.md ← 未完成の間の注意書き（看板のREADMEは最後に作る）
├── docs/
│   ├── design/              ← 設計判断と理由の記録（日本語）。note.md・README.md・history.md
│   └── reference/            ← 仕様（英語が正、*_ja.mdと2本立て）。まだ置き場所の決まりだけ
├── .devcontainer/
└── .claude/                 ← rules/、settings.json（人間が管理）
```

## まだ無いもの（実装しながら決める）

- `Makefile`、`go.mod`／`go.sum`
- `cmd/qsoku/main.go`
- `internal/` — `qsokufile`の解析・`//`の置き換え・`sh`の起動・管理用コマンド（`.init`・`.add`・`.rm`・`.list`・`.edit`・`.where`・`.help`・`.shell`）をどう層に分けるかは実装しながら決めてよい。**層を決めたら`.golangci.yaml`にdepguardのルールを足す**（層の依存の向きを機械的に検査するため。mtqgの`.claude/rules/directory-structure.md`「配置の判断基準」が参考になる）
- `e2e/` — 本物のバイナリと本物のシェル（bash・zsh・fish）で動かすテスト
- `.github/workflows/`（CI）
- `.claude/hooks/`（`Makefile`ができたら、`.go`編集後に`make build`するフックを足せる）
- `docs/tour/`・`docs/examples/`（実装完了後）

## 配置の判断基準

- **`internal/`**：Goの仕組みとして、リポジトリの外からimportできない。外との約束は、配布物（`go install`で入るバイナリ）とデータ形式（`qsokufile`の書式）だけ
- **配布物（ビルド済みバイナリ）はコミットしない**（`.gitignore`にルート直下限定で`/qsoku`・`/qsoku.exe`・`/dist/`）
- **シェル連携・補完のスクリプト**は、埋め込み（`go:embed`）で配る想定（mtqgの`internal/cli/completions/`と同じやり方が参考になる）。対応シェルはbash・zsh・fish（PowerShellは対象外）
