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
)

type (
	PoliciesHTTPSecurityModel struct {
		VirtualServiceID string `validate:"required,urn=loadBalancerVirtualService"`
		// List of HTTP Request Policies
		Policies []*PoliciesHTTPSecurityModelPolicy `validate:"required,dive"`
	}

	PoliciesHTTPSecurityModelPolicy struct {
		// Name of the rule
		Name string `validate:"required"`
		// Whether the rule is active or not.
		Active bool `validate:"omitempty"`
		// Whether to enable logging for the rule (policy)
		Logging bool `validate:"omitempty"`
		// MatchCriteria for the HTTP Request
		MatchCriteria PoliciesHTTPSecurityMatchCriteria `validate:"required"`

		// Action to take when the rule matches

		// HTTP Connection Action
		// If set, the rule will either allow or close the connection based on the action specified.
		// It can be configured in combination with other actions
		ConnectionAction PoliciesHTTPConnectionAction `validate:"omitempty,oneof=ALLOW CLOSE"`
		// HTTP Rate Limit Action
		// If set, the rule will limit the rate of requests from a client IP address.
		// It can be configured in combination with other actions
		RateLimitAction *PoliciesHTTPActionRateLimit `validate:"omitempty"`
		// HTTP Redirect to HTTPS Action
		// If set, the rule will redirect HTTP requests to HTTPS on the specified port.
		// It can be configured in combination with other actions
		RedirectToHTTPSAction *int `validate:"omitempty"`
		// HTTP Send Response Action
		// If set, the rule will send a custom response to the client.
		// It can be configured in combination with other actions
		SendResponseAction *PoliciesHTTPActionSendResponse `validate:"omitempty"`
	}

	PoliciesHTTPSecurityMatchCriteria struct {
		// Protocol
		Protocol PoliciesHTTPProtocol `validate:"omitempty,oneof=HTTP HTTPS"`
		// Client IP addresses
		ClientIPMatch *PoliciesHTTPClientIPMatch `validate:"omitempty"`
		// Service Ports
		ServicePortMatch *PoliciesHTTPServicePortMatch `validate:"omitempty"`
		// HTTP Methods
		MethodMatch *PoliciesHTTPMethodMatch `validate:"omitempty"`
		// Path Match
		PathMatch *PoliciesHTTPPathMatch `validate:"omitempty"`
		// HTTP request cookies
		CookieMatch *PoliciesHTTPCookieMatch `validate:"omitempty"`
		// HTTP request headers
		HeaderMatch PoliciesHTTPHeadersMatch `validate:"omitempty"`
		// HTTP request query strings in key=value format
		QueryMatch []string `validate:"omitempty,dive,str_key_value"`
	}
)

// * Helpers to convert PoliciesHTTPSecurityModel to and from vCD types

func (PoliciesHTTPSecurityModel) fromVCD(virtualServiceID string, rules []*govcdtypes.AlbVsHttpSecurityRule) *PoliciesHTTPSecurityModel {
	x := &PoliciesHTTPSecurityModel{
		VirtualServiceID: virtualServiceID,
	}

	for _, rule := range rules {
		x.Policies = append(x.Policies, (&PoliciesHTTPSecurityModelPolicy{}).fromVCD(rule))
	}
	return x
}

func (p *PoliciesHTTPSecurityModel) toVCD() *govcdtypes.AlbVsHttpSecurityRules {
	m := &govcdtypes.AlbVsHttpSecurityRules{}

	for _, policy := range p.Policies {
		m.Values = append(m.Values, policy.toVCD())
	}

	return m
}

// * Helpers to convert PoliciesHTTPSecurityModelPolicy to and from vCD types

func (PoliciesHTTPSecurityModelPolicy) fromVCD(rule *govcdtypes.AlbVsHttpSecurityRule) *PoliciesHTTPSecurityModelPolicy {
	return &PoliciesHTTPSecurityModelPolicy{
		Name:             rule.Name,
		Active:           rule.Active,
		Logging:          rule.Logging,
		MatchCriteria:    PoliciesHTTPSecurityMatchCriteria{}.fromVCD(rule.MatchCriteria),
		ConnectionAction: PoliciesHTTPConnectionAction(rule.AllowOrCloseConnectionAction),
		RedirectToHTTPSAction: func() *int {
			if rule.RedirectToHTTPSAction != nil {
				return &rule.RedirectToHTTPSAction.Port
			}
			return nil
		}(),
		SendResponseAction: (&PoliciesHTTPActionSendResponse{}).fromVCD(rule.LocalResponseAction),
		RateLimitAction:    (&PoliciesHTTPActionRateLimit{}).fromVCD(rule.RateLimitAction),
	}
}

func (p *PoliciesHTTPSecurityModelPolicy) toVCD() govcdtypes.AlbVsHttpSecurityRule {
	return govcdtypes.AlbVsHttpSecurityRule{
		Name:                         p.Name,
		Active:                       p.Active,
		Logging:                      p.Logging,
		MatchCriteria:                p.MatchCriteria.toVCD(),
		AllowOrCloseConnectionAction: string(p.ConnectionAction),
		RedirectToHTTPSAction: func() *govcdtypes.AlbVsHttpSecurityRuleRedirectToHTTPSAction {
			if p.RedirectToHTTPSAction != nil {
				return &govcdtypes.AlbVsHttpSecurityRuleRedirectToHTTPSAction{Port: *p.RedirectToHTTPSAction}
			}
			return nil
		}(),
		LocalResponseAction: p.SendResponseAction.toVCD(),
		RateLimitAction:     p.RateLimitAction.toVCD(),
	}
}

// * Helpers to convert PoliciesHTTPSecurityMatchCriteria to and from vCD types

func (PoliciesHTTPSecurityMatchCriteria) fromVCD(criteria govcdtypes.AlbVsHttpRequestAndSecurityRuleMatchCriteria) PoliciesHTTPSecurityMatchCriteria {
	return PoliciesHTTPSecurityMatchCriteria{
		Protocol:         PoliciesHTTPProtocol(criteria.Protocol),
		ClientIPMatch:    (&PoliciesHTTPClientIPMatch{}).fromVCD(criteria.ClientIPMatch),
		ServicePortMatch: (&PoliciesHTTPServicePortMatch{}).fromVCD(criteria.ServicePortMatch),
		MethodMatch:      (&PoliciesHTTPMethodMatch{}).fromVCD(criteria.MethodMatch),
		PathMatch:        (&PoliciesHTTPPathMatch{}).fromVCD(criteria.PathMatch),
		CookieMatch:      (&PoliciesHTTPCookieMatch{}).fromVCD(criteria.CookieMatch),
		HeaderMatch:      (&PoliciesHTTPHeadersMatch{}).fromVCD(criteria.HeaderMatch),
		QueryMatch:       criteria.QueryMatch,
	}
}

func (p PoliciesHTTPSecurityMatchCriteria) toVCD() govcdtypes.AlbVsHttpRequestAndSecurityRuleMatchCriteria {
	return govcdtypes.AlbVsHttpRequestAndSecurityRuleMatchCriteria{
		Protocol:         string(p.Protocol),
		ClientIPMatch:    p.ClientIPMatch.toVCD(),
		ServicePortMatch: p.ServicePortMatch.toVCD(),
		MethodMatch:      p.MethodMatch.toVCD(),
		PathMatch:        p.PathMatch.toVCD(),
		CookieMatch:      p.CookieMatch.toVCD(),
		HeaderMatch:      p.HeaderMatch.toVCD(),
		QueryMatch:       p.QueryMatch,
	}
}
