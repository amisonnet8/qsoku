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

## 2026-09-23　管理用コマンド（Step 6）

`docs/reference/cli.md`「Management commands」の残り（`.init` `.add` `.rm` `.list` `.names` `.edit` `.where` `.help`／引数なし）を実装した。`internal/qsokufile`に書き込み系（`SetEntry`・`RemoveEntry`。対象の1行だけ差し替える・削る、ほかは一切変えない）を足し、`internal/cli`はmtqg本体の`internal/cli/`に倣って1コマンド1ファイルにした。

**Step 6の範囲を広げた判断：** 元のtodoの文言には`.names`が無いが、`cli.md`の表では`.list`のすぐ下にあり実装もほぼ同じなので一緒に作った。

**未知の管理用コマンド・引数の数の扱い（この段階で決めた）：** `.version`〜`.help`は実装、**`.shell`は「認知しているが未実装」として`not implemented yet`・終了コード1のまま**（Step 7で実装予定と分かっているため）、それ以外の未知の`.foo`は終了コード2（`cli.md`の終了コード表どおり）。管理用コマンドの引数の数が違う場合も終了コード2。**`.help`と引数なしだけは、引数が何であっても常に成功する**（`.names`が壊れたqsokufileに強いのと同じ理由で、「困ったときに必ず助けになる」役割を優先した）。`.names`・`.list`・`.where`・`.init`・`.add`・`.rm`は通常どおり引数の数を検証する

**`.edit`の実装で決めたこと：** 新しいサブプロセス起動コードを書かず、**`internal/run.Execute`を再利用**した（`sh -c 'エディタ "$1"' qsoku <qsokufileのパス>`という形にすれば、居場所の持ち帰り・終了コードの素通しの仕組みがそのまま使える）。`$EDITOR`の値自体は引用符を付けずにスクリプト文字列へ埋め込む（`EDITOR="code --wait"`のような引数付きの値を活かすため）。対象ファイルのパスは`"$1"`経由で渡し、パスにスペースが含まれても壊れないようにした

**書き込み（`SetEntry`・`RemoveEntry`）の実装方針：** 生のファイル内容を`\n`で分割した`[]string`に対して、対象行だけを`Parse`が返す行番号で差し替える・削る、末尾に足すときはトレーリング改行の有無で場合分けする、という最小限の操作にした。元ファイルのパーミッションを`os.Stat`で読んで保つ

