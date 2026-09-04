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
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
)

// configureRetry applies a conservative retry policy to the given resty
// client: retry on network errors and on HTTP 429/503 responses, honoring
// the Retry-After header when present, otherwise falling back to resty's
// default exponential-backoff-with-jitter algorithm.
func configureRetry(c *resty.Client) *resty.Client {
	c.
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(30 * time.Second).
		AddRetryCondition(isRetryableResponse).
		SetRetryAfter(retryAfterFromHeader)

	// Log retry exhaustion so callers can detect when a failure
	// happened after multiple retry attempts rather than on the
	// first try.
	c.OnError(func(r *resty.Request, err error) {
		if r.Attempt > 0 {
			log.Printf("[WARN] cloudavenue: retrying %s %s failed (attempt %d/%d): %v",
				r.Method, r.URL, r.Attempt+1, c.RetryCount+1, err)
		}
	})

	return c
}

// idempotentMethods are the HTTP methods safe to retry after a network-level
// error, i.e. methods where retrying cannot cause an unintended duplicate
// side effect on the server.
var idempotentMethods = map[string]bool{
	http.MethodGet:    true,
	http.MethodHead:   true,
	http.MethodPut:    true,
	http.MethodDelete: true,
}

// isRetryableResponse - Returns true if the request should be retried.
//
// HTTP 429 (Too Many Requests) / 503 (Service Unavailable) are always
// retried: these indicate the server explicitly declined to process the
// request (throttling or temporary unavailability), so retrying is safe.
//
// HTTP 502 (Bad Gateway) and 504 (Gateway Timeout) are NOT retried:
// these can arrive after the upstream server has already started processing
// the request, particularly with long-running operations. Retrying a 502
// on a write operation (POST/PUT/DELETE) risks duplicating the side effect.
// For read-heavy workloads that encounter 502/504, the caller should
// implement its own retry with appropriate idempotency guarantees.
//
// Network-level errors (err != nil) are only retried for idempotent HTTP
// methods (GET, HEAD, PUT, DELETE). A network error can occur after the
// server already processed the request, so retrying a POST/PATCH risks
// creating a duplicate resource. If the request method cannot be
// determined, this fails closed (does not retry).
func isRetryableResponse(r *resty.Response, err error) bool {
	if err != nil {
		if r == nil || r.Request == nil {
			// Method can't be determined: fail closed, do not retry.
			return false
		}
		return idempotentMethods[r.Request.Method]
	}

	if r == nil {
		return false
	}

	switch r.StatusCode() {
	case http.StatusTooManyRequests, http.StatusServiceUnavailable:
		return true
	default:
		return false
	}
}

// retryAfterFromHeader - Parses the Retry-After header (RFC 7231 §7.1.3) on
// 429/503 responses, supporting both the delta-seconds and HTTP-date
// formats. Returns (0, nil) when the header is absent, unparseable, or the
// status code isn't 429/503, which tells resty to fall back to its default
// exponential-backoff-with-jitter algorithm.
func retryAfterFromHeader(_ *resty.Client, r *resty.Response) (time.Duration, error) {
	if r == nil {
		return 0, nil
	}

	switch r.StatusCode() {
	case http.StatusTooManyRequests, http.StatusServiceUnavailable:
	default:
		return 0, nil
	}

	header := r.Header().Get("Retry-After")
	if header == "" {
		return 0, nil
	}

	// delta-seconds
	if seconds, err := strconv.Atoi(header); err == nil {
		if seconds < 0 {
			return 0, nil
		}
		return nonZeroDuration(time.Duration(seconds) * time.Second), nil
	}

	// HTTP-date
	if t, err := http.ParseTime(header); err == nil {
		d := time.Until(t)
		if d < 0 {
			return 0, nil
		}
		return nonZeroDuration(d), nil
	}

	// Unparseable value: fall back to the default backoff rather than
	// aborting the retry logic.
	return 0, nil
}

// nonZeroDuration guards against resty's RetryAfterFunc contract, where
// returning exactly 0 means "no override, use default jitter backoff"
// instead of "retry immediately" as RFC 7231 intends for `Retry-After: 0`.
// A successfully parsed zero/near-zero duration is remapped to a minimal
// non-zero delay so the explicit "retry now" signal isn't silently
// converted into a 1-30s default backoff.
func nonZeroDuration(d time.Duration) time.Duration {
	if d <= 0 {
		return 1 * time.Millisecond
	}
	return d
}
