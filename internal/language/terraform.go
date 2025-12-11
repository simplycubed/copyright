package language

// Terraform handles Terraform and HCL files
type Terraform struct {
	BaseLanguage
}

// NewTerraform creates a new Terraform language handler
func NewTerraform() *Terraform {
	return &Terraform{
		BaseLanguage: BaseLanguage{
			name:         "Terraform",
			extensions:   []string{".tf", ".tfvars", ".hcl"},
			filenames:    []string{},
			commentStyle: CommentStyleLine,
			linePrefix:   "#",
		},
	}
}

// Note: Terraform supports both # and /* */ style comments
// We use # by default as it's more common in Terraform configurations
