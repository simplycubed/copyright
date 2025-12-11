package language

// JavaScript handles JavaScript and TypeScript files
type JavaScript struct {
	BaseLanguage
}

// NewJavaScript creates a new JavaScript language handler
func NewJavaScript() *JavaScript {
	return &JavaScript{
		BaseLanguage: BaseLanguage{
			name: "JavaScript",
			extensions: []string{
				".js", ".ts", ".jsx", ".tsx",
				".mjs", ".cjs", ".mts", ".cts",
			},
			filenames:    []string{},
			commentStyle: CommentStyleBlock,
			blockStart:   "/*",
			blockPrefix:  " *",
			blockEnd:     " */",
		},
	}
}

// PreserveLines handles shebang for Node.js scripts
func (j *JavaScript) PreserveLines(content []byte) []string {
	var preserve []string

	if shebang := getShebangLine(content); shebang != "" {
		preserve = append(preserve, shebang)
	}

	return preserve
}

// InsertPosition returns the byte offset where copyright should be inserted
func (j *JavaScript) InsertPosition(content []byte) int {
	return calculateInsertPosition(j.PreserveLines(content))
}
