// Package files provides file system scanning capabilities for the copyright CLI tool.
//
// The package implements a robust file scanner that:
//   - Walks directory trees efficiently
//   - Respects .gitignore files (root and nested)
//   - Supports custom exclude/include patterns with glob matching
//   - Detects binary files automatically
//   - Handles symlinks and hidden files
//   - Supports context cancellation for long-running scans
//
// Example usage:
//
//	// Create scanner with default options
//	scanner, err := files.NewScanner("/path/to/repo", nil)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Or with custom options
//	options := &files.ScanOptions{
//	    ExcludePatterns: []string{"**/test/**", "**/*.tmp"},
//	    IncludePatterns: []string{"**/*.go", "**/*.js"},
//	    ExcludedExtensions: []string{".exe", ".dll"},
//	    SkipHidden: true,
//	    FollowSymlinks: false,
//	}
//	scanner, err := files.NewScanner("/path/to/repo", options)
//
//	// Scan with context for cancellation
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
//	defer cancel()
//
//	result, err := scanner.Scan(ctx)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Process results
//	for _, file := range result.Files {
//	    if !file.IsBinary {
//	        fmt.Printf("Processing: %s\n", file.RelativePath)
//	    }
//	}
//
//	// Check for errors during scan
//	for _, err := range result.Errors {
//	    log.Printf("Warning: %v\n", err)
//	}
//
// Pattern Matching:
//
// The package supports glob patterns for include/exclude rules:
//   - * matches any characters except /
//   - ** matches any characters including /
//   - ? matches a single character
//   - [abc] character classes
//
// Examples:
//   - "*.go" matches all .go files
//   - "**/*.go" matches all .go files in any directory
//   - "test/**" matches everything under test directories
//   - "**/node_modules/**" matches all node_modules directories and their contents
//
// Binary Detection:
//
// Files are classified as binary using multiple heuristics:
//   - Common binary file extensions (.exe, .png, .zip, etc.)
//   - Magic number signatures (ELF, PE, ZIP, etc.)
//   - Null byte detection in first 8KB
//   - High ratio of non-printable characters
//
// Gitignore Support:
//
// The scanner automatically loads and respects .gitignore files:
//   - Loads .gitignore from repository root
//   - Discovers and applies nested .gitignore files during traversal
//   - Custom patterns can be added via ScanOptions.ExcludePatterns
//   - Always ignores .git directories
package files
