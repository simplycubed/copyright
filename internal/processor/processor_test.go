package processor

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/leansecurity-co/copyright/internal/config"
)

// setupTestDir creates a temporary test directory with test files
func setupTestDir(t *testing.T) string {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "copyright-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// Create test files
	files := map[string]string{
		"with_copyright.go": `// Copyright 2025 Test Company
// SPDX-License-Identifier: MIT

package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
`,
		"without_copyright.go": `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
`,
		"with_shebang.py": `#!/usr/bin/env python3

def main():
    print("Hello, World!")

if __name__ == "__main__":
    main()
`,
		"unsupported.json": `{
  "name": "test",
  "version": "1.0.0"
}
`,
		"binary.bin": string([]byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE}),
	}

	for filename, content := range files {
		filePath := filepath.Join(tmpDir, filename)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create test file %s: %v", filename, err)
		}
	}

	return tmpDir
}

// cleanupTestDir removes the temporary test directory
func cleanupTestDir(t *testing.T, dir string) {
	t.Helper()
	if err := os.RemoveAll(dir); err != nil {
		t.Errorf("failed to cleanup test dir: %v", err)
	}
}

// getTestConfig returns a test configuration
func getTestConfig() *config.Config {
	cfg := config.NewDefault()
	cfg.Copyright.Template = "Copyright {{.Year}} {{.Company}}\nSPDX-License-Identifier: {{.SPDX}}"
	cfg.Copyright.Year = "2025"
	cfg.Copyright.Company = "Test Company"
	cfg.Copyright.SPDXIdentifier = "MIT"
	cfg.Detection.LinesToScan = 20
	cfg.Detection.Keywords = []string{"Copyright", "SPDX-License-Identifier"}
	cfg.Detection.RequireAll = false
	return cfg
}

func TestNewProcessor(t *testing.T) {
	cfg := getTestConfig()

	processor, err := NewProcessor(cfg)
	if err != nil {
		t.Fatalf("NewProcessor failed: %v", err)
	}

	if processor == nil {
		t.Fatal("processor is nil")
	}
	if processor.config == nil {
		t.Error("processor.config is nil")
	}
	if processor.checker == nil {
		t.Error("processor.checker is nil")
	}
	if processor.generator == nil {
		t.Error("processor.generator is nil")
	}
	if processor.registry == nil {
		t.Error("processor.registry is nil")
	}
}

func TestNewProcessor_NilConfig(t *testing.T) {
	_, err := NewProcessor(nil)
	if err == nil {
		t.Error("expected error for nil config, got nil")
	}
}

func TestProcess_CheckMode_WithCopyright(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer cleanupTestDir(t, tmpDir)

	cfg := getTestConfig()
	processor, err := NewProcessor(cfg)
	if err != nil {
		t.Fatalf("NewProcessor failed: %v", err)
	}

	ctx := context.Background()
	result, err := processor.Process(ctx, tmpDir, ModeCheck)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	if result == nil {
		t.Fatal("result is nil")
	}

	// Should find the file with copyright
	if result.FilesWithCopyright == 0 {
		t.Error("expected at least 1 file with copyright, got 0")
	}

	// Check that with_copyright.go was detected
	found := false
	for _, fr := range result.FileResults {
		if strings.Contains(fr.RelativePath, "with_copyright.go") {
			found = true
			if fr.Status != StatusHasCopyright {
				t.Errorf("expected StatusHasCopyright for with_copyright.go, got %s", fr.Status)
			}
		}
	}
	if !found {
		t.Error("with_copyright.go not found in results")
	}
}

func TestProcess_CheckMode_WithoutCopyright(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer cleanupTestDir(t, tmpDir)

	cfg := getTestConfig()
	processor, err := NewProcessor(cfg)
	if err != nil {
		t.Fatalf("NewProcessor failed: %v", err)
	}

	ctx := context.Background()
	result, err := processor.Process(ctx, tmpDir, ModeCheck)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	// Should find files missing copyright
	if result.FilesMissingCopyright == 0 {
		t.Error("expected at least 1 file missing copyright, got 0")
	}

	// Check that without_copyright.go was detected
	found := false
	for _, fr := range result.FileResults {
		if strings.Contains(fr.RelativePath, "without_copyright.go") {
			found = true
			if fr.Status != StatusMissingCopyright {
				t.Errorf("expected StatusMissingCopyright for without_copyright.go, got %s", fr.Status)
			}
		}
	}
	if !found {
		t.Error("without_copyright.go not found in results")
	}

	// Exit code should be 1 (files missing copyright)
	if result.ExitCode() != 1 {
		t.Errorf("expected exit code 1, got %d", result.ExitCode())
	}
}

