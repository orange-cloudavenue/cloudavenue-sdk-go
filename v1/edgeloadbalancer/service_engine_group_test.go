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
	"fmt"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/internal/utils"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

func TestClient_ListServiceEngineGroups(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock client for cloudavenue.
	clientCAV := NewMockclientFake(ctrl)

	c, _ := NewFakeClient(clientCAV)

	urnServiceEngineGroup := urn.ServiceEngineGroup.String() + uuid.New().String()
	urnEdgeGateway := urn.Gateway.String() + uuid.New().String()

	tests := []struct {
		name              string
		mockFunc          func()
		expectedCertValue []*ServiceEngineGroupModel
		expectedErr       bool
		edgeGatewayID     string
		err               error
	}{
		{
			name:          "success",
			edgeGatewayID: urnEdgeGateway,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)

				v := url.Values{}
				v.Add("filter", "gatewayRef.id=="+urnEdgeGateway)
				clientCAV.EXPECT().GetAllAlbServiceEngineGroupAssignments(gomock.AssignableToTypeOf(v)).Return([]*govcd.NsxtAlbServiceEngineGroupAssignment{
					{
						NsxtAlbServiceEngineGroupAssignment: &govcdtypes.NsxtAlbServiceEngineGroupAssignment{
							ServiceEngineGroupRef: &govcdtypes.OpenApiReference{
								ID:   urnServiceEngineGroup,
								Name: "name",
							},
							GatewayRef: &govcdtypes.OpenApiReference{
								ID:   urnEdgeGateway,
								Name: "edge_name",
							},
							MaxVirtualServices:         utils.ToPTR(10),
							MinVirtualServices:         utils.ToPTR(1),
							NumDeployedVirtualServices: 2,
						},
					},
				}, nil)
			},
			expectedCertValue: []*ServiceEngineGroupModel{
				{
					ID:   urnServiceEngineGroup,
					Name: "name",
					GatewayRef: &govcdtypes.OpenApiReference{
						ID:   urnEdgeGateway,
						Name: "edge_name",
					},
					MaxVirtualServices:         utils.ToPTR(10),
					MinVirtualServices:         utils.ToPTR(1),
					NumDeployedVirtualServices: 2,
				},
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name:          "refresh-error",
			edgeGatewayID: urnEdgeGateway,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(errors.New("error"))
			},
			expectedCertValue: nil,
			expectedErr:       true,
			err:               errors.New("error"),
		},
		{
			name:          "error-get-all-certificates",
			edgeGatewayID: urnEdgeGateway,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)

				v := url.Values{}
				v.Add("filter", "gatewayRef.id=="+urnEdgeGateway)
				clientCAV.EXPECT().GetAllAlbServiceEngineGroupAssignments(gomock.AssignableToTypeOf(v)).Return(nil, errors.New("error"))
			},
			expectedCertValue: nil,
			expectedErr:       true,
		},
		{
			name:          "error-get-all-certificates-nil",
			edgeGatewayID: urnEdgeGateway,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)

				v := url.Values{}
				v.Add("filter", "gatewayRef.id=="+urnEdgeGateway)
				clientCAV.EXPECT().GetAllAlbServiceEngineGroupAssignments(gomock.AssignableToTypeOf(v)).Return(nil, nil)
			},
			expectedCertValue: nil,
			expectedErr:       true,
			err:               fmt.Errorf("no service engine group found for edge gateway %s. The service Load Balancer might not be enabled on this edge gateway. Contact the support", urnEdgeGateway),
		},
		{
			name:          "error-validation-edgeGateway-ID-empty",
			edgeGatewayID: "",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
			},
			expectedCertValue: nil,
			expectedErr:       true,
			err:               fmt.Errorf("edgeGatewayID cannot be empty"),
		},
		{
			name:          "error-validation-edgeGateway-ID-not-urn",
			edgeGatewayID: "not-urn",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
			},
			expectedCertValue: nil,
			expectedErr:       true,
			err:               fmt.Errorf("edgeGatewayID is not a valid URN"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			certificates, err := c.ListServiceEngineGroups(context.Background(), tc.edgeGatewayID)
			if !tc.expectedErr {
				assert.NoError(t, err)
				assert.NotNil(t, certificates)
			} else {
				assert.Error(t, err)
				assert.Nil(t, certificates)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCertValue, certificates)
		})
	}
}

