# Entry points for building and checking qsoku (.claude/rules/testing.md).
# The recipes assume a POSIX shell (Git Bash on Windows).

.PHONY: build fmt vet lint unit check test docs-examples race trivy shellcheck hook goreleaser-check

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

# Also reports files that the formatters would change. If golangci-lint is older
# than the Go in this container it cannot read the standard library (typecheck
# errors "export data version N is greater than ..."); say what to do about it.
lint:
	@out=$$(golangci-lint run 2>&1); rc=$$?; printf '%s\n' "$$out"; \
	if [ $$rc -ne 0 ] && printf '%s' "$$out" | grep -q 'export data version'; then \
	  echo "hint: this golangci-lint is older than the Go in this container and cannot read its standard library. Update it: go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest (.claude/rules/testing.md)" >&2; \
	fi; \
	exit $$rc

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

# The decision table of the hook that forces commands through bcuz (needs jq).
hook:
	bash .claude/hooks/require-bcuz.test.sh

# Validates .goreleaser.yaml and does a real build of every release target
# (distribution.md), without publishing anything (--skip=publish) or
# needing a tag (--snapshot). release.yml runs the real thing, only when a
# human pushes a version tag.
goreleaser-check:
	goreleaser check
	goreleaser release --snapshot --clean --skip=publish