**テストを書いていて見つけた、実装とは無関係の落とし穴：** `.edit`のテストで、サブテスト名に`"runs $EDITOR on the qsokufile in use"`と書いたところ、Goの`t.TempDir()`がテスト名をほぼそのままディレクトリ名に使うため、生成されたパスの中に文字どおり`$EDITOR`という文字列が混入した。そのパスを（`.edit`の実装どおり）引用符なしで`sh -c`のスクリプト文字列に埋め込んだところ、`sh`がパス中の`$EDITOR`を本物のシェル変数として展開してしまい、パスが壊れて「not found」エラーになった。**実装のバグではなくテストの命名が原因**——サブテスト名に`$`・`` ` ``などシェルで特別な意味を持つ文字を含めない、という教訓。`.edit`が`$EDITOR`の値をそのまま（エスケープせず）スクリプトに埋め込む設計自体は正しい（現実の`$EDITOR`は普通のパスやコマンド名で、任意のシェル構文を意図的に許すのが仕様）

**手元で確かめたこと：** `go test ./internal/qsokufile/... -v`（`SetEntry`5パターン・`RemoveEntry`3パターン・パーミッション保持）と`go test ./internal/cli/... -v`（各コマンドごとの正常系・異常系）がすべて緑。`make check`・`make race`が通る。ビルドしたバイナリで`.init`→`.add`×2→`.list`→`.names`→`.where`→`.edit`（`EDITOR=cat`）→`.rm`→`.help`（引数なし）→未知の`.foo`を一通り手で実行し、期待どおりの出力・終了コードを確認

## 2026-09-23　シェル連携と補完（Step 7）

`docs/reference/cli.md`「Shell integration」「Shell completion」を実装した。`internal/cli/shells/`に3つのシェルスクリプト（`qsoku.bash`・`qsoku.zsh`・`qsoku.fish`。`go:embed`）を置き、`qsoku .shell <shell>`がそれを出す。

**e2e/をこの段階で前倒しした判断：** 元の計画ではe2e/はStep 8の担当だが、Step 7のtodo自体が「本物のbash・zsh・fishで確かめる」ことを要求しており、これにはStep 8が用意する予定の土台（本物にビルドしたバイナリ・本物の外部シェルプロセス）がそのまま要る。Step 5・6と同じ理由で、**e2e/の骨組み（`TestMain`でバイナリを1回ビルド。mtqg本体の`e2e/e2e_test.go`と同じ形）をStep 7で作った**。`Makefile`の`test`ターゲットから`[ -d e2e ]`ガードを外し、`vet`に`go vet -tags e2e ./e2e/...`を足した（`golangci-lint run`は`.golangci.yaml`の`run.build-tags: [e2e]`で元から見ていた）。Step 8はここに、名前の実行とdocsの例の確認を足すだけになる

**補完の情報源の按分：** mtqgは`candidates`という専用コマンドを新設して全補完を動的化したが、qsokuは「独自の決まりを増やさない」方針から、**管理用コマンドの固定一覧は3つのシェルスクリプトに直接書き、qsokufileの名前だけ`qsoku .names`で動的に取る**ことにした。スクリプトはバイナリに埋め込まれ一緒に配られるので、mtqgが動的化で避けた「スクリプトと本体の食い違い」はそもそも起きない。3つのスクリプトで一覧を手で揃える必要があるが、10個程度で滅多に変わらないため許容する

**実装中に見つけた、本物のシェルでしか出ない不具合（すべてe2eで検知→直した）：**
- **zshの`status`変数**：`qsoku()`関数のローカル変数名に`status`を使うと、zshは`$status`を`$?`の別名として予約しており、`local status`での上書きが「read-only variable」エラーになる（bashでは問題ない）。`qsoku_status`に変更した
- **zshの`compdef`**：`qsoku.zsh`の末尾は`$funcstack[1]`が`_qsoku`かどうかで分岐するが、単に`eval`しただけではその条件は常に偽になり`compdef`を呼ぶ。`compinit`を読み込んでいないシェルでは`compdef: command not found`になる——mtqg本体の同じ形のzshスクリプトも同じ構造で、mtqgのe2eは`compdef() { :; }`とスタブ化して確かめている。qsokuのe2eも同じやり方にした（実運用では`.zshrc`に元々`compinit`があるので、この問題はテスト環境固有）
- **fishの補完はドットファイル扱い**：fishは、補完候補が`.`で始まる場合、打っている語自体が`.`で始まるまで隠す（パス補完のドットファイル規則が、ファイルかどうかに関わらず先頭文字だけで適用される）。qsokuの管理用コマンドは全部`.`始まりなので、**fishでは`qsoku <TAB>`で管理用コマンドが出ない**——`qsoku .<TAB>`まで打って初めて出る。これはqsoku側のバグではなくfish自体の仕様。テストを`qsoku `と`qsoku .`の2ケースに分け、`cli.md`「Shell completion」に注意書きを追記した

**手元で確かめたこと：** `go test ./internal/cli/... -v`（`.shell`の埋め込み・引数検証）が緑。`make test`（`e2e`タグ、今回から実体化）で6つのテスト（bash・zsh・fishそれぞれの居場所の持ち帰り・補完。fishは名前と管理用コマンドで別ケース）がすべて緑。`make check`・`make race`・`make shellcheck`が通る（`qsoku.bash`のSC2164・SC2207を修正）。実バイナリを本物のbash・zsh・fishに`eval`させ、`qsoku`関数の居場所の持ち帰りと補完候補を手で確認

## 2026-09-23　e2eの残りと、docs/reference/の例を実測で確かめる仕組み（Step 8）

Step 7で前倒しした`e2e/`の土台に、残り2つを足した：本物のバイナリを直接実行するテスト（`e2e/run_test.go`）と、`docs/reference/`の実行例を実測で確かめる仕組み（`e2e/examples_test.go`、`make docs-examples`）。手本はmtqg本体の`e2e/examples_test.go`と`make docs-examples`

**`run_test.go`が要る理由：** `internal/cli`の単体テストは`cli.Run`をプロセス内で呼ぶだけなので、「ビルドした実バイナリの終了コード・シグナル・`QSOKU_CWD_FILE`の受け渡し」はまだ本物のプロセス境界を越えたことが無かった。終了コードの素通し（`sh -c "exit 7"`→7）・シグナル（`kill -TERM $$`→143）・qsokufile未発見/解析エラー/未定義名/引数の数違いの終了コード（1/1/2/2）・引数への`//`置き換え・`QSOKU_CWD_FILE`と標準出力が混ざらないこと・qsoku自身のエラーでは何も書かれないこと、を表テストで押さえた

