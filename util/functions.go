package util

// Coalesce returns the dereferenced value or defaultVal if ptr is nil
func Coalesce[T any](ptr *T, defaultVal T) T {
	if ptr != nil {
		return *ptr
	}
	return defaultVal

}
