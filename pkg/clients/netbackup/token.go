/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package clientnetbackup

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/errors"
)

// maxErrorBodyLen caps the amount of raw response body embedded in a
// fallback error message, mirroring commoncloudavenue.ToError, to avoid
// unbounded error messages / log bloat when the Netbackup auth endpoint
// returns large error pages (e.g. HTML gateway pages) instead of its usual
// JSON error body.
const maxErrorBodyLen = 512

type token struct {
	baererToken string
	expiresAt   time.Time

	username string
	password string

	endpoint string

	debug bool
}

// IsExpired - Returns true if the token is expired.
func (t *token) IsExpired() bool {
	return t.expiresAt.Before(time.Now())
}

// IsSet - Returns true if the token is set.
func (t *token) IsSet() bool {
	return t.baererToken != ""
}

// GetToken - Returns the token.
func (t *token) GetToken() string {
	return t.baererToken
}

// RefreshToken - Refreshes the token.
func (t *token) RefreshToken() error {
	if !t.IsSet() || t.IsExpired() {
		c := resty.New().SetBaseURL(t.endpoint)

		criteria := url.Values{
			"grant_type": {"password"},
			"username":   {t.username},
			"password":   {t.password},
		}

		r, err := c.R().
			SetDebug(t.debug).
			SetHeader("Content-Type", "application/x-www-form-url-encoded").
			SetHeader("Accept", "application/json").
			SetResult(&authTokenResponse{}).
			SetError(&apiAuthTokenErrorResponse{}).
			SetFormDataFromValues(criteria).
			Post("/auth/token")
		if err != nil {
			return err
		}

		if r.IsError() {
			return fmt.Errorf("authentication failed: %w", ToError(r))
		}

		refreshedToken := r.Result().(*authTokenResponse)
		t.baererToken = refreshedToken.AccessToken
		t.expiresAt = time.Now().Add(time.Duration(refreshedToken.ExpiresIn) * time.Second)
	}
	return nil
}

type authTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	UserName     string `json:"userName"`
	UserGUID     string `json:"userGuid"`
	Issued       string `json:".issued"`
	Expires      string `json:".expires"`
}

type apiAuthTokenErrorResponse struct {
	Error string `json:"error"`
}

// FormatError - Formats the Netbackup auth error, omitting the field when
// empty. Returns an empty string if Error is not set.
func (e *apiAuthTokenErrorResponse) FormatError() string {
	if e.Error == "" {
		return ""
	}
	return e.Error
}

// apiCallError carries the HTTP status code alongside the formatted message
// so callers can check the status structurally instead of substring-matching
// the final error string, which may embed untrusted raw response bodies.
// Mirrors commoncloudavenue.apiCallError; kept as a small unexported
// duplicate here rather than exported/shared to avoid introducing a
// cross-package dependency between clientnetbackup and commoncloudavenue.
type apiCallError struct {
	statusCode int
	message    string
}

func (e *apiCallError) Error() string {
	return e.message
}

// ToError - Converts a resty response into an error.
// It prefers the structured apiAuthTokenErrorResponse field when available,
// falling back to the raw HTTP status and response body when the typed
// error has nothing usable (e.g. plain-text or HTML error bodies from
// rate-limiting or upstream gateways).
func ToError(r *resty.Response) error {
	statusCode := r.StatusCode()

	authErr, _ := r.Error().(*apiAuthTokenErrorResponse)
	if authErr != nil {
		if formatted := authErr.FormatError(); formatted != "" {
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
