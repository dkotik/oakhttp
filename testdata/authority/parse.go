package cueroles

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

/*
ACTION JUST URL-ENCODED?


? domain~glob=*.truthonly.com & domain~regex: glob:*.truthonly.com domain: glob:*.truthonly.com

deny [
	domain: glob:*.truthonly.com
	domain: glob:*.truthonly.com
	domain: glob:*.truthonly.com
]

Match: {
	bydomain: domain=glob:*.truthonly.com
}

[Allow]: [
	[bydomain, byuser, bydate],
]

Deny: [
	[...]
]


*/

/* CUE

#Match: {
	attribute: string
	operation: string
	pattern: string
}

#Role: {
	name: string
	allow: [...#Match]
	deny: [...#Match]
}

*/

// [service=oakacs]
// domain: *.truthonly.com
// domain *= truthonly test ~= regex test @= tag match
func stringToComparator(s string) (Comparator, error) {
	key, comparator, value := &strings.Builder{}, &strings.Builder{}, &strings.Builder{}
	current := key
	for _, c := range s {
		switch c {
		case ':':
			if current == key {
				current = comparator
				continue
			}
		case '=':
			if current != value {
				current = value
				continue
			}
		}
		if !unicode.IsSpace(c) {
			current.WriteRune(c)
		}
	}

	if key.Len() == 0 {
		return nil, fmt.Errorf("condition %q does not contain a valid key value", s)
	}

	switch comparator.String() {
	case "glob":
		return nil, errors.New("glob comparator is not yet implemented")
	case "regex":
		return nil, errors.New("regex comparator is not yet implemented")
		// case "":
		// 	return ComparatorExact(key.String(), value.String()), nil
	}
	return nil, fmt.Errorf("comparator %q is not registered", comparator.String())
}

// type jsonCondition struct {
// 	Key        string
// 	Comparator string
// 	Match      string
// }
//
// type jsonPermission struct {
// 	Deny  []jsonCondition
// 	Allow []jsonCondition
// }
//
// func jsonPermissionToChecker(p jsonPermission) func(map[string]string) error {
//     return func(map[string]string) error {
//
//     }
// }
