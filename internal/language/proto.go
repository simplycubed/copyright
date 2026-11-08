package language

// ProtocolBuffers handles Protocol Buffer files
type ProtocolBuffers struct {
	BaseLanguage
}

// NewProtocolBuffers creates a new Protocol Buffers language handler
func NewProtocolBuffers() *ProtocolBuffers {
	return &ProtocolBuffers{
		BaseLanguage: BaseLanguage{
			name:         "Protocol Buffers",
			extensions:   []string{".proto"},
			filenames:    []string{},
			commentStyle: CommentStyleLine,
			linePrefix:   "//",
		},
	}
}
