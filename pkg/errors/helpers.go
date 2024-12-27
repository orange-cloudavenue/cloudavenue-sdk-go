package errors

import "errors"

// IsNotFound - Returns true if the error contains "not found".
func IsNotFound(e error) bool {
	return errors.Is(e, ErrNotFound)
}
