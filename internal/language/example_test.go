package language_test

import (
	"fmt"

	"github.com/charlesgreen/copyright/internal/language"
)

// Example demonstrates how to use the language registry
func Example() {
	registry := language.GetRegistry()

	// Get language for a Go file
	lang := registry.GetLanguage("main.go")
	if lang != nil {
		fmt.Printf("Language: %s\n", lang.Name())
		fmt.Printf("Can have comments: %v\n", lang.CanHaveComments())

		// Format a copyright header
		copyright := "Copyright 2025 Lean Security Co\nSPDX-License-Identifier: MIT"
		header := lang.FormatHeader(copyright)
		fmt.Printf("Header:\n%s", header)
	}

	// Output:
	// Language: Go
	// Can have comments: true
	// Header:
	// // Copyright 2025 Lean Security Co
	// // SPDX-License-Identifier: MIT
}

// ExampleRegistry_GetLanguage demonstrates language detection
func ExampleRegistry_GetLanguage() {
	registry := language.GetRegistry()

	// Test various file types
	files := []string{
		"main.go",
		"script.py",
		"app.ts",
		"README.md",
		"Dockerfile",
	}

	for _, file := range files {
		lang := registry.GetLanguage(file)
		if lang != nil {
			fmt.Printf("%s -> %s\n", file, lang.Name())
		}
	}

	// Output:
	// main.go -> Go
	// script.py -> Python
	// app.ts -> JavaScript
	// README.md -> Markdown
	// Dockerfile -> Dockerfile
}

// ExamplePython_PreserveLines demonstrates shebang preservation
func ExamplePython_PreserveLines() {
	python := language.NewPython()

	content := []byte("#!/usr/bin/env python3\n# coding: utf-8\nprint('hello')")
	preserved := python.PreserveLines(content)

	fmt.Printf("Preserved lines: %d\n", len(preserved))
	for i, line := range preserved {
		fmt.Printf("Line %d: %s", i+1, line)
	}

	// Output:
	// Preserved lines: 2
	// Line 1: #!/usr/bin/env python3
	// Line 2: # coding: utf-8
}

// ExampleLanguage_FormatHeader demonstrates different comment styles
func ExampleLanguage_FormatHeader() {
	copyright := "Copyright 2025 Lean Security Co\nSPDX-License-Identifier: MIT"

	// Go - line comments
	goLang := language.NewGoLang()
	fmt.Println("Go:")
	fmt.Print(goLang.FormatHeader(copyright))

	// Java - block comments
	fmt.Println("\nJava:")
	java := language.NewJava()
	fmt.Print(java.FormatHeader(copyright))

	// Python - hash comments
	fmt.Println("\nPython:")
	python := language.NewPython()
	fmt.Print(python.FormatHeader(copyright))

	// SQL - SQL-style comments
	fmt.Println("\nSQL:")
	sql := language.NewSQL()
	fmt.Print(sql.FormatHeader(copyright))

	// Output:
	// Go:
	// // Copyright 2025 Lean Security Co
	// // SPDX-License-Identifier: MIT
	//
	// Java:
	// /*
	//  * Copyright 2025 Lean Security Co
	//  * SPDX-License-Identifier: MIT
	//  */
	//
	// Python:
	// # Copyright 2025 Lean Security Co
	// # SPDX-License-Identifier: MIT
	//
	// SQL:
	// -- Copyright 2025 Lean Security Co
	// -- SPDX-License-Identifier: MIT
}
