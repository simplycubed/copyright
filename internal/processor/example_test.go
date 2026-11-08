package processor_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/simplycubed/copyright/internal/config"
	"github.com/simplycubed/copyright/internal/processor"
)

// ExampleProcessor_check demonstrates checking files for copyright headers
func ExampleProcessor_check() {
	// Create a temporary directory with test files
	tmpDir, err := os.MkdirTemp("", "copyright-example-*")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create a test file without copyright
	testFile := filepath.Join(tmpDir, "main.go")
	content := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		log.Fatal(err)
	}

	// Configure processor
	cfg := config.NewDefault()
	cfg.Copyright.Company = "SimplyCubed, GK"
	cfg.Copyright.SPDXIdentifier = "MIT"
	cfg.Copyright.Year = "2025"

	// Create processor
	proc, err := processor.NewProcessor(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Check files for copyright
	ctx := context.Background()
	result, err := proc.Process(ctx, tmpDir, processor.ModeCheck)
	if err != nil {
		log.Fatal(err)
	}

	// Print summary
	fmt.Printf("Total files: %d\n", result.TotalFiles)
	fmt.Printf("Files with copyright: %d\n", result.FilesWithCopyright)
	fmt.Printf("Files missing copyright: %d\n", result.FilesMissingCopyright)

	// Output:
	// Total files: 1
	// Files with copyright: 0
	// Files missing copyright: 1
}

// ExampleProcessor_fix demonstrates adding copyright headers to files
func ExampleProcessor_fix() {
	// Create a temporary directory with test files
	tmpDir, err := os.MkdirTemp("", "copyright-example-*")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create a test file without copyright
	testFile := filepath.Join(tmpDir, "main.go")
	content := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		log.Fatal(err)
	}

	// Configure processor
	cfg := config.NewDefault()
	cfg.Copyright.Company = "SimplyCubed, GK"
	cfg.Copyright.SPDXIdentifier = "MIT"
	cfg.Copyright.Year = "2025"

	// Create processor
	proc, err := processor.NewProcessor(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Fix files (add copyright headers)
	ctx := context.Background()
	opts := &processor.ProcessorOptions{
		DryRun: false,
		Backup: false,
	}
	result, err := proc.ProcessWithOptions(ctx, tmpDir, processor.ModeFix, opts)
	if err != nil {
		log.Fatal(err)
	}

	// Print summary
	fmt.Printf("Total files: %d\n", result.TotalFiles)
	fmt.Printf("Files fixed: %d\n", result.FilesFixed)

	// Output:
	// Total files: 1
	// Files fixed: 1
}

// ExampleProcessor_dryRun demonstrates previewing changes without modifying files
func ExampleProcessor_dryRun() {
	// Create a temporary directory with test files
	tmpDir, err := os.MkdirTemp("", "copyright-example-*")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create a test file without copyright
	testFile := filepath.Join(tmpDir, "main.go")
	content := `package main

func main() {}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		log.Fatal(err)
	}

	// Configure processor
	cfg := config.NewDefault()
	cfg.Copyright.Company = "SimplyCubed, GK"
	cfg.Copyright.SPDXIdentifier = "MIT"

	proc, err := processor.NewProcessor(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Dry-run: preview changes without modifying files
	ctx := context.Background()
	opts := &processor.ProcessorOptions{
		DryRun: true,
		Backup: false,
	}
	result, err := proc.ProcessWithOptions(ctx, tmpDir, processor.ModeFix, opts)
	if err != nil {
		log.Fatal(err)
	}

	// Files would be fixed, but not actually modified
	fmt.Printf("Would fix %d files\n", result.FilesFixed)

	// Verify file wasn't modified
	modifiedContent, _ := os.ReadFile(testFile)
	if string(modifiedContent) == content {
		fmt.Println("File not modified (dry-run)")
	}

	// Output:
	// Would fix 1 files
	// File not modified (dry-run)
}

// ExampleProcessResult_ExitCode demonstrates getting exit codes for CLI tools
func ExampleProcessResult_ExitCode() {
	// Simulate different result scenarios
	scenarios := []struct {
		name   string
		result *processor.ProcessResult
	}{
		{
			name: "all files have copyright",
			result: &processor.ProcessResult{
				TotalFiles:         5,
				FilesWithCopyright: 5,
			},
		},
		{
			name: "some files missing copyright",
			result: &processor.ProcessResult{
				TotalFiles:            5,
				FilesWithCopyright:    3,
				FilesMissingCopyright: 2,
			},
		},
		{
			name: "some files failed",
			result: &processor.ProcessResult{
				TotalFiles:         5,
				FilesWithCopyright: 4,
				FilesFailed:        1,
			},
		},
	}

	for _, scenario := range scenarios {
		exitCode := scenario.result.ExitCode()
		fmt.Printf("%s: exit code %d\n", scenario.name, exitCode)
	}

	// Output:
	// all files have copyright: exit code 0
	// some files missing copyright: exit code 1
	// some files failed: exit code 1
}
