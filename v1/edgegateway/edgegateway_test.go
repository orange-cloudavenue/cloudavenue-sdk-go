/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgegateway

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/internal/endpoints"
	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

func TestClient_GetEdgeGateway(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	defer httpmock.DeactivateAndReset()

	// Mock client for cloudavenue.
	clientCAV := NewMockclientInterface(ctrl)

	c, _ := NewFakeClient(clientCAV)

	edgeGatewayID := urn.Gateway.String() + uuid.New().String()
	vdcID := urn.VDC.String() + uuid.New().String()

	tests := []struct {
		name                string
		edgeGatewayNameOrID string
		mockFunc            func()
		expectedEdgeGateway *EdgeGatewayModel
		expectedError       bool
		err                 error
	}{
		{
			name:                "success",
			edgeGatewayNameOrID: edgeGatewayID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetNsxtEdgeGatewayById(edgeGatewayID).Return(&govcd.NsxtEdgeGateway{
					EdgeGateway: &govcdtypes.OpenAPIEdgeGateway{
						ID:          edgeGatewayID,
						Name:        "test-edge-gateway",
						Description: "test description",
						Status:      "ACTIVE",
						OwnerRef: &govcdtypes.OpenApiReference{
							ID:   vdcID,
							Name: "test-vdc",
						},
						EdgeGatewayUplinks: []govcdtypes.EdgeGatewayUplinks{
							{
								UplinkName: "prvrf01eocb0009999allsp01",
							},
						},
					},
				}, nil)
				// mock getBandwidth
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())
					responder, err := httpmock.NewJsonResponder(200, json.RawMessage(`{"rateLimit":10}`))
					if err != nil {
						t.Fatal(err)
					}

					httpmock.RegisterResponder("GET", endpoints.InlineTemplate(endpoints.EdgeGatewayGet, map[string]string{"edge-id": urn.ExtractUUID(edgeGatewayID)}), responder)
					return clientcloudavenue.MockClient().R()
				})
				// mock getNetworkServices
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())
					responder, err := httpmock.NewJsonResponder(200, json.RawMessage(fmt.Sprintf(`
[
   {
      "type": "tier-0-vrf",
      "name": "prvrf01eocb0009999allsp01",
      "displayName": "",
      "properties": {
         "classOfService": "SHARED_STANDARD"
      },
      "children": [
         {
            "type": "edge-gateway",
            "name": "test-edge-gateway",
            "displayName": "",
            "properties": {
               "rateLimit": 5,
               "edgeUUID": "%s"
            },
            "children": [
               {
                  "type": "load-balancer",
                  "name": "737b9768-95a0-4955-bbbe-d5eab846e8dc",
                  "displayName": "v999w99eprnxcdshrdsegp99",
                  "properties": {
                     "classOfService": "PREMIUM",
                     "maxVirtualServices": 10
                  },
                  "children": []
               },
               {
                  "type": "service",
                  "name": "internet",
                  "displayName": "internet",
                  "properties": {
                     "ip": "12.123.123.12",
                     "announced": true
                  },
                  "children": [],
                  "serviceId": "ip-12-123-123-12"
               },
               {
                  "type": "service",
                  "name": "cav-services",
                  "displayName": "Cloud Avenue Services",
                  "properties": {
                     "ranges": [
                        "100.113.99.96/27"
                     ],
                     "ipCount": 16
                  },
                  "children": [],
                  "serviceId": "tn99e99ocb0009999spt199-cav-services"
               }
            ]
         }
      ]
   }
]
`, urn.ExtractUUID(edgeGatewayID))))
					if err != nil {
						t.Fatal(err)
					}

					httpmock.RegisterResponder("GET", endpoints.NetworkServiceGet, responder)
					return clientcloudavenue.MockClient().R()
				})
			},
			expectedEdgeGateway: &EdgeGatewayModel{
				ID:          edgeGatewayID,
				Name:        "test-edge-gateway",
				Description: "test description",
				OwnerRef: &govcdtypes.OpenApiReference{
					ID:   vdcID,
					Name: "test-vdc",
				},
				Status:    "ACTIVE",
				UplinkT0:  "prvrf01eocb0009999allsp01",
				Bandwidth: 10,
				Services: NetworkServicesModelSvcs{
					LoadBalancer: &NetworkServicesModelSvcLoadBalancer{
						NetworkServicesModelSvc: NetworkServicesModelSvc{
							ID:   "737b9768-95a0-4955-bbbe-d5eab846e8dc",
							Name: "v999w99eprnxcdshrdsegp99",
						},
						ClassOfService:     "PREMIUM",
						MaxVirtualServices: 10,
					},
					PublicIP: []*NetworkServicesModelSvcPublicIP{
						{
							NetworkServicesModelSvc: NetworkServicesModelSvc{
								ID:   "ip-12-123-123-12",
								Name: "12.123.123.12",
							},
							IP:        "12.123.123.12",
							Announced: true,
						},
					},
					Service: &NetworkServicesModelSvcService{
						NetworkServicesModelSvc: NetworkServicesModelSvc{
							ID:   "tn99e99ocb0009999spt199-cav-services",
							Name: "Cloud Avenue Services",
						},
						Network:               "100.113.99.96/27",
						DedicatedIPForService: "100.113.99.96",
						ServiceDetails:        ListOfServices,
					},
				},
			},
			expectedError: false,
		},
		{
			name:                "success-by-name",
			edgeGatewayNameOrID: "test-edge-gateway",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetNsxtEdgeGatewayByName("test-edge-gateway").Return(&govcd.NsxtEdgeGateway{
					EdgeGateway: &govcdtypes.OpenAPIEdgeGateway{
						ID:          edgeGatewayID,
						Name:        "test-edge-gateway",
						Description: "test description",
						Status:      "ACTIVE",
						OwnerRef: &govcdtypes.OpenApiReference{
							ID:   vdcID,
							Name: "test-vdc",
						},
						EdgeGatewayUplinks: []govcdtypes.EdgeGatewayUplinks{
							{
								UplinkName: "prvrf01eocb0009999allsp01",
							},
						},
					},
				}, nil)
				// mock getBandwidth
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())
					responder, err := httpmock.NewJsonResponder(200, json.RawMessage(`{"rateLimit":10}`))
					if err != nil {
						t.Fatal(err)
					}

					httpmock.RegisterResponder("GET", endpoints.InlineTemplate(endpoints.EdgeGatewayGet, map[string]string{"edge-id": urn.ExtractUUID(edgeGatewayID)}), responder)
					return clientcloudavenue.MockClient().R()
				})
				// mock getNetworkServices
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())
					responder, err := httpmock.NewJsonResponder(200, json.RawMessage(fmt.Sprintf(`
[
   {
      "type": "tier-0-vrf",
      "name": "prvrf01eocb0009999allsp01",
      "displayName": "",
      "properties": {
         "classOfService": "SHARED_STANDARD"
      },
      "children": [
         {
            "type": "edge-gateway",
            "name": "test-edge-gateway",
            "displayName": "",
            "properties": {
               "rateLimit": 5,
               "edgeUUID": "%s"
            },
            "children": [
               {
                  "type": "load-balancer",
                  "name": "737b9768-95a0-4955-bbbe-d5eab846e8dc",
                  "displayName": "v999w99eprnxcdshrdsegp99",
                  "properties": {
                     "classOfService": "PREMIUM",
                     "maxVirtualServices": 10
                  },
                  "children": []
               },
               {
                  "type": "service",
                  "name": "internet",
                  "displayName": "internet",
                  "properties": {
                     "ip": "12.123.123.12",
                     "announced": true
                  },
                  "children": [],
                  "serviceId": "ip-12-123-123-12"
               },
               {
                  "type": "service",
                  "name": "cav-services",
                  "displayName": "Cloud Avenue Services",
                  "properties": {
                     "ranges": [
                        "100.113.99.96/27"
                     ],
                     "ipCount": 16
                  },
                  "children": [],
                  "serviceId": "tn99e99ocb0009999spt199-cav-services"
               }
            ]
         }
      ]
   }
]
`, urn.ExtractUUID(edgeGatewayID))))
					if err != nil {
						t.Fatal(err)
					}

					httpmock.RegisterResponder("GET", endpoints.NetworkServiceGet, responder)
					return clientcloudavenue.MockClient().R()
				})
			},
			expectedEdgeGateway: &EdgeGatewayModel{
				ID:          edgeGatewayID,
				Name:        "test-edge-gateway",
				Description: "test description",
				OwnerRef: &govcdtypes.OpenApiReference{
					ID:   vdcID,
					Name: "test-vdc",
				},
				Status:    "ACTIVE",
				UplinkT0:  "prvrf01eocb0009999allsp01",
				Bandwidth: 10,
				Services: NetworkServicesModelSvcs{
					LoadBalancer: &NetworkServicesModelSvcLoadBalancer{
						NetworkServicesModelSvc: NetworkServicesModelSvc{
							ID:   "737b9768-95a0-4955-bbbe-d5eab846e8dc",
							Name: "v999w99eprnxcdshrdsegp99",
						},
						ClassOfService:     "PREMIUM",
						MaxVirtualServices: 10,
					},
					PublicIP: []*NetworkServicesModelSvcPublicIP{
						{
							NetworkServicesModelSvc: NetworkServicesModelSvc{
								ID:   "ip-12-123-123-12",
								Name: "12.123.123.12",
							},
							IP:        "12.123.123.12",
							Announced: true,
						},
					},
					Service: &NetworkServicesModelSvcService{
						NetworkServicesModelSvc: NetworkServicesModelSvc{
							ID:   "tn99e99ocb0009999spt199-cav-services",
							Name: "Cloud Avenue Services",
						},
						Network:               "100.113.99.96/27",
						DedicatedIPForService: "100.113.99.96",
						ServiceDetails:        ListOfServices,
					},
				},
			},
			expectedError: false,
		},
		{
			name:                "error-get",
			edgeGatewayNameOrID: edgeGatewayID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetNsxtEdgeGatewayById(edgeGatewayID).Return(nil, fmt.Errorf("error"))
			},
			expectedEdgeGateway: nil,
			expectedError:       true,
			err:                 fmt.Errorf("error retrieving edge gateway %s: %w", edgeGatewayID, fmt.Errorf("error")),
		},
		{
			name:                "refresh-error",
			edgeGatewayNameOrID: edgeGatewayID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(fmt.Errorf("refresh error"))
			},
			expectedEdgeGateway: nil,
			expectedError:       true,
			err:                 fmt.Errorf("refresh error"),
		},
		{
			name:                "empty-edge-gateway-name-or-id",
			edgeGatewayNameOrID: "",
			expectedEdgeGateway: nil,
			expectedError:       true,
			err:                 fmt.Errorf("edgeGatewayNameOrID is %w. Please provide a valid edgeGatewayNameOrID", errors.ErrEmpty),
			mockFunc:            func() {},
		},
		{
			name:                "error-get-services",
			edgeGatewayNameOrID: edgeGatewayID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetNsxtEdgeGatewayById(edgeGatewayID).Return(&govcd.NsxtEdgeGateway{
					EdgeGateway: &govcdtypes.OpenAPIEdgeGateway{
						ID:          edgeGatewayID,
						Name:        "test-edge-gateway",
						Description: "test description",
						Status:      "ACTIVE",
						OwnerRef: &govcdtypes.OpenApiReference{
							ID:   vdcID,
							Name: "test-vdc",
						},
					},
				}, nil)
				// mock getBandwidth
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())
					responder, err := httpmock.NewJsonResponder(200, json.RawMessage(`{"rateLimit":10}`))
					if err != nil {
						t.Fatal(err)
					}

					httpmock.RegisterResponder("GET", endpoints.InlineTemplate(endpoints.EdgeGatewayGet, map[string]string{"edge-id": urn.ExtractUUID(edgeGatewayID)}), responder)
					return clientcloudavenue.MockClient().R()
				})
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())
					responder := httpmock.NewErrorResponder(fmt.Errorf("error"))
					httpmock.RegisterResponder("GET", endpoints.NetworkServiceGet, responder)
					return clientcloudavenue.MockClient().R()
				})
			},
			expectedEdgeGateway: nil,
			expectedError:       true,
			err:                 fmt.Errorf("error getting edge gateway network services"),
		},
		{
			name:                "error-500-get-services",
			edgeGatewayNameOrID: edgeGatewayID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetNsxtEdgeGatewayById(edgeGatewayID).Return(&govcd.NsxtEdgeGateway{
					EdgeGateway: &govcdtypes.OpenAPIEdgeGateway{
						ID:          edgeGatewayID,
						Name:        "test-edge-gateway",
						Description: "test description",
						Status:      "ACTIVE",
						OwnerRef: &govcdtypes.OpenApiReference{
							ID:   vdcID,
							Name: "test-vdc",
						},
					},
				}, nil)
				// mock getBandwidth
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())
					responder, err := httpmock.NewJsonResponder(200, json.RawMessage(`{"rateLimit":10}`))
					if err != nil {
						t.Fatal(err)
					}

					httpmock.RegisterResponder("GET", endpoints.InlineTemplate(endpoints.EdgeGatewayGet, map[string]string{"edge-id": urn.ExtractUUID(edgeGatewayID)}), responder)
					return clientcloudavenue.MockClient().R()
				})
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())
					responder, err := httpmock.NewJsonResponder(500, json.RawMessage(`{"error":"error"}`))
					if err != nil {
						t.Fatal(err)
					}

					httpmock.RegisterResponder("GET", endpoints.NetworkServiceGet, responder)
					return clientcloudavenue.MockClient().R()
				})
			},
			expectedEdgeGateway: nil,
			expectedError:       true,
			err:                 fmt.Errorf("error getting edge gateway network services"),
		},
		{
			name:                "error-get-bandwidth",
			edgeGatewayNameOrID: edgeGatewayID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetNsxtEdgeGatewayById(edgeGatewayID).Return(&govcd.NsxtEdgeGateway{
					EdgeGateway: &govcdtypes.OpenAPIEdgeGateway{
						ID:          edgeGatewayID,
						Name:        "test-edge-gateway",
						Description: "test description",
						Status:      "ACTIVE",
						OwnerRef: &govcdtypes.OpenApiReference{
							ID:   vdcID,
							Name: "test-vdc",
						},
					},
				}, nil)
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())
					responder := httpmock.NewErrorResponder(fmt.Errorf("error"))
					httpmock.RegisterResponder("GET", endpoints.InlineTemplate(endpoints.EdgeGatewayGet, map[string]string{"edge-id": urn.ExtractUUID(edgeGatewayID)}), responder)
					return clientcloudavenue.MockClient().R()
				})
			},
			expectedEdgeGateway: nil,
			expectedError:       true,
			err:                 fmt.Errorf("error getting edge gateway"),
		},
		{
			name:                "error-get-bandwidth-500",
			edgeGatewayNameOrID: edgeGatewayID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetNsxtEdgeGatewayById(edgeGatewayID).Return(&govcd.NsxtEdgeGateway{
					EdgeGateway: &govcdtypes.OpenAPIEdgeGateway{
						ID:          edgeGatewayID,
						Name:        "test-edge-gateway",
						Description: "test description",
						Status:      "ACTIVE",
						OwnerRef: &govcdtypes.OpenApiReference{
							ID:   vdcID,
							Name: "test-vdc",
						},
					},
				}, nil)
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())
					responder, err := httpmock.NewJsonResponder(500, json.RawMessage(`{"error":"error"}`))
					if err != nil {
						t.Fatal(err)
					}

					httpmock.RegisterResponder("GET", endpoints.InlineTemplate(endpoints.EdgeGatewayGet, map[string]string{"edge-id": urn.ExtractUUID(edgeGatewayID)}), responder)
					return clientcloudavenue.MockClient().R()
				})
			},
			expectedEdgeGateway: nil,
			expectedError:       true,
			err:                 fmt.Errorf("error getting edge gateway"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockFunc()

			edgeGateway, err := c.GetEdgeGateway(context.Background(), test.edgeGatewayNameOrID)
			if !test.expectedError {
				assert.NoError(t, err)
				assert.NotNil(t, edgeGateway)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.err.Error())
				assert.Nil(t, edgeGateway)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, test.expectedEdgeGateway, edgeGateway.EdgeGatewayModel)
		})
	}
}

