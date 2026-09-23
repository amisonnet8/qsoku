# 改訂履歴

決めたことの時系列。書き換えて過去の理由を消さない（`README.md`「位置づけ」）。検証ログや進行管理のメモ（「手元で確かめたこと」など）はここには置かない——テストの方法自体は`testing.md`に、個々の実行結果はmtqgの記録にある。ここに残すのは、qsoku自身が**なぜ今の形になっているか**の理由だけ。

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
- **進捗管理の文書（PLAN.md相当）は作らない**と決定。現在地・次にやること・保留事項はmtqgのtodo・memoに入れる

## 2026-09-23　`docs/reference/`への仕様の書き起こし

`note.md`の内容とmtqgのquestionで決めたことを、`docs/reference/qsokufile.md`・`qsokufile_ja.md`・`cli.md`・`cli_ja.md`に仕様として書き起こした。決めたこと・`note.md`からの変更点は次のとおり（判断の経緯はmtqgのquestion、IDは`mtqg show <ID>`で見られる）。

**mtqgのquestionで確認したこと：**
- `#`コメントは、行頭の空白・タブの後に`#`があってもコメント扱いにする（コマンド末尾の`# …`はqsokuではなく`sh`の解釈に任せる）
- 名前に使える文字は英数字・`.`・`_`・`-`。先頭の`.`は禁止
- 同じ名前が2行あるとエラー。ただし**`.add`は上書きに変更**（下記）
- 補完は`.shell`の出力に含め、1回の`eval`で済ませる（mtqgの`completion`のような独立コマンドにはしない）
- `.edit`は`$EDITOR`→無ければ`nano`→無ければエラー
- 居場所の持ち帰りは`pwd`（論理パス）。シェル本来の`cd`の既定動作に合わせた（`pwd -P`にする積極的な理由が無いため）
- CIはLinux・macOS（Windowsは対象外のまま）
- qsoku自身のエラーの終了コードは、**mtqgの慣習（0成功／1実行できなかった／2コマンドラインの誤り）に揃える**。独自の体系は作らない（`CLAUDE.md`「独自の決まりを増やさない」）。qsokufileが見つからない→1、名前が無い→2、`sh`が起動できない→1。コマンド実行後は`sh`の終了コードをそのまま返す（上書きしない、note.mdの決定のまま）

**この会話中の追加の決定：**
- 名前の先頭に`-`を許す。qsokuは自分自身のコマンドラインオプションを持たない（`-h`・`--help`は無く、`.help`のみ）ので、`-`始まりの名前と衝突しない
- 補完が名前の一覧を得るための公開コマンド`.names`を新設。qsokufileが無い・壊れていても終了コード0・標準エラーに何も出さない（TABで行を壊さないため）
- バージョン表示の`.version`を管理用コマンドに追加（`distribution.md`の`debug.ReadBuildInfo`の決定を反映）

**仕様に書き下ろす際に決めた細部：**
- 行の形式：前後の空白を除き最初の`:`までが名前、残りが（先頭の空白を除いて）コマンド。`:`が無い行はエラー（行番号つき）。CRLFの`\r`は解析前に除く。行の継続は無い
- 探索はカレントから親へファイルシステムのルートまで。gitのルートで止まらない
- `.where`はqsokufileの置いてあるディレクトリを表示。`.init`はカレントに既にあればエラー（親にあっても作れる）。`.rm`は名前が無ければエラー。`.edit`・`.where`は壊れたqsokufileでも動く（直せるように）
- `qsoku`引数なしは`.help`と同じ（オプションのつづりを持たない）

## 2026-09-23　qsokufileの探索と解析

`internal/qsokufile`パッケージを作り、`docs/reference/qsokufile.md`の「置き場所と探し方」「書式」「名前」を実装した。`.golangci.yaml`にdepguardを足し、依存の向きを`cli`→`qsokufile`（逆はしない）に固定した（mtqg本体の`journal`/`model`ルールと同じ形）。

