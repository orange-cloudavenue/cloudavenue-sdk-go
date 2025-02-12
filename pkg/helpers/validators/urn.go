package validators

import (
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

// URN is a validator that checks if a string is a valid URN (Uniform Resource Name).
var URN = &CustomValidator{
	Key: "urn",
	Func: func(fl validator.FieldLevel) bool {
		fl.Param()

		u, err := urn.FindURNTypeFromString(fl.Param())
		if err != nil {
			return false
		}

		return strings.Contains(fl.Field().String(), u.String())
	},
}
