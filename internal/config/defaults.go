package config

// Default configuration values for the copyright tool.
const (
	// DefaultTemplate is the default copyright header template.
	DefaultTemplate = "Copyright {{.Year}} {{.Company}}\nSPDX-License-Identifier: {{.SPDX}}"

	// DefaultYear is the default year value (auto = current year).
	DefaultYear = "auto"

	// DefaultLinesToScan is the default number of lines to scan for copyright.
	DefaultLinesToScan = 20
)

var (
	// DefaultKeywords are the default keywords used to detect copyright headers.
	DefaultKeywords = []string{
		"copyright",
		"©",
		"spdx-license-identifier",
	}

	// DefaultExcludedExtensions are file extensions excluded from scanning by default.
	DefaultExcludedExtensions = []string{
		// Data/Config files
		".json",
		".lock",
		".sum",
		".mod",
		".toml",
		".ini",
		".cfg",

		// Images
		".png",
		".jpg",
		".jpeg",
		".gif",
		".ico",
		".svg",
		".bmp",
		".webp",

		// Fonts
		".woff",
		".woff2",
		".ttf",
		".eot",
		".otf",

		// Documents
		".pdf",
		".doc",
		".docx",

		// Archives
		".zip",
		".tar",
		".gz",
		".bz2",
		".xz",
		".7z",
		".rar",

		// Binaries/Libraries
		".exe",
		".dll",
		".so",
		".dylib",
		".a",
		".o",

		// Media
		".mp3",
		".mp4",
		".avi",
		".mov",
		".wmv",
		".flv",
		".wav",
	}

	// DefaultExcludePatterns are glob patterns excluded from scanning by default.
	DefaultExcludePatterns = []string{
		"**/node_modules/**",
		"**/vendor/**",
		"**/.git/**",
		"**/.svn/**",
		"**/.hg/**",
		"**/dist/**",
		"**/build/**",
		"**/*.min.js",
		"**/*.min.css",
	}
)

// ValidOutputFormats lists all valid output format options.
var ValidOutputFormats = []string{
	"text",
	"json",
	"junit",
	"sarif",
}

// ValidLicenseTypes lists commonly supported license types.
var ValidLicenseTypes = []string{
	"apache-2.0",
	"mit",
	"gpl-3.0",
	"gpl-2.0",
	"lgpl-3.0",
	"lgpl-2.1",
	"bsd-2-clause",
	"bsd-3-clause",
	"mpl-2.0",
	"agpl-3.0",
	"unlicense",
	"cc0-1.0",
	"custom",
}
