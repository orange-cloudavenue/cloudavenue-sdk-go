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
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/internal/utils"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

func TestClient_GetPoliciesHTTPRequest(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock client for cloudavenue.
	clientCAV := NewMockclientFake(ctrl)

	c, _ := NewFakeClient(clientCAV)

	virtualServiceID := urn.LoadBalancerVirtualService.String() + uuid.New().String()
	edgeGatewayID := urn.Gateway.String() + uuid.New().String()
	poolID := urn.LoadBalancerPool.String() + uuid.New().String()
	serviceEngineID := urn.ServiceEngineGroup.String() + uuid.New().String()

	tests := []struct {
		name             string
		virtualServiceID string
		mockFunc         func()
		expectedPolicies *PoliciesHTTPRequestModel
		expectedErr      bool
		err              error
	}{
		{
			name:             "success",
			virtualServiceID: virtualServiceID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbVirtualServiceById(virtualServiceID).Return(&govcd.NsxtAlbVirtualService{
					NsxtAlbVirtualService: &govcdtypes.NsxtAlbVirtualService{
						ID:          virtualServiceID,
						Name:        "virtualServiceName2",
						Description: "virtualServiceDescription2",
						Enabled:     utils.ToPTR(true),
						ApplicationProfile: govcdtypes.NsxtAlbVirtualServiceApplicationProfile{
							Type: "HTTPS",
						},
						GatewayRef: govcdtypes.OpenApiReference{
							ID: edgeGatewayID,
						},
						LoadBalancerPoolRef: govcdtypes.OpenApiReference{
							ID: poolID,
						},
						ServiceEngineGroupRef: govcdtypes.OpenApiReference{
							ID: serviceEngineID,
						},
						ServicePorts: []govcdtypes.NsxtAlbVirtualServicePort{
							{
								PortStart:  utils.ToPTR(443),
								PortEnd:    nil,
								SslEnabled: utils.ToPTR(true),
								TcpUdpProfile: &govcdtypes.NsxtAlbVirtualServicePortTcpUdpProfile{
									Type: "TCP_PROXY",
								},
							},
						},
						VirtualIpAddress:      "192.168.0.1",
						HealthStatus:          "UP",
						HealthMessage:         "OK",
						DetailedHealthMessage: "OK",
					},
				}, nil)
				getPoliciesHTTPRequest = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpRequestRule, error) {
					return []*govcdtypes.AlbVsHttpRequestRule{
						{
							Name:    "ruleName",
							Active:  true,
							Logging: true,
							MatchCriteria: govcdtypes.AlbVsHttpRequestAndSecurityRuleMatchCriteria{
								ClientIPMatch: &govcdtypes.AlbVsHttpRequestRuleClientIPMatch{
									MatchCriteria: "IS_IN",
									Addresses: []string{
										"12.23.34.45",
										"12.23.34.0/24",
										"12.23.34.0-12.23.34.100",
									},
								},
								ServicePortMatch: &govcdtypes.AlbVsHttpRequestRuleServicePortMatch{
									MatchCriteria: "IS_IN",
									Ports: []int{
										80,
										443,
									},
								},
								MethodMatch: &govcdtypes.AlbVsHttpRequestRuleMethodMatch{
									MatchCriteria: "IS_IN",
									Methods: []string{
										"GET",
										"POST",
									},
								},
								Protocol: "HTTP",
								PathMatch: &govcdtypes.AlbVsHttpRequestRulePathMatch{
									MatchCriteria: "BEGINS_WITH",
									MatchStrings: []string{
										"/path1",
										"/path2",
									},
								},
								QueryMatch: []string{
									"key1=value1",
									"key2=value2",
								},
								HeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{
									{
										MatchCriteria: "IS_IN",
										Key:           "User-Agent",
										Value: []string{
											"Mozilla/5.0",
											"curl/7.64.1",
										},
									},
									{
										MatchCriteria: "IS_IN",
										Key:           "Accept",
										Value: []string{
											"application/json",
											"text/html",
										},
									},
								},
								CookieMatch: &govcdtypes.AlbVsHttpRequestRuleCookieMatch{
									MatchCriteria: "BEGINS_WITH",
									Key:           "session_id",
									Value:         "abc123",
								},
							},
							RedirectAction: &govcdtypes.AlbVsHttpRequestRuleRedirectAction{
								Host:       "example.com",
								KeepQuery:  true,
								Path:       "/newpath",
								Port:       utils.ToPTR(80),
								Protocol:   "HTTP",
								StatusCode: 301,
							},
						},
					}, nil
				}
			},
			expectedPolicies: &PoliciesHTTPRequestModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPRequestModelPolicy{
					{
						Name:    "ruleName",
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPRequestMatchCriteria{
							Protocol:         "HTTP",
							ClientIPMatch:    &PoliciesHTTPClientIPMatch{Criteria: "IS_IN", Addresses: []string{"12.23.34.45", "12.23.34.0/24", "12.23.34.0-12.23.34.100"}},
							ServicePortMatch: &PoliciesHTTPServicePortMatch{Criteria: "IS_IN", Ports: []int{80, 443}},
							MethodMatch:      &PoliciesHTTPMethodMatch{Criteria: "IS_IN", Methods: []string{"GET", "POST"}},
							PathMatch:        &PoliciesHTTPPathMatch{Criteria: "BEGINS_WITH", MatchStrings: []string{"/path1", "/path2"}},
							QueryMatch:       []string{"key1=value1", "key2=value2"},
							HeaderMatch: []PoliciesHTTPHeaderMatch{
								{
									Criteria: "IS_IN",
									Name:     "User-Agent",
									Values:   []string{"Mozilla/5.0", "curl/7.64.1"},
								},
								{
									Criteria: "IS_IN",
									Name:     "Accept",
									Values:   []string{"application/json", "text/html"},
								},
							},
							CookieMatch: &PoliciesHTTPCookieMatch{Criteria: "BEGINS_WITH", Name: "session_id", Value: "abc123"},
						},
						RedirectAction: &PoliciesHTTPActionRedirect{
							Host:       "example.com",
							KeepQuery:  true,
							Path:       "/newpath",
							Port:       utils.ToPTR(80),
							Protocol:   "HTTP",
							StatusCode: 301,
						},
					},
				},
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name:             "success-only-match-protocol",
			virtualServiceID: virtualServiceID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbVirtualServiceById(virtualServiceID).Return(&govcd.NsxtAlbVirtualService{
					NsxtAlbVirtualService: &govcdtypes.NsxtAlbVirtualService{
						ID:          virtualServiceID,
						Name:        "virtualServiceName2",
						Description: "virtualServiceDescription2",
						Enabled:     utils.ToPTR(true),
						ApplicationProfile: govcdtypes.NsxtAlbVirtualServiceApplicationProfile{
							Type: "HTTPS",
						},
						GatewayRef: govcdtypes.OpenApiReference{
							ID: edgeGatewayID,
						},
						LoadBalancerPoolRef: govcdtypes.OpenApiReference{
							ID: poolID,
						},
						ServiceEngineGroupRef: govcdtypes.OpenApiReference{
							ID: serviceEngineID,
						},
						ServicePorts: []govcdtypes.NsxtAlbVirtualServicePort{
							{
								PortStart:  utils.ToPTR(443),
								PortEnd:    nil,
								SslEnabled: utils.ToPTR(true),
								TcpUdpProfile: &govcdtypes.NsxtAlbVirtualServicePortTcpUdpProfile{
									Type: "TCP_PROXY",
								},
							},
						},
						VirtualIpAddress:      "192.168.0.1",
						HealthStatus:          "UP",
						HealthMessage:         "OK",
						DetailedHealthMessage: "OK",
					},
				}, nil)
				getPoliciesHTTPRequest = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpRequestRule, error) {
					return []*govcdtypes.AlbVsHttpRequestRule{
						{
							Name:    "ruleName",
							Active:  true,
							Logging: true,
							MatchCriteria: govcdtypes.AlbVsHttpRequestAndSecurityRuleMatchCriteria{
								Protocol:    "HTTP",
								HeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{},
							},
							RedirectAction: &govcdtypes.AlbVsHttpRequestRuleRedirectAction{
								Host:       "example.com",
								KeepQuery:  true,
								Path:       "/newpath",
								Port:       utils.ToPTR(80),
								Protocol:   "HTTP",
								StatusCode: 301,
							},
						},
					}, nil
				}
			},
			expectedPolicies: &PoliciesHTTPRequestModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPRequestModelPolicy{
					{
						Name:    "ruleName",
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPRequestMatchCriteria{
							Protocol: "HTTP",
						},
						RedirectAction: &PoliciesHTTPActionRedirect{
							Host:       "example.com",
							KeepQuery:  true,
							Path:       "/newpath",
							Port:       utils.ToPTR(80),
							Protocol:   "HTTP",
							StatusCode: 301,
						},
					},
				},
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name:             "success-action-modify-headers-and-rewrite-url",
			virtualServiceID: virtualServiceID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbVirtualServiceById(virtualServiceID).Return(&govcd.NsxtAlbVirtualService{
					NsxtAlbVirtualService: &govcdtypes.NsxtAlbVirtualService{
						ID:          virtualServiceID,
						Name:        "virtualServiceName2",
						Description: "virtualServiceDescription2",
						Enabled:     utils.ToPTR(true),
						ApplicationProfile: govcdtypes.NsxtAlbVirtualServiceApplicationProfile{
							Type: "HTTPS",
						},
						GatewayRef: govcdtypes.OpenApiReference{
							ID: edgeGatewayID,
						},
						LoadBalancerPoolRef: govcdtypes.OpenApiReference{
							ID: poolID,
						},
						ServiceEngineGroupRef: govcdtypes.OpenApiReference{
							ID: serviceEngineID,
						},
						ServicePorts: []govcdtypes.NsxtAlbVirtualServicePort{
							{
								PortStart:  utils.ToPTR(443),
								PortEnd:    nil,
								SslEnabled: utils.ToPTR(true),
								TcpUdpProfile: &govcdtypes.NsxtAlbVirtualServicePortTcpUdpProfile{
									Type: "TCP_PROXY",
								},
							},
						},
						VirtualIpAddress:      "192.168.0.1",
						HealthStatus:          "UP",
						HealthMessage:         "OK",
						DetailedHealthMessage: "OK",
					},
				}, nil)
				getPoliciesHTTPRequest = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpRequestRule, error) {
					return []*govcdtypes.AlbVsHttpRequestRule{
						{
							Name:    "ruleName",
							Active:  true,
							Logging: true,
							MatchCriteria: govcdtypes.AlbVsHttpRequestAndSecurityRuleMatchCriteria{
								ClientIPMatch: &govcdtypes.AlbVsHttpRequestRuleClientIPMatch{
									MatchCriteria: "IS_IN",
									Addresses: []string{
										"12.23.34.45",
										"12.23.34.0/24",
										"12.23.34.0-12.23.34.100",
									},
								},
								ServicePortMatch: &govcdtypes.AlbVsHttpRequestRuleServicePortMatch{
									MatchCriteria: "IS_IN",
									Ports: []int{
										80,
										443,
									},
								},
								MethodMatch: &govcdtypes.AlbVsHttpRequestRuleMethodMatch{
									MatchCriteria: "IS_IN",
									Methods: []string{
										"GET",
										"POST",
									},
								},
								Protocol: "HTTP",
								PathMatch: &govcdtypes.AlbVsHttpRequestRulePathMatch{
									MatchCriteria: "BEGINS_WITH",
									MatchStrings: []string{
										"/path1",
										"/path2",
									},
								},
								QueryMatch: []string{
									"key1=value1",
									"key2=value2",
								},
								HeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{
									{
										MatchCriteria: "IS_IN",
										Key:           "User-Agent",
										Value: []string{
											"Mozilla/5.0",
											"curl/7.64.1",
										},
									},
									{
										MatchCriteria: "IS_IN",
										Key:           "Accept",
										Value: []string{
											"application/json",
											"text/html",
										},
									},
								},
								CookieMatch: &govcdtypes.AlbVsHttpRequestRuleCookieMatch{
									MatchCriteria: "BEGINS_WITH",
									Key:           "session_id",
									Value:         "abc123",
								},
							},
							HeaderActions: []*govcdtypes.AlbVsHttpRequestRuleHeaderActions{
								{
									Action: "ADD",
									Name:   "X-Forwarded-For",
									Value:  "test",
								},
								{
									Action: "REMOVE",
									Name:   "X-Forwarded-Proto",
									Value:  "",
								},
							},
							RewriteURLAction: &govcdtypes.AlbVsHttpRequestRuleRewriteURLAction{
								Host:      "example.com",
								Path:      "/newpath",
								KeepQuery: true,
							},
						},
					}, nil
				}
			},
			expectedPolicies: &PoliciesHTTPRequestModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPRequestModelPolicy{
					{
						Name:    "ruleName",
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPRequestMatchCriteria{
							Protocol:         "HTTP",
							ClientIPMatch:    &PoliciesHTTPClientIPMatch{Criteria: "IS_IN", Addresses: []string{"12.23.34.45", "12.23.34.0/24", "12.23.34.0-12.23.34.100"}},
							ServicePortMatch: &PoliciesHTTPServicePortMatch{Criteria: "IS_IN", Ports: []int{80, 443}},
							MethodMatch:      &PoliciesHTTPMethodMatch{Criteria: "IS_IN", Methods: []string{"GET", "POST"}},
							PathMatch:        &PoliciesHTTPPathMatch{Criteria: "BEGINS_WITH", MatchStrings: []string{"/path1", "/path2"}},
							QueryMatch:       []string{"key1=value1", "key2=value2"},
							HeaderMatch: []PoliciesHTTPHeaderMatch{
								{
									Criteria: "IS_IN",
									Name:     "User-Agent",
									Values:   []string{"Mozilla/5.0", "curl/7.64.1"},
								},
								{
									Criteria: "IS_IN",
									Name:     "Accept",
									Values:   []string{"application/json", "text/html"},
								},
							},
							CookieMatch: &PoliciesHTTPCookieMatch{Criteria: "BEGINS_WITH", Name: "session_id", Value: "abc123"},
						},
						HeaderRewriteActions: PoliciesHTTPActionHeadersRewrite{
							{
								Action: "ADD",
								Name:   "X-Forwarded-For",
								Value:  "test",
							},
							{
								Action: "REMOVE",
								Name:   "X-Forwarded-Proto",
								Value:  "",
							},
						},
						URLRewriteAction: &PoliciesHTTPActionURLRewrite{
							HostHeader: "example.com",
							Path:       "/newpath",
							KeepQuery:  true,
						},
					},
				},
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name:             "success-no-policies",
			virtualServiceID: virtualServiceID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbVirtualServiceById(virtualServiceID).Return(&govcd.NsxtAlbVirtualService{
					NsxtAlbVirtualService: &govcdtypes.NsxtAlbVirtualService{
						ID:          virtualServiceID,
						Name:        "virtualServiceName2",
						Description: "virtualServiceDescription2",
						Enabled:     utils.ToPTR(true),
						ApplicationProfile: govcdtypes.NsxtAlbVirtualServiceApplicationProfile{
							Type: "HTTPS",
						},
						GatewayRef: govcdtypes.OpenApiReference{
							ID: edgeGatewayID,
						},
						LoadBalancerPoolRef: govcdtypes.OpenApiReference{
							ID: poolID,
						},
						ServiceEngineGroupRef: govcdtypes.OpenApiReference{
							ID: serviceEngineID,
						},
						ServicePorts: []govcdtypes.NsxtAlbVirtualServicePort{
							{
								PortStart:  utils.ToPTR(443),
								PortEnd:    nil,
								SslEnabled: utils.ToPTR(true),
								TcpUdpProfile: &govcdtypes.NsxtAlbVirtualServicePortTcpUdpProfile{
									Type: "TCP_PROXY",
								},
							},
						},
						VirtualIpAddress:      "192.168.0.1",
						HealthStatus:          "UP",
						HealthMessage:         "OK",
						DetailedHealthMessage: "OK",
					},
				}, nil)
				getPoliciesHTTPRequest = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpRequestRule, error) {
					return nil, nil
				}
			},
			expectedPolicies: &PoliciesHTTPRequestModel{
				VirtualServiceID: virtualServiceID,
				Policies:         nil,
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name:             "error-refresh",
			expectedErr:      true,
			virtualServiceID: virtualServiceID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(errors.New("error"))
			},
			err: errors.New("error"),
		},
		{
			name:             "error-virtualserviceValidation",
			expectedErr:      true,
			virtualServiceID: "",
			mockFunc: func() {
			},
			err: errors.New("virtualServiceID is empty. Please provide a valid virtualServiceID"),
		},
		{
			name:             "error-getVirtualService",
			expectedErr:      true,
			virtualServiceID: virtualServiceID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbVirtualServiceById(virtualServiceID).Return(nil, errors.New("error"))
			},
			err: errors.New("error"),
		},
		{
			name:             "error-getPoliciesHTTPRequest",
			expectedErr:      true,
			virtualServiceID: virtualServiceID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbVirtualServiceById(virtualServiceID).Return(&govcd.NsxtAlbVirtualService{}, nil)
				getPoliciesHTTPRequest = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpRequestRule, error) {
					return nil, errors.New("error")
				}
			},
			err: errors.New("error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			policies, err := c.GetPoliciesHTTPRequest(t.Context(), tc.virtualServiceID)
			if !tc.expectedErr {
				assert.NoError(t, err)
				assert.NotNil(t, policies)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.err.Error())
				assert.Nil(t, policies)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedPolicies, policies)
		})
	}
}