func TestClient_ListEdgeGateway(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock client for cloudavenue.
	clientCAV := NewMockclientInterface(ctrl)

	c, _ := NewFakeClient(clientCAV)

	edgeGatewayID := urn.Gateway.String() + uuid.New().String()
	edgeGatewayID2 := urn.Gateway.String() + uuid.New().String()
	vdcID := urn.VDC.String() + uuid.New().String()

	tests := []struct {
		name                string
		mockFunc            func()
		expectedEdgeGateway []*EdgeGatewayModel
		expectedError       bool
		err                 error
	}{
		{
			name: "success",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAllNsxtEdgeGateways(nil).Return([]*govcd.NsxtEdgeGateway{
					{
						EdgeGateway: &govcdtypes.OpenAPIEdgeGateway{
							ID:          edgeGatewayID,
							Name:        "test-edge-gateway",
							Description: "test description",
							Status:      "ACTIVE",
							OwnerRef: &govcdtypes.OpenApiReference{
								ID:   vdcID,
								Name: "test-vdc",
							},
						},
					},
					{
						EdgeGateway: &govcdtypes.OpenAPIEdgeGateway{
							ID:          edgeGatewayID2,
							Name:        "test-edge-gateway-2",
							Description: "test description 2",
							Status:      "ACTIVE",
							OwnerRef: &govcdtypes.OpenApiReference{
								ID:   vdcID,
								Name: "test-vdc",
							},
						},
					},
				}, nil)
				// mock getBandwidth test-edge-gateway
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())
					responder, err := httpmock.NewJsonResponder(200, json.RawMessage(`{"rateLimit":10}`))
					if err != nil {
						t.Fatal(err)
					}

					httpmock.RegisterResponder("GET", endpoints.InlineTemplate(endpoints.EdgeGatewayGet, map[string]string{"edge-id": urn.ExtractUUID(edgeGatewayID)}), responder)
					return clientcloudavenue.MockClient().R()
				})
				// mock getBandwidth test-edge-gateway-2
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())
					responder, err := httpmock.NewJsonResponder(200, json.RawMessage(`{"rateLimit":100}`))
					if err != nil {
						t.Fatal(err)
					}

					httpmock.RegisterResponder("GET", endpoints.InlineTemplate(endpoints.EdgeGatewayGet, map[string]string{"edge-id": urn.ExtractUUID(edgeGatewayID2)}), responder)
					return clientcloudavenue.MockClient().R()
				})
			},
			expectedEdgeGateway: []*EdgeGatewayModel{
				{
					ID:          edgeGatewayID,
					Name:        "test-edge-gateway",
					Description: "test description",
					Status:      "ACTIVE",
					OwnerRef: &govcdtypes.OpenApiReference{
						ID:   vdcID,
						Name: "test-vdc",
					},
					Bandwidth: 10,
				},
				{
					ID:          edgeGatewayID2,
					Name:        "test-edge-gateway-2",
					Description: "test description 2",
					Status:      "ACTIVE",
					OwnerRef: &govcdtypes.OpenApiReference{
						ID:   vdcID,
						Name: "test-vdc",
					},
					Bandwidth: 100,
				},
			},
			expectedError: false,
		},
		{
			name: "error-get-all-edge-gateways",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAllNsxtEdgeGateways(nil).Return(nil, fmt.Errorf("error"))
			},
			expectedEdgeGateway: nil,
			expectedError:       true,
			err:                 fmt.Errorf("error"),
		},
		{
			name: "refresh-error",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(fmt.Errorf("refresh error"))
			},
			expectedEdgeGateway: nil,
			expectedError:       true,
			err:                 fmt.Errorf("refresh error"),
		},
		{
			name: "error-get-bandwidth",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAllNsxtEdgeGateways(nil).Return([]*govcd.NsxtEdgeGateway{
					{
						EdgeGateway: &govcdtypes.OpenAPIEdgeGateway{
							ID:          edgeGatewayID,
							Name:        "test-edge-gateway",
							Description: "test description",
							Status:      "ACTIVE",
							OwnerRef: &govcdtypes.OpenApiReference{
								ID:   vdcID,
								Name: "test-vdc",
							},
						},
					},
					{
						EdgeGateway: &govcdtypes.OpenAPIEdgeGateway{
							ID:          edgeGatewayID2,
							Name:        "test-edge-gateway-2",
							Description: "test description 2",
							Status:      "ACTIVE",
							OwnerRef: &govcdtypes.OpenApiReference{
								ID:   vdcID,
								Name: "test-vdc",
							},
						},
					},
				}, nil)
				// mock getBandwidth test-edge-gateway
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())
					responder, err := httpmock.NewJsonResponder(200, json.RawMessage(`{"rateLimit":10}`))
					if err != nil {
						t.Fatal(err)
					}

					httpmock.RegisterResponder("GET", endpoints.InlineTemplate(endpoints.EdgeGatewayGet, map[string]string{"edge-id": urn.ExtractUUID(edgeGatewayID)}), responder)
					return clientcloudavenue.MockClient().R()
				})

				// mock getBandwidth test-edge-gateway-2
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())
					responder := httpmock.NewErrorResponder(fmt.Errorf("error"))
					httpmock.RegisterResponder("GET", endpoints.InlineTemplate(endpoints.EdgeGatewayGet, map[string]string{"edge-id": urn.ExtractUUID(edgeGatewayID2)}), responder)
					return clientcloudavenue.MockClient().R()
				})
			},
			expectedEdgeGateway: nil,
			expectedError:       true,
			err:                 fmt.Errorf("error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockFunc()

			allEdgeGateways, err := c.ListEdgeGateway(context.Background())
			if !test.expectedError {
				assert.NoError(t, err)
				assert.NotNil(t, allEdgeGateways)
				assert.Equal(t, test.expectedEdgeGateway, allEdgeGateways)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.err.Error())
				assert.Nil(t, allEdgeGateways)
				return
			}
		})
	}
}

