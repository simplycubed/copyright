package files

import (
	"path/filepath"
	"strings"
)

// FilterByExtension returns files matching the given extensions
func FilterByExtension(files []FileInfo, extensions ...string) []FileInfo {
	if len(extensions) == 0 {
		return files
	}

	// Normalize extensions
	normalizedExts := make(map[string]bool)
	for _, ext := range extensions {
		if ext != "" && !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		normalizedExts[strings.ToLower(ext)] = true
	}

	var result []FileInfo
	for _, file := range files {
		if normalizedExts[strings.ToLower(file.Extension)] {
			result = append(result, file)
		}
	}

	return result
}

// FilterTextFiles returns only non-binary files
func FilterTextFiles(files []FileInfo) []FileInfo {
	var result []FileInfo
	for _, file := range files {
		if !file.IsBinary {
			result = append(result, file)
		}
	}
	return result
}

// FilterBinaryFiles returns only binary files
func FilterBinaryFiles(files []FileInfo) []FileInfo {
	var result []FileInfo
	for _, file := range files {
		if file.IsBinary {
			result = append(result, file)
		}
	}
	return result
}

// GroupByExtension groups files by their extension
func GroupByExtension(files []FileInfo) map[string][]FileInfo {
	groups := make(map[string][]FileInfo)
	for _, file := range files {
		ext := file.Extension
		if ext == "" {
			ext = "(no extension)"
		}
		groups[ext] = append(groups[ext], file)
	}
	return groups
}

// GroupByDirectory groups files by their parent directory
func GroupByDirectory(files []FileInfo) map[string][]FileInfo {
	groups := make(map[string][]FileInfo)
	for _, file := range files {
		dir := filepath.Dir(file.RelativePath)
		if dir == "." {
			dir = "(root)"
		}
		groups[dir] = append(groups[dir], file)
	}
	return groups
}

// GetRelativePaths returns just the relative paths from a slice of FileInfo
func GetRelativePaths(files []FileInfo) []string {
	paths := make([]string, len(files))
	for i, file := range files {
		paths[i] = file.RelativePath
	}
	return paths
}

// GetAbsolutePaths returns just the absolute paths from a slice of FileInfo
func GetAbsolutePaths(files []FileInfo) []string {
	paths := make([]string, len(files))
	for i, file := range files {
		paths[i] = file.Path
	}
	return paths
}

// TotalSize calculates the total size of all files
func TotalSize(files []FileInfo) int64 {
	var total int64
	for _, file := range files {
		total += file.Size
	}
	return total
}

// CountByExtension returns a count of files per extension
func CountByExtension(files []FileInfo) map[string]int {
	counts := make(map[string]int)
	for _, file := range files {
		ext := file.Extension
		if ext == "" {
			ext = "(no extension)"
		}
		counts[ext]++
	}
	return counts
}

// FindLargestFiles returns the N largest files
func FindLargestFiles(files []FileInfo, n int) []FileInfo {
	if n <= 0 || len(files) == 0 {
		return nil
	}

	// Make a copy to avoid modifying the original slice
	sorted := make([]FileInfo, len(files))
	copy(sorted, files)

	// Simple selection sort for finding top N
	for i := 0; i < n && i < len(sorted); i++ {
		maxIdx := i
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].Size > sorted[maxIdx].Size {
				maxIdx = j
			}
		}
		sorted[i], sorted[maxIdx] = sorted[maxIdx], sorted[i]
	}

	if n > len(sorted) {
		n = len(sorted)
	}
	return sorted[:n]
}

// FilterBySize returns files within the given size range (inclusive)
// Use -1 for min or max to ignore that boundary
func FilterBySize(files []FileInfo, minSize, maxSize int64) []FileInfo {
	var result []FileInfo
	for _, file := range files {
		if (minSize < 0 || file.Size >= minSize) &&
			(maxSize < 0 || file.Size <= maxSize) {
			result = append(result, file)
		}
	}
	return result
}

// IsHidden returns true if the file or any of its parent directories is hidden
func IsHidden(path string) bool {
	parts := strings.Split(filepath.ToSlash(path), "/")
	for _, part := range parts {
		if strings.HasPrefix(part, ".") && part != "." && part != ".." {
			return true
		}
	}
	return false
}
