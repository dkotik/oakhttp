package oakacs

// ErrAuthorization happens when access to a resource is denied for any reason.
type ErrAuthorization struct {
	Action string
	Cause  error
}

func (e *ErrAuthorization) Unwrap() error {
	return e.Cause
}

func (e *ErrAuthorization) Error() string { return "access denied" }