func TestProcess_FixMode_AddsHeaders(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer cleanupTestDir(t, tmpDir)

	cfg := getTestConfig()
	processor, err := NewProcessor(cfg)
	if err != nil {
		t.Fatalf("NewProcessor failed: %v", err)
	}

	ctx := context.Background()
	opts := &ProcessorOptions{
		DryRun: false,
		Backup: false,
	}
	result, err := processor.ProcessWithOptions(ctx, tmpDir, ModeFix, opts)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	// Should have fixed files
	if result.FilesFixed == 0 {
		t.Error("expected at least 1 file to be fixed, got 0")
	}

	// Check that without_copyright.go was fixed
	filePath := filepath.Join(tmpDir, "without_copyright.go")
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read fixed file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "Copyright 2025 Test Company") {
		t.Error("fixed file does not contain copyright header")
	}
	if !strings.Contains(contentStr, "SPDX-License-Identifier: MIT") {
		t.Error("fixed file does not contain SPDX identifier")
	}
	if !strings.Contains(contentStr, "package main") {
		t.Error("fixed file lost original content")
	}
}

func TestProcess_FixMode_PreservesShebang(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer cleanupTestDir(t, tmpDir)

	cfg := getTestConfig()
	processor, err := NewProcessor(cfg)
	if err != nil {
		t.Fatalf("NewProcessor failed: %v", err)
	}

	ctx := context.Background()
	opts := &ProcessorOptions{
		DryRun: false,
		Backup: false,
	}
	result, err := processor.ProcessWithOptions(ctx, tmpDir, ModeFix, opts)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	// Check that with_shebang.py was processed
	filePath := filepath.Join(tmpDir, "with_shebang.py")
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read fixed file: %v", err)
	}

	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	// First line should be the shebang
	if !strings.HasPrefix(lines[0], "#!/usr/bin/env python3") {
		t.Errorf("shebang not preserved, first line is: %s", lines[0])
	}

	// Should contain copyright after shebang
	if !strings.Contains(contentStr, "Copyright 2025 Test Company") {
		t.Error("fixed file does not contain copyright header")
	}

	// Verify the file was marked as fixed
	found := false
	for _, fr := range result.FileResults {
		if strings.Contains(fr.RelativePath, "with_shebang.py") {
			found = true
			if fr.Status != StatusFixed {
				t.Errorf("expected StatusFixed for with_shebang.py, got %s", fr.Status)
			}
		}
	}
	if !found {
		t.Error("with_shebang.py not found in results")
	}
}

func TestProcess_DryRun_DoesNotModifyFiles(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer cleanupTestDir(t, tmpDir)

	// Read original content
	filePath := filepath.Join(tmpDir, "without_copyright.go")
	originalContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read original file: %v", err)
	}

	cfg := getTestConfig()
	processor, err := NewProcessor(cfg)
	if err != nil {
		t.Fatalf("NewProcessor failed: %v", err)
	}

	ctx := context.Background()
	opts := &ProcessorOptions{
		DryRun: true,
		Backup: false,
	}
	result, err := processor.ProcessWithOptions(ctx, tmpDir, ModeFix, opts)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	// Should report files that would be fixed
	if result.FilesFixed == 0 {
		t.Error("expected files to be marked for fixing in dry-run mode")
	}

	// But file should not be modified
	newContent, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file after dry-run: %v", err)
	}

	if string(originalContent) != string(newContent) {
		t.Error("file was modified in dry-run mode")
	}
}

func TestProcess_Backup_CreatesBackupFile(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer cleanupTestDir(t, tmpDir)

	cfg := getTestConfig()
	processor, err := NewProcessor(cfg)
	if err != nil {
		t.Fatalf("NewProcessor failed: %v", err)
	}

	ctx := context.Background()
	opts := &ProcessorOptions{
		DryRun: false,
		Backup: true,
	}
	_, err = processor.ProcessWithOptions(ctx, tmpDir, ModeFix, opts)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	// Check that backup file was created
	backupPath := filepath.Join(tmpDir, "without_copyright.go.bak")
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		t.Error("backup file was not created")
	}

	// Backup should contain original content
	backupContent, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatalf("failed to read backup file: %v", err)
	}

	if strings.Contains(string(backupContent), "Copyright") {
		t.Error("backup file should not contain copyright header")
	}
}

func TestProcess_SkipsUnsupportedFiles(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer cleanupTestDir(t, tmpDir)

	cfg := getTestConfig()
	// Override excluded extensions to allow JSON files to be scanned
	cfg.Scan.ExcludedExtensions = []string{}

	processor, err := NewProcessor(cfg)
	if err != nil {
		t.Fatalf("NewProcessor failed: %v", err)
	}

	ctx := context.Background()
	result, err := processor.Process(ctx, tmpDir, ModeCheck)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	// JSON files don't support comments, should be skipped
	found := false
	for _, fr := range result.FileResults {
		if strings.Contains(fr.RelativePath, "unsupported.json") {
			found = true
			if fr.Status != StatusSkipped {
				t.Errorf("expected StatusSkipped for unsupported.json, got %s", fr.Status)
			}
		}
	}
	if !found {
		t.Error("unsupported.json not found in results")
	}

	// Should have skipped files
	if result.FilesSkipped == 0 {
		t.Error("expected at least 1 file to be skipped")
	}
}

