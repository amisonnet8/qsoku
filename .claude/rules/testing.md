# テスト方針

## 開発環境の前提

開発者の手元環境は**Linux（devcontainer）のみ**。実機のWindowsは無いので、Windows固有の挙動（`pwd -W`のパス形式、シグナルの伝わり方など）はCI（`windows-latest`）でしか確かめられない（2026-09-27に対応済み。以前は「Windowsは対象外」としていたが、方針を改めた）。

**CI（`.github/workflows/`）の対象OSはLinux・macOS・Windows**（`check`ジョブのみ。`race`はLinux・macOSのまま、`shellcheck`・`trivy`・`goreleaser`はLinuxのみのまま）。**windows-latestランナーには`make`が無い**（`Makefile`は`sh`前提）ため、`ci.yml`の`check (windows)`・`test (windows)`ステップは、`make check`・`make test`の中身（`go vet`・`golangci-lint run`・`go test`）を直接呼ぶ

| 対象 | 進め方 |
| :--- | :--- |
| `qsokufile`の解析・`//`の置き換え・管理用コマンド | Linux・Goだけで単体テストを書く。字句解析（引用符の状態を追う部分）は、URL・引用符の中・括弧の中・`;`の後を表で押さえる |
| シェル連携・補完（bash・zsh・fish・pwsh） | 4つとも devcontainer に入れてある（`postCreate.sh`。pwshは公式feature`powershell`から）。**本物のシェルで動かして確かめる**（スクリプトを模したものに置き換えない）。`e2e/shell_test.go`（`go test -tags e2e`、`make test`）が実施。手元に無いシェルは`shellPath`ヘルパーで`exec.LookPath`が失敗すればSkipする（例：Windowsではzsh・fishが無いのでSkip） |
| 本物のバイナリを直接実行（終了コード・シグナル・`//`置き換え・`QSOKU_CWD_FILE`） | `e2e/run_test.go`（`make test`）。`internal/cli`の単体テストは`Run`をプロセス内で呼ぶだけなので、ビルドした実バイナリでしか見えない部分をここで押さえる |
| `docs/reference/`・`README.md`・`docs/examples/README.md`の実行例 | 手で書かない。`$ qsoku ...`の直前に`<!-- qsoku:example dir=<fixture> ... -->`を置くと、`e2e/examples_test.go`が実際にビルドしたqsokuで実行し、`make docs-examples`が出力を文書へ書き込む（`docs/reference/README.md`参照）。対象文書は`e2e/examples_test.go`の`documentPairs`。比較だけなら`make test`が実施 |
| `docs/examples/`のqsokufile | 実行はしない（`go`・`npm`など、このリポジトリに無いツールを呼ぶため）。`e2e/run_test.go`の`TestDocsExamplesQsokufilesParse`が、各ディレクトリで`qsoku .list`が成功すること（解析できること）だけを確かめる |
| 依存の脆弱性・ライセンス | `trivy`（`trivy.yaml`） |
| シェルスクリプト | `shellcheck`（追跡中の`*.sh`・`*.bash`のみ） |
| bcuzを強制するフック（`.claude/hooks/require-bcuz.sh`） | `make hook`（`require-bcuz.test.sh`の判定の表。偽の`bcuz`を`PATH`に置くので、bcuzが無くても動く。`jq`が要る） |

## 実装後の動作確認

実装したら、ビルド確認に加えて実際の動作確認を行うこと。ロジック上正しそうに見えても、動かして初めて見つかる不具合はある。

- `make check`（vet・lint・単体テスト。整形の差異は`lint`が報告する）は作業の区切りで必ず通す。シェル連携・実行に関わる変更は`make test`（e2e）も通す
- 並行性・排他制御に関わる変更は`make race`も通す
- 依存を足したら`make trivy`（既知の脆弱性・ライセンス）を通す
- シェルスクリプトを足したら`make shellcheck`を通す（`git add`してから走らせる。追跡中のものだけが対象）
- `.claude/hooks/require-bcuz.sh`を変えたら`make hook`を通す（判定を変えたら、`require-bcuz.test.sh`の表も同じ変更の中で直す）
- ビルド確認に加えて、実際に`qsoku`コマンド（`make build`が作る`./qsoku`）を打って確かめること
- `qsokufile`のコマンドを実行するテストは、**本物の`sh`**で動かす（sh実行を模したものに置き換えない）