**仕様に明記が無く、実装時に決めた判断：**
- **`qsokufile`という名前の**ディレクトリ**があった場合**：無いものとして扱い、親へ探索を続ける（エラーにしない）。大文字小文字を区別しないファイルシステムの注意（`qsokufile.md`）と同じく、実運用で起こりうる紛れ込みに対して落ちないようにする判断
- `Find`が`os.Stat`で権限エラーなどそれ以外のエラーに遭遇したら、探索を続けずその場でエラーを返す（黙って親へ進むと、本当は読めるはずの`qsokufile`を見落とす恐れがあるため）
- `Parse`は最初に見つかった不正行で止まる（仕様の「エラーで終了する」の文言どおり、複数エラーをまとめて報告する仕組みは作らない）

## 2026-09-23　`//`の置き換え

`internal/qsokufile`に`Substitute`（qsokufileのコマンド文字列向け）・`SubstituteArg`（CLI引数向け）を追加し、`docs/reference/qsokufile.md`「`//`: the qsokufile's location」を実装した。

**実装中に見つけた、公開済み仕様の食い違い：** `qsokufile.md`の`//`の表の最終行（`build: (cd //; make build) # see //docs`は置き換わらない、理由は「コメントなので読まれない」）は、`Format`節の実際の決定（行末の`# …`はqsoku側で剥がさず、コマンド文字列にそのまま残る）と噛み合っていなかった。素朴な「引用符と単語の先頭だけを見る字句解析」だと、コメント部分の中の`//`も単語の先頭であれば置き換えてしまい、表の「No」と矛盾する。**`Substitute`に「引用符の外・単語の先頭で`#`に出会ったら、そこから行末までは走査を止める」という一手を足して、表の例をそのまま正しい動きにした**（`sh`自身がコメントの開始を判定する規則――単語の先頭にある`#`――と同じ条件を流用しただけなので、新しい決まりを増やしてはいない）。仕様書の文言自体は変更していない（今回の実装がその文言どおりの動きになった）。

## 2026-09-23　実行

`docs/reference/cli.md`「Running a name」を実装した。`internal/run`（`sh`を起動するだけの層。qsokufileが何かを知らない）を新設し、`internal/cli`を`qsokufile`・`run`の両方に配線して、`qsoku <名前>`が実際にqsokufileのコマンドを実行するようにした。

**仕様の実装中に見つけた食い違い：** `cli.md`「Bringing the working directory back」は「`QSOKU_CWD_FILE`が未設定でも末尾の`pwd`行は常に足し、読むものが無いだけで無害」としていたが、実際に未設定なら`"$QSOKU_CWD_FILE"`は空文字列に展開され、`pwd > ""`は`sh`が「作れない」という**エラーを標準エラー出力に出す**（無害ではない）。**`QSOKU_CWD_FILE`が呼び出し元の環境に無いときは、持ち帰り行そのものを付けない**ことにして直した（結果として「無害」という仕様の意図どおりになる）。`cli.md`・`cli_ja.md`の該当箇所（`sh -c`の例、`QSOKU_CWD_FILE`の表の行）を実装に合わせて書き直した。

**テストを書いていて見つけた、仕様に無かった注意点：** qsokufileの項目のコマンド自身が`exit`を呼ぶと、qsokuが末尾に足す持ち帰り行（同じスクリプトの続き）に到達する前にスクリプト全体が終わってしまう（`exit`はシェルを即座に終了させるため）。これはバグではなく実際の`sh`の挙動どおりだが、qsokufileを書く人には自明ではないので、`cli.md`「Running a name」に「項目は`exit`を呼ばず、最後のコマンド自身の終了コードに委ねること」という注意を追記した（`make`のレシピと同じ考え方）。

**Load・Lookupの設計：** `qsokufile.Load`のエラーに、`fmt.Errorf("%s: %w", path, err)`でファイルパスを含めるようにした。CLI層でエラーメッセージを組み立てる段になり、`Format`節の「ファイル名と行番号を示して終了する」をそのまま満たすには`Load`自身が持つ方が自然だった。あわせて`File.Lookup(name)`を追加（実行に名前引きが要るため）。

## 2026-09-23　管理用コマンド

`docs/reference/cli.md`「Management commands」の残り（`.init` `.add` `.rm` `.list` `.names` `.edit` `.where` `.help`／引数なし）を実装した。`internal/qsokufile`に書き込み系（`SetEntry`・`RemoveEntry`。対象の1行だけ差し替える・削る、ほかは一切変えない）を足し、`internal/cli`はmtqg本体の`internal/cli/`に倣って1コマンド1ファイルにした。

