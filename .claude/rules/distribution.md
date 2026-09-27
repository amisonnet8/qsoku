# 配布方法

## `go install`（基本の入れ方）

```
go install github.com/amisonnet8/qsoku/cmd/qsoku@latest
```

モジュールパスは`github.com/amisonnet8/qsoku`。

- **タグを打つこと自体は人間が行う**（Claude Codeの`git push`は拒否設定でもある）。タグを打つとGoのモジュールプロキシ（`proxy.golang.org`）にそのバージョンが記録され、後から消せない——未完成の版を誤って`@latest`にしないよう、公開したい状態が整うまでタグは打たない（v0.1.0を打つまでの間に踏んだ判断。mtqg本体の「公開の2段階」方針を踏襲。`docs/design/history.md` 2026-09-23）

## GoReleaserでのバイナリ配布（2026-09-23に用意）

`.goreleaser.yaml`（リポジトリ直下）と`.github/workflows/release.yml`で、`go install`が使えない・使いたくない人向けに、ビルド済みバイナリを配れるようにしてある。**動くのは人間が`v*`のタグをpushしたときだけ**（`release.yml`のトリガー）。

- **対象OS・アーキテクチャはlinux・darwin・windowsのamd64・arm64**（Windowsはv0.2.0から。`testing.md`のCI検証OSと揃えた。`qsokufile`のコマンド自体は引き続きWindowsでも`sh`実行——Git for Windowsが要る——のまま。経緯は`docs/design/history.md` 2026-09-27）
- 純粋なGo（`CGO_ENABLED=0`）でクロスコンパイルする（本節の「純粋なGoにする」のまま）
- 配布物は、windows向けだけ`.zip`（`format_overrides`）、それ以外は`.tar.gz`1本（アーカイブの中身はGoReleaserの既定glob——`LICENSE*`・`README*`——に任せている。`LICENSE`・`README.md`・`README_ja.md`、windowsでは`qsoku.exe`が入る）
- **設定が壊れていないかは、タグを打つ前にCIで分かる。** `.github/workflows/ci.yml`の`goreleaser`ジョブが、通常のpush・PRのたびに`make goreleaser-check`（`goreleaser check`＋`goreleaser release --snapshot --clean --skip=publish`。タグ不要・公開なし）を実行する。手元でも同じコマンドで確認できる（`goreleaser`は`postCreate.sh`で入る）

## 純粋なGoにする（cgoを使わない）

- devcontainerは`CGO_ENABLED=0`
- cgoを要する依存を足さない（`-race`だけは例外）

## バージョンの取り方

- `go install`で入れた`qsoku`のバージョンは`runtime/debug.ReadBuildInfo`で取る（`go install`では`-ldflags`の埋め込みが効かないため、それだけに頼らない）
- **GoReleaserが作るバイナリは`-ldflags`でバージョンを埋め込む**（`internal/cli/version.go`の`var version string`、mtqg本体の`internal/cli/version.go`と同じ形。`.goreleaser.yaml`が`{{.Tag}}`を渡す）。`buildVersion()`はこの`version`を最優先で見て、空なら従来どおり`ReadBuildInfo`に落ちる——`go install`側の挙動は変えていない

## リポジトリにバイナリをコミットしない

- ビルドした`qsoku`、`dist/`は`.gitignore`に入れる（済み）
