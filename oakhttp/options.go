package oakhttp

// TODO: rewrite with dual-capability options.

type ClientOption interface {
	applyToClientOptions(*clientOptions) error
}

type ServerOption interface {
	applyToServerOptions(*serverOptions) error
}
