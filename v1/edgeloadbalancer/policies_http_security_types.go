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
		// It can be configured in combination with other actions
		ConnectionAction string `validate:"omitempty,oneof=ALLOW CLOSE"`
		// HTTP Rate Limit Action
		// It can be configured in combination with other actions
		RateLimitAction *PoliciesHTTPActionRateLimit `validate:"omitempty"`
		// HTTP Redirect to HTTPS Action
		// It can be configured in combination with other actions
		// RedirectToHTTPSAction *PoliciesHTTPActionRedirectToHTTPS `validate:"omitempty"`
		RedirectToHTTPSAction *int `validate:"omitempty"`
		// HTTP Send Response Action
		// It can be configured in combination with other actions
		SendResponseAction *PoliciesHTTPActionSendResponse `validate:"omitempty"`
	}

	PoliciesHTTPSecurityMatchCriteria struct {
		// Protocol
		Protocol string `validate:"omitempty,oneof=HTTP HTTPS"`
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

func (PoliciesHTTPSecurityModel) fromVCD(virtualServiceID string, rules *govcdtypes.AlbVsHttpSecurityRules) *PoliciesHTTPSecurityModel {
	x := &PoliciesHTTPSecurityModel{
		VirtualServiceID: virtualServiceID,
	}

	for _, rule := range rules.Values {
		x.Policies = append(x.Policies, (&PoliciesHTTPSecurityModelPolicy{}).fromVCD(&rule))
	}
	return x
}

func (p *PoliciesHTTPSecurityModel) toVCD() *govcdtypes.AlbVsHttpSecurityRules {
	m := &govcdtypes.AlbVsHttpSecurityRules{}

	var rules []govcdtypes.AlbVsHttpSecurityRule
	for _, policy := range p.Policies {
		rules = append(rules, policy.toVCD())
	}

	m.Values = rules

	return m
}

// * Helpers to convert PoliciesHTTPSecurityModelPolicy to and from vCD types

func (PoliciesHTTPSecurityModelPolicy) fromVCD(rule *govcdtypes.AlbVsHttpSecurityRule) *PoliciesHTTPSecurityModelPolicy {
	return &PoliciesHTTPSecurityModelPolicy{
		Name:             rule.Name,
		Active:           rule.Active,
		Logging:          rule.Logging,
		MatchCriteria:    PoliciesHTTPSecurityMatchCriteria{}.fromVCD(rule.MatchCriteria),
		ConnectionAction: rule.AllowOrCloseConnectionAction,
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
		AllowOrCloseConnectionAction: p.ConnectionAction,
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
		Protocol:         criteria.Protocol,
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
		Protocol:         p.Protocol,
		ClientIPMatch:    p.ClientIPMatch.toVCD(),
		ServicePortMatch: p.ServicePortMatch.toVCD(),
		MethodMatch:      p.MethodMatch.toVCD(),
		PathMatch:        p.PathMatch.toVCD(),
		CookieMatch:      p.CookieMatch.toVCD(),
		HeaderMatch:      p.HeaderMatch.toVCD(),
		QueryMatch:       p.QueryMatch,
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
		CloseConnectionAction: func() string {
			if action.CloseConnectionAction != "" {
				return "Close_Connection"
			}
			return ""
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
			return p.CloseConnectionAction
		}(),
		RedirectAction: func() *govcdtypes.AlbVsHttpRequestRuleRedirectAction {
			if p.RedirectAction != nil {
				return p.RedirectAction.toVCD()
			}
			return nil
		}(),
		LocalResponseAction: func() *govcdtypes.AlbVsHttpSecurityRuleRateLimitLocalResponseAction {
			if p.LocalResponseAction != nil {
				return p.LocalResponseAction.toVCD()
			}
			return nil
		}(),
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
