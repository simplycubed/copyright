package files

// FileInfo represents metadata about a discovered file
type FileInfo struct {
	Path         string // Absolute file path
	RelativePath string // Relative to repository root
	Extension    string // Including dot (e.g., ".go")
	Filename     string // Just the filename
	Size         int64
	IsBinary     bool
}

// ScanResult contains the results of a directory scan
type ScanResult struct {
	Files        []FileInfo
	SkippedFiles []string
	Errors       []error
}

// ScanOptions configures scanning behavior
type ScanOptions struct {
	ExcludePatterns    []string
	IncludePatterns    []string
	ExcludedExtensions []string
	SkipHidden         bool
	FollowSymlinks     bool
}

// DefaultScanOptions returns a new ScanOptions with sensible defaults
func DefaultScanOptions() *ScanOptions {
	return &ScanOptions{
		SkipHidden:     true,
		FollowSymlinks: false,
		ExcludePatterns: []string{
			"**/.git/**",
			"**/node_modules/**",
			"**/vendor/**",
		},
	}
}
