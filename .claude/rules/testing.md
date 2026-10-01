# テスト方針

## 開発環境の前提

開発者の手元環境は**Linux（devcontainer）のみ**。実機のWindowsは無いので、Windows固有の挙動（`pwd -W`のパス形式、シグナルの伝わり方など）はCI（`windows-latest`）でしか確かめられない（2026-09-27に対応済み。以前は「Windowsは対象外」としていたが、方針を改めた。経緯は`docs/design/history.md`参照）。

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

## 実装後の動作確認

- ビルド確認に加えて、実際に`qsoku`コマンドを打って確かめること
- `qsokufile`のコマンドを実行するテストは、**本物の`sh`**で動かす（sh実行を模したものに置き換えない）

## mtqg固有の検証項目（qsokuに置き換えたもの）

- **常に`sh`で実行すること**：`()`がサブシェルとして扱われ、外の居場所を変えないこと。`cd`を含むコマンドが、成功でも失敗でも実行後の居場所を正しく持ち帰ること
- **`//`の置き換え**：単語の先頭で引用符の外にあるときだけ置き換わり、URL（`http://`）・引用符の中（シングル・ダブル）では置き換わらないこと
- **居場所の受け渡し**：`QSOKU_CWD_FILE`に書き出す1行が、コマンド自身の標準出力（`make build`のログなど）と混ざらないこと。qsoku自身のエラー（名前がない、`sh`が起動できない）では一時ファイルに何も書かれず、シェル関数側で`cd`が起きないこと
- **終了コード**：コマンドを実行した後は`sh`の終了コードをそのまま返すこと（qsoku自身の決まりで上書きしない）。コマンドを実行する前のqsoku自身のエラーは0/1/2（`docs/reference/cli.md`「Exit codes」）を表テストで押さえる
- **シェル補完**：登録されている名前を補完する。bash・zsh・fish・pwshで確認する（mtqgの`e2e/completion_test.go`のやり方が参考になる）。**fishは同じ候補にならない**：`.`で始まる候補（管理用コマンド）は、打っている語自体が`.`で始まるまでfishが隠すため（ドットファイルの扱いと同じ規則）。`qsoku <TAB>`は定義済みの名前だけ、`qsoku .<TAB>`で管理用コマンドが出ることを別々に確かめる（`docs/reference/cli.md`「Shell completion」に明記）。pwshは`Register-ArgumentCompleter`（bash・zshと同じく`.`の隠し規則は無い）を`TabExpansion2`で呼び出して確かめる

## Windows対応（2026-09-27）

- **qsokufileのコマンドはWindowsでも常に`sh`で実行する**（CLAUDE.md「実行は常にsh」を維持。Git for Windowsの`sh.exe`が`PATH`にある前提）。PowerShellは呼び出し元シェルとしてのみ対応する（`qsoku .shell pwsh`、`pwsh`＝PowerShell 7+を全OS共通で対応。Windows標準の旧PowerShell 5.1は検証しない）
- 手元（Linux devcontainer）ではWindows実機の検証はできない。**Windows固有の変更はクロスコンパイル（`GOOS=windows go build ./...`）で構文・型チェックを行い、実際の動作確認はCIの`windows-latest`で行う**——手元での確認とCIでの確認は別物と考える

### Windowsでの落とし穴（実機CIで反復修正して分かったこと。経緯は`docs/design/history.md`）

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

- `go get`で依存を足したら、必ずTrivyを通す。依存は原則Go標準ライブラリと`golang.org/x/`だけ（`CLAUDE.md`）
- シェルスクリプトのコメント行を、行頭が小文字の`# shellcheck`で始まる文にしない（ShellCheckがインラインディレクティブとして誤解釈し、SC1072/SC1073でパースが止まる）

## `-race`の運用

- devcontainerは`CGO_ENABLED=0`。`-race`が要るときは`make race`（`CGO_ENABLED=1`にして走らせる）で分けている。CIの`race`ジョブはLinux・macOSのみ（Windowsはcgoにmingwが要るため対象外）

## Bashサンドボックスの落とし穴

> devcontainer内でClaude Codeが実行するBashコマンドは、既定でOSレベルのサンドボックス（Linux bubblewrap）にかかる。`.claude/settings.json`の`sandbox`でファイルシステムの書き込み先とネットワーク接続先を許可リストで絞っている（2026-09-30導入）。**設定の経緯・判断は人間が持つ**（`CLAUDE.md`「権限・自動化について」）。ここは、その設定の下で`make check`・`make test`・`make trivy`・`make shellcheck`・`make goreleaser-check`を実際に動かして見つかった落とし穴の記録（一部はmtqg本体の同種の記録を参考にした。2026-10-01）。

- **Trivyの脆弱性DBの取得先は`ghcr.io`ではなく`mirror.gcr.io`（`mirror.gcr.io/aquasec/trivy-db:2`）。** 名前から`ghcr.io`（GitHub Container Registry）を許可すればよいと思い込むと`make trivy`が`sandbox_violations`で失敗する。`trivy --debug`の`--db-repository`既定値で実際の取得先を確認できる
- **Trivyの脆弱性DBは`~/.cache/trivy`に書き込む。** `sandbox.filesystem.allowWrite`にこれが無いと、ネットワークを許可してもダウンロード後の書き込みで`read-only file system`になる（`make trivy`の初回失敗として実際に発生。`.devcontainer/postCreate.sh`が`settings.json`の`allowWrite`の全パスを事前に`mkdir -p`するので対策済み——サンドボックスの書き込み許可は、既存パスへのbindマウント式らしく、存在しないパスを新しく掘れないため）
- **golangci-lintのキャッシュは`~/.cache/golangci-lint`に書き込む。** `allowWrite`にこれが無いと書き込みが`read-only file system`になる。qsokuが固定している`golangci-lint` v2.13.2は、この書き込み失敗を警告なしで黙って無視する（`0 issues`・終了コード0のまま変わらない）——結果は壊れないが、キャッシュが効かず`make lint`・`make check`のたびに毎回フルスキャンになる
- **`check.trivy.dev`への接続は、許可リストに無くても`make trivy`の結果・終了コードには影響しない。** Trivy自身のバージョン確認機能への接続で、`sandbox_violations`として拒否されるが、スキャン自体（脆弱性・ライセンスのレポート）は正常に完了する（mtqg本体の記録からの情報。qsoku側では`make trivy`がまだ最後まで通っていないため未検証）
- **作業ディレクトリ直下に`.bashrc`・`.gitconfig`・`.claude/agents`などのダミーファイルが`git status`の未追跡ファイルとして現れることがある。** `ls -la`で見ると`/dev/null`相当のキャラクタデバイス（`crw-rw-rw-`）で、本来ホームディレクトリ側を指すはずのサンドボックスの書き込み保護パスが、作業ディレクトリ基準で解決されてしまったものと見られる。中身は空。`.gitignore`のドットエントリのホワイトリスト（`/.*`を無視して追跡中のものだけ`!`で戻す、2026-10-01）で`git status`には出なくなった。**新しく追跡するドットファイル・ドットディレクトリを足すときは、`.gitignore`に`!`行を足すこと**（足し忘れると`git add`が無視される）。qsoku・mtqgの実装のどちらにも起因しない、サンドボックス機構側の挙動
