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
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type APIErrorResponse struct {
	Code    string `json:"code"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

func (e *APIErrorResponse) UnmarshalJSON(data []byte) error {
	type Alias APIErrorResponse

	aux := &struct {
		Code json.RawMessage `json:"code"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// Try string first.
	var s string
	if err := json.Unmarshal(aux.Code, &s); err == nil {
		e.Code = s
		return nil
	}

	// Then try integer.
	var i int64
	if err := json.Unmarshal(aux.Code, &i); err == nil {
		e.Code = strconv.FormatInt(i, 10)
		return nil
	}

	return fmt.Errorf("invalid code field: %s", aux.Code)
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
