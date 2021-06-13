package oakacs

import "context"

// Registration creates an Entity within ACS backend.
type Registration func(context.Context, Identity) error
