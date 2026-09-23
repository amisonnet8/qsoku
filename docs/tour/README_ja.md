# 入門

*[English](README.md) | **日本語***

空のリポジトリから、シェル連携・補完までを手を動かしながら追う入門。
正式な仕様は[docs/reference/](../reference/)を、そのままコピーして使える
実例は[docs/examples/](../examples/)を見ること。

## 何もない状態から始める

まだ`qsokufile`が無いリポジトリ：

<!-- qsoku:example dir=none -->
```
$ qsoku .init
Created /home/you/project/qsokufile
```

`qsoku .init`はカレントディレクトリに空の`qsokufile`を作る。ただのテキスト
ファイルなので、好きなエディタで開いてもよいし、`qsoku .add`（次の節）を
使ってもよい。

## 近道を足す

<!-- qsoku:example dir=none -->
```
$ qsoku .init
Created /home/you/project/qsokufile
$ qsoku .add build 'echo "building the project"'
$ qsoku .add hello 'echo "hello, $1"'
$ qsoku .list
build: echo "building the project"
hello: echo "hello, $1"
```

`qsokufile`の各行は`名前: コマンド`。`qsoku .add`はこれを1行足す（すでに
その名前があれば、その場で置き換える）だけで、手で書いたコメントやほかの
行の並び順にはいっさい触れない。

## 実行する

<!-- qsoku:example dir=tour -->
```
$ qsoku build
building the project
$ qsoku run everyone
running with arg: everyone
```

名前の後の引数はすべて、`$1`・`$2`……として、そのまま`sh`自身が展開する
のと同じようにコマンドへ渡る。コマンドは`qsoku`をどのシェルから打っても
必ず`sh`で実行されるので、`qsokufile`はチーム全員にとって同じように動く。

## `//`はqsokufile自身の置き場所

<!-- qsoku:example dir=tour -->
```
$ qsoku test
testing from /home/you/project
```

`test`の項目は`(cd //; echo "testing from $QSOKU_ROOT")`：`//`は使われて
いる`qsokufile`が置いてあるディレクトリへ展開されるので、どこから実行
しても近道は必ずリポジトリのルートへ戻る道を知っている。括弧はその
`cd`をサブシェルの中に閉じ込める（下も参照）。

## どのサブディレクトリからでも同じqsokufile

<!-- qsoku:example dir=tour cwd=src -->
```
$ qsoku build
building the project
```

qsokuはカレントディレクトリから親へたどって`qsokufile`を見つける——gitの
`.git/`と同じように、その下のどこからでも同じ1つのファイルが使われる。

## 居場所を持ち帰る

`qsokufile`の項目は別プロセスで動くので、中の普通の`cd`だけでは、打った
人自身のシェルの居場所は本来変わらない。下のシェル連携（`qsoku .shell`）
がこれを解決する：名前を実行した後、連携用の関数がqsokuの終わった場所を
読み戻し、自分で`cd`する（成功でも失敗でも）。上の`test`のように、近道の
中だけで移動させたい`cd`は括弧でくくる。

正確な規則（`cd`が失敗したとき、シンボリックリンクなど）は
[cli_ja.md](../reference/cli_ja.md#居場所を持ち帰る)を参照。

## シェル連携と補完

シェルの起動ファイルに1行足す：

```sh
eval "$(qsoku .shell bash)"   # ~/.bashrc（zshは~/.zshrc）
```

```fish
qsoku .shell fish | source    # ~/.config/fish/config.fish
```

これは`qsoku`というシェル関数を定義し、本物のバイナリを実行し、居場所を
持ち帰り、さらに使われている`qsokufile`の名前（と管理用コマンド）への
TAB補完も有効にする。関数が実際にしていることは
[cli_ja.md](../reference/cli_ja.md#シェル連携)、シェルごとの補完の挙動は
[cli_ja.md](../reference/cli_ja.md#シェル補完)を参照（fishは特に、先頭が
`.`を打つまで管理用コマンドを隠す）。

## 直す、場所を確かめる

<!-- qsoku:example dir=tour -->
```
$ qsoku .where
/home/you/project
$ qsoku .rm run
$ qsoku .list
build: echo "building the project"
test: (cd //; echo "testing from $QSOKU_ROOT")
```

`qsoku .where`は使われている`qsokufile`が置いてあるディレクトリ（`//`が
展開される先）を表示し、`qsoku .edit`は`$EDITOR`でそれを開く。どちらも、
今の`qsokufile`が解析できない状態でも動く——直すために必要だから。

## コミットする

`qsokufile`は`Makefile`と同じく、コミットされることを前提にしている。
一度リポジトリに入れば、チームの誰もが——そのリポジトリで動くAIエージェント
も含めて——同じ近道の一覧を見て、同じように呼べる。

## この先

- 名前に使える文字・コメント・終了コードなどの正式な仕様は
  [docs/reference/](../reference/)
- 自分のプロジェクトにコピーして使える`qsokufile`は
  [docs/examples/](../examples/)
