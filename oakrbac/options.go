package oakrbac

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/exp/slog"
)

type options struct {
	roleRepository       RoleRepository
	contextRoleExtractor ContextRoleExtractor
	roles                []Role
	listeners            []Listener
}

type Option func(*options) error

func WithOptions(withOptions ...Option) Option {
	return func(o *options) (err error) {
		for _, option := range withOptions {
			if err = option(o); err != nil {
				return err
			}
		}
		return nil
	}
}

func WithDefaultOptions() Option {
	return func(o *options) (err error) {
		defer func() {
			if err != nil {
				err = fmt.Errorf("could not set a default option: %w", err)
			}
		}()

		if o.roleRepository == nil {
			if err = WithListRoleRepository()(o); err != nil {
				return err
			}
		}

		if o.contextRoleExtractor == nil {
			if err = WithContextRoleExtractor(RoleNameFromContext)(o); err != nil {
				return err
			}
		}

		if len(o.listeners) == 0 {
			if err = WithAuthorizationGrantLogger(
				slog.New(slog.NewTextHandler(os.Stderr)),
				slog.LevelInfo,
			)(o); err != nil {
				return err
			}
		}

		return nil
	}
}

func WithRoleRepository(repo RoleRepository) Option {
	return func(o *options) error {
		if o.roleRepository != nil {
			return errors.New("role repository is already set")
		}
		if repo == nil {
			return errors.New("cannot use a <nil> role repository")
		}
		o.roleRepository = repo
		return nil
	}
}

func WithListRoleRepository() Option {
	return WithRoleRepository(&ListRepository{})
}

func WithMapRoleRepository() Option {
	return WithRoleRepository(&MapRepository{})
}

func WithContextRoleExtractor(extractor ContextRoleExtractor) Option {
	return func(o *options) error {
		if o.contextRoleExtractor != nil {
			return errors.New("context role extractor is already set")
		}
		if extractor == nil {
			return errors.New("cannot use a <nil> context role extractor")
		}
		o.contextRoleExtractor = extractor
		return nil
	}
}

func WithCustomRole(r Role) Option {
	return func(o *options) error {
		if r == nil {
			return errors.New("cannot use a <nil> role")
		}
		name := r.Name()
		if name == "" {
			return errors.New("cannot use an empty role name")
		}
		// for _, role := range o.roles {
		// 	if name == role.Name() {
		// 		return fmt.Errorf("role names must be unique: %s", name)
		// 	}
		// }
		o.roles = append(o.roles, r)
		return nil
	}
}

func WithRole(name string, ps ...Policy) Option {
	return func(o *options) error {
		if len(ps) == 0 {
			return fmt.Errorf("role %q must incluse at least one policy", name)
		}
		for _, policy := range ps {
			if policy == nil {
				return fmt.Errorf("role %q cannot use a <nil> policy", name)
			}
		}
		return WithCustomRole(&basicRole{
			name:     name,
			policies: ps,
		})(o)
	}
}

func WithOmnipotentRole(name string) Option {
	return WithCustomRole(&omnipotentRole{name: name})
}

func WithImpotentRole(name string) Option {
	return WithCustomRole(&impotentRole{name: name})
}

func WithListener(l Listener) Option {
	return func(o *options) error {
		if l == nil {
			return errors.New("cannot use a <nil> event listener")
		}
		o.listeners = append(o.listeners, l)
		return nil
	}
}
