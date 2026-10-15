# CLAUDE.md

このファイルは qsoku プロジェクトで作業する際、Claude Code が常に踏まえるべき前提を示す。

## qsokuとは

**リポジトリ固有のコマンドの近道**を`qsokufile`に書いておき、リポジトリのどこにいても`qsoku <名前>`で呼べるCLIツール。makeの`make build`から着想を得たが、依存関係やターゲットの仕組みは持たず、「短い名前 → コマンド」の対応表だけに徹する。

- `alias`との違い：`alias`は自分の環境全体に効き、リポジトリごとに変えられず、共有もできない。qsokuは`qsokufile`をコミットすれば、チーム全員が同じ近道を使える
- makeとの違い：makeはビルドの道具で、タブの決まりや`.PHONY`などの作法があり、`cd`のような移動は苦手。qsokuは近道の一覧に割り切り、`cd`も効く
- AIエージェントにも効く：`qsokufile`がそのリポジトリのコマンドの一覧になり、人間とAIが同じ名前で呼べる
- 独自の決まりは3つだけ（`qsokufile`の探し方、`//`、居場所の持ち帰り）。残りはすべて`sh`の知識で書ける

**位置づけ：** mtqg（`github.com/amisonnet8/mtqg`）の段階2のサンプルPJとして始まったが、本当に便利な道具になりそうなので、記録は使い捨てにせず残す方向（詳細は`docs/design/note.md`）。

## 参照すべきファイル

作業を始める前に、以下を確認すること。

- **`docs/design/`** — 設計判断と、その理由の記録（日本語）。`note.md`が出発点（設計メモの原文）、`README.md`に位置づけと索引がある。経緯（日付ごとの流れ）は置かず、mtqgに記録する（`.claude/rules/documentation.md`）
- **`docs/reference/`** — データ形式・コマンドの仕様（英語が正、日本語版`*_ja.md`と2本立て）。実行例は`e2e/examples_test.go`（`make docs-examples`）が実測して書き込む。新しい仕様を書くときも、実装に入る前にまずここへ書くこと
- **`.claude/rules/bcuz-usage.md`** — bcuzの使い方。**Bashで実行するコマンドは、必ず`bcuz "理由" コマンド…`の形で打つ**（フックが、そうでないものを断る）。理由は「なぜ」を書く。**コマンドを打つ前に読む**
- **`.claude/rules/mtqg-usage.md`** — このリポジトリでのmtqgの使い方（記録者、種類の使い分け、こまめに記録する。下の「記録の使い分け」の詳細）
- **`.claude/rules/srwr-usage.md`** — srwrの使い方（**厳格モード**）。**`Edit`・`Write`は禁止されており、ファイルの変更は`look`／`edit`／`new`（理由`why`つき）で行う。作業のまとまりごとに`session`で題つきのテープを始め、コミットの前にテープを閉じる。ファイルを変える前に読む**
- **`.claude/rules/qsoku-usage.md`** — このリポジトリの開発は`make`（`Makefile`が正。名前の一覧、`qsokufile`は見本）。**ビルド・テスト・lintは`go`などを直接打たず、`make <名前>`で実行する**
- **`.claude/rules/directory-structure.md`** — ディレクトリ構成（今あるものと、まだ無いもの）と、ファイルの置き場所の判断基準。**新しいファイルを追加するときに読む**
- **`.claude/rules/testing.md`** — テスト方針（開発環境の前提、動作確認の粒度、Trivy、ShellCheck、`-race`、Windows、Bashサンドボックス・CIの落とし穴）
- **`.claude/rules/datetime.md`** — 日時の扱い（データは UTC で持ち、表示だけ利用者のタイムゾーンに直す。瞬間で比べる、時計を差し替える）。**日時を保存・表示・比較するコードを書くときに読む**
- **`.claude/rules/working-with-human.md`** — 人間とのやり取り（頼み方の3種類、小さなつまずきの頼み方、確認の自動化、見え方を変えるときの了承）。**人間に何かを頼む前に読む**
- **`.claude/rules/naming.md`** — 命名規則（未記入。決まったら書く）
- **`.claude/rules/distribution.md`** — 配布方法（`go install`、GoReleaser、リリースの流れ）
- **`.claude/rules/documentation.md`** — ドキュメントの書き方（`docs/design/`に書く経緯の線引き。そのほかは未記入で、決まったら書く）

