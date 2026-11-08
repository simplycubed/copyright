package files

import (
	"path/filepath"
	"strings"
)

// MatchPattern matches a glob pattern against a path
// Supports:
// - * matches any characters except /
// - ** matches any characters including /
// - ? matches a single character
// - [abc] character classes
func MatchPattern(pattern, path string) (bool, error) {
	// Normalize path separators
	pattern = filepath.ToSlash(pattern)
	path = filepath.ToSlash(path)

	// Handle ** patterns specially
	if strings.Contains(pattern, "**") {
		return matchDoubleStarPattern(pattern, path)
	}

	// Use filepath.Match for simple patterns
	matched, err := filepath.Match(pattern, path)
	if err != nil {
		return false, err
	}

	// Also try matching just the base name
	if !matched {
		matched, err = filepath.Match(pattern, filepath.Base(path))
	}

	return matched, err
}

// MatchAnyPattern returns true if the path matches any of the provided patterns
func MatchAnyPattern(patterns []string, path string) bool {
	for _, pattern := range patterns {
		matched, err := MatchPattern(pattern, path)
		if err != nil {
			// If pattern is invalid, skip it
			continue
		}
		if matched {
			return true
		}
	}
	return false
}

// matchDoubleStarPattern handles patterns containing **
func matchDoubleStarPattern(pattern, path string) (bool, error) {
	// Split pattern by **
	parts := strings.Split(pattern, "**")

	// If pattern starts with **, match any prefix
	if strings.HasPrefix(pattern, "**") {
		// Remove leading **
		if len(parts) > 1 && parts[0] == "" {
			parts = parts[1:]
		}
	}

	// If pattern ends with **, match any suffix
	if strings.HasSuffix(pattern, "**") {
		// Remove trailing **
		if len(parts) > 0 && parts[len(parts)-1] == "" {
			parts = parts[:len(parts)-1]
		}
	}

	// Handle single ** (matches everything)
	if len(parts) == 0 || (len(parts) == 1 && parts[0] == "") {
		return true, nil
	}

	// For multiple parts, we need to match them in sequence
	return matchPatternParts(parts, path)
}

// matchPatternParts matches a path against pattern parts separated by **
func matchPatternParts(parts []string, path string) (bool, error) {
	if len(parts) == 0 {
		return true, nil
	}

	// Clean up empty parts and trim slashes
	var cleanParts []string
	for _, part := range parts {
		part = strings.Trim(part, "/")
		if part != "" {
			cleanParts = parts
		}
	}

	if len(cleanParts) == 0 {
		return true, nil
	}

	// For simplicity with ** patterns, check if the path contains all parts in order
	currentPath := path
	for i, part := range cleanParts {
		part = strings.Trim(part, "/")
		if part == "" {
			continue
		}

		// For the first part, check if it matches the beginning
		if i == 0 && strings.HasPrefix(path, part) {
			// Check if this is an exact segment match
			matched, err := matchPathSegment(part, currentPath, true)
			if err != nil {
				return false, err
			}
			if matched {
				// Move past this part
				idx := strings.Index(currentPath, "/")
				if idx > 0 {
					currentPath = currentPath[idx+1:]
				} else {
					currentPath = ""
				}
				continue
			}
		}

		// For middle parts, find them anywhere
		matched := false
		segments := strings.Split(currentPath, "/")
		for j, segment := range segments {
			segmentMatched, err := filepath.Match(part, segment)
			if err != nil {
				return false, err
			}
			if segmentMatched {
				matched = true
				// Move past this segment
				if j+1 < len(segments) {
					currentPath = strings.Join(segments[j+1:], "/")
				} else {
					currentPath = ""
				}
				break
			}
		}

		if !matched {
			return false, nil
		}
	}

	return true, nil
}

// matchPathSegment checks if a pattern matches a segment at the start of the path
func matchPathSegment(pattern, path string, atStart bool) (bool, error) {
	segments := strings.Split(path, "/")
	if len(segments) == 0 {
		return false, nil
	}

	if atStart {
		return filepath.Match(pattern, segments[0])
	}

	// Check any segment
	for _, segment := range segments {
		matched, err := filepath.Match(pattern, segment)
		if err != nil {
			return false, err
		}
		if matched {
			return true, nil
		}
	}

	return false, nil
}
