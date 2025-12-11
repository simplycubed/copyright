package language

// CSharp handles C# source files
type CSharp struct {
	BaseLanguage
}

// NewCSharp creates a new C# language handler
func NewCSharp() *CSharp {
	return &CSharp{
		BaseLanguage: BaseLanguage{
			name:         "C#",
			extensions:   []string{".cs"},
			filenames:    []string{},
			commentStyle: CommentStyleBlock,
			blockStart:   "/*",
			blockPrefix:  " *",
			blockEnd:     " */",
		},
	}
}
