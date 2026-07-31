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
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
)

// doGet performs a resty GET against the given server and returns the raw
// *resty.Response, with retries disabled so the response reflects a single
// attempt.
func doGet(t *testing.T, server *httptest.Server) *resty.Response {
	t.Helper()

	r, err := resty.New().R().Get(server.URL)
	if err != nil {
		t.Fatalf("unexpected request error: %v", err)
	}

	return r
}

// doNetworkError performs a request with the given method against an
// unreachable address, forcing a network-level error. Resty still populates
// the returned *resty.Response.Request (including Method) even though there
// is no HTTP response body — verified against go-resty/resty/v2 v2.17.2's
// Client.execute, which builds `&Response{Request: req, RawResponse: resp}`
// before checking the transport error.
func doNetworkError(t *testing.T, method string) (*resty.Response, error) {
	t.Helper()

	const unreachableAddr = "http://127.0.0.1:1/unreachable"

	req := resty.New().SetRetryCount(0).R()

	var (
		r   *resty.Response
		err error
	)
	switch method {
	case http.MethodGet:
		r, err = req.Get(unreachableAddr)
	case http.MethodHead:
		r, err = req.Head(unreachableAddr)
	case http.MethodPost:
		r, err = req.Post(unreachableAddr)
	case http.MethodPut:
		r, err = req.Put(unreachableAddr)
	case http.MethodDelete:
		r, err = req.Delete(unreachableAddr)
	case http.MethodPatch:
		r, err = req.Patch(unreachableAddr)
	default:
		t.Fatalf("unsupported method in test helper: %s", method)
	}

	if err == nil {
		t.Fatalf("expected a network error, got nil (is 127.0.0.1:1 reachable in this environment?)")
	}

	return r, err
}

