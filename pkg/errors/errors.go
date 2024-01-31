package errors

import (
	"errors"
	"fmt"
)

var (

	// * Generic
	ErrNotFound      = errors.New("not found")
	ErrEmpty         = errors.New("empty")
	ErrInvalidFormat = errors.New("invalid format")

	// * Client

	ErrConfigureVmwareClient       = errors.New("unable to configure vmware client")
	ErrOrganizationFormatIsInvalid = fmt.Errorf("organization has an %w", ErrInvalidFormat)
)
