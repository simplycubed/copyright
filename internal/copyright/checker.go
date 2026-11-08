package copyright

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
)

// CheckerOptions controls how a file is scanned for an existing copyright notice.
type CheckerOptions struct {
	LinesToScan int
	Keywords    []string
	RequireAll  bool
}

// CheckResult is the outcome of scanning one file.
type CheckResult struct {
	HasCopyright bool
	LineFound    int
	Error        error
}

// Checker scans the start of a file for configured copyright keywords.
type Checker struct {
	opts CheckerOptions
}

// NewChecker returns a checker. A nil options pointer scans nothing.
func NewChecker(opts *CheckerOptions) *Checker {
	c := &Checker{}
	if opts != nil {
		c.opts = *opts
	}
	return c
}

// Check reads the first LinesToScan lines of path and reports whether the
// configured keywords are present. Matching is case-insensitive. LineFound is
// 1-based and is the first scanned line that contains a keyword.
func (c *Checker) Check(ctx context.Context, path string) (*CheckResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() { _ = file.Close() }()

	result := &CheckResult{}
	if c.opts.LinesToScan <= 0 || len(c.opts.Keywords) == 0 {
		return result, nil
	}

	keywords := make([]string, len(c.opts.Keywords))
	for i, keyword := range c.opts.Keywords {
		keywords[i] = strings.ToLower(keyword)
	}
	found := make([]bool, len(keywords))

	scanner := bufio.NewScanner(file)
	lineNo := 0
	for scanner.Scan() {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		lineNo++
		if lineNo > c.opts.LinesToScan {
			break
		}
		line := strings.ToLower(scanner.Text())
		for i, keyword := range keywords {
			if keyword == "" || !strings.Contains(line, keyword) {
				continue
			}
			found[i] = true
			if result.LineFound == 0 {
				result.LineFound = lineNo
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	if c.opts.RequireAll {
		result.HasCopyright = true
		for _, ok := range found {
			if !ok {
				result.HasCopyright = false
				result.LineFound = 0
				break
			}
		}
		return result, nil
	}

	result.HasCopyright = result.LineFound > 0
	return result, nil
}