## Goのテスト

- **期待値（件数・ID・桁など）は、手で書かず、実際の出力を読んで正しいと確かめてから固定する。** 手で数えた件数や写し間違えたIDは、正しい実装の出力と食い違い、実装を疑う遠回りになる。テストが落ちたときは、実装と期待値のどちらが誤りかを先に確かめる
- **gitを使うテストは、開発者の設定の影響を受けないようにする。** `HOME`と`GIT_CONFIG_GLOBAL`を一時ディレクトリに向け、`user.name`・`user.email`はテストの中で決める（`t.Setenv`）
- **Goのソースに、見えない文字（ゼロ幅、U+2028・U+2029、文字の向きを変える制御文字、そのほかの制御文字）をそのまま書かない。** エスケープ（`\u2028`・`\x00`など）で書く。書き込みの途中で本物の文字に置き換わると、見た目では気づけない不具合になる

## 壊して確かめる（mutation-check）

通るだけのテストは、効いているとは限らない。実装を1か所ずつ壊して、テストが落ちるかを見る。手順・スクリプトは、mtqgのリポジトリ（https://github.com/amisonnet8/mtqg ）の`.claude/skills/mutation-check/`（`mutate.sh`が、変異を入れて、テストを走らせ、必ず元に戻す）を参考にする。生き残った変異は、握りつぶさずに読み、テストを足すか、同じ意味だと理由を書いて残す。

## mtqg固有の検証項目（qsokuに置き換えたもの）

- **常に`sh`で実行すること**：`()`がサブシェルとして扱われ、外の居場所を変えないこと。`cd`を含むコマンドが、成功でも失敗でも実行後の居場所を正しく持ち帰ること
- **`//`の置き換え**：単語の先頭で引用符の外にあるときだけ置き換わり、URL（`http://`）・引用符の中（シングル・ダブル）では置き換わらないこと
- **居場所の受け渡し**：`QSOKU_CWD_FILE`に書き出す1行が、コマンド自身の標準出力（`make build`のログなど）と混ざらないこと。qsoku自身のエラー（名前がない、`sh`が起動できない）では一時ファイルに何も書かれず、シェル関数側で`cd`が起きないこと
- **終了コード**：コマンドを実行した後は`sh`の終了コードをそのまま返すこと（qsoku自身の決まりで上書きしない）。コマンドを実行する前のqsoku自身のエラーは0/1/2（`docs/reference/cli.md`「Exit codes」）を表テストで押さえる
- **シェル補完**：登録されている名前を補完する。bash・zsh・fish・pwshで確認する（mtqgの`e2e/completion_test.go`のやり方が参考になる）。**fishは同じ候補にならない**：`.`で始まる候補（管理用コマンド）は、打っている語自体が`.`で始まるまでfishが隠すため（ドットファイルの扱いと同じ規則）。`qsoku <TAB>`は定義済みの名前だけ、`qsoku .<TAB>`で管理用コマンドが出ることを別々に確かめる（`docs/reference/cli.md`「Shell completion」に明記）。pwshは`Register-ArgumentCompleter`（bash・zshと同じく`.`の隠し規則は無い）を`TabExpansion2`で呼び出して確かめる

## Windows対応（2026-09-27）

- **qsokufileのコマンドはWindowsでも常に`sh`で実行する**（CLAUDE.md「実行は常にsh」を維持。Git for Windowsの`sh.exe`が`PATH`にある前提）。PowerShellは呼び出し元シェルとしてのみ対応する（`qsoku .shell pwsh`、`pwsh`＝PowerShell 7+を全OS共通で対応。Windows標準の旧PowerShell 5.1は検証しない）
- 手元（Linux devcontainer）ではWindows実機の検証はできない。**Windows固有の変更はクロスコンパイル（`GOOS=windows go build ./...`）で構文・型チェックを行い、実際の動作確認はCIの`windows-latest`で行う**——手元での確認とCIでの確認は別物と考える

### Windowsでの落とし穴（実機CIで反復修正して分かったこと）

