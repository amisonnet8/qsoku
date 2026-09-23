# 配布方法

## v0.xの間は`go install`だけ

```
go install github.com/amisonnet8/qsoku/cmd/qsoku@latest
```

モジュールパスは`github.com/amisonnet8/qsoku`。

- **v0.1相当までタグを打たない。** リポジトリは最初からpublicで、タグを打つとGoのモジュールプロキシ（`proxy.golang.org`）にバージョンが記録され、後から消せない。未完成の版が`@latest`で入ってしまう（mtqg本体の「公開の2段階」方針を踏襲。`docs/design/history.md` 2026-09-23）。**タグを打つこと自体は人間が行う**（Claude Codeの`git push`は拒否設定でもある）

## GoReleaserでのバイナリ配布（2026-09-23、Step 9で用意）

`.goreleaser.yaml`（リポジトリ直下）と`.github/workflows/release.yml`で、`go install`が使えない・使いたくない人向けに、ビルド済みバイナリを配れるようにしてある。**動くのは人間が`v*`のタグをpushしたときだけ**（`release.yml`のトリガー）。

- **対象OS・アーキテクチャはlinux・darwinのamd64・arm64のみ（Windowsは対象外）。** CIの検証OS（`testing.md`）、`qsokufile`が常に`sh`実行前提であること（Windowsでの利用にはGit Bash/WSLが要る）と揃えた
- 純粋なGo（`CGO_ENABLED=0`）でクロスコンパイルする（本節の「純粋なGoにする」のまま）
- 配布物は`.tar.gz`1本（アーカイブの中身はGoReleaserの既定glob——`LICENSE*`・`README*`——に任せている。`LICENSE`・`README.md`・`README_ja.md`が実際に入ることを`make goreleaser-check`で確認済み）
- **設定が壊れていないかは、タグを打つ前にCIで分かる。** `.github/workflows/ci.yml`の`goreleaser`ジョブが、通常のpush・PRのたびに`make goreleaser-check`（`goreleaser check`＋`goreleaser release --snapshot --clean --skip=publish`。タグ不要・公開なし）を実行する。手元でも同じコマンドで確認できる（`goreleaser`は`postCreate.sh`で入る）

## 純粋なGoにする（cgoを使わない）

- devcontainerは`CGO_ENABLED=0`
- cgoを要する依存を足さない（`-race`だけは例外）

## バージョンの取り方

- `go install`で入れた`qsoku`のバージョンは`runtime/debug.ReadBuildInfo`で取る（`go install`では`-ldflags`の埋め込みが効かないため、それだけに頼らない）
- **GoReleaserが作るバイナリは`-ldflags`でバージョンを埋め込む**（`internal/cli/version.go`の`var version string`、mtqg本体の`internal/cli/version.go`と同じ形。`.goreleaser.yaml`が`{{.Tag}}`を渡す）。`buildVersion()`はこの`version`を最優先で見て、空なら従来どおり`ReadBuildInfo`に落ちる——`go install`側の挙動は変えていない

## リポジトリにバイナリをコミットしない

- ビルドした`qsoku`、`dist/`は`.gitignore`に入れる（済み）
