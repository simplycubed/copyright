// Package processor provides the core orchestration system for copyright management.
//
// The processor package coordinates all the components needed to check and fix
// copyright headers in source files. It brings together file scanning, copyright
// detection, header generation, and file modification in a cohesive workflow.
//
// # Architecture
//
// The processor package is the main orchestrator that integrates:
//   - internal/files: File discovery and scanning
//   - internal/language: Language-specific comment formatting
//   - internal/copyright: Copyright detection and header generation
//   - internal/license: LICENSE file management
//   - internal/config: Configuration management
//
// # Usage
//
// Basic usage with default configuration:
//
//	cfg := config.NewDefault()
//	cfg.Copyright.Company = "SimplyCubed, GK"
//	cfg.Copyright.SPDXIdentifier = "MIT"
//
//	proc, err := processor.NewProcessor(cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	ctx := context.Background()
//	result, err := proc.Process(ctx, "./src", processor.ModeCheck)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	fmt.Printf("Files checked: %d\n", result.TotalFiles)
//	fmt.Printf("Missing copyright: %d\n", result.FilesMissingCopyright)
//	os.Exit(result.ExitCode())
//
// # Modes
//
// The processor supports two operation modes:
//
// ModeCheck: Scans files and reports copyright status without modifying files.
// Returns exit code 1 if any files are missing copyright headers.
//
// ModeFix: Scans files and automatically adds copyright headers to files
// that are missing them. Preserves special lines like shebangs and encoding
// declarations.
//
// # Advanced Options
//
// For advanced use cases, use ProcessWithOptions:
//
//	opts := &processor.ProcessorOptions{
//	    DryRun: true,  // Preview changes without modifying files
//	    Backup: true,  // Create .bak files before modifying
//	}
//
//	result, err := proc.ProcessWithOptions(ctx, "./src", processor.ModeFix, opts)
//
// # Context Support
//
// All processing operations support context cancellation and timeouts:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//
//	result, err := proc.Process(ctx, "./src", processor.ModeCheck)
//	if err == context.DeadlineExceeded {
//	    fmt.Println("Processing timed out")
//	}
//
// # File Processing Workflow
//
// For each file discovered by the scanner:
//
//  1. Skip if binary file
//  2. Get language handler from registry
//  3. Skip if language doesn't support comments
//  4. Check if file has copyright using checker
//  5. In ModeCheck: record result
//  6. In ModeFix: add copyright if missing
//     a. Generate copyright text from template
//     b. Format header for language (comments)
//     c. Preserve special lines (shebang, encoding)
//     d. Insert header at correct position
//     e. Write modified content to file
//
// # Shebang and Special Line Preservation
//
// The processor correctly handles files with special leading lines:
//
//   - Shell scripts with shebangs (#!/bin/bash)
//   - Python files with encoding declarations
//   - Ruby files with encoding declarations
//   - XML files with XML declarations
//
// These lines are always preserved before the copyright header.
//
// # Result Types
//
// FileResult represents the processing result for a single file:
//   - StatusHasCopyright: File already has copyright
//   - StatusMissingCopyright: File is missing copyright (check mode)
//   - StatusFixed: Copyright was added to file (fix mode)
//   - StatusSkipped: File was skipped (unsupported type, binary, etc.)
//   - StatusFailed: Processing failed with an error
//
// ProcessResult aggregates results for all files and provides:
//   - Total file counts by status
//   - Detailed results for each file
//   - LICENSE file status
//   - ExitCode() for use with CLI tools
//   - Summary() for human-readable output
//
// # Error Handling
//
// The processor continues processing files even if individual files fail.
// Failed files are tracked in FileResults with their errors. The ProcessResult
// includes both successful and failed operations.
//
// Fatal errors (scanner creation, invalid paths) return immediately with an error.
// Non-fatal errors (individual file processing failures) are recorded in FileResults.
//
// # LICENSE File Management
//
// In ModeFix, the processor can automatically create a LICENSE file if:
//   - config.License.AutoCreate is true
//   - No LICENSE file exists
//   - A license type or custom text is configured
//
// The LICENSE file status is included in ProcessResult.
package processor
