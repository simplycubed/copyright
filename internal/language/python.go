package language

import (
	"bytes"
	"strings"
)

// Python handles Python source files
type Python struct {
	BaseLanguage
}

// NewPython creates a new Python language handler
func NewPython() *Python {
	return &Python{
		BaseLanguage: BaseLanguage{
			name:         "Python",
			extensions:   []string{".py", ".pyw"},
			filenames:    []string{},
			commentStyle: CommentStyleLine,
			linePrefix:   "#",
		},
	}
}

// PreserveLines handles shebang and encoding declarations for Python
func (p *Python) PreserveLines(content []byte) []string {
	var preserve []string

	// Check for shebang
	if shebang := getShebangLine(content); shebang != "" {
		preserve = append(preserve, shebang)
	}

	// Check for encoding declaration (must be on first or second line)
	lines := bytes.SplitN(content, []byte("\n"), 3)

	// Check first line if not shebang
	if len(preserve) == 0 && len(lines) > 0 {
		line := strings.TrimSpace(string(lines[0]))
		if hasPythonEncoding(line) {
			preserve = append(preserve, string(lines[0])+"\n")
		}
	}

	// Check second line for encoding
	if len(lines) > 1 {
		line := strings.TrimSpace(string(lines[1]))
		if hasPythonEncoding(line) {
			preserve = append(preserve, string(lines[1])+"\n")
		}
	}

	return preserve
}

// InsertPosition returns the byte offset where copyright should be inserted
func (p *Python) InsertPosition(content []byte) int {
	return calculateInsertPosition(p.PreserveLines(content))
}
