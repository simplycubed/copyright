package language

// SQL handles SQL files
type SQL struct {
	BaseLanguage
}

// NewSQL creates a new SQL language handler
func NewSQL() *SQL {
	return &SQL{
		BaseLanguage: BaseLanguage{
			name:         "SQL",
			extensions:   []string{".sql"},
			filenames:    []string{},
			commentStyle: CommentStyleSQL,
			linePrefix:   "--",
		},
	}
}
