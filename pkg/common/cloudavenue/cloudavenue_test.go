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
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/go-resty/resty/v2"
)

func TestAPIErrorResponse_FormatError(t *testing.T) {
	tests := []struct {
		name string
		err  APIErrorResponse
		want string
	}{
		{
			name: "all fields empty",
			err:  APIErrorResponse{},
			want: "",
		},
		{
			name: "only message set",
			err:  APIErrorResponse{Message: "boom"},
			want: "ErrorMessage:boom",
		},
		{
			name: "only code set",
			err:  APIErrorResponse{Code: "404"},
			want: "ErrorCode:404",
		},
		{
			name: "code and message set",
			err:  APIErrorResponse{Code: "500", Message: "boom"},
			want: "ErrorCode:500 - ErrorMessage:boom",
		},
		{
			name: "all fields set",
			err:  APIErrorResponse{Code: "400", Reason: "badRequest", Message: "invalid input"},
			want: "ErrorCode:400 - ErrorReason:badRequest - ErrorMessage:invalid input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.FormatError(); got != tt.want {
				t.Errorf("FormatError() = %q, want %q", got, tt.want)
			}
		})
	}
}

// doRequest performs a resty GET request against the given test server and
// returns the raw *resty.Response with SetError wired to APIErrorResponse.
func doRequest(t *testing.T, server *httptest.Server) *resty.Response {
	t.Helper()

	r, err := resty.New().R().
		SetError(&APIErrorResponse{}).
		Get(server.URL)
	if err != nil {
		t.Fatalf("unexpected request error: %v", err)
	}

	return r
}

func TestToError(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		body           string
		contentType    string
		wantExact      string // when set, ToError().Error() must match exactly
		wantContains   []string
		wantNotContain []string
	}{
		{
			name:        "valid full JSON error body",
			statusCode:  http.StatusBadRequest,
			body:        `{"code":"400","reason":"badRequest","message":"invalid input"}`,
			contentType: "application/json",
			wantExact:   "ErrorCode:400 - ErrorReason:badRequest - ErrorMessage:invalid input",
		},
		{
			// Exact-match regression test for the duplicated-status-code bug:
			// r.Status() already returns "<code> <reason>", so pairing it
			// with a separate "%d" of r.StatusCode() produced garbage like
			// "HTTPCode:503 503 Service Unavailable".
			name:        "empty body falls back to HTTP status (exact match)",
			statusCode:  http.StatusServiceUnavailable,
			body:        "",
			contentType: "application/json",
			wantExact:   "HTTPCode:503 Service Unavailable",
		},
		{
			name:        "malformed non-JSON body falls back to HTTP status and body",
			statusCode:  http.StatusTooManyRequests,
			body:        "<html><body>Too Many Requests</body></html>",
			contentType: "text/html",
			wantExact:   "HTTPCode:429 Too Many Requests - body: <html><body>Too Many Requests</body></html>",
		},
		{
			name:        "empty JSON object falls back to HTTP status and raw body",
			statusCode:  http.StatusInternalServerError,
			body:        `{}`,
			contentType: "application/json",
			wantExact:   "HTTPCode:500 Internal Server Error - body: {}",
		},
		{
			name:        "oversized body is truncated",
			statusCode:  http.StatusServiceUnavailable,
			body:        strings.Repeat("x", 1000),
			contentType: "text/plain",
			wantContains: []string{
				"HTTPCode:503",
				strings.Repeat("x", 512) + "... (truncated)",
			},
			wantNotContain: []string{strings.Repeat("x", 513)},
		},
		{
			// Regression test for byte-slice truncation splitting a
			// multi-byte UTF-8 sequence: fill the body with a 3-byte rune
			// ("é" is 2 bytes; use "€" which is 3 bytes: U+20AC) so that a
			// naive body[:512] slice lands mid-rune, and assert the result
			// is still valid UTF-8 with no replacement-character garbage.
			name:        "oversized non-ASCII body truncates on a valid UTF-8 boundary",
			statusCode:  http.StatusServiceUnavailable,
			body:        strings.Repeat("€", 200), // 600 bytes, well past MaxErrorBodyLen (512)
			contentType: "text/plain; charset=utf-8",
			wantContains: []string{
				"HTTPCode:503",
				"... (truncated)",
			},
			wantNotContain: []string{"\uFFFD"}, // unicode replacement character
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", tt.contentType)
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer server.Close()

			r := doRequest(t, server)

			got := ToError(r).Error()

			if !utf8.ValidString(got) {
				t.Errorf("ToError() = %q, is not valid UTF-8", got)
			}

			if tt.wantExact != "" {
				if got != tt.wantExact {
					t.Errorf("ToError() = %q, want exact %q", got, tt.wantExact)
				}
				return
			}

			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("ToError() = %q, want substring %q", got, want)
				}
			}
			for _, notWant := range tt.wantNotContain {
				if strings.Contains(got, notWant) {
					t.Errorf("ToError() = %q, should not contain %q", got, notWant)
				}
			}
		})
	}
}

func TestIsNotFound_NilError(t *testing.T) {
	if IsNotFound(nil) {
		t.Error("IsNotFound(nil) = true, want false")
	}
}

func TestIsNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"code":"404","reason":"notFound","message":"resource not found"}`))
	}))
	defer server.Close()

	r := doRequest(t, server)
	err := ToError(r)

	if !IsNotFound(err) {
		t.Errorf("IsNotFound() = false, want true for error: %v", err)
	}
}

func TestIsNotFound_FallbackPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("not found"))
	}))
	defer server.Close()

	r := doRequest(t, server)
	err := ToError(r)

	if !IsNotFound(err) {
		t.Errorf("IsNotFound() = false, want true for fallback error: %v", err)
	}
}

// TestIsNotFound_NotSpoofableByBodyContent is a regression test proving
// IsNotFound checks the real HTTP status structurally instead of
// substring-matching the final (possibly untrusted, body-derived) error
// string. A raw response body that happens to literally contain
// "ErrorCode:404" must NOT cause a false positive when the actual HTTP
// status is not 404.
func TestIsNotFound_NotSpoofableByBodyContent(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{name: "spoofed body with 200 OK status", statusCode: http.StatusOK},
		{name: "spoofed body with 500 status", statusCode: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte("unrelated text containing ErrorCode:404 by coincidence"))
			}))
			defer server.Close()

			r := doRequest(t, server)

			// Only exercise the fallback/spoofing scenario when resty
			// actually treats the response as an error (IsError() gates
			// every real call site); for the 200 case we still want to
			// prove ToError()+IsNotFound() aren't fooled by body content
			// if ever called directly.
			err := ToError(r)

			if IsNotFound(err) {
				t.Errorf("IsNotFound() = true, want false for status %d with spoofed body: %v", tt.statusCode, err)
			}
		})
	}
}
