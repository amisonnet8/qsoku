# qsoku

*[English](README.md) | **日本語***

> **開発中です。まだ使えません。**
> qsokuは開発の途中です。リリースされたバージョンはまだなく（`@latest`は
> タグではなく最新のコミットを入れます）、データ形式やコマンドは予告なく
> 変わります。当面は、ご自身のプロジェクトで使わないでください。

**qsoku**は、リポジトリ固有のコマンドの近道を`qsokufile`に書いておくツール。
そのリポジトリのどこにいても（人間でもAIでも）`qsoku <名前>`で呼べる。

- シェルの`alias`と違い、`qsokufile`はリポジトリと一緒にコミットされるので、
  チーム全員が同じ近道を使え、ほかのプロジェクトに漏れ出すこともない
- makeと違い、依存関係のグラフも`.PHONY`もタブの決まりも無い——ただの
  名前の一覧で、近道の中で`cd`も普通に効く
- AIエージェントがリポジトリを読んでも、人間が見るのと同じ`qsokufile`が
  見える——どちらも同じ名前で呼べる一覧
- qsoku自身が決めているのは3つだけ：`qsokufile`の探し方、`//`の意味、
  居場所の持ち帰り方。残りはすべて普通の`sh`

## 例

<!-- qsoku:example dir=basic -->
```
$ qsoku hello world
hello, world
```

```
# qsokufile
hello: echo "hello, $1"
```

## インストール

```sh
go install github.com/amisonnet8/qsoku/cmd/qsoku@latest
```

シェルの起動ファイルに1行足す（bash・zsh・fishそれぞれの詳細は
[cli_ja.md](docs/reference/cli_ja.md#シェル連携)を参照）：

```sh
eval "$(qsoku .shell bash)"
```

これで、`qsokufile`に定義した名前のシェル補完も一緒に有効になる。

## もっと知る

- [docs/tour/](docs/tour/) — 空の`qsokufile`からシェル連携・補完までを
  手を動かしながら追う入門
- [docs/reference/](docs/reference/) — `qsokufile`の書式と`qsoku`
  コマンドラインの正式な仕様
- [docs/examples/](docs/examples/) — 自分のリポジトリにコピーして使える
  `qsokufile`の実例

## ライセンス

MIT（[LICENSE](LICENSE)参照）。
