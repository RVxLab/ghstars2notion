package utils

func MapGetKeys[K comparable, V any](m map[K]V) []K {
	var keys []K

	for k, _ := range m {
		keys = append(keys, k)
	}

	return keys
}
