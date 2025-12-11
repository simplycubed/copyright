package language

import (
	"bytes"
	"strings"
)

// HTML handles HTML and XML files
type HTML struct {
	BaseLanguage
}

// NewHTML creates a new HTML/XML language handler
func NewHTML() *HTML {
	return &HTML{
		BaseLanguage: BaseLanguage{
			name:         "HTML/XML",
			extensions:   []string{".html", ".htm", ".xml", ".xhtml", ".svg"},
			filenames:    []string{},
			commentStyle: CommentStyleHTML,
			blockStart:   "<!--",
			blockPrefix:  "  ",
			blockEnd:     "-->",
		},
	}
}

// PreserveLines handles XML declaration for XML files
func (h *HTML) PreserveLines(content []byte) []string {
	var preserve []string

	// Check for XML declaration (<?xml version="1.0" ?>)
	if bytes.HasPrefix(bytes.TrimSpace(content), []byte("<?xml")) {
		idx := bytes.Index(content, []byte("?>"))
		if idx != -1 {
			// Find the newline after the declaration
			endIdx := idx + 2
			newlineIdx := bytes.IndexByte(content[endIdx:], '\n')
			if newlineIdx != -1 {
				preserve = append(preserve, string(content[:endIdx+newlineIdx+1]))
			} else {
				preserve = append(preserve, string(content[:endIdx])+"\n")
			}
		}
	}

	// Check for DOCTYPE declaration
	trimmed := bytes.TrimSpace(content)
	if bytes.HasPrefix(trimmed, []byte("<!DOCTYPE")) || bytes.HasPrefix(trimmed, []byte("<!doctype")) {
		// Find the position in original content
		doctypeStart := bytes.Index(content, []byte("<!DOCTYPE"))
		if doctypeStart == -1 {
			doctypeStart = bytes.Index(content, []byte("<!doctype"))
		}

		if doctypeStart != -1 {
			remaining := content[doctypeStart:]
			idx := bytes.IndexByte(remaining, '>')
			if idx != -1 {
				// Find the newline after the declaration
				endIdx := doctypeStart + idx + 1
				newlineIdx := bytes.IndexByte(content[endIdx:], '\n')
				if newlineIdx != -1 {
					doctype := string(content[:endIdx+newlineIdx+1])
					// Avoid duplicates
					if len(preserve) == 0 || !strings.Contains(preserve[0], "DOCTYPE") {
						preserve = append(preserve, doctype)
					}
				}
			}
		}
	}

	return preserve
}

// InsertPosition returns the byte offset where copyright should be inserted
func (h *HTML) InsertPosition(content []byte) int {
	return calculateInsertPosition(h.PreserveLines(content))
}