func TestClient_UpdatePoliciesHTTPRequest(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock client for cloudavenue.
	clientCAV := NewMockclientFake(ctrl)

	c, _ := NewFakeClient(clientCAV)

	virtualServiceID := urn.LoadBalancerVirtualService.String() + uuid.New().String()
	edgeGatewayID := urn.Gateway.String() + uuid.New().String()
	poolID := urn.LoadBalancerPool.String() + uuid.New().String()
	serviceEngineID := urn.ServiceEngineGroup.String() + uuid.New().String()

	tests := []struct {
		name        string
		policies    *PoliciesHTTPRequestModel
		mockFunc    func()
		expectedErr bool
		err         error
	}{
		{
			name: "success",
			policies: &PoliciesHTTPRequestModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPRequestModelPolicy{
					{
						Name:    "ruleName2",
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPRequestMatchCriteria{
							Protocol:         "HTTP",
							ClientIPMatch:    &PoliciesHTTPClientIPMatch{Criteria: "IS_IN", Addresses: []string{"12.23.34.45", "12.23.34.0/24", "12.23.34.0-12.23.34.100"}},
							ServicePortMatch: &PoliciesHTTPServicePortMatch{Criteria: "IS_IN", Ports: []int{80, 443}},
							MethodMatch:      &PoliciesHTTPMethodMatch{Criteria: "IS_IN", Methods: []string{"GET", "POST"}},
							PathMatch:        &PoliciesHTTPPathMatch{Criteria: "BEGINS_WITH", MatchStrings: []string{"/path1", "/path2"}},
							QueryMatch:       []string{"key1=value1", "key2=value2"},
							HeaderMatch: []PoliciesHTTPHeaderMatch{
								{
									Criteria: "IS_IN",
									Name:     "User-Agent",
									Values:   []string{"Mozilla/5.0", "curl/7.64.1"},
								},
								{
									Criteria: "IS_IN",
									Name:     "Accept",
									Values:   []string{"application/json", "text/html"},
								},
							},
							CookieMatch: &PoliciesHTTPCookieMatch{Criteria: "BEGINS_WITH", Name: "session_id", Value: "abc123"},
						},
						RedirectAction: &PoliciesHTTPActionRedirect{
							Host:       "example.com",
							KeepQuery:  true,
							Path:       "/newpath",
							Port:       utils.ToPTR(80),
							Protocol:   "HTTP",
							StatusCode: 301,
						},
					},
				},
			},
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbVirtualServiceById(virtualServiceID).Return(&govcd.NsxtAlbVirtualService{
					NsxtAlbVirtualService: &govcdtypes.NsxtAlbVirtualService{
						ID:          virtualServiceID,
						Name:        "virtualServiceName2",
						Description: "virtualServiceDescription2",
						Enabled:     utils.ToPTR(true),
						ApplicationProfile: govcdtypes.NsxtAlbVirtualServiceApplicationProfile{
							Type: "HTTPS",
						},
						GatewayRef: govcdtypes.OpenApiReference{
							ID: edgeGatewayID,
						},
						LoadBalancerPoolRef: govcdtypes.OpenApiReference{
							ID: poolID,
						},
						ServiceEngineGroupRef: govcdtypes.OpenApiReference{
							ID: serviceEngineID,
						},
						ServicePorts: []govcdtypes.NsxtAlbVirtualServicePort{
							{
								PortStart:  utils.ToPTR(443),
								PortEnd:    nil,
								SslEnabled: utils.ToPTR(true),
								TcpUdpProfile: &govcdtypes.NsxtAlbVirtualServicePortTcpUdpProfile{
									Type: "TCP_PROXY",
								},
							},
						},
						VirtualIpAddress:      "192.168.0.1",
						HealthStatus:          "UP",
						HealthMessage:         "OK",
						DetailedHealthMessage: "OK",
					},
				}, nil)
				getPoliciesHTTPRequest = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpRequestRule, error) {
					return []*govcdtypes.AlbVsHttpRequestRule{
						{
							Name:    "ruleName",
							Active:  true,
							Logging: true,
							MatchCriteria: govcdtypes.AlbVsHttpRequestAndSecurityRuleMatchCriteria{
								ClientIPMatch: &govcdtypes.AlbVsHttpRequestRuleClientIPMatch{
									MatchCriteria: "IS_IN",
									Addresses: []string{
										"12.23.34.45",
										"12.23.34.0/24",
										"12.23.34.0-12.23.34.100",
									},
								},
								ServicePortMatch: &govcdtypes.AlbVsHttpRequestRuleServicePortMatch{
									MatchCriteria: "IS_IN",
									Ports: []int{
										80,
										443,
									},
								},
								MethodMatch: &govcdtypes.AlbVsHttpRequestRuleMethodMatch{
									MatchCriteria: "IS_IN",
									Methods: []string{
										"GET",
										"POST",
									},
								},
								Protocol: "HTTP",
								PathMatch: &govcdtypes.AlbVsHttpRequestRulePathMatch{
									MatchCriteria: "BEGINS_WITH",
									MatchStrings: []string{
										"/path1",
										"/path2",
									},
								},
								QueryMatch: []string{
									"key1=value1",
									"key2=value2",
								},
								HeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{
									{
										MatchCriteria: "IS_IN",
										Key:           "User-Agent",
										Value: []string{
											"Mozilla/5.0",
											"curl/7.64.1",
										},
									},
									{
										MatchCriteria: "IS_IN",
										Key:           "Accept",
										Value: []string{
											"application/json",
											"text/html",
										},
									},
								},
								CookieMatch: &govcdtypes.AlbVsHttpRequestRuleCookieMatch{
									MatchCriteria: "BEGINS_WITH",
									Key:           "session_id",
									Value:         "abc123",
								},
							},
							RedirectAction: &govcdtypes.AlbVsHttpRequestRuleRedirectAction{
								Host:       "example.com",
								KeepQuery:  true,
								Path:       "/newpath",
								Port:       utils.ToPTR(80),
								Protocol:   "HTTP",
								StatusCode: 301,
							},
						},
					}, nil
				}
				updatePoliciesHTTPRequest = func(_ fakeVirtualServiceClient, v *govcdtypes.AlbVsHttpRequestRules) (*govcdtypes.AlbVsHttpRequestRules, error) {
					return v, nil
				}
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name: "success-minimal",
			policies: &PoliciesHTTPRequestModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPRequestModelPolicy{
					{
						Name:    "ruleName2",
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPRequestMatchCriteria{
							Protocol:   "HTTP",
							QueryMatch: []string{"key1=value1", "key2=value2"},
						},
						RedirectAction: &PoliciesHTTPActionRedirect{
							Host:       "example.com",
							KeepQuery:  true,
							Path:       "/newpath",
							Port:       utils.ToPTR(80),
							Protocol:   "HTTP",
							StatusCode: 301,
						},
					},
				},
			},
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbVirtualServiceById(virtualServiceID).Return(&govcd.NsxtAlbVirtualService{
					NsxtAlbVirtualService: &govcdtypes.NsxtAlbVirtualService{
						ID:          virtualServiceID,
						Name:        "virtualServiceName2",
						Description: "virtualServiceDescription2",
						Enabled:     utils.ToPTR(true),
						ApplicationProfile: govcdtypes.NsxtAlbVirtualServiceApplicationProfile{
							Type: "HTTPS",
						},
						GatewayRef: govcdtypes.OpenApiReference{
							ID: edgeGatewayID,
						},
						LoadBalancerPoolRef: govcdtypes.OpenApiReference{
							ID: poolID,
						},
						ServiceEngineGroupRef: govcdtypes.OpenApiReference{
							ID: serviceEngineID,
						},
						ServicePorts: []govcdtypes.NsxtAlbVirtualServicePort{
							{
								PortStart:  utils.ToPTR(443),
								PortEnd:    nil,
								SslEnabled: utils.ToPTR(true),
								TcpUdpProfile: &govcdtypes.NsxtAlbVirtualServicePortTcpUdpProfile{
									Type: "TCP_PROXY",
								},
							},
						},
						VirtualIpAddress:      "192.168.0.1",
						HealthStatus:          "UP",
						HealthMessage:         "OK",
						DetailedHealthMessage: "OK",
					},
				}, nil)
				getPoliciesHTTPRequest = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpRequestRule, error) {
					return []*govcdtypes.AlbVsHttpRequestRule{
						{
							Name:    "ruleName",
							Active:  true,
							Logging: true,
							MatchCriteria: govcdtypes.AlbVsHttpRequestAndSecurityRuleMatchCriteria{
								ClientIPMatch: &govcdtypes.AlbVsHttpRequestRuleClientIPMatch{
									MatchCriteria: "IS_IN",
									Addresses: []string{
										"12.23.34.45",
										"12.23.34.0/24",
										"12.23.34.0-12.23.34.100",
									},
								},
								ServicePortMatch: &govcdtypes.AlbVsHttpRequestRuleServicePortMatch{
									MatchCriteria: "IS_IN",
									Ports: []int{
										80,
										443,
									},
								},
								MethodMatch: &govcdtypes.AlbVsHttpRequestRuleMethodMatch{
									MatchCriteria: "IS_IN",
									Methods: []string{
										"GET",
										"POST",
									},
								},
								Protocol: "HTTP",
								PathMatch: &govcdtypes.AlbVsHttpRequestRulePathMatch{
									MatchCriteria: "BEGINS_WITH",
									MatchStrings: []string{
										"/path1",
										"/path2",
									},
								},
								QueryMatch: []string{
									"key1=value1",
									"key2=value2",
								},
								HeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{
									{
										MatchCriteria: "IS_IN",
										Key:           "User-Agent",
										Value: []string{
											"Mozilla/5.0",
											"curl/7.64.1",
										},
									},
									{
										MatchCriteria: "IS_IN",
										Key:           "Accept",
										Value: []string{
											"application/json",
											"text/html",
										},
									},
								},
								CookieMatch: &govcdtypes.AlbVsHttpRequestRuleCookieMatch{
									MatchCriteria: "BEGINS_WITH",
									Key:           "session_id",
									Value:         "abc123",
								},
							},
							RedirectAction: &govcdtypes.AlbVsHttpRequestRuleRedirectAction{
								Host:       "example.com",
								KeepQuery:  true,
								Path:       "/newpath",
								Port:       utils.ToPTR(80),
								Protocol:   "HTTP",
								StatusCode: 301,
							},
						},
					}, nil
				}
				updatePoliciesHTTPRequest = func(_ fakeVirtualServiceClient, v *govcdtypes.AlbVsHttpRequestRules) (*govcdtypes.AlbVsHttpRequestRules, error) {
					return v, nil
				}
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name: "success-action-url-rewrite",
			policies: &PoliciesHTTPRequestModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPRequestModelPolicy{
					{
						Name:    "ruleName2",
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPRequestMatchCriteria{
							Protocol:   "HTTP",
							QueryMatch: []string{"key1=value1", "key2=value2"},
						},
						URLRewriteAction: &PoliciesHTTPActionURLRewrite{
							HostHeader: "example.com",
							Path:       "/newpath",
							KeepQuery:  true,
						},
					},
				},
			},
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbVirtualServiceById(virtualServiceID).Return(&govcd.NsxtAlbVirtualService{
					NsxtAlbVirtualService: &govcdtypes.NsxtAlbVirtualService{
						ID:          virtualServiceID,
						Name:        "virtualServiceName2",
						Description: "virtualServiceDescription2",
						Enabled:     utils.ToPTR(true),
						ApplicationProfile: govcdtypes.NsxtAlbVirtualServiceApplicationProfile{
							Type: "HTTPS",
						},
						GatewayRef: govcdtypes.OpenApiReference{
							ID: edgeGatewayID,
						},
						LoadBalancerPoolRef: govcdtypes.OpenApiReference{
							ID: poolID,
						},
						ServiceEngineGroupRef: govcdtypes.OpenApiReference{
							ID: serviceEngineID,
						},
						ServicePorts: []govcdtypes.NsxtAlbVirtualServicePort{
							{
								PortStart:  utils.ToPTR(443),
								PortEnd:    nil,
								SslEnabled: utils.ToPTR(true),
								TcpUdpProfile: &govcdtypes.NsxtAlbVirtualServicePortTcpUdpProfile{
									Type: "TCP_PROXY",
								},
							},
						},
						VirtualIpAddress:      "192.168.0.1",
						HealthStatus:          "UP",
						HealthMessage:         "OK",
						DetailedHealthMessage: "OK",
					},
				}, nil)
				getPoliciesHTTPRequest = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpRequestRule, error) {
					return []*govcdtypes.AlbVsHttpRequestRule{
						{
							Name:    "ruleName",
							Active:  true,
							Logging: true,
							MatchCriteria: govcdtypes.AlbVsHttpRequestAndSecurityRuleMatchCriteria{
								ClientIPMatch: &govcdtypes.AlbVsHttpRequestRuleClientIPMatch{
									MatchCriteria: "IS_IN",
									Addresses: []string{
										"12.23.34.45",
										"12.23.34.0/24",
										"12.23.34.0-12.23.34.100",
									},
								},
								ServicePortMatch: &govcdtypes.AlbVsHttpRequestRuleServicePortMatch{
									MatchCriteria: "IS_IN",
									Ports: []int{
										80,
										443,
									},
								},
								MethodMatch: &govcdtypes.AlbVsHttpRequestRuleMethodMatch{
									MatchCriteria: "IS_IN",
									Methods: []string{
										"GET",
										"POST",
									},
								},
								Protocol: "HTTP",
								PathMatch: &govcdtypes.AlbVsHttpRequestRulePathMatch{
									MatchCriteria: "BEGINS_WITH",
									MatchStrings: []string{
										"/path1",
										"/path2",
									},
								},
								QueryMatch: []string{
									"key1=value1",
									"key2=value2",
								},
								HeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{
									{
										MatchCriteria: "IS_IN",
										Key:           "User-Agent",
										Value: []string{
											"Mozilla/5.0",
											"curl/7.64.1",
										},
									},
									{
										MatchCriteria: "IS_IN",
										Key:           "Accept",
										Value: []string{
											"application/json",
											"text/html",
										},
									},
								},
								CookieMatch: &govcdtypes.AlbVsHttpRequestRuleCookieMatch{
									MatchCriteria: "BEGINS_WITH",
									Key:           "session_id",
									Value:         "abc123",
								},
							},
							RewriteURLAction: &govcdtypes.AlbVsHttpRequestRuleRewriteURLAction{
								Host:      "example.com",
								Path:      "/newpath",
								KeepQuery: true,
							},
						},
					}, nil
				}
				updatePoliciesHTTPRequest = func(_ fakeVirtualServiceClient, v *govcdtypes.AlbVsHttpRequestRules) (*govcdtypes.AlbVsHttpRequestRules, error) {
					return v, nil
				}
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name: "error-validation-model",
			policies: &PoliciesHTTPRequestModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPRequestModelPolicy{
					{
						Name:    "ruleName2",
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPRequestMatchCriteria{
							Protocol:         "HTTP",
							ClientIPMatch:    &PoliciesHTTPClientIPMatch{Criteria: "IS_IN", Addresses: []string{"12.23.34.45", "12.23.34.0/24", "12.23.34.0-12.23.34.100"}},
							ServicePortMatch: &PoliciesHTTPServicePortMatch{Criteria: "IS_IN", Ports: []int{80, 443}},
							MethodMatch:      &PoliciesHTTPMethodMatch{Criteria: "IS_IN", Methods: []string{"GET", "POST"}},
							PathMatch:        &PoliciesHTTPPathMatch{Criteria: "IS_IN", MatchStrings: []string{"/path1", "/path2"}},
							QueryMatch:       []string{"key1=value1", "key2=value2"},
							HeaderMatch: []PoliciesHTTPHeaderMatch{
								{
									Criteria: "IS_IN",
									Name:     "User-Agent",
									Values:   []string{"Mozilla/5.0", "curl/7.64.1"},
								},
								{
									Criteria: "IS_IN",
									Name:     "Accept",
									Values:   []string{"application/json", "text/html"},
								},
							},
							CookieMatch: &PoliciesHTTPCookieMatch{Criteria: "IS_IN", Name: "session_id", Value: "abc123"},
						},
						RedirectAction: &PoliciesHTTPActionRedirect{
							Host:       "example.com",
							KeepQuery:  true,
							Path:       "/newpath",
							Port:       utils.ToPTR(80),
							Protocol:   "HTTP",
							StatusCode: 301,
						},
					},
				},
			},
			mockFunc:    func() {},
			expectedErr: true,
			err:         errors.New("Error:Field validation"),
		},
		{
			name:        "error-refresh",
			expectedErr: true,
			policies: &PoliciesHTTPRequestModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPRequestModelPolicy{
					{
						Name:    "ruleName2",
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPRequestMatchCriteria{
							Protocol:         "HTTP",
							ClientIPMatch:    &PoliciesHTTPClientIPMatch{Criteria: "IS_IN", Addresses: []string{"12.23.34.45", "12.23.34.0/24", "12.23.34.0-12.23.34.100"}},
							ServicePortMatch: &PoliciesHTTPServicePortMatch{Criteria: "IS_IN", Ports: []int{80, 443}},
							MethodMatch:      &PoliciesHTTPMethodMatch{Criteria: "IS_IN", Methods: []string{"GET", "POST"}},
							PathMatch:        &PoliciesHTTPPathMatch{Criteria: "BEGINS_WITH", MatchStrings: []string{"/path1", "/path2"}},
							QueryMatch:       []string{"key1=value1", "key2=value2"},
							HeaderMatch: []PoliciesHTTPHeaderMatch{
								{
									Criteria: "IS_IN",
									Name:     "User-Agent",
									Values:   []string{"Mozilla/5.0", "curl/7.64.1"},
								},
								{
									Criteria: "IS_IN",
									Name:     "Accept",
									Values:   []string{"application/json", "text/html"},
								},
							},
							CookieMatch: &PoliciesHTTPCookieMatch{Criteria: "BEGINS_WITH", Name: "session_id", Value: "abc123"},
						},
						RedirectAction: &PoliciesHTTPActionRedirect{
							Host:       "example.com",
							KeepQuery:  true,
							Path:       "/newpath",
							Port:       utils.ToPTR(80),
							Protocol:   "HTTP",
							StatusCode: 301,
						},
					},
				},
			},
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(errors.New("error"))
			},
			err: errors.New("error"),
		},
		{
			name:        "error-getVirtualService",
			expectedErr: true,
			policies: &PoliciesHTTPRequestModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPRequestModelPolicy{
					{
						Name:    "ruleName2",
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPRequestMatchCriteria{
							Protocol:         "HTTP",
							ClientIPMatch:    &PoliciesHTTPClientIPMatch{Criteria: "IS_IN", Addresses: []string{"12.23.34.45", "12.23.34.0/24", "12.23.34.0-12.23.34.100"}},
							ServicePortMatch: &PoliciesHTTPServicePortMatch{Criteria: "IS_IN", Ports: []int{80, 443}},
							MethodMatch:      &PoliciesHTTPMethodMatch{Criteria: "IS_IN", Methods: []string{"GET", "POST"}},
							PathMatch:        &PoliciesHTTPPathMatch{Criteria: "BEGINS_WITH", MatchStrings: []string{"/path1", "/path2"}},
							QueryMatch:       []string{"key1=value1", "key2=value2"},
							HeaderMatch: []PoliciesHTTPHeaderMatch{
								{
									Criteria: "IS_IN",
									Name:     "User-Agent",
									Values:   []string{"Mozilla/5.0", "curl/7.64.1"},
								},
								{
									Criteria: "IS_IN",
									Name:     "Accept",
									Values:   []string{"application/json", "text/html"},
								},
							},
							CookieMatch: &PoliciesHTTPCookieMatch{Criteria: "BEGINS_WITH", Name: "session_id", Value: "abc123"},
						},
						RedirectAction: &PoliciesHTTPActionRedirect{
							Host:       "example.com",
							KeepQuery:  true,
							Path:       "/newpath",
							Port:       utils.ToPTR(80),
							Protocol:   "HTTP",
							StatusCode: 301,
						},
					},
				},
			},
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbVirtualServiceById(virtualServiceID).Return(nil, errors.New("error"))
			},
			err: errors.New("error"),
		},
		{
			name:        "error-updatePoliciesHTTPRequest",
			expectedErr: true,
			policies: &PoliciesHTTPRequestModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPRequestModelPolicy{
					{
						Name:    "ruleName2",
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPRequestMatchCriteria{
							Protocol:         "HTTP",
							ClientIPMatch:    &PoliciesHTTPClientIPMatch{Criteria: "IS_IN", Addresses: []string{"12.23.34.45", "12.23.34.0/24", "12.23.34.0-12.23.34.100"}},
							ServicePortMatch: &PoliciesHTTPServicePortMatch{Criteria: "IS_IN", Ports: []int{80, 443}},
							MethodMatch:      &PoliciesHTTPMethodMatch{Criteria: "IS_IN", Methods: []string{"GET", "POST"}},
							PathMatch:        &PoliciesHTTPPathMatch{Criteria: "BEGINS_WITH", MatchStrings: []string{"/path1", "/path2"}},
							QueryMatch:       []string{"key1=value1", "key2=value2"},
							HeaderMatch: []PoliciesHTTPHeaderMatch{
								{
									Criteria: "IS_IN",
									Name:     "User-Agent",
									Values:   []string{"Mozilla/5.0", "curl/7.64.1"},
								},
								{
									Criteria: "IS_IN",
									Name:     "Accept",
									Values:   []string{"application/json", "text/html"},
								},
							},
							CookieMatch: &PoliciesHTTPCookieMatch{Criteria: "BEGINS_WITH", Name: "session_id", Value: "abc123"},
						},
						RedirectAction: &PoliciesHTTPActionRedirect{
							Host:       "example.com",
							KeepQuery:  true,
							Path:       "/newpath",
							Port:       utils.ToPTR(80),
							Protocol:   "HTTP",
							StatusCode: 301,
						},
					},
				},
			},
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbVirtualServiceById(virtualServiceID).Return(&govcd.NsxtAlbVirtualService{}, nil)
				updatePoliciesHTTPRequest = func(_ fakeVirtualServiceClient, _ *govcdtypes.AlbVsHttpRequestRules) (*govcdtypes.AlbVsHttpRequestRules, error) {
					return nil, errors.New("error")
				}
			},
			err: errors.New("error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			updatedPolicies, err := c.UpdatePoliciesHTTPRequest(t.Context(), tc.policies)
			if !tc.expectedErr {
				assert.NoError(t, err)
				assert.NotNil(t, tc.policies)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.policies, updatedPolicies)
		})
	}
}

