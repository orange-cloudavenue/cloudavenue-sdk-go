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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/consoles"
	cloudavenueerrors "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/errors"
)

func clearCloudavenueEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		"CLOUDAVENUE_CORE_API",
		"CLOUDAVENUE_URL",
		"CLOUDAVENUE_USERNAME",
		"CLOUDAVENUE_PASSWORD",
		"CLOUDAVENUE_ORG",
		"CLOUDAVENUE_VDC",
	} {
		t.Setenv(key, "")
	}
	t.Setenv("CLOUDAVENUE_DEBUG", "false")
	t.Setenv("CLOUDAVENUE_DEV", "false")
}

func TestOptsValidateDefaultsCoreAPI(t *testing.T) {
	clearCloudavenueEnv(t)
	t.Setenv("CLOUDAVENUE_DEV", "true")

	opts := &Opts{
		URL:      "https://vcd.example.com",
		Username: "username",
		Password: "password",
		Org:      "org",
	}

	err := opts.Validate()

	assert.NoError(t, err)
	assert.Equal(t, "https://vcd.example.com", opts.URL)
	assert.Equal(t, consoles.CerberusAPIEndpoint, opts.CoreAPI)
}

func TestOptsValidateRejectsInvalidCoreAPI(t *testing.T) {
	tests := []struct {
		name    string
		coreAPI string
	}{
		{name: "malformed URL", coreAPI: "://invalid"},
		{name: "relative URL", coreAPI: "/backend"},
		{name: "http URL", coreAPI: "http://core-api.example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearCloudavenueEnv(t)
			t.Setenv("CLOUDAVENUE_DEV", "true")

			opts := &Opts{
				URL:      "https://vcd.example.com",
				Username: "username",
				Password: "password",
				Org:      "org",
				CoreAPI:  tt.coreAPI,
			}

			err := opts.Validate()

			assert.ErrorIs(t, err, cloudavenueerrors.ErrInvalidFormat)
			assert.EqualError(t, err, "the core API \""+tt.coreAPI+"\" has an invalid format")
		})
	}
}

func TestInitKeepsVMwareURLSeparatedFromCoreAPI(t *testing.T) {
	clearCloudavenueEnv(t)
	t.Setenv("CLOUDAVENUE_DEV", "true")
	resetClientState(t)

	opts := &Opts{
		URL:      "https://vcd.example.com",
		Username: "username",
		Password: "password",
		Org:      "org",
		CoreAPI:  "https://core-api.example.com",
	}

	err := Init(opts)

	assert.NoError(t, err)
	assert.Equal(t, "https://vcd.example.com", c.token.GetEndpoint())
	assert.Equal(t, "https://core-api.example.com", c.token.effectiveCoreAPI())
	assert.NotEqual(t, c.token.GetEndpoint(), c.token.effectiveCoreAPI())
}

func TestTokenUsesCoreAPIForAuthAndBackendClient(t *testing.T) {
	var mu sync.Mutex
	var paths []string
	var authHeaders []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		paths = append(paths, r.URL.Path)
		authHeaders = append(authHeaders, r.Header.Get("Authorization"))
		mu.Unlock()

		switch r.URL.Path {
		case "/auth/v1/user/token":
			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(map[string]any{
				"access_token": "access-token",
				"token_type":   "Bearer",
				"expires_in":   3600,
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				assert.NoError(t, err)
			}
		case "/infrapicustomerproxy/v2.0/configurations":
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	token := token{
		clientID:     "username",
		clientSecret: "password",
		org:          "org",
		coreAPI:      server.URL,
	}

	err := token.RefreshToken()
	assert.NoError(t, err)

	backendClient := token.newBackendClient()
	assert.Equal(t, server.URL, backendClient.BaseURL)

	resp, err := backendClient.R().Get("/infrapicustomerproxy/v2.0/configurations")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode())

	// All HTTP calls above complete synchronously; no concurrent writers remain at this point.
	mu.Lock()
	defer mu.Unlock()
	if assert.Len(t, paths, 2) {
		assert.Equal(t, []string{"/auth/v1/user/token", "/infrapicustomerproxy/v2.0/configurations"}, paths)
	}
	if assert.Len(t, authHeaders, 2) {
		assert.Equal(t, "", authHeaders[0]) // auth endpoint uses no Authorization header — credentials are in the request body
		assert.Equal(t, "Bearer access-token", authHeaders[1])
	}
}

func TestTokenFallsBackToDefaultCoreAPIForAuthAndBackendClient(t *testing.T) {
	token := token{
		clientID:     "username",
		clientSecret: "password",
		org:          "org",
	}

	authClient := token.newAuthClient()
	backendClient := token.newBackendClient()

	assert.Equal(t, consoles.CerberusAPIEndpoint, authClient.BaseURL)
	assert.Equal(t, consoles.CerberusAPIEndpoint, backendClient.BaseURL)
	assert.Equal(t, authClient.BaseURL, backendClient.BaseURL)
}

func resetClientState(t *testing.T) {
	t.Helper()

	previousClient := c
	previousCache := cache

	c = &internalClient{}
	cache = nil

	t.Cleanup(func() {
		c = previousClient
		cache = previousCache
	})
}