- **パスの比較は、生の文字列ではなく`filepath.EvalSymlinks`（`os.SameFile`でもよい）を通してから行う。** 同じディレクトリでも、macOSの`/tmp`→`/private/tmp`のようなシンボリックリンクのエイリアスや、Windowsの8.3短縮名（`RUNNER~1`）と実際の長い名前（`runneradmin`）など、表記の違うパスが混在する
- **pwshは物理パス（シンボリックリンクを解決した実体）を返し、bash・zsh・fishは論理パス（`cd`した通りの文字列）を保つ。** 同じ操作でも、どのシェルで確かめるかによって期待値の作り方を変える必要がある
- **Windowsのbash・zshの`pwd`はMSYS形式（`/c/Users/...`）を返す。** これは`internal/run`が使う`pwd -W`（Windows形式`C:/Users/...`）とは別物なので、テストで比較する際は変換が要る（`e2e/shell_test.go`の`toMsysPath`）
- **ネイティブなWindowsパス（バックスラッシュ）を、引用符なしで`sh`のスクリプト文字列に直接埋め込まない。** `sh`は引用符の外のバックスラッシュを次の1文字へのエスケープとして読むため、`C:\Users\...`が`C:Users...`のように文字化けする
- **パスの区切り文字にまつわる変更は、`internal/*/*_test.go`も含めて横断的にgrepして洗い出す。** Windows以外ではバックスラッシュとフォワードスラッシュが一致してしまうため、`e2e/`だけ直して安心すると、Linux・macOSでは何度実行しても検知できない見落としが残る
- **pwshの出力には、DECCKMなどの端末制御シーケンス（`\x1b[?1h`等）や行末の`\r`が混ざることがある。** 標準出力がパイプでも起きる。`e2e/shell_test.go`の`cleanPwshOutput`で除去してからパースする
- **Git Bash（MSYS）は、ネイティブ（MSYS非対応）プログラムに渡す前に、`/`で始まる引数を書き換えることがある。** `docs/reference/`の実行例をGit Bash経由で動かすテストには`MSYS_NO_PATHCONV=1`が要る（`e2e/examples_test.go`）
- **WindowsにはPOSIXパーミッションビットが無い。** `os.Chmod`は読み取り専用属性の切り替えにしかならないため、パーミッションを検証するテストはWindowsで`t.Skip`にする
- **シグナル系のテストはWindowsで一旦Skip**：`kill -TERM $$`がGit Bashの`sh.exe`経由でGoの`syscall.WaitStatus`にどう見えるか、実機Windowsで未確認のため（`internal/run/run_test.go`・`e2e/run_test.go`）。確認・Skip解除はtodo `6fa5476e`

## Trivy・ShellCheckの運用

- `trivy fs`で、依存モジュールの既知の脆弱性（CVE）とライセンスを検査する（`make trivy`）
- 何を失敗とみなすか（深刻度、禁止するライセンスの種類）は`trivy.yaml`で固定する
- Trivyは脆弱性データベースの取得に通信が要る（devcontainerのBashサンドボックスを使う場合は下記「Bashサンドボックスの落とし穴」参照）
- シェルスクリプトは最小限にする。込み入った処理はGoで書く
- `make shellcheck`は`git ls-files '*.sh' '*.bash'`で追跡中のものを列挙してShellCheckにかける
- `go get`で依存を足したら、必ずTrivyを通す。依存は原則Go標準ライブラリと`golang.org/x/`だけ（`CLAUDE.md`）
- シェルスクリプトのコメント行を、行頭が小文字の`# shellcheck`で始まる文にしない（ShellCheckがインラインディレクティブとして誤解釈し、SC1072/SC1073でパースが止まる）
- **`postCreate.sh`の中で`. /etc/os-release`をsourceしない。** ShellCheckのSC1091（追えない形で指定されたファイル）が情報の水準で報告され、終了コードが0でなくなる。ディストリビューションの版は`lsb_release -rs`で取る

## `-race`の運用

- devcontainerは`CGO_ENABLED=0`。`-race`が要るときは`make race`（`CGO_ENABLED=1`にして走らせる）で分けている。CIの`race`ジョブはLinux・macOSのみ（Windowsはcgoにmingwが要るため対象外。`windows-latest`には標準でCコンパイラがない）

## Bashサンドボックスの落とし穴

