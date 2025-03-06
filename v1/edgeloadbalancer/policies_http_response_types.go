package edgeloadbalancer

import (
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"
)

type (
	PoliciesHTTPResponseModel struct {
		VirtualServiceID string `validate:"required,urn=loadBalancerVirtualService"`
		// List of HTTP Response Policies
		Policies []*PoliciesHTTPResponseModelPolicy `validate:"required,dive"`
	}

	PoliciesHTTPResponseModelPolicy struct {
		// Name of the rule
		Name string `validate:"required"`
		// Whether the rule is active or not.
		Active bool `validate:"omitempty"`
		// Whether to enable logging with headers on rule match or not
		Logging bool `validate:"omitempty"`
		// MatchCriteria for the HTTP Response
		MatchCriteria PoliciesHTTPResponseMatchCriteria `validate:"required"`
		// HTTP header rewrite action
		// It can be configured in combination with rewrite URL action
		HeaderRewriteActions PoliciesHTTPActionHeadersRewrite `validate:"omitempty"`
		// HTTP location rewrite action
		LocationRewriteAction *PoliciesHTTPActionLocationRewrite `validate:"omitempty"`
	}

	PoliciesHTTPResponseMatchCriteria struct {
		// Client IP addresses.
		ClientIPMatch *PoliciesHTTPClientIPMatch `validate:"omitempty"`
		// Virtual service ports.
		ServicePortMatch *PoliciesHTTPServicePortMatch `validate:"omitempty"`
		// HTTP methods such as GET, PUT, DELETE, POST etc.
		MethodMatch *PoliciesHTTPMethodMatch `validate:"omitempty"`
		// Protocol
		// var PoliciesHTTPPathMatchCriteria and PoliciesHTTPPathMatchCriteriaString are defined get the list of valid values.
		Protocol string `validate:"omitempty,oneof=HTTP HTTPS"`
		// Path Match
		PathMatch *PoliciesHTTPPathMatch `validate:"omitempty"`
		// HTTP request query strings in key=value format
		QueryMatch []string `validate:"omitempty,dive,str_key_value"`
		// HTTP request cookies
		CookieMatch *PoliciesHTTPCookieMatch `validate:"omitempty"`

		// Defines match criteria based on response location header
		LocationMatch *PoliciesHTTPLocationMatch `validate:"omitempty"`
		// Defines match criteria based on the request headers
		RequestHeaderMatch PoliciesHTTPHeadersMatch `validate:"omitempty"`
		// Defines match criteria based on the response headers
		ResponseHeaderMatch PoliciesHTTPHeadersMatch `validate:"omitempty"`
		// Defines match criteria based on response status codes
		StatusCodeMatch *PoliciesHTTPStatusCodeMatch `validate:"omitempty"`
	}
)

// * Helpers to convert PoliciesHTTPResponseModel to and from vCD types

// fromVCD converts the vCD type to the model.
func (p *PoliciesHTTPResponseModel) fromVCD(virtualServiceID string, rules *govcdtypes.AlbVsHttpResponseRules) *PoliciesHTTPResponseModel {
	x := &PoliciesHTTPResponseModel{
		VirtualServiceID: virtualServiceID,
	}

	for _, rule := range rules.Values {
		x.Policies = append(x.Policies, (&PoliciesHTTPResponseModelPolicy{}).fromVCD(rule))
	}

	return x
}

// toVCD converts the model to the vCD type.
func (p *PoliciesHTTPResponseModel) toVCD() *govcdtypes.AlbVsHttpResponseRules {
	m := &govcdtypes.AlbVsHttpResponseRules{}

	var rules []govcdtypes.AlbVsHttpResponseRule
	for _, policy := range p.Policies {
		rules = append(rules, policy.toVCD())
	}

	m.Values = rules

	return m
}