これらのルールファイルは、実装中に得た細かい気づき・教訓を Claude Code 自身が育てていくものである。新しく気づいたルール・踏んだ落とし穴があれば、該当するファイルに追記すること。どのファイルにも当てはまらない新しい種類の気づきであれば、新しいルールファイルを作ってよい（作ったらこの一覧にも足すこと）。

## 記録の使い分け（mtqgの運用）

このリポジトリは、**mtqgを使ってCLIだけで作業の過程を記録しながら**開発する（`.claude/rules/mtqg-usage.md`に詳細）。

- **進捗管理の文書（PLAN.md相当）は作らない。** 現在地・次にやること・保留事項は**mtqgのtodoとmemo**に入れる
- **作業を始める前に、必ず`mtqg context`を読むこと。** 今の状態はここから拾う
- 作業中に生まれた気づき→`m add`（memo）、やること→`t add`（todo）、迷って人間に判断を仰いだこと→`q add`（question）、不具合→`b add`（bug）、言葉の合意→`g add`（glossary）、チームの決まり事→`r add`（rule）
- **少なくとも一度は、ブランチか別の端末で並行して記録すること。** mtqgのunionマージと「並行した状態変更」の表示を実地で試すため（サンプルPJとしての目的の一つ）
- 記録者は環境変数`MTQG_AUTHOR_KIND=ai`・`MTQG_AUTHOR_NAME=claude-code`を設定して使う
- `.mtqg/`には決まり事（rule）を9件登録済み——`mtqg context`で確認できる
- **記録は日本語で書く**（2026-09-27、mtqg rule。詳細は`.claude/rules/mtqg-usage.md`）
- **mtqgの使い勝手について気づいたことは、その場でmemoかbugに残すこと。** コマンドが分かりにくい、表示が足りない、データ形式が足りない、逆にAIが記録を忘れる・`context`を読まない・引数を間違える、といった気づきが、mtqg本体への報告（「CLIで直すもの」「MCP・フックで解くもの」の2分類）の材料になる

## 開発の進め方

- **設計判断を先に決め、`docs/reference/`（仕様）を実装より先に書く。** 実装が先行して仕様が後追いになると、どちらが正なのか分からなくなる
- **Bashのコマンドは、すべて`bcuz "理由" コマンド…`で実行する**（`bcuz "lintの確認" make lint`）。理由は「何をするか」でなく「なぜするか」。`bcuz`で始まらない・`&&`や`;`でつなぐコマンドは、`.claude/settings.json`の`PreToolUse`フック（`.claude/hooks/require-bcuz.sh`）が実行前に断る。**`git commit`だけは`bcuz`を通さない**（通すとコミットの後に記録が増えるため。理由はコミットメッセージに書く）。実行の記録は`.bcuz/history.jsonl`に残り（コミットする）、`bcuz -H`で見られる。**トークンなどの秘密を、コマンドの引数・`env`の指定・理由に書かない**（記録に残り、コミットされる）（`.claude/rules/bcuz-usage.md`）
- **ファイルの変更は、srwrの道具（`look`で見て`edit`、新規は`new`）で行い、すべてに理由（`why`）を添える。** 組み込みの`Edit`・`Write`は`.claude/settings.json`で禁止されている（厳格モード）。Bashでの書き換え（`sed -i`・リダイレクト）もしない
- **作業のまとまり（todo・bug・人間の依頼・作業の種類の変わり目）ごとに、最初に`session`で題つきのテープを始める。** 1つの作業＝1本のテープ、コミットも同じ区切り
- **コミットの前に、テープを閉じる**（`srwr tapes check`のあと`srwr tapes new`）。書き込み中の生のテープは`.gitignore`で対象外で、閉じて圧縮した`.srwr/tapes/*.tape.jsonl.gz`だけがgitに入る（`.claude/rules/srwr-usage.md`）
- 実装後は`.claude/rules/testing.md`に従い、ビルド確認に加えて実際の動作確認（`bcuz "…" make check`、必要に応じて`make test`・`make race`・`make trivy`・`make shellcheck`）を行うこと。判断に迷ったら実行する側に倒すこと
- **依存ライブラリは、Goの標準ライブラリと`golang.org/x/`だけを使う。** それ以外は「どうしても必要な場合」の例外とし、理由を説明して確認を取り、mtqgに記録すること。追加したらTrivyでライセンスと脆弱性を確認する
- **看板としてのREADME（`README.md`・`README_ja.md`）は実装完了後に作った。** v0.1.0のリリース後、人間の判断で冒頭の「開発初期」の注意書きは外した

