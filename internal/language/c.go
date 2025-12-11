package language

// C handles C and C++ source files
type C struct {
	BaseLanguage
}

// NewC creates a new C/C++ language handler
func NewC() *C {
	return &C{
		BaseLanguage: BaseLanguage{
			name: "C/C++",
			extensions: []string{
				".c", ".h",
				".cpp", ".hpp",
				".cc", ".hh",
				".cxx", ".hxx",
				".c++", ".h++",
			},
			filenames:    []string{},
			commentStyle: CommentStyleBlock,
			blockStart:   "/*",
			blockPrefix:  " *",
			blockEnd:     " */",
		},
	}
}
