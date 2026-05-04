/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package urn

import "testing"

const (
	testEmptyStringName    = "EmptyString"
	testInvalidVMURN       = "urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d479"
	testInvalidVMURNLong   = "urn:vcloud:vm:f47ac10b-58cc-4372-a567-0e02b2c3d4791"
	testInvalidUserURN     = "urn:vcloud:user:f47ac10b-58cc-4372-a567-0e02b2c3d479"
	testInvalidUserURNLong = "urn:vcloud:user:f47ac10b-58cc-4372-a567-0e02b2c3d4791"
	testExtractUUID        = "123e4567-e89b-12d3-a456-426614174000"
)

func TestNormalize(t *testing.T) {
	type args struct {
		prefix URN
		uuid   string
	}
	tests := []struct {
		name string
		args args
		want URN
	}{
		{
			name: "Normalize",
			args: args{
				prefix: VM,
				uuid:   validUUIDv4,
			},
			want: URN(VM.String() + validUUIDv4),
		},
		// Check prefix is empty
		{
			name: "EmptyPrefix",
			args: args{
				prefix: "",
				uuid:   validUUIDv4,
			},
			want: "",
		},
		// Check uuid is already an URN
		{
			name: "AlreadyURN",
			args: args{
				prefix: VM,
				uuid:   VM.String() + validUUIDv4,
			},
			want: URN(VM.String() + validUUIDv4),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Normalize(tt.args.prefix, tt.args.uuid); got != tt.want {
				t.Errorf("Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractUUID(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "ValidUUID",
			input: "This is a string with a UUID " + testExtractUUID + " inside",
			want:  testExtractUUID,
		},
		{
			name:  "NoUUID",
			input: "This string does not contain a UUID",
			want:  "",
		},
		{
			name:  "ValidUUIDInURN",
			input: "urn:vcloud:vm:" + testExtractUUID,
			want:  testExtractUUID,
		},
		{
			name:  "ValidUUID in URL",
			input: "https://example.com/" + testExtractUUID,
			want:  testExtractUUID,
		},
		{
			name:  "ValidUUID in URL with URN format",
			input: "https://example.com/urn:vcloud:vm:" + testExtractUUID,
			want:  testExtractUUID,
		},
		{
			name:  "ValidUUID in URL with URN format and query",
			input: "https://example.com/urn:vcloud:vm:" + testExtractUUID + "?page=10",
			want:  testExtractUUID,
		},
		{
			name:  "InvalidUUIDFormat",
			input: "This string contains an invalid UUID 123e4567-e89b-12d3-a456-42661417400Z",
			want:  "",
		},
		{
			name:  testEmptyStringName,
			input: "",
			want:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractUUID(tt.input); got != tt.want {
				t.Errorf("ExtractUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}
