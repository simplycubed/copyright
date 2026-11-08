package license

import (
	"embed"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

//go:embed templates/*.txt
var licenseTemplates embed.FS

// GetTemplate returns template content for a license type
func GetTemplate(licenseType string) (string, error) {
	filename := fmt.Sprintf("templates/%s.txt", licenseType)
	content, err := licenseTemplates.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("license template not found: %s", licenseType)
	}
	return string(content), nil
}

// ListTemplates returns available license types
func ListTemplates() []string {
	entries, err := licenseTemplates.ReadDir("templates")
	if err != nil {
		return []string{}
	}

	var templates []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(name, ".txt") {
			// Remove .txt extension
			templates = append(templates, strings.TrimSuffix(name, filepath.Ext(name)))
		}
	}

	sort.Strings(templates)
	return templates
}
