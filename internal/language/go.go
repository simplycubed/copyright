package language

// GoLang handles Go source files
type GoLang struct {
	BaseLanguage
}

// NewGoLang creates a new Go language handler
func NewGoLang() *GoLang {
	return &GoLang{
		BaseLanguage: BaseLanguage{
			name:         "Go",
			extensions:   []string{".go"},
			filenames:    []string{},
			commentStyle: CommentStyleLine,
			linePrefix:   "//",
		},
	}
}
