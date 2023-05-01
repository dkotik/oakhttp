package document

type Document struct {
	Name        string
	Description string
	Package     string
	Roles       []RoleDefinition
	SilentRoles []RoleDefinition
	DefaultRole string
	Policies    []PolicyDefinition
}
