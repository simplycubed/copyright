// Package language provides language-specific copyright header formatting.
//
// This package supports automatic detection of programming languages and file types
// based on file extensions and specific filenames, and provides language-appropriate
// comment formatting for copyright headers.
//
// Supported Languages:
//   - Go (.go)
//   - JavaScript/TypeScript (.js, .ts, .jsx, .tsx, .mjs, .cjs)
//   - Java (.java)
//   - C# (.cs)
//   - Rust (.rs)
//   - C/C++ (.c, .h, .cpp, .hpp, .cc, .cxx)
//   - Python (.py, .pyw)
//   - Ruby (.rb, .rake, Rakefile, Gemfile)
//   - Shell (.sh, .bash, .zsh, .fish)
//   - YAML/TOML (.yaml, .yml, .toml)
//   - HTML/XML (.html, .htm, .xml, .svg)
//   - CSS (.css, .scss, .sass, .less)
//   - SQL (.sql)
//   - Markdown (.md, .markdown) - no comments supported
//   - JSON (.json) - no comments supported
//   - Protocol Buffers (.proto)
//   - Makefile (Makefile, makefile, GNUmakefile)
//   - Dockerfile (Dockerfile, Containerfile)
//   - Terraform (.tf, .tfvars, .hcl)
//
// Usage:
//
//	registry := language.GetRegistry()
//	lang := registry.GetLanguage("main.go")
//	if lang != nil && lang.CanHaveComments() {
//	    copyright := "Copyright 2025 Lean Security Co\nSPDX-License-Identifier: MIT"
//	    header := lang.FormatHeader(copyright)
//	    // header contains properly formatted comment
//	}
//
// Special Features:
//   - Preserves shebang lines (#!) for scripts
//   - Preserves Python encoding declarations
//   - Preserves Ruby encoding declarations
//   - Preserves XML declarations
//   - Case-insensitive filename matching for special files
package language

import (
	"bytes"
	"strings"
)

// CommentStyle defines the type of comment syntax a language uses
type CommentStyle int

const (
	// CommentStyleNone indicates the file type doesn't support comments
	CommentStyleNone CommentStyle = iota
	// CommentStyleLine indicates line comments (e.g., //, #)
	CommentStyleLine
	// CommentStyleBlock indicates block comments (e.g., /* */)
	CommentStyleBlock
	// CommentStyleHTML indicates HTML-style comments (e.g., <!-- -->)
	CommentStyleHTML
	// CommentStyleSQL indicates SQL-style comments (e.g., --)
	CommentStyleSQL
)

// Language defines the interface for language-specific copyright header handling
type Language interface {
	// Name returns the human-readable name of the language
	Name() string

	// Extensions returns the file extensions this language handles (e.g., ".go")
	Extensions() []string

	// Filenames returns specific filenames this language handles (e.g., "Makefile")
	Filenames() []string

	// CanHaveComments returns whether this file type supports comments
	CanHaveComments() bool

	// FormatHeader formats the copyright text as a comment block for this language
	FormatHeader(copyrightText string) string

	// PreserveLines returns lines that should be preserved before the copyright
	// (e.g., shebang lines, encoding declarations)
	PreserveLines(content []byte) []string

	// InsertPosition returns the byte offset where the copyright should be inserted
	InsertPosition(content []byte) int
}

// BaseLanguage provides common implementation for language handlers
type BaseLanguage struct {
	name         string
	extensions   []string
	filenames    []string
	commentStyle CommentStyle
	linePrefix   string // For line and hash comments
	blockStart   string // For block comments
	blockPrefix  string // For block comments (continuation)
	blockEnd     string // For block comments
}

// Name returns the language name
func (b *BaseLanguage) Name() string {
	return b.name
}

// Extensions returns the file extensions
func (b *BaseLanguage) Extensions() []string {
	return b.extensions
}

// Filenames returns specific filenames
func (b *BaseLanguage) Filenames() []string {
	return b.filenames
}

// CanHaveComments returns whether the language supports comments
func (b *BaseLanguage) CanHaveComments() bool {
	return b.commentStyle != CommentStyleNone
}

// FormatHeader formats the copyright text as a comment
func (b *BaseLanguage) FormatHeader(copyrightText string) string {
	lines := strings.Split(strings.TrimSpace(copyrightText), "\n")

	switch b.commentStyle {
	case CommentStyleLine:
		// Line comments like // or #
		var formatted []string
		for _, line := range lines {
			if line == "" {
				formatted = append(formatted, b.linePrefix)
			} else {
				formatted = append(formatted, b.linePrefix+" "+line)
			}
		}
		return strings.Join(formatted, "\n") + "\n"

	case CommentStyleBlock, CommentStyleHTML:
		// Block comments like /* */ or <!-- -->
		var formatted []string
		formatted = append(formatted, b.blockStart)
		for _, line := range lines {
			if line == "" {
				formatted = append(formatted, b.blockPrefix)
			} else {
				formatted = append(formatted, b.blockPrefix+" "+line)
			}
		}
		formatted = append(formatted, b.blockEnd)
		return strings.Join(formatted, "\n") + "\n"

	case CommentStyleSQL:
		// SQL-style comments with --
		var formatted []string
		for _, line := range lines {
			if line == "" {
				formatted = append(formatted, b.linePrefix)
			} else {
				formatted = append(formatted, b.linePrefix+" "+line)
			}
		}
		return strings.Join(formatted, "\n") + "\n"

	case CommentStyleNone:
		// No comment support
		return ""

	default:
		return ""
	}
}

// PreserveLines returns lines to preserve (default: none)
func (b *BaseLanguage) PreserveLines(content []byte) []string {
	return nil
}

// InsertPosition returns the position to insert the copyright
// Note: This method cannot be implemented in BaseLanguage because it needs to call
// the overridden PreserveLines method from the specific language implementation.
// Each language that overrides PreserveLines must also override InsertPosition.
func (b *BaseLanguage) InsertPosition(content []byte) int {
	// For languages without special preserved lines, insert at the beginning
	return 0
}

// hasShebang checks if the content starts with a shebang line
func hasShebang(content []byte) bool {
	return bytes.HasPrefix(content, []byte("#!"))
}

// getShebangLine extracts the shebang line if present
func getShebangLine(content []byte) string {
	if !hasShebang(content) {
		return ""
	}

	idx := bytes.IndexByte(content, '\n')
	if idx == -1 {
		return string(content)
	}
	return string(content[:idx+1])
}

// hasPythonEncoding checks if a line is a Python encoding declaration
func hasPythonEncoding(line string) bool {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "#") {
		return false
	}

	// Python encoding: # -*- coding: utf-8 -*-
	// or: # coding: utf-8
	// or: # coding=utf-8
	return strings.Contains(line, "coding:") ||
		strings.Contains(line, "coding=") ||
		strings.Contains(line, "-*- coding")
}

// calculateInsertPosition calculates byte position after preserved lines
func calculateInsertPosition(preserved []string) int {
	position := 0
	for _, line := range preserved {
		position += len(line)
	}
	return position
}
