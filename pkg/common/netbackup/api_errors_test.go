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

package commonnetbackup

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/go-resty/resty/v2"
)

func TestAPIError_FormatError(t *testing.T) {
	tests := []struct {
		name string
		err  APIError
		want string
	}{
		{
			name: "no entries",
			err:  APIError{},
			want: "",
		},
		{
			name: "single entry with property and message",
			err:  APIError{{PropertyName: "Name", Message: "is required"}},
			want: "Name: is required",
		},
		{
			name: "single entry with message only",
			err:  APIError{{Message: "unexpected error"}},
			want: "unexpected error",
		},
		{
			name: "single entry with property only",
			err:  APIError{{PropertyName: "Name"}},
			want: "Name",
		},
		{
			name: "entry with both fields empty is skipped",
			err:  APIError{{}},
			want: "",
		},
		{
			name: "multiple entries joined",
			err: APIError{
				{PropertyName: "Name", Message: "is required"},
				{PropertyName: "Email", Message: "is invalid"},
			},
			want: "Name: is required; Email: is invalid",
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

// doNetbackupRequest performs a resty GET request against the given test
// server and returns the raw *resty.Response with SetError wired to
// APIError.
func doNetbackupRequest(t *testing.T, server *httptest.Server) *resty.Response {
	t.Helper()

	r, err := resty.New().R().
		SetError(&APIError{}).
		Get(server.URL)
	if err != nil {
		t.Fatalf("unexpected request error: %v", err)
	}

	return r
}

func TestNetbackupToError(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		body         string
		contentType  string
		wantExact    string
		wantContains []string
	}{
		{
			name:        "valid JSON error body",
			statusCode:  http.StatusBadRequest,
			body:        `[{"PropertyName":"Name","Message":"is required"}]`,
			contentType: "application/json",
			wantExact:   "Name: is required",
		},
		{
			name:        "empty body falls back to HTTP status",
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
			name:        "empty JSON array falls back to HTTP status and raw body",
			statusCode:  http.StatusInternalServerError,
			body:        `[]`,
			contentType: "application/json",
			wantExact:   "HTTPCode:500 Internal Server Error - body: []",
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

			r := doNetbackupRequest(t, server)

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
		})
	}
}
