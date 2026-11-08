package language

// Dockerfile handles Dockerfile and Containerfile files
type Dockerfile struct {
	BaseLanguage
}

// NewDockerfile creates a new Dockerfile language handler
func NewDockerfile() *Dockerfile {
	return &Dockerfile{
		BaseLanguage: BaseLanguage{
			name:         "Dockerfile",
			extensions:   []string{".dockerfile"},
			filenames:    []string{"Dockerfile", "Containerfile", "dockerfile", "containerfile"},
			commentStyle: CommentStyleLine,
			linePrefix:   "#",
		},
	}
}
