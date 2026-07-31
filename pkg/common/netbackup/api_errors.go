/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package commonnetbackup

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/errors"
)

// APIError is the error body shape returned by the Netbackup API.
type APIError []struct {
	PropertyName string `json:"PropertyName"`
	Message      string `json:"Message"`
}

// FormatError - Formats the Netbackup error, omitting empty entries.
// Each entry is rendered as "PropertyName: Message" (or just "Message" when
// PropertyName is empty), multiple entries joined by "; ". Returns an empty
// string when there are no usable entries.
func (e APIError) FormatError() string {
	var parts []string
	for _, entry := range e {
		switch {
		case entry.PropertyName != "" && entry.Message != "":
			parts = append(parts, fmt.Sprintf("%s: %s", entry.PropertyName, entry.Message))
		case entry.Message != "":
			parts = append(parts, entry.Message)
		case entry.PropertyName != "":
			parts = append(parts, entry.PropertyName)
		}
	}

	return strings.Join(parts, "; ")
}

// apiCallError carries the HTTP status code alongside the formatted message
// so callers can check the status structurally instead of substring-matching
// the final error string, which may embed untrusted raw response bodies.
// Mirrors commoncloudavenue.apiCallError; kept as a small unexported
// duplicate here rather than exported/shared to avoid introducing a
// cross-package dependency between commonnetbackup and commoncloudavenue.
type apiCallError struct {
	statusCode int
	message    string
}

func (e *apiCallError) Error() string {
	return e.message
}

// ToError - Converts a resty response into an error.
// It prefers the structured APIError entries when available, falling back
// to the raw HTTP status and response body when the typed error has
// nothing usable (e.g. plain-text or HTML error bodies from rate-limiting
// or upstream gateways).
func ToError(r *resty.Response) error {
	statusCode := r.StatusCode()

	apiErr, _ := r.Error().(*APIError)
	if apiErr != nil {
		if formatted := apiErr.FormatError(); formatted != "" {
			return &apiCallError{statusCode: statusCode, message: formatted}
		}
	}

	body := strings.TrimSpace(r.String())
	if body == "" {
		return &apiCallError{statusCode: statusCode, message: fmt.Sprintf("HTTPCode:%s", r.Status())}
	}

	body = errors.TruncateBody(body, errors.MaxErrorBodyLen)

	return &apiCallError{statusCode: statusCode, message: fmt.Sprintf("HTTPCode:%s - body: %s", r.Status(), body)}
}