**`examples_test.go`はmtqgの仕組みを簡素化して移植した：**
- **プロセス内実行の経路は作らなかった。** mtqgは`cli.Run`をプロセス内で呼ぶ経路（時計・環境を差し替えられる）と、実バイナリの経路の両方を持ち、`binary`オプションで両者の出力が一致することも確かめる。qsokuの`Run`は時計や環境の差し替え口をそもそも持たず、出力も時刻に依存しないので、**実バイナリだけで足りる**と判断した（比較対象が無いので`binary=`オプション自体も無い）
- **パスの置き換えは常時オン。** mtqgは`path=`オプションで選択制だが、qsokuの例は`.where`・`.init`・エラー文で絶対パスが出るものが大半なので、**すべての例で一時ディレクトリのパスを`/home/you/project`へ常に置き換える**ことにした
- **`$ `行は`qsoku`で始まるものだけに制限した。** mtqgは`git | mtqg`のようなパイプや`>/dev/null`も許すが、qsokuの例にはその必要が無いので、`checkInvocation`で単純に絞った
- **fixtureは言語非依存にした。** mtqgは`_ja`サフィックス付きのfixtureを別に持つ（記録の中身が日本語だとJSON Linesがそのまま日本語を含むため）。qsokuのfixture（qsokufile・付随ファイル）は中身が英語のコメント程度なので、英日どちらの文書からも同じfixtureを指す

**印の形式：** `<!-- qsoku:example dir=<fixture> [cwd=<subdir>] [skip="理由"] -->`（mtqgの`mtqg:example`と同じ形だが、`repo=`ではなく`dir=`、`ids=any`・`author=`はqsoku自身に該当する概念が無いので無し）。`TestDocExamplesAreMarkedAndMatch`が印の無い例・宙に浮いた印を検知し、英語版と`_ja`版で例の数・コマンド行が一致することを確かめる

**書き足した例：** `cli.md`「Running a name」（`qsoku hello world`・未定義名のエラー・サブディレクトリからの探索）「Management commands」（`.init`→`.add`×2→`.list`→`.names`→`.where`→`.rm`→`.list`の一続き）、`qsokufile.md`「Format」「Names」（`:`の無い行・名前の重複のエラー）「`//`」（表の実演：`url`・`grepq`・`echoq`で置き換わらないこと、`show //src`で引数が置き換わること）。出力はすべて`make docs-examples`が書き込んだもので、手では書いていない

**変異確認：** 文書の出力を1行わざと書き換えたら`TestDocExamples`が落ちる、印を1つ消したら`TestDocExamplesAreMarkedAndMatch`が落ちる（英語版とjaの例数が食い違うところまで検知する）ことを確認し、元に戻した