// TODO Wait migration VDC and VDCGroup to new SDK
// func TestClient_CreateEdgeGateway(t *testing.T) {

// 	// Mock controller.
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	defer httpmock.DeactivateAndReset()

// 	// Mock client for cloudavenue.
// 	clientCAV := NewMockclientInterface(ctrl)

// 	c, _ := NewFakeClient(clientCAV)

// 	vdcID := urn.VDC.String() + uuid.New().String()
// 	vdcGroupID := urn.VDCGroup.String() + uuid.New().String()

// 	tests := []struct {
// 		name                string
// 		mockFunc            func()
// 		expectedEdgeGateway *EdgeGatewayModel
// 		expectedError       bool
// 		err                 error
// 	}{
// 		{
// 			name: "create-edge-gateway-vdc-success",
// 			mockFunc: func() {
// 				clientCAV.EXPECT().Refresh().Return(nil)
// 				clientCAV.EXPECT().GetVDCById(vdcID, true).Return(&govcd.Vdc{
// 					Vdc: &govcdtypes.OpenAPIVdc{
// 						ID:   vdcID,
// 						Name: "test-vdc",
// 					},
// 				}, nil)

// }

func TestClient_UpdateEdgeGateway(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	defer httpmock.DeactivateAndReset()

	// Mock client for cloudavenue.
	clientCAV := NewMockclientInterface(ctrl)

	c, _ := NewFakeClient(clientCAV)

	edgeGatewayID := urn.Gateway.String() + uuid.New().String()

	jobStatusID := uuid.New().String()

	tests := []struct {
		name          string
		edgeID        string
		mockFunc      func()
		bandwidth     int
		expectedError bool
		err           error
	}{
		{
			name:   "success",
			edgeID: edgeGatewayID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)

				// mock getBandwidth test-edge-gateway
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())

					// * mock getBandwidth
					responder, err := httpmock.NewJsonResponder(200, json.RawMessage(`{"message":"Job successfully created","jobId":"`+jobStatusID+`"}`))
					if err != nil {
						t.Fatal(err)
					}

					httpmock.RegisterResponder("PUT", endpoints.InlineTemplate(endpoints.EdgeGatewayUpdate, map[string]string{"edge-id": urn.ExtractUUID(edgeGatewayID)}), responder)

					// * mock getJobStatus
					responderJob, err := httpmock.NewJsonResponder(200, json.RawMessage(`[{"actions":[],"description":"string","name":"string","status":"DONE"}]`))
					if err != nil {
						t.Fatal(err)
					}

					httpmock.RegisterResponder("GET", endpoints.InlineTemplate(endpoints.JobStatusGet, map[string]string{"job-id": jobStatusID}), responderJob)
					return clientcloudavenue.MockClient().R()
				})
			},
			bandwidth:     10,
			expectedError: false,
		},
		{
			name:   "error-no-id",
			edgeID: "",
			mockFunc: func() {
			},
			bandwidth:     10,
			expectedError: true,
			err:           fmt.Errorf("Error:Field validation for 'ID' failed on the 'required' tag"),
		},
		{
			name:   "error-update-edge-gateway",
			edgeID: edgeGatewayID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)

				// mock getBandwidth test-edge-gateway
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())

					// * mock getBandwidth
					responder := httpmock.NewErrorResponder(fmt.Errorf("error"))
					httpmock.RegisterResponder("PUT", endpoints.InlineTemplate(endpoints.EdgeGatewayUpdate, map[string]string{"edge-id": urn.ExtractUUID(edgeGatewayID)}), responder)
					return clientcloudavenue.MockClient().R()
				})
			},
			bandwidth:     10,
			expectedError: true,
			err:           fmt.Errorf("error"),
		},
		{
			name:   "error-refresh",
			edgeID: edgeGatewayID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(fmt.Errorf("error"))
			},
			bandwidth:     10,
			expectedError: true,
			err:           fmt.Errorf("error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockFunc()

			err := c.UpdateEdgeGateway(context.Background(), &EdgeGatewayModelUpdate{
				ID:        test.edgeID,
				Bandwidth: test.bandwidth,
			})

			if !test.expectedError {
				assert.NoError(t, err)
				return
			}
			assert.Error(t, err)
			assert.Contains(t, err.Error(), test.err.Error())
		})
	}
}

