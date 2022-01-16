package oakacs

import "context"

func (acs *AccessControlSystem) Register(
	ctx context.Context,
	request *AuthenticationRequest,
) (*Session, error) {
	// sessions should be promiscious
	// use session values to keep registration / authentication state
	return nil, nil
}
