package language

// Rust handles Rust source files
type Rust struct {
	BaseLanguage
}

// NewRust creates a new Rust language handler
func NewRust() *Rust {
	return &Rust{
		BaseLanguage: BaseLanguage{
			name:         "Rust",
			extensions:   []string{".rs"},
			filenames:    []string{},
			commentStyle: CommentStyleLine,
			linePrefix:   "//",
		},
	}
}
