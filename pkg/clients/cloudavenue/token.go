/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package clientcloudavenue

import (
	"fmt"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
)

// token holds the OAuth2 authentication state for Cerberus API.
type token struct {
	// OAuth2 token fields
	accessToken string
	tokenType   string // "Bearer"
	expiresAt   time.Time

	// Credentials
	clientID     string // username
	clientSecret string // password
	org          string

	// Legacy fields maintained for compatibility
	orgID string
	vdc   string

	// Connection settings
	endpoint string
	debug    bool
}

// GetOrganization - Returns the organization.
func (t *token) GetOrganization() string {
	return t.org
}

// GetEndpoint - Returns the API endpoint.
func (t *token) GetEndpoint() string {
	return t.endpoint
}

// GetEndpointURL - Returns the API endpoint URL.
func (t *token) GetEndpointURL() url.URL {
	u, _ := url.Parse(t.endpoint)
	return *u
}

// IsExpired - Returns true if the token is expired.
// Includes a 30-second buffer to prevent edge cases.
func (t *token) IsExpired() bool {
	return t.expiresAt.Add(-30 * time.Second).Before(time.Now())
}

// IsSet - Returns true if the token is set.
func (t *token) IsSet() bool {
	return t.accessToken != ""
}

// GetToken - Returns the access token.
func (t *token) GetToken() string {
	return t.accessToken
}

// GetTokenType - Returns the token type (Bearer).
func (t *token) GetTokenType() string {
	return t.tokenType
}

// GetOrgID - Returns the organization ID.
// Note: OrgID is no longer available directly from Cerberus auth response.
// It must be fetched separately via /infrapicustomerproxy/v2.0/configurations.
func (t *token) GetOrgID() string {
	return t.orgID
}

// SetOrgID - Sets the organization ID.
// Used to set OrgID after fetching it from configurations endpoint.
func (t *token) SetOrgID(orgID string) {
	t.orgID = orgID
}

// RefreshToken - Authenticates to Cerberus API using OAuth2 Client Credentials.
// POST /auth/v1/user/token
// Content-Type: application/x-www-form-urlencoded
// Body: grant_type=client_credentials&client_id={username}&client_secret={password}&scope=tenant:{org}
func (t *token) RefreshToken() error {
	if t.IsSet() && !t.IsExpired() {
		return nil
	}

	c := resty.New().SetBaseURL("https://api1.cloudavenue.orange-business.com")

	r, err := c.R().
		SetDebug(t.debug).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetHeader("Accept", "application/json").
		SetFormData(map[string]string{
			"grant_type":    "client_credentials",
			"client_id":     t.clientID,
			"client_secret": t.clientSecret,
			"scope":         "tenant:" + t.org,
		}).
		SetResult(&cerberusAuthResponse{}).
		SetError(&CerberusErrorResponse{}).
		Post("/auth/v1/user/token")
	if err != nil {
		return fmt.Errorf("authentication request failed: %w", err)
	}

	if r.IsError() {
		cerberusErr, ok := r.Error().(*CerberusErrorResponse)
		if ok && cerberusErr != nil {
			return fmt.Errorf("authentication failed: HTTPCode:%s - %s", r.Status(), cerberusErr.FormatError())
		}
		return fmt.Errorf("authentication failed: HTTPCode:%s", r.Status())
	}

	// Parse the OAuth2 response
	authResp, ok := r.Result().(*cerberusAuthResponse)
	if !ok || authResp == nil {
		return fmt.Errorf("authentication failed: invalid response format")
	}

	if authResp.AccessToken == "" {
		return fmt.Errorf("authentication failed: empty access token received")
	}

	// Set the token
	t.accessToken = authResp.AccessToken
	t.tokenType = authResp.Type
	if t.tokenType == "" {
		t.tokenType = "Bearer" // Default to Bearer if not specified
	}

	// Calculate the expiration date (expires_in is in seconds)
	t.expiresAt = time.Now().Add(time.Duration(authResp.ExpiresIn) * time.Second)

	return nil
}

// cerberusAuthResponse - OAuth2 token response from Cerberus API.
// Response from POST /auth/v1/user/token
type cerberusAuthResponse struct {
	AccessToken string `json:"access_token"`
	Type        string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"` // seconds
}

// CerberusErrorResponse - Error response from Cerberus API.
type CerberusErrorResponse struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description,omitempty"`
}

// FormatError - Formats the Cerberus error.
func (e *CerberusErrorResponse) FormatError() string {
	if e.Description != "" {
		return fmt.Sprintf("ErrorCode:%d - Message:%s - Description:%s", e.Code, e.Message, e.Description)
	}
	return fmt.Sprintf("ErrorCode:%d - Message:%s", e.Code, e.Message)
}

// APIErrorResponse - Generic API error response.
// Kept for backward compatibility with existing code.
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
