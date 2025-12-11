package language

// Shell handles shell script files
type Shell struct {
	BaseLanguage
}

// NewShell creates a new Shell language handler
func NewShell() *Shell {
	return &Shell{
		BaseLanguage: BaseLanguage{
			name:         "Shell",
			extensions:   []string{".sh", ".bash", ".zsh", ".fish", ".ksh"},
			filenames:    []string{".bashrc", ".zshrc", ".bash_profile", ".profile"},
			commentStyle: CommentStyleLine,
			linePrefix:   "#",
		},
	}
}

// PreserveLines handles shebang for shell scripts
func (s *Shell) PreserveLines(content []byte) []string {
	var preserve []string

	if shebang := getShebangLine(content); shebang != "" {
		preserve = append(preserve, shebang)
	}

	return preserve
}

// InsertPosition returns the byte offset where copyright should be inserted
func (s *Shell) InsertPosition(content []byte) int {
	return calculateInsertPosition(s.PreserveLines(content))
}
