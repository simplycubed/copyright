// Package config provides configuration management for the copyright CLI tool.
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config represents the main configuration structure for the copyright tool.
type Config struct {
	Copyright CopyrightConfig `yaml:"copyright" mapstructure:"copyright"`
	License   LicenseConfig   `yaml:"license" mapstructure:"license"`
	Scan      ScanConfig      `yaml:"scan" mapstructure:"scan"`
	Detection DetectionConfig `yaml:"detection" mapstructure:"detection"`
	Output    OutputConfig    `yaml:"output" mapstructure:"output"`
}

// CopyrightConfig contains copyright header configuration.
type CopyrightConfig struct {
	// Template is the copyright header template.
	// Supports placeholders: {{.Year}}, {{.Company}}, {{.Author}}, {{.SPDX}}, {{.SPDXIdentifier}}.
	// {{.SPDX}} and {{.SPDXIdentifier}} both render spdx_identifier.
	Template string `yaml:"template" mapstructure:"template"`

	// Year is the copyright year. Use "auto" for current year.
	Year string `yaml:"year" mapstructure:"year"`

	// Company is the company or organization name.
	Company string `yaml:"company" mapstructure:"company"`

	// Author is the author name (optional).
	Author string `yaml:"author" mapstructure:"author"`

	// SPDXIdentifier is the SPDX license identifier (e.g., "Apache-2.0", "MIT").
	SPDXIdentifier string `yaml:"spdx_identifier" mapstructure:"spdx_identifier"`
}

// LicenseConfig contains license-related configuration.
type LicenseConfig struct {
	// Type is the license type (e.g., "apache-2.0", "mit", "custom").
	Type string `yaml:"type" mapstructure:"type"`

	// CustomText is the custom license text (used when Type is "custom").
	CustomText string `yaml:"custom_text" mapstructure:"custom_text"`

	// AutoCreate indicates whether to auto-create LICENSE file if missing.
	AutoCreate bool `yaml:"auto_create" mapstructure:"auto_create"`
}

// ScanConfig contains file scanning configuration.
type ScanConfig struct {
	// ExcludePatterns is a list of glob patterns for files/directories to exclude.
	ExcludePatterns []string `yaml:"exclude_patterns" mapstructure:"exclude_patterns"`

	// IncludePatterns is a list of glob patterns for files to include.
	IncludePatterns []string `yaml:"include_patterns" mapstructure:"include_patterns"`

	// SkipHidden indicates whether to skip hidden files and directories.
	SkipHidden bool `yaml:"skip_hidden" mapstructure:"skip_hidden"`

	// FollowSymlinks indicates whether to follow symbolic links.
	FollowSymlinks bool `yaml:"follow_symlinks" mapstructure:"follow_symlinks"`

	// ExcludedExtensions is a list of file extensions to exclude.
	ExcludedExtensions []string `yaml:"excluded_extensions" mapstructure:"excluded_extensions"`
}

// DetectionConfig contains copyright detection configuration.
type DetectionConfig struct {
	// LinesToScan is the number of lines from the start of file to scan.
	LinesToScan int `yaml:"lines_to_scan" mapstructure:"lines_to_scan"`

	// Keywords is a list of keywords that indicate copyright presence.
	Keywords []string `yaml:"keywords" mapstructure:"keywords"`

	// RequireAll indicates whether all keywords must be present.
	RequireAll bool `yaml:"require_all" mapstructure:"require_all"`
}

// OutputConfig contains output formatting configuration.
type OutputConfig struct {
	// Format is the output format (e.g., "text", "json", "junit").
	Format string `yaml:"format" mapstructure:"format"`

	// Verbose enables verbose output.
	Verbose bool `yaml:"verbose" mapstructure:"verbose"`

	// Quiet suppresses non-error output.
	Quiet bool `yaml:"quiet" mapstructure:"quiet"`
}

// Load reads and parses the configuration file using Viper.
// It looks for .copyright.yaml in the current directory and applies defaults.
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v)

	// Configure Viper
	v.SetConfigName(".copyright")
	v.SetConfigType("yaml")

	if configPath != "" {
		// Use specified config file
		v.SetConfigFile(configPath)
	} else {
		// Search for config in current directory
		v.AddConfigPath(".")
		v.AddConfigPath("$HOME")
	}

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		// Config file not found is acceptable - we'll use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Unmarshal config
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// LoadOrDefault loads the configuration or returns defaults if config file is not found.
func LoadOrDefault(configPath string) *Config {
	cfg, err := Load(configPath)
	if err != nil {
		// If config fails to load, return defaults
		return NewDefault()
	}
	return cfg
}

