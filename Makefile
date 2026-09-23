# Entry points for building and checking qsoku (.claude/rules/testing.md).
# The recipes assume a POSIX shell (Git Bash on Windows).

.PHONY: build fmt vet lint unit check test race trivy shellcheck

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

lint:
	golangci-lint run

unit:
	go test ./...

check: vet lint unit

# End-to-end tests: the real qsoku binary against a real sh (e2e/, built
# with the tag e2e, so make check does not run them). e2e/ does not exist
# yet (added in Step 8 of the implementation plan; .mtqg todo fd80a62a28).
test:
	@if [ -d e2e ]; then \
		go test -tags e2e -count=1 ./e2e/...; \
	else \
		echo "e2e/ does not exist yet; nothing to run"; \
	fi

# -race needs cgo, which the container turns off (.claude/rules/testing.md).
race:
	CGO_ENABLED=1 go test -race -count=1 ./...

trivy:
	trivy fs --config trivy.yaml .

shellcheck:
	git ls-files '*.sh' '*.bash' | xargs -r shellcheck