func TestClient_EnableNetworkService(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	defer httpmock.DeactivateAndReset()

	// Mock client for cloudavenue.
	clientCAV := NewMockclientInterface(ctrl)

	e := newFakeEdgeGatewayClient(clientCAV)

	edgeGatewayID := urn.Gateway.String() + uuid.New().String()
	vdcID := urn.VDC.String() + uuid.New().String()
	jobStatusID := uuid.New().String()

	tests := []struct {
		name          string
		mockFunc      func()
		expectedError bool
		err           error
	}{
		{
			name: "enable-network-service",
			mockFunc: func() {
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())

					// * mock enableNetworkService
					responder, err := httpmock.NewJsonResponder(200, json.RawMessage(`{"message":"Job successfully created","jobId":"`+jobStatusID+`"}`))
					if err != nil {
						t.Fatal(err)
					}

					httpmock.RegisterResponder("POST", endpoints.NetworkServiceCreate, responder)
					// 	return clientcloudavenue.MockClient().R()
					// })

					// clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					// * mock getJobStatus
					responderJob, err := httpmock.NewJsonResponder(200, json.RawMessage(`[{"actions":[],"description":"string","name":"string","status":"DONE"}]`))
					if err != nil {
						t.Fatal(err)
					}

					httpmock.RegisterResponder("GET", endpoints.InlineTemplate(endpoints.JobStatusGet, map[string]string{"job-id": jobStatusID}), responderJob)

					return clientcloudavenue.MockClient().R()
				})

				// mock getNetworkServices
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())
					responder, err := httpmock.NewJsonResponder(200, json.RawMessage(fmt.Sprintf(`
[
   {
      "type": "tier-0-vrf",
      "name": "prvrf01eocb0009999allsp01",
      "displayName": "",
      "properties": {
         "classOfService": "SHARED_STANDARD"
      },
      "children": [
         {
            "type": "edge-gateway",
            "name": "test-edge-gateway",
            "displayName": "",
            "properties": {
               "rateLimit": 5,
               "edgeUUID": "%s"
            },
            "children": [
               {
                  "type": "load-balancer",
                  "name": "737b9768-95a0-4955-bbbe-d5eab846e8dc",
                  "displayName": "v999w99eprnxcdshrdsegp99",
                  "properties": {
                     "classOfService": "PREMIUM",
                     "maxVirtualServices": 10
                  },
                  "children": []
               },
               {
                  "type": "service",
                  "name": "internet",
                  "displayName": "internet",
                  "properties": {
                     "ip": "12.123.123.12",
                     "announced": true
                  },
                  "children": [],
                  "serviceId": "ip-12-123-123-12"
               },
               {
                  "type": "service",
                  "name": "cav-services",
                  "displayName": "Cloud Avenue Services",
                  "properties": {
                     "ranges": [
                        "100.113.99.96/27"
                     ],
                     "ipCount": 16
                  },
                  "children": [],
                  "serviceId": "tn99e99ocb0009999spt199-cav-services"
               }
            ]
         }
      ]
   }
]
`, urn.ExtractUUID(edgeGatewayID))))
					if err != nil {
						t.Fatal(err)
					}

					httpmock.RegisterResponder("GET", endpoints.NetworkServiceGet, responder)
					return clientcloudavenue.MockClient().R()
				})
			},
			expectedError: false,
		},
		{
			name: "error-enable-network-service",
			mockFunc: func() {
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())

					// * mock enableNetworkService
					responder := httpmock.NewErrorResponder(fmt.Errorf("error"))
					httpmock.RegisterResponder("POST", endpoints.NetworkServiceCreate, responder)
					return clientcloudavenue.MockClient().R()
				})
			},
			expectedError: true,
			err:           fmt.Errorf("error"),
		},
		{
			name: "error-500-enable-network-service",
			mockFunc: func() {
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())

					// * mock enableNetworkService
					responder, err := httpmock.NewJsonResponder(500, json.RawMessage(`{"message":"error"}`))
					if err != nil {
						t.Fatal(err)
					}

					httpmock.RegisterResponder("POST", endpoints.NetworkServiceCreate, responder)
					return clientcloudavenue.MockClient().R()
				})
			},
			expectedError: true,
			err:           fmt.Errorf("error"),
		},
		{
			name: "error-get-job-status",
			mockFunc: func() {
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())

					// * mock enableNetworkService
					responder, err := httpmock.NewJsonResponder(200, json.RawMessage(`{"message":"Job successfully created","jobId":"`+jobStatusID+`"}`))
					if err != nil {
						t.Fatal(err)
					}

					httpmock.RegisterResponder("POST", endpoints.NetworkServiceCreate, responder)

					// * mock getJobStatus
					responderJob := httpmock.NewErrorResponder(fmt.Errorf("error"))
					httpmock.RegisterResponder("GET", endpoints.InlineTemplate(endpoints.JobStatusGet, map[string]string{"job-id": jobStatusID}), responderJob)

					return clientcloudavenue.MockClient().R()
				})
			},
			expectedError: true,
			err:           fmt.Errorf("error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockFunc()

			e.EdgeGatewayModel = &EdgeGatewayModel{
				ID:   edgeGatewayID,
				Name: "test-edge-gateway",
				OwnerRef: &govcdtypes.OpenApiReference{
					ID:   vdcID,
					Name: "test-vdc",
				},
			}

			err := e.EnableNetworkService(context.Background())
			if !test.expectedError {
				assert.NoError(t, err)
				return
			}
			assert.Error(t, err)
			assert.Contains(t, err.Error(), test.err.Error())
		})
	}
}