func TestClient_GetServiceEngineGroup(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock client for cloudavenue.
	clientCAV := NewMockclientFake(ctrl)

	c, _ := NewFakeClient(clientCAV)

	urnServiceEngineGroup := urn.ServiceEngineGroup.String() + uuid.New().String()
	urnEdgeGateway := urn.Gateway.String() + uuid.New().String()

	tests := []struct {
		name              string
		mockFunc          func()
		expectedCertValue *ServiceEngineGroupModel
		expectedErr       bool
		edgeGatewayID     string
		err               error
		nameOrID          string
	}{
		{
			name:          "success",
			edgeGatewayID: urnEdgeGateway,
			nameOrID:      urnServiceEngineGroup,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)

				v := url.Values{}
				v.Add("filter", "gatewayRef.id=="+urnEdgeGateway+";serviceEngineGroupRef.id=="+urnServiceEngineGroup)
				clientCAV.EXPECT().GetAllAlbServiceEngineGroupAssignments(gomock.AssignableToTypeOf(v)).Return([]*govcd.NsxtAlbServiceEngineGroupAssignment{
					{
						NsxtAlbServiceEngineGroupAssignment: &govcdtypes.NsxtAlbServiceEngineGroupAssignment{
							ServiceEngineGroupRef: &govcdtypes.OpenApiReference{
								ID:   urnServiceEngineGroup,
								Name: "name",
							},
							GatewayRef: &govcdtypes.OpenApiReference{
								ID:   urnEdgeGateway,
								Name: "edge_name",
							},
							MaxVirtualServices:         utils.ToPTR(10),
							MinVirtualServices:         utils.ToPTR(1),
							NumDeployedVirtualServices: 2,
						},
					},
				}, nil)
			},
			expectedCertValue: &ServiceEngineGroupModel{
				ID:   urnServiceEngineGroup,
				Name: "name",
				GatewayRef: &govcdtypes.OpenApiReference{
					ID:   urnEdgeGateway,
					Name: "edge_name",
				},
				MaxVirtualServices:         utils.ToPTR(10),
				MinVirtualServices:         utils.ToPTR(1),
				NumDeployedVirtualServices: 2,
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name:          "success-name",
			edgeGatewayID: urnEdgeGateway,
			nameOrID:      "name",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)

				v := url.Values{}
				v.Add("filter", "gatewayRef.id=="+urnEdgeGateway)
				clientCAV.EXPECT().GetAllAlbServiceEngineGroupAssignments(gomock.AssignableToTypeOf(v)).Return([]*govcd.NsxtAlbServiceEngineGroupAssignment{
					{
						NsxtAlbServiceEngineGroupAssignment: &govcdtypes.NsxtAlbServiceEngineGroupAssignment{
							ServiceEngineGroupRef: &govcdtypes.OpenApiReference{
								ID:   urnServiceEngineGroup,
								Name: "name",
							},
							GatewayRef: &govcdtypes.OpenApiReference{
								ID:   urnEdgeGateway,
								Name: "edge_name",
							},
							MaxVirtualServices:         utils.ToPTR(10),
							MinVirtualServices:         utils.ToPTR(1),
							NumDeployedVirtualServices: 2,
						},
					},
				}, nil)
			},
			expectedCertValue: &ServiceEngineGroupModel{
				ID:   urnServiceEngineGroup,
				Name: "name",
				GatewayRef: &govcdtypes.OpenApiReference{
					ID:   urnEdgeGateway,
					Name: "edge_name",
				},
				MaxVirtualServices:         utils.ToPTR(10),
				MinVirtualServices:         utils.ToPTR(1),
				NumDeployedVirtualServices: 2,
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name:          "refresh-error",
			edgeGatewayID: urnEdgeGateway,
			nameOrID:      urnServiceEngineGroup,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(errors.New("error"))
			},
			expectedCertValue: nil,
			expectedErr:       true,
			err:               errors.New("error"),
		},
		{
			name:          "error-get-all-certificates",
			edgeGatewayID: urnEdgeGateway,
			nameOrID:      urnServiceEngineGroup,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)

				v := url.Values{}
				v.Add("filter", "gatewayRef.id=="+urnEdgeGateway+";serviceEngineGroupRef.id=="+urnServiceEngineGroup)
				clientCAV.EXPECT().GetAllAlbServiceEngineGroupAssignments(gomock.AssignableToTypeOf(v)).Return(nil, errors.New("error"))
			},
			expectedCertValue: nil,
			expectedErr:       true,
		},
		{
			name:          "error-service-engine-group-not-found",
			edgeGatewayID: urnEdgeGateway,
			nameOrID:      "notfound",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)

				v := url.Values{}
				v.Add("filter", "gatewayRef.id=="+urnEdgeGateway+";serviceEngineGroupRef.id=="+urnServiceEngineGroup)
				clientCAV.EXPECT().GetAllAlbServiceEngineGroupAssignments(gomock.AssignableToTypeOf(v)).Return([]*govcd.NsxtAlbServiceEngineGroupAssignment{
					{
						NsxtAlbServiceEngineGroupAssignment: &govcdtypes.NsxtAlbServiceEngineGroupAssignment{
							ServiceEngineGroupRef: &govcdtypes.OpenApiReference{
								ID:   urnServiceEngineGroup,
								Name: "name",
							},
							GatewayRef: &govcdtypes.OpenApiReference{
								ID:   urnEdgeGateway,
								Name: "edge_name",
							},
							MaxVirtualServices:         utils.ToPTR(10),
							MinVirtualServices:         utils.ToPTR(1),
							NumDeployedVirtualServices: 2,
						},
					},
				}, nil)
			},
			expectedCertValue: nil,
			expectedErr:       true,
			err:               fmt.Errorf("the service engine group %s was not found for edge gateway %s", "notfound", urnEdgeGateway),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			certificates, err := c.GetServiceEngineGroup(context.Background(), tc.edgeGatewayID, tc.nameOrID)
			if !tc.expectedErr {
				assert.NoError(t, err)
				assert.NotNil(t, certificates)
			} else {
				assert.Error(t, err)
				assert.Nil(t, certificates)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedCertValue, certificates)
		})
	}
}
