/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgeloadbalancer

import govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

type (
	PoliciesHTTPProtocol                  string
	PoliciesHTTPMethod                    string
	PoliciesHTTPMatchCriteriaCriteria     string
	PoliciesHTTPActionHeaderRewriteAction string
)

type (
	PoliciesHTTPClientIPMatch struct {
		// Criteria to use for IP address matching the HTTP request.
		// var PoliciesHTTPClientIPMatchCriteria and PoliciesHTTPClientIPMatchCriteriaString are defined get the list of valid values.
		Criteria string `validate:"required,oneof=IS_IN IS_NOT_IN"`
		// Either a single IP address, a range of IP addresses or a network CIDR. Must contain at least
		// one item.
		Addresses []string `validate:"required,dive,ipv4|cidr|ipv4_range"`
	}

	PoliciesHTTPServicePortMatch struct {
		// Criteria to use for port matching the HTTP request.
		// var PoliciesHTTPServicePortMatchCriteria and PoliciesHTTPServicePortMatchCriteriaString are defined get the list of valid values.
		Criteria string `validate:"required,oneof=IS_IN IS_NOT_IN"`
		// Listening TCP ports.
		Ports []int `validate:"required,dive,tcp_udp_port"`
	}

	PoliciesHTTPMethodMatch struct {
		// Criteria to use for HTTP method matching the HTTP request.
		// var PoliciesHTTPMethodMatchCriteria and PoliciesHTTPMethodMatchCriteriaString are defined get the list of valid values.
		Criteria string `validate:"required,oneof=IS_IN IS_NOT_IN"`
		// HTTP methods.
		// var PoliciesHTTPMethodsMatch and PoliciesHTTPMethodsMatchString are defined get the list of valid values.
		Methods []string `validate:"required,dive,oneof=GET POST PUT DELETE PATCH OPTIONS TRACE CONNECT PROPFIND PROPPATCH MKCOL COPY MOVE LOCK UNLOCK"`
	}

	PoliciesHTTPPathMatch struct {
		// Criteria to use for path matching the HTTP request.
		// var PoliciesHTTPPathMatchCriteria and PoliciesHTTPPathMatchCriteriaString are defined get the list of valid values.
		Criteria string `validate:"required,oneof=BEGINS_WITH DOES_NOT_BEGIN_WITH CONTAINS DOES_NOT_CONTAIN ENDS_WITH DOES_NOT_END_WITH EQUALS DOES_NOT_EQUAL REGEX_MATCH REGEX_DOES_NOT_MATCH"`
		// String values to match the path
		MatchStrings []string `validate:"required"`
	}

	PoliciesHTTPHeadersMatch []PoliciesHTTPHeaderMatch

	PoliciesHTTPHeaderMatch struct {
		// Criteria to use for header matching the HTTP request.
		// var PoliciesHTTPHeaderMatchCriteria and PoliciesHTTPHeaderMatchCriteriaString are defined get the list of valid values.
		Criteria string `validate:"required,oneof=BEGINS_WITH DOES_NOT_BEGIN_WITH CONTAINS DOES_NOT_CONTAIN ENDS_WITH DOES_NOT_END_WITH EQUALS DOES_NOT_EQUAL EXISTS DOES_NOT_EXIST"`
		// Name of the HTTP header whose value is to be matched
		Name string `validate:"required"`
		// String values to match for an HTTP header
		Values []string `validate:"required|excluded_if=Criteria EXISTS DOES_NOT_EXIST"`
	}

	PoliciesHTTPCookieMatch struct {
		// Criteria to use for cookie matching the HTTP request.
		// var PoliciesHTTPCookieMatchCriteria and PoliciesHTTPCookieMatchCriteriaString are defined get the list of valid values.
		Criteria string `validate:"required,oneof=BEGINS_WITH DOES_NOT_BEGIN_WITH CONTAINS DOES_NOT_CONTAIN ENDS_WITH DOES_NOT_END_WITH EQUALS DOES_NOT_EQUAL"`
		// Name of the cookie whose value is to be matched
		Name string `validate:"required"`
		// String values to match for a cookie.Value length should be less than 10240
		Value string `validate:"required|excluded_if=Criteria EXISTS DOES_NOT_EXIST"`
	}

	PoliciesHTTPStatusCodeMatch struct {
		// Criteria to use for matching the HTTP response status code
		Criteria string `validate:"required,oneof=IS_IN IS_NOT_IN"`
		// StatusCode is single status codes or ranges of status codes separated by a hyphen.
		// For example, "200-299" will match all status codes between 200 and 299.
		StatusCodes []string `validate:"required,dive,http_status_code|http_status_code_range"`
	}

	PoliciesHTTPLocationMatch struct {
		// Criteria to use for location header matching the HTTP response.
		// var PoliciesHTTPResponseLocationHeaderMatchCriteria and PoliciesHTTPResponseLocationHeaderMatchCriteriaString are defined get the list of valid values.
		Criteria string `validate:"required,oneof=BEGINS_WITH DOES_NOT_BEGIN_WITH CONTAINS DOES_NOT_CONTAIN ENDS_WITH DOES_NOT_END_WITH EQUALS DOES_NOT_EQUAL REGEX_MATCH REGEX_DOES_NOT_MATCH"`
		// String values to match for an HTTP header
		Values []string `validate:"required"`
	}
)