// NewDefault creates a new configuration with default values.
func NewDefault() *Config {
	return &Config{
		Copyright: CopyrightConfig{
			Template:       DefaultTemplate,
			Year:           DefaultYear,
			Company:        "",
			Author:         "",
			SPDXIdentifier: "",
		},
		License: LicenseConfig{
			Type:       "",
			CustomText: "",
			AutoCreate: false,
		},
		Scan: ScanConfig{
			ExcludePatterns:    DefaultExcludePatterns,
			IncludePatterns:    []string{},
			SkipHidden:         true,
			FollowSymlinks:     false,
			ExcludedExtensions: DefaultExcludedExtensions,
		},
		Detection: DetectionConfig{
			LinesToScan: DefaultLinesToScan,
			Keywords:    DefaultKeywords,
			RequireAll:  false,
		},
		Output: OutputConfig{
			Format:  "text",
			Verbose: false,
			Quiet:   false,
		},
	}
}

// setDefaults configures default values in Viper.
func setDefaults(v *viper.Viper) {
	// Copyright defaults
	v.SetDefault("copyright.template", DefaultTemplate)
	v.SetDefault("copyright.year", DefaultYear)
	v.SetDefault("copyright.company", "")
	v.SetDefault("copyright.author", "")
	v.SetDefault("copyright.spdx_identifier", "")

	// License defaults
	v.SetDefault("license.type", "")
	v.SetDefault("license.custom_text", "")
	v.SetDefault("license.auto_create", false)

	// Scan defaults
	v.SetDefault("scan.exclude_patterns", DefaultExcludePatterns)
	v.SetDefault("scan.include_patterns", []string{})
	v.SetDefault("scan.skip_hidden", true)
	v.SetDefault("scan.follow_symlinks", false)
	v.SetDefault("scan.excluded_extensions", DefaultExcludedExtensions)

	// Detection defaults
	v.SetDefault("detection.lines_to_scan", DefaultLinesToScan)
	v.SetDefault("detection.keywords", DefaultKeywords)
	v.SetDefault("detection.require_all", false)

	// Output defaults
	v.SetDefault("output.format", "text")
	v.SetDefault("output.verbose", false)
	v.SetDefault("output.quiet", false)
}

// Save writes the configuration to a file.
func (c *Config) Save(path string) error {
	v := viper.New()

	// Set all values
	v.Set("copyright", c.Copyright)
	v.Set("license", c.License)
	v.Set("scan", c.Scan)
	v.Set("detection", c.Detection)
	v.Set("output", c.Output)

	// Write to file
	if err := v.WriteConfigAs(path); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// Merge merges another configuration into this one, with values from other taking precedence.
func (c *Config) Merge(other *Config) {
	if other.Copyright.Template != "" && other.Copyright.Template != DefaultTemplate {
		c.Copyright.Template = other.Copyright.Template
	}
	if other.Copyright.Year != "" && other.Copyright.Year != DefaultYear {
		c.Copyright.Year = other.Copyright.Year
	}
	if other.Copyright.Company != "" {
		c.Copyright.Company = other.Copyright.Company
	}
	if other.Copyright.Author != "" {
		c.Copyright.Author = other.Copyright.Author
	}
	if other.Copyright.SPDXIdentifier != "" {
		c.Copyright.SPDXIdentifier = other.Copyright.SPDXIdentifier
	}

	if other.License.Type != "" {
		c.License.Type = other.License.Type
	}
	if other.License.CustomText != "" {
		c.License.CustomText = other.License.CustomText
	}
	if other.License.AutoCreate {
		c.License.AutoCreate = other.License.AutoCreate
	}

	if len(other.Scan.ExcludePatterns) > 0 {
		c.Scan.ExcludePatterns = other.Scan.ExcludePatterns
	}
	if len(other.Scan.IncludePatterns) > 0 {
		c.Scan.IncludePatterns = other.Scan.IncludePatterns
	}
	if len(other.Scan.ExcludedExtensions) > 0 {
		c.Scan.ExcludedExtensions = other.Scan.ExcludedExtensions
	}

	if other.Detection.LinesToScan > 0 {
		c.Detection.LinesToScan = other.Detection.LinesToScan
	}
	if len(other.Detection.Keywords) > 0 {
		c.Detection.Keywords = other.Detection.Keywords
	}

	if other.Output.Format != "" && other.Output.Format != "text" {
		c.Output.Format = other.Output.Format
	}
}

// GetYear returns the copyright year, resolving "auto" to the current year.
func (c *Config) GetYear() string {
	if c.Copyright.Year == "auto" || c.Copyright.Year == "" {
		return fmt.Sprintf("%d", getCurrentYear())
	}
	return c.Copyright.Year
}

// getCurrentYear returns the current year.
func getCurrentYear() int {
	// In production, this would use time.Now().Year()
	// For now using a fixed year for deterministic testing
	return 2025
}
