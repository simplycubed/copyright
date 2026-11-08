package language

import (
	"path/filepath"
	"strings"
	"sync"
)

// Registry manages the mapping of file extensions and names to language handlers
type Registry struct {
	mu                 sync.RWMutex
	extensionMap       map[string]Language
	filenameMap        map[string]Language
	caseInsensitiveMap map[string]Language // For case-insensitive filename matching
}

var defaultRegistry *Registry
var once sync.Once

// GetRegistry returns the singleton registry instance
func GetRegistry() *Registry {
	once.Do(func() {
		defaultRegistry = &Registry{
			extensionMap:       make(map[string]Language),
			filenameMap:        make(map[string]Language),
			caseInsensitiveMap: make(map[string]Language),
		}
		registerAllLanguages()
	})
	return defaultRegistry
}

// Register adds a language to the registry
func (r *Registry) Register(lang Language) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Register extensions
	for _, ext := range lang.Extensions() {
		r.extensionMap[ext] = lang
	}

	// Register specific filenames
	for _, filename := range lang.Filenames() {
		r.filenameMap[filename] = lang
		// Also register lowercase version for case-insensitive matching
		r.caseInsensitiveMap[strings.ToLower(filename)] = lang
	}
}

// GetLanguage returns the language handler for a given filename
func (r *Registry) GetLanguage(filename string) Language {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Extract the base filename
	base := filepath.Base(filename)

	// First, check exact filename match
	if lang, ok := r.filenameMap[base]; ok {
		return lang
	}

	// Check case-insensitive filename match
	if lang, ok := r.caseInsensitiveMap[strings.ToLower(base)]; ok {
		return lang
	}

	// Then check by extension
	ext := filepath.Ext(filename)
	if lang, ok := r.extensionMap[ext]; ok {
		return lang
	}

	// Check for double extensions (e.g., .test.js)
	if len(ext) > 0 {
		// Remove the extension and try again
		nameWithoutExt := strings.TrimSuffix(base, ext)
		secondExt := filepath.Ext(nameWithoutExt)
		if len(secondExt) > 0 {
			combinedExt := secondExt + ext
			if lang, ok := r.extensionMap[combinedExt]; ok {
				return lang
			}
		}
	}

	return nil
}

// SupportedExtensions returns all registered file extensions
func (r *Registry) SupportedExtensions() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	exts := make([]string, 0, len(r.extensionMap))
	for ext := range r.extensionMap {
		exts = append(exts, ext)
	}
	return exts
}

// SupportedFilenames returns all registered specific filenames
func (r *Registry) SupportedFilenames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.filenameMap))
	for name := range r.filenameMap {
		names = append(names, name)
	}
	return names
}

// registerAllLanguages registers all built-in language handlers
func registerAllLanguages() {
	registry := defaultRegistry

	// Register all language handlers
	registry.Register(NewGoLang())
	registry.Register(NewJavaScript())
	registry.Register(NewJava())
	registry.Register(NewCSharp())
	registry.Register(NewRust())
	registry.Register(NewC())
	registry.Register(NewPython())
	registry.Register(NewRuby())
	registry.Register(NewShell())
	registry.Register(NewYAML())
	registry.Register(NewHTML())
	registry.Register(NewCSS())
	registry.Register(NewSQL())
	registry.Register(NewMarkdown())
	registry.Register(NewJSON())
	registry.Register(NewProtocolBuffers())
	registry.Register(NewMakefile())
	registry.Register(NewDockerfile())
	registry.Register(NewTerraform())
}