func TestClient_DeletePoliciesHTTPRequest(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock client for cloudavenue.
	clientCAV := NewMockclientFake(ctrl)

	c, _ := NewFakeClient(clientCAV)

	virtualServiceID := urn.LoadBalancerVirtualService.String() + uuid.New().String()
	edgeGatewayID := urn.Gateway.String() + uuid.New().String()
	poolID := urn.LoadBalancerPool.String() + uuid.New().String()
	serviceEngineID := urn.ServiceEngineGroup.String() + uuid.New().String()

	tests := []struct {
		name             string
		virtualServiceID string
		mockFunc         func()
		expectedErr      bool
		err              error
	}{
		{
			name:             "success",
			virtualServiceID: virtualServiceID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbVirtualServiceById(virtualServiceID).Return(&govcd.NsxtAlbVirtualService{
					NsxtAlbVirtualService: &govcdtypes.NsxtAlbVirtualService{
						ID:          virtualServiceID,
						Name:        "virtualServiceName2",
						Description: "virtualServiceDescription2",
						Enabled:     utils.ToPTR(true),
						ApplicationProfile: govcdtypes.NsxtAlbVirtualServiceApplicationProfile{
							Type: "HTTPS",
						},
						GatewayRef: govcdtypes.OpenApiReference{
							ID: edgeGatewayID,
						},
						LoadBalancerPoolRef: govcdtypes.OpenApiReference{
							ID: poolID,
						},
						ServiceEngineGroupRef: govcdtypes.OpenApiReference{
							ID: serviceEngineID,
						},
						ServicePorts: []govcdtypes.NsxtAlbVirtualServicePort{
							{
								PortStart:  utils.ToPTR(443),
								PortEnd:    nil,
								SslEnabled: utils.ToPTR(true),
								TcpUdpProfile: &govcdtypes.NsxtAlbVirtualServicePortTcpUdpProfile{
									Type: "TCP_PROXY",
								},
							},
						},
						VirtualIpAddress:      "192.168.0.1",
						HealthStatus:          "UP",
						HealthMessage:         "OK",
						DetailedHealthMessage: "OK",
					},
				}, nil)
				updatePoliciesHTTPRequest = func(_ fakeVirtualServiceClient, _ *govcdtypes.AlbVsHttpRequestRules) (*govcdtypes.AlbVsHttpRequestRules, error) {
					return &govcdtypes.AlbVsHttpRequestRules{}, nil
				}
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name:             "error-delete",
			virtualServiceID: virtualServiceID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbVirtualServiceById(virtualServiceID).Return(&govcd.NsxtAlbVirtualService{}, nil)
				updatePoliciesHTTPRequest = func(_ fakeVirtualServiceClient, _ *govcdtypes.AlbVsHttpRequestRules) (*govcdtypes.AlbVsHttpRequestRules, error) {
					return &govcdtypes.AlbVsHttpRequestRules{}, errors.New("error")
				}
			},
			expectedErr: true,
			err:         errors.New("error"),
		},
		{
			name:             "error-refresh",
			expectedErr:      true,
			virtualServiceID: virtualServiceID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(errors.New("error"))
			},
			err: errors.New("error"),
		},
		{
			name:             "error-virtualserviceValidation",
			expectedErr:      true,
			virtualServiceID: "",
			mockFunc: func() {
			},
			err: errors.New("virtualServiceID is empty. Please provide a valid virtualServiceID"),
		},
		{
			name:             "error-getVirtualService",
			expectedErr:      true,
			virtualServiceID: virtualServiceID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbVirtualServiceById(virtualServiceID).Return(nil, errors.New("error"))
			},
			err: errors.New("error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			err := c.DeletePoliciesHTTPRequest(t.Context(), tc.virtualServiceID)
			if !tc.expectedErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.err.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}
