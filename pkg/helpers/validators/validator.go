package validators

import "github.com/go-playground/validator/v10"

// New creates a new validator.
func New() *validator.Validate {
	v := validator.New()
	_ = v.RegisterValidation(DisallowUpper.Key, DisallowUpper.Func)
	_ = v.RegisterValidation(DisallowSpace.Key, DisallowSpace.Func)

	return v
}
