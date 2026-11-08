package license

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewManager(t *testing.T) {
	data := &LicenseData{
		Year:    "2024",
		Company: "SimplyCubed, GK",
		Author:  "",
	}

	manager := NewManager("mit", "", data)
	if manager == nil {
		t.Fatal("NewManager returned nil")
	}
	if manager.licenseType != "mit" {
		t.Errorf("expected license type 'mit', got '%s'", manager.licenseType)
	}
	if manager.data != data {
		t.Error("data not set correctly")
	}
}

func TestManager_FindLicense(t *testing.T) {
	tests := []struct {
		name          string
		createFiles   []string
		expectedFound bool
		expectedName  string
	}{
		{
			name:          "no license file",
			createFiles:   []string{"README.md"},
			expectedFound: false,
		},
		{
			name:          "LICENSE file exists",
			createFiles:   []string{"LICENSE"},
			expectedFound: true,
			expectedName:  "LICENSE",
		},
		{
			name:          "LICENSE.txt exists",
			createFiles:   []string{"LICENSE.txt"},
			expectedFound: true,
			expectedName:  "LICENSE.txt",
		},
		{
			name:          "LICENSE.md exists",
			createFiles:   []string{"LICENSE.md"},
			expectedFound: true,
			expectedName:  "LICENSE.md",
		},
		{
			name:          "COPYING exists",
			createFiles:   []string{"COPYING"},
			expectedFound: true,
			expectedName:  "COPYING",
		},
		{
			name:          "multiple license files - finds first",
			createFiles:   []string{"LICENSE", "LICENSE.txt"},
			expectedFound: true,
			expectedName:  "LICENSE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tmpDir := t.TempDir()

			// Create test files
			for _, filename := range tt.createFiles {
				path := filepath.Join(tmpDir, filename)
				if err := os.WriteFile(path, []byte("test content"), 0644); err != nil {
					t.Fatalf("failed to create test file: %v", err)
				}
			}

			// Test FindLicense
			manager := NewManager("mit", "", nil)
			path, found := manager.FindLicense(tmpDir)

			if found != tt.expectedFound {
				t.Errorf("expected found=%v, got %v", tt.expectedFound, found)
			}

			if tt.expectedFound {
				expectedPath := filepath.Join(tmpDir, tt.expectedName)
				if path != expectedPath {
					t.Errorf("expected path '%s', got '%s'", expectedPath, path)
				}
			}
		})
	}
}

