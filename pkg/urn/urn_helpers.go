package urn

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Special functions for the terraform provider test.
// TestIsType returns true if the URN is of the specified type.
func TestIsType(urnType URN) resource.CheckResourceAttrWithFunc {
	return func(value string) error {
		if value == "" {
			return nil
		}

		if !URN(value).IsType(urnType) {
			return fmt.Errorf("urn %s is not of type %s", value, urnType)
		}
		return nil
	}
}
