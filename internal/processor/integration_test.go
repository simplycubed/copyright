package processor

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/charlesgreen/copyright/internal/config"
)

// TestIntegration_FullWorkflow tests the complete workflow from check to fix
func TestIntegration_FullWorkflow(t *testing.T) {
	// Create a temporary directory with multiple file types
	tmpDir, err := os.MkdirTemp("", "copyright-integration-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create test files in different languages
	files := map[string]string{
		"main.go": `package main

import "fmt"

func main() {
	fmt.Println("Hello from Go")
}
`,
		"script.py": `#!/usr/bin/env python3

def main():
    print("Hello from Python")

if __name__ == "__main__":
    main()
`,
		"app.js": `function main() {
    console.log("Hello from JavaScript");
}

main();
`,
		"styles.css": `body {
    margin: 0;
    padding: 0;
}
`,
		"README.md": `# Test Project

This is a test project.
`,
	}

	for filename, content := range files {
		filePath := filepath.Join(tmpDir, filename)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create test file %s: %v", filename, err)
		}
	}

	// Configure processor
	cfg := config.NewDefault()
	cfg.Copyright.Template = "Copyright {{.Year}} {{.Company}}\nSPDX-License-Identifier: {{.SPDX}}"
	cfg.Copyright.Year = "2025"
	cfg.Copyright.Company = "Integration Test Corp"
	cfg.Copyright.SPDXIdentifier = "Apache-2.0"
	cfg.Detection.LinesToScan = 20
	cfg.Detection.Keywords = []string{"Copyright", "SPDX-License-Identifier"}
	cfg.Detection.RequireAll = false
	cfg.Scan.ExcludedExtensions = []string{} // Don't exclude anything for this test

	processor, err := NewProcessor(cfg)
	if err != nil {
		t.Fatalf("NewProcessor failed: %v", err)
	}

	ctx := context.Background()

	// Step 1: Check mode - all files should be missing copyright
	t.Run("initial_check", func(t *testing.T) {
		result, err := processor.Process(ctx, tmpDir, ModeCheck)
		if err != nil {
			t.Fatalf("Process check failed: %v", err)
		}

		// Markdown files don't support comments, so they should be skipped
		// All other files should be missing copyright
		expectedMissing := 0
		expectedSkipped := 0
		for filename := range files {
			if strings.HasSuffix(filename, ".md") {
				expectedSkipped++
			} else {
				expectedMissing++
			}
		}

		if result.FilesMissingCopyright != expectedMissing {
			t.Errorf("expected %d files missing copyright, got %d", expectedMissing, result.FilesMissingCopyright)
		}

		if result.FilesSkipped != expectedSkipped {
			t.Errorf("expected %d files skipped, got %d", expectedSkipped, result.FilesSkipped)
		}

		if result.ExitCode() != 1 {
			t.Error("expected exit code 1 for missing copyright")
		}
	})

	// Step 2: Fix mode - add copyright headers
	t.Run("fix", func(t *testing.T) {
		opts := &ProcessorOptions{
			DryRun: false,
			Backup: true,
		}

		result, err := processor.ProcessWithOptions(ctx, tmpDir, ModeFix, opts)
		if err != nil {
			t.Fatalf("Process fix failed: %v", err)
		}

		// Should have fixed all files except markdown
		expectedFixed := 0
		for filename := range files {
			if !strings.HasSuffix(filename, ".md") {
				expectedFixed++
			}
		}

		if result.FilesFixed != expectedFixed {
			t.Errorf("expected %d files fixed, got %d", expectedFixed, result.FilesFixed)
		}

		// Verify backup files were created
		for filename := range files {
			if !strings.HasSuffix(filename, ".md") {
				backupPath := filepath.Join(tmpDir, filename+".bak")
				if _, err := os.Stat(backupPath); os.IsNotExist(err) {
					t.Errorf("backup file not created for %s", filename)
				}
			}
		}
	})

	// Step 3: Verify file contents
	t.Run("verify_contents", func(t *testing.T) {
		// Check Go file
		goContent, err := os.ReadFile(filepath.Join(tmpDir, "main.go"))
		if err != nil {
			t.Fatalf("failed to read main.go: %v", err)
		}
		goStr := string(goContent)
		if !strings.Contains(goStr, "Copyright 2025 Integration Test Corp") {
			t.Error("main.go missing copyright")
		}
		if !strings.Contains(goStr, "SPDX-License-Identifier: Apache-2.0") {
			t.Error("main.go missing SPDX identifier")
		}
		if !strings.Contains(goStr, "package main") {
			t.Error("main.go lost original content")
		}

		// Check Python file - shebang should be preserved
		pyContent, err := os.ReadFile(filepath.Join(tmpDir, "script.py"))
		if err != nil {
			t.Fatalf("failed to read script.py: %v", err)
		}
		pyStr := string(pyContent)
		lines := strings.Split(pyStr, "\n")
		if !strings.HasPrefix(lines[0], "#!/usr/bin/env python3") {
			t.Error("script.py lost shebang")
		}
		if !strings.Contains(pyStr, "Copyright 2025 Integration Test Corp") {
			t.Error("script.py missing copyright")
		}

		// Check JavaScript file
		jsContent, err := os.ReadFile(filepath.Join(tmpDir, "app.js"))
		if err != nil {
			t.Fatalf("failed to read app.js: %v", err)
		}
		jsStr := string(jsContent)
		if !strings.Contains(jsStr, "Copyright 2025 Integration Test Corp") {
			t.Error("app.js missing copyright")
		}
		if !strings.Contains(jsStr, "function main()") {
			t.Error("app.js lost original content")
		}

		// Check CSS file
		cssContent, err := os.ReadFile(filepath.Join(tmpDir, "styles.css"))
		if err != nil {
			t.Fatalf("failed to read styles.css: %v", err)
		}
		cssStr := string(cssContent)
		if !strings.Contains(cssStr, "Copyright 2025 Integration Test Corp") {
			t.Error("styles.css missing copyright")
		}
		if !strings.Contains(cssStr, "body {") {
			t.Error("styles.css lost original content")
		}
	})

	// Step 4: Check mode again - all files should have copyright
	t.Run("final_check", func(t *testing.T) {
		result, err := processor.Process(ctx, tmpDir, ModeCheck)
		if err != nil {
			t.Fatalf("Process check failed: %v", err)
		}

		expectedWithCopyright := 0
		for filename := range files {
			if !strings.HasSuffix(filename, ".md") {
				expectedWithCopyright++
			}
		}

		if result.FilesWithCopyright != expectedWithCopyright {
			t.Errorf("expected %d files with copyright, got %d", expectedWithCopyright, result.FilesWithCopyright)
		}

		if result.FilesMissingCopyright != 0 {
			t.Errorf("expected 0 files missing copyright, got %d", result.FilesMissingCopyright)
		}

		if result.ExitCode() != 0 {
			t.Error("expected exit code 0 after fixing")
		}
	})

	// Step 5: Fix mode again should be idempotent (no changes)
	t.Run("idempotent_fix", func(t *testing.T) {
		result, err := processor.Process(ctx, tmpDir, ModeFix)
		if err != nil {
			t.Fatalf("Process fix failed: %v", err)
		}

		if result.FilesFixed != 0 {
			t.Errorf("expected 0 files to be fixed (idempotent), got %d", result.FilesFixed)
		}

		expectedWithCopyright := 0
		for filename := range files {
			if !strings.HasSuffix(filename, ".md") {
				expectedWithCopyright++
			}
		}

		if result.FilesWithCopyright != expectedWithCopyright {
			t.Errorf("expected %d files with copyright, got %d", expectedWithCopyright, result.FilesWithCopyright)
		}
	})
}