> devcontainer内でClaude Codeが実行するBashコマンドは、既定でOSレベルのサンドボックス（Linux bubblewrap）にかかる。`.claude/settings.json`の`sandbox`でファイルシステムの書き込み先とネットワーク接続先を許可リストで絞っている（2026-09-30導入）。**設定の経緯・判断は人間が持つ**（`CLAUDE.md`「権限・自動化について」）。ここは、その設定の下で`make check`・`make test`・`make trivy`・`make shellcheck`・`make goreleaser-check`を実際に動かして見つかった落とし穴の記録（一部はmtqg本体の同種の記録を参考にした。2026-10-01）。

- **Trivyの脆弱性DBの取得先は`ghcr.io`ではなく`mirror.gcr.io`（`mirror.gcr.io/aquasec/trivy-db:2`）。** 名前から`ghcr.io`（GitHub Container Registry）を許可すればよいと思い込むと`make trivy`が`sandbox_violations`で失敗する。`trivy --debug`の`--db-repository`既定値で実際の取得先を確認できる
- **Trivyの脆弱性DBは`~/.cache/trivy`に書き込む。** `sandbox.filesystem.allowWrite`にこれが無いと、ネットワークを許可してもダウンロード後の書き込みで`read-only file system`になる（`make trivy`の初回失敗として実際に発生。`.devcontainer/postCreate.sh`が`settings.json`の`allowWrite`の全パスを事前に`mkdir -p`するので対策済み——サンドボックスの書き込み許可は、既存パスへのbindマウント式らしく、存在しないパスを新しく掘れないため）
- **`allowWrite`に書いたパスは、ディレクトリが既に存在しないと効かない。** 新しいコンテナには`~/.cache/golangci-lint`・`~/.cache/trivy`が無く、`~/.cache`自体は読み取り専用なので、`make lint`・`make check`は`failed to initialize build cache ... read-only file system`で**失敗**し、`make trivy`もDBを保存できず終了コード1になる。`postCreate.sh`が`.claude/settings.json`の`allowWrite`を`jq`で読み、`~/`と絶対パス（`/go`など。`/dev`・`/proc`・`/sys`と相対パスは除く。作れなければ警告だけ）を`mkdir -p`して先に作る（`allowWrite`を変えても`postCreate.sh`は直さなくてよい）。**すでにあるコンテナを直すには、サンドボックスの外で該当のディレクトリを`mkdir -p`する**
- **Linuxでは`bubblewrap`（`bwrap`）と`socat`が無いと、`sandbox.enabled: true`でも黙って無効のままになる。** エラーも出ないので、サンドボックスが効いているつもりで効いていない状態になる。`postCreate.sh`が、無いものだけ`apt`で入れる（ベースイメージが持っていれば何もしない）
- **golangci-lintは、Goより古いと標準ライブラリを読めず、`typecheck`で落ちる**（`could not load export data … export data version 5 is greater than maximum supported version 4`）。golangci-lintは、自分が抱えるx/toolsの版までのexport dataしか読めない。コンテナのGoが新しいパッチ・版に上がると、固定した版（ビルド済みの`install.sh`の版も）が突然落ちる。**同じ版をコンテナのGoで作り直しても直らない**（実測）。そのため、版を固定せず`go install …/golangci-lint/v2/cmd/golangci-lint@latest`で入れる（`.devcontainer/postCreate.sh`、CIの`version: latest`）。Goが上がった直後で、対応するlintがまだ無いときは落ちたままになる。`make lint`が、このエラーを見つけたら更新のしかたを案内する（`Makefile`）。直し方は、`go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest`を動かす（`proxy.golang.org`と`storage.googleapis.com`への接続が要る）。固定しない代わりに、新しいlintの指摘で赤くなることがある（そのときは設定か指摘を直す）
- **golangci-lintのキャッシュは`~/.cache/golangci-lint`に書き込む。** `filesystem.allowWrite`にこれが無いと、`make lint`・`make check`が`Failed to persist facts to cache ... read-only file system`の警告を出す。結果（`0 issues`）と終了コードには影響せず、キャッシュが効かなくなるだけ
- **`check.trivy.dev`への接続は、許可リストに無くても`make trivy`の結果・終了コードには影響しない。** Trivy自身のバージョン確認機能への接続で、`sandbox_violations`として拒否されるが、スキャン自体（脆弱性・ライセンスのレポート）は正常に完了する（mtqg本体の記録からの情報。qsoku側では`make trivy`がまだ最後まで通っていないため未検証）
- **許可リストに無いものに当たったら、自分で回り道を探し続けない。** `.claude/rules/working-with-human.md` 2章の形で、何をどこに足すかを短く頼む
- **`git fetch`・`git pull`と`gh pr create`には`github.com`・`api.github.com`への許可が要る**
- **`sudo`は、サンドボックスの中では承認しても失敗する。** サンドボックスのuidの付け替え（特権の無いユーザー名前空間）で、本来rootが持つファイルがnobodyの持ち物に見え、`sudo`自身の所有者の確認に落ちる。`sudo`が要るときは、人間のターミナルで動かすよう頼む（`.claude/rules/working-with-human.md` 2章）
- **`srwr init`は、書き換えるものがあるときは、サンドボックスの中では`.mcp.json`・`.claude/settings.json`を置き換えられない（書き換えが無ければ動く。設定はこのリポジトリに入っているので、ふだんは動かさなくてよい）**（一時ファイルを`rename`する先が書き込み保護で、`device or resource busy`。`.srwr/`の鍵だけが先にできる）。人間のターミナルで動かす（`.claude/rules/srwr-usage.md`「落とし穴」）。`srwr view`は`~/.cache/srwr`に書くので、同じく人間のターミナルで動かす
- **作業ディレクトリ直下に`.bashrc`・`.gitconfig`・`.claude/agents`などのダミーファイルが`git status`の未追跡ファイルとして現れることがある。** `ls -la`で見ると`/dev/null`相当のキャラクタデバイス（`crw-rw-rw-`）で、本来ホームディレクトリ側を指すはずのサンドボックスの書き込み保護パスが、作業ディレクトリ基準で解決されてしまったものと見られる。中身は空。`.gitignore`のドットエントリのホワイトリスト（`/.*`を無視して追跡中のものだけ`!`で戻す、2026-10-01）で`git status`には出なくなった。**新しく追跡するドットファイル・ドットディレクトリを足すときは、`.gitignore`に`!`行を足すこと**（足し忘れると`git add`が無視される）。qsoku・mtqgの実装のどちらにも起因しない、サンドボックス機構側の挙動

