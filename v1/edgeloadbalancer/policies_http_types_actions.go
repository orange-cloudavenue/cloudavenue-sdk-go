/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgeloadbalancer

import (
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/internal/utils"
)

// * Action

type (
	PoliciesHTTPActionHeadersRewrite []*PoliciesHTTPActionHeaderRewrite
	PoliciesHTTPActionHeaderRewrite  struct {
		// Action for the chosen header
		// var PoliciesHTTPActionHeaderRewriteActions and PoliciesHTTPActionHeaderRewriteActionsString are defined get the list of valid values.
		Action string `validate:"required,oneof=ADD REMOVE REPLACE"`
		// Name of HTTP header
		Name string `validate:"required"`
		// Value of HTTP header
		Value string `validate:"required_if=Action ADD|REPLACE,excluded_if=Action REMOVE"`
	}

	PoliciesHTTPActionLocationRewrite struct {
		// Protocol is HTTP or HTTPS
		Protocol string `validate:"required,oneof=HTTP HTTPS"`
		// Host to which redirect the request. Default is the original host
		Host string `validate:"omitempty"`
		// Port to which redirect the request.
		Port *int `validate:"omitempty,tcp_udp_port"`
		// Path to which redirect the request. Default is the original path
		Path string `validate:"omitempty"`
		// Keep or drop the query of the incoming request URI in the redirected URI
		KeepQuery bool `validate:"omitempty"`
	}

	PoliciesHTTPActionRedirect struct {
		// Host to which redirect the request. Default is the original host
		Host string `validate:"omitempty"`
		// Keep or drop the query of the incoming request URI in the redirected URI
		KeepQuery bool
		// Path to which redirect the request. Default is the original path
		Path string `validate:"omitempty"`
		// Port to which redirect the request.
		Port *int `validate:"required,tcp_udp_port"`
		// HTTP or HTTPS protocol
		Protocol string `validate:"required,oneof=HTTP HTTPS"`
		// One of the redirect status codes - 301, 302, 307
		StatusCode int `validate:"required,oneof=301 302 307"`
	}

	PoliciesHTTPActionURLRewrite struct {
		// Host header to use for the rewritten URL.
		HostHeader string `validate:"required"`
		// Path to use for the rewritten URL.
		Path string `validate:"required"`
		// Query string to use or append to the existing query string in the rewritten URL.
		Query string `validate:"omitempty"`
		// Whether or not to keep the existing query string when rewriting the URL. Defaults to true.
		KeepQuery bool `validate:"omitempty"`
	}

	PoliciesHTTPActionRateLimit struct {
		//
		// number of requests per period allowed 1 to 1000000000
		// Default is 1000 requests
		// 1 request is the minimum
		// 1000000000 requests is the maximum
		Count int `default:"1000" validate:"min=1,max=1000000000"`
		//
		// Time period in seconds for the rate limit 1 to 1000000000
		// Default is 60 seconds
		// 1 second is the minimum period
		// 1000000000 seconds is the maximum period
		Period int `default:"60" validate:"min=1,max=1000000000"`
		//
		// Action to do an HTTP redirect when the rate limit is exceeded
		// It can't be configured in combination with other actions below
		RedirectAction *PoliciesHTTPActionRedirect `validate:"omitempty"`
		//
		// Action to close the connection HTTP when the rate limit is exceeded
		// The network connection is closed (no error http return).
		// It can't be configured in combination with other actions
		CloseConnectionAction *bool `validate:"omitempty"`
		//
		// You can use this action to send a custom response to the client when the rate limit is exceeded.
		// It can't be configured in combination with other actions
		LocalResponseAction *PoliciesHTTPActionSendResponse `validate:"omitempty"`
	}

	PoliciesHTTPActionSendResponse struct {
		// HTTP status code to return
		StatusCode int `validate:"required,oneof=200 204 403 404 429 501"`
		// Content type of the response
		ContentType string `validate:"required,oneof=application/json text/html text/plain"`
		// Content of the response - base64 encoded string
		Content string `validate:"required,base64"`
	}
)

// * Helpers to convert PoliciesHTTPActionHeaderRewrite to and from vCD types

func (PoliciesHTTPActionHeaderRewrite) fromVCD(action *govcdtypes.AlbVsHttpRequestRuleHeaderActions) *PoliciesHTTPActionHeaderRewrite {
	if action == nil {
		return nil
	}

	return &PoliciesHTTPActionHeaderRewrite{
		Action: action.Action,
		Name:   action.Name,
		Value:  action.Value,
	}
}

func (p *PoliciesHTTPActionHeaderRewrite) toVCD() *govcdtypes.AlbVsHttpRequestRuleHeaderActions {
	if p == nil {
		return nil
	}

	return &govcdtypes.AlbVsHttpRequestRuleHeaderActions{
		Action: p.Action,
		Name:   p.Name,
		Value:  p.Value,
	}
}

// * Helpers to convert PoliciesHTTPActionHeadersRewrite to and from vCD types

func (PoliciesHTTPActionHeadersRewrite) fromVCD(action []*govcdtypes.AlbVsHttpRequestRuleHeaderActions) PoliciesHTTPActionHeadersRewrite {
	var headers []*PoliciesHTTPActionHeaderRewrite
	for _, h := range action {
		headers = append(headers, (&PoliciesHTTPActionHeaderRewrite{}).fromVCD(h))
	}

	return headers
}

