package cueroles

import "fmt"

func MethodExact(v interface{}, p string) error {
	v, ok := v.(string)
	if ok && v == p {
		return nil
	}
	return fmt.Errorf("%v does not exactly match %q", v, p)
}

// func MethodExact(attribute, against string) Method {
// 	return func(annotation map[string]string) bool {
// 		value, ok := annotation[attribute]
// 		if !ok {
// 			return false
// 		}
// 		return value == against
// 	}
// }

// func MethodOlderThan(v interface{}, p string) error {
// 	t, ok := v.(time.Time)
// 	if !ok {
// 		return errors.New("attribute type mistmatch error")
// 	}
// 	// parse p into duration
// 	// compare
// 	return nil
// }
