# 改訂履歴

決めたことの時系列。書き換えて過去の理由を消さない（`README.md`「位置づけ」）。

## 2026-09-22　設計メモの作成

qsokuの設計が`note.md`としてまとまった。要点：

- 名前は`quick`と日本語の「速」を重ねた造語。先に検討した`soku`は既存のCLI（Node.js製・Go製の2つ）と衝突したため避けた。`sok`（学術用語SoKと紛れる）・`skq`（Pythonのライブラリと衝突）も見送り、`qsoku`に決めた
- `qsokufile`は1つだけ探す（複数候補があったときの優先順位という決まり自体を無くすため）。名前は小文字固定、`.`で始まる名前は管理用コマンドのために予約
- コマンドは常に`sh`で実行し、実行後の居場所を一時ファイル（`QSOKU_CWD_FILE`）経由でシェル関数に持ち帰らせる。`eval`でシェルコードを流し込む方式（zoxide・direnv）は、コマンド自身の標準出力と混ざるため採らなかった
- `//`は`qsokufile`の置いてある場所を表す。引用符の外・単語の先頭でだけ置き換える（URL・引用符の中は対象外）
- 実装はGo（純粋、cgoなし）。シェル側は居場所を持ち帰る数行だけ
- 対応シェルはbash・zsh・fish。**PowerShellは対象外**（`sh`実行が前提のため）

## 2026-09-23　mtqg段階2の題材として決定、初期構成の準備

- mtqg側で、サンプルPJの題材を（検討していた日本の営業日計算やGoコード可視化ツールではなく）qsokuに決定。理由はmtqg本体の`PLAN.md`「段階2：サンプルPJ」に記録
- モジュールパスを`github.com/amisonnet8/qsoku`に決定
- 公開範囲は、mtqg本体と同じ「公開の2段階」方針（最初からpublic、v0.1相当のタグを打つまではタグを打たない）を踏襲すると決定
- mtqgの導入は、qsokuのdevcontainerに人間がこのリポジトリをcloneして用意すると決定（初期構成の作業では入れない）
- **進捗管理の文書（PLAN.md相当）は作らない**と決定。現在地・次にやること・保留事項はmtqgのtodo・memoに入れる（このリポジトリの開発自体をmtqgで記録しながら進めることが、サンプルPJとしての目的の一つであるため）
- 初期構成（`CLAUDE.md`、`docs/`、ルートの設定ファイル、`.devcontainer/`、`.claude/`）を、mtqg本体のリポジトリのセッションが作成した

## 2026-09-23　`docs/reference/`への仕様の書き起こし（Step 1）

PJの計画（mtqgのtodo、9ステップ）のStep 1として、`note.md`の内容とmtqgのquestionで決めたことを、`docs/reference/qsokufile.md`・`qsokufile_ja.md`・`cli.md`・`cli_ja.md`に仕様として書き起こした。決めたこと・`note.md`からの変更点は次のとおり（判断の経緯はmtqgのquestion、IDは`mtqg show <ID>`で見られる）。

**mtqgのquestionで確認したこと：**
- `#`コメントは、行頭の空白・タブの後に`#`があってもコメント扱いにする（コマンド末尾の`# …`はqsokuではなく`sh`の解釈に任せる）
- 名前に使える文字は英数字・`.`・`_`・`-`。先頭の`.`は禁止
- 同じ名前が2行あるとエラー。ただし**`.add`は上書きに変更**（下記）
- 補完は`.shell`の出力に含め、1回の`eval`で済ませる（mtqgの`completion`のような独立コマンドにはしない）
- `.edit`は`$EDITOR`→無ければ`nano`→無ければエラー
- 居場所の持ち帰りは`pwd`（論理パス）。シェル本来の`cd`の既定動作に合わせた（`pwd -P`にする積極的な理由が無いため）
- CIはLinux・macOS（Windowsは対象外のまま）
- qsoku自身のエラーの終了コードは、**mtqgの慣習（0成功／1実行できなかった／2コマンドラインの誤り）に揃える**。独自の体系は作らない（`CLAUDE.md`「独自の決まりを増やさない」）。qsokufileが見つからない→1、名前が無い→2、`sh`が起動できない→1。コマンド実行後は`sh`の終了コードをそのまま返す（上書きしない、note.mdの決定のまま）

