package language

import (
	"bytes"
	"strings"
	"testing"
)

// TestIntegrationRealWorldScenarios tests real-world file processing scenarios
func TestIntegrationRealWorldScenarios(t *testing.T) {
	registry := GetRegistry()
	copyright := "Copyright 2025 Lean Security Co\nSPDX-License-Identifier: MIT"

	tests := []struct {
		name            string
		filename        string
		existingContent string
		wantHeader      string
		wantPreserved   int // number of preserved lines
	}{
		{
			name:     "Go file simple",
			filename: "main.go",
			existingContent: `package main

func main() {
	println("hello")
}`,
			wantHeader:    "// Copyright 2025 Lean Security Co\n// SPDX-License-Identifier: MIT\n",
			wantPreserved: 0,
		},
		{
			name:     "Python with shebang",
			filename: "script.py",
			existingContent: `#!/usr/bin/env python3
# -*- coding: utf-8 -*-
print("hello")`,
			wantHeader:    "# Copyright 2025 Lean Security Co\n# SPDX-License-Identifier: MIT\n",
			wantPreserved: 2,
		},
		{
			name:     "Shell script",
			filename: "deploy.sh",
			existingContent: `#!/bin/bash
set -e
echo "deploying..."`,
			wantHeader:    "# Copyright 2025 Lean Security Co\n# SPDX-License-Identifier: MIT\n",
			wantPreserved: 1,
		},
		{
			name:     "JavaScript with shebang (Node)",
			filename: "cli.js",
			existingContent: `#!/usr/bin/env node
console.log("hello");`,
			wantHeader: `/*
 * Copyright 2025 Lean Security Co
 * SPDX-License-Identifier: MIT
 */
`,
			wantPreserved: 1,
		},
		{
			name:     "Ruby with encoding",
			filename: "app.rb",
			existingContent: `#!/usr/bin/env ruby
# encoding: utf-8
puts "hello"`,
			wantHeader:    "# Copyright 2025 Lean Security Co\n# SPDX-License-Identifier: MIT\n",
			wantPreserved: 2,
		},
		{
			name:     "Java file",
			filename: "Main.java",
			existingContent: `public class Main {
    public static void main(String[] args) {
        System.out.println("hello");
    }
}`,
			wantHeader: `/*
 * Copyright 2025 Lean Security Co
 * SPDX-License-Identifier: MIT
 */
`,
			wantPreserved: 0,
		},
		{
			name:     "HTML with DOCTYPE",
			filename: "index.html",
			existingContent: `<!DOCTYPE html>
<html>
<head>
    <title>Test</title>
</head>
</html>`,
			wantHeader: `<!--
   Copyright 2025 Lean Security Co
   SPDX-License-Identifier: MIT
-->
`,
			wantPreserved: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lang := registry.GetLanguage(tt.filename)
			if lang == nil {
				t.Fatalf("No language handler for %s", tt.filename)
			}

			if !lang.CanHaveComments() {
				t.Skipf("Language %s doesn't support comments", lang.Name())
			}

			// Test header formatting
			header := lang.FormatHeader(copyright)
			if header != tt.wantHeader {
				t.Errorf("FormatHeader() =\n%q\nwant:\n%q", header, tt.wantHeader)
			}

			// Test preserved lines
			preserved := lang.PreserveLines([]byte(tt.existingContent))
			if len(preserved) != tt.wantPreserved {
				t.Errorf("PreserveLines() returned %d lines, want %d\nGot: %v",
					len(preserved), tt.wantPreserved, preserved)
			}

			// Test insert position
			insertPos := lang.InsertPosition([]byte(tt.existingContent))
			expectedPos := 0
			for _, line := range preserved {
				expectedPos += len(line)
			}
			if insertPos != expectedPos {
				t.Errorf("InsertPosition() = %d, want %d", insertPos, expectedPos)
			}

			// Simulate full copyright insertion
			content := []byte(tt.existingContent)
			var result bytes.Buffer

			// Write preserved lines
			for _, line := range preserved {
				result.WriteString(line)
			}

			// Write copyright header
			result.WriteString(header)

			// Add blank line
			result.WriteString("\n")

			// Write remaining content
			remaining := content[insertPos:]
			result.Write(remaining)

			// Verify the result contains all parts
			resultStr := result.String()
			if !strings.Contains(resultStr, "Copyright 2025 Lean Security Co") {
				t.Errorf("Result doesn't contain copyright text")
			}

			// Verify preserved content is still at the beginning
			for i, line := range preserved {
				if !strings.HasPrefix(strings.Split(resultStr, "\n")[i], strings.TrimSuffix(line, "\n")) {
					t.Errorf("Preserved line %d not at correct position", i)
				}
			}
		})
	}
}

// TestFileExtensionEdgeCases tests edge cases in file extension matching
func TestFileExtensionEdgeCases(t *testing.T) {
	registry := GetRegistry()

	tests := []struct {
		name     string
		filename string
		wantLang string
		wantNil  bool
	}{
		{"Rakefile no extension", "Rakefile", "Ruby", false},
		{"Dockerfile no extension", "Dockerfile", "Dockerfile", false},
		{"Makefile case sensitive", "Makefile", "Makefile", false},
		{"makefile lowercase", "makefile", "Makefile", false},
		{"Path with Makefile", "/path/to/project/Makefile", "Makefile", false},
		{"Hidden bash config", ".bashrc", "Shell", false},
		{"TypeScript definition", "types.d.ts", "JavaScript", false},
		{"React component", "App.jsx", "JavaScript", false},
		{"Multiple dots", "test.spec.ts", "JavaScript", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lang := registry.GetLanguage(tt.filename)

			if tt.wantNil {
				if lang != nil {
					t.Errorf("GetLanguage(%q) = %v, want nil", tt.filename, lang.Name())
				}
				return
			}

			if lang == nil {
				t.Errorf("GetLanguage(%q) = nil, want %q", tt.filename, tt.wantLang)
				return
			}

			if got := lang.Name(); got != tt.wantLang {
				t.Errorf("GetLanguage(%q).Name() = %q, want %q", tt.filename, got, tt.wantLang)
			}
		})
	}
}

// TestMultilineContent tests handling of complex multi-line content
func TestMultilineContent(t *testing.T) {
	python := NewPython()

	content := `#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Module docstring
"""

import sys

def main():
    print("hello")
`

	preserved := python.PreserveLines([]byte(content))
	if len(preserved) != 2 {
		t.Errorf("Expected 2 preserved lines, got %d", len(preserved))
	}

	if preserved[0] != "#!/usr/bin/env python3\n" {
		t.Errorf("First preserved line incorrect: %q", preserved[0])
	}

	if preserved[1] != "# -*- coding: utf-8 -*-\n" {
		t.Errorf("Second preserved line incorrect: %q", preserved[1])
	}
}

// BenchmarkGetLanguage benchmarks language lookup performance
func BenchmarkGetLanguage(b *testing.B) {
	registry := GetRegistry()
	files := []string{
		"main.go",
		"app.js",
		"script.py",
		"style.css",
		"index.html",
		"Dockerfile",
		"Makefile",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, file := range files {
			_ = registry.GetLanguage(file)
		}
	}
}

// BenchmarkFormatHeader benchmarks header formatting
func BenchmarkFormatHeader(b *testing.B) {
	lang := NewGoLang()
	copyright := "Copyright 2025 Lean Security Co\nSPDX-License-Identifier: MIT"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = lang.FormatHeader(copyright)
	}
}
