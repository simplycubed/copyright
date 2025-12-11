package language

// CSS handles CSS and preprocessor files
type CSS struct {
	BaseLanguage
}

// NewCSS creates a new CSS language handler
func NewCSS() *CSS {
	return &CSS{
		BaseLanguage: BaseLanguage{
			name:         "CSS",
			extensions:   []string{".css", ".scss", ".sass", ".less"},
			filenames:    []string{},
			commentStyle: CommentStyleBlock,
			blockStart:   "/*",
			blockPrefix:  " *",
			blockEnd:     " */",
		},
	}
}
