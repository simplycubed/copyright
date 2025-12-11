package files

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
)

const (
	// bufferSize is the number of bytes to read for binary detection
	bufferSize = 8192
)

// commonBinarySignatures contains magic numbers for common binary formats
var commonBinarySignatures = [][]byte{
	{0x7F, 0x45, 0x4C, 0x46},                   // ELF executable
	{0x4D, 0x5A},                               // DOS/Windows executable
	{0x50, 0x4B, 0x03, 0x04},                   // ZIP archive
	{0x50, 0x4B, 0x05, 0x06},                   // ZIP archive (empty)
	{0x50, 0x4B, 0x07, 0x08},                   // ZIP archive (spanned)
	{0x1F, 0x8B},                               // GZIP compressed
	{0x42, 0x5A, 0x68},                         // BZIP2 compressed
	{0xFD, 0x37, 0x7A, 0x58, 0x5A, 0x00},       // XZ compressed
	{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A}, // PNG image
	{0xFF, 0xD8, 0xFF},                         // JPEG image
	{0x47, 0x49, 0x46, 0x38},                   // GIF image
	{0x25, 0x50, 0x44, 0x46},                   // PDF
	{0xCA, 0xFE, 0xBA, 0xBE},                   // Java class file
	{0xCE, 0xFA, 0xED, 0xFE},                   // Mach-O binary (32-bit)
	{0xCF, 0xFA, 0xED, 0xFE},                   // Mach-O binary (64-bit)
	{0xFE, 0xED, 0xFA, 0xCE},                   // Mach-O binary (reverse byte ordering)
	{0xFE, 0xED, 0xFA, 0xCF},                   // Mach-O binary (reverse byte ordering)
}

// commonBinaryExtensions is a list of file extensions commonly associated with binary files
var commonBinaryExtensions = map[string]bool{
	".exe":    true,
	".dll":    true,
	".so":     true,
	".dylib":  true,
	".a":      true,
	".o":      true,
	".bin":    true,
	".dat":    true,
	".db":     true,
	".sqlite": true,
	".png":    true,
	".jpg":    true,
	".jpeg":   true,
	".gif":    true,
	".bmp":    true,
	".ico":    true,
	".svg":    true,
	".webp":   true,
	".pdf":    true,
	".zip":    true,
	".tar":    true,
	".gz":     true,
	".bz2":    true,
	".xz":     true,
	".7z":     true,
	".rar":    true,
	".woff":   true,
	".woff2":  true,
	".ttf":    true,
	".eot":    true,
	".mp3":    true,
	".mp4":    true,
	".avi":    true,
	".mov":    true,
	".wav":    true,
	".class":  true,
	".jar":    true,
	".war":    true,
	".ear":    true,
	".pyc":    true,
	".pyo":    true,
}

// IsBinary determines if a file is binary by checking:
// 1. File extension
// 2. Magic number signatures
// 3. Null byte detection in first 8KB
func IsBinary(path string) (bool, error) {
	// Quick check: common binary extensions
	ext := filepath.Ext(path)
	if commonBinaryExtensions[ext] {
		return true, nil
	}

	// Open the file
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer func() { _ = f.Close() }()

	// Read first bufferSize bytes
	buffer := make([]byte, bufferSize)
	n, err := f.Read(buffer)
	if err != nil && err != io.EOF {
		return false, err
	}
	buffer = buffer[:n]

	// Empty files are not binary
	if n == 0 {
		return false, nil
	}

	// Check for magic number signatures
	if hasBinarySignature(buffer) {
		return true, nil
	}

	// Check for null bytes (strong indicator of binary content)
	if bytes.Contains(buffer, []byte{0}) {
		return true, nil
	}

	// Check for high ratio of non-printable characters
	nonPrintable := 0
	for _, b := range buffer {
		// Count non-ASCII printable characters
		// ASCII printable: 32-126, plus common whitespace: 9 (tab), 10 (LF), 13 (CR)
		if (b < 32 || b > 126) && b != 9 && b != 10 && b != 13 {
			nonPrintable++
		}
	}

	// If more than 30% of characters are non-printable, likely binary
	if float64(nonPrintable)/float64(n) > 0.3 {
		return true, nil
	}

	return false, nil
}

// hasBinarySignature checks if the buffer starts with a known binary signature
func hasBinarySignature(buffer []byte) bool {
	for _, signature := range commonBinarySignatures {
		if len(buffer) >= len(signature) {
			if bytes.HasPrefix(buffer, signature) {
				return true
			}
		}
	}
	return false
}
