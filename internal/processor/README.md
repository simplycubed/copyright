# Processor Package

The `processor` package provides the core orchestration system for the copyright CLI tool. It coordinates file scanning, copyright detection, header generation, and file modification in a unified workflow.

## Features

- **Two Operation Modes**: Check mode (read-only) and Fix mode (adds headers)
- **Smart File Processing**: Automatically detects file types and applies appropriate comment styles
- **Special Line Preservation**: Preserves shebangs, encoding declarations, and XML declarations
- **Context Support**: Full support for cancellation and timeouts
- **Dry-Run Mode**: Preview changes without modifying files
- **Backup Support**: Optionally create .bak files before modifications
- **Comprehensive Reporting**: Detailed results for each file and aggregate statistics
- **LICENSE Management**: Optional automatic LICENSE file creation

## Architecture

The processor integrates several components:

```
┌─────────────┐
│  Processor  │
└──────┬──────┘
       │
       ├──> Scanner (files package)
       │    └─> Discovers files to process
       │
       ├──> Registry (language package)
       │    └─> Determines file language
       │
       ├──> Checker (copyright package)
       │    └─> Detects copyright headers
       │
       ├──> Generator (copyright package)
       │    └─> Generates copyright text
       │
       ├──> Writer (processor package)
       │    └─> Writes headers to files
       │
       └──> Manager (license package)
            └─> Manages LICENSE files
```

## Usage

### Basic Check Mode

Check files for copyright headers without modifying them:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/simplycubed/copyright/internal/config"
    "github.com/simplycubed/copyright/internal/processor"
)

func main() {
    // Load or create configuration
    cfg := config.NewDefault()
    cfg.Copyright.Company = "SimplyCubed, GK"
    cfg.Copyright.SPDXIdentifier = "MIT"

    // Create processor
    proc, err := processor.NewProcessor(cfg)
    if err != nil {
        log.Fatal(err)
    }

    // Process files in check mode
    ctx := context.Background()
    result, err := proc.Process(ctx, "./src", processor.ModeCheck)
    if err != nil {
        log.Fatal(err)
    }

    // Display results
    fmt.Println(result.Summary())
    os.Exit(result.ExitCode())
}
```

### Fix Mode

Automatically add copyright headers to files missing them:

```go
opts := &processor.ProcessorOptions{
    DryRun: false,
    Backup: true, // Create .bak files
}

result, err := proc.ProcessWithOptions(ctx, "./src", processor.ModeFix, opts)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Fixed %d files\n", result.FilesFixed)
```

### Dry-Run Mode

Preview changes without modifying files:

```go
opts := &processor.ProcessorOptions{
    DryRun: true,
    Backup: false,
}

result, err := proc.ProcessWithOptions(ctx, "./src", processor.ModeFix, opts)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Would fix %d files\n", result.FilesFixed)
```

### With Context Timeout

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

result, err := proc.Process(ctx, "./src", processor.ModeCheck)
if err == context.DeadlineExceeded {
    fmt.Println("Processing timed out")
}
```

## Processing Workflow

For each file discovered:

1. **Skip Binary Files**: Binary files are automatically skipped
2. **Language Detection**: File type is determined from extension/name
3. **Comment Support Check**: Files without comment support are skipped
4. **Copyright Detection**: File is scanned for copyright keywords
5. **Mode-Specific Action**:
   - **Check Mode**: Record whether copyright is present
   - **Fix Mode**: Add copyright header if missing
6. **Special Line Handling**: Preserve shebangs and encoding declarations
7. **Header Insertion**: Insert at correct position with proper formatting

## Result Types

### FileStatus

```go
const (
    StatusHasCopyright      // File has copyright header
    StatusMissingCopyright  // File missing copyright (check mode)
    StatusFixed             // Copyright added (fix mode)
    StatusSkipped           // File skipped (unsupported/binary)
    StatusFailed            // Processing failed
)
```

### FileResult

Result for a single file:

```go
type FileResult struct {
    FilePath     string      // Absolute path
    RelativePath string      // Relative to scan root
    Status       FileStatus  // Processing status
    Message      string      // Human-readable message
    Error        error       // Error if failed
}
```

### ProcessResult

Aggregate results:

```go
type ProcessResult struct {
    TotalFiles            int
    FilesWithCopyright    int
    FilesMissingCopyright int
    FilesFixed            int
    FilesSkipped          int
    FilesFailed           int
    FileResults           []FileResult
    LicenseCreated        bool
    LicenseStatus         string
}
```

Methods:
- `ExitCode() int`: Returns 0 if all files have copyright, 1 otherwise
- `Summary() string`: Human-readable summary

## Special Line Preservation

The processor correctly handles files with special leading lines:

### Shell Scripts (Shebang)

```bash
#!/bin/bash
# Copyright 2025 SimplyCubed, GK
# SPDX-License-Identifier: MIT

echo "Hello"
```

### Python (Encoding + Shebang)

```python
#!/usr/bin/env python3
# -*- coding: utf-8 -*-
# Copyright 2025 SimplyCubed, GK
# SPDX-License-Identifier: MIT

def main():
    pass
```

### Ruby (Encoding)

```ruby
# encoding: utf-8
# Copyright 2025 SimplyCubed, GK
# SPDX-License-Identifier: MIT

puts "Hello"
```

## Configuration Integration

The processor uses configuration from the `config` package:

```yaml
copyright:
  template: "Copyright {{.Year}} {{.Company}}\nSPDX-License-Identifier: {{.SPDX}}"
  year: "auto"
  company: "SimplyCubed, GK"
  spdx_identifier: "MIT"

scan:
  exclude_patterns:
    - "**/node_modules/**"
    - "**/vendor/**"
  skip_hidden: true

detection:
  lines_to_scan: 20
  keywords:
    - "copyright"
    - "spdx-license-identifier"
```

## Error Handling

### Fatal Errors

Return immediately with error:
- Invalid root path
- Failed to create scanner
- Context cancelled during setup

### Non-Fatal Errors

Recorded in FileResults, processing continues:
- Failed to read individual file
- Failed to write header to file
- Binary content detection errors

### Example Error Handling

```go
result, err := proc.Process(ctx, "./src", processor.ModeCheck)
if err != nil {
    // Fatal error - couldn't complete processing
    log.Fatal(err)
}

// Check for individual file failures
for _, fr := range result.FileResults {
    if fr.Status == processor.StatusFailed {
        log.Printf("Failed to process %s: %v", fr.RelativePath, fr.Error)
    }
}

// Exit with appropriate code
os.Exit(result.ExitCode())
```

## Testing

The package includes comprehensive tests:

```bash
# Run all tests
go test ./internal/processor/...

# Run with coverage
go test -cover ./internal/processor/...

# Run examples
go test -v -run Example ./internal/processor/
```

## Performance Considerations

- **Parallel Processing**: Files are processed sequentially (can be made parallel if needed)
- **Memory Usage**: Only one file is held in memory at a time
- **I/O Optimization**: Minimal disk I/O - single read, single write per file
- **Context Checking**: Regular context checks for responsive cancellation

## Supported File Types

The processor supports all file types handled by the language registry:

- Go, JavaScript/TypeScript, Java, C#, Rust, C/C++
- Python, Ruby, Shell scripts
- YAML, HTML, CSS, SQL, Protocol Buffers
- Terraform, Dockerfiles, Makefiles

Files without comment support (JSON, Markdown) are automatically skipped.

## Future Enhancements

Potential improvements:
- Parallel file processing for better performance
- Progress callbacks for long-running operations
- Incremental processing (only changed files)
- Custom header templates per file type
- Header update mode (refresh existing headers)