func (p PoliciesHTTPActionHeadersRewrite) toVCD() []*govcdtypes.AlbVsHttpRequestRuleHeaderActions {
	var headers []*govcdtypes.AlbVsHttpRequestRuleHeaderActions
	for _, h := range p {
		headers = append(headers, h.toVCD())
	}

	return headers
}

// * Helpers to convert PoliciesHTTPActionURLRewrite to and from vCD types

func (PoliciesHTTPActionURLRewrite) fromVCD(action *govcdtypes.AlbVsHttpRequestRuleRewriteURLAction) *PoliciesHTTPActionURLRewrite {
	if action == nil {
		return nil
	}

	return &PoliciesHTTPActionURLRewrite{
		HostHeader: action.Host,
		Path:       action.Path,
		Query:      action.Query,
		KeepQuery:  action.KeepQuery,
	}
}

func (p *PoliciesHTTPActionURLRewrite) toVCD() *govcdtypes.AlbVsHttpRequestRuleRewriteURLAction {
	if p == nil {
		return nil
	}

	return &govcdtypes.AlbVsHttpRequestRuleRewriteURLAction{
		Host:      p.HostHeader,
		Path:      p.Path,
		Query:     p.Query,
		KeepQuery: p.KeepQuery,
	}
}

// * Helpers to convert PoliciesHTTPActionLocationRewrite to and from vCD types.

func (p *PoliciesHTTPActionLocationRewrite) fromVCD(action *govcdtypes.AlbVsHttpRespRuleRewriteLocationHeaderAction) *PoliciesHTTPActionLocationRewrite {
	if action == nil {
		return nil
	}

	return &PoliciesHTTPActionLocationRewrite{
		Protocol:  action.Protocol,
		Host:      action.Host,
		Port:      action.Port,
		Path:      action.Path,
		KeepQuery: action.KeepQuery,
	}
}

func (p *PoliciesHTTPActionLocationRewrite) toVCD() *govcdtypes.AlbVsHttpRespRuleRewriteLocationHeaderAction {
	if p == nil {
		return nil
	}

	return &govcdtypes.AlbVsHttpRespRuleRewriteLocationHeaderAction{
		Protocol:  p.Protocol,
		Host:      p.Host,
		Port:      p.Port,
		Path:      p.Path,
		KeepQuery: p.KeepQuery,
	}
}

// * Helpers to convert PoliciesHTTPActionRateLimit to and from vCD types

func (PoliciesHTTPActionRateLimit) fromVCD(action *govcdtypes.AlbVsHttpSecurityRuleRateLimitAction) *PoliciesHTTPActionRateLimit {
	if action == nil {
		return nil
	}
	return &PoliciesHTTPActionRateLimit{
		Count:  action.Count,
		Period: action.Period,
		RedirectAction: func() *PoliciesHTTPActionRedirect {
			if action.RedirectAction != nil {
				return (&PoliciesHTTPActionRedirect{}).fromVCD(action.RedirectAction)
			}
			return nil
		}(),
		LocalResponseAction: func() *PoliciesHTTPActionSendResponse {
			if action.LocalResponseAction != nil {
				return (&PoliciesHTTPActionSendResponse{}).fromVCD(action.LocalResponseAction)
			}
			return nil
		}(),
		CloseConnectionAction: func() *bool {
			if action.CloseConnectionAction == string(PoliciesHTTPConnectionActionCLOSE) {
				return utils.ToPTR(true)
			}
			return nil
		}(),
	}
}

func (p *PoliciesHTTPActionRateLimit) toVCD() *govcdtypes.AlbVsHttpSecurityRuleRateLimitAction {
	if p == nil {
		return nil
	}
	return &govcdtypes.AlbVsHttpSecurityRuleRateLimitAction{
		Count:  p.Count,
		Period: p.Period,
		CloseConnectionAction: func() string {
			if p.CloseConnectionAction != nil && *p.CloseConnectionAction {
				return string(PoliciesHTTPConnectionActionCLOSE)
			}
			return ""
		}(),
		RedirectAction:      p.RedirectAction.toVCD(),
		LocalResponseAction: p.LocalResponseAction.toVCD(),
	}
}

// * Helpers to convert PoliciesHTTPActionSendResponse to and from vCD types

func (PoliciesHTTPActionSendResponse) fromVCD(action *govcdtypes.AlbVsHttpSecurityRuleRateLimitLocalResponseAction) *PoliciesHTTPActionSendResponse {
	if action == nil {
		return nil
	}
	return &PoliciesHTTPActionSendResponse{
		StatusCode:  action.StatusCode,
		ContentType: action.ContentType,
		Content:     action.Content,
	}
}

func (p *PoliciesHTTPActionSendResponse) toVCD() *govcdtypes.AlbVsHttpSecurityRuleRateLimitLocalResponseAction {
	if p == nil {
		return nil
	}
	return &govcdtypes.AlbVsHttpSecurityRuleRateLimitLocalResponseAction{
		Content:     p.Content,
		ContentType: p.ContentType,
		StatusCode:  p.StatusCode,
	}
}
