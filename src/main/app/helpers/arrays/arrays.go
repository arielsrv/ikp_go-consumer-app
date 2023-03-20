package arrays

func IsEmpty[T comparable](elements []T) bool {
	return len(elements) == 0
}

func Contains[T comparable](elements []T, target T) bool {
	for _, element := range elements {
		if target == element {
			return true
		}
	}
	return false
}

func Find[T comparable](elements []T, predicate func(element T) bool) *T {
	for _, element := range elements {
		if predicate(element) {
			return &element
		}
	}
	return nil
}