**gosecの扱い：** mtqg本体は`.golangci.yaml`で`_test.go$`全体からgosecを除外しているが、qsokuは既存方針（G204などをサイトごとに`//nolint:gosec`で理由付きにする、`.golangci.yaml`のコメント参照）を踏襲し、**除外ルールは追加せず**、`e2e/examples_test.go`・`e2e/run_test.go`の該当箇所に理由付きの`//nolint:gosec`を個別に付けた

**手元で確かめたこと：** `go test ./... `・`go test -tags e2e ./e2e/...`・`make check`・`make race`・`make shellcheck`がすべて緑。`make docs-examples`を2回連続で実行し、2回目で差分が出ない（出力が安定している）ことを確認

## 2026-09-23　仕上げと報告（Step 9）

実装の最後のステップとして、利用者向けの文書（看板README英日・`docs/tour/`・`docs/examples/`）、公開前の同名チェック、mtqgの「並行した状態変更の表示」の実地確認、そしてmtqg本体への報告を行った。

**公開前の同名チェック（2026-09-23）：** GitHub（`gh search repos qsoku`）は自分のリポジトリのみ、npm（`registry.npmjs.org/qsoku`）は404で空き、Go（`pkg.go.dev/search?q=qsoku`）は0件。紛らわしい先客がないことを確認し、`qsoku`のまま進めた

**並行した状態変更の表示の実地確認：** 前回（Step 3〜4の頃、memo `eca203d5a1`）はブランチとmainで別々の記録に触れただけで、`mtqg review`は何も検知しなかった。今回は意図的に、同じtodo（worktree側で`done`、main側で`edit`）と、既存のglossary語`qsokufile`（worktree側・main側でそれぞれ別の定義文を`g add`）という**同じ記録**を両側で変更してから`git merge --no-ff`した。結果：glossaryの重複定義は`mtqg review`の「Duplicate glossary definitions」にそのまま検知され、3つの定義が並んで表示された。一方、**同じtodoへの並行した変更（statusとedit）自体には、`mtqg review`・`mtqg status`のどちらにも印が付かない**——`mtqg show`で履歴を読んで初めて、2つの変更が別々の枝から来たかもしれないと気づける。journalのマージ自体は`.gitattributes`の`merge=union`どおり衝突なく自動マージされた。詳細はmemo `2b4deeb094`。実験用の記録（重複定義・検証用todo）は`mtqg delete`で片付けた（journalと履歴には残る）

**docs/reference/の実測の仕組みを広げた：** `e2e/examples_test.go`の`documents`（`docs/reference/`固定の4ファイル）を`documentPairs`（リポジトリルートからの相対パスの英日ペアの一覧）に一般化し、`README.md`/`README_ja.md`・`docs/tour/README.md`/`README_ja.md`を対象に追加した。仕組み自体（印・`qsoku`始まりの限定・パス置き換え）はStep 8のまま変えていない

**`docs/tour/`の構成：** 「何もない状態から始める」（`.init`）→「近道を足す」（`.add`・`.list`）→「実行する」（引数）→「`//`」→「どのサブディレクトリからでも」→「居場所を持ち帰る」（シェル連携の説明、非実行）→「シェル連携と補完」（非実行、cli.mdへリンク）→「直す・場所を確かめる」（`.where`・`.rm`）→「コミットする」の順。新しいfixture`tour`（`build`・`test`・`run`の3項目、`src/`つき）を使い、`.init`直後の2例だけは既存の`none`fixtureを再利用した

**`docs/examples/`の実例：** `go`・`node`・`monorepo`の3つ。`go`・`node`はよくある近道の一覧、`monorepo`は「括弧なしで意図的に移動する項目」と「括弧つきで移動しない項目」の対比、および`qsoku`から`qsoku`を呼ぶ項目を含む。これらは`go`・`npm`のような、このリポジトリ自身が依存しないツールを呼ぶため、**実行はしない**——`e2e/run_test.go`に`TestDocsExamplesQsokufilesParse`を足し、各ディレクトリで`qsoku .list`が終了コード0であること（解析できること）だけを確かめる

