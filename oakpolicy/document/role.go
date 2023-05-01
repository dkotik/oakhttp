package document

import (
	_ "embed"
	"html/template"
	"io"
)

var (
	//go:embed role.tmpl
	roleTemplateSource string

	roleTemplate = template.Must(template.New("role template").Parse(roleTemplateSource))
)

type RoleDefinition struct {
	Name        string
	Description string
	Policies    []string
}

func (r *RoleDefinition) Validate() error {
	return nil
}

func (r *RoleDefinition) Render(w io.Writer) error {
	return roleTemplate.Execute(w, r)
}
