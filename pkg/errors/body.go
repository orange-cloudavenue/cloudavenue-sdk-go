/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package errors

import "strings"

// MaxErrorBodyLen caps the amount of raw response body embedded in a
// fallback error message, to avoid unbounded error messages / log bloat when
// upstream returns large error pages (e.g. HTML gateway pages).
const MaxErrorBodyLen = 512

// TruncateBody truncates body to at most maxLen bytes while preserving valid
// UTF-8 (stripping any partial/invalid trailing rune left by the byte-slice
// truncation, and any invalid sequences within the retained prefix, via
// strings.ToValidUTF8), appending "... (truncated)" when the body was cut. The
// body is returned unchanged when len(body) <= maxLen or when maxLen <= 0.
func TruncateBody(body string, maxLen int) string {
	if maxLen <= 0 || len(body) <= maxLen {
		return body
	}

	// strings.ToValidUTF8 replaces every invalid byte sequence in the
	// retained prefix — including the partial/invalid rune left dangling by
	// the byte-slice truncation — instead of producing garbled replacement
	// characters for non-ASCII upstream bodies.
	return strings.ToValidUTF8(body[:maxLen], "") + "... (truncated)"
}
