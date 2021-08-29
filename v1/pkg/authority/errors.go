package authority

// ErrAccessDenied happens when access to a resource is denied for any reason.
type ErrAccessDenied struct {
	Role   string
	Action string
	Cause  error
}

func (e *ErrAccessDenied) Unwrap() error {
	return e.Cause
}

func (e *ErrAccessDenied) Error() string { return "access denied" }
