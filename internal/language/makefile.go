package language

// Makefile handles Makefile files
type Makefile struct {
	BaseLanguage
}

// NewMakefile creates a new Makefile language handler
func NewMakefile() *Makefile {
	return &Makefile{
		BaseLanguage: BaseLanguage{
			name:         "Makefile",
			extensions:   []string{".mk"},
			filenames:    []string{"Makefile", "makefile", "GNUmakefile"},
			commentStyle: CommentStyleLine,
			linePrefix:   "#",
		},
	}
}