// * Helpers to convert PoliciesHTTPClientIPMatch to and from vCD types

func (PoliciesHTTPClientIPMatch) fromVCD(match *govcdtypes.AlbVsHttpRequestRuleClientIPMatch) *PoliciesHTTPClientIPMatch {
	if match == nil {
		return nil
	}

	return &PoliciesHTTPClientIPMatch{
		Criteria:  match.MatchCriteria,
		Addresses: match.Addresses,
	}
}

func (p *PoliciesHTTPClientIPMatch) toVCD() *govcdtypes.AlbVsHttpRequestRuleClientIPMatch {
	if p == nil {
		return nil
	}

	return &govcdtypes.AlbVsHttpRequestRuleClientIPMatch{
		MatchCriteria: p.Criteria,
		Addresses:     p.Addresses,
	}
}

// * Helpers to convert PoliciesHTTPServicePortMatch to and from vCD types

func (PoliciesHTTPServicePortMatch) fromVCD(match *govcdtypes.AlbVsHttpRequestRuleServicePortMatch) *PoliciesHTTPServicePortMatch {
	if match == nil {
		return nil
	}

	return &PoliciesHTTPServicePortMatch{
		Criteria: match.MatchCriteria,
		Ports:    match.Ports,
	}
}

func (p *PoliciesHTTPServicePortMatch) toVCD() *govcdtypes.AlbVsHttpRequestRuleServicePortMatch {
	if p == nil {
		return nil
	}

	return &govcdtypes.AlbVsHttpRequestRuleServicePortMatch{
		MatchCriteria: p.Criteria,
		Ports:         p.Ports,
	}
}

// * Helpers to convert PoliciesHTTPMethodMatch to and from vCD types

func (PoliciesHTTPMethodMatch) fromVCD(match *govcdtypes.AlbVsHttpRequestRuleMethodMatch) *PoliciesHTTPMethodMatch {
	if match == nil {
		return nil
	}

	return &PoliciesHTTPMethodMatch{
		Criteria: match.MatchCriteria,
		Methods:  match.Methods,
	}
}

func (p *PoliciesHTTPMethodMatch) toVCD() *govcdtypes.AlbVsHttpRequestRuleMethodMatch {
	if p == nil {
		return nil
	}

	return &govcdtypes.AlbVsHttpRequestRuleMethodMatch{
		MatchCriteria: p.Criteria,
		Methods:       p.Methods,
	}
}

// * Helpers to convert PoliciesHTTPPathMatch to and from vCD types

func (PoliciesHTTPPathMatch) fromVCD(match *govcdtypes.AlbVsHttpRequestRulePathMatch) *PoliciesHTTPPathMatch {
	if match == nil {
		return nil
	}

	return &PoliciesHTTPPathMatch{
		Criteria:     match.MatchCriteria,
		MatchStrings: match.MatchStrings,
	}
}

func (p *PoliciesHTTPPathMatch) toVCD() *govcdtypes.AlbVsHttpRequestRulePathMatch {
	if p == nil {
		return nil
	}

	return &govcdtypes.AlbVsHttpRequestRulePathMatch{
		MatchCriteria: p.Criteria,
		MatchStrings:  p.MatchStrings,
	}
}

