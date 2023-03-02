package oakhttp

import "fmt"

func ExampleUnwrap() {
	err := fmt.Errorf("some execution context: %w", NewNotFoundError("resource"))

	if unwrapped, ok := Unwrap(err); ok {
		fmt.Println("Message:", unwrapped.Error())
		fmt.Println("StatusCode:", unwrapped.HTTPStatusCode())
	}

	// Output:
	// Message: resource "resource" was not found
	// StatusCode: 404
}
