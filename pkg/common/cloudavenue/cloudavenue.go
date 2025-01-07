/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package commoncloudavenue

import (
	"fmt"
	"strings"
)

type APIErrorResponse struct {
	Code    string `json:"code"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

// FormatError - Formats the error.
func (e *APIErrorResponse) FormatError() string {
	return fmt.Sprintf("ErrorCode:%s - ErrorReason:%s - ErrorMessage:%s", e.Code, e.Reason, e.Message)
}

// ToError - Converts an APIErrorResponse to an error.
func ToError(e *APIErrorResponse) error {
	return fmt.Errorf("error on API call: %s", e.FormatError())
}

// IsNotFound - Returns true if the error is a 404.
func IsNotFound(e error) bool {
	return strings.Contains(e.Error(), "ErrorCode:404")
}
