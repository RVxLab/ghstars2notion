package utils

func Map[TValue any](slice []TValue, mappingFunc func(value TValue) TValue) []TValue {
	mapped := make([]TValue, len(slice))

	for i, value := range slice {
		mapped[i] = mappingFunc(value)
	}

	return mapped
}
