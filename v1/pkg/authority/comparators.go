package authority

func ComparatorExact(attribute, against string) Comparator {
	return func(annotation map[string]string) bool {
		value, ok := annotation[attribute]
		if !ok {
			return false
		}
		return value == against
	}
}