func TestIsRetryableResponse(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{
			name:       "429 Too Many Requests is retryable",
			statusCode: http.StatusTooManyRequests,
			want:       true,
		},
		{
			name:       "503 Service Unavailable is retryable",
			statusCode: http.StatusServiceUnavailable,
			want:       true,
		},
		{
			name:       "400 Bad Request is not retryable",
			statusCode: http.StatusBadRequest,
			want:       false,
		},
		{
			name:       "401 Unauthorized is not retryable",
			statusCode: http.StatusUnauthorized,
			want:       false,
		},
		{
			name:       "404 Not Found is not retryable",
			statusCode: http.StatusNotFound,
			want:       false,
		},
		{
			name:       "500 Internal Server Error is not retryable",
			statusCode: http.StatusInternalServerError,
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			r := doGet(t, server)

			if got := isRetryableResponse(r, nil); got != tt.want {
				t.Errorf("isRetryableResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsRetryableResponse_NetworkErrorRespectsIdempotency covers the
// data-integrity fix: a network error must only be retried for idempotent
// HTTP methods (GET/HEAD/PUT/DELETE), never for POST/PATCH, since a network
// error can occur after the server already processed a mutating request and
// retrying risks creating a duplicate resource.
func TestIsRetryableResponse_NetworkErrorRespectsIdempotency(t *testing.T) {
	tests := []struct {
		method string
		want   bool
	}{
		{method: http.MethodGet, want: true},
		{method: http.MethodHead, want: true},
		{method: http.MethodPut, want: true},
		{method: http.MethodDelete, want: true},
		{method: http.MethodPost, want: false},
		{method: http.MethodPatch, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			r, err := doNetworkError(t, tt.method)

			if got := isRetryableResponse(r, err); got != tt.want {
				t.Errorf("isRetryableResponse() for network error + %s = %v, want %v", tt.method, got, tt.want)
			}
		})
	}
}

// TestIsRetryableResponse_StatusRetryableRegardlessOfMethod covers the
// other half of the data-integrity fix: 429/503 status responses are always
// retryable, including for POST, because a status response means the
// request was rejected/throttled before processing.
func TestIsRetryableResponse_StatusRetryableRegardlessOfMethod(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	r, err := resty.New().R().Post(server.URL)
	if err != nil {
		t.Fatalf("unexpected request error: %v", err)
	}

	if got := isRetryableResponse(r, nil); got != true {
		t.Errorf("isRetryableResponse() for 429 + POST = %v, want true", got)
	}
}

// TestIsRetryableResponse_FailsClosedWhenMethodUnknown covers the
// fail-closed guarantee when the request method truly cannot be determined
// (e.g. both response and request are nil).
func TestIsRetryableResponse_FailsClosedWhenMethodUnknown(t *testing.T) {
	if got := isRetryableResponse(nil, errors.New("connection reset")); got != false {
		t.Errorf("isRetryableResponse(nil, err) = %v, want false (fail closed)", got)
	}
}

func TestRetryAfterFromHeader(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		retryAfter  string // empty means header not set
		wantSeconds float64
	}{
		{
			name:        "delta-seconds parsed correctly",
			statusCode:  http.StatusTooManyRequests,
			retryAfter:  "5",
			wantSeconds: 5,
		},
		{
			// RFC 7231 treats "Retry-After: 0" as "retry immediately". A
			// resty RetryAfterFunc returning exactly 0 means "no override,
			// use default jitter backoff" instead, so a successfully
			// parsed zero must be remapped to a minimal non-zero delay
			// rather than silently falling back to the 1-30s default.
			name:        "Retry-After: 0 means retry immediately, not default backoff",
			statusCode:  http.StatusTooManyRequests,
			retryAfter:  "0",
			wantSeconds: 0, // asserted precisely below, not via tolerance
		},
		{
			name:        "HTTP-date parsed correctly",
			statusCode:  http.StatusServiceUnavailable,
			retryAfter:  time.Now().Add(10 * time.Second).UTC().Format(http.TimeFormat),
			wantSeconds: 10,
		},
		{
			name:        "header absent falls back to default backoff",
			statusCode:  http.StatusTooManyRequests,
			retryAfter:  "",
			wantSeconds: 0,
		},
		{
			name:        "status not 429/503 falls back to default backoff",
			statusCode:  http.StatusInternalServerError,
			retryAfter:  "5",
			wantSeconds: 0,
		},
		{
			name:        "unparseable value falls back to default backoff without error",
			statusCode:  http.StatusTooManyRequests,
			retryAfter:  "not-a-valid-value",
			wantSeconds: 0,
		},
		{
			name:        "negative delta-seconds falls back to default backoff",
			statusCode:  http.StatusTooManyRequests,
			retryAfter:  "-5",
			wantSeconds: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				if tt.retryAfter != "" {
					w.Header().Set("Retry-After", tt.retryAfter)
				}
				w.WriteHeader(tt.statusCode)
			}))
			defer server.Close()

			r := doGet(t, server)

			d, err := retryAfterFromHeader(nil, r)
			if err != nil {
				t.Fatalf("retryAfterFromHeader() returned unexpected error: %v", err)
			}

			if tt.retryAfter == "0" {
				// Must be a strictly positive, minimal delay ("retry
				// immediately"), never exactly 0 (which resty interprets
				// as "no override, use default backoff").
				if d <= 0 {
					t.Errorf("retryAfterFromHeader() for Retry-After:0 = %v, want > 0 (retry immediately)", d)
				}
				if d > 100*time.Millisecond {
					t.Errorf("retryAfterFromHeader() for Retry-After:0 = %v, want a minimal near-zero delay", d)
				}
				return
			}

			wantDuration := time.Duration(tt.wantSeconds * float64(time.Second))

			// Allow a small tolerance for the HTTP-date case since it's
			// computed via time.Until relative to "now".
			diff := d - wantDuration
			if diff < 0 {
				diff = -diff
			}
			if diff > 2*time.Second {
				t.Errorf("retryAfterFromHeader() = %v, want ~%v", d, wantDuration)
			}
		})
	}
}
