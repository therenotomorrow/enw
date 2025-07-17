set dotenv-load := false
set shell := ["sh", "-cu"]

BIN := justfile_directory() / ".bin"

[private]
default:
    @just --list

# ---- golangci-lint

GOLANGCI_LINT_VERSION := 'v2.2.2'
GOLANGCI_LINT_PATH := BIN / 'golangci-lint'
GOLANGCI_LINT := GOLANGCI_LINT_PATH + '@' + GOLANGCI_LINT_VERSION

[private]
install-golangci-lint:
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b {{ BIN }} {{ GOLANGCI_LINT_VERSION }}
    mv {{ GOLANGCI_LINT_PATH }} {{ GOLANGCI_LINT }}

[doc('Run static analysis using `golangci-lint` to detect code issues')]
lint:
    @if test ! -e {{ GOLANGCI_LINT }}; then just install-golangci-lint; fi
    {{ GOLANGCI_LINT }} run ./...

# ---- fieldaligment

FIELDALIGNMENT_VERSION := 'v0.35.0'
FIELDALIGNMENT_PATH := BIN / 'fieldalignment'
FIELDALIGNMENT := FIELDALIGNMENT_PATH + '@' + FIELDALIGNMENT_VERSION

[private]
install-fieldaligment:
    GOBIN={{ BIN }} go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@{{ FIELDALIGNMENT_VERSION }}
    mv {{ FIELDALIGNMENT_PATH }} {{ FIELDALIGNMENT }}

[doc('Reorder struct fields using `fieldalignment` to improve memory layout')]
align:
    @if test ! -e {{ FIELDALIGNMENT }}; then just install-fieldaligment; fi
    {{ FIELDALIGNMENT }} --fix ./...

# ---- godoc

GODOC_VERSION := 'v0.35.0'
GODOC_PATH := BIN / 'godoc'
GODOC := GODOC_PATH + '@' + GODOC_VERSION

[private]
install-godoc:
    GOBIN={{ BIN }} go install golang.org/x/tools/cmd/godoc@{{ GODOC_VERSION }}
    mv {{ GODOC_PATH }} {{ GODOC }}

[doc('Run documentation server using `godoc`')]
docs port='8000':
    @if test ! -e {{ GODOC }}; then just install-godoc; fi
    @echo http://127.0.0.1:{{ port }}/pkg/github.com/therenotomorrow/ex/
    {{ GODOC }} -http=:{{ port }}

# ---- testing

[private]
smoke:
    go test ./...

[private]
cover:
    go test -count 1 -race -coverprofile=coverage.out
    go tool cover -func coverage.out

# ---- shortcuts

[doc('Run all code quality tools')]
code: align lint

[doc('Run tests by type: `smoke` for quick checks, `cover` for detailed analysis')]
test type='smoke':
    just {{ type }}
