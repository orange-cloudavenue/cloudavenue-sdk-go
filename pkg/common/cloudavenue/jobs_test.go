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
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-resty/resty/v2"

	caverrors "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/errors"
)

func doRawRequest(t *testing.T, server *httptest.Server) *resty.Response {
	t.Helper()

	r, err := resty.New().R().Get(server.URL)
	if err != nil {
		t.Fatalf("unexpected request error: %v", err)
	}

	return r
}

func TestMissingJobIDError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"message":"capacity exceeded"}`)
	}))
	defer srv.Close()

	err := MissingJobIDError(doRawRequest(t, srv))

	if !errors.Is(err, ErrJobIDEmpty) {
		t.Fatalf("errors.Is(err, ErrJobIDEmpty) = false, want true (err: %v)", err)
	}

	msg := err.Error()
	for _, want := range []string{"job id is empty", "HTTP 201 Created", "capacity exceeded"} {
		if !strings.Contains(msg, want) {
			t.Errorf("error message %q does not contain %q", msg, want)
		}
	}

	var apiErr *apiCallError
	if !errors.As(err, &apiErr) {
		t.Fatalf("errors.As(err, **apiCallError) = false, want true")
	}
	if apiErr.statusCode != http.StatusCreated {
		t.Errorf("statusCode = %d, want %d", apiErr.statusCode, http.StatusCreated)
	}
}

func TestMissingJobIDErrorTruncatesBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, strings.Repeat("x", caverrors.MaxErrorBodyLen+100))
	}))
	defer srv.Close()

	msg := MissingJobIDError(doRawRequest(t, srv)).Error()
	if !strings.HasSuffix(msg, "... (truncated)") {
		t.Errorf("error message should end with the truncation marker, got: %q", msg)
	}
}

func TestRefreshEmptyJobIDWrapsSentinel(t *testing.T) {
	err := (&JobStatus{}).Refresh()

	if !errors.Is(err, ErrJobIDEmpty) {
		t.Fatalf("errors.Is(err, ErrJobIDEmpty) = false, want true (err: %v)", err)
	}
	// Keep the exact historical message: callers substring-match on it.
	const want = "cannot refresh job status: job id is empty"
	if err.Error() != want {
		t.Errorf("message = %q, want %q", err.Error(), want)
	}
}