## 守ること（設計の芯）

実装の細部より優先する。これを崩す変更は、実装する前に必ず相談すること。

- **`qsokufile`の解釈はバイナリの1か所に集める。** シェルごとにコマンドの文法を焼き直さない（`//`の置き換え、`qsokufile`の探索、居場所の持ち帰りは、対応する各シェルのスクリプトではなくGoのコードで行う）
- **実行は常に`sh`。** どのシェルから呼ばれても、`qsokufile`のコマンドは`sh`で実行する（利用者のシェルによって意味が変わらないように）
- **独自の決まりを増やさない。** `//`の逃がし方は引用符（新しい記号を作らない）。管理用コマンドは`.`で始まるものに閉じる
- **`qsokufile`は確認なしで実行する。** 利用者が`qsoku`と打ったときだけ動くため（direnvのような自動フックとは違う）
- **makeの真似をしない。** 依存関係・ターゲット・`.PHONY`のような仕組みは持たない

## ドキュメントの言語

- **`CLAUDE.md`・`.claude/rules/`・`docs/design/`は日本語のみ。** 英語版は作らない
  - 例外：本ファイル末尾の`mtqg init --agent claude-code`が挿入した1行（mtqgの配線コマンドが冪等性の判定に使う定型文で、書き換えると再実行のたびに重複追記されうるため、英語のまま残す。2026-09-25）
- **`docs/reference/`は英語版（`*.md`）と日本語版（`*_ja.md`）の2本立て。** 英語版が正。同じ変更の中で両方直す
- **コード・コメント・コミットメッセージ・CLI出力は英語。** 対話は日本語で行うが、これらは英語で書くこと

## gitの扱い（Claude Codeの作業として）

- **pushは人間が行う。** Claude Codeはコミットまでで止めること。`git push`は`.claude/settings.json`で拒否（deny）されている
- 作業の区切りでこまめにコミットしてよい（消してしまったときに戻せるように）
- `git reset --hard`・`git clean`は、コミットしていない作業を消すので、使う前に必ず理由を説明すること
- **コミットの前に`git config user.email`／`user.name`を確かめる。** 正しい値は`65074207+amisonnet8@users.noreply.github.com`／`amisonnet8`（GitHubのプライバシー用noreplyアドレス）。違う値なら、`.git/config`の`[user]`（ローカルの上書き）を疑う。ローカルに設定が無いと、コンテナのglobalの設定（noreplyでないことがある）が使われる

## 権限・自動化について

`.claude/settings.json`（人間が管理する）により、次の設定になっている。

- **サンドボックス化**（2026-09-30、mtqg本体の`settings.json`を参考に変更）：ネットワークは許可したドメイン（`proxy.golang.org`・`sum.golang.org`・`mirror.gcr.io`・`github.com`・`api.github.com`）だけに制限され、書き込みは`~/.cache/go-build`・`/go`・`~/.cache/trivy`・`~/.cache/golangci-lint`だけに限られる（`~/.ssh`・`~/.aws`・`~/.config/gh`は読み取りも不可）。このおかげで`curl`・`wget`・`go get`は個別の確認なしで動く。このコンテナのGOPATHは`/go`（`go env GOPATH`）
- 拒否（deny）：`git push`、`gh pr merge`、`gh release create`（タグを打つのは人間、という`distribution.md`の方針と揃えた）、`Edit`・`Write`・`MultiEdit`・`NotebookEdit`（srwrの厳格モード、下）
- **`bcuz`を通さないBashコマンドの拒否**：`PreToolUse`フック（`.claude/hooks/require-bcuz.sh`）が、`bcuz`で始まらないコマンド（`git commit`を除く。逆に`bcuz`を通した`git commit`を断る）と、後ろに別のコマンドをつなぐもの（`&&`・`||`・`;`・`&`・改行・コマンド置換）を、実行前に断る。直し方は断りのメッセージと`.claude/rules/bcuz-usage.md`に書いてある
- 実行前に確認（ask）：`git reset --hard`、`git clean`、`sudo`、`gh pr create`

