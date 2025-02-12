/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package urn

import (
	"errors"
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

// FindURNTypeFromString returns the URN type from a string.
func FindURNTypeFromString(value string) (URN, error) {
	if value == "" {
		return "", errors.New("value doesn't contains an URN type provided")
	}

	if u, ok := URNByNames[value]; ok {
		return u, nil
	}

	return "", fmt.Errorf("URN type %s doesn't exist by package urn", value)
}
