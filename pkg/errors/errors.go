package errors

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrEmpty         = errors.New("empty")
	ErrInvalidFormat = errors.New("invalid format")

	ErrOrganizationFormatIsInvalid = fmt.Errorf("organization has an %w", ErrInvalidFormat)
)
