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
				getPoliciesHTTPSecurity = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpSecurityRule, error) {
					return []*govcdtypes.AlbVsHttpSecurityRule{
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
						Name:    testRuleName,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPSecurityMatchCriteria{
							Protocol:         PoliciesHTTPProtocolHTTP,
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
				getPoliciesHTTPSecurity = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpSecurityRule, error) {
					return []*govcdtypes.AlbVsHttpSecurityRule{
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
							LocalResponseAction: &govcdtypes.AlbVsHttpSecurityRuleRateLimitLocalResponseAction{
								StatusCode:  404,
								ContentType: testContentTypeJSON,
								Content:     testJSONBody,
							},
						},
					}, nil
				}
			},
			expectedPolicies: &PoliciesHTTPSecurityModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPSecurityModelPolicy{
					{
						Name:    testRuleName,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPSecurityMatchCriteria{
							Protocol:         PoliciesHTTPProtocolHTTP,
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
						SendResponseAction: &PoliciesHTTPActionSendResponse{
							StatusCode:  404,
							ContentType: testContentTypeJSON,
							Content:     testJSONBody,
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
				getPoliciesHTTPSecurity = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpSecurityRule, error) {
					return []*govcdtypes.AlbVsHttpSecurityRule{
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
							RateLimitAction: &govcdtypes.AlbVsHttpSecurityRuleRateLimitAction{
								Count:  100,
								Period: 60,
								RedirectAction: &govcdtypes.AlbVsHttpRequestRuleRedirectAction{
									Host:       testDomain,
									KeepQuery:  true,
									Protocol:   string(PoliciesHTTPProtocolHTTPS),
									Port:       utils.ToPTR(443),
									StatusCode: 302,
									Path:       testNewPath,
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
						Name:    testRuleName,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPSecurityMatchCriteria{
							Protocol:         PoliciesHTTPProtocolHTTP,
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
						RateLimitAction: &PoliciesHTTPActionRateLimit{
							Count:  100,
							Period: 60,
							RedirectAction: &PoliciesHTTPActionRedirect{
								Host:       testDomain,
								KeepQuery:  true,
								Protocol:   string(PoliciesHTTPProtocolHTTPS),
								Port:       utils.ToPTR(443),
								StatusCode: 302,
								Path:       testNewPath,
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
				getPoliciesHTTPSecurity = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpSecurityRule, error) {
					return []*govcdtypes.AlbVsHttpSecurityRule{
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
							RateLimitAction: &govcdtypes.AlbVsHttpSecurityRuleRateLimitAction{
								Count:  100,
								Period: 60,
								LocalResponseAction: &govcdtypes.AlbVsHttpSecurityRuleRateLimitLocalResponseAction{
									StatusCode:  404,
									ContentType: testContentTypeJSON,
									Content:     testJSONBody,
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
						Name:    testRuleName,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPSecurityMatchCriteria{
							Protocol:         PoliciesHTTPProtocolHTTP,
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
						RateLimitAction: &PoliciesHTTPActionRateLimit{
							Count:  100,
							Period: 60,
							LocalResponseAction: &PoliciesHTTPActionSendResponse{
								StatusCode:  404,
								ContentType: testContentTypeJSON,
								Content:     testJSONBody,
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
		name             string
		policies         *PoliciesHTTPSecurityModel
		expectedPolicies *PoliciesHTTPSecurityModel
		mockFunc         func()
		expectedErr      bool
		err              error
	}{
		// ? ------------------------------------------------------------------------------
		// ? ------------------------------ SUCCESS CASES ---------------------------------
		{
			name: testSuccess,
			policies: &PoliciesHTTPSecurityModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPSecurityModelPolicy{
					{
						Name:    testRuleName2,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPSecurityMatchCriteria{
							Protocol:         PoliciesHTTPProtocolHTTP,
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
						RedirectToHTTPSAction: utils.ToPTR(8443),
					},
				},
			},
			expectedPolicies: &PoliciesHTTPSecurityModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPSecurityModelPolicy{
					{
						Name:    testRuleName2,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPSecurityMatchCriteria{
							Protocol:         PoliciesHTTPProtocolHTTP,
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
						RedirectToHTTPSAction: utils.ToPTR(8443),
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
				getPoliciesHTTPSecurity = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpSecurityRule, error) {
					return []*govcdtypes.AlbVsHttpSecurityRule{
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
		{
			name: "success_default",
			policies: &PoliciesHTTPSecurityModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPSecurityModelPolicy{
					{
						Name:    testRuleName3,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPSecurityMatchCriteria{
							Protocol:  PoliciesHTTPProtocolHTTP,
							PathMatch: &PoliciesHTTPPathMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), MatchStrings: []string{testPath1, testPath2}},
						},
						// RedirectToHTTPSAction: utils.ToPTR(8443),
						RateLimitAction: &PoliciesHTTPActionRateLimit{
							CloseConnectionAction: utils.ToPTR(true),
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
				getPoliciesHTTPSecurity = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpSecurityRule, error) {
					return []*govcdtypes.AlbVsHttpSecurityRule{
						{
							Name:    testRuleName3,
							Active:  true,
							Logging: true,
							MatchCriteria: govcdtypes.AlbVsHttpRequestAndSecurityRuleMatchCriteria{
								Protocol: string(PoliciesHTTPProtocolHTTP),
								PathMatch: &govcdtypes.AlbVsHttpRequestRulePathMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH),
									MatchStrings: []string{
										testPath1,
										testPath2,
									},
								},
							},
							RateLimitAction: &govcdtypes.AlbVsHttpSecurityRuleRateLimitAction{
								// Count:                 1000,
								// Period:                60,
								CloseConnectionAction: "CLOSE",
							},
						},
					}, nil
				}
				updatePoliciesHTTPSecurity = func(_ fakeVirtualServiceClient, v *govcdtypes.AlbVsHttpSecurityRules) (*govcdtypes.AlbVsHttpSecurityRules, error) {
					return v, nil
				}
			},
			expectedPolicies: &PoliciesHTTPSecurityModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPSecurityModelPolicy{
					{
						Name:    testRuleName3,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPSecurityMatchCriteria{
							Protocol:  PoliciesHTTPProtocolHTTP,
							PathMatch: &PoliciesHTTPPathMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), MatchStrings: []string{testPath1, testPath2}},
						},
						// RedirectToHTTPSAction: utils.ToPTR(8443),
						RateLimitAction: &PoliciesHTTPActionRateLimit{
							Count:                 1000,
							Period:                60,
							CloseConnectionAction: utils.ToPTR(true),
						},
					},
				},
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name: "success-rate-limit-action-redirect",
			policies: &PoliciesHTTPSecurityModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPSecurityModelPolicy{
					{
						Name:    testRuleName3,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPSecurityMatchCriteria{
							Protocol: PoliciesHTTPProtocolHTTP,
						},
						RateLimitAction: &PoliciesHTTPActionRateLimit{
							Count:  100,
							Period: 60,
							RedirectAction: &PoliciesHTTPActionRedirect{
								Host:       testDomain,
								KeepQuery:  true,
								Protocol:   string(PoliciesHTTPProtocolHTTPS),
								Port:       utils.ToPTR(443),
								StatusCode: 302,
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
				getPoliciesHTTPSecurity = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpSecurityRule, error) {
					return []*govcdtypes.AlbVsHttpSecurityRule{
						{
							Name:    testRuleName3,
							Active:  true,
							Logging: true,
							MatchCriteria: govcdtypes.AlbVsHttpRequestAndSecurityRuleMatchCriteria{
								Protocol: string(PoliciesHTTPProtocolHTTP),
							},
							RateLimitAction: &govcdtypes.AlbVsHttpSecurityRuleRateLimitAction{
								Count:  100,
								Period: 60,
								RedirectAction: &govcdtypes.AlbVsHttpRequestRuleRedirectAction{
									Host:       testDomain,
									KeepQuery:  true,
									Protocol:   string(PoliciesHTTPProtocolHTTPS),
									Port:       utils.ToPTR(443),
									StatusCode: 302,
								},
							},
						},
					}, nil
				}
				updatePoliciesHTTPSecurity = func(_ fakeVirtualServiceClient, v *govcdtypes.AlbVsHttpSecurityRules) (*govcdtypes.AlbVsHttpSecurityRules, error) {
					return v, nil
				}
			},
			expectedPolicies: &PoliciesHTTPSecurityModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPSecurityModelPolicy{
					{
						Name:    testRuleName3,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPSecurityMatchCriteria{
							Protocol: PoliciesHTTPProtocolHTTP,
						},
						RateLimitAction: &PoliciesHTTPActionRateLimit{
							Count:  100,
							Period: 60,
							RedirectAction: &PoliciesHTTPActionRedirect{
								Host:       testDomain,
								KeepQuery:  true,
								Protocol:   string(PoliciesHTTPProtocolHTTPS),
								Port:       utils.ToPTR(443),
								StatusCode: 302,
							},
						},
					},
				},
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name: "success-rate-limit-action-local-response",
			policies: &PoliciesHTTPSecurityModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPSecurityModelPolicy{
					{
						Name:    testRuleName3,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPSecurityMatchCriteria{
							Protocol: PoliciesHTTPProtocolHTTP,
						},
						RateLimitAction: &PoliciesHTTPActionRateLimit{
							Count:  100,
							Period: 60,
							LocalResponseAction: &PoliciesHTTPActionSendResponse{
								StatusCode:  404,
								ContentType: testContentTypeJSON,
								Content:     testBase64Body, // Note: JSON string : {"key":"value"}
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
				getPoliciesHTTPSecurity = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpSecurityRule, error) {
					return []*govcdtypes.AlbVsHttpSecurityRule{
						{
							Name:    testRuleName3,
							Active:  true,
							Logging: true,
							MatchCriteria: govcdtypes.AlbVsHttpRequestAndSecurityRuleMatchCriteria{
								Protocol: string(PoliciesHTTPProtocolHTTP),
							},
							RateLimitAction: &govcdtypes.AlbVsHttpSecurityRuleRateLimitAction{
								Count:  100,
								Period: 60,
								LocalResponseAction: &govcdtypes.AlbVsHttpSecurityRuleRateLimitLocalResponseAction{
									StatusCode:  404,
									ContentType: testContentTypeJSON,
									Content:     testBase64Body, // Note: JSON string : {"key":"value"}
								},
							},
						},
					}, nil
				}
				updatePoliciesHTTPSecurity = func(_ fakeVirtualServiceClient, v *govcdtypes.AlbVsHttpSecurityRules) (*govcdtypes.AlbVsHttpSecurityRules, error) {
					return v, nil
				}
			},
			expectedPolicies: &PoliciesHTTPSecurityModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPSecurityModelPolicy{
					{
						Name:    testRuleName3,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPSecurityMatchCriteria{
							Protocol: PoliciesHTTPProtocolHTTP,
						},
						RateLimitAction: &PoliciesHTTPActionRateLimit{
							Count:  100,
							Period: 60,
							LocalResponseAction: &PoliciesHTTPActionSendResponse{
								StatusCode:  404,
								ContentType: testContentTypeJSON,
								Content:     testBase64Body, // Note: JSON string : {"key":"value"}
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
			name:        testErrorRefresh,
			expectedErr: true,
			policies: &PoliciesHTTPSecurityModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPSecurityModelPolicy{
					{
						Name:    testRuleName3,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPSecurityMatchCriteria{
							Protocol:         PoliciesHTTPProtocolHTTP,
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
						RedirectToHTTPSAction: utils.ToPTR(8443),
					},
				},
			},
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(errors.New("error"))
			},
			err: errors.New("error"),
		},
		{
			name: testErrorValidationModel,
			policies: &PoliciesHTTPSecurityModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPSecurityModelPolicy{
					{
						Name:    testRuleName3,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPSecurityMatchCriteria{
							Protocol:         PoliciesHTTPProtocolHTTP,
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
						RedirectToHTTPSAction: utils.ToPTR(8443),
					},
				},
			},
			mockFunc:    func() {},
			expectedErr: true,
			err:         errors.New("Error:Field validation"),
		},
		{
			name:        testErrorGetVS,
			expectedErr: true,
			policies: &PoliciesHTTPSecurityModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPSecurityModelPolicy{
					{
						Name:    testRuleName2,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPSecurityMatchCriteria{
							Protocol:         PoliciesHTTPProtocolHTTP,
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
						RedirectToHTTPSAction: utils.ToPTR(8443),
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
			policies: &PoliciesHTTPSecurityModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPSecurityModelPolicy{
					{
						Name:    testRuleName2,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPSecurityMatchCriteria{
							Protocol:         PoliciesHTTPProtocolHTTP,
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
						RedirectToHTTPSAction: utils.ToPTR(8443),
					},
				},
			},
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbVirtualServiceById(virtualServiceID).Return(&govcd.NsxtAlbVirtualService{}, nil)
				updatePoliciesHTTPSecurity = func(_ fakeVirtualServiceClient, _ *govcdtypes.AlbVsHttpSecurityRules) (*govcdtypes.AlbVsHttpSecurityRules, error) {
					return nil, errors.New("error")
				}
			},
			err: errors.New("error"),
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
			assert.Equal(t, tc.expectedPolicies, updatedPolicies)
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
