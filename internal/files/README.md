# Files Package

The `files` package provides robust file system scanning capabilities for the copyright CLI tool.

## Features

- **Directory Tree Walking**: Efficiently traverses directory structures using `filepath.WalkDir`
- **Gitignore Support**: Automatically respects `.gitignore` files (root and nested)
- **Pattern Matching**: Supports glob patterns for include/exclude rules
- **Binary Detection**: Automatically identifies binary files using multiple heuristics
- **Context Support**: Respects context cancellation for long-running operations
- **Error Handling**: Gracefully handles permission errors and other issues
- **Hidden Files**: Optionally skip hidden files and directories
- **Symlink Handling**: Configurable symlink following behavior

## Installation

The package is already integrated into the project. The required dependency is:

```bash
go get github.com/sabhiram/go-gitignore
```

## Usage

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/charlesgreen/copyright/internal/files"
)

func main() {
    // Create scanner with default options
    scanner, err := files.NewScanner(".", nil)
    if err != nil {
        log.Fatal(err)
    }

    // Scan the directory
    result, err := scanner.Scan(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    // Process results
    for _, file := range result.Files {
        if !file.IsBinary {
            fmt.Printf("Processing: %s\n", file.RelativePath)
        }
    }
}
```

### Custom Options

```go
options := &files.ScanOptions{
    ExcludePatterns: []string{
        "**/test/**",       // Exclude test directories
        "**/*.tmp",         // Exclude temporary files
        "**/node_modules/**", // Exclude node_modules
    },
    IncludePatterns: []string{
        "**/*.go",          // Only include Go files
        "**/*.js",          // and JavaScript files
    },
    ExcludedExtensions: []string{
        ".exe",             // Skip executables
        ".dll",
        ".so",
    },
    SkipHidden: true,       // Skip hidden files/directories
    FollowSymlinks: false,  // Don't follow symlinks
}

scanner, err := files.NewScanner("/path/to/repo", options)
```

### With Context and Timeout

```go
import (
    "context"
    "time"
)

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

result, err := scanner.Scan(ctx)
if err != nil {
    log.Fatal(err)
}
```

### Error Handling

```go
result, err := scanner.Scan(context.Background())
if err != nil {
    log.Fatal(err)
}

// Check for warnings during scan
if len(result.Errors) > 0 {
    fmt.Printf("Scan completed with %d warnings:\n", len(result.Errors))
    for _, err := range result.Errors {
        fmt.Printf("  - %v\n", err)
    }
}

// See what was skipped
fmt.Printf("Skipped %d files:\n", len(result.SkippedFiles))
for _, path := range result.SkippedFiles {
    fmt.Printf("  - %s\n", path)
}
```

## Pattern Matching

The package supports glob patterns for flexible file matching:

### Patterns

- `*` - Matches any characters except `/`
- `**` - Matches any characters including `/` (recursive)
- `?` - Matches a single character
- `[abc]` - Character classes

### Examples

```go
"*.go"                  // All .go files in current directory
"**/*.go"               // All .go files recursively
"test/**"               // Everything under test directories
"**/node_modules/**"    // All node_modules directories
"src/**/*.js"           // All .js files under src
"*.{go,js,py}"          // Multiple extensions (if supported by matcher)
```

## Binary Detection

Files are classified as binary using multiple methods:

1. **File Extension Check**: Quick check against known binary extensions (`.exe`, `.png`, `.zip`, etc.)
2. **Magic Number Signatures**: Detects common binary formats (ELF, PE, ZIP, PNG, JPEG, PDF, etc.)
3. **Null Byte Detection**: Checks first 8KB for null bytes
4. **Non-Printable Character Ratio**: If >30% of characters are non-printable, classified as binary

### Usage

```go
isBinary, err := files.IsBinary("/path/to/file")
if err != nil {
    log.Fatal(err)
}

if isBinary {
    fmt.Println("File is binary")
} else {
    fmt.Println("File is text")
}
```

## Gitignore Support

The scanner automatically:

- Loads `.gitignore` from repository root
- Discovers nested `.gitignore` files during traversal
- Applies patterns from all relevant `.gitignore` files
- Supports custom patterns via `ScanOptions.ExcludePatterns`
- Always ignores `.git` directories

## Types

### FileInfo

```go
type FileInfo struct {
    Path         string  // Absolute file path
    RelativePath string  // Relative to repository root
    Extension    string  // Including dot (e.g., ".go")
    Filename     string  // Just the filename
    Size         int64   // File size in bytes
    IsBinary     bool    // Whether file is binary
}
```

### ScanResult

```go
type ScanResult struct {
    Files        []FileInfo  // Successfully scanned files
    SkippedFiles []string    // Files that were skipped
    Errors       []error     // Errors encountered during scan
}
```

### ScanOptions

```go
type ScanOptions struct {
    ExcludePatterns    []string  // Patterns to exclude
    IncludePatterns    []string  // Patterns to include (if set, only these)
    ExcludedExtensions []string  // Extensions to skip
    SkipHidden         bool      // Skip hidden files/directories
    FollowSymlinks     bool      // Follow symbolic links
}
```

## Default Behavior

When creating a scanner with `nil` options, the following defaults are used:

```go
&ScanOptions{
    SkipHidden:     true,
    FollowSymlinks: false,
    ExcludePatterns: []string{
        "**/.git/**",
        "**/node_modules/**",
        "**/vendor/**",
    },
}
```

## Performance Considerations

- The scanner uses `filepath.WalkDir` which is efficient for large directory trees
- Binary detection reads only the first 8KB of files
- Directory pruning happens early to avoid unnecessary traversal
- Context cancellation is checked at each directory level

## Error Handling

The scanner handles errors gracefully:

- **Permission Denied**: Logged as errors but scanning continues
- **Invalid Patterns**: Return errors immediately
- **Context Cancellation**: Stops scanning immediately
- **File Access Errors**: Logged per-file but don't stop the scan

## Testing

Run the examples:

```bash
go test -v -run Example ./internal/files
```

## Integration

This package is designed to integrate with other components:

- **Config Package**: Load scan options from configuration
- **Language Package**: Determine file types for language-specific processing
- **Processor Package**: Process discovered files for copyright headers
- **Reporter Package**: Generate reports on scan results
