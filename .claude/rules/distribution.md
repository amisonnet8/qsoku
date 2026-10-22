# 配布方法

## `go install`（基本の入れ方）

```
go install github.com/amisonnet8/qsoku/cmd/qsoku@latest
```

モジュールパスは`github.com/amisonnet8/qsoku`。

- **タグを打つこと自体は人間が行う**（Claude Codeの`git push`は拒否設定でもある）。タグを打つとGoのモジュールプロキシ（`proxy.golang.org`）にそのバージョンが記録され、後から消せない——未完成の版を誤って`@latest`にしないよう、公開したい状態が整うまでタグは打たない（v0.1.0を打つまでの間に踏んだ判断。mtqg本体の「公開の2段階」方針を踏襲）

## 配布物の作り方（goreleaser）

`go install`が使えない・使いたくない人向けに、`cmd/qsoku`のビルド済みの配布物を、`.goreleaser.yaml`・`.github/workflows/release.yml`・`tools/licenses`で作る（テンプレートと同じ形。見本は crube（https://github.com/amisonnet8/crube ））。

- **cgo を使わない**（`CGO_ENABLED=0`。下の「純粋なGoにする」）
- 配布物は linux・darwin・windows の amd64・arm64。linux・darwin は`.tar.gz`、windows は`.zip`。名前は`<名前>_<版>_<OS>_<アーキテクチャ>`（版は`v`なし）。`checksums.txt`を添える。アーカイブには、実行ファイル・`LICENSE`・`README.md`・`README_ja.md`・`licenses/`（依存のライセンス文）を入れる
- **版**：goreleaser が`-X main.version=<タグ>`（`v0.1.0`の形）を埋める。`main`パッケージに`var version string`を置き、空のとき（`go install`で入れたもの）は`runtime/debug.ReadBuildInfo`のモジュールの版を出すと、どちらでも同じ`v0.1.0`になる（crube の`cli.BuildVersion`。下の「バージョンの取り方」）
- **依存のライセンス文**を配布物に同梱する。`tools/licenses`が、3つの OS で実行ファイルに組み込まれるモジュールを`go list -deps`で求め、そのライセンスのファイルと Go の標準ライブラリの`LICENSE`を`bin/licenses/<モジュールのパス>/`に写す（goreleaser の`before`の手順）。ライセンスのファイルが無いモジュールがあれば失敗し、リリースが止まる
  - 依存を足す・外すと、`tools/licenses`の`TestRun`の一覧（`./cmd/qsoku`のもの。今は`go/LICENSE`だけ）が落ちる。新しいモジュールのライセンスを確かめてから一覧を直す（Trivy の`make trivy`も通す）
- **Windows でも**、`qsokufile`のコマンドは`sh`で動く（Git for Windows が要る）

## リリースの流れ

ビルド済みのバイナリなどを GitHub Releases で配るときは、次の流れにする。どのプロジェクトも同じ。

1. **人間が**タグ（`v0.1.0`の形）を打って push する。人間がするのはこれだけ
2. タグで動く workflow（`.github/workflows/release.yml`）が、検査を通してから配布物を作り、Release に添付して**公開する**（本文は goreleaser が前のタグからのコミットで作る）

- **Claude Codeは、タグも Release も作らない。** `gh release create`は`.claude/settings.json`で拒否（deny）されている。Release を作るのは workflow の中だけ
- `release.yml`は、まず`ci.yml`をそのまま呼んで（`workflow_call`）普段の CI の検査をすべて通し、そのあと goreleaser（`.goreleaser.yaml`。`goreleaser-action`の`version: "~> v2"`で、v2 の最新）で配布物を作る。書き込みの権限（`contents: write`）は、Release を作るジョブだけに与える
- goreleaser はいったん**下書き**に添付し、すべて添付し終えてから、最後の手順（`gh release edit --draft=false`）で公開する。途中で落ちると下書きだけが残り、公開された Release に配布物が欠けることはない。同じタグでもう一度動かすと、その下書きを置き換える（`replace_existing_draft`）。落ちたときは、原因を直したコミットに新しいタグを打つか、同じコミットのままなら workflow を再実行する
- **公開済みの Release があるタグでは、何も変えずに失敗する**（最初の手順）。出し直すときは、新しいタグを打つ
- **タグを打つ前に、手元で`make goreleaser-check`（テンプレートの`qsoku snapshot`に当たる。このリポジトリの開発は make）を動かして中身を見る**。`.github/workflows/ci.yml`の`goreleaser`ジョブも、普段の push・PR のたびに同じものを動かす。タグなしで`dist/`に同じ配布物を作る（版は前のタグの次のパッチに`-snapshot`、バイナリの版は前のタグ。v0.2.0のあとなら`0.2.1-snapshot`と`v0.2.0`で、2026-10-22に実測。タグが1つも無ければ`0.0.1-snapshot`と`v0.0.0`）。goreleaser は`.devcontainer/postCreate.sh`が入れる。コンテナを作り直す前に要るときは、`GOTOOLCHAIN=auto go install github.com/goreleaser/goreleaser/v2@latest`（goreleaser がコンテナより新しい Go を求めるとき、その Go を取ってきてビルドする。Bashサンドボックスの中では、`proxy.golang.org`と`storage.googleapis.com`への接続が要る）。`.goreleaser.yaml`を変えたら`make goreleaser-check`で確かめる
- ビルド済みのバイナリは、リポジトリにコミットしない（`.claude/rules/directory-structure.md`「配布物」）

## 純粋なGoにする（cgoを使わない）

- devcontainerは`CGO_ENABLED=0`
- cgoを要する依存を足さない（`-race`だけは例外）

## バージョンの取り方

- `go install`で入れた`qsoku`のバージョンは`runtime/debug.ReadBuildInfo`で取る（`go install`では`-ldflags`の埋め込みが効かないため、それだけに頼らない）
- **GoReleaserが作るバイナリは`-ldflags`でバージョンを埋め込む**（goreleaser の`-X main.version={{ .Tag }}`。上の「配布物の作り方」）。`cmd/qsoku`の`main`が`version`を`cli.Version`に渡し、`buildVersion()`はこれを最優先で見て、空なら従来どおり`ReadBuildInfo`に落ちる——`go install`側の挙動は変えていない

## リポジトリにバイナリをコミットしない

- ビルドした`qsoku`、`dist/`は`.gitignore`に入れる（済み）
