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

## 2026-09-23　qsokufileの探索と解析（Step 3）

`internal/qsokufile`パッケージを作り、`docs/reference/qsokufile.md`の「置き場所と探し方」「書式」「名前」を実装した（`//`の置き換えはStep 4）。`internal/`が初めて2パッケージになったので、`.golangci.yaml`にdepguardを足した（依存の向きは`cli`→`qsokufile`、逆はしない。mtqg本体の`journal`/`model`ルールと同じ形）。実際に`internal/cli`をimportさせてlintがdenyすることを確認してから戻した。

**仕様に明記が無く、実装時に決めた判断（history.mdに記録。仕様自体は変えない、Step 1の「自明な延長」の扱いを踏襲）：**
- **`qsokufile`という名前の**ディレクトリ**があった場合**：無いものとして扱い、親へ探索を続ける（エラーにしない）。迷ったら緩い方を選んだ。大文字小文字を区別しないファイルシステムの注意（`qsokufile.md`）と同じく、実運用で起こりうる紛れ込みに対して落ちないようにする判断
- `Find`が`os.Stat`で権限エラーなどそれ以外のエラーに遭遇したら、探索を続けずその場でエラーを返す（黙って親へ進むと、本当は読めるはずの`qsokufile`を見落とす恐れがあるため）
- `Parse`は最初に見つかった不正行で止まる（仕様の「エラーで終了する」の文言どおり、複数エラーをまとめて報告する仕組みは作らない）

**mtqg本体との比較：** `internal/journal/find.go`の`locate`は`.git`境界で止まる（qsokufileは境界で止まらない仕様なので実装はより単純）。「見つからない」ケースのテストは、mtqg自身も`t.TempDir()`から実際にファイルシステムのルートまで歩かせて確認しており（`/`近辺に紛れ込みが無い前提）、qsokuの`find_test.go`も同じ前提を踏襲した

**手元で確かめたこと：** `make check`・`make race`が通る。`go test ./internal/qsokufile/... -v`で表テスト（`Parse`14パターン、`Find`5パターン）が緑。カバレッジ77.1%（数値目標は設けていない。`filepath.Abs`の失敗などOSレベルの経路は未到達のまま残した）

## 2026-09-23　`//`の置き換え（Step 4）

`internal/qsokufile`に`Substitute`（qsokufileのコマンド文字列向け）・`SubstituteArg`（CLI引数向け）を追加し、`docs/reference/qsokufile.md`「`//`: the qsokufile's location」を実装した。

**実装中に見つけた、公開済み仕様の食い違い：** `qsokufile.md`の`//`の表の最終行（`build: (cd //; make build) # see //docs`は置き換わらない、理由は「コメントなので読まれない」）は、`Format`節の実際の決定（行末の`# …`はqsoku側で剥がさず、コマンド文字列にそのまま残る）と噛み合っていなかった。素朴な「引用符と単語の先頭だけを見る字句解析」だと、コメント部分の中の`//`も単語の先頭であれば置き換えてしまい、表の「No」と矛盾する。**`Substitute`に「引用符の外・単語の先頭で`#`に出会ったら、そこから行末までは走査を止める」という一手を足して、表の例をそのまま正しい動きにした**（`sh`自身がコメントの開始を判定する規則――単語の先頭にある`#`――と同じ条件を流用しただけなので、新しい決まりを増やしてはいない）。仕様書の文言自体は変更していない（今回の実装がその文言どおりの動きになった）。

**手作業での変異確認（`.claude/skills/mutation-check`相当。qsokuにはまだこのスキルを作っていないので、Step 3のdepguard確認と同じやり方で手で行った）：** ①`//`置き換え条件から`atWordStart`を外す→`curl http://example.com`等のテストが検知。②`#`の早期終了を外す→コメント関連のテストが検知。③`escape`の扱いを外す→**最初はテストが検知できなかった**（ダブルクォート内でエスケープされた`"`が誤って文字列を終わらせても、たまたま出力が変わらない入力例だったため）。テスト（`an escaped double quote does not close the string`）を、エスケープされた引用符の**後**に単語の先頭の`//`が来る例（`echo "x\" //y" z`）に作り直したところ、正しく検知するようになった。**手で変異確認をすると、見た目は妥当でも実は何も検証していないテストが見つかることがある、という実例**（mtqg本体への報告候補：スキル化した`mutate.sh`があれば、この種の「生き残ったが実は無意味」ケースをもっと機械的に洗い出せていた可能性がある）

