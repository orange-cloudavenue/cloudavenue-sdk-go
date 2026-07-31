/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package errors

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestTruncateBody(t *testing.T) {
	tests := []struct {
		name   string
		body   string
		maxLen int
		want   string
	}{
		{
			name:   "body shorter than maxLen",
			body:   "short",
			maxLen: 512,
			want:   "short",
		},
		{
			name:   "body exactly maxLen",
			body:   strings.Repeat("x", 512),
			maxLen: 512,
			want:   strings.Repeat("x", 512),
		},
		{
			name:   "body longer than maxLen",
			body:   strings.Repeat("x", 1000),
			maxLen: 512,
			want:   strings.Repeat("x", 512) + "... (truncated)",
		},
		{
			name:   "multi-byte rune split at boundary stays valid UTF-8",
			body:   strings.Repeat("€", 200), // 600 bytes, 3-byte rune
			maxLen: 512,
			want:   strings.Repeat("€", 170) + "... (truncated)", // 510 bytes + suffix
		},
		{
			name:   "maxLen zero",
			body:   strings.Repeat("x", 1000),
			maxLen: 0,
			want:   strings.Repeat("x", 1000),
		},
		{
			name:   "negative maxLen",
			body:   strings.Repeat("x", 1000),
			maxLen: -1,
			want:   strings.Repeat("x", 1000),
		},
		{
			name:   "empty body",
			body:   "",
			maxLen: 512,
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TruncateBody(tt.body, tt.maxLen)

			if got != tt.want {
				t.Errorf("TruncateBody() = %q, want %q", got, tt.want)
			}

			if !utf8.ValidString(got) {
				t.Errorf("TruncateBody() = %q, is not valid UTF-8", got)
			}
		})
	}
}