**未知の管理用コマンド・引数の数の扱い：** 未知の`.foo`は終了コード2（`cli.md`の終了コード表どおり）。管理用コマンドの引数の数が違う場合も終了コード2。**`.help`と引数なしだけは、引数が何であっても常に成功する**（`.names`が壊れたqsokufileに強いのと同じ理由で、「困ったときに必ず助けになる」役割を優先した）。`.names`・`.list`・`.where`・`.init`・`.add`・`.rm`は通常どおり引数の数を検証する。

**`.edit`の実装で決めたこと：** 新しいサブプロセス起動コードを書かず、**`internal/run.Execute`を再利用**した（`sh -c 'エディタ "$1"' qsoku <qsokufileのパス>`という形にすれば、居場所の持ち帰り・終了コードの素通しの仕組みがそのまま使える）。`$EDITOR`の値自体は引用符を付けずにスクリプト文字列へ埋め込む（`EDITOR="code --wait"`のような引数付きの値を活かすため）。対象ファイルのパスは`"$1"`経由で渡し、パスにスペースが含まれても壊れないようにした。

**書き込み（`SetEntry`・`RemoveEntry`）の実装方針：** 生のファイル内容を`\n`で分割した`[]string`に対して、対象行だけを`Parse`が返す行番号で差し替える・削る、末尾に足すときはトレーリング改行の有無で場合分けする、という最小限の操作にした。元ファイルのパーミッションを`os.Stat`で読んで保つ。

## 2026-09-23　シェル連携と補完

`docs/reference/cli.md`「Shell integration」「Shell completion」を実装した。`internal/cli/shells/`に3つのシェルスクリプト（`qsoku.bash`・`qsoku.zsh`・`qsoku.fish`。`go:embed`）を置き、`qsoku .shell <shell>`がそれを出す。

**補完の情報源の按分：** mtqgは`candidates`という専用コマンドを新設して全補完を動的化したが、qsokuは「独自の決まりを増やさない」方針から、**管理用コマンドの固定一覧は3つのシェルスクリプトに直接書き、qsokufileの名前だけ`qsoku .names`で動的に取る**ことにした。スクリプトはバイナリに埋め込まれ一緒に配られるので、mtqgが動的化で避けた「スクリプトと本体の食い違い」はそもそも起きない。3つのスクリプトで一覧を手で揃える必要があるが、10個程度で滅多に変わらないため許容する。

**実装中に見つけた、本物のシェルでしか出ない不具合（すべて修正済み）：**
- **zshの`status`変数**：`qsoku()`関数のローカル変数名に`status`を使うと、zshは`$status`を`$?`の別名として予約しており、`local status`での上書きが「read-only variable」エラーになる（bashでは問題ない）。`qsoku_status`に変更した
- **zshの`compdef`**：`qsoku.zsh`の末尾は`$funcstack[1]`が`_qsoku`かどうかで分岐するが、単に`eval`しただけではその条件は常に偽になり`compdef`を呼ぶ。`compinit`を読み込んでいないシェルでは`compdef: command not found`になる——mtqg本体の同じ形のzshスクリプトも同じ構造で、mtqgのe2eは`compdef() { :; }`とスタブ化して確かめている。qsokuのe2eも同じやり方にした（実運用では`.zshrc`に元々`compinit`があるので、この問題はテスト環境固有）
- **fishの補完はドットファイル扱い**：fishは、補完候補が`.`で始まる場合、打っている語自体が`.`で始まるまで隠す（パス補完のドットファイル規則が、ファイルかどうかに関わらず先頭文字だけで適用される）。qsokuの管理用コマンドは全部`.`始まりなので、**fishでは`qsoku <TAB>`で管理用コマンドが出ない**——`qsoku .<TAB>`まで打って初めて出る。これはqsoku側のバグではなくfish自体の仕様。`cli.md`「Shell completion」に注意書きを追記した

## 2026-09-23　e2eとdocs/reference/の実測の仕組み

本物のバイナリを直接実行するテスト（`e2e/run_test.go`）と、`docs/reference/`の実行例を実測で確かめる仕組み（`e2e/examples_test.go`、`make docs-examples`）を作った。手本はmtqg本体の`e2e/examples_test.go`と`make docs-examples`。

