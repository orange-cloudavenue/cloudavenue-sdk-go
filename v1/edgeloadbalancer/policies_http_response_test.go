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

func TestClient_GetPoliciesHTTPResponse(t *testing.T) {
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
		expectedPolicies *PoliciesHTTPResponseModel
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
				getPoliciesHTTPResponse = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpResponseRule, error) {
					return []*govcdtypes.AlbVsHttpResponseRule{
						{
							Name:    "ruleName",
							Active:  true,
							Logging: true,
							MatchCriteria: govcdtypes.AlbVsHttpResponseRuleMatchCriteria{
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
								RequestHeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{
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
								ResponseHeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{
									{
										MatchCriteria: "IS_IN",
										Key:           "X-Custom-Header",
										Value:         []string{"value1", "value2"},
									},
								},
								CookieMatch: &govcdtypes.AlbVsHttpRequestRuleCookieMatch{
									MatchCriteria: "BEGINS_WITH",
									Key:           "session_id",
									Value:         "abc123",
								},
								LocationHeaderMatch: &govcdtypes.AlbVsHttpResponseLocationHeaderMatch{
									MatchCriteria: "BEGINS_WITH",
									Value: []string{
										"example.com",
									},
								},
								StatusCodeMatch: &govcdtypes.AlbVsHttpRuleStatusCodeMatch{
									MatchCriteria: "IS_IN",
									StatusCodes: []string{
										"200",
										"301-303",
									},
								},
							},
							RewriteLocationHeaderAction: &govcdtypes.AlbVsHttpRespRuleRewriteLocationHeaderAction{
								Host:      "example.com",
								KeepQuery: true,
								Path:      "/newpath",
								Port:      utils.ToPTR(80),
								Protocol:  "HTTP",
							},
						},
					}, nil
				}
			},
			expectedPolicies: &PoliciesHTTPResponseModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPResponseModelPolicy{
					{
						Name:    "ruleName",
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPResponseMatchCriteria{
							Protocol:         "HTTP",
							ClientIPMatch:    &PoliciesHTTPClientIPMatch{Criteria: "IS_IN", Addresses: []string{"12.23.34.45", "12.23.34.0/24", "12.23.34.0-12.23.34.100"}},
							ServicePortMatch: &PoliciesHTTPServicePortMatch{Criteria: "IS_IN", Ports: []int{80, 443}},
							MethodMatch:      &PoliciesHTTPMethodMatch{Criteria: "IS_IN", Methods: []string{"GET", "POST"}},
							PathMatch:        &PoliciesHTTPPathMatch{Criteria: "BEGINS_WITH", MatchStrings: []string{"/path1", "/path2"}},
							QueryMatch:       []string{"key1=value1", "key2=value2"},
							CookieMatch:      &PoliciesHTTPCookieMatch{Criteria: "BEGINS_WITH", Name: "session_id", Value: "abc123"},
							LocationMatch:    &PoliciesHTTPLocationMatch{Criteria: "BEGINS_WITH", Values: []string{"example.com"}},
							RequestHeaderMatch: PoliciesHTTPHeadersMatch{
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
							ResponseHeaderMatch: PoliciesHTTPHeadersMatch{
								{
									Criteria: "IS_IN",
									Name:     "X-Custom-Header",
									Values:   []string{"value1", "value2"},
								},
							},
							StatusCodeMatch: &PoliciesHTTPStatusCodeMatch{Criteria: "IS_IN", StatusCodes: []string{"200", "301-303"}},
						},
						LocationRewriteAction: &PoliciesHTTPActionLocationRewrite{
							Host:      "example.com",
							KeepQuery: true,
							Path:      "/newpath",
							Port:      utils.ToPTR(80),
							Protocol:  "HTTP",
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
				getPoliciesHTTPResponse = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpResponseRule, error) {
					return []*govcdtypes.AlbVsHttpResponseRule{
						{
							Name:    "ruleName",
							Active:  true,
							Logging: true,
							MatchCriteria: govcdtypes.AlbVsHttpResponseRuleMatchCriteria{
								Protocol:           "HTTP",
								RequestHeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{},
							},
							RewriteLocationHeaderAction: &govcdtypes.AlbVsHttpRespRuleRewriteLocationHeaderAction{
								Host:      "example.com",
								KeepQuery: true,
								Path:      "/newpath",
								Port:      utils.ToPTR(80),
								Protocol:  "HTTP",
							},
						},
					}, nil
				}
			},
			expectedPolicies: &PoliciesHTTPResponseModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPResponseModelPolicy{
					{
						Name:    "ruleName",
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPResponseMatchCriteria{
							Protocol: "HTTP",
						},
						LocationRewriteAction: &PoliciesHTTPActionLocationRewrite{
							Host:      "example.com",
							KeepQuery: true,
							Path:      "/newpath",
							Port:      utils.ToPTR(80),
							Protocol:  "HTTP",
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
				getPoliciesHTTPResponse = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpResponseRule, error) {
					return []*govcdtypes.AlbVsHttpResponseRule{
						{
							Name:    "ruleName",
							Active:  true,
							Logging: true,
							MatchCriteria: govcdtypes.AlbVsHttpResponseRuleMatchCriteria{
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
								RequestHeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{
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
								ResponseHeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{
									{
										MatchCriteria: "IS_IN",
										Key:           "X-Custom-Header",
										Value:         []string{"value1", "value2"},
									},
								},
								CookieMatch: &govcdtypes.AlbVsHttpRequestRuleCookieMatch{
									MatchCriteria: "BEGINS_WITH",
									Key:           "session_id",
									Value:         "abc123",
								},
								LocationHeaderMatch: &govcdtypes.AlbVsHttpResponseLocationHeaderMatch{
									MatchCriteria: "BEGINS_WITH",
									Value: []string{
										"example.com",
									},
								},
								StatusCodeMatch: &govcdtypes.AlbVsHttpRuleStatusCodeMatch{
									MatchCriteria: "IS_IN",
									StatusCodes: []string{
										"200",
										"301-303",
									},
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
							RewriteLocationHeaderAction: &govcdtypes.AlbVsHttpRespRuleRewriteLocationHeaderAction{
								Protocol:  "HTTP",
								Host:      "example.com",
								Port:      utils.ToPTR(80),
								Path:      "/newpath",
								KeepQuery: true,
							},
						},
					}, nil
				}
			},
			expectedPolicies: &PoliciesHTTPResponseModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPResponseModelPolicy{
					{
						Name:    "ruleName",
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPResponseMatchCriteria{
							Protocol:         "HTTP",
							QueryMatch:       []string{"key1=value1", "key2=value2"},
							ClientIPMatch:    &PoliciesHTTPClientIPMatch{Criteria: "IS_IN", Addresses: []string{"12.23.34.45", "12.23.34.0/24", "12.23.34.0-12.23.34.100"}},
							ServicePortMatch: &PoliciesHTTPServicePortMatch{Criteria: "IS_IN", Ports: []int{80, 443}},
							MethodMatch:      &PoliciesHTTPMethodMatch{Criteria: "IS_IN", Methods: []string{"GET", "POST"}},
							PathMatch:        &PoliciesHTTPPathMatch{Criteria: "BEGINS_WITH", MatchStrings: []string{"/path1", "/path2"}},
							CookieMatch:      &PoliciesHTTPCookieMatch{Criteria: "BEGINS_WITH", Name: "session_id", Value: "abc123"},
							LocationMatch:    &PoliciesHTTPLocationMatch{Criteria: "BEGINS_WITH", Values: []string{"example.com"}},
							RequestHeaderMatch: PoliciesHTTPHeadersMatch{
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
							ResponseHeaderMatch: PoliciesHTTPHeadersMatch{
								{
									Criteria: "IS_IN",
									Name:     "X-Custom-Header",
									Values:   []string{"value1", "value2"},
								},
							},
							StatusCodeMatch: &PoliciesHTTPStatusCodeMatch{Criteria: "IS_IN", StatusCodes: []string{"200", "301-303"}},
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
						LocationRewriteAction: &PoliciesHTTPActionLocationRewrite{
							Protocol:  "HTTP",
							Host:      "example.com",
							Port:      utils.ToPTR(80),
							Path:      "/newpath",
							KeepQuery: true,
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
				getPoliciesHTTPResponse = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpResponseRule, error) {
					return nil, nil
				}
			},
			expectedPolicies: &PoliciesHTTPResponseModel{
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
			name:             "error-getPoliciesHTTPResponse",
			expectedErr:      true,
			virtualServiceID: virtualServiceID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbVirtualServiceById(virtualServiceID).Return(&govcd.NsxtAlbVirtualService{}, nil)
				getPoliciesHTTPResponse = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpResponseRule, error) {
					return nil, errors.New("error")
				}
			},
			err: errors.New("error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			policies, err := c.GetPoliciesHTTPResponse(context.Background(), tc.virtualServiceID)
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

func TestClient_UpdatePoliciesHTTPResponse(t *testing.T) {
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
		policies    *PoliciesHTTPResponseModel
		mockFunc    func()
		expectedErr bool
		err         error
	}{
		{
			name: "success",
			policies: &PoliciesHTTPResponseModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPResponseModelPolicy{
					{
						Name:    "ruleName2",
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPResponseMatchCriteria{
							Protocol:         "HTTP",
							ClientIPMatch:    &PoliciesHTTPClientIPMatch{Criteria: "IS_IN", Addresses: []string{"12.23.34.45", "12.23.34.0/24", "12.23.34.0-12.23.34.100"}},
							ServicePortMatch: &PoliciesHTTPServicePortMatch{Criteria: "IS_IN", Ports: []int{80, 443}},
							MethodMatch:      &PoliciesHTTPMethodMatch{Criteria: "IS_IN", Methods: []string{"GET", "POST"}},
							PathMatch:        &PoliciesHTTPPathMatch{Criteria: "BEGINS_WITH", MatchStrings: []string{"/path1", "/path2"}},
							QueryMatch:       []string{"key1=value1", "key2=value2"},
							CookieMatch:      &PoliciesHTTPCookieMatch{Criteria: "BEGINS_WITH", Name: "session_id", Value: "abc123"},
							LocationMatch:    &PoliciesHTTPLocationMatch{Criteria: "BEGINS_WITH", Values: []string{"example.com"}},
							RequestHeaderMatch: PoliciesHTTPHeadersMatch{
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
							ResponseHeaderMatch: PoliciesHTTPHeadersMatch{
								{
									Criteria: "IS_IN",
									Name:     "X-Custom-Header",
									Values:   []string{"value1", "value2"},
								},
							},
							StatusCodeMatch: &PoliciesHTTPStatusCodeMatch{Criteria: "IS_IN", StatusCodes: []string{"200", "301-303"}},
						},
						LocationRewriteAction: &PoliciesHTTPActionLocationRewrite{
							Host:      "example.com",
							KeepQuery: true,
							Path:      "/newpath",
							Port:      utils.ToPTR(80),
							Protocol:  "HTTP",
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
				getPoliciesHTTPResponse = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpResponseRule, error) {
					return []*govcdtypes.AlbVsHttpResponseRule{
						{
							Name:    "ruleName",
							Active:  true,
							Logging: true,
							MatchCriteria: govcdtypes.AlbVsHttpResponseRuleMatchCriteria{
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
								RequestHeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{
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
								ResponseHeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{
									{
										MatchCriteria: "IS_IN",
										Key:           "X-Custom-Header",
										Value:         []string{"value1", "value2"},
									},
								},
								CookieMatch: &govcdtypes.AlbVsHttpRequestRuleCookieMatch{
									MatchCriteria: "BEGINS_WITH",
									Key:           "session_id",
									Value:         "abc123",
								},
								LocationHeaderMatch: &govcdtypes.AlbVsHttpResponseLocationHeaderMatch{
									MatchCriteria: "BEGINS_WITH",
									Value: []string{
										"example.com",
									},
								},
								StatusCodeMatch: &govcdtypes.AlbVsHttpRuleStatusCodeMatch{
									MatchCriteria: "IS_IN",
									StatusCodes: []string{
										"200",
										"301-303",
									},
								},
							},
							RewriteLocationHeaderAction: &govcdtypes.AlbVsHttpRespRuleRewriteLocationHeaderAction{
								Host:      "example.com",
								KeepQuery: true,
								Path:      "/newpath",
								Port:      utils.ToPTR(80),
								Protocol:  "HTTP",
							},
						},
					}, nil
				}
				updatePoliciesHTTPResponse = func(_ fakeVirtualServiceClient, v *govcdtypes.AlbVsHttpResponseRules) (*govcdtypes.AlbVsHttpResponseRules, error) {
					return v, nil
				}
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name: "success-header-rewrite",
			policies: &PoliciesHTTPResponseModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPResponseModelPolicy{
					{
						Name:    "ruleName2",
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPResponseMatchCriteria{
							Protocol:         "HTTP",
							ClientIPMatch:    &PoliciesHTTPClientIPMatch{Criteria: "IS_IN", Addresses: []string{"12.23.34.45", "12.23.34.0/24", "12.23.34.0-12.23.34.100"}},
							ServicePortMatch: &PoliciesHTTPServicePortMatch{Criteria: "IS_IN", Ports: []int{80, 443}},
							MethodMatch:      &PoliciesHTTPMethodMatch{Criteria: "IS_IN", Methods: []string{"GET", "POST"}},
							PathMatch:        &PoliciesHTTPPathMatch{Criteria: "BEGINS_WITH", MatchStrings: []string{"/path1", "/path2"}},
							QueryMatch:       []string{"key1=value1", "key2=value2"},
							CookieMatch:      &PoliciesHTTPCookieMatch{Criteria: "BEGINS_WITH", Name: "session_id", Value: "abc123"},
							LocationMatch:    &PoliciesHTTPLocationMatch{Criteria: "BEGINS_WITH", Values: []string{"example.com"}},
							RequestHeaderMatch: PoliciesHTTPHeadersMatch{
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
							ResponseHeaderMatch: PoliciesHTTPHeadersMatch{
								{
									Criteria: "IS_IN",
									Name:     "X-Custom-Header",
									Values:   []string{"value1", "value2"},
								},
							},
							StatusCodeMatch: &PoliciesHTTPStatusCodeMatch{Criteria: "IS_IN", StatusCodes: []string{"200", "301-303"}},
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
							},
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
				getPoliciesHTTPResponse = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpResponseRule, error) {
					return []*govcdtypes.AlbVsHttpResponseRule{
						{
							Name:    "ruleName",
							Active:  true,
							Logging: true,
							MatchCriteria: govcdtypes.AlbVsHttpResponseRuleMatchCriteria{
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
								RequestHeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{
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
								ResponseHeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{
									{
										MatchCriteria: "IS_IN",
										Key:           "X-Custom-Header",
										Value:         []string{"value1", "value2"},
									},
								},
								CookieMatch: &govcdtypes.AlbVsHttpRequestRuleCookieMatch{
									MatchCriteria: "BEGINS_WITH",
									Key:           "session_id",
									Value:         "abc123",
								},
								LocationHeaderMatch: &govcdtypes.AlbVsHttpResponseLocationHeaderMatch{
									MatchCriteria: "BEGINS_WITH",
									Value: []string{
										"example.com",
									},
								},
								StatusCodeMatch: &govcdtypes.AlbVsHttpRuleStatusCodeMatch{
									MatchCriteria: "IS_IN",
									StatusCodes: []string{
										"200",
										"301-303",
									},
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
								},
							},
						},
					}, nil
				}
				updatePoliciesHTTPResponse = func(_ fakeVirtualServiceClient, v *govcdtypes.AlbVsHttpResponseRules) (*govcdtypes.AlbVsHttpResponseRules, error) {
					return v, nil
				}
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name: "success-minimal",
			policies: &PoliciesHTTPResponseModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPResponseModelPolicy{
					{
						Name:    "ruleName2",
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPResponseMatchCriteria{
							Protocol:   "HTTP",
							QueryMatch: []string{"key1=value1", "key2=value2"},
						},
						LocationRewriteAction: &PoliciesHTTPActionLocationRewrite{
							Host:      "example.com",
							KeepQuery: true,
							Path:      "/newpath",
							Port:      utils.ToPTR(80),
							Protocol:  "HTTP",
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
				getPoliciesHTTPResponse = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpResponseRule, error) {
					return []*govcdtypes.AlbVsHttpResponseRule{
						{
							Name:    "ruleName",
							Active:  true,
							Logging: true,
							MatchCriteria: govcdtypes.AlbVsHttpResponseRuleMatchCriteria{
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
								RequestHeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{
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
								ResponseHeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{
									{
										MatchCriteria: "IS_IN",
										Key:           "X-Custom-Header",
										Value:         []string{"value1", "value2"},
									},
								},
								CookieMatch: &govcdtypes.AlbVsHttpRequestRuleCookieMatch{
									MatchCriteria: "BEGINS_WITH",
									Key:           "session_id",
									Value:         "abc123",
								},
								LocationHeaderMatch: &govcdtypes.AlbVsHttpResponseLocationHeaderMatch{
									MatchCriteria: "BEGINS_WITH",
									Value: []string{
										"example.com",
									},
								},
								StatusCodeMatch: &govcdtypes.AlbVsHttpRuleStatusCodeMatch{
									MatchCriteria: "IS_IN",
									StatusCodes: []string{
										"200",
										"301-303",
									},
								},
							},
							RewriteLocationHeaderAction: &govcdtypes.AlbVsHttpRespRuleRewriteLocationHeaderAction{
								Host:      "example.com",
								KeepQuery: true,
								Path:      "/newpath",
								Port:      utils.ToPTR(80),
								Protocol:  "HTTP",
							},
						},
					}, nil
				}
				updatePoliciesHTTPResponse = func(_ fakeVirtualServiceClient, v *govcdtypes.AlbVsHttpResponseRules) (*govcdtypes.AlbVsHttpResponseRules, error) {
					return v, nil
				}
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name: "error-validation-model",
			policies: &PoliciesHTTPResponseModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPResponseModelPolicy{
					{
						Name:    "ruleName2",
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPResponseMatchCriteria{
							Protocol:         "HTTP",
							ClientIPMatch:    &PoliciesHTTPClientIPMatch{Criteria: "IS_IN", Addresses: []string{"12.23.34.45", "12.23.34.0/24", "12.23.34.0-12.23.34.100"}},
							ServicePortMatch: &PoliciesHTTPServicePortMatch{Criteria: "IS_IN", Ports: []int{80, 443}},
							MethodMatch:      &PoliciesHTTPMethodMatch{Criteria: "IS_IN", Methods: []string{"GET", "POST"}},
							PathMatch:        &PoliciesHTTPPathMatch{Criteria: "IS_IN", MatchStrings: []string{"/path1", "/path2"}},
							QueryMatch:       []string{"key1=value1", "key2=value2"},
							CookieMatch:      &PoliciesHTTPCookieMatch{Criteria: "IS_IN", Name: "session_id", Value: "abc123"},
						},
						LocationRewriteAction: &PoliciesHTTPActionLocationRewrite{
							Host:      "example.com",
							KeepQuery: true,
							Path:      "/newpath",
							Port:      utils.ToPTR(80),
							Protocol:  "HTTP",
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
			policies: &PoliciesHTTPResponseModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPResponseModelPolicy{
					{
						Name:    "ruleName2",
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPResponseMatchCriteria{
							Protocol:         "HTTP",
							ClientIPMatch:    &PoliciesHTTPClientIPMatch{Criteria: "IS_IN", Addresses: []string{"12.23.34.45", "12.23.34.0/24", "12.23.34.0-12.23.34.100"}},
							ServicePortMatch: &PoliciesHTTPServicePortMatch{Criteria: "IS_IN", Ports: []int{80, 443}},
							MethodMatch:      &PoliciesHTTPMethodMatch{Criteria: "IS_IN", Methods: []string{"GET", "POST"}},
							PathMatch:        &PoliciesHTTPPathMatch{Criteria: "BEGINS_WITH", MatchStrings: []string{"/path1", "/path2"}},
							QueryMatch:       []string{"key1=value1", "key2=value2"},
							CookieMatch:      &PoliciesHTTPCookieMatch{Criteria: "BEGINS_WITH", Name: "session_id", Value: "abc123"},
							LocationMatch:    &PoliciesHTTPLocationMatch{Criteria: "BEGINS_WITH", Values: []string{"example.com"}},
							RequestHeaderMatch: PoliciesHTTPHeadersMatch{
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
							ResponseHeaderMatch: PoliciesHTTPHeadersMatch{
								{
									Criteria: "IS_IN",
									Name:     "X-Custom-Header",
									Values:   []string{"value1", "value2"},
								},
							},
							StatusCodeMatch: &PoliciesHTTPStatusCodeMatch{Criteria: "IS_IN", StatusCodes: []string{"200", "301-303"}},
						},
						LocationRewriteAction: &PoliciesHTTPActionLocationRewrite{
							Host:      "example.com",
							KeepQuery: true,
							Path:      "/newpath",
							Port:      utils.ToPTR(80),
							Protocol:  "HTTP",
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
			policies: &PoliciesHTTPResponseModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPResponseModelPolicy{
					{
						Name:    "ruleName2",
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPResponseMatchCriteria{
							Protocol:         "HTTP",
							ClientIPMatch:    &PoliciesHTTPClientIPMatch{Criteria: "IS_IN", Addresses: []string{"12.23.34.45", "12.23.34.0/24", "12.23.34.0-12.23.34.100"}},
							ServicePortMatch: &PoliciesHTTPServicePortMatch{Criteria: "IS_IN", Ports: []int{80, 443}},
							MethodMatch:      &PoliciesHTTPMethodMatch{Criteria: "IS_IN", Methods: []string{"GET", "POST"}},
							PathMatch:        &PoliciesHTTPPathMatch{Criteria: "BEGINS_WITH", MatchStrings: []string{"/path1", "/path2"}},
							QueryMatch:       []string{"key1=value1", "key2=value2"},
							CookieMatch:      &PoliciesHTTPCookieMatch{Criteria: "BEGINS_WITH", Name: "session_id", Value: "abc123"},
							LocationMatch:    &PoliciesHTTPLocationMatch{Criteria: "BEGINS_WITH", Values: []string{"example.com"}},
							RequestHeaderMatch: PoliciesHTTPHeadersMatch{
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
							ResponseHeaderMatch: PoliciesHTTPHeadersMatch{
								{
									Criteria: "IS_IN",
									Name:     "X-Custom-Header",
									Values:   []string{"value1", "value2"},
								},
							},
							StatusCodeMatch: &PoliciesHTTPStatusCodeMatch{Criteria: "IS_IN", StatusCodes: []string{"200", "301-303"}},
						},
						LocationRewriteAction: &PoliciesHTTPActionLocationRewrite{
							Host:      "example.com",
							KeepQuery: true,
							Path:      "/newpath",
							Port:      utils.ToPTR(80),
							Protocol:  "HTTP",
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
			name:        "error-updatePoliciesHTTPResponse",
			expectedErr: true,
			policies: &PoliciesHTTPResponseModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPResponseModelPolicy{
					{
						Name:    "ruleName2",
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPResponseMatchCriteria{
							Protocol:         "HTTP",
							ClientIPMatch:    &PoliciesHTTPClientIPMatch{Criteria: "IS_IN", Addresses: []string{"12.23.34.45", "12.23.34.0/24", "12.23.34.0-12.23.34.100"}},
							ServicePortMatch: &PoliciesHTTPServicePortMatch{Criteria: "IS_IN", Ports: []int{80, 443}},
							MethodMatch:      &PoliciesHTTPMethodMatch{Criteria: "IS_IN", Methods: []string{"GET", "POST"}},
							PathMatch:        &PoliciesHTTPPathMatch{Criteria: "BEGINS_WITH", MatchStrings: []string{"/path1", "/path2"}},
							QueryMatch:       []string{"key1=value1", "key2=value2"},
							CookieMatch:      &PoliciesHTTPCookieMatch{Criteria: "BEGINS_WITH", Name: "session_id", Value: "abc123"},
							LocationMatch:    &PoliciesHTTPLocationMatch{Criteria: "BEGINS_WITH", Values: []string{"example.com"}},
							RequestHeaderMatch: PoliciesHTTPHeadersMatch{
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
							ResponseHeaderMatch: PoliciesHTTPHeadersMatch{
								{
									Criteria: "IS_IN",
									Name:     "X-Custom-Header",
									Values:   []string{"value1", "value2"},
								},
							},
							StatusCodeMatch: &PoliciesHTTPStatusCodeMatch{Criteria: "IS_IN", StatusCodes: []string{"200", "301-303"}},
						},
						LocationRewriteAction: &PoliciesHTTPActionLocationRewrite{
							Host:      "example.com",
							KeepQuery: true,
							Path:      "/newpath",
							Port:      utils.ToPTR(80),
							Protocol:  "HTTP",
						},
					},
				},
			},
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbVirtualServiceById(virtualServiceID).Return(&govcd.NsxtAlbVirtualService{}, nil)
				updatePoliciesHTTPResponse = func(_ fakeVirtualServiceClient, _ *govcdtypes.AlbVsHttpResponseRules) (*govcdtypes.AlbVsHttpResponseRules, error) {
					return nil, errors.New("error")
				}
			},
			err: errors.New("error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			updatedPolicies, err := c.UpdatePoliciesHTTPResponse(context.Background(), tc.policies)
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

func TestClient_DeletePoliciesHTTPResponse(t *testing.T) {
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
				updatePoliciesHTTPResponse = func(_ fakeVirtualServiceClient, _ *govcdtypes.AlbVsHttpResponseRules) (*govcdtypes.AlbVsHttpResponseRules, error) {
					return &govcdtypes.AlbVsHttpResponseRules{}, nil
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
				updatePoliciesHTTPResponse = func(_ fakeVirtualServiceClient, _ *govcdtypes.AlbVsHttpResponseRules) (*govcdtypes.AlbVsHttpResponseRules, error) {
					return &govcdtypes.AlbVsHttpResponseRules{}, errors.New("error")
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

			err := c.DeletePoliciesHTTPResponse(context.Background(), tc.virtualServiceID)
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
