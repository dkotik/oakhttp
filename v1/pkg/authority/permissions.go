package authority

import "errors"

type (
	// Checker    func(map[string]string) error
	Comparator func(map[string]string) bool
)

func newChecker(deny, allow []Comparator) func(description map[string]string) error {
	return func(actionAnnotations map[string]string) error {
		for _, comp := range deny {
			if comp(actionAnnotations) {
				return errors.New("denied")
			}
		}
		for _, comp := range allow {
			if comp(actionAnnotations) {
				return nil // this action is allowed
			}
		}
		return errors.New("denied by default")
	}
}
