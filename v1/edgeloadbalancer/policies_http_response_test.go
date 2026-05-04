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
				getPoliciesHTTPResponse = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpResponseRule, error) {
					return []*govcdtypes.AlbVsHttpResponseRule{
						{
							Name:    testRuleName,
							Active:  true,
							Logging: true,
							MatchCriteria: govcdtypes.AlbVsHttpResponseRuleMatchCriteria{
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
								RequestHeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{
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
								ResponseHeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{
									{
										MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
										Key:           testHeaderXCustom,
										Value:         []string{testHeaderValue1, testHeaderValue2},
									},
								},
								CookieMatch: &govcdtypes.AlbVsHttpRequestRuleCookieMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH),
									Key:           testCookieName,
									Value:         testCookieValue,
								},
								LocationHeaderMatch: &govcdtypes.AlbVsHttpResponseLocationHeaderMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH),
									Value: []string{
										testDomain,
									},
								},
								StatusCodeMatch: &govcdtypes.AlbVsHttpRuleStatusCodeMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									StatusCodes: []string{
										"200",
										testRedirectCode,
									},
								},
							},
							RewriteLocationHeaderAction: &govcdtypes.AlbVsHttpRespRuleRewriteLocationHeaderAction{
								Host:      testDomain,
								KeepQuery: true,
								Path:      testNewPath,
								Port:      utils.ToPTR(80),
								Protocol:  string(PoliciesHTTPProtocolHTTP),
							},
						},
					}, nil
				}
			},
			expectedPolicies: &PoliciesHTTPResponseModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPResponseModelPolicy{
					{
						Name:    testRuleName,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPResponseMatchCriteria{
							Protocol:         string(PoliciesHTTPProtocolHTTP),
							ClientIPMatch:    &PoliciesHTTPClientIPMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Addresses: []string{testIPSingle, testIPCIDR, testIPRange}},
							ServicePortMatch: &PoliciesHTTPServicePortMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Ports: []int{80, 443}},
							MethodMatch:      &PoliciesHTTPMethodMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Methods: []string{string(PoliciesHTTPMethodGET), string(PoliciesHTTPMethodPOST)}},
							PathMatch:        &PoliciesHTTPPathMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), MatchStrings: []string{testPath1, testPath2}},
							QueryMatch:       []string{testQuery1, testQuery2},
							CookieMatch:      &PoliciesHTTPCookieMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), Name: testCookieName, Value: testCookieValue},
							LocationMatch:    &PoliciesHTTPLocationMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), Values: []string{testDomain}},
							RequestHeaderMatch: PoliciesHTTPHeadersMatch{
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
							ResponseHeaderMatch: PoliciesHTTPHeadersMatch{
								{
									Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Name:     testHeaderXCustom,
									Values:   []string{testHeaderValue1, testHeaderValue2},
								},
							},
							StatusCodeMatch: &PoliciesHTTPStatusCodeMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), StatusCodes: []string{"200", testRedirectCode}},
						},
						LocationRewriteAction: &PoliciesHTTPActionLocationRewrite{
							Host:      testDomain,
							KeepQuery: true,
							Path:      testNewPath,
							Port:      utils.ToPTR(80),
							Protocol:  string(PoliciesHTTPProtocolHTTP),
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
				getPoliciesHTTPResponse = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpResponseRule, error) {
					return []*govcdtypes.AlbVsHttpResponseRule{
						{
							Name:    testRuleName,
							Active:  true,
							Logging: true,
							MatchCriteria: govcdtypes.AlbVsHttpResponseRuleMatchCriteria{
								Protocol:           string(PoliciesHTTPProtocolHTTP),
								RequestHeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{},
							},
							RewriteLocationHeaderAction: &govcdtypes.AlbVsHttpRespRuleRewriteLocationHeaderAction{
								Host:      testDomain,
								KeepQuery: true,
								Path:      testNewPath,
								Port:      utils.ToPTR(80),
								Protocol:  string(PoliciesHTTPProtocolHTTP),
							},
						},
					}, nil
				}
			},
			expectedPolicies: &PoliciesHTTPResponseModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPResponseModelPolicy{
					{
						Name:    testRuleName,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPResponseMatchCriteria{
							Protocol: string(PoliciesHTTPProtocolHTTP),
						},
						LocationRewriteAction: &PoliciesHTTPActionLocationRewrite{
							Host:      testDomain,
							KeepQuery: true,
							Path:      testNewPath,
							Port:      utils.ToPTR(80),
							Protocol:  string(PoliciesHTTPProtocolHTTP),
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
				getPoliciesHTTPResponse = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpResponseRule, error) {
					return []*govcdtypes.AlbVsHttpResponseRule{
						{
							Name:    testRuleName,
							Active:  true,
							Logging: true,
							MatchCriteria: govcdtypes.AlbVsHttpResponseRuleMatchCriteria{
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
								RequestHeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{
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
								ResponseHeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{
									{
										MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
										Key:           testHeaderXCustom,
										Value:         []string{testHeaderValue1, testHeaderValue2},
									},
								},
								CookieMatch: &govcdtypes.AlbVsHttpRequestRuleCookieMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH),
									Key:           testCookieName,
									Value:         testCookieValue,
								},
								LocationHeaderMatch: &govcdtypes.AlbVsHttpResponseLocationHeaderMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH),
									Value: []string{
										testDomain,
									},
								},
								StatusCodeMatch: &govcdtypes.AlbVsHttpRuleStatusCodeMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									StatusCodes: []string{
										"200",
										testRedirectCode,
									},
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
							RewriteLocationHeaderAction: &govcdtypes.AlbVsHttpRespRuleRewriteLocationHeaderAction{
								Protocol:  string(PoliciesHTTPProtocolHTTP),
								Host:      testDomain,
								Port:      utils.ToPTR(80),
								Path:      testNewPath,
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
						Name:    testRuleName,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPResponseMatchCriteria{
							Protocol:         string(PoliciesHTTPProtocolHTTP),
							QueryMatch:       []string{testQuery1, testQuery2},
							ClientIPMatch:    &PoliciesHTTPClientIPMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Addresses: []string{testIPSingle, testIPCIDR, testIPRange}},
							ServicePortMatch: &PoliciesHTTPServicePortMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Ports: []int{80, 443}},
							MethodMatch:      &PoliciesHTTPMethodMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Methods: []string{string(PoliciesHTTPMethodGET), string(PoliciesHTTPMethodPOST)}},
							PathMatch:        &PoliciesHTTPPathMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), MatchStrings: []string{testPath1, testPath2}},
							CookieMatch:      &PoliciesHTTPCookieMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), Name: testCookieName, Value: testCookieValue},
							LocationMatch:    &PoliciesHTTPLocationMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), Values: []string{testDomain}},
							RequestHeaderMatch: PoliciesHTTPHeadersMatch{
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
							ResponseHeaderMatch: PoliciesHTTPHeadersMatch{
								{
									Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Name:     testHeaderXCustom,
									Values:   []string{testHeaderValue1, testHeaderValue2},
								},
							},
							StatusCodeMatch: &PoliciesHTTPStatusCodeMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), StatusCodes: []string{"200", testRedirectCode}},
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
						LocationRewriteAction: &PoliciesHTTPActionLocationRewrite{
							Protocol:  string(PoliciesHTTPProtocolHTTP),
							Host:      testDomain,
							Port:      utils.ToPTR(80),
							Path:      testNewPath,
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
			name: testSuccess,
			policies: &PoliciesHTTPResponseModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPResponseModelPolicy{
					{
						Name:    testRuleName2,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPResponseMatchCriteria{
							Protocol:         string(PoliciesHTTPProtocolHTTP),
							ClientIPMatch:    &PoliciesHTTPClientIPMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Addresses: []string{testIPSingle, testIPCIDR, testIPRange}},
							ServicePortMatch: &PoliciesHTTPServicePortMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Ports: []int{80, 443}},
							MethodMatch:      &PoliciesHTTPMethodMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Methods: []string{string(PoliciesHTTPMethodGET), string(PoliciesHTTPMethodPOST)}},
							PathMatch:        &PoliciesHTTPPathMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), MatchStrings: []string{testPath1, testPath2}},
							QueryMatch:       []string{testQuery1, testQuery2},
							CookieMatch:      &PoliciesHTTPCookieMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), Name: testCookieName, Value: testCookieValue},
							LocationMatch:    &PoliciesHTTPLocationMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), Values: []string{testDomain}},
							RequestHeaderMatch: PoliciesHTTPHeadersMatch{
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
							ResponseHeaderMatch: PoliciesHTTPHeadersMatch{
								{
									Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Name:     testHeaderXCustom,
									Values:   []string{testHeaderValue1, testHeaderValue2},
								},
							},
							StatusCodeMatch: &PoliciesHTTPStatusCodeMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), StatusCodes: []string{"200", testRedirectCode}},
						},
						LocationRewriteAction: &PoliciesHTTPActionLocationRewrite{
							Host:      testDomain,
							KeepQuery: true,
							Path:      testNewPath,
							Port:      utils.ToPTR(80),
							Protocol:  string(PoliciesHTTPProtocolHTTP),
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
				getPoliciesHTTPResponse = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpResponseRule, error) {
					return []*govcdtypes.AlbVsHttpResponseRule{
						{
							Name:    testRuleName,
							Active:  true,
							Logging: true,
							MatchCriteria: govcdtypes.AlbVsHttpResponseRuleMatchCriteria{
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
								RequestHeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{
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
								ResponseHeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{
									{
										MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
										Key:           testHeaderXCustom,
										Value:         []string{testHeaderValue1, testHeaderValue2},
									},
								},
								CookieMatch: &govcdtypes.AlbVsHttpRequestRuleCookieMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH),
									Key:           testCookieName,
									Value:         testCookieValue,
								},
								LocationHeaderMatch: &govcdtypes.AlbVsHttpResponseLocationHeaderMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH),
									Value: []string{
										testDomain,
									},
								},
								StatusCodeMatch: &govcdtypes.AlbVsHttpRuleStatusCodeMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									StatusCodes: []string{
										"200",
										testRedirectCode,
									},
								},
							},
							RewriteLocationHeaderAction: &govcdtypes.AlbVsHttpRespRuleRewriteLocationHeaderAction{
								Host:      testDomain,
								KeepQuery: true,
								Path:      testNewPath,
								Port:      utils.ToPTR(80),
								Protocol:  string(PoliciesHTTPProtocolHTTP),
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
						Name:    testRuleName2,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPResponseMatchCriteria{
							Protocol:         string(PoliciesHTTPProtocolHTTP),
							ClientIPMatch:    &PoliciesHTTPClientIPMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Addresses: []string{testIPSingle, testIPCIDR, testIPRange}},
							ServicePortMatch: &PoliciesHTTPServicePortMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Ports: []int{80, 443}},
							MethodMatch:      &PoliciesHTTPMethodMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Methods: []string{string(PoliciesHTTPMethodGET), string(PoliciesHTTPMethodPOST)}},
							PathMatch:        &PoliciesHTTPPathMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), MatchStrings: []string{testPath1, testPath2}},
							QueryMatch:       []string{testQuery1, testQuery2},
							CookieMatch:      &PoliciesHTTPCookieMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), Name: testCookieName, Value: testCookieValue},
							LocationMatch:    &PoliciesHTTPLocationMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), Values: []string{testDomain}},
							RequestHeaderMatch: PoliciesHTTPHeadersMatch{
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
							ResponseHeaderMatch: PoliciesHTTPHeadersMatch{
								{
									Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Name:     testHeaderXCustom,
									Values:   []string{testHeaderValue1, testHeaderValue2},
								},
							},
							StatusCodeMatch: &PoliciesHTTPStatusCodeMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), StatusCodes: []string{"200", testRedirectCode}},
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
				getPoliciesHTTPResponse = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpResponseRule, error) {
					return []*govcdtypes.AlbVsHttpResponseRule{
						{
							Name:    testRuleName,
							Active:  true,
							Logging: true,
							MatchCriteria: govcdtypes.AlbVsHttpResponseRuleMatchCriteria{
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
								RequestHeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{
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
								ResponseHeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{
									{
										MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
										Key:           testHeaderXCustom,
										Value:         []string{testHeaderValue1, testHeaderValue2},
									},
								},
								CookieMatch: &govcdtypes.AlbVsHttpRequestRuleCookieMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH),
									Key:           testCookieName,
									Value:         testCookieValue,
								},
								LocationHeaderMatch: &govcdtypes.AlbVsHttpResponseLocationHeaderMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH),
									Value: []string{
										testDomain,
									},
								},
								StatusCodeMatch: &govcdtypes.AlbVsHttpRuleStatusCodeMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									StatusCodes: []string{
										"200",
										testRedirectCode,
									},
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
						Name:    testRuleName2,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPResponseMatchCriteria{
							Protocol:   string(PoliciesHTTPProtocolHTTP),
							QueryMatch: []string{testQuery1, testQuery2},
						},
						LocationRewriteAction: &PoliciesHTTPActionLocationRewrite{
							Host:      testDomain,
							KeepQuery: true,
							Path:      testNewPath,
							Port:      utils.ToPTR(80),
							Protocol:  string(PoliciesHTTPProtocolHTTP),
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
				getPoliciesHTTPResponse = func(_ fakeVirtualServiceClient) ([]*govcdtypes.AlbVsHttpResponseRule, error) {
					return []*govcdtypes.AlbVsHttpResponseRule{
						{
							Name:    testRuleName,
							Active:  true,
							Logging: true,
							MatchCriteria: govcdtypes.AlbVsHttpResponseRuleMatchCriteria{
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
								RequestHeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{
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
								ResponseHeaderMatch: []govcdtypes.AlbVsHttpRequestRuleHeaderMatch{
									{
										MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
										Key:           testHeaderXCustom,
										Value:         []string{testHeaderValue1, testHeaderValue2},
									},
								},
								CookieMatch: &govcdtypes.AlbVsHttpRequestRuleCookieMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH),
									Key:           testCookieName,
									Value:         testCookieValue,
								},
								LocationHeaderMatch: &govcdtypes.AlbVsHttpResponseLocationHeaderMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH),
									Value: []string{
										testDomain,
									},
								},
								StatusCodeMatch: &govcdtypes.AlbVsHttpRuleStatusCodeMatch{
									MatchCriteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									StatusCodes: []string{
										"200",
										testRedirectCode,
									},
								},
							},
							RewriteLocationHeaderAction: &govcdtypes.AlbVsHttpRespRuleRewriteLocationHeaderAction{
								Host:      testDomain,
								KeepQuery: true,
								Path:      testNewPath,
								Port:      utils.ToPTR(80),
								Protocol:  string(PoliciesHTTPProtocolHTTP),
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
			name: testErrorValidationModel,
			policies: &PoliciesHTTPResponseModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPResponseModelPolicy{
					{
						Name:    testRuleName2,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPResponseMatchCriteria{
							Protocol:         string(PoliciesHTTPProtocolHTTP),
							ClientIPMatch:    &PoliciesHTTPClientIPMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Addresses: []string{testIPSingle, testIPCIDR, testIPRange}},
							ServicePortMatch: &PoliciesHTTPServicePortMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Ports: []int{80, 443}},
							MethodMatch:      &PoliciesHTTPMethodMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Methods: []string{string(PoliciesHTTPMethodGET), string(PoliciesHTTPMethodPOST)}},
							PathMatch:        &PoliciesHTTPPathMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), MatchStrings: []string{testPath1, testPath2}},
							QueryMatch:       []string{testQuery1, testQuery2},
							CookieMatch:      &PoliciesHTTPCookieMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Name: testCookieName, Value: testCookieValue},
						},
						LocationRewriteAction: &PoliciesHTTPActionLocationRewrite{
							Host:      testDomain,
							KeepQuery: true,
							Path:      testNewPath,
							Port:      utils.ToPTR(80),
							Protocol:  string(PoliciesHTTPProtocolHTTP),
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
			policies: &PoliciesHTTPResponseModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPResponseModelPolicy{
					{
						Name:    testRuleName2,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPResponseMatchCriteria{
							Protocol:         string(PoliciesHTTPProtocolHTTP),
							ClientIPMatch:    &PoliciesHTTPClientIPMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Addresses: []string{testIPSingle, testIPCIDR, testIPRange}},
							ServicePortMatch: &PoliciesHTTPServicePortMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Ports: []int{80, 443}},
							MethodMatch:      &PoliciesHTTPMethodMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Methods: []string{string(PoliciesHTTPMethodGET), string(PoliciesHTTPMethodPOST)}},
							PathMatch:        &PoliciesHTTPPathMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), MatchStrings: []string{testPath1, testPath2}},
							QueryMatch:       []string{testQuery1, testQuery2},
							CookieMatch:      &PoliciesHTTPCookieMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), Name: testCookieName, Value: testCookieValue},
							LocationMatch:    &PoliciesHTTPLocationMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), Values: []string{testDomain}},
							RequestHeaderMatch: PoliciesHTTPHeadersMatch{
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
							ResponseHeaderMatch: PoliciesHTTPHeadersMatch{
								{
									Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Name:     testHeaderXCustom,
									Values:   []string{testHeaderValue1, testHeaderValue2},
								},
							},
							StatusCodeMatch: &PoliciesHTTPStatusCodeMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), StatusCodes: []string{"200", testRedirectCode}},
						},
						LocationRewriteAction: &PoliciesHTTPActionLocationRewrite{
							Host:      testDomain,
							KeepQuery: true,
							Path:      testNewPath,
							Port:      utils.ToPTR(80),
							Protocol:  string(PoliciesHTTPProtocolHTTP),
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
			policies: &PoliciesHTTPResponseModel{
				VirtualServiceID: virtualServiceID,
				Policies: []*PoliciesHTTPResponseModelPolicy{
					{
						Name:    testRuleName2,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPResponseMatchCriteria{
							Protocol:         string(PoliciesHTTPProtocolHTTP),
							ClientIPMatch:    &PoliciesHTTPClientIPMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Addresses: []string{testIPSingle, testIPCIDR, testIPRange}},
							ServicePortMatch: &PoliciesHTTPServicePortMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Ports: []int{80, 443}},
							MethodMatch:      &PoliciesHTTPMethodMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Methods: []string{string(PoliciesHTTPMethodGET), string(PoliciesHTTPMethodPOST)}},
							PathMatch:        &PoliciesHTTPPathMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), MatchStrings: []string{testPath1, testPath2}},
							QueryMatch:       []string{testQuery1, testQuery2},
							CookieMatch:      &PoliciesHTTPCookieMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), Name: testCookieName, Value: testCookieValue},
							LocationMatch:    &PoliciesHTTPLocationMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), Values: []string{testDomain}},
							RequestHeaderMatch: PoliciesHTTPHeadersMatch{
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
							ResponseHeaderMatch: PoliciesHTTPHeadersMatch{
								{
									Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Name:     testHeaderXCustom,
									Values:   []string{testHeaderValue1, testHeaderValue2},
								},
							},
							StatusCodeMatch: &PoliciesHTTPStatusCodeMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), StatusCodes: []string{"200", testRedirectCode}},
						},
						LocationRewriteAction: &PoliciesHTTPActionLocationRewrite{
							Host:      testDomain,
							KeepQuery: true,
							Path:      testNewPath,
							Port:      utils.ToPTR(80),
							Protocol:  string(PoliciesHTTPProtocolHTTP),
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
						Name:    testRuleName2,
						Active:  true,
						Logging: true,
						MatchCriteria: PoliciesHTTPResponseMatchCriteria{
							Protocol:         string(PoliciesHTTPProtocolHTTP),
							ClientIPMatch:    &PoliciesHTTPClientIPMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Addresses: []string{testIPSingle, testIPCIDR, testIPRange}},
							ServicePortMatch: &PoliciesHTTPServicePortMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Ports: []int{80, 443}},
							MethodMatch:      &PoliciesHTTPMethodMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), Methods: []string{string(PoliciesHTTPMethodGET), string(PoliciesHTTPMethodPOST)}},
							PathMatch:        &PoliciesHTTPPathMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), MatchStrings: []string{testPath1, testPath2}},
							QueryMatch:       []string{testQuery1, testQuery2},
							CookieMatch:      &PoliciesHTTPCookieMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), Name: testCookieName, Value: testCookieValue},
							LocationMatch:    &PoliciesHTTPLocationMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaBEGINSWITH), Values: []string{testDomain}},
							RequestHeaderMatch: PoliciesHTTPHeadersMatch{
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
							ResponseHeaderMatch: PoliciesHTTPHeadersMatch{
								{
									Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN),
									Name:     testHeaderXCustom,
									Values:   []string{testHeaderValue1, testHeaderValue2},
								},
							},
							StatusCodeMatch: &PoliciesHTTPStatusCodeMatch{Criteria: string(PoliciesHTTPMatchCriteriaCriteriaISIN), StatusCodes: []string{"200", testRedirectCode}},
						},
						LocationRewriteAction: &PoliciesHTTPActionLocationRewrite{
							Host:      testDomain,
							KeepQuery: true,
							Path:      testNewPath,
							Port:      utils.ToPTR(80),
							Protocol:  string(PoliciesHTTPProtocolHTTP),
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
				updatePoliciesHTTPResponse = func(_ fakeVirtualServiceClient, _ *govcdtypes.AlbVsHttpResponseRules) (*govcdtypes.AlbVsHttpResponseRules, error) {
					return &govcdtypes.AlbVsHttpResponseRules{}, nil
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
				updatePoliciesHTTPResponse = func(_ fakeVirtualServiceClient, _ *govcdtypes.AlbVsHttpResponseRules) (*govcdtypes.AlbVsHttpResponseRules, error) {
					return &govcdtypes.AlbVsHttpResponseRules{}, errors.New("error")
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
