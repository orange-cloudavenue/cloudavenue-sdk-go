// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0

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
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestRefreshToken_ConcurrentCallsSingleFlight covers the thundering-herd
// fix: concurrent RefreshToken calls around token expiry must result in
// exactly one actual HTTP round-trip to the auth endpoint, with all other
// goroutines observing the fresh token via the double-check under the
// mutex, instead of each independently retrying against an already
// throttled endpoint.
func TestRefreshToken_ConcurrentCallsSingleFlight(t *testing.T) {
	var authCalls int64

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/auth/v1/user/token" {
			atomic.AddInt64(&authCalls, 1)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"access_token": "access-token",
				"token_type":   "Bearer",
				"expires_in":   3600,
			})
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	tok := &token{
		clientID:     testUsername,
		clientSecret: testPassword,
		org:          testOrg,
		coreAPI:      server.URL,
	}

	const goroutines = 20

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for range goroutines {
		go func() {
			defer wg.Done()
			assert.NoError(t, tok.RefreshToken())
		}()
	}

	wg.Wait()

	assert.Equal(t, int64(1), atomic.LoadInt64(&authCalls), "expected exactly one auth round-trip across concurrent RefreshToken calls")
	assert.True(t, tok.IsSet())
}

// TestConcurrentTokenAccessorsRaceAgainstRefresh is a regression test for
// the incomplete thundering-herd fix: it exercises the read-side accessors
// (IsExpired, IsSet, GetToken) — exactly what New()/GetBearerToken() call —
// concurrently against an in-flight RefreshToken(), which is the real
// production usage pattern (~28 SDK entry points call New(), which reads
// these fields before/independently of RefreshToken()'s own locking).
//
// The previous single-flight test only ever raced RefreshToken() against
// itself, never against these read accessors, so it could not catch a
// mutex that was only acquired inside RefreshToken()'s body.
//
// The auth endpoint always returns an already-expired token (expires_in: 0
// combined with the 30s freshness buffer in IsExpired), forcing many real
// refresh cycles for the duration of the test instead of a single one, to
// maximize the race window. Run with `go test -race` to prove there is no
// data race on the token's fields.
func TestConcurrentTokenAccessorsRaceAgainstRefresh(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/auth/v1/user/token" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token": "access-token",
			"token_type":   "Bearer",
			"expires_in":   0, // always "expired" (minus the 30s buffer), forcing repeated refreshes
		})
	}))
	defer server.Close()

	tok := &token{
		clientID:     testUsername,
		clientSecret: testPassword,
		org:          testOrg,
		coreAPI:      server.URL,
	}

	const (
		refreshers = 5
		readers    = 20
	)

	stop := make(chan struct{})
	var wg sync.WaitGroup

	// Goroutines continuously triggering refreshes (the write side).
	wg.Add(refreshers)
	for range refreshers {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					_ = tok.RefreshToken()
				}
			}
		}()
	}

	// Goroutines continuously reading the token state (the read side used
	// by New() and GetBearerToken() in production).
	wg.Add(readers)
	for range readers {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					_ = tok.IsExpired()
					_ = tok.IsSet()
					_ = tok.GetToken()
				}
			}
		}()
	}

	time.Sleep(100 * time.Millisecond)
	close(stop)
	wg.Wait()

	assert.True(t, tok.IsSet())
}

// TestGetBearerToken_ConcurrentWithRefresh exercises the actual package-level
// production call path (GetBearerToken -> c.token.GetToken()) concurrently
// against an in-flight c.token.RefreshToken(), proving there is no data race
// between the package-level singleton token and callers that read it via
// GetBearerToken() without going through RefreshToken() first (the pattern
// New() uses before deciding whether to refresh).
func TestGetBearerToken_ConcurrentWithRefresh(t *testing.T) {
	resetClientState(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/auth/v1/user/token" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token": "access-token",
			"token_type":   "Bearer",
			"expires_in":   0, // always "expired", forcing repeated refreshes
		})
	}))
	defer server.Close()

	c.token = token{
		clientID:     testUsername,
		clientSecret: testPassword,
		org:          testOrg,
		coreAPI:      server.URL,
	}

	const (
		refreshers = 5
		readers    = 20
	)

	stop := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(refreshers)
	for range refreshers {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					_ = c.token.RefreshToken()
				}
			}
		}()
	}

	wg.Add(readers)
	for range readers {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					_ = GetBearerToken()
				}
			}
		}()
	}

	time.Sleep(100 * time.Millisecond)
	close(stop)
	wg.Wait()

	assert.NotEmpty(t, GetBearerToken())
}
