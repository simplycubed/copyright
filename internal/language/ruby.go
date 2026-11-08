package language

import (
	"bytes"
	"strings"
)

// Ruby handles Ruby source files
type Ruby struct {
	BaseLanguage
}

// NewRuby creates a new Ruby language handler
func NewRuby() *Ruby {
	return &Ruby{
		BaseLanguage: BaseLanguage{
			name:         "Ruby",
			extensions:   []string{".rb", ".rake", ".gemspec"},
			filenames:    []string{"Rakefile", "Gemfile"},
			commentStyle: CommentStyleLine,
			linePrefix:   "#",
		},
	}
}

// PreserveLines handles shebang and encoding declarations for Ruby
func (r *Ruby) PreserveLines(content []byte) []string {
	var preserve []string

	// Check for shebang
	if shebang := getShebangLine(content); shebang != "" {
		preserve = append(preserve, shebang)
	}

	// Check for encoding declaration (must be on first or second line)
	// Ruby encoding: # encoding: utf-8 or # coding: utf-8
	lines := bytes.SplitN(content, []byte("\n"), 3)

	// Check first line if not shebang
	if len(preserve) == 0 && len(lines) > 0 {
		line := strings.TrimSpace(string(lines[0]))
		if hasRubyEncoding(line) {
			preserve = append(preserve, string(lines[0])+"\n")
		}
	}

	// Check second line for encoding
	if len(lines) > 1 {
		line := strings.TrimSpace(string(lines[1]))
		if hasRubyEncoding(line) {
			preserve = append(preserve, string(lines[1])+"\n")
		}
	}

	return preserve
}

// InsertPosition returns the byte offset where copyright should be inserted
func (r *Ruby) InsertPosition(content []byte) int {
	return calculateInsertPosition(r.PreserveLines(content))
}

// hasRubyEncoding checks if a line is a Ruby encoding declaration
func hasRubyEncoding(line string) bool {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "#") {
		return false
	}

	// Ruby encoding: # encoding: utf-8
	// or: # coding: utf-8
	return strings.Contains(line, "encoding:") ||
		strings.Contains(line, "coding:")
}
