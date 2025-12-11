// Package language provides comprehensive language-specific copyright header support.
//
// # Overview
//
// This package implements automatic language detection and copyright header formatting
// for 19+ programming languages and file types. It intelligently handles:
//   - Multiple comment styles (line, block, HTML, SQL)
//   - Shebang preservation for scripts
//   - Encoding declarations (Python, Ruby)
//   - XML/DOCTYPE declarations
//   - Case-insensitive filename matching
//
// # Architecture
//
// The package consists of:
//   - Language interface: Defines the contract for all language handlers
//   - Registry: Maps file extensions and names to language handlers
//   - BaseLanguage: Provides common implementation for most languages
//   - Individual handlers: One per language/file type
//
// # Comment Styles
//
// Line Comments (//, #):
//   - Go, Rust, Protocol Buffers
//   - Python, Ruby, Shell
//   - YAML, TOML, Terraform
//   - Makefile, Dockerfile
//
// Block Comments (/* */):
//   - JavaScript/TypeScript, Java, C#, C/C++
//   - CSS, SCSS, SASS, LESS
//
// HTML Comments (<!-- -->):
//   - HTML, XML, SVG
//
// SQL Comments (--):
//   - SQL
//
// No Comments:
//   - Markdown, JSON
//
// # Example Output
//
// Go, Rust, Protocol Buffers:
//
//	// Copyright 2025 Lean Security Co
//	// SPDX-License-Identifier: MIT
//
// Java, C/C++, JavaScript, CSS:
//
//	/*
//	 * Copyright 2025 Lean Security Co
//	 * SPDX-License-Identifier: MIT
//	 */
//
// Python, Ruby, Shell, YAML:
//
//	# Copyright 2025 Lean Security Co
//	# SPDX-License-Identifier: MIT
//
// HTML, XML:
//
//	<!--
//	  Copyright 2025 Lean Security Co
//	  SPDX-License-Identifier: MIT
//	-->
//
// SQL:
//
//	-- Copyright 2025 Lean Security Co
//	-- SPDX-License-Identifier: MIT
//
// # Usage Examples
//
// Basic usage:
//
//	registry := language.GetRegistry()
//	lang := registry.GetLanguage("main.go")
//	if lang != nil && lang.CanHaveComments() {
//	    copyright := "Copyright 2025 Lean Security Co\nSPDX-License-Identifier: MIT"
//	    header := lang.FormatHeader(copyright)
//	    fmt.Print(header)
//	}
//
// With shebang preservation:
//
//	python := language.NewPython()
//	content := []byte("#!/usr/bin/env python3\nprint('hello')")
//	preserved := python.PreserveLines(content)
//	insertPos := python.InsertPosition(content)
//	// Insert copyright at insertPos, preserving shebang
//
// # Extension
//
// To add a new language:
//
//  1. Create a new file (e.g., mylang.go)
//  2. Define a struct embedding BaseLanguage
//  3. Implement a NewMyLang() constructor
//  4. Register in registerAllLanguages()
//  5. Override PreserveLines() if needed for special handling
//
// Example:
//
//	type MyLang struct {
//	    BaseLanguage
//	}
//
//	func NewMyLang() *MyLang {
//	    return &MyLang{
//	        BaseLanguage: BaseLanguage{
//	            name:         "MyLang",
//	            extensions:   []string{".ml"},
//	            commentStyle: CommentStyleLine,
//	            linePrefix:   "//",
//	        },
//	    }
//	}
package language
