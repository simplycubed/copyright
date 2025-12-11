package license

import (
	"strings"
	"testing"
)

func TestGetTemplate(t *testing.T) {
	tests := []struct {
		name         string
		licenseType  string
		wantError    bool
		checkContent []string // Strings that should be present in the template
	}{
		{
			name:         "MIT license",
			licenseType:  "mit",
			wantError:    false,
			checkContent: []string{"MIT License", "{{.Year}}", "{{.Company}}", "Permission is hereby granted"},
		},
		{
			name:         "Apache 2.0 license",
			licenseType:  "apache-2.0",
			wantError:    false,
			checkContent: []string{"Apache License", "{{.Year}}", "{{.Company}}", "http://www.apache.org/licenses/LICENSE-2.0"},
		},
		{
			name:         "BSD 3-Clause license",
			licenseType:  "bsd-3-clause",
			wantError:    false,
			checkContent: []string{"BSD 3-Clause License", "{{.Year}}", "{{.Company}}", "Redistribution and use"},
		},
		{
			name:         "GPL 3.0 license",
			licenseType:  "gpl-3.0",
			wantError:    false,
			checkContent: []string{"GNU General Public License", "{{.Year}}", "{{.Company}}", "free software"},
		},
		{
			name:         "MPL 2.0 license",
			licenseType:  "mpl-2.0",
			wantError:    false,
			checkContent: []string{"Mozilla Public License", "{{.Year}}", "{{.Company}}"},
		},
		{
			name:        "invalid license type",
			licenseType: "invalid-license",
			wantError:   true,
		},
		{
			name:        "empty license type",
			licenseType: "",
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := GetTemplate(tt.licenseType)

			if tt.wantError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if content == "" {
				t.Error("expected non-empty content")
			}

			// Check for expected content
			for _, expected := range tt.checkContent {
				if !strings.Contains(content, expected) {
					t.Errorf("expected content to contain '%s'", expected)
				}
			}
		})
	}
}

func TestListTemplates(t *testing.T) {
	templates := ListTemplates()

	if len(templates) == 0 {
		t.Fatal("expected at least one template")
	}

	// Check that templates are sorted
	for i := 1; i < len(templates); i++ {
		if templates[i-1] >= templates[i] {
			t.Errorf("templates are not sorted: %v", templates)
			break
		}
	}

	// Expected templates
	expectedTemplates := map[string]bool{
		"apache-2.0":   true,
		"bsd-3-clause": true,
		"gpl-3.0":      true,
		"mit":          true,
		"mpl-2.0":      true,
	}

	// Check that all expected templates are present
	foundTemplates := make(map[string]bool)
	for _, tmpl := range templates {
		foundTemplates[tmpl] = true
	}

	for expected := range expectedTemplates {
		if !foundTemplates[expected] {
			t.Errorf("expected template '%s' not found in list", expected)
		}
	}

	// Verify no .txt extension in template names
	for _, tmpl := range templates {
		if strings.HasSuffix(tmpl, ".txt") {
			t.Errorf("template name should not have .txt extension: %s", tmpl)
		}
	}
}

func TestTemplateFilesHavePlaceholders(t *testing.T) {
	templates := ListTemplates()

	for _, tmpl := range templates {
		t.Run(tmpl, func(t *testing.T) {
			content, err := GetTemplate(tmpl)
			if err != nil {
				t.Fatalf("failed to get template: %v", err)
			}

			// Check for placeholders
			hasYear := strings.Contains(content, "{{.Year}}")
			hasCompany := strings.Contains(content, "{{.Company}}")

			if !hasYear {
				t.Error("template missing {{.Year}} placeholder")
			}
			if !hasCompany {
				t.Error("template missing {{.Company}} placeholder")
			}

			// Verify content is not empty
			if len(content) < 50 {
				t.Error("template content seems too short")
			}
		})
	}
}

func TestGetTemplate_ErrorMessage(t *testing.T) {
	_, err := GetTemplate("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent template")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "license template not found") {
		t.Errorf("expected error message to contain 'license template not found', got: %s", errMsg)
	}
	if !strings.Contains(errMsg, "nonexistent") {
		t.Errorf("expected error message to contain license type 'nonexistent', got: %s", errMsg)
	}
}
