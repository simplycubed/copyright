package language

// JSON handles JSON files
// Note: JSON doesn't support comments, so copyright headers cannot be added
type JSON struct {
	BaseLanguage
}

// NewJSON creates a new JSON language handler
func NewJSON() *JSON {
	return &JSON{
		BaseLanguage: BaseLanguage{
			name:         "JSON",
			extensions:   []string{".json", ".jsonc"},
			filenames:    []string{},
			commentStyle: CommentStyleNone,
		},
	}
}
