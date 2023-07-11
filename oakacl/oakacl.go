package oakacl

import (
	"github.com/dkotik/oakacs/oakpolicy"
	"golang.org/x/exp/slog"
)

type ACL struct {
	loggerAllowLevel            slog.Level
	loggerDenyLevel             slog.Level
	loggerAllow                 *slog.Logger
	loggerDeny                  *slog.Logger
	contextPermissionsExtractor ContextPermissionsExtractor
}

func (acl *ACL) match(ps []Permission, a oakpolicy.Action, dpath oakpolicy.DomainPath) bool {
	// for _, p := range ps {
	// 	if p.Action.Match(a) && dpath.Match(p.ResourceMask) {
	// 		return true
	// 	}
	// }
	return false
}

// func (acl *ACL) IsAuthorized(ctx context.Context) oakpolicy.AuthorizedContext {
// 	permissions, _ := PermissionsFromContext(ctx)
// 	return authorizedContext(permissions)
// }

// func (acl *ACL) Authorize(ctx context.Context, a oakpolicy.Action, r oakpolicy.Resource) error {
// 	permissions, err := acl.contextPermissionsExtractor(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	if !acl.match(permissions, a, r.DomainPath()) {
// 		return oakpolicy.Deny
// 	}
// 	return nil
// }

// AllowActionsForResourcesMatching authorizes any action from the provided list to any resource matching provided masks. This a helper method for debugging. Prefer to use generated policies.
// func AllowActionsForResourcesMatching(actions []Action, resourceMasks [][]string) Policy {
// 	return func(ctx context.Context, i Intention) error {
// 		for _, resourceMask := range resourceMasks {
// 			if i.ResourcePath().Match(resourceMask...) {
// 				if i.Action().In(actions...) {
// 					return Allow
// 				}
// 			}
// 		}
// 		return nil
// 	}
// }
