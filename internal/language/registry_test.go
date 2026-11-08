package language

import (
	"testing"
)

func TestGetLanguage(t *testing.T) {
	registry := GetRegistry()

	tests := []struct {
		name     string
		filename string
		wantLang string
		wantNil  bool
	}{
		// Go
		{"Go file", "main.go", "Go", false},
		{"Go file with path", "/path/to/main.go", "Go", false},

		// JavaScript/TypeScript
		{"JavaScript file", "app.js", "JavaScript", false},
		{"TypeScript file", "app.ts", "JavaScript", false},
		{"JSX file", "component.jsx", "JavaScript", false},
		{"TSX file", "component.tsx", "JavaScript", false},
		{"MJS file", "module.mjs", "JavaScript", false},
		{"CJS file", "module.cjs", "JavaScript", false},

		// Java
		{"Java file", "Main.java", "Java", false},

		// C#
		{"C# file", "Program.cs", "C#", false},

		// Rust
		{"Rust file", "main.rs", "Rust", false},

		// C/C++
		{"C file", "main.c", "C/C++", false},
		{"C header", "main.h", "C/C++", false},
		{"C++ file", "main.cpp", "C/C++", false},
		{"C++ header", "main.hpp", "C/C++", false},

		// Python
		{"Python file", "script.py", "Python", false},
		{"Python w file", "script.pyw", "Python", false},

		// Ruby
		{"Ruby file", "script.rb", "Ruby", false},
		{"Rake file", "task.rake", "Ruby", false},
		{"Rakefile", "Rakefile", "Ruby", false},
		{"Gemfile", "Gemfile", "Ruby", false},

		// Shell
		{"Shell file", "script.sh", "Shell", false},
		{"Bash file", "script.bash", "Shell", false},
		{"Zsh file", "script.zsh", "Shell", false},
		{"Fish file", "script.fish", "Shell", false},

		// Config files
		{"YAML file", "config.yaml", "YAML/TOML", false},
		{"YML file", "config.yml", "YAML/TOML", false},
		{"TOML file", "config.toml", "YAML/TOML", false},

		// HTML/XML
		{"HTML file", "index.html", "HTML/XML", false},
		{"XML file", "config.xml", "HTML/XML", false},
		{"SVG file", "icon.svg", "HTML/XML", false},

		// CSS
		{"CSS file", "style.css", "CSS", false},
		{"SCSS file", "style.scss", "CSS", false},
		{"SASS file", "style.sass", "CSS", false},

		// SQL
		{"SQL file", "schema.sql", "SQL", false},

		// Markdown
		{"Markdown file", "README.md", "Markdown", false},
		{"Markdown file alt", "doc.markdown", "Markdown", false},

		// JSON
		{"JSON file", "package.json", "JSON", false},

		// Protocol Buffers
		{"Proto file", "api.proto", "Protocol Buffers", false},

		// Makefile
		{"Makefile", "Makefile", "Makefile", false},
		{"makefile", "makefile", "Makefile", false},
		{"GNUmakefile", "GNUmakefile", "Makefile", false},

		// Dockerfile
		{"Dockerfile", "Dockerfile", "Dockerfile", false},
		{"Containerfile", "Containerfile", "Dockerfile", false},

		// Terraform
		{"Terraform file", "main.tf", "Terraform", false},
		{"Terraform vars", "variables.tfvars", "Terraform", false},

		// Unknown
		{"Unknown extension", "file.unknown", "", true},
		{"No extension", "file", "", true},
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

func TestFormatHeader(t *testing.T) {
	tests := []struct {
		name      string
		lang      Language
		copyright string
		want      string
	}{
		{
			name:      "Go line comments",
			lang:      &GoLang{BaseLanguage: BaseLanguage{commentStyle: CommentStyleLine, linePrefix: "//"}},
			copyright: "Copyright 2025 SimplyCubed, GK\nSPDX-License-Identifier: MIT",
			want:      "// Copyright 2025 SimplyCubed, GK\n// SPDX-License-Identifier: MIT\n",
		},
		{
			name:      "Python hash comments",
			lang:      &Python{BaseLanguage: BaseLanguage{commentStyle: CommentStyleLine, linePrefix: "#"}},
			copyright: "Copyright 2025 SimplyCubed, GK\nSPDX-License-Identifier: MIT",
			want:      "# Copyright 2025 SimplyCubed, GK\n# SPDX-License-Identifier: MIT\n",
		},
		{
			name: "Java block comments",
			lang: &Java{BaseLanguage: BaseLanguage{
				commentStyle: CommentStyleBlock,
				blockStart:   "/*",
				blockPrefix:  " *",
				blockEnd:     " */",
			}},
			copyright: "Copyright 2025 SimplyCubed, GK\nSPDX-License-Identifier: MIT",
			want:      "/*\n * Copyright 2025 SimplyCubed, GK\n * SPDX-License-Identifier: MIT\n */\n",
		},
		{
			name: "HTML comments",
			lang: &HTML{BaseLanguage: BaseLanguage{
				commentStyle: CommentStyleHTML,
				blockStart:   "<!--",
				blockPrefix:  "  ",
				blockEnd:     "-->",
			}},
			copyright: "Copyright 2025 SimplyCubed, GK\nSPDX-License-Identifier: MIT",
			want:      "<!--\n   Copyright 2025 SimplyCubed, GK\n   SPDX-License-Identifier: MIT\n-->\n",
		},
		{
			name:      "SQL comments",
			lang:      &SQL{BaseLanguage: BaseLanguage{commentStyle: CommentStyleSQL, linePrefix: "--"}},
			copyright: "Copyright 2025 SimplyCubed, GK\nSPDX-License-Identifier: MIT",
			want:      "-- Copyright 2025 SimplyCubed, GK\n-- SPDX-License-Identifier: MIT\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.lang.FormatHeader(tt.copyright)
			if got != tt.want {
				t.Errorf("FormatHeader() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPreserveLines(t *testing.T) {
	tests := []struct {
		name    string
		lang    Language
		content string
		want    []string
	}{
		{
			name:    "Python with shebang",
			lang:    NewPython(),
			content: "#!/usr/bin/env python3\nprint('hello')",
			want:    []string{"#!/usr/bin/env python3\n"},
		},
		{
			name:    "Python with encoding",
			lang:    NewPython(),
			content: "# -*- coding: utf-8 -*-\nprint('hello')",
			want:    []string{"# -*- coding: utf-8 -*-\n"},
		},
		{
			name:    "Python with shebang and encoding",
			lang:    NewPython(),
			content: "#!/usr/bin/env python3\n# coding: utf-8\nprint('hello')",
			want:    []string{"#!/usr/bin/env python3\n", "# coding: utf-8\n"},
		},
		{
			name:    "Shell with shebang",
			lang:    NewShell(),
			content: "#!/bin/bash\necho 'hello'",
			want:    []string{"#!/bin/bash\n"},
		},
		{
			name:    "Ruby with shebang",
			lang:    NewRuby(),
			content: "#!/usr/bin/env ruby\nputs 'hello'",
			want:    []string{"#!/usr/bin/env ruby\n"},
		},
		{
			name:    "Go no preserve",
			lang:    NewGoLang(),
			content: "package main\n",
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.lang.PreserveLines([]byte(tt.content))
			if len(got) != len(tt.want) {
				t.Errorf("PreserveLines() returned %d lines, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("PreserveLines()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestCanHaveComments(t *testing.T) {
	tests := []struct {
		name string
		lang Language
		want bool
	}{
		{"Go can have comments", NewGoLang(), true},
		{"Python can have comments", NewPython(), true},
		{"JavaScript can have comments", NewJavaScript(), true},
		{"Markdown cannot have comments", NewMarkdown(), false},
		{"JSON cannot have comments", NewJSON(), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.lang.CanHaveComments(); got != tt.want {
				t.Errorf("CanHaveComments() = %v, want %v", got, tt.want)
			}
		})
	}
}