func TestClient_DisableNetworkService(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	defer httpmock.DeactivateAndReset()

	// Mock client for cloudavenue.
	clientCAV := NewMockclientInterface(ctrl)

	e := newFakeEdgeGatewayClient(clientCAV)

	edgeGatewayID := urn.Gateway.String() + uuid.New().String()
	vdcID := urn.VDC.String() + uuid.New().String()
	jobStatusID := uuid.New().String()

	tests := []struct {
		name          string
		mockFunc      func()
		expectedError bool
		err           error
	}{
		{
			name: "disable-network-service",
			mockFunc: func() {
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())

					// * mock disableNetworkService
					responder, err := httpmock.NewJsonResponder(200, json.RawMessage(`{"message":"Job successfully created","jobId":"`+jobStatusID+`"}`))
					if err != nil {
						t.Fatal(err)
					}

					httpmock.RegisterResponder("DELETE", endpoints.InlineTemplate(endpoints.NetworkServiceDelete, map[string]string{"service-id": "tn99e99ocb0009999spt199-cav-services"}), responder)

					// * mock getJobStatus
					responderJob, err := httpmock.NewJsonResponder(200, json.RawMessage(`[{"actions":[],"description":"string","name":"string","status":"DONE"}]`))
					if err != nil {
						t.Fatal(err)
					}

					httpmock.RegisterResponder("GET", endpoints.InlineTemplate(endpoints.JobStatusGet, map[string]string{"job-id": jobStatusID}), responderJob)
					return clientcloudavenue.MockClient().R()
				})

				// mock getNetworkServices
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())
					responder, err := httpmock.NewJsonResponder(200, json.RawMessage(fmt.Sprintf(`
[
   {
      "type": "tier-0-vrf",
      "name": "prvrf01eocb0009999allsp01",
      "displayName": "",
      "properties": {
         "classOfService": "SHARED_STANDARD"
      },
      "children": [
         {
            "type": "edge-gateway",
            "name": "test-edge-gateway",
            "displayName": "",
            "properties": {
               "rateLimit": 5,
               "edgeUUID": "%s"
            },
            "children": [
               {
                  "type": "load-balancer",
                  "name": "737b9768-95a0-4955-bbbe-d5eab846e8dc",
                  "displayName": "v999w99eprnxcdshrdsegp99",
                  "properties": {
                     "classOfService": "PREMIUM",
                     "maxVirtualServices": 10
                  },
                  "children": []
               },
               {
                  "type": "service",
                  "name": "internet",
                  "displayName": "internet",
                  "properties": {
                     "ip": "12.123.123.12",
                     "announced": true
                  },
                  "children": [],
                  "serviceId": "ip-12-123-123-12"
               },
               {
                  "type": "service",
                  "name": "cav-services",
                  "displayName": "Cloud Avenue Services",
                  "properties": {
                     "ranges": [
                        "100.113.99.96/27"
                     ],
                     "ipCount": 16
                  },
                  "children": [],
                  "serviceId": "tn99e99ocb0009999spt199-cav-services"
               }
            ]
         }
      ]
   }
]`, urn.ExtractUUID(edgeGatewayID))))
					if err != nil {
						t.Fatal(err)
					}

					httpmock.RegisterResponder("GET", endpoints.NetworkServiceGet, responder)
					return clientcloudavenue.MockClient().R()
				})
			},
			expectedError: false,
		},
		{
			name: "error-disable-network-service",
			mockFunc: func() {
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())

					// * mock disableNetworkService
					responder := httpmock.NewErrorResponder(fmt.Errorf("error"))
					httpmock.RegisterResponder("DELETE", endpoints.InlineTemplate(endpoints.NetworkServiceDelete, map[string]string{"service-id": "tn99e99ocb0009999spt199-cav-services"}), responder)
					return clientcloudavenue.MockClient().R()
				})
			},
			expectedError: true,
			err:           fmt.Errorf("error"),
		},
		{
			name: "error-500-disable-network-service",
			mockFunc: func() {
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())

					// * mock disableNetworkService
					responder, err := httpmock.NewJsonResponder(500, json.RawMessage(`{"message":"error"}`))
					if err != nil {
						t.Fatal(err)
					}

					httpmock.RegisterResponder("DELETE", endpoints.InlineTemplate(endpoints.NetworkServiceDelete, map[string]string{"service-id": "tn99e99ocb0009999spt199-cav-services"}), responder)
					return clientcloudavenue.MockClient().R()
				})
			},
			expectedError: true,
			err:           fmt.Errorf("error"),
		},
		{
			name: "error-get-job-status",
			mockFunc: func() {
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())

					// * mock disableNetworkService
					responder, err := httpmock.NewJsonResponder(200, json.RawMessage(`{"message":"Job successfully created","jobId":"`+jobStatusID+`"}`))
					if err != nil {
						t.Fatal(err)
					}

					httpmock.RegisterResponder("DELETE", endpoints.InlineTemplate(endpoints.NetworkServiceDelete, map[string]string{"service-id": "tn99e99ocb0009999spt199-cav-services"}), responder)

					// * mock getJobStatus
					responderJob := httpmock.NewErrorResponder(fmt.Errorf("error"))
					httpmock.RegisterResponder("GET", endpoints.InlineTemplate(endpoints.JobStatusGet, map[string]string{"job-id": jobStatusID}), responderJob)

					return clientcloudavenue.MockClient().R()
				})
			},
			expectedError: true,
			err:           fmt.Errorf("error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockFunc()

			e.EdgeGatewayModel = &EdgeGatewayModel{
				ID:   edgeGatewayID,
				Name: "test-edge-gateway",
				OwnerRef: &govcdtypes.OpenApiReference{
					ID:   vdcID,
					Name: "test-vdc",
				},
				Services: NetworkServicesModelSvcs{
					Service: &NetworkServicesModelSvcService{
						NetworkServicesModelSvc: NetworkServicesModelSvc{
							ID:   "tn99e99ocb0009999spt199-cav-services",
							Name: "Cloud Avenue Services",
						},
					},
				},
			}

			err := e.DisableNetworkService(context.Background())
			if !test.expectedError {
				assert.NoError(t, err)
				return
			}
			assert.Error(t, err)
			assert.Contains(t, err.Error(), test.err.Error())
		})
	}
}

