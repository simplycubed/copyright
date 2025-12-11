package files

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Scanner walks a directory tree and discovers files according to configured options
type Scanner struct {
	rootPath  string
	gitignore *GitignoreMatcher
	options   *ScanOptions
}

// NewScanner creates a new file scanner for the given root path
func NewScanner(rootPath string, options *ScanOptions) (*Scanner, error) {
	// Get absolute path
	absRoot, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Verify root path exists and is a directory
	info, err := os.Stat(absRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to stat root path: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("root path is not a directory: %s", absRoot)
	}

	// Use default options if none provided
	if options == nil {
		options = DefaultScanOptions()
	}

	// Initialize gitignore matcher
	gitignoreMatcher, err := NewGitignoreMatcher(absRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize gitignore matcher: %w", err)
	}

	// Add custom exclude patterns to gitignore
	if len(options.ExcludePatterns) > 0 {
		if err := gitignoreMatcher.AddPatterns(options.ExcludePatterns); err != nil {
			return nil, fmt.Errorf("failed to add exclude patterns: %w", err)
		}
	}

	return &Scanner{
		rootPath:  absRoot,
		gitignore: gitignoreMatcher,
		options:   options,
	}, nil
}

// Scan walks the directory tree and returns discovered files
func (s *Scanner) Scan(ctx context.Context) (*ScanResult, error) {
	result := &ScanResult{
		Files:        make([]FileInfo, 0),
		SkippedFiles: make([]string, 0),
		Errors:       make([]error, 0),
	}

	// Walk the directory tree
	err := filepath.WalkDir(s.rootPath, func(path string, d fs.DirEntry, err error) error {
		// Check for context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Handle errors from WalkDir
		if err != nil {
			// Log permission errors but continue
			if os.IsPermission(err) {
				result.Errors = append(result.Errors, fmt.Errorf("permission denied: %s", path))
				return nil
			}
			return err
		}

		// Get relative path
		relPath, err := filepath.Rel(s.rootPath, path)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("failed to get relative path for %s: %w", path, err))
			return nil
		}

		// Handle directories
		if d.IsDir() {
			// Skip .git directory always
			if d.Name() == ".git" {
				return filepath.SkipDir
			}

			// Skip hidden directories if configured
			if s.options.SkipHidden && strings.HasPrefix(d.Name(), ".") && path != s.rootPath {
				result.SkippedFiles = append(result.SkippedFiles, relPath)
				return filepath.SkipDir
			}

			// Check gitignore patterns
			if s.gitignore.ShouldIgnoreDir(path) {
				result.SkippedFiles = append(result.SkippedFiles, relPath)
				return filepath.SkipDir
			}

			// Try to load nested .gitignore
			_ = s.gitignore.LoadNestedGitignore(path)

			return nil
		}

		// Handle symlinks
		if d.Type()&os.ModeSymlink != 0 {
			if !s.options.FollowSymlinks {
				result.SkippedFiles = append(result.SkippedFiles, relPath)
				return nil
			}
			// If following symlinks, resolve and check if it's a file
			realPath, err := filepath.EvalSymlinks(path)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Errorf("failed to resolve symlink %s: %w", path, err))
				return nil
			}
			info, err := os.Stat(realPath)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Errorf("failed to stat symlink target %s: %w", realPath, err))
				return nil
			}
			if info.IsDir() {
				// Don't follow directory symlinks even if FollowSymlinks is true
				result.SkippedFiles = append(result.SkippedFiles, relPath)
				return nil
			}
		}

		// Process files only
		if !d.Type().IsRegular() {
			return nil
		}

		// Skip hidden files if configured
		if s.options.SkipHidden && strings.HasPrefix(d.Name(), ".") {
			result.SkippedFiles = append(result.SkippedFiles, relPath)
			return nil
		}

		// Check gitignore patterns
		if s.gitignore.Match(path) {
			result.SkippedFiles = append(result.SkippedFiles, relPath)
			return nil
		}

		// Check excluded extensions
		ext := filepath.Ext(path)
		if s.isExcludedExtension(ext) {
			result.SkippedFiles = append(result.SkippedFiles, relPath)
			return nil
		}

		// Check include patterns (if specified, only include matching files)
		if len(s.options.IncludePatterns) > 0 {
			if !MatchAnyPattern(s.options.IncludePatterns, relPath) {
				result.SkippedFiles = append(result.SkippedFiles, relPath)
				return nil
			}
		}

		// Get file info
		info, err := d.Info()
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("failed to get file info for %s: %w", path, err))
			return nil
		}

		// Check if binary
		isBinary, err := IsBinary(path)
		if err != nil {
			// If we can't determine, assume it's binary and skip
			result.Errors = append(result.Errors, fmt.Errorf("failed to check if binary %s: %w", path, err))
			result.SkippedFiles = append(result.SkippedFiles, relPath)
			return nil
		}

		// Create FileInfo
		fileInfo := FileInfo{
			Path:         path,
			RelativePath: relPath,
			Extension:    ext,
			Filename:     d.Name(),
			Size:         info.Size(),
			IsBinary:     isBinary,
		}

		result.Files = append(result.Files, fileInfo)
		return nil
	})

	if err != nil && err != context.Canceled && err != context.DeadlineExceeded {
		return nil, fmt.Errorf("failed to walk directory tree: %w", err)
	}

	// If context was canceled, include that as an error
	if err == context.Canceled || err == context.DeadlineExceeded {
		result.Errors = append(result.Errors, err)
	}

	return result, nil
}

// isExcludedExtension checks if a file extension is in the excluded list
func (s *Scanner) isExcludedExtension(ext string) bool {
	if len(s.options.ExcludedExtensions) == 0 {
		return false
	}

	// Normalize extension (ensure it starts with .)
	if ext != "" && !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	for _, excluded := range s.options.ExcludedExtensions {
		// Normalize excluded extension
		if excluded != "" && !strings.HasPrefix(excluded, ".") {
			excluded = "." + excluded
		}
		if strings.EqualFold(ext, excluded) {
			return true
		}
	}

	return false
}

// GetRootPath returns the scanner's root path
func (s *Scanner) GetRootPath() string {
	return s.rootPath
}
