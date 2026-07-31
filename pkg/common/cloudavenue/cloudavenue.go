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
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"

	caverrors "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/errors"
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

// FormatError - Formats the error, omitting any fields that are empty.
// Returns an empty string if none of the fields are set.
func (e *APIErrorResponse) FormatError() string {
	var parts []string
	if e.Code != "" {
		parts = append(parts, fmt.Sprintf("ErrorCode:%s", e.Code))
	}
	if e.Reason != "" {
		parts = append(parts, fmt.Sprintf("ErrorReason:%s", e.Reason))
	}
	if e.Message != "" {
		parts = append(parts, fmt.Sprintf("ErrorMessage:%s", e.Message))
	}

	return strings.Join(parts, " - ")
}

// apiCallError carries the HTTP status code alongside the formatted message
// so callers (e.g. IsNotFound) can check the status structurally instead of
// substring-matching the final error string, which may embed untrusted raw
// response bodies.
type apiCallError struct {
	statusCode int
	message    string
}

func (e *apiCallError) Error() string {
	return e.message
}

// ToError - Converts a resty response into an error.
// It prefers the structured APIErrorResponse fields when available, falling
// back to the raw HTTP status and response body when the typed error has
// nothing usable (e.g. plain-text or HTML error bodies from rate-limiting
// or upstream gateways). The returned error always carries the actual HTTP
// status code for structural checks (see IsNotFound).
func ToError(r *resty.Response) error {
	statusCode := r.StatusCode()

	apiErr, _ := r.Error().(*APIErrorResponse)
	if apiErr != nil {
		if formatted := apiErr.FormatError(); formatted != "" {
			return &apiCallError{statusCode: statusCode, message: formatted}
		}
	}

	body := strings.TrimSpace(r.String())
	if body == "" {
		return &apiCallError{statusCode: statusCode, message: fmt.Sprintf("HTTPCode:%s", r.Status())}
	}

	body = caverrors.TruncateBody(body, caverrors.MaxErrorBodyLen)

	return &apiCallError{statusCode: statusCode, message: fmt.Sprintf("HTTPCode:%s - body: %s", r.Status(), body)}
}

// IsNotFound - Returns true if the error is a 404.
func IsNotFound(e error) bool {
	if e == nil {
		return false
	}

	var apiErr *apiCallError
	if errors.As(e, &apiErr) {
		return apiErr.statusCode == 404
	}

	// Legacy fallback for errors not produced via ToError.
	return strings.Contains(e.Error(), "ErrorCode:404")
}
