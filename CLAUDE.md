# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Test Commands

```bash
make build          # Build the copyright binary
make test           # Run all tests with verbose output
make test-unit      # Run unit tests only (excludes integration tests)
make test-race      # Run tests with race condition detection
make lint           # Run golangci-lint (auto-installs if missing)
make fmt            # Format code with gofmt
make vet            # Run go vet
make check          # Run fmt, vet, lint, and test (local development)
make ci             # Run fmt-check, vet, lint, and test (CI environment)
make coverage       # Generate HTML coverage report
```

Run a single test:
```bash
go test -v -run TestName ./internal/package/...
```

## Architecture

This is a Go CLI tool for managing copyright headers in source files. It uses Cobra for CLI commands and Viper for configuration.

### Core Packages

- **`main.go`** - CLI entry point with Cobra commands (check, fix, init, version)
- **`internal/processor`** - Orchestrates the workflow: scans files, checks for copyright, adds headers when missing
- **`internal/config`** - Loads `.copyright.yaml` configuration using Viper, defines all config structs
- **`internal/files`** - File system scanning with gitignore support, binary detection, pattern matching
- **`internal/language`** - Registry of 20+ language handlers; each knows its extensions, comment style, and how to format headers
- **`internal/copyright`** - Checker (keyword detection in first N lines) and HeaderGenerator (template rendering)
- **`internal/license`** - LICENSE file detection and generation from templates

### Data Flow

1. `Processor` creates a `Scanner` which walks the directory tree respecting gitignore and exclude patterns
2. For each file, `Registry.GetLanguage()` returns the appropriate language handler based on extension/filename
3. `Checker.Check()` scans the first N lines for copyright keywords
4. In fix mode, `HeaderGenerator.Generate()` renders the template, `Language.FormatHeader()` wraps it in comments, and `Writer.WriteHeader()` inserts it (preserving shebangs/encoding declarations)

### Language Handler Pattern

Each language implements the `Language` interface (`internal/language/language.go`):
- `FormatHeader()` wraps copyright text in language-appropriate comments
- `PreserveLines()` identifies lines that must stay before copyright (shebangs, encoding declarations)
- `InsertPosition()` returns byte offset where copyright should be inserted

New languages are registered in `registerAllLanguages()` in `internal/language/registry.go`.

### Configuration

The tool reads `.copyright.yaml` from the current directory. Key sections:
- `copyright` - Template, year, company, SPDX identifier
- `scan` - Exclude patterns, skip hidden files, follow symlinks
- `detection` - Lines to scan, keywords, require_all mode