func TestManager_Exists(t *testing.T) {
	tmpDir := t.TempDir()

	manager := NewManager("mit", "", nil)

	// Should not exist initially
	if manager.Exists(tmpDir) {
		t.Error("expected Exists to return false for empty directory")
	}

	// Create a LICENSE file
	licensePath := filepath.Join(tmpDir, "LICENSE")
	if err := os.WriteFile(licensePath, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Should exist now
	if !manager.Exists(tmpDir) {
		t.Error("expected Exists to return true after creating LICENSE")
	}
}

func TestManager_Generate_MIT(t *testing.T) {
	data := &LicenseData{
		Year:    "2024",
		Company: "SimplyCubed, GK",
		Author:  "",
	}

	manager := NewManager("mit", "", data)
	content, err := manager.Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Check that the content contains expected values
	if !strings.Contains(content, "MIT License") {
		t.Error("expected content to contain 'MIT License'")
	}
	if !strings.Contains(content, "2024") {
		t.Error("expected content to contain year '2024'")
	}
	if !strings.Contains(content, "SimplyCubed, GK") {
		t.Error("expected content to contain company 'SimplyCubed, GK'")
	}
	if !strings.Contains(content, "Permission is hereby granted") {
		t.Error("expected content to contain MIT license text")
	}
}

func TestManager_Generate_Apache(t *testing.T) {
	data := &LicenseData{
		Year:    "2024",
		Company: "SimplyCubed, GK",
	}

	manager := NewManager("apache-2.0", "", data)
	content, err := manager.Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if !strings.Contains(content, "Apache License") {
		t.Error("expected content to contain 'Apache License'")
	}
	if !strings.Contains(content, "2024") {
		t.Error("expected content to contain year '2024'")
	}
	if !strings.Contains(content, "SimplyCubed, GK") {
		t.Error("expected content to contain company 'SimplyCubed, GK'")
	}
}

func TestManager_Generate_BSD(t *testing.T) {
	data := &LicenseData{
		Year:    "2024",
		Company: "SimplyCubed, GK",
	}

	manager := NewManager("bsd-3-clause", "", data)
	content, err := manager.Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if !strings.Contains(content, "BSD 3-Clause License") {
		t.Error("expected content to contain 'BSD 3-Clause License'")
	}
	if !strings.Contains(content, "2024") {
		t.Error("expected content to contain year '2024'")
	}
	if !strings.Contains(content, "SimplyCubed, GK") {
		t.Error("expected content to contain company 'SimplyCubed, GK'")
	}
}

func TestManager_Generate_GPL(t *testing.T) {
	data := &LicenseData{
		Year:    "2024",
		Company: "SimplyCubed, GK",
	}

	manager := NewManager("gpl-3.0", "", data)
	content, err := manager.Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if !strings.Contains(content, "GNU General Public License") {
		t.Error("expected content to contain 'GNU General Public License'")
	}
	if !strings.Contains(content, "2024") {
		t.Error("expected content to contain year '2024'")
	}
	if !strings.Contains(content, "SimplyCubed, GK") {
		t.Error("expected content to contain company 'SimplyCubed, GK'")
	}
}

func TestManager_Generate_MPL(t *testing.T) {
	data := &LicenseData{
		Year:    "2024",
		Company: "SimplyCubed, GK",
	}

	manager := NewManager("mpl-2.0", "", data)
	content, err := manager.Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if !strings.Contains(content, "Mozilla Public License") {
		t.Error("expected content to contain 'Mozilla Public License'")
	}
	if !strings.Contains(content, "2024") {
		t.Error("expected content to contain year '2024'")
	}
	if !strings.Contains(content, "SimplyCubed, GK") {
		t.Error("expected content to contain company 'SimplyCubed, GK'")
	}
}

func TestManager_Generate_Custom(t *testing.T) {
	customText := "Custom License Text\nCopyright 2024"

	manager := NewManager("custom", customText, nil)
	content, err := manager.Generate()
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if content != customText {
		t.Errorf("expected content '%s', got '%s'", customText, content)
	}
}

func TestManager_Generate_CustomWithoutText(t *testing.T) {
	manager := NewManager("custom", "", nil)
	_, err := manager.Generate()
	if err == nil {
		t.Error("expected error for custom license without text")
	}
	if !strings.Contains(err.Error(), "custom license text is required") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestManager_Generate_InvalidLicense(t *testing.T) {
	data := &LicenseData{
		Year:    "2024",
		Company: "SimplyCubed, GK",
	}

	manager := NewManager("invalid-license", "", data)
	_, err := manager.Generate()
	if err == nil {
		t.Error("expected error for invalid license type")
	}
	if !strings.Contains(err.Error(), "license template not found") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestManager_Create(t *testing.T) {
	tmpDir := t.TempDir()

	data := &LicenseData{
		Year:    "2024",
		Company: "SimplyCubed, GK",
	}

	manager := NewManager("mit", "", data)
	if err := manager.Create(tmpDir); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Check that the file was created
	licensePath := filepath.Join(tmpDir, "LICENSE")
	if _, err := os.Stat(licensePath); err != nil {
		t.Errorf("LICENSE file was not created: %v", err)
	}

	// Read and verify content
	content, err := os.ReadFile(licensePath)
	if err != nil {
		t.Fatalf("failed to read LICENSE file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "MIT License") {
		t.Error("LICENSE file does not contain expected content")
	}
	if !strings.Contains(contentStr, "2024") {
		t.Error("LICENSE file does not contain year")
	}
	if !strings.Contains(contentStr, "SimplyCubed, GK") {
		t.Error("LICENSE file does not contain company")
	}
}

func TestManager_Create_FileAlreadyExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Create an existing LICENSE file
	licensePath := filepath.Join(tmpDir, "LICENSE")
	if err := os.WriteFile(licensePath, []byte("existing"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	data := &LicenseData{
		Year:    "2024",
		Company: "SimplyCubed, GK",
	}

	manager := NewManager("mit", "", data)
	err := manager.Create(tmpDir)
	if err == nil {
		t.Error("expected error when LICENSE file already exists")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestManager_Create_WithDifferentLicenseName(t *testing.T) {
	tmpDir := t.TempDir()

	// Create an existing LICENSE.txt file
	licensePath := filepath.Join(tmpDir, "LICENSE.txt")
	if err := os.WriteFile(licensePath, []byte("existing"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	data := &LicenseData{
		Year:    "2024",
		Company: "SimplyCubed, GK",
	}

	manager := NewManager("mit", "", data)
	err := manager.Create(tmpDir)
	if err == nil {
		t.Error("expected error when LICENSE.txt file already exists")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestManager_Create_AllLicenseTypes(t *testing.T) {
	licenseTypes := []string{"mit", "apache-2.0", "bsd-3-clause", "gpl-3.0", "mpl-2.0"}

	for _, licenseType := range licenseTypes {
		t.Run(licenseType, func(t *testing.T) {
			tmpDir := t.TempDir()

			data := &LicenseData{
				Year:    "2024",
				Company: "SimplyCubed, GK",
			}

			manager := NewManager(licenseType, "", data)
			if err := manager.Create(tmpDir); err != nil {
				t.Fatalf("Create failed for %s: %v", licenseType, err)
			}

			// Verify file exists
			licensePath := filepath.Join(tmpDir, "LICENSE")
			if _, err := os.Stat(licensePath); err != nil {
				t.Errorf("LICENSE file was not created for %s: %v", licenseType, err)
			}

			// Verify content is not empty
			content, err := os.ReadFile(licensePath)
			if err != nil {
				t.Fatalf("failed to read LICENSE file: %v", err)
			}
			if len(content) == 0 {
				t.Errorf("LICENSE file is empty for %s", licenseType)
			}
		})
	}
}
