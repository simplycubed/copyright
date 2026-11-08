# Copyright

A Golang CLI tool for automated copyright header management across multi-language codebases.

## Features

- Check source files for missing or non-compliant copyright headers
- Automatically add copyright headers to files without modification risk
- Support for 20+ programming languages and configuration file formats (see [Supported Languages](#supported-languages))
- Smart preservation of shebangs, encoding declarations, and other first-line directives
- Configurable templates with variable substitution (year, company, author, SPDX)
- LICENSE file detection and generation for common open-source licenses
- Output in human-readable text or JSON format
- Dry-run mode to preview changes before applying

## Installation

Requires Go 1.27.3.

```bash
go install github.com/simplycubed/copyright@latest
```

To uninstall:

```bash
rm $(go env GOPATH)/bin/copyright
```

Or build from source:

```bash
git clone https://github.com/simplycubed/copyright.git
cd copyright
make build
```

## Usage

```bash
copyright <command> [options] [path]

Commands:
  check    Verify files have copyright headers
  fix      Add missing copyright headers to files
  init     Create .copyright.yaml configuration file

Options:
  -v, --verbose    Verbose output (show all processed files)
  --dry-run        Preview changes without modifying files (fix only)
  --backup         Create .bak files before modifying (fix only)
  --json           Output results as JSON

Examples:
  copyright check
  copyright check ./src
  copyright fix --dry-run
  copyright fix --backup
  copyright check -v --json > results.json
```

## Quick Start

Initialize configuration in your project:

```bash
cd /path/to/your/project
copyright init
```

Edit `.copyright.yaml` to customize your copyright header.

Check for files missing headers:

```bash
copyright check
```

Preview what changes would be made:

```bash
copyright fix --dry-run
```

Apply copyright headers:

```bash
copyright fix
```

## Configuration

The `.copyright.yaml` file controls header generation and scanning behavior:

```yaml
copyright:
  template: 'Copyright {{.Year}} {{.Company}}. All rights reserved.'
  year: auto # "auto" uses current year, or specify "2024" or "2020-2024"
  company: 'SimplyCubed, GK'
  author: '' # Optional author field
  spdx_identifier: 'Apache-2.0' # Optional SPDX license identifier

license:
  type: 'apache-2.0' # Supported: mit, apache-2.0, bsd-3-clause, gpl-3.0, mpl-2.0, custom
  auto_create: false # Automatically create LICENSE file if missing

scan:
  exclude_patterns:
    - 'vendor/**'
    - 'node_modules/**'
    - '*.generated.*'
    - 'dist/**'
    - 'build/**'
  skip_hidden: true # Skip hidden files and directories
  follow_symlinks: false # Don't follow symbolic links

detection:
  lines_to_scan: 20 # Number of lines to scan for existing copyright
  keywords: # Keywords indicating copyright presence
    - 'copyright'
    - '©'
    - 'spdx-license-identifier'
    - '(c)'
  require_all: false # If true, all keywords must be present

output:
  format: text # Output format: text or json
  verbose: false # Show detailed processing information
  quiet: false # Suppress non-error output
```

### Template Variables

Available variables for the `copyright.template` field:

- `{{.Year}}` - Year(s) from the `year` configuration
- `{{.Company}}` - Company name from the `company` configuration
- `{{.Author}}` - Author name from the `author` configuration (optional)
- `{{.SPDX}}` and `{{.SPDXIdentifier}}` - SPDX license identifier from `spdx_identifier`

Example templates:

```yaml
# Simple format
template: "Copyright {{.Year}} {{.Company}}"

# With SPDX identifier
template: "Copyright {{.Year}} {{.Company}}\nSPDX-License-Identifier: {{.SPDXIdentifier}}"

# Detailed format
template: "Copyright (c) {{.Year}} {{.Company}} <contact@example.com>\nLicensed under the Apache License, Version 2.0"
```

## Supported Languages

| Language   | File Extensions                                          | Comment Style |
| ---------- | -------------------------------------------------------- | ------------- |
| Go         | `.go`                                                    | `//`          |
| Python     | `.py`, `.pyw`                                            | `#`           |
| JavaScript | `.js`, `.mjs`, `.cjs`, `.jsx`, `.mts`, `.cts`            | `//`, `/* */` |
| TypeScript | `.ts`, `.tsx`                                            | `//`, `/* */` |
| Java       | `.java`                                                  | `//`, `/* */` |
| C#         | `.cs`                                                    | `//`, `/* */` |
| Rust       | `.rs`                                                    | `//`, `/* */` |
| C/C++      | `.c`, `.h`, `.cpp`, `.hpp`, `.cc`, `.cxx`, `.hh`, `.hxx` | `//`, `/* */` |
| Ruby       | `.rb`, `.rake`, `.gemspec`, `Rakefile`, `Gemfile`        | `#`           |
| Shell      | `.sh`, `.bash`, `.zsh`, `.fish`, `.ksh`                  | `#`           |
| SQL        | `.sql`                                                   | `--`          |
| HTML/XML   | `.html`, `.htm`, `.xml`, `.xhtml`, `.svg`                | `<!-- -->`    |
| CSS        | `.css`, `.scss`, `.sass`, `.less`                        | `/* */`       |
| YAML       | `.yaml`, `.yml`                                          | `#`           |
| TOML       | `.toml`                                                  | `#`           |
| JSON       | `.json`, `.jsonc`                                        | N/A           |
| Markdown   | `.md`, `.markdown`                                       | `<!-- -->`    |
| Protobuf   | `.proto`                                                 | `//`          |
| Terraform  | `.tf`, `.tfvars`, `.hcl`                                 | `#`, `//`     |
| Dockerfile | `Dockerfile`, `Containerfile`                            | `#`           |
| Makefile   | `Makefile`, `makefile`, `GNUmakefile`, `.mk`             | `#`           |

## How It Works

The tool operates in two phases:

1. **Scanning**: Walks the directory tree identifying source code files based on extension matching, respecting exclusion patterns and hidden file settings

2. **Processing**:
   - **Check mode**: Scans the first N lines (configurable) for copyright keywords and reports files without valid headers
   - **Fix mode**: For files without headers, generates language-appropriate comments and inserts them while preserving special first-line directives (shebangs, encoding
     declarations, XML declarations)

Copyright detection uses keyword matching against configurable patterns. The tool does NOT validate copyright content accuracy, only presence/absence.

## Example Output

### Text Output (Check)

```text
Checking files for copyright headers...

✗ src/main.go - Missing copyright header
✗ src/utils/helper.go - Missing copyright header
✓ src/config/config.go - Has copyright header

Summary: 2 files missing headers (1 file OK)
```

### Text Output (Fix)

```text
Adding copyright headers...

✓ src/main.go - Copyright header added
✓ src/utils/helper.go - Copyright header added
⊘ src/config/config.go - Already has header (skipped)

Summary: 2 files updated (1 file skipped)
```

### JSON Output

```json
{
  "files_checked": 3,
  "files_missing_headers": 2,
  "files_with_headers": 1,
  "results": [
    {
      "path": "src/main.go",
      "has_header": false,
      "language": "go"
    },
    {
      "path": "src/utils/helper.go",
      "has_header": false,
      "language": "go"
    },
    {
      "path": "src/config/config.go",
      "has_header": true,
      "language": "go"
    }
  ]
}
```

## Exit Codes

| Code | Description                                                                        |
| ---- | ---------------------------------------------------------------------------------- |
| 0    | All files have copyright headers (check) or operation completed successfully (fix) |
| 1    | Some files missing headers (check) or operation failed with errors (fix)           |

Use exit codes for CI/CD integration:

```bash
# Fail build if any files are missing copyright headers
copyright check || exit 1
```

## Limitations

- Detects existing headers via keyword matching only (doesn't validate header content)
- JSON files receive copyright in a special metadata comment field (not standard JSON)
- Binary files are automatically skipped but may require explicit exclusion patterns
- Symlinks are not followed by default (configurable)
- Multi-line copyright blocks must contain at least one configured keyword in the scanned line range

## Development

```bash
# Run tests
make test

# Run linter
make lint

# Run all checks (test + lint)
make check

# Build binary
make build

# Clean build artifacts
make clean
```

## License

Apache License 2.0 - see [LICENSE](LICENSE) file for details.
