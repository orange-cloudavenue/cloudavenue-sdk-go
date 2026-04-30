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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/consoles"
	cloudavenueerrors "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/errors"
)

func TestOptsValidateDefaultsCoreAPI(t *testing.T) {
	opts := &Opts{
		URL:      "https://vcd.example.com",
		Username: "username",
		Password: "password",
		Org:      "org",
		Dev:      true,
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
			opts := &Opts{
				URL:      "https://vcd.example.com",
				Username: "username",
				Password: "password",
				Org:      "org",
				Dev:      true,
				CoreAPI:  tt.coreAPI,
			}

			err := opts.Validate()

			assert.ErrorIs(t, err, cloudavenueerrors.ErrInvalidFormat)
			assert.EqualError(t, err, "the core API \""+tt.coreAPI+"\" has an invalid format")
		})
	}
}

func TestInitKeepsVMwareURLSeparatedFromCoreAPI(t *testing.T) {
	resetClientState(t)

	opts := &Opts{
		URL:      "https://vcd.example.com",
		Username: "username",
		Password: "password",
		Org:      "org",
		Dev:      true,
		CoreAPI:  "https://core-api.example.com",
	}

	err := Init(opts)

	assert.NoError(t, err)
	assert.Equal(t, "https://vcd.example.com", c.token.GetEndpoint())
	assert.Equal(t, "https://core-api.example.com", c.token.effectiveCoreAPI())
	assert.NotEqual(t, c.token.GetEndpoint(), c.token.effectiveCoreAPI())
}

func TestTokenUsesCoreAPIForAuthAndBackendClient(t *testing.T) {
	var paths []string
	var authHeaders []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		authHeaders = append(authHeaders, r.Header.Get("Authorization"))

		switch r.URL.Path {
		case "/auth/v1/user/token":
			w.Header().Set("Content-Type", "application/json")
			assert.NoError(t, json.NewEncoder(w).Encode(map[string]any{
				"access_token": "access-token",
				"token_type":   "Bearer",
				"expires_in":   3600,
			}))
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

	assert.Equal(t, []string{"/auth/v1/user/token", "/infrapicustomerproxy/v2.0/configurations"}, paths)
	assert.Equal(t, "", authHeaders[0])
	assert.Equal(t, "Bearer access-token", authHeaders[1])
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
