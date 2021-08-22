package oakacs

import (
	"context"
	"fmt"
)

// TODO: I should treat all secrets as tokens, hashed passwords should just present a token-like API.
// TODO: AccessControlSystem Authenticate should also indicate the METHOD OF AUTH! and be matched with an authenticator. JUST ONE TOKEN PER AUTHENTICATOR?

type Authenticator interface {
	Generate(ctx context.Context, tokenOrPassword string) (*Secret, error)
	Compare(ctx context.Context, tokenOrPassword string, secret []*Secret) error
}

// Authenticate matches provided user and password to an indentity.
func (acs *AccessControlSystem) Authenticate(ctx context.Context, user, tokenOrPassword, authenticator string) (s *Session, err error) {
	defer func() {
		event := Event{
			ctx:  ctx,
			Type: EventTypeAuthenticationSuccess,
		}
		if err != nil {
			err = fmt.Errorf("authentication error: %w", err)
			event.Error = err
			event.Type = EventTypeAuthenticationFailure
		}
		acs.Broadcast(event)
	}()
	// throttle attempts by user

	// auth, ok := acs.authenticators[authenticator]
	// if !ok {
	// 	return nil, errors.New("chosen authenticator is not active")
	// }

	// // retrieve identity
	// identity, err := acs.persistent.RetrieveIdentity(ctx, user)
	// if err != nil {
	// 	// modulate time here to avoid betraying proof of existance
	// 	// compare to a random string to modulate?
	// 	return nil, err
	// }
	//
	// secrets, err := acs.persistent.RetrieveSecret(ctx, identity.UUID, authenticator)
	// if err != nil {
	// 	return nil, err
	// }
	// // if secret.Expires.After(time.Now()) {
	// // 	return nil, errors.New("existing access token is expired")
	// // }

	// if err = auth.Compare(ctx, tokenOrPassword, secrets); err != nil {
	// 	// TODO: this is a deep error, should it be a different type?
	// 	// for time-modulation?
	// 	return nil, err
	// }
	// // secret.Used = time.Now()
	// // if err = acs.backend.UpdateSecret(ctx, secret); err != nil {
	// // 	return
	// // }

	// start session
	// attach role to session
	// issue event
	return
}