**手元で確かめたこと：** `make check`・`make test`・`make race`・`make shellcheck`・`make trivy`がすべて緑。`make docs-examples`を2回連続で実行し、2回目で差分が出ない。変異確認：`README.md`の出力を1行書き換えたら`TestDocExamples`が落ちる、`docs/examples/go/qsokufile`に`:`の無い行を足したら`TestDocsExamplesQsokufilesParse`が落ちる、両方確認して復元。実バイナリ・本物のbashで`docs/tour/`の手順（`.init`〜補完）を頭から1回手でなぞり、`src/`からの実行・TAB補完を含めて記載どおりであることを確認

**mtqg本体への報告：** `/home/vscode/mtqg-report.md`に書いた（人間の指示、q&a `c5dd00dfa8`。qsokuリポジトリにはコミットしない）

## 2026-09-23　正式公開の準備（GoReleaser・GitHubメタデータ）

Step 9完了後、人間から「正式公開する、基本的に今はやらないは無い、全部やる。トップのREADMEだけちゃんと考えて後で作る」との指示を受け、公開直前チェックリストのうち、トップの`README.md`/`README_ja.md`の内容変更を除く全項目を実施した（タグ付け自体は`distribution.md`の方針どおり人間の判断のまま）。

**GoReleaserの導入：** `.goreleaser.yaml`（新設）でlinux・darwinのamd64・arm64向けにクロスコンパイルする設定を書いた。**Windowsは対象外**——`testing.md`のCI検証OSの判断、`qsokufile`が`sh`実行前提でWindowsではGit Bash/WSLが要るという既存の判断と揃えた。`CGO_ENABLED=0`（distribution.mdの「純粋なGoにする」のまま）。アーカイブの中身は`LICENSE`・`README.md`・`README_ja.md`——GoReleaserの既定globに任せた（明示的な`files:`指定は書いていない）。実際にこの中身が入ることは、この会話のシェルへ`go install github.com/goreleaser/goreleaser/v2@v2.18.2`して`goreleaser release --snapshot --clean --skip=publish`を実行し、4アーキテクチャ分のtar.gzを展開して確認した（`dist/`はビルド後に削除、`.gitignore`済みでコミットはされない）。

**バージョンの埋め込み：** `internal/cli/version.go`にmtqg本体の`internal/cli/version.go`と同じ形で`var version string`を足し、`-ldflags`で埋め込まれていればそれを最優先、空なら従来どおり`debug.ReadBuildInfo()`に落ちる、という順にした。`go install`側の挙動（`ReadBuildInfo`に頼る理由）は変えていない。スナップショットビルドで`.version`が`v0.0.0`（GoReleaserがタグ無し状態に振る仮のタグ）を返すことを実機で確認した。`version_test.go`を新設し、`version`変数を設定した場合に`buildVersion()`がそれを返すことを確認する小テストを足した。

**CIの安全網：** `.github/workflows/ci.yml`に`goreleaser`ジョブを足し、通常のpush・PRのたびに`make goreleaser-check`（`goreleaser check`＋`--snapshot --skip=publish`のビルド）を実行するようにした。**タグを打つ前に設定の壊れを検知できる**——`.github/workflows/release.yml`（新設、`v*`タグのpushだけで動く）が実際に公開するのはタグを打った後なので、それより先にCIで気づけるようにする狙い（e2e/docs-examplesと同じ「実測で確かめる」方針を配布設定にも適用した）。`goreleaser-action`のバージョンは、実装時点の最新リリース（`gh api repos/goreleaser/goreleaser-action/releases/latest`で確認、`v7.2.3`）に合わせて`@v7`、CLI本体は`v2.18.2`に固定した（golangci-lint-action・setup-trivyと同じ「アクション本体はメジャー版タグ、中のツールは正確なバージョンを固定」という既存の流儀）。`Makefile`に`goreleaser-check`ターゲットを足し、`postCreate.sh`にも`goreleaser`のインストールを足した（次回のコンテナ再構築から手元でも使える）。

