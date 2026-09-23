# 配布方法

## v0.xの間は`go install`だけ

```
go install github.com/amisonnet8/qsoku/cmd/qsoku@latest
```

モジュールパスは`github.com/amisonnet8/qsoku`。

- **v0.1相当までタグを打たない。** リポジトリは最初からpublicで、タグを打つとGoのモジュールプロキシ（`proxy.golang.org`）にバージョンが記録され、後から消せない。未完成の版が`@latest`で入ってしまう（mtqg本体の「公開の2段階」方針を踏襲。`docs/design/history.md` 2026-09-23）
- GoReleaserなどでのバイナリ配布は、公開の段階で検討する。今は作らない

## 純粋なGoにする（cgoを使わない）

- devcontainerは`CGO_ENABLED=0`
- cgoを要する依存を足さない（`-race`だけは例外）

## バージョンの取り方

- `qsoku`のバージョンは`runtime/debug.ReadBuildInfo`で取る（`go install`では`-ldflags`の埋め込みが効かないため、それだけに頼らない）

## リポジトリにバイナリをコミットしない

- ビルドした`qsoku`、`dist/`は`.gitignore`に入れる（済み）
