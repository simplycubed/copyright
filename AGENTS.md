# AGENTS.md: `github.com/simplycubed/copyright`

- `AGENTS.md` is the repo-local source of truth for agent instructions, and the contract for this
  tool's design decisions. There is no `CLAUDE.md` or other shim: agents read this file
  directly. Where this file states a decision, implement that decision rather than re-deriving an
  alternative. If a decision turns out to be wrong, change this file in the same PR that changes
  the code, with the reason.
- Usage and configuration live in [`README.md`](README.md). This file holds the rules and decisions.

A Go CLI that checks source files for a copyright header and inserts one when it is missing. The
module path is `github.com/simplycubed/copyright`. The copyright holder is SimplyCubed, GK. Commands
are Cobra; configuration is Viper reading `.copyright.yaml`.

## The one rule that shapes everything

**The tool detects whether a header is present and can add one. It does not decide whether the
notice is legally correct.**

Everything that looks like restraint follows from this:

- Detection is keyword matching in the first N lines. It does not parse the notice, compare it to
  the template, or check that the year, company, or SPDX identifier are right.
- Fix mode inserts the configured template. It does not rewrite a header that is already present.
- The holder, year, and SPDX identifier come from configuration. The binary does not assume
  SimplyCubed, GK for projects it scans.

## Decisions

1. **Four commands.** `check` reports, `fix` writes missing headers, `init` creates
   `.copyright.yaml` only when that file is absent, and `version` prints the version. `fix` accepts
   `--dry-run` and `--backup`. A missing config file falls back to defaults.
2. **`internal/processor` owns the workflow.** It scans, detects, generates, and writes. It does
   not grow language-specific comment rules of its own.
3. **One language handler per language, registered in one place.** Each handler implements
   `Language` in `internal/language/language.go`: `FormatHeader`, `PreserveLines`, and
   `InsertPosition`. New languages are added in `registerAllLanguages()` in
   `internal/language/registry.go`.
4. **Lines that are not the program stay above the header.** Shebangs, encoding declarations, and
   XML declarations are preserved. The header is inserted at the byte offset
   `InsertPosition` returns.
5. **Configuration has three sections that change behavior.** `copyright` (template, year, company,
   author, SPDX identifier), `scan` (exclude patterns, skip hidden files, follow symlinks), and
   `detection` (lines to scan, keywords, `require_all`).
6. **LICENSE generation is separate from header insertion.** `internal/license` detects and renders
   LICENSE files from templates. A source-file header and a repository LICENSE file are not the
   same write.

## Gate

`make ci` = fmt-check → vet → lint → test. CI runs that, then `make build`. Locally, `make check`
is the same sequence with `make fmt` instead of `fmt-check`.

**Gate-first, every feature:** no feature work begins until its gate exists; extend the gate first,
then build the feature against it. A gate nobody has watched fire is a false clear, so every new
check lands with a negative control.

```bash
make build          # build the copyright binary
make test           # all tests, verbose
make test-unit      # unit tests only (excludes integration tests)
make test-race      # tests with the race detector
make lint           # golangci-lint (installs the binary if it is missing)
make fmt            # gofmt
make vet            # go vet
make check          # fmt, vet, lint, and test
make ci             # fmt-check, vet, lint, and test
make coverage       # HTML coverage report
```

Run one test:

```bash
go test -v -run TestName ./internal/package/...
```

## Conventions

- **No attribution in commits or PRs:** never a `Co-Authored-By` trailer, a "Generated with…" line,
  or any tool-attribution text.
- Format with `gofmt -s` (`make fmt`). Do not hand-format around it.
- Keep the module path `github.com/simplycubed/copyright`. The notice in `LICENSE` names
  SimplyCubed, GK.
- Sample header text in tests and package docs is fixture data. Change both sides of an assertion
  together.
