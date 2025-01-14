/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package validators

import (
	"unicode"

	"github.com/go-playground/validator/v10"
)

var (

	// DisallowUpper is a validator that disallows uppercase characters.
	DisallowUpper = &CustomValidator{
		Key: "disallow_upper",
		Func: func(fl validator.FieldLevel) bool {
			for _, r := range fl.Field().String() {
				if unicode.IsUpper(r) {
					return false
				}
			}
			return true
		},
	}

	// DisallowSpace is a validator that disallows spaces.
	DisallowSpace = &CustomValidator{
		Key: "disallow_space",
		Func: func(fl validator.FieldLevel) bool {
			for _, r := range fl.Field().String() {
				if unicode.IsSpace(r) {
					return false
				}
			}
			return true
		},
	}
)