## GitHub Actions CIの落とし穴

- **`windows-latest`で`qsoku`を`go install`しても、内部で呼ぶ`sh`がそのままPATHに乗るとは限らない。** このリポジトリのCIは`qsoku`を使わない（`make`。Windowsでは中身を直接呼ぶ）が、qsokuの利用者向けの注意として残す
- **`/dev/stderr`等のUnix固有のパスはWindows（Git Bash）で壊れる。** `tee /dev/stderr`などを使わず、標準のリダイレクトだけで書く
- **Windowsのcheckoutで改行がCRLFになると、`gofmt -l`が全ファイルを未整形と誤検知する。** ルートの`.gitattributes`（`* text=auto eol=lf`）で防ぐ
- **`uses: owner/repo@TAG`はタグ名と厳密に一致しないと失敗する。** `v`の有無を見落としやすい。書く前に実際のタグ名を確かめる
- **Windowsには実行ビットの概念がない。** `os.Chmod`後に実行ビットを確かめるテストは、`runtime.GOOS != "windows"`でガードする
- **実行権限の無いファイルで「起動できない（126）」を作るテストは、Windowsでは「見つからない（127）」になる。** 拡張子の無い名前（`script`・`build.sh`）は`PATHEXT`に合わず、コマンドとして探されないため。Windowsでは`.exe`の名前で、中身が実行形式でないファイルにすると、起動に失敗して126になる
- **出力に出るパスと比べるテストは、一時ディレクトリを`filepath.EvalSymlinks`で解決してから比べる。** 同じディレクトリでも、macOSは`/var`→`/private/var`のシンボリックリンク、Windowsは8.3の短縮名（`RUNNER~1`）と長い名前が食い違う。`t.TempDir()`の値のまま比べると、macOS・WindowsのCIだけで落ちる
- **Windowsでは、開いたままのファイルを置き換え（`rename`）られないことがあり、開いたハンドルのあるディレクトリは消せない。** 一時ファイル経由の置き換えや、失敗したときの後片付け（`os.RemoveAll`）の前に、先に`Close`する。Windowsのランナーでだけ落ちる
