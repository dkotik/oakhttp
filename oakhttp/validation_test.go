package oakhttp

import (
	"errors"
	"testing"
)

func TestInvalidRequestError(t *testing.T) {
	// var err error
	err := NewInvalidRequestError(errors.Join(
		errors.New("first error"),
		errors.New("second error"),
	))

	unwrappable, ok := err.Unwrap().(interface{ Unwrap() []error })
	if !ok {
		t.Fatal("inside error does not implement `Unwrap() []error` method")
	}

	errs := unwrappable.Unwrap()
	if len(errs) != 2 {
		t.Fatal("error count does not match", len(errs), 2)
	}
	// t.Fatalf("%#v", errs)
}