// TestIntegration_SubdirectoryProcessing tests processing nested directories
func TestIntegration_SubdirectoryProcessing(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "copyright-subdirs-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create nested directory structure
	srcDir := filepath.Join(tmpDir, "src")
	libDir := filepath.Join(tmpDir, "src", "lib")
	if err := os.MkdirAll(libDir, 0755); err != nil {
		t.Fatalf("failed to create lib dir: %v", err)
	}

	// Create files in different directories
	files := map[string]string{
		filepath.Join(srcDir, "main.go"):     `package main`,
		filepath.Join(libDir, "helper.go"):   `package lib`,
		filepath.Join(tmpDir, "README.md"):   `# Project`,
		filepath.Join(tmpDir, "config.json"): `{"name": "test"}`,
	}

	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create file %s: %v", path, err)
		}
	}

	cfg := config.NewDefault()
	cfg.Copyright.Company = "Test"
	cfg.Copyright.SPDXIdentifier = "MIT"
	cfg.Scan.ExcludedExtensions = []string{} // Process all files

	processor, err := NewProcessor(cfg)
	if err != nil {
		t.Fatalf("NewProcessor failed: %v", err)
	}

	ctx := context.Background()

	// Process from root
	result, err := processor.Process(ctx, tmpDir, ModeFix)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	// Should process Go files from both directories
	if result.FilesFixed < 2 {
		t.Errorf("expected at least 2 files fixed, got %d", result.FilesFixed)
	}

	// Verify both Go files were processed
	for path := range files {
		if strings.HasSuffix(path, ".go") {
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("failed to read %s: %v", path, err)
			}
			if !strings.Contains(string(content), "Copyright") {
				t.Errorf("file %s missing copyright", path)
			}
		}
	}
}
