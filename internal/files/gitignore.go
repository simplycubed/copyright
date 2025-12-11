package files

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
)

// GitignoreMatcher handles gitignore pattern matching
type GitignoreMatcher struct {
	rootPath string
	matchers map[string]*ignore.GitIgnore // path -> matcher for nested .gitignore files
	global   *ignore.GitIgnore            // global patterns
}

// NewGitignoreMatcher creates a new gitignore matcher for the given root path
func NewGitignoreMatcher(rootPath string) (*GitignoreMatcher, error) {
	absRoot, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	gm := &GitignoreMatcher{
		rootPath: absRoot,
		matchers: make(map[string]*ignore.GitIgnore),
	}

	// Load root .gitignore if it exists
	if err := gm.loadGitignoreFile(absRoot); err != nil {
		// If .gitignore doesn't exist, that's okay
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	// Initialize global matcher with empty patterns (will be populated via AddPatterns)
	gm.global = ignore.CompileIgnoreLines()

	return gm, nil
}

// LoadNestedGitignore loads a .gitignore file from a subdirectory
func (g *GitignoreMatcher) LoadNestedGitignore(dirPath string) error {
	return g.loadGitignoreFile(dirPath)
}

// loadGitignoreFile loads a .gitignore file from a directory
func (g *GitignoreMatcher) loadGitignoreFile(dirPath string) error {
	gitignorePath := filepath.Join(dirPath, ".gitignore")

	// Check if .gitignore exists
	info, err := os.Stat(gitignorePath)
	if err != nil {
		return err
	}

	// Don't process if it's not a regular file
	if !info.Mode().IsRegular() {
		return nil
	}

	// Compile the gitignore file
	matcher, err := ignore.CompileIgnoreFile(gitignorePath)
	if err != nil {
		return fmt.Errorf("failed to compile gitignore file %s: %w", gitignorePath, err)
	}

	// Store the matcher for this directory
	g.matchers[dirPath] = matcher

	return nil
}

// Match checks if a path should be ignored based on gitignore rules
func (g *GitignoreMatcher) Match(path string) bool {
	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	// Get relative path from root
	relPath, err := filepath.Rel(g.rootPath, absPath)
	if err != nil {
		return false
	}

	// Normalize path separators for gitignore matching
	relPath = filepath.ToSlash(relPath)

	// Check global patterns first
	if g.global != nil && g.global.MatchesPath(relPath) {
		return true
	}

	// Check each .gitignore from root to the file's directory
	segments := strings.Split(filepath.Dir(absPath), string(filepath.Separator))

	for i := len(strings.Split(g.rootPath, string(filepath.Separator))); i <= len(segments); i++ {
		checkDir := strings.Join(segments[:i], string(filepath.Separator))
		if checkDir == "" {
			continue
		}

		// Get relative path from this directory
		relFromDir, err := filepath.Rel(checkDir, absPath)
		if err != nil {
			continue
		}
		relFromDir = filepath.ToSlash(relFromDir)

		// Check if this directory has a gitignore matcher
		if matcher, exists := g.matchers[checkDir]; exists {
			if matcher.MatchesPath(relFromDir) {
				return true
			}
			// Also check with just the filename for patterns like "*.log"
			if matcher.MatchesPath(filepath.Base(absPath)) {
				return true
			}
		}
	}

	return false
}

// AddPatterns adds custom ignore patterns to the global matcher
func (g *GitignoreMatcher) AddPatterns(patterns []string) error {
	if len(patterns) == 0 {
		return nil
	}

	// Compile patterns
	g.global = ignore.CompileIgnoreLines(patterns...)

	return nil
}

// ShouldIgnoreDir checks if a directory should be ignored
// This is useful for early pruning of directory trees
func (g *GitignoreMatcher) ShouldIgnoreDir(dirPath string) bool {
	// Always ignore .git directories
	if filepath.Base(dirPath) == ".git" {
		return true
	}

	// Check gitignore patterns
	return g.Match(dirPath)
}
