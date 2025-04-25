/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package validators_test

import (
	"reflect"
	"testing"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/helpers/validators"
)

// validators_test.go test functions
// This file contains the test functions for the validators package.
// It includes tests for the Default validator, which sets default values for fields
// based on their type and the provided parameter.
// The tests cover various scenarios, including string, int, float, and bool types.
// The tests also check for invalid cases where the field type does not match the expected type.

type TestTypeString struct {
	StringField string `validate:"default=default"`
}
type TestTypeInt struct {
	// IntField     int     `/ validate:"default=42"`
	// Int8Field    int8    `validate:"default=42"`
	// Int16Field   int16   `validate:"default=42"`
	// Int32Field   int32   `validate:"default=42"`
	Int64Field int64 `validate:"default=42"`
}
type TestTypeUInt struct {
	// UintField    uint    `validate:"default=42"`
	// Uint8Field   uint8   `validate:"default=42"`
	// Uint16Field  uint16  `validate:"default=42"`
	// Uint32Field  uint32  `validate:"default=42"`
	Uint64Field uint64 `validate:"default=42"`
}
type TestTypeFloat struct {
	// Float32Field float32 `validate:"default=42.0"`
	Float64Field float64 `validate:"default=42.0"`
}
type TestTypeBool struct {
	BoolField bool `validate:"default=true"`
}
type TestTypeInvalidInt struct {
	// IntField     int     `validate:"default=invalid"`
	// Int8Field    int8    `validate:"default=invalid"`
	// Int16Field   int16   `validate:"default=invalid"`
	// Int32Field   int32   `validate:"default=invalid"`
	Int64Field int64 `validate:"default=invalid"`
}
type TestTypeInvalidUInt struct {
	// UintField    uint    `validate:"default=invalid"`
	// Uint8Field   uint8   `validate:"default=invalid"`
	// Uint16Field  uint16  `validate:"default=invalid"`
	// Uint32Field  uint32  `validate:"default=invalid"`
	Uint64Field uint64 `validate:"default=invalid"`
}
type TestTypeInvalidFloat struct {
	// Float32Field float32 `validate:"default=invalid"`
	Float64Field float64 `validate:"default=invalid"`
}
type TestTypeInvalidBool struct {
	BoolField bool `validate:"default=invalid"`
}

type TestTypeNotvalid struct {
	ArrayField []string `validate:"default=[]{\"default\"}"`
	// MapField     map[string]string `validate:"default=invalid"`
	// StructField  TestTypeString `validate:"default=invalid"`
	// InterfaceField interface{} `validate:"default=invalid"`
	// ChanField    chan string `validate:"default=invalid"`
	// FuncField    func() `validate:"default=invalid"`
	// ComplexField complex128 `validate:"default=invalid"`
	// PtrField     *string `validate:"default=invalid"`
}

func TestDefaultValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name        string
		input       interface{}
		expected    interface{}
		expectError bool
	}

	tests := []testCase{
		{
			name: "Valid default string",
			input: &TestTypeString{
				StringField: "",
			},
			expected: &TestTypeString{
				StringField: "default",
			},
			expectError: false,
		},
		{
			name: "Valid default int64",
			input: &TestTypeInt{
				Int64Field: 0,
			},
			expected: &TestTypeInt{
				Int64Field: 42,
			},
			expectError: false,
		},
		{
			name: "Valid default uint64",
			input: &TestTypeUInt{
				Uint64Field: 0,
			},
			expected: &TestTypeUInt{
				Uint64Field: 42,
			},
			expectError: false,
		},
		{
			name: "Valid default float64",
			input: &TestTypeFloat{
				Float64Field: 0.0,
			},
			expected: &TestTypeFloat{
				Float64Field: 42.0,
			},
			expectError: false,
		},
		{
			name: "Valid default bool",
			input: &TestTypeBool{
				BoolField: false,
			},
			expected: &TestTypeBool{
				BoolField: true,
			},
			expectError: false,
		},
		{
			name: "Invalid default int64",
			input: &TestTypeInvalidInt{
				Int64Field: 0,
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "Invalid default bool",
			input: &TestTypeInvalidBool{
				BoolField: false,
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "Invalid default uint64",
			input: &TestTypeInvalidUInt{
				Uint64Field: 0,
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "Invalid default float64",
			input: &TestTypeInvalidFloat{
				Float64Field: 0.0,
			},
			expected:    0.0,
			expectError: true,
		},
		{
			name: "Invalid default not valid",
			input: &TestTypeNotvalid{
				ArrayField: []string{},
			},
			expected: &TestTypeNotvalid{
				ArrayField: []string{},
			},
			expectError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			validate := validators.New()
			err := validate.Struct(test.input)

			// Check if the validation error is expected
			if test.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			// If we expect no error, but got one, fail the test
			if err != nil {
				// If we expect no error, but got one, fail the test
				// and print the error message
				t.Errorf("expected no error, got %v", err)
			}

			// Check if the input matches the expected value
			if !reflect.DeepEqual(test.input, test.expected) {
				t.Errorf("expected %v, got %v", test.expected, test.input)
			}
		})
	}
}
