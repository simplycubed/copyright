package config

import (
	"fmt"
	"strings"
)

// Validate performs validation on the configuration.
func (c *Config) Validate() error {
	if err := c.validateCopyright(); err != nil {
		return fmt.Errorf("copyright config: %w", err)
	}

	if err := c.validateLicense(); err != nil {
		return fmt.Errorf("license config: %w", err)
	}

	if err := c.validateScan(); err != nil {
		return fmt.Errorf("scan config: %w", err)
	}

	if err := c.validateDetection(); err != nil {
		return fmt.Errorf("detection config: %w", err)
	}

	if err := c.validateOutput(); err != nil {
		return fmt.Errorf("output config: %w", err)
	}

	return nil
}

// validateCopyright validates the copyright configuration.
func (c *Config) validateCopyright() error {
	if c.Copyright.Template == "" {
		return fmt.Errorf("template cannot be empty")
	}

	// Validate year format
	if c.Copyright.Year != "" && c.Copyright.Year != "auto" {
		// Check if it's a valid year format (either single year or range)
		if !isValidYearFormat(c.Copyright.Year) {
			return fmt.Errorf("invalid year format: %s (expected format: YYYY, YYYY-YYYY, or 'auto')", c.Copyright.Year)
		}
	}

	return nil
}

// validateLicense validates the license configuration.
func (c *Config) validateLicense() error {
	if c.License.Type == "" {
		// License type is optional
		return nil
	}

	// Normalize license type to lowercase
	licenseType := strings.ToLower(c.License.Type)

	// Check if it's a valid license type
	if !isValidLicenseType(licenseType) {
		return fmt.Errorf("invalid license type: %s (valid types: %s)", c.License.Type, strings.Join(ValidLicenseTypes, ", "))
	}

	// If custom license, ensure custom text is provided
	if licenseType == "custom" && c.License.CustomText == "" {
		return fmt.Errorf("custom_text must be provided when license type is 'custom'")
	}

	return nil
}

// validateScan validates the scan configuration.
func (c *Config) validateScan() error {
	// Validate exclude patterns
	for _, pattern := range c.Scan.ExcludePatterns {
		if pattern == "" {
			return fmt.Errorf("exclude_patterns cannot contain empty strings")
		}
	}

	// Validate include patterns
	for _, pattern := range c.Scan.IncludePatterns {
		if pattern == "" {
			return fmt.Errorf("include_patterns cannot contain empty strings")
		}
	}

	// Validate excluded extensions format
	for _, ext := range c.Scan.ExcludedExtensions {
		if ext == "" {
			return fmt.Errorf("excluded_extensions cannot contain empty strings")
		}
		if !strings.HasPrefix(ext, ".") {
			return fmt.Errorf("excluded extension must start with '.': %s", ext)
		}
	}

	return nil
}

// validateDetection validates the detection configuration.
func (c *Config) validateDetection() error {
	if c.Detection.LinesToScan <= 0 {
		return fmt.Errorf("lines_to_scan must be positive, got: %d", c.Detection.LinesToScan)
	}

	if c.Detection.LinesToScan > 1000 {
		return fmt.Errorf("lines_to_scan is too large (max: 1000), got: %d", c.Detection.LinesToScan)
	}

	if len(c.Detection.Keywords) == 0 {
		return fmt.Errorf("keywords list cannot be empty")
	}

	// Validate keywords are not empty
	for i, keyword := range c.Detection.Keywords {
		if strings.TrimSpace(keyword) == "" {
			return fmt.Errorf("keyword at index %d cannot be empty or whitespace", i)
		}
	}

	return nil
}

// validateOutput validates the output configuration.
func (c *Config) validateOutput() error {
	if c.Output.Format == "" {
		return fmt.Errorf("output format cannot be empty")
	}

	// Validate output format
	if !isValidOutputFormat(c.Output.Format) {
		return fmt.Errorf("invalid output format: %s (valid formats: %s)", c.Output.Format, strings.Join(ValidOutputFormats, ", "))
	}

	// Verbose and quiet are mutually exclusive
	if c.Output.Verbose && c.Output.Quiet {
		return fmt.Errorf("verbose and quiet options are mutually exclusive")
	}

	return nil
}

// isValidYearFormat checks if the year string is in a valid format.
func isValidYearFormat(year string) bool {
	// Allow "auto"
	if year == "auto" {
		return true
	}

	// Allow single year (YYYY)
	if len(year) == 4 {
		for _, c := range year {
			if c < '0' || c > '9' {
				return false
			}
		}
		return true
	}

	// Allow year range (YYYY-YYYY)
	if len(year) == 9 && year[4] == '-' {
		startYear := year[0:4]
		endYear := year[5:9]

		for _, c := range startYear {
			if c < '0' || c > '9' {
				return false
			}
		}
		for _, c := range endYear {
			if c < '0' || c > '9' {
				return false
			}
		}
		return true
	}

	return false
}

// isValidLicenseType checks if the license type is valid.
func isValidLicenseType(licenseType string) bool {
	normalized := strings.ToLower(licenseType)
	for _, valid := range ValidLicenseTypes {
		if normalized == valid {
			return true
		}
	}
	return false
}

// isValidOutputFormat checks if the output format is valid.
func isValidOutputFormat(format string) bool {
	normalized := strings.ToLower(format)
	for _, valid := range ValidOutputFormats {
		if normalized == valid {
			return true
		}
	}
	return false
}
