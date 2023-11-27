package utils

// ToPTR is a generic function that takes a value and returns a pointer to it.
func ToPTR[T any](v T) *T {
	return &v
}