**`run_test.go`が要る理由：** `internal/cli`の単体テストは`cli.Run`をプロセス内で呼ぶだけなので、「ビルドした実バイナリの終了コード・シグナル・`QSOKU_CWD_FILE`の受け渡し」は本物のプロセス境界を越えて確かめる必要があった。

**`examples_test.go`はmtqgの仕組みを簡素化して移植した：**
- **プロセス内実行の経路は作らなかった。** mtqgは`cli.Run`をプロセス内で呼ぶ経路と、実バイナリの経路の両方を持つが、qsokuの`Run`は時計や環境の差し替え口をそもそも持たず、出力も時刻に依存しないので、**実バイナリだけで足りる**と判断した
- **パスの置き換えは常時オン。** mtqgは`path=`オプションで選択制だが、qsokuの例は`.where`・`.init`・エラー文で絶対パスが出るものが大半なので、**すべての例で一時ディレクトリのパスを`/home/you/project`へ常に置き換える**ことにした
- **`$ `行は`qsoku`で始まるものだけに制限した。** mtqgは`git | mtqg`のようなパイプや`>/dev/null`も許すが、qsokuの例にはその必要が無いので、`checkInvocation`で単純に絞った
- **fixtureは言語非依存にした。** mtqgは`_ja`サフィックス付きのfixtureを別に持つ（記録の中身が日本語だとJSON Linesがそのまま日本語を含むため）。qsokuのfixture（qsokufile・付随ファイル）は中身が英語のコメント程度なので、英日どちらの文書からも同じfixtureを指す

**印の形式：** `<!-- qsoku:example dir=<fixture> [cwd=<subdir>] [skip="理由"] -->`（mtqgの`mtqg:example`と同じ形だが、`repo=`ではなく`dir=`、`ids=any`・`author=`はqsoku自身に該当する概念が無いので無し）。`TestDocExamplesAreMarkedAndMatch`が印の無い例・宙に浮いた印を検知し、英語版と`_ja`版で例の数・コマンド行が一致することを確かめる。

**gosecの扱い：** mtqg本体は`.golangci.yaml`で`_test.go$`全体からgosecを除外しているが、qsokuは既存方針（G204などをサイトごとに`//nolint:gosec`で理由付きにする）を踏襲し、**除外ルールは追加せず**、`e2e/examples_test.go`・`e2e/run_test.go`の該当箇所に理由付きの`//nolint:gosec`を個別に付けた。

## 2026-09-23　仕上げ：利用者向け文書と実測範囲の拡張

利用者向けの文書（看板README英日・`docs/tour/`・`docs/examples/`）を整えた。

**`docs/reference/`の実測の仕組みを広げた：** `e2e/examples_test.go`の`documents`（`docs/reference/`固定の4ファイル）を`documentPairs`（リポジトリルートからの相対パスの英日ペアの一覧）に一般化し、`README.md`/`README_ja.md`・`docs/tour/README.md`/`README_ja.md`も対象に加えた。仕組み自体（印・`qsoku`始まりの限定・パス置き換え）は変えていない。

**`docs/examples/`の実例：** `go`・`node`・`monorepo`の3つ。これらは`go`・`npm`のような、このリポジトリ自身が依存しないツールを呼ぶため、**実行はしない**——`e2e/run_test.go`の`TestDocsExamplesQsokufilesParse`が、各ディレクトリで`qsoku .list`が終了コード0であること（解析できること）だけを確かめる。

## 2026-09-23　正式公開の準備：GoReleaser・GitHubメタデータ

**GoReleaserの導入：** `.goreleaser.yaml`でlinux・darwinのamd64・arm64向けにクロスコンパイルする設定を書いた。**Windowsは対象外**——`testing.md`のCI検証OSの判断、`qsokufile`が`sh`実行前提でWindowsではGit Bash/WSLが要るという既存の判断と揃えた。`CGO_ENABLED=0`（distribution.mdの「純粋なGoにする」のまま）。アーカイブの中身はGoReleaserの既定glob（`LICENSE*`・`README*`）に任せている。

