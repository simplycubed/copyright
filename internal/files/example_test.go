package files_test

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/charlesgreen/copyright/internal/files"
)

// ExampleScanner_basic demonstrates basic usage of the file scanner
func ExampleScanner_basic() {
	// Create scanner with default options
	scanner, err := files.NewScanner(".", nil)
	if err != nil {
		log.Fatal(err)
	}

	// Scan the directory
	ctx := context.Background()
	result, err := scanner.Scan(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Process results
	fmt.Printf("Found %d files\n", len(result.Files))
	fmt.Printf("Skipped %d files\n", len(result.SkippedFiles))
	fmt.Printf("Encountered %d errors\n", len(result.Errors))
}

// ExampleScanner_withOptions demonstrates scanner with custom options
func ExampleScanner_withOptions() {
	// Create custom options
	options := &files.ScanOptions{
		ExcludePatterns: []string{
			"**/test/**",
			"**/*.tmp",
			"**/build/**",
		},
		IncludePatterns: []string{
			"**/*.go",
			"**/*.js",
			"**/*.py",
		},
		ExcludedExtensions: []string{".exe", ".dll", ".so"},
		SkipHidden:         true,
		FollowSymlinks:     false,
	}

	scanner, err := files.NewScanner(".", options)
	if err != nil {
		log.Fatal(err)
	}

	// Scan with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	result, err := scanner.Scan(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Filter for non-binary files
	var sourceFiles []files.FileInfo
	for _, file := range result.Files {
		if !file.IsBinary {
			sourceFiles = append(sourceFiles, file)
		}
	}

	fmt.Printf("Found %d source files\n", len(sourceFiles))
}

// ExampleScanner_errorHandling demonstrates error handling during scan
func ExampleScanner_errorHandling() {
	scanner, err := files.NewScanner(".", nil)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	result, err := scanner.Scan(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Check for permission errors and other issues
	if len(result.Errors) > 0 {
		fmt.Printf("Scan completed with %d warnings:\n", len(result.Errors))
		for _, err := range result.Errors {
			fmt.Printf("  - %v\n", err)
		}
	}

	// Process files
	for _, file := range result.Files {
		fmt.Printf("Processing: %s (%d bytes)\n", file.RelativePath, file.Size)
	}
}

// ExampleScanner_filterByExtension demonstrates filtering by extension
func ExampleScanner_filterByExtension() {
	scanner, err := files.NewScanner(".", nil)
	if err != nil {
		log.Fatal(err)
	}

	result, err := scanner.Scan(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Group files by extension
	byExtension := make(map[string][]files.FileInfo)
	for _, file := range result.Files {
		ext := file.Extension
		if ext == "" {
			ext = "(no extension)"
		}
		byExtension[ext] = append(byExtension[ext], file)
	}

	// Print summary
	for ext, files := range byExtension {
		fmt.Printf("%s: %d files\n", ext, len(files))
	}
}

// ExampleIsBinary demonstrates binary file detection
func ExampleIsBinary() {
	// Check various file types
	testFiles := []string{
		"main.go",
		"image.png",
		"document.pdf",
		"script.sh",
	}

	for _, file := range testFiles {
		isBinary, err := files.IsBinary(file)
		if err != nil {
			fmt.Printf("%s: error - %v\n", file, err)
			continue
		}
		fileType := "text"
		if isBinary {
			fileType = "binary"
		}
		fmt.Printf("%s: %s\n", file, fileType)
	}
}

// ExampleMatchPattern demonstrates glob pattern matching
func ExampleMatchPattern() {
	patterns := []string{
		"*.go",
		"**/*.js",
		"test/**",
		"**/node_modules/**",
	}

	testPaths := []string{
		"main.go",
		"src/app.js",
		"test/unit/test.go",
		"vendor/lib/code.go",
		"node_modules/pkg/index.js",
	}

	for _, pattern := range patterns {
		fmt.Printf("\nPattern: %s\n", pattern)
		for _, path := range testPaths {
			matched, err := files.MatchPattern(pattern, path)
			if err != nil {
				fmt.Printf("  %s: error - %v\n", path, err)
				continue
			}
			if matched {
				fmt.Printf("  %s: matched\n", path)
			}
		}
	}
}
