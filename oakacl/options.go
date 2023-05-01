package oakacl

import "golang.org/x/exp/slog"

type options struct {
	logger                      *slog.Logger
	logAllowedActionsOnly       bool
	contextPermissionsExtractor ContextPermissionsExtractor
}
