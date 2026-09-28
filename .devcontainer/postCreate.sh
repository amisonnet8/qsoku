#!/usr/bin/env bash
set -euo pipefail

# make:           Build and test entry points (once a Makefile exists).
# wget, gnupg,
# lsb-release:    Adding the Trivy and GitHub CLI apt repositories below.
# gcc:            make race (CGO_ENABLED=1 go test -race) needs a C compiler.
#                 The container itself runs with CGO_ENABLED=0 (.claude/rules/distribution.md).
# jq:             Inspecting devcontainer.json / settings.json and --json output while debugging.
# ShellCheck:     Static analysis of tracked *.sh and *.bash files.
#                 Comment lines must not start with the lowercase directive word,
#                 or ShellCheck parses them as directives (SC1072/SC1073).
# zsh, fish:      Two of the shells qsoku's completion and shell integration
#                 support (.claude/rules/testing.md). PowerShell (pwsh) is
#                 the third; it is installed separately below (Microsoft's
#                 own apt repository, not Debian's).
sudo apt-get update
sudo apt-get install -y make wget gnupg lsb-release gcc jq shellcheck zsh fish apt-transport-https software-properties-common

# Trivy: known vulnerabilities (CVE) and license compatibility of dependencies.
# Installed from the official apt repository.
wget -qO - https://aquasecurity.github.io/trivy-repo/deb/public.key | gpg --dearmor | sudo tee /usr/share/keyrings/trivy.gpg >/dev/null
echo "deb [signed-by=/usr/share/keyrings/trivy.gpg] https://aquasecurity.github.io/trivy-repo/deb $(lsb_release -sc) main" | sudo tee /etc/apt/sources.list.d/trivy.list >/dev/null
sudo apt-get update
sudo apt-get install -y trivy

# gh: GitHub CLI, for checking issues, pull requests and Actions runs.
# Installed from the official apt repository (same pattern as Trivy).
sudo mkdir -p -m 755 /etc/apt/keyrings
wget -qO - https://cli.github.com/packages/githubcli-archive-keyring.gpg | sudo tee /etc/apt/keyrings/githubcli-archive-keyring.gpg >/dev/null
sudo chmod go+r /etc/apt/keyrings/githubcli-archive-keyring.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list >/dev/null
sudo apt-get update
sudo apt-get install -y gh

# PowerShell (pwsh): the third shell qsoku's shell integration and
# completion support (docs/reference/cli.md "Shell integration"). It is
# only ever the *caller's* shell -- qsokufile commands themselves always
# run under sh, on every OS (CLAUDE.md "実行は常にsh"). Installed from
# Microsoft's own apt repository (same pattern as Trivy and gh above).
wget -q "https://packages.microsoft.com/config/debian/12/packages-microsoft-prod.deb" -O /tmp/packages-microsoft-prod.deb
sudo dpkg -i /tmp/packages-microsoft-prod.deb
rm /tmp/packages-microsoft-prod.deb
sudo apt-get update
sudo apt-get install -y powershell

# golangci-lint: lint (.golangci.yaml). The official install script puts the
# binary into GOPATH/bin. The version is pinned so that lint results do not
# change when the container is rebuilt.
wget -qO - https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b "$(go env GOPATH)/bin" v2.13.2

go install golang.org/x/tools/gopls@latest
go install golang.org/x/tools/cmd/goimports@latest

# goreleaser: builds and validates the release (.goreleaser.yaml,
# distribution.md, "make goreleaser-check"). @latest here only ever resolves
# within v2 (the major version is part of the Go import path), matching
# .github/workflows/'s own "~> v2" constraint on goreleaser-action -- both
# stay current automatically without drifting to a v3 that might break
# .goreleaser.yaml's v2-schema config.
go install github.com/goreleaser/goreleaser/v2@latest

# mtqg: this repository's own process-recording tool (CLAUDE.md,
# .claude/rules/mtqg.md). Also wires up the Claude Code hook and MCP server
# (.claude/settings.json, .mcp.json; "mtqg init --agent claude-code"). mtqg
# has tagged v1 (2026-09-29); the VS Code extension (amisonnet8.mtqg, listed
# above) is installed by the devcontainer's extensions list, not here.
go install github.com/amisonnet8/mtqg/cmd/mtqg@latest

mkdir -p ~/.local/share/bash-completion/completions
mtqg completion bash >~/.local/share/bash-completion/completions/mtqg

# qsoku: this repository's own binary, built from the local source (not a
# published version) so that it always matches the code in this checkout.
# Same pattern mtqg's own devcontainer uses for itself. postCreate runs once
# at container creation, so this does not track edits made afterward --
# .claude/hooks/build.sh keeps the separate ./qsoku at the repository root
# (used for manual testing, the demo GIF, and so on) up to date on every Go
# edit instead.
go install ./cmd/qsoku

# Wire up qsoku's own shell integration and completion for this repository's
# qsokufile (the "sample" one at the repository root -- CLAUDE.md, directory-
# structure.md; real development still goes through make). This is for
# dogfooding and checking that the integration itself actually works, not a
# replacement for make.
# shellcheck disable=SC2016 # single-quoted on purpose: written literally so it expands at bash startup, not now
grep -qF 'qsoku .shell bash' ~/.bashrc 2>/dev/null || echo 'eval "$(qsoku .shell bash)"' >>~/.bashrc