確認を求められた場合、無理に実行しようとせず、指示を仰ぐこと。サンドボックスの制限（許可ドメイン外のネットワーク・許可外のファイル書き込み）に引っかかった場合は、`.claude/rules/working-with-human.md` 2章の形（原因と、`.claude/settings.json` のどこに何を足すか）で**短く頼む**（許可を広げるかは人間が判断する）。`.claude/settings.json`を変更したいときは、変更案を差分で示すにとどめること（人間が貼るだけで済む形にする）。落とし穴は`.claude/rules/testing.md`「Bashサンドボックスの落とし穴」参照

**srwrの厳格モード**：`permissions.deny`に`Edit`・`Write`・`MultiEdit`・`NotebookEdit`があり、`.mcp.json`にsrwrのMCPサーバー、`PostToolUse`に`srwr hook`（読み取りと外部変更の記録）がある。テープ（`.srwr/tapes/`）はgitで共有する（鍵`.srwr/key`は入れない）。許可を広げる（緩いモードにする）のは人間の判断。詳細は`.claude/rules/srwr-usage.md`。

ビルドの自動フック（`PostToolUse`）が設定されている。srwrの`edit`・`new`で`.go`・`go.mod`・`go.sum`を変えると`.claude/hooks/build.sh`が`make build`を実行し（直下の`./qsoku`も作り直す）、失敗すると理由がClaude Codeに返る。複数ファイルにまたがる編集の途中は、まだ書いていないファイルを参照して一時的に失敗することがある。

**既知の制約（2026-09-30、サンドボックス化に伴って判明）：**

- `sudo`は承認されても実際には失敗する。サンドボックスのuid remapping（アンプリビレッジドなユーザー名前空間）で、本来root所有のファイルがnobody所有に見えてしまい、sudo自身の所有者チェックに落ちるため。`ask`に残しているのは、将来サンドボックスの実装が変わった場合に備えてで、今のところ`sudo`が要る操作（追加のapt installなど）は代わりに人間に依頼すること
- サンドボックスの書き込み許可（`sandbox.filesystem.allowWrite`）は、パスが既に存在する前提の仕組みらしく、存在しないパスを新規に掘ることはできない。新しくキャッシュディレクトリ等が要る場合は、`allowWrite`に足すだけでよい（`.devcontainer/postCreate.sh`が`allowWrite`の全パスをサンドボックス外で`mkdir -p`する。コンテナの再作成後に有効）

## まだ無いもの

初期構成にわざと入れていない。作業の中で必要になった時点で作る。各項目が何のためのものかは`.claude/rules/directory-structure.md`（ツリー全体と`docs/`の4つの違い）を見ること。

- `go.sum`（依存ができたら足す。依存が無いためまだ無い）

## ルール・スキルの提案

作業を進める中で、以下に気づいたら、都度こちらから提案すること（提案するだけで、勝手に作成・適用はしない。判断はこちらが行う）。

- **新しいルールにした方がよさそうな知見**: 同じ種類の判断や落とし穴に複数回遭遇した、既存の`.claude/rules/`のどれにも当てはまらない、といった場合。どのファイルに追記すべきか、あるいは新規ファイルが必要かも合わせて提案する
- **Skill化した方が効率的そうな作業**: 「同じ手順を3回目繰り返している」など、パターン化できそうな一連の作業に気づいた場合。まだ実装が薄い段階では無理に提案する必要はない。パターンが実際に繰り返されてから提案すること

提案は気づいたタイミングで随時行ってよく、まとめて報告するために貯めておく必要はない。

This repository records its development with mtqg: run `mtqg context` at the start of a session, and see `.mtqg/SCHEMA.md` for the data format.
