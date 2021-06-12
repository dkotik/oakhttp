package oakacs

// Action represents something an Identity can do.
type Action struct {
	Name   Name
	Verb   Name
	Target Name
}
