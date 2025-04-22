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
	"context"
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

func TestClient_GetPoliciesHTTPSecurity(t *testing.T) {
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
		expectedPolicies *PoliciesHTTPSecurityModel
		expectedErr      bool
		err              error
	}{
		// ? ------------------------------------------------------------------------------
		// ? ------------------------------ SUCCESS CASES ---------------------------------
		{
			name:             "success-action-redirect-to-https",
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
				getPoliciesHTTPSecurity = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpSecurityRule, error) {
					return []*govcdtypes.AlbVsHttpSecurityRule{
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
							RedirectToHTTPSAction: &govcdtypes.AlbVsHttpSecurityRuleRedirectToHTTPSAction{
								Port: 8443,
							},
						},
					}, nil
				}
			},
			expectedPolicies: &PoliciesHTTPSecurityModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPSecurityModelPolicy{
					{
						Name:    "ruleName",
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPSecurityMatchCriteria{
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
						RedirectToHTTPSAction: utils.ToPTR(8443),
					},
				},
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name:             "success-action-send-response",
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
				getPoliciesHTTPSecurity = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpSecurityRule, error) {
					return []*govcdtypes.AlbVsHttpSecurityRule{
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
							LocalResponseAction: &govcdtypes.AlbVsHttpSecurityRuleRateLimitLocalResponseAction{
								StatusCode:  404,
								ContentType: "application/json",
								Content:     "{\"key\":\"value\"}",
							},
						},
					}, nil
				}
			},
			expectedPolicies: &PoliciesHTTPSecurityModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPSecurityModelPolicy{
					{
						Name:    "ruleName",
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPSecurityMatchCriteria{
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
						SendResponseAction: &PoliciesHTTPActionSendResponse{
							StatusCode:  404,
							ContentType: "application/json",
							Content:     "{\"key\":\"value\"}",
						},
					},
				},
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name:             "success-action-rate-limit-action-redirect",
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
				getPoliciesHTTPSecurity = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpSecurityRule, error) {
					return []*govcdtypes.AlbVsHttpSecurityRule{
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
							RateLimitAction: &govcdtypes.AlbVsHttpSecurityRuleRateLimitAction{
								Count:  100,
								Period: 60,
								RedirectAction: &govcdtypes.AlbVsHttpRequestRuleRedirectAction{
									Host:       "example.com",
									KeepQuery:  true,
									Protocol:   "HTTPS",
									Port:       utils.ToPTR(443),
									StatusCode: 302,
									Path:       "/newpath",
								},
							},
						},
					}, nil
				}
			},
			expectedPolicies: &PoliciesHTTPSecurityModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPSecurityModelPolicy{
					{
						Name:    "ruleName",
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPSecurityMatchCriteria{
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
						RateLimitAction: &PoliciesHTTPActionRateLimit{
							Count:  100,
							Period: 60,
							RedirectAction: &PoliciesHTTPActionRedirect{
								Host:       "example.com",
								KeepQuery:  true,
								Protocol:   "HTTPS",
								Port:       utils.ToPTR(443),
								StatusCode: 302,
								Path:       "/newpath",
							},
						},
					},
				},
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name:             "success-action-rate-limit-action-local-response",
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
				getPoliciesHTTPSecurity = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpSecurityRule, error) {
					return []*govcdtypes.AlbVsHttpSecurityRule{
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
							RateLimitAction: &govcdtypes.AlbVsHttpSecurityRuleRateLimitAction{
								Count:  100,
								Period: 60,
								LocalResponseAction: &govcdtypes.AlbVsHttpSecurityRuleRateLimitLocalResponseAction{
									StatusCode:  404,
									ContentType: "application/json",
									Content:     "{\"key\":\"value\"}",
								},
							},
						},
					}, nil
				}
			},
			expectedPolicies: &PoliciesHTTPSecurityModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPSecurityModelPolicy{
					{
						Name:    "ruleName",
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPSecurityMatchCriteria{
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
						RateLimitAction: &PoliciesHTTPActionRateLimit{
							Count:  100,
							Period: 60,
							LocalResponseAction: &PoliciesHTTPActionSendResponse{
								StatusCode:  404,
								ContentType: "application/json",
								Content:     "{\"key\":\"value\"}",
							},
						},
					},
				},
			},
			expectedErr: false,
			err:         nil,
		},
		// ? ------------------------------------------------------------------------------
		// ? ------------------------------ ERROR CASES -----------------------------------
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
			name:             "error-getPoliciesHTTPSecurity",
			expectedErr:      true,
			virtualServiceID: virtualServiceID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbVirtualServiceById(virtualServiceID).Return(&govcd.NsxtAlbVirtualService{}, nil)
				getPoliciesHTTPSecurity = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpSecurityRule, error) {
					return nil, errors.New("error")
				}
			},
			err: errors.New("error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			policies, err := c.GetPoliciesHTTPSecurity(context.Background(), tc.virtualServiceID)
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

func TestClient_UpdatePoliciesHTTPSecurity(t *testing.T) {
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
		policies    *PoliciesHTTPSecurityModel
		mockFunc    func()
		expectedErr bool
		err         error
	}{
		// ? ------------------------------------------------------------------------------
		// ? ------------------------------ SUCCESS CASES ---------------------------------
		{
			name: "success",
			policies: &PoliciesHTTPSecurityModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPSecurityModelPolicy{
					{
						Name:    "ruleName2",
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPSecurityMatchCriteria{
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
						RedirectToHTTPSAction: utils.ToPTR(8443),
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
				getPoliciesHTTPSecurity = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpSecurityRule, error) {
					return []*govcdtypes.AlbVsHttpSecurityRule{
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
							RedirectToHTTPSAction: &govcdtypes.AlbVsHttpSecurityRuleRedirectToHTTPSAction{
								Port: 8443,
							},
						},
					}, nil
				}
				updatePoliciesHTTPSecurity = func(_ fakeVirtualServiceClient, v *govcdtypes.AlbVsHttpSecurityRules) (*govcdtypes.AlbVsHttpSecurityRules, error) {
					return v, nil
				}
			},
			expectedErr: false,
			err:         nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			updatedPolicies, err := c.UpdatePoliciesHTTPSecurity(context.Background(), tc.policies)
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

func TestClient_DeletePoliciesHTTPSecurity(t *testing.T) {
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
		// ? ------------------------------------------------------------------------------
		// ? ------------------------------ SUCCESS CASES ---------------------------------
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
				updatePoliciesHTTPSecurity = func(_ fakeVirtualServiceClient, _ *govcdtypes.AlbVsHttpSecurityRules) (*govcdtypes.AlbVsHttpSecurityRules, error) {
					return &govcdtypes.AlbVsHttpSecurityRules{}, nil
				}
			},
			expectedErr: false,
			err:         nil,
		},
		// ? ------------------------------------------------------------------------------
		// ? ------------------------------ ERROR CASES -----------------------------------
		{
			name:             "error-delete",
			virtualServiceID: virtualServiceID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbVirtualServiceById(virtualServiceID).Return(&govcd.NsxtAlbVirtualService{}, nil)
				updatePoliciesHTTPSecurity = func(_ fakeVirtualServiceClient, _ *govcdtypes.AlbVsHttpSecurityRules) (*govcdtypes.AlbVsHttpSecurityRules, error) {
					return &govcdtypes.AlbVsHttpSecurityRules{}, errors.New("error")
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

			err := c.DeletePoliciesHTTPSecurity(context.Background(), tc.virtualServiceID)
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
