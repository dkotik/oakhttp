package oakacl

import "context"

type Permission struct {
	ResourceMask []string
	Action       oakpolicy.Action
}

func (acl *ACL) NewPolicy() oakpolicy.Policy {
	return func(ctx context.Context, a oakpolicy.Action, r oakpolicy.Resource) error {
		permissions, err := acl.contextPermissionsExtractor(ctx)
		if err != nil {
			return err
		}
		authorizationPath := r.ResourcePath()
		for _, permission := range permissions {
			if !authorizationPath.Match(permission.ResourceMask) {
				continue
			}
			if permission.Action != a && permission.Action != oakpolicy.ActionAny {
				continue
			}
			return oakpolicy.Allow
		}
		return oakpolicy.Deny
	}
}