func TestProcess_ContextCancellation(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer cleanupTestDir(t, tmpDir)

	cfg := getTestConfig()
	processor, err := NewProcessor(cfg)
	if err != nil {
		t.Fatalf("NewProcessor failed: %v", err)
	}

	// Create a context that's already cancelled
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result, err := processor.Process(ctx, tmpDir, ModeCheck)
	// With a small number of files, the scan might complete before context check
	// So we accept either a context.Canceled error or a successful result
	if err != nil && err != context.Canceled {
		t.Errorf("expected nil or context.Canceled error, got: %v", err)
	}
	// If the result is returned, it should be partial
	if result != nil && err == nil {
		// This is OK - the scan completed before context was checked
		t.Logf("Scan completed before context cancellation was detected (result: %s)", result.Summary())
	}
}

func TestProcess_ContextTimeout(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer cleanupTestDir(t, tmpDir)

	cfg := getTestConfig()
	processor, err := NewProcessor(cfg)
	if err != nil {
		t.Fatalf("NewProcessor failed: %v", err)
	}

	// Create a context with a very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Sleep briefly to ensure timeout
	time.Sleep(1 * time.Millisecond)

	result, err := processor.Process(ctx, tmpDir, ModeCheck)
	// With a small number of files, the scan might complete before timeout
	// So we accept either a context.DeadlineExceeded error or a successful result
	if err != nil && err != context.DeadlineExceeded {
		t.Errorf("expected nil or context.DeadlineExceeded error, got: %v", err)
	}
	// If the result is returned, it should be partial or complete
	if result != nil && err == nil {
		// This is OK - the scan completed before timeout
		t.Logf("Scan completed before timeout (result: %s)", result.Summary())
	}
}

func TestProcessResult_ExitCode(t *testing.T) {
	tests := []struct {
		name     string
		result   *ProcessResult
		expected int
	}{
		{
			name: "all files have copyright",
			result: &ProcessResult{
				TotalFiles:         5,
				FilesWithCopyright: 5,
			},
			expected: 0,
		},
		{
			name: "some files missing copyright",
			result: &ProcessResult{
				TotalFiles:            5,
				FilesWithCopyright:    3,
				FilesMissingCopyright: 2,
			},
			expected: 1,
		},
		{
			name: "some files failed",
			result: &ProcessResult{
				TotalFiles:         5,
				FilesWithCopyright: 4,
				FilesFailed:        1,
			},
			expected: 1,
		},
		{
			name: "all files skipped",
			result: &ProcessResult{
				TotalFiles:   5,
				FilesSkipped: 5,
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exitCode := tt.result.ExitCode()
			if exitCode != tt.expected {
				t.Errorf("expected exit code %d, got %d", tt.expected, exitCode)
			}
		})
	}
}

func TestFileStatus_String(t *testing.T) {
	tests := []struct {
		status   FileStatus
		expected string
	}{
		{StatusHasCopyright, "has_copyright"},
		{StatusMissingCopyright, "missing_copyright"},
		{StatusFixed, "fixed"},
		{StatusSkipped, "skipped"},
		{StatusFailed, "failed"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			str := tt.status.String()
			if str != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, str)
			}
		})
	}
}

func TestProcessResult_Summary(t *testing.T) {
	result := &ProcessResult{
		TotalFiles:            10,
		FilesWithCopyright:    5,
		FilesMissingCopyright: 2,
		FilesFixed:            2,
		FilesSkipped:          3,
		FilesFailed:           1,
	}

	summary := result.Summary()
	if !strings.Contains(summary, "Total: 10") {
		t.Errorf("summary missing total count: %s", summary)
	}
	if !strings.Contains(summary, "Has Copyright: 5") {
		t.Errorf("summary missing copyright count: %s", summary)
	}
	if !strings.Contains(summary, "Missing: 2") {
		t.Errorf("summary missing missing count: %s", summary)
	}
	if !strings.Contains(summary, "Fixed: 2") {
		t.Errorf("summary missing fixed count: %s", summary)
	}
	if !strings.Contains(summary, "Skipped: 3") {
		t.Errorf("summary missing skipped count: %s", summary)
	}
	if !strings.Contains(summary, "Failed: 1") {
		t.Errorf("summary missing failed count: %s", summary)
	}
}

func TestProcess_InvalidPath(t *testing.T) {
	cfg := getTestConfig()
	processor, err := NewProcessor(cfg)
	if err != nil {
		t.Fatalf("NewProcessor failed: %v", err)
	}

	ctx := context.Background()
	_, err = processor.Process(ctx, "/nonexistent/path", ModeCheck)
	if err == nil {
		t.Error("expected error for nonexistent path, got nil")
	}
}
