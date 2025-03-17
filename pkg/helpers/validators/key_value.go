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
	"regexp"

	"github.com/go-playground/validator/v10"
)

// KeyValue is a validator that checks if a string is a valid key=value pair.
var KeyValue = &CustomValidator{
	Key: "str_key_value",
	Func: func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`^([a-zA-Z0-9_]+=[a-zA-Z0-9_]+)$`).MatchString(fl.Field().String())
	},
}
