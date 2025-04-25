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
	"reflect"
	"strconv"

	"github.com/go-playground/validator/v10"
)

// URN is a validator that checks if a string is a valid URN (Uniform Resource Name).
var Default = &CustomValidator{
	Key: "default",
	Func: func(fl validator.FieldLevel) bool {
		if !fl.Field().IsZero() {
			return true
		}

		k := fl.Field().Type().Kind()

		switch k {
		case reflect.String:
			fl.Field().SetString(fl.Param())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			i, err := strconv.Atoi(fl.Param())
			if err != nil {
				return false
			}
			fl.Field().SetInt(int64(i))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			i, err := strconv.ParseUint(fl.Param(), 10, 64)
			if err != nil {
				return false
			}
			fl.Field().SetUint(i)
		case reflect.Float32, reflect.Float64:
			f, err := strconv.ParseFloat(fl.Param(), 64)
			if err != nil {
				return false
			}
			fl.Field().SetFloat(f)
		case reflect.Bool:
			b, err := strconv.ParseBool(fl.Param())
			if err != nil {
				return false
			}
			fl.Field().SetBool(b)
		default:

			// If the field is not a string, int, or float, we can't set a default value
			// and we should return false to indicate that the validation failed.
			// Set the field to its nil value.
			// fl.Field().Set(reflect.Zero(fl.Field().Type()))
			return false
		}

		return true
	},
}
