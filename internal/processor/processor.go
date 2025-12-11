package processor

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/leansecurity-co/copyright/internal/config"
	"github.com/leansecurity-co/copyright/internal/copyright"
	"github.com/leansecurity-co/copyright/internal/files"
	"github.com/leansecurity-co/copyright/internal/language"
	"github.com/leansecurity-co/copyright/internal/license"
)

// Processor orchestrates the copyright management workflow
type Processor struct {
	config    *config.Config
	scanner   *files.Scanner
	checker   *copyright.Checker
	generator *copyright.HeaderGenerator
	writer    *Writer
	registry  *language.Registry
	license   *license.Manager
}

// ProcessorOptions configures processor behavior
type ProcessorOptions struct {
	DryRun bool // If true, don't modify files
	Backup bool // If true, create .bak files before modifying
}

// NewProcessor creates a new processor with the given configuration
func NewProcessor(cfg *config.Config) (*Processor, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Create checker options from config
	checkerOpts := &copyright.CheckerOptions{
		LinesToScan: cfg.Detection.LinesToScan,
		Keywords:    cfg.Detection.Keywords,
		RequireAll:  cfg.Detection.RequireAll,
	}

	// Create license manager with config values
	licenseData := &license.LicenseData{
		Year:    cfg.Copyright.Year,
		Company: cfg.Copyright.Company,
		Author:  cfg.Copyright.Author,
	}

	return &Processor{
		config:    cfg,
		checker:   copyright.NewChecker(checkerOpts),
		generator: copyright.NewHeaderGenerator(cfg),
		registry:  language.GetRegistry(),
		license:   license.NewManager(cfg.License.Type, cfg.License.CustomText, licenseData),
	}, nil
}

// Process runs the copyright check/fix workflow on the given root path
func (p *Processor) Process(ctx context.Context, rootPath string, mode Mode) (*ProcessResult, error) {
	return p.ProcessWithOptions(ctx, rootPath, mode, &ProcessorOptions{})
}

// ProcessWithOptions runs the copyright check/fix workflow with custom options
func (p *Processor) ProcessWithOptions(ctx context.Context, rootPath string, mode Mode, opts *ProcessorOptions) (*ProcessResult, error) {
	// Initialize result
	result := &ProcessResult{
		FileResults: make([]FileResult, 0),
	}

	// Get absolute path
	absPath, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Verify path exists
	if _, err := os.Stat(absPath); err != nil {
		return nil, fmt.Errorf("path does not exist: %w", err)
	}

	// Create scanner options from config
	scanOpts := &files.ScanOptions{
		ExcludePatterns:    p.config.Scan.ExcludePatterns,
		IncludePatterns:    p.config.Scan.IncludePatterns,
		ExcludedExtensions: p.config.Scan.ExcludedExtensions,
		SkipHidden:         p.config.Scan.SkipHidden,
		FollowSymlinks:     p.config.Scan.FollowSymlinks,
	}

	// Create scanner
	scanner, err := files.NewScanner(absPath, scanOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to create scanner: %w", err)
	}
	p.scanner = scanner

	// Create writer with options
	p.writer = NewWriter(opts.DryRun, opts.Backup)

	// Scan for files
	scanResult, err := scanner.Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to scan files: %w", err)
	}

	// Process each file
	for _, file := range scanResult.Files {
		// Check for context cancellation
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
		}

		fileResult := p.processFile(ctx, file, mode)
		result.FileResults = append(result.FileResults, *fileResult)

		// Update counters
		result.TotalFiles++
		switch fileResult.Status {
		case StatusHasCopyright:
			result.FilesWithCopyright++
		case StatusMissingCopyright:
			result.FilesMissingCopyright++
		case StatusFixed:
			result.FilesFixed++
		case StatusSkipped:
			result.FilesSkipped++
		case StatusFailed:
			result.FilesFailed++
		}
	}

	// Handle LICENSE file
	if mode == ModeFix {
		licensePath, licenseExists := p.license.FindLicense(absPath)
		if licenseExists {
			result.LicenseStatus = fmt.Sprintf("LICENSE exists at %s", licensePath)
		} else {
			// Try to create if configured
			if p.config.License.AutoCreate {
				if err := p.license.Create(absPath); err != nil {
					result.LicenseStatus = fmt.Sprintf("Failed to create LICENSE: %v", err)
				} else {
					result.LicenseCreated = true
					result.LicenseStatus = "LICENSE created"
				}
			} else {
				result.LicenseStatus = "No LICENSE file found (auto-create disabled)"
			}
		}
	}

	return result, nil
}

// processFile handles processing of a single file
func (p *Processor) processFile(ctx context.Context, file files.FileInfo, mode Mode) *FileResult {
	result := &FileResult{
		FilePath:     file.Path,
		RelativePath: file.RelativePath,
	}

	// Skip binary files
	if file.IsBinary {
		result.Status = StatusSkipped
		result.Message = "binary file"
		return result
	}

	// Get language handler
	lang := p.registry.GetLanguage(file.Filename)
	if lang == nil {
		result.Status = StatusSkipped
		result.Message = "unsupported file type"
		return result
	}

	// Skip files that don't support comments
	if !lang.CanHaveComments() {
		result.Status = StatusSkipped
		result.Message = "file type does not support comments"
		return result
	}

	// Check if file has copyright
	checkResult, err := p.checker.Check(ctx, file.Path)
	if err != nil {
		result.Status = StatusFailed
		result.Error = err
		result.Message = fmt.Sprintf("failed to check file: %v", err)
		return result
	}

	if checkResult.Error != nil {
		result.Status = StatusFailed
		result.Error = checkResult.Error
		result.Message = fmt.Sprintf("check error: %v", checkResult.Error)
		return result
	}

	// If file has copyright, we're done
	if checkResult.HasCopyright {
		result.Status = StatusHasCopyright
		result.Message = fmt.Sprintf("copyright found at line %d", checkResult.LineFound)
		return result
	}

	// File is missing copyright
	if mode == ModeCheck {
		result.Status = StatusMissingCopyright
		result.Message = "copyright header missing"
		return result
	}

	// Mode is ModeFix - add the copyright header
	if err := p.addCopyrightHeader(file.Path, lang); err != nil {
		result.Status = StatusFailed
		result.Error = err
		result.Message = fmt.Sprintf("failed to add copyright: %v", err)
		return result
	}

	result.Status = StatusFixed
	result.Message = "copyright header added"
	return result
}

// addCopyrightHeader adds a copyright header to a file
func (p *Processor) addCopyrightHeader(filePath string, lang language.Language) error {
	// Generate the copyright text
	copyrightText, err := p.generator.Generate()
	if err != nil {
		return fmt.Errorf("failed to generate copyright: %w", err)
	}

	// Format the header for the specific language
	header := lang.FormatHeader(copyrightText)

	// Write the header to the file
	if err := p.writer.WriteHeader(filePath, header, lang); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	return nil
}