// * Helpers to convert PoliciesHTTPResponseModelPolicy to and from vCD types.
func (p *PoliciesHTTPResponseModelPolicy) fromVCD(rule govcdtypes.AlbVsHttpResponseRule) *PoliciesHTTPResponseModelPolicy {
	return &PoliciesHTTPResponseModelPolicy{
		Name:                  rule.Name,
		Active:                rule.Active,
		Logging:               rule.Logging,
		MatchCriteria:         (PoliciesHTTPResponseMatchCriteria{}).fromVCD(rule.MatchCriteria),
		HeaderRewriteActions:  PoliciesHTTPActionHeadersRewrite{}.fromVCD(rule.HeaderActions),
		LocationRewriteAction: (&PoliciesHTTPActionLocationRewrite{}).fromVCD(rule.RewriteLocationHeaderAction),
	}
}

func (p *PoliciesHTTPResponseModelPolicy) toVCD() govcdtypes.AlbVsHttpResponseRule {
	return govcdtypes.AlbVsHttpResponseRule{
		Name:                        p.Name,
		Active:                      p.Active,
		Logging:                     p.Logging,
		MatchCriteria:               p.MatchCriteria.toVCD(),
		HeaderActions:               p.HeaderRewriteActions.toVCD(),
		RewriteLocationHeaderAction: p.LocationRewriteAction.toVCD(),
	}
}

// * Helpers to convert PoliciesHTTPResponseMatchCriteria to and from vCD types.
func (p PoliciesHTTPResponseMatchCriteria) fromVCD(criteria govcdtypes.AlbVsHttpResponseRuleMatchCriteria) PoliciesHTTPResponseMatchCriteria {
	return PoliciesHTTPResponseMatchCriteria{
		ClientIPMatch:       (&PoliciesHTTPClientIPMatch{}).fromVCD(criteria.ClientIPMatch),
		ServicePortMatch:    (&PoliciesHTTPServicePortMatch{}).fromVCD(criteria.ServicePortMatch),
		MethodMatch:         (&PoliciesHTTPMethodMatch{}).fromVCD(criteria.MethodMatch),
		Protocol:            criteria.Protocol,
		PathMatch:           (&PoliciesHTTPPathMatch{}).fromVCD(criteria.PathMatch),
		QueryMatch:          criteria.QueryMatch,
		CookieMatch:         (&PoliciesHTTPCookieMatch{}).fromVCD(criteria.CookieMatch),
		LocationMatch:       (&PoliciesHTTPLocationMatch{}).fromVCD(criteria.LocationHeaderMatch),
		RequestHeaderMatch:  (&PoliciesHTTPHeadersMatch{}).fromVCD(criteria.RequestHeaderMatch),
		ResponseHeaderMatch: (&PoliciesHTTPHeadersMatch{}).fromVCD(criteria.ResponseHeaderMatch),
		StatusCodeMatch:     PoliciesHTTPStatusCodeMatch{}.fromVCD(criteria.StatusCodeMatch),
	}
}

func (p *PoliciesHTTPResponseMatchCriteria) toVCD() govcdtypes.AlbVsHttpResponseRuleMatchCriteria {
	return govcdtypes.AlbVsHttpResponseRuleMatchCriteria{
		ClientIPMatch:       p.ClientIPMatch.toVCD(),
		ServicePortMatch:    p.ServicePortMatch.toVCD(),
		MethodMatch:         p.MethodMatch.toVCD(),
		Protocol:            p.Protocol,
		PathMatch:           p.PathMatch.toVCD(),
		QueryMatch:          p.QueryMatch,
		CookieMatch:         p.CookieMatch.toVCD(),
		LocationHeaderMatch: p.LocationMatch.toVCD(),
		RequestHeaderMatch:  p.RequestHeaderMatch.toVCD(),
		ResponseHeaderMatch: p.ResponseHeaderMatch.toVCD(),
		StatusCodeMatch:     p.StatusCodeMatch.toVCD(),
	}
}
