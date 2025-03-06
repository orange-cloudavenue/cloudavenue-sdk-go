package edgeloadbalancer

import (
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"
)

type (
	PoliciesHTTPRequestModel struct {
		VirtualServiceID string `validate:"required,urn=loadBalancerVirtualService"`
		// List of HTTP Request Policies
		Policies []*PoliciesHTTPRequestModelPolicy `validate:"required,dive"`
	}

	PoliciesHTTPRequestModelPolicy struct {
		// Name of the rule
		Name string `validate:"required"`
		// Whether the rule is active or not.
		Active bool `validate:"omitempty"`
		// Whether to enable logging for the rule (policy)
		Logging bool `validate:"omitempty"`
		// MatchCriteria for the HTTP Request
		MatchCriteria PoliciesHTTPRequestMatchCriteria `validate:"required"`

		// Action to take when the rule matches

		// HTTP Redirect Action
		// It cannot be configured in combination with other actions
		RedirectAction *PoliciesHTTPActionRedirect `validate:"omitempty"`
		// HTTP header rewrite action
		// It can be configured in combination with rewrite URL action
		HeaderRewriteActions PoliciesHTTPActionHeadersRewrite `validate:"omitempty"`
		// HTTP request URL rewrite action
		// It can be configured in combination with multiple header actions
		URLRewriteAction *PoliciesHTTPActionURLRewrite `validate:"omitempty"`
	}

	PoliciesHTTPRequestMatchCriteria struct {
		// Protocol
		// var PoliciesHTTPMatchCriteriaProtocols and PoliciesHTTPMatchCriteriaProtocolsString are defined get the list of valid values.
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

// * Helpers to convert PoliciesHTTPRequestModel to and from vCD types

func (PoliciesHTTPRequestModel) fromVCD(virtualServiceID string, rules *govcdtypes.AlbVsHttpRequestRules) *PoliciesHTTPRequestModel {
	x := &PoliciesHTTPRequestModel{
		VirtualServiceID: virtualServiceID,
	}

	for _, rule := range rules.Values {
		x.Policies = append(x.Policies, (&PoliciesHTTPRequestModelPolicy{}).fromVCD(&rule))
	}
	return x
}

func (p *PoliciesHTTPRequestModel) toVCD() *govcdtypes.AlbVsHttpRequestRules {
	m := &govcdtypes.AlbVsHttpRequestRules{}

	var rules []govcdtypes.AlbVsHttpRequestRule
	for _, policy := range p.Policies {
		rules = append(rules, policy.toVCD())
	}

	m.Values = rules

	return m
}

// * Helpers to convert PoliciesHTTPRequestModelPolicy to and from vCD types

func (PoliciesHTTPRequestModelPolicy) fromVCD(rule *govcdtypes.AlbVsHttpRequestRule) *PoliciesHTTPRequestModelPolicy {
	return &PoliciesHTTPRequestModelPolicy{
		Name:                 rule.Name,
		Active:               rule.Active,
		Logging:              rule.Logging,
		MatchCriteria:        PoliciesHTTPRequestMatchCriteria{}.fromVCD(rule.MatchCriteria),
		RedirectAction:       (&PoliciesHTTPActionRedirect{}).fromVCD(rule.RedirectAction),
		HeaderRewriteActions: PoliciesHTTPActionHeadersRewrite{}.fromVCD(rule.HeaderActions),
		URLRewriteAction:     (&PoliciesHTTPActionURLRewrite{}).fromVCD(rule.RewriteURLAction),
	}
}

func (p *PoliciesHTTPRequestModelPolicy) toVCD() govcdtypes.AlbVsHttpRequestRule {
	return govcdtypes.AlbVsHttpRequestRule{
		Name:             p.Name,
		Active:           p.Active,
		Logging:          p.Logging,
		MatchCriteria:    p.MatchCriteria.toVCD(),
		RedirectAction:   p.RedirectAction.toVCD(),
		HeaderActions:    p.HeaderRewriteActions.toVCD(),
		RewriteURLAction: p.URLRewriteAction.toVCD(),
	}
}

// * Helpers to convert PoliciesHTTPRequestMatchCriteria to and from vCD types

func (PoliciesHTTPRequestMatchCriteria) fromVCD(criteria govcdtypes.AlbVsHttpRequestAndSecurityRuleMatchCriteria) PoliciesHTTPRequestMatchCriteria {
	return PoliciesHTTPRequestMatchCriteria{
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

func (p PoliciesHTTPRequestMatchCriteria) toVCD() govcdtypes.AlbVsHttpRequestAndSecurityRuleMatchCriteria {
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

// * Helpers to convert PoliciesHTTPActionRedirect to and from vCD types

func (PoliciesHTTPActionRedirect) fromVCD(action *govcdtypes.AlbVsHttpRequestRuleRedirectAction) *PoliciesHTTPActionRedirect {
	if action == nil {
		return nil
	}
	return &PoliciesHTTPActionRedirect{
		Host:       action.Host,
		KeepQuery:  action.KeepQuery,
		Path:       action.Path,
		Port:       action.Port,
		Protocol:   action.Protocol,
		StatusCode: action.StatusCode,
	}
}

func (p *PoliciesHTTPActionRedirect) toVCD() *govcdtypes.AlbVsHttpRequestRuleRedirectAction {
	if p == nil {
		return nil
	}
	return &govcdtypes.AlbVsHttpRequestRuleRedirectAction{
		Host:       p.Host,
		KeepQuery:  p.KeepQuery,
		Path:       p.Path,
		Port:       p.Port,
		Protocol:   p.Protocol,
		StatusCode: p.StatusCode,
	}
}
