package backend

// Ephemeral backend tracks sessions and tokens.
type Ephemeral interface {
	// CRUD Session
	// CRUD Tokens?
	// Open(ctx context.Context, who *oakacs.Identity, as *oakacs.Role) (sessionUUID string, err error)
	// Close(ctx context.Context, sessionUUID string) (err error)
	// Authorize(context.Context, string, *oakacs.Permission) (err error)
}

// Persistent backend tracks identities, groups, and roles.
type Persistent interface {
	// CRUD Identity
	// CRUD Group
	// CRUD Role?
}
