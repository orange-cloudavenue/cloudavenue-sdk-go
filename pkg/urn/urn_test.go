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
	"testing"
)

const (
	validUUIDv4 = "12345678-1234-1234-1234-123456789012"
)

func TestURN_ContainsPrefix(t *testing.T) {
	tests := []struct {
		name string
		urn  URN
		want bool
	}{
		{
			name: "ContainsPrefix",
			urn:  URN(VM.String() + validUUIDv4),
			want: true,
		},
		{
			name: "DoesNotContainPrefix",
			urn:  URN("urn:vm:" + validUUIDv4),
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			urn:  URN(""),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.urn.ContainsPrefix(); got != tt.want {
				t.Errorf("URN.ContainsPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isURNV4(t *testing.T) {
	type args struct {
		urn string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ValidURN",
			args: args{
				urn: validUUIDv4,
			},
			want: true,
		},
		{
			name: "InvalidURN",
			args: args{
				urn: "f47ac10b-58cddc-43-a567-0e02b2c3d4791",
			},
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			args: args{
				urn: "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isUUIDV4(tt.args.urn); got != tt.want {
				t.Errorf("isUUIDV4() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestURN_IsType(t *testing.T) {
	type args struct {
		prefix URN
	}
	tests := []struct {
		name string
		urn  URN
		args args
		want bool
	}{
		{
			name: "IsType",
			urn:  URN(VM.String() + validUUIDv4),
			args: args{
				prefix: VM,
			},
			want: true,
		},
		{
			name: "IsNotType",
			urn:  URN(VM.String() + validUUIDv4),
			args: args{
				prefix: User,
			},
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			urn:  URN(""),
			args: args{
				prefix: VM,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.urn.IsType(tt.args.prefix); got != tt.want {
				t.Errorf("URN.IsType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractURNv4(t *testing.T) {
	type args struct {
		urn    string
		prefix URN
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ExtractURN",
			args: args{
				urn:    VM.String() + validUUIDv4,
				prefix: VM,
			},
			want: validUUIDv4,
		},
		{
			name: "EmptyString",
			args: args{
				urn:    "",
				prefix: VM,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractUUIDv4(tt.args.urn, tt.args.prefix); got != tt.want {
				t.Errorf("extractURNv4() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	type args struct {
		urn string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ValidURN",
			args: args{
				urn: VM.String() + validUUIDv4,
			},
			want: true,
		},
		{
			name: "InvalidURN",
			args: args{
				urn: "f47ac10b-58cddc-43-a567-0e02b2c3d4791",
			},
			want: false,
		},
		{
			name: "InvalidPrefix",
			args: args{
				urn: "urn:vm:f47ac10b-58cddc-43-a567-0e02b2c3d4791",
			},
			want: false,
		},
		{ // Empty string
			name: "EmptyString",
			args: args{
				urn: "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValid(tt.args.urn); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsType tests the TestIsType function.
func TestTestIsType(t *testing.T) {
	testCases := []struct {
		name    string
		urnType URN
		urn     URN
		want    bool
	}{
		{
			name:    "valid urn",
			urnType: VM,
			urn:     URN(VM.String() + validUUIDv4),
			want:    true,
		},
		{
			name:    "invalid urn",
			urnType: VM,
			urn:     "invalid-urn",
			want:    false,
		},
		{
			name:    "empty value",
			urnType: VM,
			urn:     "",
			want:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := TestIsType(tc.urnType)(tc.urn.String())
			if tc.want && err != nil {
				t.Errorf("TestIsType() = %v, want %v", err, tc.want)
			}
		})
	}
}

func TestIsUUIDV4(t *testing.T) {
	tests := []struct {
		name string
		urn  string
		want bool
	}{
		{
			name: "ValidUUIDV4",
			urn:  validUUIDv4,
			want: true,
		},
		{
			name: "InvalidUUIDV4",
			urn:  "12345678-1234-1234-1234-12345678901Z",
			want: false,
		},
		{
			name: "InvalidFormat",
			urn:  "12345678-1234-1234-1234-1234567890123",
			want: false,
		},
		{
			name: "EmptyString",
			urn:  "",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUUIDV4(tt.urn); got != tt.want {
				t.Errorf("IsUUIDV4() = %v, want %v", got, tt.want)
			}
		})
	}
}
