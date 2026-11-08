package language

// Markdown handles Markdown files
// Note: Markdown doesn't support comments, so copyright headers cannot be added
type Markdown struct {
	BaseLanguage
}

// NewMarkdown creates a new Markdown language handler
func NewMarkdown() *Markdown {
	return &Markdown{
		BaseLanguage: BaseLanguage{
			name:         "Markdown",
			extensions:   []string{".md", ".markdown"},
			filenames:    []string{},
			commentStyle: CommentStyleNone,
		},
	}
}
