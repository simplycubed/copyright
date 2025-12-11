package language

// YAML handles YAML and TOML configuration files
type YAML struct {
	BaseLanguage
}

// NewYAML creates a new YAML/TOML language handler
func NewYAML() *YAML {
	return &YAML{
		BaseLanguage: BaseLanguage{
			name:         "YAML/TOML",
			extensions:   []string{".yaml", ".yml", ".toml"},
			filenames:    []string{},
			commentStyle: CommentStyleLine,
			linePrefix:   "#",
		},
	}
}
