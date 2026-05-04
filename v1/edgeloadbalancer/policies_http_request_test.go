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
			name:             testSuccess,
			virtualServiceID: virtualServiceID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbVirtualServiceById(virtualServiceID).Return(&govcd.NsxtAlbVirtualService{
					NsxtAlbVirtualService: &govcdtypes.NsxtAlbVirtualService{
						ID:          virtualServiceID,
						Name:        testVirtualServiceName2,
						Description: testVirtualServiceDesc2,
						Enabled:     utils.ToPTR(true),
						ApplicationProfile: govcdtypes.NsxtAlbVirtualServiceApplicationProfile{
							Type: string(PoliciesHTTPProtocolHTTPS),
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
									Type: string(virtualServiceServicePortTypeTCPProxy),
								},
							},
						},
						VirtualIpAddress:      testIPAddress,
						HealthStatus:          "UP",
						HealthMessage:         "OK",
						DetailedHealthMessage: "OK",
					},
				}, nil)
				getPoliciesHTTPRequest = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpRequestRule, error) {
					return []*govcdtypes.AlbVsHttpRequestRule{
						{
							Name:    testRuleName,
							Active:  true,
							Logging: true,
							MatchCriteria: govcdtypes.AlbVsHttpRequestAndSecurityRuleMatchCriteria{
								ClientIPMatch: &govcdtypes.AlbVsHttpRequestRuleClientIPMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Addresses: []string{
										testIPSingle,
										testIPCIDR,
										testIPRange,
									},
								},
								ServicePortMatch: &govcdtypes.AlbVsHttpRequestRuleServicePortMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Ports: []int{
										80,
										443,
									},
								},
								MethodMatch: &govcdtypes.AlbVsHttpRequestRuleMethodMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Methods: []string{
										string(PoliciesHTTPMethodGET),
										string(PoliciesHTTPMethodPOST),
									},
								},
								Protocol: string(PoliciesHTTPProtocolHTTP),
								PathMatch: &govcdtypes.AlbVsHttpRequestRulePathMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH),
									MatchStrings: []string{
										testPath1,
										testPath2,
									},
								},
								QueryMatch: []string{
									testQuery1,
									testQuery2,
								},
								HeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{
									{
										MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
										Key:           testHeaderUserAgent,
										Value: []string{
											testHeaderValueMozilla,
											testHeaderValueCurl,
										},
									},
									{
										MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
										Key:           testHeaderAccept,
										Value: []string{
											testContentTypeJSON,
											testContentTypeHTML,
										},
									},
								},
								CookieMatch: &govcdtypes.AlbVsHttpRequestRuleCookieMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH),
									Key:           testCookieName,
									Value:         testCookieValue,
								},
							},
							RedirectAction: &govcdtypes.AlbVsHttpRequestRuleRedirectAction{
								Host:       testDomain,
								KeepQuery:  true,
								Path:       testNewPath,
								Port:       utils.ToPTR(80),
								Protocol:   string(PoliciesHTTPProtocolHTTP),
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
						Name:    testRuleName,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPRequestMatchCriteria{
							Protocol:         string(PoliciesHTTPProtocolHTTP),
							ClientIPMatch:    &PoliciesHTTPClientIPMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Addresses: []string{testIPSingle, testIPCIDR, testIPRange}},
							ServicePortMatch: &PoliciesHTTPServicePortMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Ports: []int{80, 443}},
							MethodMatch:      &PoliciesHTTPMethodMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Methods: []string{string(PoliciesHTTPMethodGET), string(PoliciesHTTPMethodPOST)}},
							PathMatch:        &PoliciesHTTPPathMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), MatchStrings: []string{testPath1, testPath2}},
							QueryMatch:       []string{testQuery1, testQuery2},
							HeaderMatch: []PoliciesHTTPHeaderMatch{
								{
									Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Name:     testHeaderUserAgent,
									Values:   []string{testHeaderValueMozilla, testHeaderValueCurl},
								},
								{
									Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Name:     testHeaderAccept,
									Values:   []string{testContentTypeJSON, testContentTypeHTML},
								},
							},
							CookieMatch: &PoliciesHTTPCookieMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), Name: testCookieName, Value: testCookieValue},
						},
						RedirectAction: &PoliciesHTTPActionRedirect{
							Host:       testDomain,
							KeepQuery:  true,
							Path:       testNewPath,
							Port:       utils.ToPTR(80),
							Protocol:   string(PoliciesHTTPProtocolHTTP),
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
						Name:        testVirtualServiceName2,
						Description: testVirtualServiceDesc2,
						Enabled:     utils.ToPTR(true),
						ApplicationProfile: govcdtypes.NsxtAlbVirtualServiceApplicationProfile{
							Type: string(PoliciesHTTPProtocolHTTPS),
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
									Type: string(virtualServiceServicePortTypeTCPProxy),
								},
							},
						},
						VirtualIpAddress:      testIPAddress,
						HealthStatus:          "UP",
						HealthMessage:         "OK",
						DetailedHealthMessage: "OK",
					},
				}, nil)
				getPoliciesHTTPRequest = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpRequestRule, error) {
					return []*govcdtypes.AlbVsHttpRequestRule{
						{
							Name:    testRuleName,
							Active:  true,
							Logging: true,
							MatchCriteria: govcdtypes.AlbVsHttpRequestAndSecurityRuleMatchCriteria{
								Protocol:    string(PoliciesHTTPProtocolHTTP),
								HeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{},
							},
							RedirectAction: &govcdtypes.AlbVsHttpRequestRuleRedirectAction{
								Host:       testDomain,
								KeepQuery:  true,
								Path:       testNewPath,
								Port:       utils.ToPTR(80),
								Protocol:   string(PoliciesHTTPProtocolHTTP),
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
						Name:    testRuleName,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPRequestMatchCriteria{
							Protocol: string(PoliciesHTTPProtocolHTTP),
						},
						RedirectAction: &PoliciesHTTPActionRedirect{
							Host:       testDomain,
							KeepQuery:  true,
							Path:       testNewPath,
							Port:       utils.ToPTR(80),
							Protocol:   string(PoliciesHTTPProtocolHTTP),
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
						Name:        testVirtualServiceName2,
						Description: testVirtualServiceDesc2,
						Enabled:     utils.ToPTR(true),
						ApplicationProfile: govcdtypes.NsxtAlbVirtualServiceApplicationProfile{
							Type: string(PoliciesHTTPProtocolHTTPS),
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
									Type: string(virtualServiceServicePortTypeTCPProxy),
								},
							},
						},
						VirtualIpAddress:      testIPAddress,
						HealthStatus:          "UP",
						HealthMessage:         "OK",
						DetailedHealthMessage: "OK",
					},
				}, nil)
				getPoliciesHTTPRequest = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpRequestRule, error) {
					return []*govcdtypes.AlbVsHttpRequestRule{
						{
							Name:    testRuleName,
							Active:  true,
							Logging: true,
							MatchCriteria: govcdtypes.AlbVsHttpRequestAndSecurityRuleMatchCriteria{
								ClientIPMatch: &govcdtypes.AlbVsHttpRequestRuleClientIPMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Addresses: []string{
										testIPSingle,
										testIPCIDR,
										testIPRange,
									},
								},
								ServicePortMatch: &govcdtypes.AlbVsHttpRequestRuleServicePortMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Ports: []int{
										80,
										443,
									},
								},
								MethodMatch: &govcdtypes.AlbVsHttpRequestRuleMethodMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Methods: []string{
										string(PoliciesHTTPMethodGET),
										string(PoliciesHTTPMethodPOST),
									},
								},
								Protocol: string(PoliciesHTTPProtocolHTTP),
								PathMatch: &govcdtypes.AlbVsHttpRequestRulePathMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH),
									MatchStrings: []string{
										testPath1,
										testPath2,
									},
								},
								QueryMatch: []string{
									testQuery1,
									testQuery2,
								},
								HeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{
									{
										MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
										Key:           testHeaderUserAgent,
										Value: []string{
											testHeaderValueMozilla,
											testHeaderValueCurl,
										},
									},
									{
										MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
										Key:           testHeaderAccept,
										Value: []string{
											testContentTypeJSON,
											testContentTypeHTML,
										},
									},
								},
								CookieMatch: &govcdtypes.AlbVsHttpRequestRuleCookieMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH),
									Key:           testCookieName,
									Value:         testCookieValue,
								},
							},
							HeaderActions: []*govcdtypes.AlbVsHttpRequestRuleHeaderActions{
								{
									Action: string(PoliciesHTTPActionHeaderRewriteActionADD),
									Name:   testHeaderXForwardedFor,
									Value:  testHeaderValueTest,
								},
								{
									Action: string(PoliciesHTTPActionHeaderRewriteActionREMOVE),
									Name:   testHeaderXForwardedProto,
									Value:  "",
								},
							},
							RewriteURLAction: &govcdtypes.AlbVsHttpRequestRuleRewriteURLAction{
								Host:      testDomain,
								Path:      testNewPath,
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
						Name:    testRuleName,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPRequestMatchCriteria{
							Protocol:         string(PoliciesHTTPProtocolHTTP),
							ClientIPMatch:    &PoliciesHTTPClientIPMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Addresses: []string{testIPSingle, testIPCIDR, testIPRange}},
							ServicePortMatch: &PoliciesHTTPServicePortMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Ports: []int{80, 443}},
							MethodMatch:      &PoliciesHTTPMethodMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Methods: []string{string(PoliciesHTTPMethodGET), string(PoliciesHTTPMethodPOST)}},
							PathMatch:        &PoliciesHTTPPathMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), MatchStrings: []string{testPath1, testPath2}},
							QueryMatch:       []string{testQuery1, testQuery2},
							HeaderMatch: []PoliciesHTTPHeaderMatch{
								{
									Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Name:     testHeaderUserAgent,
									Values:   []string{testHeaderValueMozilla, testHeaderValueCurl},
								},
								{
									Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Name:     testHeaderAccept,
									Values:   []string{testContentTypeJSON, testContentTypeHTML},
								},
							},
							CookieMatch: &PoliciesHTTPCookieMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), Name: testCookieName, Value: testCookieValue},
						},
						HeaderRewriteActions: PoliciesHTTPActionHeadersRewrite{
							{
								Action: string(PoliciesHTTPActionHeaderRewriteActionADD),
								Name:   testHeaderXForwardedFor,
								Value:  testHeaderValueTest,
							},
							{
								Action: string(PoliciesHTTPActionHeaderRewriteActionREMOVE),
								Name:   testHeaderXForwardedProto,
								Value:  "",
							},
						},
						URLRewriteAction: &PoliciesHTTPActionURLRewrite{
							HostHeader: testDomain,
							Path:       testNewPath,
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
						Name:        testVirtualServiceName2,
						Description: testVirtualServiceDesc2,
						Enabled:     utils.ToPTR(true),
						ApplicationProfile: govcdtypes.NsxtAlbVirtualServiceApplicationProfile{
							Type: string(PoliciesHTTPProtocolHTTPS),
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
									Type: string(virtualServiceServicePortTypeTCPProxy),
								},
							},
						},
						VirtualIpAddress:      testIPAddress,
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
			name:             testErrorRefresh,
			expectedErr:      true,
			virtualServiceID: virtualServiceID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(errors.New("error"))
			},
			err: errors.New("error"),
		},
		{
			name:             testErrorVSValidation,
			expectedErr:      true,
			virtualServiceID: "",
			mockFunc: func() {
			},
			err: errors.New("virtualServiceID is empty. Please provide a valid virtualServiceID"),
		},
		{
			name:             testErrorGetVS,
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

			policies, err := c.GetPoliciesHTTPRequest(context.Background(), tc.virtualServiceID)
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
			name: testSuccess,
			policies: &PoliciesHTTPRequestModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPRequestModelPolicy{
					{
						Name:    testRuleName2,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPRequestMatchCriteria{
							Protocol:         string(PoliciesHTTPProtocolHTTP),
							ClientIPMatch:    &PoliciesHTTPClientIPMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Addresses: []string{testIPSingle, testIPCIDR, testIPRange}},
							ServicePortMatch: &PoliciesHTTPServicePortMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Ports: []int{80, 443}},
							MethodMatch:      &PoliciesHTTPMethodMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Methods: []string{string(PoliciesHTTPMethodGET), string(PoliciesHTTPMethodPOST)}},
							PathMatch:        &PoliciesHTTPPathMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), MatchStrings: []string{testPath1, testPath2}},
							QueryMatch:       []string{testQuery1, testQuery2},
							HeaderMatch: []PoliciesHTTPHeaderMatch{
								{
									Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Name:     testHeaderUserAgent,
									Values:   []string{testHeaderValueMozilla, testHeaderValueCurl},
								},
								{
									Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Name:     testHeaderAccept,
									Values:   []string{testContentTypeJSON, testContentTypeHTML},
								},
							},
							CookieMatch: &PoliciesHTTPCookieMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), Name: testCookieName, Value: testCookieValue},
						},
						RedirectAction: &PoliciesHTTPActionRedirect{
							Host:       testDomain,
							KeepQuery:  true,
							Path:       testNewPath,
							Port:       utils.ToPTR(80),
							Protocol:   string(PoliciesHTTPProtocolHTTP),
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
						Name:        testVirtualServiceName2,
						Description: testVirtualServiceDesc2,
						Enabled:     utils.ToPTR(true),
						ApplicationProfile: govcdtypes.NsxtAlbVirtualServiceApplicationProfile{
							Type: string(PoliciesHTTPProtocolHTTPS),
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
									Type: string(virtualServiceServicePortTypeTCPProxy),
								},
							},
						},
						VirtualIpAddress:      testIPAddress,
						HealthStatus:          "UP",
						HealthMessage:         "OK",
						DetailedHealthMessage: "OK",
					},
				}, nil)
				getPoliciesHTTPRequest = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpRequestRule, error) {
					return []*govcdtypes.AlbVsHttpRequestRule{
						{
							Name:    testRuleName,
							Active:  true,
							Logging: true,
							MatchCriteria: govcdtypes.AlbVsHttpRequestAndSecurityRuleMatchCriteria{
								ClientIPMatch: &govcdtypes.AlbVsHttpRequestRuleClientIPMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Addresses: []string{
										testIPSingle,
										testIPCIDR,
										testIPRange,
									},
								},
								ServicePortMatch: &govcdtypes.AlbVsHttpRequestRuleServicePortMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Ports: []int{
										80,
										443,
									},
								},
								MethodMatch: &govcdtypes.AlbVsHttpRequestRuleMethodMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Methods: []string{
										string(PoliciesHTTPMethodGET),
										string(PoliciesHTTPMethodPOST),
									},
								},
								Protocol: string(PoliciesHTTPProtocolHTTP),
								PathMatch: &govcdtypes.AlbVsHttpRequestRulePathMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH),
									MatchStrings: []string{
										testPath1,
										testPath2,
									},
								},
								QueryMatch: []string{
									testQuery1,
									testQuery2,
								},
								HeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{
									{
										MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
										Key:           testHeaderUserAgent,
										Value: []string{
											testHeaderValueMozilla,
											testHeaderValueCurl,
										},
									},
									{
										MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
										Key:           testHeaderAccept,
										Value: []string{
											testContentTypeJSON,
											testContentTypeHTML,
										},
									},
								},
								CookieMatch: &govcdtypes.AlbVsHttpRequestRuleCookieMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH),
									Key:           testCookieName,
									Value:         testCookieValue,
								},
							},
							RedirectAction: &govcdtypes.AlbVsHttpRequestRuleRedirectAction{
								Host:       testDomain,
								KeepQuery:  true,
								Path:       testNewPath,
								Port:       utils.ToPTR(80),
								Protocol:   string(PoliciesHTTPProtocolHTTP),
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
						Name:    testRuleName2,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPRequestMatchCriteria{
							Protocol:   string(PoliciesHTTPProtocolHTTP),
							QueryMatch: []string{testQuery1, testQuery2},
						},
						RedirectAction: &PoliciesHTTPActionRedirect{
							Host:       testDomain,
							KeepQuery:  true,
							Path:       testNewPath,
							Port:       utils.ToPTR(80),
							Protocol:   string(PoliciesHTTPProtocolHTTP),
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
						Name:        testVirtualServiceName2,
						Description: testVirtualServiceDesc2,
						Enabled:     utils.ToPTR(true),
						ApplicationProfile: govcdtypes.NsxtAlbVirtualServiceApplicationProfile{
							Type: string(PoliciesHTTPProtocolHTTPS),
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
									Type: string(virtualServiceServicePortTypeTCPProxy),
								},
							},
						},
						VirtualIpAddress:      testIPAddress,
						HealthStatus:          "UP",
						HealthMessage:         "OK",
						DetailedHealthMessage: "OK",
					},
				}, nil)
				getPoliciesHTTPRequest = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpRequestRule, error) {
					return []*govcdtypes.AlbVsHttpRequestRule{
						{
							Name:    testRuleName,
							Active:  true,
							Logging: true,
							MatchCriteria: govcdtypes.AlbVsHttpRequestAndSecurityRuleMatchCriteria{
								ClientIPMatch: &govcdtypes.AlbVsHttpRequestRuleClientIPMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Addresses: []string{
										testIPSingle,
										testIPCIDR,
										testIPRange,
									},
								},
								ServicePortMatch: &govcdtypes.AlbVsHttpRequestRuleServicePortMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Ports: []int{
										80,
										443,
									},
								},
								MethodMatch: &govcdtypes.AlbVsHttpRequestRuleMethodMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Methods: []string{
										string(PoliciesHTTPMethodGET),
										string(PoliciesHTTPMethodPOST),
									},
								},
								Protocol: string(PoliciesHTTPProtocolHTTP),
								PathMatch: &govcdtypes.AlbVsHttpRequestRulePathMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH),
									MatchStrings: []string{
										testPath1,
										testPath2,
									},
								},
								QueryMatch: []string{
									testQuery1,
									testQuery2,
								},
								HeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{
									{
										MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
										Key:           testHeaderUserAgent,
										Value: []string{
											testHeaderValueMozilla,
											testHeaderValueCurl,
										},
									},
									{
										MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
										Key:           testHeaderAccept,
										Value: []string{
											testContentTypeJSON,
											testContentTypeHTML,
										},
									},
								},
								CookieMatch: &govcdtypes.AlbVsHttpRequestRuleCookieMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH),
									Key:           testCookieName,
									Value:         testCookieValue,
								},
							},
							RedirectAction: &govcdtypes.AlbVsHttpRequestRuleRedirectAction{
								Host:       testDomain,
								KeepQuery:  true,
								Path:       testNewPath,
								Port:       utils.ToPTR(80),
								Protocol:   string(PoliciesHTTPProtocolHTTP),
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
						Name:    testRuleName2,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPRequestMatchCriteria{
							Protocol:   string(PoliciesHTTPProtocolHTTP),
							QueryMatch: []string{testQuery1, testQuery2},
						},
						URLRewriteAction: &PoliciesHTTPActionURLRewrite{
							HostHeader: testDomain,
							Path:       testNewPath,
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
						Name:        testVirtualServiceName2,
						Description: testVirtualServiceDesc2,
						Enabled:     utils.ToPTR(true),
						ApplicationProfile: govcdtypes.NsxtAlbVirtualServiceApplicationProfile{
							Type: string(PoliciesHTTPProtocolHTTPS),
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
									Type: string(virtualServiceServicePortTypeTCPProxy),
								},
							},
						},
						VirtualIpAddress:      testIPAddress,
						HealthStatus:          "UP",
						HealthMessage:         "OK",
						DetailedHealthMessage: "OK",
					},
				}, nil)
				getPoliciesHTTPRequest = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpRequestRule, error) {
					return []*govcdtypes.AlbVsHttpRequestRule{
						{
							Name:    testRuleName,
							Active:  true,
							Logging: true,
							MatchCriteria: govcdtypes.AlbVsHttpRequestAndSecurityRuleMatchCriteria{
								ClientIPMatch: &govcdtypes.AlbVsHttpRequestRuleClientIPMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Addresses: []string{
										testIPSingle,
										testIPCIDR,
										testIPRange,
									},
								},
								ServicePortMatch: &govcdtypes.AlbVsHttpRequestRuleServicePortMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Ports: []int{
										80,
										443,
									},
								},
								MethodMatch: &govcdtypes.AlbVsHttpRequestRuleMethodMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Methods: []string{
										string(PoliciesHTTPMethodGET),
										string(PoliciesHTTPMethodPOST),
									},
								},
								Protocol: string(PoliciesHTTPProtocolHTTP),
								PathMatch: &govcdtypes.AlbVsHttpRequestRulePathMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH),
									MatchStrings: []string{
										testPath1,
										testPath2,
									},
								},
								QueryMatch: []string{
									testQuery1,
									testQuery2,
								},
								HeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{
									{
										MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
										Key:           testHeaderUserAgent,
										Value: []string{
											testHeaderValueMozilla,
											testHeaderValueCurl,
										},
									},
									{
										MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
										Key:           testHeaderAccept,
										Value: []string{
											testContentTypeJSON,
											testContentTypeHTML,
										},
									},
								},
								CookieMatch: &govcdtypes.AlbVsHttpRequestRuleCookieMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH),
									Key:           testCookieName,
									Value:         testCookieValue,
								},
							},
							RewriteURLAction: &govcdtypes.AlbVsHttpRequestRuleRewriteURLAction{
								Host:      testDomain,
								Path:      testNewPath,
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
			name: testErrorValidationModel,
			policies: &PoliciesHTTPRequestModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPRequestModelPolicy{
					{
						Name:    testRuleName2,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPRequestMatchCriteria{
							Protocol:         string(PoliciesHTTPProtocolHTTP),
							ClientIPMatch:    &PoliciesHTTPClientIPMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Addresses: []string{testIPSingle, testIPCIDR, testIPRange}},
							ServicePortMatch: &PoliciesHTTPServicePortMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Ports: []int{80, 443}},
							MethodMatch:      &PoliciesHTTPMethodMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Methods: []string{string(PoliciesHTTPMethodGET), string(PoliciesHTTPMethodPOST)}},
							PathMatch:        &PoliciesHTTPPathMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), MatchStrings: []string{testPath1, testPath2}},
							QueryMatch:       []string{testQuery1, testQuery2},
							HeaderMatch: []PoliciesHTTPHeaderMatch{
								{
									Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Name:     testHeaderUserAgent,
									Values:   []string{testHeaderValueMozilla, testHeaderValueCurl},
								},
								{
									Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Name:     testHeaderAccept,
									Values:   []string{testContentTypeJSON, testContentTypeHTML},
								},
							},
							CookieMatch: &PoliciesHTTPCookieMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Name: testCookieName, Value: testCookieValue},
						},
						RedirectAction: &PoliciesHTTPActionRedirect{
							Host:       testDomain,
							KeepQuery:  true,
							Path:       testNewPath,
							Port:       utils.ToPTR(80),
							Protocol:   string(PoliciesHTTPProtocolHTTP),
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
			name:        testErrorRefresh,
			expectedErr: true,
			policies: &PoliciesHTTPRequestModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPRequestModelPolicy{
					{
						Name:    testRuleName2,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPRequestMatchCriteria{
							Protocol:         string(PoliciesHTTPProtocolHTTP),
							ClientIPMatch:    &PoliciesHTTPClientIPMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Addresses: []string{testIPSingle, testIPCIDR, testIPRange}},
							ServicePortMatch: &PoliciesHTTPServicePortMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Ports: []int{80, 443}},
							MethodMatch:      &PoliciesHTTPMethodMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Methods: []string{string(PoliciesHTTPMethodGET), string(PoliciesHTTPMethodPOST)}},
							PathMatch:        &PoliciesHTTPPathMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), MatchStrings: []string{testPath1, testPath2}},
							QueryMatch:       []string{testQuery1, testQuery2},
							HeaderMatch: []PoliciesHTTPHeaderMatch{
								{
									Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Name:     testHeaderUserAgent,
									Values:   []string{testHeaderValueMozilla, testHeaderValueCurl},
								},
								{
									Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Name:     testHeaderAccept,
									Values:   []string{testContentTypeJSON, testContentTypeHTML},
								},
							},
							CookieMatch: &PoliciesHTTPCookieMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), Name: testCookieName, Value: testCookieValue},
						},
						RedirectAction: &PoliciesHTTPActionRedirect{
							Host:       testDomain,
							KeepQuery:  true,
							Path:       testNewPath,
							Port:       utils.ToPTR(80),
							Protocol:   string(PoliciesHTTPProtocolHTTP),
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
			name:        testErrorGetVS,
			expectedErr: true,
			policies: &PoliciesHTTPRequestModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPRequestModelPolicy{
					{
						Name:    testRuleName2,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPRequestMatchCriteria{
							Protocol:         string(PoliciesHTTPProtocolHTTP),
							ClientIPMatch:    &PoliciesHTTPClientIPMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Addresses: []string{testIPSingle, testIPCIDR, testIPRange}},
							ServicePortMatch: &PoliciesHTTPServicePortMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Ports: []int{80, 443}},
							MethodMatch:      &PoliciesHTTPMethodMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Methods: []string{string(PoliciesHTTPMethodGET), string(PoliciesHTTPMethodPOST)}},
							PathMatch:        &PoliciesHTTPPathMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), MatchStrings: []string{testPath1, testPath2}},
							QueryMatch:       []string{testQuery1, testQuery2},
							HeaderMatch: []PoliciesHTTPHeaderMatch{
								{
									Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Name:     testHeaderUserAgent,
									Values:   []string{testHeaderValueMozilla, testHeaderValueCurl},
								},
								{
									Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Name:     testHeaderAccept,
									Values:   []string{testContentTypeJSON, testContentTypeHTML},
								},
							},
							CookieMatch: &PoliciesHTTPCookieMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), Name: testCookieName, Value: testCookieValue},
						},
						RedirectAction: &PoliciesHTTPActionRedirect{
							Host:       testDomain,
							KeepQuery:  true,
							Path:       testNewPath,
							Port:       utils.ToPTR(80),
							Protocol:   string(PoliciesHTTPProtocolHTTP),
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
						Name:    testRuleName2,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPRequestMatchCriteria{
							Protocol:         string(PoliciesHTTPProtocolHTTP),
							ClientIPMatch:    &PoliciesHTTPClientIPMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Addresses: []string{testIPSingle, testIPCIDR, testIPRange}},
							ServicePortMatch: &PoliciesHTTPServicePortMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Ports: []int{80, 443}},
							MethodMatch:      &PoliciesHTTPMethodMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Methods: []string{string(PoliciesHTTPMethodGET), string(PoliciesHTTPMethodPOST)}},
							PathMatch:        &PoliciesHTTPPathMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), MatchStrings: []string{testPath1, testPath2}},
							QueryMatch:       []string{testQuery1, testQuery2},
							HeaderMatch: []PoliciesHTTPHeaderMatch{
								{
									Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Name:     testHeaderUserAgent,
									Values:   []string{testHeaderValueMozilla, testHeaderValueCurl},
								},
								{
									Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Name:     testHeaderAccept,
									Values:   []string{testContentTypeJSON, testContentTypeHTML},
								},
							},
							CookieMatch: &PoliciesHTTPCookieMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), Name: testCookieName, Value: testCookieValue},
						},
						RedirectAction: &PoliciesHTTPActionRedirect{
							Host:       testDomain,
							KeepQuery:  true,
							Path:       testNewPath,
							Port:       utils.ToPTR(80),
							Protocol:   string(PoliciesHTTPProtocolHTTP),
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

			updatedPolicies, err := c.UpdatePoliciesHTTPRequest(context.Background(), tc.policies)
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
			name:             testSuccess,
			virtualServiceID: virtualServiceID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbVirtualServiceById(virtualServiceID).Return(&govcd.NsxtAlbVirtualService{
					NsxtAlbVirtualService: &govcdtypes.NsxtAlbVirtualService{
						ID:          virtualServiceID,
						Name:        testVirtualServiceName2,
						Description: testVirtualServiceDesc2,
						Enabled:     utils.ToPTR(true),
						ApplicationProfile: govcdtypes.NsxtAlbVirtualServiceApplicationProfile{
							Type: string(PoliciesHTTPProtocolHTTPS),
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
									Type: string(virtualServiceServicePortTypeTCPProxy),
								},
							},
						},
						VirtualIpAddress:      testIPAddress,
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
			name:             testErrorDelete,
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
			name:             testErrorRefresh,
			expectedErr:      true,
			virtualServiceID: virtualServiceID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(errors.New("error"))
			},
			err: errors.New("error"),
		},
		{
			name:             testErrorVSValidation,
			expectedErr:      true,
			virtualServiceID: "",
			mockFunc: func() {
			},
			err: errors.New("virtualServiceID is empty. Please provide a valid virtualServiceID"),
		},
		{
			name:             testErrorGetVS,
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

			err := c.DeletePoliciesHTTPRequest(context.Background(), tc.virtualServiceID)
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