**バージョンの埋め込み：** `internal/cli/version.go`にmtqg本体の`internal/cli/version.go`と同じ形で`var version string`を足し、`-ldflags`で埋め込まれていればそれを最優先、空なら従来どおり`debug.ReadBuildInfo()`に落ちる、という順にした。`go install`側の挙動（`ReadBuildInfo`に頼る理由）は変えていない。

**CIの安全網：** `.github/workflows/ci.yml`に`goreleaser`ジョブを足し、通常のpush・PRのたびに`make goreleaser-check`（`goreleaser check`＋`--snapshot --skip=publish`のビルド）を実行するようにした。**タグを打つ前に設定の壊れを検知できる**——`.github/workflows/release.yml`（`v*`タグのpushだけで動く）が実際に公開するのはタグを打った後なので、それより先にCIで気づけるようにする狙い（e2e/docs-examplesと同じ「実測で確かめる」方針を配布設定にも適用した）。

**GitHubリポジトリのメタデータ：** `gh repo edit`でdescriptionとTopics（`note.md`の下書き一覧12個）を設定した。`ai-agents`は見送った（AIエージェント向けの機能・説明が実際にできてから足す）。

## 2026-09-23　看板READMEの作り直しと、リポジトリ自身のqsokufile

v0.1.0の正式リリース後、トップの`README.md`/`README_ja.md`を、ロゴ・バッジ・特徴一覧・デモ・折りたたみ式の実行例を持つ形に作り直した。デモには実際の録画（`vhs`。`docs/assets/demo.gif`）を使っている。

**リポジトリ自身に本物の`qsokufile`を新設した（ドッグフーディング）：** デモの録画用に、リポジトリ直下へ`build: (cd //; go build ./...)` / `test: (cd //; go test ./...)`という2行の`qsokufile`を追加した。**`make`を経由させない**——`qsoku`が`make`の薄いラッパーに見えると、「makeの真似をしない」というqsoku自身の立ち位置（`CLAUDE.md`）と矛盾して見えるため、Goの実コマンドを直接`//`経由で呼ぶ形にした。`internal/cli/`という2階層下から`qsoku test`を実行すると、`//`がリポジトリのルートを見つけて`go test ./...`が実際に走る——録画用の小道具ではなく、このリポジトリで今後も使う本物の近道として残した。

**READMEの「## Example」節を実践的な内容に変更：** 冒頭の実行例が、`hello: echo "hello, $1"`という`hello world`のおもちゃの例に、位置引数`$1`という初見には分かりにくい要素が混ざっていた。`docs/tour/`で使っている`tour`fixture（`build`・`test`という実プロジェクトらしい名前の近道、`$1`を使わない）に差し替えた。`test`の項目（`(cd //; echo "testing from $QSOKU_ROOT")`）は`//`を使っている。

## 2026-09-23　リポジトリ自身のqsokufileは、CIには組み込まない

v0.1.1リリース後、「このリポジトリ自身のqsokufileをCIに導入するか」を検討し、**しない**と決めた。技術的には難しくない（`qsoku`を先にビルドしてPATHに通し、`ci.yml`の該当ステップを`qsoku build`・`qsoku test`に置き換えるだけ）が、2つの理由で見送った。

- **循環のリスク：** CIの本筋（build/test）がqsoku自身を経由すると、qsoku自身の不具合（終了コードの素通し・`//`置き換えなど）が原因でCIが落ちたときに、「テスト対象が悪いのか、qsokuが悪いのか」の切り分けが難しくなる。qsokuはリリースしたばかりで実運用実績が薄く、この時点で自分のCIの本筋を自分に依存させるのは時期尚早
- **`qsokufile`が`build`・`test`の2エントリしか無い：** CIが実際に走らせているのは`vet`・`lint`・`unit`・race・e2e・shellcheck・trivy・`goreleaser-check`の7系統で、置き換えられるのはごく一部。全体を置き換えるには`Makefile`のターゲットを`qsokufile`にも書き写すことになり、二重管理になる

**`make`→`qsoku`への移行を試すなら、qsoku以外の、もっと枯れたプロジェクトで試すべき**という結論になった。qsoku自身は、実開発では引き続き`make`を使う（`CLAUDE.md`・`testing.md`のまま）。リポジトリ直下の`qsokufile`は、デモ・手元で試す用の本物の近道として残すが、CIには組み込まない。
