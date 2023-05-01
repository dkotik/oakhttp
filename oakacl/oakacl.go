package oakacl

import "context"

type ACL struct {
	contextPermissionsExtractor ContextPermissionsExtractor
}

// AllowActionsForResourcesMatching authorizes any action from the provided list to any resource matching provided masks. This a helper method for debugging. Prefer to use generated policies.
func AllowActionsForResourcesMatching(actions []Action, resourceMasks [][]string) Policy {
	return func(ctx context.Context, i Intention) error {
		for _, resourceMask := range resourceMasks {
			if i.ResourcePath().Match(resourceMask...) {
				if i.Action().In(actions...) {
					return Allow
				}
			}
		}
		return nil
	}
}
