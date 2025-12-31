package processor

import (
	"fmt"
	"os"
	"strings"

	"github.com/charlesgreen/copyright/internal/language"
)

// Writer handles writing copyright headers to files
type Writer struct {
	dryRun bool
	backup bool
}

// NewWriter creates a new writer with the specified options
func NewWriter(dryRun, backup bool) *Writer {
	return &Writer{
		dryRun: dryRun,
		backup: backup,
	}
}

// WriteHeader inserts a copyright header into a file
// It uses the language's PreserveLines and InsertPosition to correctly
// place the header after shebangs/encoding declarations
func (w *Writer) WriteHeader(filePath string, header string, lang language.Language) error {
	// Read the current file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Get lines that should be preserved (shebang, encoding, etc.)
	preservedLines := lang.PreserveLines(content)

	// Build the new content
	var newContent strings.Builder

	// Write preserved lines
	for _, line := range preservedLines {
		newContent.WriteString(line)
	}

	// Add a blank line after preserved content if there are preserved lines
	if len(preservedLines) > 0 {
		// Check if the last preserved line already ends with newline
		lastLine := preservedLines[len(preservedLines)-1]
		if !strings.HasSuffix(lastLine, "\n") {
			newContent.WriteString("\n")
		}
		newContent.WriteString("\n")
	}

	// Write the copyright header
	newContent.WriteString(header)

	// Add a blank line after the header
	if !strings.HasSuffix(header, "\n\n") {
		if strings.HasSuffix(header, "\n") {
			newContent.WriteString("\n")
		} else {
			newContent.WriteString("\n\n")
		}
	}

	// Get the insert position (byte offset after preserved lines)
	insertPos := lang.InsertPosition(content)

	// Write the remaining content after the insert position
	if insertPos < len(content) {
		remainingContent := content[insertPos:]
		// Skip leading blank lines in remaining content to avoid too much whitespace
		remainingStr := strings.TrimLeft(string(remainingContent), "\n")
		newContent.WriteString(remainingStr)
	}

	// Ensure file ends with a newline
	finalContent := newContent.String()
	if !strings.HasSuffix(finalContent, "\n") {
		finalContent += "\n"
	}

	// Write the new content to the file
	return w.WriteFile(filePath, []byte(finalContent))
}

// WriteFile writes content to a file (with backup if enabled)
func (w *Writer) WriteFile(filePath string, content []byte) error {
	// If dry-run mode, don't actually write
	if w.dryRun {
		return nil
	}

	// Create backup if enabled
	if w.backup {
		backupPath := filePath + ".bak"
		originalContent, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read original file for backup: %w", err)
		}
		if err := os.WriteFile(backupPath, originalContent, 0644); err != nil {
			return fmt.Errorf("failed to create backup file: %w", err)
		}
	}

	// Write the new content
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
