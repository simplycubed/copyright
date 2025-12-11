package license

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

// LicenseData for template rendering
type LicenseData struct {
	Year    string
	Company string
	Author  string
}

// Manager handles LICENSE file operations
type Manager struct {
	licenseType string
	customText  string
	data        *LicenseData
}

// NewManager creates a new license manager
func NewManager(licenseType, customText string, data *LicenseData) *Manager {
	return &Manager{
		licenseType: licenseType,
		customText:  customText,
		data:        data,
	}
}

// commonLicenseNames lists common LICENSE file names to check
var commonLicenseNames = []string{
	"LICENSE",
	"LICENSE.txt",
	"LICENSE.md",
	"COPYING",
	"COPYING.txt",
	"License",
	"license",
	"license.txt",
}

// Exists checks if a LICENSE file exists in the directory
func (m *Manager) Exists(dir string) bool {
	_, exists := m.FindLicense(dir)
	return exists
}

// FindLicense searches for a LICENSE file in the directory
// Returns the path to the LICENSE file and whether it was found
func (m *Manager) FindLicense(dir string) (string, bool) {
	for _, name := range commonLicenseNames {
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); err == nil {
			return path, true
		}
	}
	return "", false
}

// Generate generates the license content based on the license type
func (m *Manager) Generate() (string, error) {
	// Handle custom license type
	if m.licenseType == "custom" {
		if m.customText == "" {
			return "", fmt.Errorf("custom license text is required for custom license type")
		}
		return m.customText, nil
	}

	// Get the template
	templateContent, err := GetTemplate(m.licenseType)
	if err != nil {
		return "", err
	}

	// Parse and execute the template
	tmpl, err := template.New("license").Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse license template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, m.data); err != nil {
		return "", fmt.Errorf("failed to execute license template: %w", err)
	}

	return buf.String(), nil
}

// Create creates a LICENSE file in the specified directory
func (m *Manager) Create(dir string) error {
	// Check if a LICENSE file already exists
	if existingPath, exists := m.FindLicense(dir); exists {
		return fmt.Errorf("license file already exists: %s", existingPath)
	}

	// Generate the license content
	content, err := m.Generate()
	if err != nil {
		return err
	}

	// Create the LICENSE file
	licensePath := filepath.Join(dir, "LICENSE")
	if err := os.WriteFile(licensePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write LICENSE file: %w", err)
	}

	return nil
}
