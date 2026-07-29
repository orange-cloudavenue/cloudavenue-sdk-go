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
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/consoles"
)

// maxErrorBodyLen caps the amount of raw response body embedded in a
// fallback error message, mirroring commoncloudavenue.ToError, to avoid
// unbounded error messages / log bloat when Cerberus returns large error
// pages (e.g. HTML gateway pages) instead of its usual JSON error body.
const maxErrorBodyLen = 512

// bearerTokenType is the default OAuth2 token type used when the server
// does not return one explicitly.
const bearerTokenType = "Bearer"

// token holds the OAuth2 authentication state for Cerberus API.
type token struct {
	// mu guards all reads/writes of the OAuth2 token fields below
	// (accessToken, tokenType, expiresAt). Without it, concurrent unguarded
	// reads (e.g. via IsExpired/IsSet/GetToken from New()/GetBearerToken())
	// racing against RefreshToken()'s writes are a genuine data race under
	// Go's memory model (torn reads on the multi-word time.Time and the
	// string header), not just a staleness window. It also prevents the
	// thundering-herd problem: multiple goroutines racing an expiring token
	// would otherwise each retry the (now throttled) auth endpoint
	// independently, amplifying load instead of reducing it.
	mu sync.RWMutex

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
	coreAPI  string
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

func (t *token) effectiveCoreAPI() string {
	if t.coreAPI == "" {
		return consoles.CerberusAPIEndpoint
	}

	return t.coreAPI
}

func (t *token) newBackendClient() *resty.Client {
	return configureRetry(resty.New().
		SetDebug(t.debug).
		SetHeader("Accept", "application/json").
		SetBaseURL(t.effectiveCoreAPI()).
		SetAuthScheme(bearerTokenType).
		OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
			return t.RefreshToken()
		}).
		SetAuthToken(t.GetToken()).
		SetHeader("User-Agent", "Cloudavenue-SDK-v1"))
}

func (t *token) newAuthClient() *resty.Client {
	return configureRetry(resty.New().SetBaseURL(t.effectiveCoreAPI()))
}

// GetEndpointURL - Returns the API endpoint URL.
func (t *token) GetEndpointURL() url.URL {
	u, _ := url.Parse(t.endpoint)
	return *u
}

// IsExpired - Returns true if the token is expired.
// Includes a 30-second buffer to prevent edge cases.
func (t *token) IsExpired() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.isExpiredLocked()
}

// isExpiredLocked is the lock-free variant of IsExpired, for internal use by
// callers (namely RefreshToken) that already hold t.mu. sync.RWMutex is not
// reentrant, so RefreshToken must not call the public IsExpired while
// holding the write lock.
func (t *token) isExpiredLocked() bool {
	return t.expiresAt.Add(-30 * time.Second).Before(time.Now())
}

// IsSet - Returns true if the token is set.
func (t *token) IsSet() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.isSetLocked()
}

// isSetLocked is the lock-free variant of IsSet, for internal use by callers
// (namely RefreshToken) that already hold t.mu.
func (t *token) isSetLocked() bool {
	return t.accessToken != ""
}

// GetToken - Returns the access token.
func (t *token) GetToken() string {
	t.mu.RLock()
	defer t.mu.RUnlock()

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
	t.mu.Lock()
	defer t.mu.Unlock()

	// Use the lock-free variants here: the public IsSet()/IsExpired() would
	// try to re-acquire t.mu (RLock), which deadlocks since sync.RWMutex is
	// not reentrant and we're already holding the write lock above.
	if t.isSetLocked() && !t.isExpiredLocked() {
		return nil
	}

	c := t.newAuthClient()

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
		return fmt.Errorf("authentication failed: %w", ToError(r))
	}

	// Parse the OAuth2 response
	authResp, ok := r.Result().(*cerberusAuthResponse)
	if !ok || authResp == nil {
		return errors.New("authentication failed: invalid response format")
	}

	if authResp.AccessToken == "" {
		return errors.New("authentication failed: empty access token received")
	}

	// Set the token
	t.accessToken = authResp.AccessToken
	t.tokenType = authResp.Type
	if t.tokenType == "" {
		t.tokenType = bearerTokenType // Default to Bearer if not specified
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

// FormatError - Formats the Cerberus error, omitting any fields that are
// empty. Code is Cerberus's own internal error code (distinct from the HTTP
// status, which is already carried separately by ToError/apiCallError), so
// it is only included when non-zero. Returns an empty string if none of the
// fields are set.
func (e *CerberusErrorResponse) FormatError() string {
	var parts []string
	if e.Code != 0 {
		parts = append(parts, fmt.Sprintf("code:%d", e.Code))
	}
	if e.Message != "" {
		parts = append(parts, fmt.Sprintf("message:%s", e.Message))
	}
	if e.Description != "" {
		parts = append(parts, fmt.Sprintf("description:%s", e.Description))
	}

	return strings.Join(parts, ": ")
}

// apiCallError carries the HTTP status code alongside the formatted message
// so callers can check the status structurally instead of substring-matching
// the final error string, which may embed untrusted raw response bodies.
// Mirrors commoncloudavenue.apiCallError; kept as a small unexported
// duplicate here rather than exported/shared to avoid introducing a
// cross-package dependency between clientcloudavenue and commoncloudavenue.
type apiCallError struct {
	statusCode int
	message    string
}

func (e *apiCallError) Error() string {
	return e.message
}

// ToError - Converts a resty response into an error.
// It prefers the structured CerberusErrorResponse fields when available,
// falling back to the raw HTTP status and response body when the typed
// error has nothing usable (e.g. plain-text or HTML error bodies from
// rate-limiting or upstream gateways).
func ToError(r *resty.Response) error {
	statusCode := r.StatusCode()

	cerberusErr, _ := r.Error().(*CerberusErrorResponse)
	if cerberusErr != nil {
		if formatted := cerberusErr.FormatError(); formatted != "" {
			return &apiCallError{statusCode: statusCode, message: formatted}
		}
	}

	body := strings.TrimSpace(r.String())
	if body == "" {
		return &apiCallError{statusCode: statusCode, message: fmt.Sprintf("HTTPCode:%s", r.Status())}
	}

	if len(body) > maxErrorBodyLen {
		// strings.ToValidUTF8 strips any partial/invalid trailing rune left
		// dangling by the byte-slice truncation, instead of producing
		// garbled replacement characters for non-ASCII upstream bodies.
		body = strings.ToValidUTF8(body[:maxErrorBodyLen], "") + "... (truncated)"
	}

	return &apiCallError{statusCode: statusCode, message: fmt.Sprintf("HTTPCode:%s - body: %s", r.Status(), body)}
}
