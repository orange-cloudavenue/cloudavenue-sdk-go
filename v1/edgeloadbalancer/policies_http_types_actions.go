package edgeloadbalancer

import govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

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
