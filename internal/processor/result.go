package processor

import "fmt"

// Mode represents the operation mode
type Mode int

const (
	// ModeCheck checks if files have copyright headers
	ModeCheck Mode = iota
	// ModeFix adds missing copyright headers to files
	ModeFix
)

// FileStatus represents the processing status of a file
type FileStatus int

const (
	// StatusHasCopyright indicates the file already has a copyright header
	StatusHasCopyright FileStatus = iota
	// StatusMissingCopyright indicates the file is missing a copyright header
	StatusMissingCopyright
	// StatusFixed indicates the copyright header was added to the file
	StatusFixed
	// StatusSkipped indicates the file was skipped (e.g., unsupported language)
	StatusSkipped
	// StatusFailed indicates processing failed for the file
	StatusFailed
)

// String returns the string representation of FileStatus
func (s FileStatus) String() string {
	switch s {
	case StatusHasCopyright:
		return "has_copyright"
	case StatusMissingCopyright:
		return "missing_copyright"
	case StatusFixed:
		return "fixed"
	case StatusSkipped:
		return "skipped"
	case StatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

// FileResult represents the result of processing a single file
type FileResult struct {
	FilePath     string
	RelativePath string
	Status       FileStatus
	Message      string
	Error        error
}

// ProcessResult represents the overall processing result
type ProcessResult struct {
	TotalFiles            int
	FilesWithCopyright    int
	FilesMissingCopyright int
	FilesFixed            int
	FilesSkipped          int
	FilesFailed           int
	FileResults           []FileResult
	LicenseCreated        bool
	LicenseStatus         string
}

// ExitCode returns the appropriate exit code based on the process result
// Returns 0 if all files have copyright headers, 1 if any are missing or failed
func (r *ProcessResult) ExitCode() int {
	if r.FilesMissingCopyright > 0 || r.FilesFailed > 0 {
		return 1
	}
	return 0
}

// Summary returns a human-readable summary of the process result
func (r *ProcessResult) Summary() string {
	return fmt.Sprintf(
		"Total: %d, Has Copyright: %d, Missing: %d, Fixed: %d, Skipped: %d, Failed: %d",
		r.TotalFiles,
		r.FilesWithCopyright,
		r.FilesMissingCopyright,
		r.FilesFixed,
		r.FilesSkipped,
		r.FilesFailed,
	)
}
