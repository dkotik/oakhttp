package oakrouter

import "fmt"

type adaptorError struct {
	name  string
	cause error
}

func (e *adaptorError) Unwrap() error {
	return e.cause
}

func (e *adaptorError) Error() string {
	return fmt.Sprintf("cannot adapt the request %q handler: %v", e.name, e.cause)
}