func TestClient_DeleteEdgeGateway(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	defer httpmock.DeactivateAndReset()

	// Mock client for cloudavenue.
	clientCAV := NewMockclientInterface(ctrl)

	c, _ := NewFakeClient(clientCAV)

	edgeGatewayID := urn.Gateway.String() + uuid.New().String()
	jobStatusID := uuid.New().String()
	vdcID := urn.VDC.String() + uuid.New().String()

	tests := []struct {
		name          string
		mockFunc      func()
		expectedError bool
		err           error
	}{
		{
			name: "delete-edge-gateway",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetNsxtEdgeGatewayById(edgeGatewayID).Return(&govcd.NsxtEdgeGateway{
					EdgeGateway: &govcdtypes.OpenAPIEdgeGateway{
						ID:          edgeGatewayID,
						Name:        "test-edge-gateway",
						Description: "test description",
						Status:      "ACTIVE",
						OwnerRef: &govcdtypes.OpenApiReference{
							ID:   vdcID,
							Name: "test-vdc",
						},
						EdgeGatewayUplinks: []govcdtypes.EdgeGatewayUplinks{
							{
								UplinkName: "prvrf01eocb0009999allsp01",
							},
						},
					},
				}, nil)
				// mock getBandwidth
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())
					responder, err := httpmock.NewJsonResponder(200, json.RawMessage(`{"rateLimit":10}`))
					if err != nil {
						t.Fatal(err)
					}

					httpmock.RegisterResponder("GET", endpoints.InlineTemplate(endpoints.EdgeGatewayGet, map[string]string{"edge-id": urn.ExtractUUID(edgeGatewayID)}), responder)
					return clientcloudavenue.MockClient().R()
				})
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())

					// * mock deleteEdgeGateway
					responder, err := httpmock.NewJsonResponder(200, json.RawMessage(`{"message":"Job successfully created","jobId":"`+jobStatusID+`"}`))
					if err != nil {
						t.Fatal(err)
					}

					httpmock.RegisterResponder("DELETE", endpoints.InlineTemplate(endpoints.EdgeGatewayDelete, map[string]string{"edge-id": urn.ExtractUUID(edgeGatewayID)}), responder)

					// * mock getJobStatus
					responderJob, err := httpmock.NewJsonResponder(200, json.RawMessage(`[{"actions":[],"description":"string","name":"string","status":"DONE"}]`))
					if err != nil {
						t.Fatal(err)
					}

					httpmock.RegisterResponder("GET", endpoints.InlineTemplate(endpoints.JobStatusGet, map[string]string{"job-id": jobStatusID}), responderJob)
					return clientcloudavenue.MockClient().R()
				})
			},
			expectedError: false,
		},
		{
			name: "error-500-delete-edge-gateway",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetNsxtEdgeGatewayById(edgeGatewayID).Return(&govcd.NsxtEdgeGateway{
					EdgeGateway: &govcdtypes.OpenAPIEdgeGateway{
						ID:          edgeGatewayID,
						Name:        "test-edge-gateway",
						Description: "test description",
						Status:      "ACTIVE",
						OwnerRef: &govcdtypes.OpenApiReference{
							ID:   vdcID,
							Name: "test-vdc",
						},
						EdgeGatewayUplinks: []govcdtypes.EdgeGatewayUplinks{
							{
								UplinkName: "prvrf01eocb0009999allsp01",
							},
						},
					},
				}, nil)
				// mock getBandwidth
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())
					responder, err := httpmock.NewJsonResponder(200, json.RawMessage(`{"rateLimit":10}`))
					if err != nil {
						t.Fatal(err)
					}

					httpmock.RegisterResponder("GET", endpoints.InlineTemplate(endpoints.EdgeGatewayGet, map[string]string{"edge-id": urn.ExtractUUID(edgeGatewayID)}), responder)
					return clientcloudavenue.MockClient().R()
				})
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())

					// * mock deleteEdgeGateway
					responder, err := httpmock.NewJsonResponder(500, json.RawMessage(`{"message":"error"}`))
					if err != nil {
						t.Fatal(err)
					}
					httpmock.RegisterResponder("DELETE", endpoints.InlineTemplate(endpoints.EdgeGatewayDelete, map[string]string{"edge-id": urn.ExtractUUID(edgeGatewayID)}), responder)
					return clientcloudavenue.MockClient().R()
				})
			},
			expectedError: true,
			err:           fmt.Errorf("error"),
		},
		{
			name: "error-delete-edge-gateway",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetNsxtEdgeGatewayById(edgeGatewayID).Return(&govcd.NsxtEdgeGateway{
					EdgeGateway: &govcdtypes.OpenAPIEdgeGateway{
						ID:          edgeGatewayID,
						Name:        "test-edge-gateway",
						Description: "test description",
						Status:      "ACTIVE",
						OwnerRef: &govcdtypes.OpenApiReference{
							ID:   vdcID,
							Name: "test-vdc",
						},
						EdgeGatewayUplinks: []govcdtypes.EdgeGatewayUplinks{
							{
								UplinkName: "prvrf01eocb0009999allsp01",
							},
						},
					},
				}, nil)
				// mock getBandwidth
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())
					responder, err := httpmock.NewJsonResponder(200, json.RawMessage(`{"rateLimit":10}`))
					if err != nil {
						t.Fatal(err)
					}

					httpmock.RegisterResponder("GET", endpoints.InlineTemplate(endpoints.EdgeGatewayGet, map[string]string{"edge-id": urn.ExtractUUID(edgeGatewayID)}), responder)
					return clientcloudavenue.MockClient().R()
				})
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())

					// * mock deleteEdgeGateway
					responder := httpmock.NewErrorResponder(fmt.Errorf("error"))
					httpmock.RegisterResponder("DELETE", endpoints.InlineTemplate(endpoints.EdgeGatewayDelete, map[string]string{"edge-id": urn.ExtractUUID(edgeGatewayID)}), responder)
					return clientcloudavenue.MockClient().R()
				})
			},
			expectedError: true,
			err:           fmt.Errorf("error"),
		},
		{
			name: "error-get-job-status",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetNsxtEdgeGatewayById(edgeGatewayID).Return(&govcd.NsxtEdgeGateway{
					EdgeGateway: &govcdtypes.OpenAPIEdgeGateway{
						ID:          edgeGatewayID,
						Name:        "test-edge-gateway",
						Description: "test description",
						Status:      "ACTIVE",
						OwnerRef: &govcdtypes.OpenApiReference{
							ID:   vdcID,
							Name: "test-vdc",
						},
						EdgeGatewayUplinks: []govcdtypes.EdgeGatewayUplinks{
							{
								UplinkName: "prvrf01eocb0009999allsp01",
							},
						},
					},
				}, nil)
				// mock getBandwidth
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())
					responder, err := httpmock.NewJsonResponder(200, json.RawMessage(`{"rateLimit":10}`))
					if err != nil {
						t.Fatal(err)
					}

					httpmock.RegisterResponder("GET", endpoints.InlineTemplate(endpoints.EdgeGatewayGet, map[string]string{"edge-id": urn.ExtractUUID(edgeGatewayID)}), responder)
					return clientcloudavenue.MockClient().R()
				})
				clientCAV.EXPECT().R().DoAndReturn(func() *resty.Request {
					httpmock.ActivateNonDefault(clientcloudavenue.MockClient().GetClient())

					// * mock deleteEdgeGateway
					responder, err := httpmock.NewJsonResponder(200, json.RawMessage(`{"message":"Job successfully created","jobId":"`+jobStatusID+`"}`))
					if err != nil {
						t.Fatal(err)
					}

					httpmock.RegisterResponder("DELETE", endpoints.InlineTemplate(endpoints.EdgeGatewayDelete, map[string]string{"edge-id": urn.ExtractUUID(edgeGatewayID)}), responder)

					// * mock getJobStatus
					responderJob := httpmock.NewErrorResponder(fmt.Errorf("error"))
					httpmock.RegisterResponder("GET", endpoints.InlineTemplate(endpoints.JobStatusGet, map[string]string{"job-id": jobStatusID}), responderJob)

					return clientcloudavenue.MockClient().R()
				})
			},
			expectedError: true,
			err:           fmt.Errorf("error"),
		},
		{
			name: "error-refresh",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(fmt.Errorf("error"))
			},
			expectedError: true,
			err:           fmt.Errorf("error"),
		},
		{
			name: "error-get-edge-gateway",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetNsxtEdgeGatewayById(edgeGatewayID).Return(nil, fmt.Errorf("error"))
			},
			expectedError: true,
			err:           fmt.Errorf("error"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockFunc()

			err := c.DeleteEdgeGateway(context.Background(), edgeGatewayID)
			if !test.expectedError {
				assert.NoError(t, err)
				return
			}
			assert.Error(t, err)
			assert.Contains(t, err.Error(), test.err.Error())
		})
	}
}
