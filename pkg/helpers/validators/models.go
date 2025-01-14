package validators

import "github.com/go-playground/validator/v10"

type (
	CustomValidator struct {
		Key  string
		Func validator.Func
	}
)