**この会話中の追加の決定（`AskUserQuestion`、mtqgにも記録済み）：**
- 名前の先頭に`-`を許す。qsokuは自分自身のコマンドラインオプションを持たない（`-h`・`--help`は無く、`.help`のみ）ので、`-`始まりの名前と衝突しない
- 補完が名前の一覧を得るための公開コマンド`.names`を新設。qsokufileが無い・壊れていても終了コード0・標準エラーに何も出さない（TABで行を壊さないため）
- バージョン表示の`.version`を管理用コマンドに追加（`distribution.md`の`debug.ReadBuildInfo`の決定を反映）

**仕様に書き下ろす際に決めた細部（自明な延長）：**
- 行の形式：前後の空白を除き最初の`:`までが名前、残りが（先頭の空白を除いて）コマンド。`:`が無い行はエラー（行番号つき）。CRLFの`\r`は解析前に除く。行の継続は無い
- 探索はカレントから親へファイルシステムのルートまで。gitのルートで止まらない
- `.where`はqsokufileの置いてあるディレクトリを表示。`.init`はカレントに既にあればエラー（親にあっても作れる）。`.rm`は名前が無ければエラー。`.edit`・`.where`は壊れたqsokufileでも動く（直せるように）
- `qsoku`引数なしは`.help`と同じ（オプションのつづりを持たない）
- **仕様のこの版には、qsokuの出力例（エラー文言など）は載せていない。** 実装後、Step 8で実測して確かめたものを載せる決まり（`docs/reference/README.md`「実行例と出力例は実際に動かして確かめたもの」）のため

**この会話で決めた運用（`.claude/rules/mtqg.md`にも反映）：** 人間に何かを尋ねたら、回答を受けたその場で`q add`→`q add <id> <回答>`→`q done`まで行い、ターンを返さずに続ける。`mtqg q list`を対話中の決定の完全な記録にするため（人間の指示、2026-09-23）。

## 2026-09-23　足場（Step 2）

mtqgのtodo（9ステップ）のStep 2として、ビルド・検査・CIの入口をmtqg本体（`/home/vscode/mtqg`）の段階1 Step 1に倣って作った。

- **決めたこと（mtqgのquestionで確認）：** `.claude/hooks/build.sh`（`.go`・`go.mod`・`go.sum`編集後に`make build`を走らせ、失敗をClaude Codeに返す）を入れ、`.claude/settings.json`（本来は人間が管理）へのフック追記もClaude Codeが行ってよい
- `go.mod`（`github.com/amisonnet8/qsoku`、`go 1.27`）、`cmd/qsoku/main.go`（引数を渡すだけ）、`internal/cli`（`Run`は`.version`だけ実装、それ以外は`qsoku: not implemented yet`で終了コード1という仮の振る舞い。Step 3〜7で本実装に置き換える）、`Makefile`（mtqg本体と同じターゲット構成）、`.github/workflows/ci.yml`（Linux・macOSのマトリクス。Windowsの`choco install make`はqsokuには無いので外した）
- **`e2e/`がまだ無いため、`make test`は`if [ -d e2e ]`で存在確認してから実行するよう仮のガードを入れた。** Step 8で`e2e/`ができたらガードを外す（`.mtqg`のtodo fd80a62a28）
- バージョンは`runtime/debug.ReadBuildInfo`のみに頼る形にした（mtqgの`version.go`にある`-ldflags`用の変数は、qsokuにはまだリリースビルドの仕組みが無いため今回は入れていない。要れば`distribution.md`のGoReleaser検討時に足す）
- `.claude/rules/directory-structure.md`・`CLAUDE.md`「まだ無いもの」を更新。◯になったもの：`Makefile`・`go.mod`・`cmd/qsoku/`・`internal/cli/`・`.github/workflows/`・`.claude/hooks/`
- **手元で確かめたこと：** `make build`・`make check`・`make race`・`make trivy`・`make shellcheck`（対象ファイル無し）が通る。`golangci-lint config verify`が通る。フックはわざと構文エラーを入れて失敗を検知することを確認済み。`./qsoku .version`が`v0.0.0-<日時>-<コミットハッシュ>+dirty`を出し終了コード0、未実装の名前は終了コード1
