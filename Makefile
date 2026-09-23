# Entry points for building and checking qsoku (.claude/rules/testing.md).
# The recipes assume a POSIX shell (Git Bash on Windows).

.PHONY: build fmt vet lint unit check test docs-examples race trivy shellcheck goreleaser-check

# Compile every package first: a package that cmd/qsoku does not import yet
# would otherwise be skipped, and so would its build errors.
build:
	go build ./...
	go build ./cmd/qsoku

# Rewrite files with the formatters enabled in .golangci.yaml.
fmt:
	golangci-lint fmt

vet:
	go vet ./...
	go vet -tags e2e ./e2e/...

lint:
	golangci-lint run

unit:
	go test ./...

check: vet lint unit

# End-to-end tests: the real qsoku binary against real shells (e2e/, built
# with the tag e2e, so make check does not run them).
test:
	go test -tags e2e -count=1 ./e2e/...

# Writes what qsoku actually prints into the examples of docs/reference/
# (docs/reference/README.md). Nobody writes an example's output by hand.
docs-examples:
	go test -tags e2e -count=1 -run '^TestDocExamples$$' ./e2e/... -update

# -race needs cgo, which the container turns off (.claude/rules/testing.md).
race:
	CGO_ENABLED=1 go test -race -count=1 ./...

trivy:
	trivy fs --config trivy.yaml .

shellcheck:
	git ls-files '*.sh' '*.bash' | xargs -r shellcheck

# Validates .goreleaser.yaml and does a real build of every release target
# (distribution.md), without publishing anything (--skip=publish) or
# needing a tag (--snapshot). release.yml runs the real thing, only when a
# human pushes a version tag.
goreleaser-check:
	goreleaser check
	goreleaser release --snapshot --clean --skip=publish
