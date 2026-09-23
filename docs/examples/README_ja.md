# 実例

*[English](README.md) | **日本語***

自分のリポジトリにそのままコピーして調整できる、動く`qsokufile`。qsoku
自体の使い方は[docs/tour/](../tour/)、正式な仕様は
[docs/reference/](../reference/)を参照。

| 実例 | 何を示すか |
|---|---|
| [go/qsokufile](go/qsokufile) | よくあるGoの近道（`build`・`test`・`race`・`lint`・`run`）と、どこにいても`//`でリポジトリのルートで`go mod tidy`を実行する例 |
| [node/qsokufile](node/qsokufile) | よくあるnpmの近道と、`//`でリポジトリのルートからビルド成果物を安全に`rm -rf`する例 |
| [monorepo/qsokufile](monorepo/qsokufile) | 複数パッケージのリポジトリ：意図的に移動する括弧なしの項目、移動せずに各パッケージをビルド・テストする括弧つきの項目、`qsoku`でほかの項目を呼ぶ項目 |