**GitHubリポジトリのメタデータ：** `gh repo edit`でdescription（"Per-repository command shortcuts: define them once in a qsokufile, run qsoku <name> from anywhere in the repo."）とTopics（`note.md`の下書き一覧12個：`cli`・`go`・`shell`・`bash`・`zsh`・`fish`・`developer-tools`・`productivity`・`shortcuts`・`command-runner`・`task-runner`・`dotfiles`）を設定した。`ai-agents`は`note.md`の元の判断どおり見送った（AIエージェント向けの機能・説明が実際にできてから足す）。

**`note.md`のチェック修正：** 「未決事項」の同名チェック項目が、Step 9で実施済みにもかかわらず`[ ]`のままだったので`[x]`に直した。

**手元で確かめたこと：** `make check`・`make test`・`make race`・`make shellcheck`・`make trivy`・新設の`make goreleaser-check`がすべて緑。`gh repo view --json description,repositoryTopics`でGitHub側の反映を確認。git tagは一切作っていない・pushしていない。

## 2026-09-23　看板READMEの作り直し（v0.1.0公開後）

v0.1.0の正式リリース後、唯一保留にしていたトップの`README.md`/`README_ja.md`を作り直した。人間から「映え対応」の具体的な指示（ヘッダー画像・バッジ・視覚コンテンツ・絵文字付き特徴一覧・折りたたみ・目次＋フッター）を受け、日本語版から作って完成させ、その後に英語版を作った。

**画像はすべて手書きのSVGにした：** devcontainerにラスター画像生成ツール（ImageMagick・PIL等）もvhs/asciinemaも無いことを確認済み。GitHubはREADME内のSVGをそのまま描画するため、`docs/assets/logo.svg`（ヘッダーのロゴ＋タグライン）と`docs/assets/demo.svg`（擬似ターミナルウィンドウ）の2つを新設した。どちらも`<style>`に`@media (prefers-color-scheme: dark)`を埋め込み、GitHubのライト/ダークどちらのテーマでも読める配色にした（ファイルを2本用意する`<picture>`方式より軽い）。

**`demo.svg`は「録画の代わりの仮置き」：** 本物のターミナル録画（GIF）は用意できないため、`docs/reference/cli.md`ですでに`make docs-examples`が実測・確認済みの`.init`→`.add`→`.list`の出力をそのまま使った静止画のSVGにした。中身は事実（検証済みの出力）そのものだが、実際の録画ではないという意味で人間の指示どおり「ダミー」として最後に作った。

**バッジ（Shields.io・pkg.go.dev）：** CI・最新リリース・ライセンス・Goバージョン・pkg.go.devの5つ。実装前に`curl`で全URLが200を返すことを確認した。

**折りたたみ（`<details>`）の中身も実測にした：** 「近道を足す一連の流れ」（`.init`→`.add`×2→`.list`）を新しく`<!-- qsoku:example dir=none -->`で印を付け、`docs/tour/`の同じ例と同一内容にした——README内のコード出力も手で書かない、という既存の方針をここにも適用した。

**開発中の注記の扱い：** v0.1.0を公開した後も、「`qsokufile`の書式やコマンドはまだ変わりうる」という警告文言自体は残した（タグを打った＝安定した、ではないため）。「バージョンが無い」という文だけを、`go install ...@latest`が今`v0.1.0`を指すという事実に合わせて書き換えた。

**確認したこと：** リンク切れチェック（既存のスクリプト。プレースホルダーの`xxx.md`以外に問題なし）、`TestDocExamplesAreMarkedAndMatch`（英日で例の数・コマンド文が一致）、`make docs-examples`→`TestDocExamples`（新設した折りたたみ内の例も含め全緑、2回目の実行で差分なし）、`make check`・`make test`・`make race`・`make shellcheck`。SVG 2つは`python3`の`xml.dom.minidom`で構文の妥当性を確認した（ブラウザでの見た目そのものはこの環境では確認できていない）。

## 2026-09-23　デモGIFへの差し替え（vhs）