**手元で確かめたこと：** `make check`・`make race`が通る。`go test ./internal/qsokufile/... -v`で`Substitute`17パターン・`SubstituteArg`6パターンが緑。変異確認3件がすべて（テスト強化後）検知される

## 2026-09-23　実行（Step 5）

`docs/reference/cli.md`「Running a name」を実装した。`internal/run`（`sh`を起動するだけの層。qsokufileが何かを知らない）を新設し、`internal/cli`を`qsokufile`・`run`の両方に配線して、`qsoku <名前>`が実際にqsokufileのコマンドを実行するようにした。

**Step 5の範囲をtodoの文言より広げた判断：** 元の9ステップの切り方では、「`.`で始まらない普通の名前を実際に実行する」配線がどの回にも明記されていなかった（Step 6は管理用コマンドのみ）。Step 5を「`sh`を起動する層の実装」だけでなく「`qsoku <名前>`を実際に動かす配線」まで含めることにした。理由：`Run`の居場所の持ち帰り・終了コードの素通しは、本物のCLI経路を通さないと実地で検証できない（`testing.md`の要求）。この結果、Step 6は管理用コマンドを追加するだけの純増分になる

**仕様の実装中に見つけた食い違い：** `cli.md`「Bringing the working directory back」は「`QSOKU_CWD_FILE`が未設定でも末尾の`pwd`行は常に足し、読むものが無いだけで無害」としていたが、実際に未設定なら`"$QSOKU_CWD_FILE"`は空文字列に展開され、`pwd > ""`は`sh`が「作れない」という**エラーを標準エラー出力に出す**（無害ではない）。**`QSOKU_CWD_FILE`が呼び出し元の環境に無いときは、持ち帰り行そのものを付けない**ことにして直した（結果として「無害」という仕様の意図どおりになる）。`cli.md`・`cli_ja.md`の該当箇所（`sh -c`の例、`QSOKU_CWD_FILE`の表の行）を実装に合わせて書き直した

**テストを書いていて見つけた、仕様に無かった注意点：** qsokufileの項目のコマンド自身が`exit`を呼ぶと、qsokuが末尾に足す持ち帰り行（同じスクリプトの続き）に到達する前にスクリプト全体が終わってしまう（`exit`はシェルを即座に終了させるため）。これはバグではなく実際の`sh`の挙動どおりだが、qsokufileを書く人には自明ではないので、`cli.md`「Running a name」に「項目は`exit`を呼ばず、最後のコマンド自身の終了コードに委ねること」という注意を追記した（`make`のレシピと同じ考え方）。テスト自身もこれにハマり、`exit 3`と書いていた箇所を`sh -c "exit 3"`（入れ子の`sh`が自然に返す終了コード）に直して気づいた

**Step 3の設計の見直し：** `qsokufile.Load`のエラーに、`fmt.Errorf("%s: %w", path, err)`でファイルパスを含めるよう変更した。Step 3時点では「pathは呼び出し側が付ける」としていたが、実際にCLI層でエラーメッセージを組み立てる段になり、`Format`節の「ファイル名と行番号を示して終了する」をそのまま満たすには`Load`自身が持つ方が自然だった。あわせて`File.Lookup(name)`を追加（実行に名前引きが要るため）

**手元で確かめたこと：** `go test ./internal/run/... -v`（本物の`sh`。成功・失敗・`$QSOKU_ROOT`の伝播・居場所の持ち帰り＝括弧あり/なし・`QSOKU_CWD_FILE`未設定時に標準エラー出力が空・シグナル終了で128+n・`PATH`を壊すと`sh`が見つからない扱いになる、の7系統）と`go test ./internal/cli/... -v`（実在するqsokufileでの実行・終了コードの素通し・名前が無い＝2・qsokufileが無い＝1）がすべて緑。`make check`・`make race`が通る。ビルドしたバイナリで手動確認：`QSOKU_CWD_FILE`無しで`//`の置き換えが効くこと、`QSOKU_CWD_FILE`ありで`cd //`の居場所が正しく持ち帰られること
