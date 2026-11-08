// Package exitcode defines exit codes for the copyright CLI tool.
package exitcode

const (
	// Success indicates all files have copyright headers.
	Success = 0

	// FilesMissingCopyright indicates some files are missing copyright headers (check mode).
	FilesMissingCopyright = 1

	// ConfigError indicates a configuration error occurred.
	ConfigError = 2

	// FileSystemError indicates a file system access error occurred.
	FileSystemError = 3

	// InternalError indicates an unexpected internal error occurred.
	InternalError = 4
)
