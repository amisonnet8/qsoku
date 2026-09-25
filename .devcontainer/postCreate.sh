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
# zsh, fish:      Two of the three shells qsoku's completion and shell integration
#                 support (.claude/rules/testing.md). PowerShell is not a target
#                 (qsokufile commands always run under sh).
sudo apt-get update
sudo apt-get install -y make wget gnupg lsb-release gcc jq shellcheck zsh fish

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

# golangci-lint: lint (.golangci.yaml). The official install script puts the
# binary into GOPATH/bin. The version is pinned so that lint results do not
# change when the container is rebuilt.
wget -qO - https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b "$(go env GOPATH)/bin" v2.13.2

go install golang.org/x/tools/gopls@latest
go install golang.org/x/tools/cmd/goimports@latest

# goreleaser: builds and validates the release (.goreleaser.yaml,
# distribution.md, "make goreleaser-check"). Same version .github/workflows/
# pins, so a local check matches CI.
go install github.com/goreleaser/goreleaser/v2@v2.18.2

# mtqg: this repository's own process-recording tool (CLAUDE.md,
# .claude/rules/mtqg.md). Also wires up the Claude Code hook and MCP server
# (.claude/settings.json, .mcp.json; "mtqg init --agent claude-code").
# Pinned to a commit, not a version tag: mtqg has not tagged a release yet.
go install github.com/amisonnet8/mtqg/cmd/mtqg@56b331c65e1406bd17b6d357eb04341973d36579