// * Helpers to convert PoliciesHTTPHeaderMatch to and from vCD types

func (PoliciesHTTPHeaderMatch) fromVCD(match *govcdtypes.AlbVsHttpRequestRuleHeaderMatch) *PoliciesHTTPHeaderMatch {
	if match == nil {
		return nil
	}

	return &PoliciesHTTPHeaderMatch{
		Criteria: match.MatchCriteria,
		Name:     match.Key,
		Values:   match.Value,
	}
}

func (p *PoliciesHTTPHeaderMatch) toVCD() *govcdtypes.AlbVsHttpRequestRuleHeaderMatch {
	if p == nil {
		return nil
	}

	return &govcdtypes.AlbVsHttpRequestRuleHeaderMatch{
		MatchCriteria: p.Criteria,
		Key:           p.Name,
		Value:         p.Values,
	}
}

// * Helpers to convert PoliciesHTTPHeadersMatch to and from vCD types

func (PoliciesHTTPHeadersMatch) fromVCD(match []govcdtypes.AlbVsHttpRequestRuleHeaderMatch) PoliciesHTTPHeadersMatch {
	var headers []PoliciesHTTPHeaderMatch
	for _, h := range match {
		headers = append(headers, *PoliciesHTTPHeaderMatch{}.fromVCD(&h))
	}

	return headers
}

func (p PoliciesHTTPHeadersMatch) toVCD() []govcdtypes.AlbVsHttpRequestRuleHeaderMatch {
	var headers []govcdtypes.AlbVsHttpRequestRuleHeaderMatch
	for _, h := range p {
		headers = append(headers, *h.toVCD())
	}

	return headers
}

// * Helpers to convert PoliciesHTTPCookieMatch to and from vCD types

func (PoliciesHTTPCookieMatch) fromVCD(match *govcdtypes.AlbVsHttpRequestRuleCookieMatch) *PoliciesHTTPCookieMatch {
	if match == nil {
		return nil
	}

	return &PoliciesHTTPCookieMatch{
		Criteria: match.MatchCriteria,
		Name:     match.Key,
		Value:    match.Value,
	}
}

func (p *PoliciesHTTPCookieMatch) toVCD() *govcdtypes.AlbVsHttpRequestRuleCookieMatch {
	if p == nil {
		return nil
	}

	return &govcdtypes.AlbVsHttpRequestRuleCookieMatch{
		MatchCriteria: p.Criteria,
		Key:           p.Name,
		Value:         p.Value,
	}
}

// * Helpers to convert PoliciesHTTPLocationMatch to and from vCD types.
func (p *PoliciesHTTPLocationMatch) fromVCD(criteria *govcdtypes.AlbVsHttpResponseLocationHeaderMatch) *PoliciesHTTPLocationMatch {
	if criteria == nil {
		return nil
	}

	return &PoliciesHTTPLocationMatch{
		Criteria: criteria.MatchCriteria,
		Values:   criteria.Value,
	}
}

func (p *PoliciesHTTPLocationMatch) toVCD() *govcdtypes.AlbVsHttpResponseLocationHeaderMatch {
	if p == nil {
		return nil
	}

	return &govcdtypes.AlbVsHttpResponseLocationHeaderMatch{
		MatchCriteria: p.Criteria,
		Value:         p.Values,
	}
}

// * Helpers to convert PoliciesHTTPStatusCodeMatch to and from vCD types

func (p PoliciesHTTPStatusCodeMatch) fromVCD(match *govcdtypes.AlbVsHttpRuleStatusCodeMatch) *PoliciesHTTPStatusCodeMatch {
	if match == nil {
		return nil
	}

	return &PoliciesHTTPStatusCodeMatch{
		Criteria:    match.MatchCriteria,
		StatusCodes: match.StatusCodes,
	}
}

func (p *PoliciesHTTPStatusCodeMatch) toVCD() *govcdtypes.AlbVsHttpRuleStatusCodeMatch {
	if p == nil {
		return nil
	}

	return &govcdtypes.AlbVsHttpRuleStatusCodeMatch{
		MatchCriteria: p.Criteria,
		StatusCodes:   p.StatusCodes,
	}
}
