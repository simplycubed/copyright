package language

// Java handles Java source files
type Java struct {
	BaseLanguage
}

// NewJava creates a new Java language handler
func NewJava() *Java {
	return &Java{
		BaseLanguage: BaseLanguage{
			name:         "Java",
			extensions:   []string{".java"},
			filenames:    []string{},
			commentStyle: CommentStyleBlock,
			blockStart:   "/*",
			blockPrefix:  " *",
			blockEnd:     " */",
		},
	}
}