看板READMEの「デモ」節を、静止画のSVG（`docs/assets/demo.svg`、実測済み出力を使った擬似ターミナル）から、`vhs`（charmbracelet/vhs）で録画した本物のGIF（`docs/assets/demo.gif`）に差し替えた。手順は人間が別プロジェクト（san-db-ox）で詰まった際にまとめた個人メモ`/home/vscode/vhs-setup-notes.md`にそのまま沿った：`ttyd`（GitHub Releasesの静的バイナリ）・`ffmpeg`・ヘッドレスChromium用共有ライブラリ・`xvfb`をaptで導入、`vhs`は`@latest`（v0.12.0、GIFを生成しない既知の回帰バグがある）ではなく`v0.11.0`を明示指定、実行は`VHS_NO_SANDBOX=true timeout 90 xvfb-run -a vhs demo.tape`（`--no-sandbox`だけではヘッドレスChromiumがハングする）。このセットアップ自体はdevcontainerの再構築で消える一回限りのもので、`.claude/rules/`には入れていない（vhs-setup-notes.md自身の方針どおり）。

**録画内容は人間の指示で2回作り直した：** 最初は「近道を足す一連の流れ」（`.init`→`.add`→`.list`）を考えたが、①「実践的で`//`を使う場面がよい」、②「録画に`make`を出さないほうがよい」（`qsoku`が`make`の薄いラッパーに見え、「makeの真似をしない」というqsoku自身の立ち位置と矛盾して見えるため）との指摘を受け、最終的に**このリポジトリ自身に本物の`qsokufile`（`build: (cd //; go build ./...)` / `test: (cd //; go test ./...)`）を新設し、それをドッグフーディングとして録画に使う**形にした。`internal/cli/`という2階層下のディレクトリから`qsoku test`を実行し、`//`がリポジトリのルートを見つけて`go test ./...`（3パッケージぶんの`ok`）が実際に走るところを見せる——`//`の一番実用的な使い方をそのまま実演する。この`qsokufile`は録画用の小道具ではなく、このリポジトリで今後も使える本物の近道としてコミットに残した。

**目視確認できた：** SVGのときはこの環境に画像レンダリング手段が無く構文チェックしかできなかったが、今回は`ffmpeg`でGIFから複数時点のフレームをPNGとして抜き出し、**Readツールで実際に画像として見て**、文字が読める・レイアウトが崩れていないことを確認した。最初の録画は`Set Height`が短すぎて上端が見切れていたため、高さを調整して録り直した。

**確認したこと：** `ffmpeg -i demo.gif`で長さ（約4.8秒）・解像度を確認、ファイルサイズは約77KB。`make check`・`make test`・既存のリンク切れチェックスクリプトが、新設した`qsokufile`・差し替えたGIFの影響を受けず緑のまま。tapeファイル自体はリポジトリに残していない。

## 2026-09-23　看板READMEの「## Example」節を実践的な内容に変更

人間から、看板README冒頭の「## Example」（デモGIFではなく、その前にある実行例）が「単純すぎて、かつトリッキー」との指摘。具体的には`hello: echo "hello, $1"`（`hello world`のおもちゃの例に、位置引数`$1`という初見には分かりにくい要素が混ざっていた）。

`basic`fixtureから、既に`docs/tour/`で使っている`tour`fixture（`build`・`test`という実プロジェクトらしい名前の近道、`$1`を使わない）に差し替えた。`test`の項目（`(cd //; echo "testing from $QSOKU_ROOT")`）が`//`を使っているので、「実践的で`//`を使う場面」というデモGIFのときと同じ人間の意向にも沿う。`.add`で`hello`を足す`<details>`節（`echo "hello, $1"`が残っている）は、`.add`の引数の渡し方を教える別の目的の節なので変更対象にしていない——今回指摘があったのは冒頭の「## Example」節のみ。

`make docs-examples`で出力を実測し直し、`make check`・`make test`（`TestDocExamplesAreMarkedAndMatch`含む）が緑であることを確認した。
