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

func TestClient_GetPool(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock client for cloudavenue.
	clientCAV := NewMockclientFake(ctrl)

	c, _ := NewFakeClient(clientCAV)

	urnEdgeGateway := urn.Gateway.String() + uuid.New().String()
	poolID := urn.LoadBalancerPool.String() + uuid.New().String()

	tests := []struct {
		name          string
		mockFunc      func()
		expectedValue *PoolModel
		expectedErr   bool
		edgeGatewayID string
		byNameOrID    string
		poolID        string
		poolName      string
		err           error
	}{
		{
			name:          "success-http-by-name",
			edgeGatewayID: urnEdgeGateway,
			poolID:        poolID,
			poolName:      "pool1",
			byNameOrID:    "name",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbPoolByName(urnEdgeGateway, "pool1").Return(&govcd.NsxtAlbPool{
					NsxtAlbPool: &govcdtypes.NsxtAlbPool{
						ID:                       poolID,
						Name:                     "pool1",
						Description:              "pool1 description",
						GatewayRef:               govcdtypes.OpenApiReference{ID: urnEdgeGateway},
						Enabled:                  utils.ToPTR(true),
						Algorithm:                "LEAST_CONNECTIONS",
						DefaultPort:              utils.ToPTR(80),
						GracefulTimeoutPeriod:    utils.ToPTR(10),
						PassiveMonitoringEnabled: utils.ToPTR(true),
						HealthMonitors: []govcdtypes.NsxtAlbPoolHealthMonitor{
							{
								Name: "monitor HTTP",
								Type: "HTTP",
							},
							{
								Name: "monitor TCP",
								Type: "TCP",
							},
						},
						Members: []govcdtypes.NsxtAlbPoolMember{
							{
								Enabled:               true,
								IpAddress:             "192.168.0.1",
								Port:                  80,
								Ratio:                 utils.ToPTR(1),
								MarkedDownBy:          nil,
								HealthStatus:          "UP",
								DetailedHealthMessage: "",
							},
						},
						MemberGroupRef:         nil, // Only one of MemberGroupRef or Members should be set.
						CaCertificateRefs:      nil,
						CommonNameCheckEnabled: utils.ToPTR(false), // false because CaCertificateRefs is nil.
						DomainNames:            nil,                // nil because CommonNameCheckEnabled is false.
						PersistenceProfile: &govcdtypes.NsxtAlbPoolPersistenceProfile{
							Name:  "persistence profile",
							Type:  "CLIENT_IP",
							Value: "",
						},
						MemberCount:        1,
						UpMemberCount:      1,
						HealthMessage:      "All members are up",
						VirtualServiceRefs: nil,
						SslEnabled:         utils.ToPTR(false),
					},
				}, nil)
			},
			expectedValue: &PoolModel{
				ID:                       poolID,
				Name:                     "pool1",
				Description:              "pool1 description",
				GatewayRef:               govcdtypes.OpenApiReference{ID: urnEdgeGateway},
				Enabled:                  utils.ToPTR(true),
				Algorithm:                PoolAlgorithmLeastConnections,
				DefaultPort:              utils.ToPTR(80),
				GracefulTimeoutPeriod:    utils.ToPTR(10),
				PassiveMonitoringEnabled: utils.ToPTR(true),
				HealthMonitors: []PoolModelHealthMonitor{
					{
						Name: "monitor HTTP",
						Type: PoolHealthMonitorTypeHTTP,
					},
					{
						Name: "monitor TCP",
						Type: PoolHealthMonitorTypeTCP,
					},
				},
				Members: []PoolModelMember{
					{
						Enabled:               true,
						IPAddress:             "192.168.0.1",
						Port:                  80,
						Ratio:                 utils.ToPTR(1),
						MarkedDownBy:          nil,
						HealthStatus:          "UP",
						DetailedHealthMessage: "",
					},
				},
				MemberGroupRef:         nil,
				CaCertificateRefs:      nil,
				CommonNameCheckEnabled: utils.ToPTR(false),
				DomainNames:            nil,
				PersistenceProfile: &PoolModelPersistenceProfile{
					Name:  "persistence profile",
					Type:  PoolPersistenceProfileTypeClientIP,
					Value: "",
				},
				MemberCount:        1,
				UpMemberCount:      1,
				HealthMessage:      "All members are up",
				VirtualServiceRefs: nil,
				SSLEnabled:         utils.ToPTR(false),
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name:          "success-http-by-id",
			edgeGatewayID: urnEdgeGateway,
			poolID:        poolID,
			poolName:      "pool1",
			byNameOrID:    "id",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbPoolById(poolID).Return(&govcd.NsxtAlbPool{
					NsxtAlbPool: &govcdtypes.NsxtAlbPool{
						ID:                       poolID,
						Name:                     "pool1",
						Description:              "pool1 description",
						GatewayRef:               govcdtypes.OpenApiReference{ID: urnEdgeGateway},
						Enabled:                  utils.ToPTR(true),
						Algorithm:                "LEAST_CONNECTIONS",
						DefaultPort:              utils.ToPTR(80),
						GracefulTimeoutPeriod:    utils.ToPTR(10),
						PassiveMonitoringEnabled: utils.ToPTR(true),
						HealthMonitors: []govcdtypes.NsxtAlbPoolHealthMonitor{
							{
								Name: "monitor HTTP",
								Type: "HTTP",
							},
							{
								Name: "monitor TCP",
								Type: "TCP",
							},
						},
						Members: []govcdtypes.NsxtAlbPoolMember{
							{
								Enabled:               true,
								IpAddress:             "192.168.0.1",
								Port:                  80,
								Ratio:                 utils.ToPTR(1),
								MarkedDownBy:          nil,
								HealthStatus:          "UP",
								DetailedHealthMessage: "",
							},
						},
						MemberGroupRef:         nil, // Only one of MemberGroupRef or Members should be set.
						CaCertificateRefs:      nil,
						CommonNameCheckEnabled: utils.ToPTR(false), // false because CaCertificateRefs is nil.
						DomainNames:            nil,                // nil because CommonNameCheckEnabled is false.
						PersistenceProfile:     nil,
						MemberCount:            1,
						UpMemberCount:          1,
						HealthMessage:          "All members are up",
						VirtualServiceRefs:     nil,
						SslEnabled:             utils.ToPTR(false),
					},
				}, nil)
			},
			expectedValue: &PoolModel{
				ID:                       poolID,
				Name:                     "pool1",
				Description:              "pool1 description",
				GatewayRef:               govcdtypes.OpenApiReference{ID: urnEdgeGateway},
				Enabled:                  utils.ToPTR(true),
				Algorithm:                PoolAlgorithmLeastConnections,
				DefaultPort:              utils.ToPTR(80),
				GracefulTimeoutPeriod:    utils.ToPTR(10),
				PassiveMonitoringEnabled: utils.ToPTR(true),
				HealthMonitors: []PoolModelHealthMonitor{
					{
						Name: "monitor HTTP",
						Type: PoolHealthMonitorTypeHTTP,
					},
					{
						Name: "monitor TCP",
						Type: PoolHealthMonitorTypeTCP,
					},
				},
				Members: []PoolModelMember{
					{
						Enabled:               true,
						IPAddress:             "192.168.0.1",
						Port:                  80,
						Ratio:                 utils.ToPTR(1),
						MarkedDownBy:          nil,
						HealthStatus:          "UP",
						DetailedHealthMessage: "",
					},
				},
				MemberGroupRef:         nil,
				CaCertificateRefs:      nil,
				CommonNameCheckEnabled: utils.ToPTR(false),
				DomainNames:            nil,
				PersistenceProfile:     nil,
				MemberCount:            1,
				UpMemberCount:          1,
				HealthMessage:          "All members are up",
				VirtualServiceRefs:     nil,
				SSLEnabled:             utils.ToPTR(false),
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name:          "refresh-error",
			edgeGatewayID: urnEdgeGateway,
			poolID:        poolID,
			byNameOrID:    "id",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(errors.New("error"))
			},
			expectedValue: nil,
			expectedErr:   true,
			err:           errors.New("error"),
		},
		{
			name:          "error-get",
			edgeGatewayID: urnEdgeGateway,
			poolID:        poolID,
			poolName:      "pool1",
			byNameOrID:    "id",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbPoolById(poolID).Return(nil, errors.New("error"))
			},
			expectedValue: &PoolModel{},
			expectedErr:   true,
			err:           errors.New("error"),
		},
		{
			name:          "param-edgeGatewayID-empty",
			edgeGatewayID: "",
			poolID:        poolID,
			poolName:      "pool1",
			byNameOrID:    "id",
			mockFunc: func() {
			},
			expectedValue: &PoolModel{},
			expectedErr:   true,
			err:           errors.New("edgeGatewayID is empty. Please provide a valid edgeGatewayID"),
		},
		{
			name:          "param-edgeGatewayID-invalid-id",
			edgeGatewayID: "1234",
			poolID:        "1234",
			poolName:      "pool1",
			byNameOrID:    "id",
			mockFunc: func() {
			},
			expectedValue: &PoolModel{},
			expectedErr:   true,
			err:           errors.New("edgeGatewayID has invalid format. Please provide a valid edgeGatewayID"),
		},
		{
			name:          "param-poolNameOrID-empty",
			edgeGatewayID: urnEdgeGateway,
			poolID:        "",
			byNameOrID:    "name",
			mockFunc: func() {
			},
			expectedValue: &PoolModel{},
			expectedErr:   true,
			err:           errors.New("poolNameOrID is empty. Please provide a valid poolNameOrID"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			nameOrID := tc.poolName
			if tc.byNameOrID == "id" {
				nameOrID = tc.poolID
			}

			pool, err := c.GetPool(context.Background(), tc.edgeGatewayID, nameOrID)
			if !tc.expectedErr {
				assert.NoError(t, err)
				assert.NotNil(t, pool)
			} else {
				assert.Error(t, err)
				assert.Nil(t, pool)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedValue, pool)
		})
	}
}

func TestClient_ListPools(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock client for cloudavenue.
	clientCAV := NewMockclientFake(ctrl)

	c, _ := NewFakeClient(clientCAV)

	urnEdgeGateway := urn.Gateway.String() + uuid.New().String()
	poolID := urn.LoadBalancerPool.String() + uuid.New().String()
	poolID2 := urn.LoadBalancerPool.String() + uuid.New().String()

	tests := []struct {
		name          string
		mockFunc      func()
		expectedValue []*PoolModel
		expectedErr   bool
		edgeGatewayID string
		err           error
	}{
		{
			name:          "success",
			edgeGatewayID: urnEdgeGateway,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil).Times(3)
				clientCAV.EXPECT().GetAllAlbPoolSummaries(urnEdgeGateway, gomock.AssignableToTypeOf(url.Values{})).Return([]*govcd.NsxtAlbPool{
					{
						NsxtAlbPool: &govcdtypes.NsxtAlbPool{
							ID:          poolID,
							Name:        "pool1",
							Description: "pool1 description",
							GatewayRef:  govcdtypes.OpenApiReference{ID: urnEdgeGateway},
						},
					},
					{
						NsxtAlbPool: &govcdtypes.NsxtAlbPool{
							ID:          poolID2,
							Name:        "pool2",
							Description: "pool2 description",
							GatewayRef:  govcdtypes.OpenApiReference{ID: urnEdgeGateway},
						},
					},
				}, nil)
				clientCAV.EXPECT().GetAlbPoolById(poolID).Return(&govcd.NsxtAlbPool{
					NsxtAlbPool: &govcdtypes.NsxtAlbPool{
						ID:                       poolID,
						Name:                     "pool1",
						Description:              "pool1 description",
						GatewayRef:               govcdtypes.OpenApiReference{ID: urnEdgeGateway},
						Enabled:                  utils.ToPTR(true),
						Algorithm:                "LEAST_CONNECTIONS",
						DefaultPort:              utils.ToPTR(80),
						GracefulTimeoutPeriod:    utils.ToPTR(10),
						PassiveMonitoringEnabled: utils.ToPTR(true),
						HealthMonitors: []govcdtypes.NsxtAlbPoolHealthMonitor{
							{
								Name: "monitor HTTP",
								Type: "HTTP",
							},
							{
								Name: "monitor TCP",
								Type: "TCP",
							},
						},
						Members: []govcdtypes.NsxtAlbPoolMember{
							{
								Enabled:               true,
								IpAddress:             "192.168.0.1",
								Port:                  80,
								Ratio:                 utils.ToPTR(1),
								MarkedDownBy:          nil,
								HealthStatus:          "UP",
								DetailedHealthMessage: "",
							},
						},
						MemberGroupRef:         nil, // Only one of MemberGroupRef or Members should be set.
						CaCertificateRefs:      nil,
						CommonNameCheckEnabled: utils.ToPTR(false), // false because CaCertificateRefs is nil.
						DomainNames:            nil,                // nil because CommonNameCheckEnabled is false.
						PersistenceProfile: &govcdtypes.NsxtAlbPoolPersistenceProfile{
							Name:  "persistence profile",
							Type:  "CLIENT_IP",
							Value: "",
						},
						MemberCount:        1,
						UpMemberCount:      1,
						HealthMessage:      "All members are up",
						VirtualServiceRefs: nil,
						SslEnabled:         utils.ToPTR(false),
					},
				}, nil)
				clientCAV.EXPECT().GetAlbPoolById(poolID2).Return(&govcd.NsxtAlbPool{
					NsxtAlbPool: &govcdtypes.NsxtAlbPool{
						ID:                       poolID2,
						Name:                     "pool2",
						Description:              "pool2 description",
						GatewayRef:               govcdtypes.OpenApiReference{ID: urnEdgeGateway},
						Enabled:                  utils.ToPTR(true),
						Algorithm:                "LEAST_CONNECTIONS",
						DefaultPort:              utils.ToPTR(80),
						GracefulTimeoutPeriod:    utils.ToPTR(10),
						PassiveMonitoringEnabled: utils.ToPTR(true),
						HealthMonitors: []govcdtypes.NsxtAlbPoolHealthMonitor{
							{
								Name: "monitor HTTP",
								Type: "HTTP",
							},
							{
								Name: "monitor TCP",
								Type: "TCP",
							},
						},
						Members: []govcdtypes.NsxtAlbPoolMember{
							{
								Enabled:               true,
								IpAddress:             "192.168.0.2",
								Port:                  80,
								Ratio:                 utils.ToPTR(1),
								MarkedDownBy:          nil,
								HealthStatus:          "UP",
								DetailedHealthMessage: "",
							},
						},
						MemberGroupRef:         nil, // Only one of MemberGroupRef or Members should be set.
						CaCertificateRefs:      nil,
						CommonNameCheckEnabled: utils.ToPTR(false), // false because CaCertificateRefs is nil.
						DomainNames:            nil,                // nil because CommonNameCheckEnabled is false.
						PersistenceProfile: &govcdtypes.NsxtAlbPoolPersistenceProfile{
							Name:  "persistence profile",
							Type:  "CLIENT_IP",
							Value: "",
						},
						MemberCount:        1,
						UpMemberCount:      1,
						HealthMessage:      "All members are up",
						VirtualServiceRefs: nil,
						SslEnabled:         utils.ToPTR(false),
					},
				}, nil)
			},
			expectedValue: []*PoolModel{
				{
					ID:                       poolID,
					Name:                     "pool1",
					Description:              "pool1 description",
					GatewayRef:               govcdtypes.OpenApiReference{ID: urnEdgeGateway},
					Enabled:                  utils.ToPTR(true),
					Algorithm:                PoolAlgorithmLeastConnections,
					DefaultPort:              utils.ToPTR(80),
					GracefulTimeoutPeriod:    utils.ToPTR(10),
					PassiveMonitoringEnabled: utils.ToPTR(true),
					HealthMonitors: []PoolModelHealthMonitor{
						{
							Name: "monitor HTTP",
							Type: PoolHealthMonitorTypeHTTP,
						},
						{
							Name: "monitor TCP",
							Type: PoolHealthMonitorTypeTCP,
						},
					},
					Members: []PoolModelMember{
						{
							Enabled:               true,
							IPAddress:             "192.168.0.1",
							Port:                  80,
							Ratio:                 utils.ToPTR(1),
							MarkedDownBy:          nil,
							HealthStatus:          "UP",
							DetailedHealthMessage: "",
						},
					},
					MemberGroupRef:         nil,
					CaCertificateRefs:      nil,
					CommonNameCheckEnabled: utils.ToPTR(false),
					DomainNames:            nil,
					PersistenceProfile: &PoolModelPersistenceProfile{
						Name:  "persistence profile",
						Type:  PoolPersistenceProfileTypeClientIP,
						Value: "",
					},
					MemberCount:        1,
					UpMemberCount:      1,
					HealthMessage:      "All members are up",
					VirtualServiceRefs: nil,
					SSLEnabled:         utils.ToPTR(false),
				},
				{
					ID:                       poolID2,
					Name:                     "pool2",
					Description:              "pool2 description",
					GatewayRef:               govcdtypes.OpenApiReference{ID: urnEdgeGateway},
					Enabled:                  utils.ToPTR(true),
					Algorithm:                PoolAlgorithmLeastConnections,
					DefaultPort:              utils.ToPTR(80),
					GracefulTimeoutPeriod:    utils.ToPTR(10),
					PassiveMonitoringEnabled: utils.ToPTR(true),
					HealthMonitors: []PoolModelHealthMonitor{
						{
							Name: "monitor HTTP",
							Type: PoolHealthMonitorTypeHTTP,
						},
						{
							Name: "monitor TCP",
							Type: PoolHealthMonitorTypeTCP,
						},
					},
					Members: []PoolModelMember{
						{
							Enabled:               true,
							IPAddress:             "192.168.0.2",
							Port:                  80,
							Ratio:                 utils.ToPTR(1),
							MarkedDownBy:          nil,
							HealthStatus:          "UP",
							DetailedHealthMessage: "",
						},
					},
					MemberGroupRef:         nil,
					CaCertificateRefs:      nil,
					CommonNameCheckEnabled: utils.ToPTR(false),
					DomainNames:            nil,
					PersistenceProfile: &PoolModelPersistenceProfile{
						Name:  "persistence profile",
						Type:  PoolPersistenceProfileTypeClientIP,
						Value: "",
					},
					MemberCount:        1,
					UpMemberCount:      1,
					HealthMessage:      "All members are up",
					VirtualServiceRefs: nil,
					SSLEnabled:         utils.ToPTR(false),
				},
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name:          "param-edgeGatewayID-empty",
			edgeGatewayID: "",
			mockFunc: func() {
			},
			expectedValue: []*PoolModel{},
			expectedErr:   true,
			err:           errors.New("edgeGatewayID is empty. Please provide a valid edgeGatewayID"),
		},
		{
			name:          "param-edgeGatewayID-invalid-id",
			edgeGatewayID: "1234",
			mockFunc: func() {
			},
			expectedValue: []*PoolModel{},
			expectedErr:   true,
			err:           errors.New("edgeGatewayID has invalid format. Please provide a valid edgeGatewayID"),
		},
		{
			name:          "refresh-error",
			edgeGatewayID: urnEdgeGateway,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(errors.New("error"))
			},
			expectedValue: nil,
			expectedErr:   true,
			err:           errors.New("error"),
		},
		{
			name:          "error-get-all-pools",
			edgeGatewayID: urnEdgeGateway,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAllAlbPoolSummaries(urnEdgeGateway, gomock.AssignableToTypeOf(url.Values{})).Return(nil, errors.New("error"))
			},
			expectedValue: []*PoolModel{},
			expectedErr:   true,
			err:           errors.New("error retrieving all ALB Pool summaries: error"),
		},
		{
			name:          "error-list-pool",
			edgeGatewayID: urnEdgeGateway,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil).Times(3)
				clientCAV.EXPECT().GetAllAlbPoolSummaries(urnEdgeGateway, gomock.AssignableToTypeOf(url.Values{})).Return([]*govcd.NsxtAlbPool{
					{
						NsxtAlbPool: &govcdtypes.NsxtAlbPool{
							ID:          poolID,
							Name:        "pool1",
							Description: "pool1 description",
							GatewayRef:  govcdtypes.OpenApiReference{ID: urnEdgeGateway},
						},
					},
					{
						NsxtAlbPool: &govcdtypes.NsxtAlbPool{
							ID:          poolID2,
							Name:        "pool2",
							Description: "pool2 description",
							GatewayRef:  govcdtypes.OpenApiReference{ID: urnEdgeGateway},
						},
					},
				}, nil)
				clientCAV.EXPECT().GetAlbPoolById(poolID).Return(&govcd.NsxtAlbPool{
					NsxtAlbPool: &govcdtypes.NsxtAlbPool{
						ID:                       poolID,
						Name:                     "pool1",
						Description:              "pool1 description",
						GatewayRef:               govcdtypes.OpenApiReference{ID: urnEdgeGateway},
						Enabled:                  utils.ToPTR(true),
						Algorithm:                "LEAST_CONNECTIONS",
						DefaultPort:              utils.ToPTR(80),
						GracefulTimeoutPeriod:    utils.ToPTR(10),
						PassiveMonitoringEnabled: utils.ToPTR(true),
						HealthMonitors: []govcdtypes.NsxtAlbPoolHealthMonitor{
							{
								Name: "monitor HTTP",
								Type: "HTTP",
							},
							{
								Name: "monitor TCP",
								Type: "TCP",
							},
						},
						Members: []govcdtypes.NsxtAlbPoolMember{
							{
								Enabled:               true,
								IpAddress:             "192.168.0.1",
								Port:                  80,
								Ratio:                 utils.ToPTR(1),
								MarkedDownBy:          nil,
								HealthStatus:          "UP",
								DetailedHealthMessage: "",
							},
						},
						MemberGroupRef:         nil, // Only one of MemberGroupRef or Members should be set.
						CaCertificateRefs:      nil,
						CommonNameCheckEnabled: utils.ToPTR(false), // false because CaCertificateRefs is nil.
						DomainNames:            nil,                // nil because CommonNameCheckEnabled is false.
						PersistenceProfile: &govcdtypes.NsxtAlbPoolPersistenceProfile{
							Name:  "persistence profile",
							Type:  "CLIENT_IP",
							Value: "",
						},
						MemberCount:        1,
						UpMemberCount:      1,
						HealthMessage:      "All members are up",
						VirtualServiceRefs: nil,
						SslEnabled:         utils.ToPTR(false),
					},
				}, nil)
				clientCAV.EXPECT().GetAlbPoolById(poolID2).Return(nil, errors.New("error"))
			},
			expectedValue: []*PoolModel{},
			expectedErr:   true,
			err:           errors.New("error retrieving complete ALB Pool: error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			pools, err := c.ListPools(context.Background(), tc.edgeGatewayID)
			if !tc.expectedErr {
				assert.NoError(t, err)
				assert.NotNil(t, pools)
			} else {
				assert.Error(t, err)
				assert.Nil(t, pools)
				return
			}

			assert.NoError(t, err)

			for i := range pools {
				for j := range tc.expectedValue {
					if pools[i].ID == tc.expectedValue[j].ID {
						assert.Equal(t, tc.expectedValue[j], pools[i])
					}
				}
			}
		})
	}
}

func TestClient_CreatePool(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock client for cloudavenue.
	clientCAV := NewMockclientFake(ctrl)

	c, _ := NewFakeClient(clientCAV)

	urnEdgeGateway := urn.Gateway.String() + uuid.New().String()
	urnPool := urn.LoadBalancerPool.String() + uuid.New().String()
	urnIPSet := urn.SecurityGroup.String() + uuid.New().String()

	tests := []struct {
		name          string
		mockFunc      func()
		expectedValue *PoolModel
		pool          PoolModelRequest
		expectedErr   bool
		err           error
	}{
		{
			name: "success",
			pool: PoolModelRequest{
				Name:                     "pool1",
				Description:              "pool1 description",
				EdgeGatewayID:            urnEdgeGateway,
				Enabled:                  utils.ToPTR(true),
				Algorithm:                "LEAST_CONNECTIONS",
				DefaultPort:              utils.ToPTR(80),
				GracefulTimeoutPeriod:    utils.ToPTR(10),
				PassiveMonitoringEnabled: utils.ToPTR(true),
				HealthMonitors: []PoolModelHealthMonitor{
					{
						Name: "monitor HTTP",
						Type: PoolHealthMonitorTypeHTTP,
					},
					{
						Name: "monitor TCP",
						Type: PoolHealthMonitorTypeTCP,
					},
				},
				Members: []PoolModelMember{
					{
						Enabled:               true,
						IPAddress:             "192.168.0.1",
						Port:                  80,
						Ratio:                 utils.ToPTR(1),
						MarkedDownBy:          nil,
						HealthStatus:          "UP",
						DetailedHealthMessage: "",
					},
				},
				MemberGroupRef:         nil, // Only one of MemberGroupRef or Members should be set.
				CaCertificateRefs:      nil,
				CommonNameCheckEnabled: utils.ToPTR(false), // false because CaCertificateRefs is nil.
				DomainNames:            nil,                // nil because CommonNameCheckEnabled is false.
				PersistenceProfile: &PoolModelPersistenceProfile{
					Name:  "persistence profile",
					Type:  PoolPersistenceProfileTypeClientIP,
					Value: "",
				},
			},
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().CreateNsxtAlbPool(gomock.AssignableToTypeOf(&govcdtypes.NsxtAlbPool{})).Return(&govcd.NsxtAlbPool{
					NsxtAlbPool: &govcdtypes.NsxtAlbPool{
						ID:                       urnPool,
						Name:                     "pool1",
						Description:              "pool1 description",
						GatewayRef:               govcdtypes.OpenApiReference{ID: urnEdgeGateway},
						Enabled:                  utils.ToPTR(true),
						Algorithm:                "LEAST_CONNECTIONS",
						DefaultPort:              utils.ToPTR(80),
						GracefulTimeoutPeriod:    utils.ToPTR(10),
						PassiveMonitoringEnabled: utils.ToPTR(true),
						HealthMonitors: []govcdtypes.NsxtAlbPoolHealthMonitor{
							{
								Name: "monitor HTTP",
								Type: "HTTP",
							},
							{
								Name: "monitor TCP",
								Type: "TCP",
							},
						},
						Members: []govcdtypes.NsxtAlbPoolMember{
							{
								Enabled:               true,
								IpAddress:             "192.168.0.1",
								Port:                  80,
								Ratio:                 utils.ToPTR(1),
								MarkedDownBy:          nil,
								HealthStatus:          "UP",
								DetailedHealthMessage: "",
							},
						},
						MemberGroupRef:         nil, // Only one of MemberGroupRef or Members should be set.
						CaCertificateRefs:      nil,
						CommonNameCheckEnabled: utils.ToPTR(false), // false because CaCertificateRefs is nil.
						DomainNames:            nil,                // nil because CommonNameCheckEnabled is false.
						PersistenceProfile: &govcdtypes.NsxtAlbPoolPersistenceProfile{
							Name:  "persistence profile",
							Type:  "CLIENT_IP",
							Value: "",
						},
						MemberCount:        1,
						UpMemberCount:      1,
						HealthMessage:      "All members are up",
						VirtualServiceRefs: nil,
						SslEnabled:         utils.ToPTR(false),
					},
				}, nil)
			},
			expectedValue: &PoolModel{
				ID:                       urnPool,
				Name:                     "pool1",
				Description:              "pool1 description",
				GatewayRef:               govcdtypes.OpenApiReference{ID: urnEdgeGateway},
				Enabled:                  utils.ToPTR(true),
				Algorithm:                "LEAST_CONNECTIONS",
				DefaultPort:              utils.ToPTR(80),
				GracefulTimeoutPeriod:    utils.ToPTR(10),
				PassiveMonitoringEnabled: utils.ToPTR(true),
				HealthMonitors: []PoolModelHealthMonitor{
					{
						Name: "monitor HTTP",
						Type: PoolHealthMonitorTypeHTTP,
					},
					{
						Name: "monitor TCP",
						Type: PoolHealthMonitorTypeTCP,
					},
				},
				Members: []PoolModelMember{
					{
						Enabled:               true,
						IPAddress:             "192.168.0.1",
						Port:                  80,
						Ratio:                 utils.ToPTR(1),
						MarkedDownBy:          nil,
						HealthStatus:          "UP",
						DetailedHealthMessage: "",
					},
				},
				MemberGroupRef:         nil, // Only one of MemberGroupRef or Members should be set.
				CaCertificateRefs:      nil,
				CommonNameCheckEnabled: utils.ToPTR(false), // false because CaCertificateRefs is nil.
				DomainNames:            nil,                // nil because CommonNameCheckEnabled is false.
				PersistenceProfile: &PoolModelPersistenceProfile{
					Name:  "persistence profile",
					Type:  PoolPersistenceProfileTypeClientIP,
					Value: "",
				},
				MemberCount:        1,
				UpMemberCount:      1,
				HealthMessage:      "All members are up",
				VirtualServiceRefs: nil,
				SSLEnabled:         utils.ToPTR(false),
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name: "success-member-ref-and-no-persistence-profile",
			pool: PoolModelRequest{
				Name:                     "pool1",
				Description:              "pool1 description",
				EdgeGatewayID:            urnEdgeGateway,
				Enabled:                  utils.ToPTR(true),
				Algorithm:                "LEAST_CONNECTIONS",
				DefaultPort:              utils.ToPTR(80),
				GracefulTimeoutPeriod:    utils.ToPTR(10),
				PassiveMonitoringEnabled: utils.ToPTR(true),
				HealthMonitors: []PoolModelHealthMonitor{
					{
						Name: "monitor HTTP",
						Type: PoolHealthMonitorTypeHTTP,
					},
					{
						Name: "monitor TCP",
						Type: PoolHealthMonitorTypeTCP,
					},
				},
				Members:                nil,
				MemberGroupRef:         &govcdtypes.OpenApiReference{ID: urnIPSet},
				CaCertificateRefs:      nil,
				CommonNameCheckEnabled: utils.ToPTR(false), // false because CaCertificateRefs is nil.
				DomainNames:            nil,                // nil because CommonNameCheckEnabled is false.
				PersistenceProfile:     nil,
			},
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().CreateNsxtAlbPool(gomock.AssignableToTypeOf(&govcdtypes.NsxtAlbPool{})).Return(&govcd.NsxtAlbPool{
					NsxtAlbPool: &govcdtypes.NsxtAlbPool{
						ID:                       urnPool,
						Name:                     "pool1",
						Description:              "pool1 description",
						GatewayRef:               govcdtypes.OpenApiReference{ID: urnEdgeGateway},
						Enabled:                  utils.ToPTR(true),
						Algorithm:                "LEAST_CONNECTIONS",
						DefaultPort:              utils.ToPTR(80),
						GracefulTimeoutPeriod:    utils.ToPTR(10),
						PassiveMonitoringEnabled: utils.ToPTR(true),
						HealthMonitors: []govcdtypes.NsxtAlbPoolHealthMonitor{
							{
								Name: "monitor HTTP",
								Type: "HTTP",
							},
							{
								Name: "monitor TCP",
								Type: "TCP",
							},
						},
						Members:                nil,
						MemberGroupRef:         &govcdtypes.OpenApiReference{ID: urnIPSet},
						CaCertificateRefs:      nil,
						CommonNameCheckEnabled: utils.ToPTR(false), // false because CaCertificateRefs is nil.
						DomainNames:            nil,                // nil because CommonNameCheckEnabled is false.
						PersistenceProfile: &govcdtypes.NsxtAlbPoolPersistenceProfile{
							Name:  "persistence profile",
							Type:  "CLIENT_IP",
							Value: "",
						},
						MemberCount:        1,
						UpMemberCount:      1,
						HealthMessage:      "All members are up",
						VirtualServiceRefs: nil,
						SslEnabled:         utils.ToPTR(false),
					},
				}, nil)
			},
			expectedValue: &PoolModel{
				ID:                       urnPool,
				Name:                     "pool1",
				Description:              "pool1 description",
				GatewayRef:               govcdtypes.OpenApiReference{ID: urnEdgeGateway},
				Enabled:                  utils.ToPTR(true),
				Algorithm:                "LEAST_CONNECTIONS",
				DefaultPort:              utils.ToPTR(80),
				GracefulTimeoutPeriod:    utils.ToPTR(10),
				PassiveMonitoringEnabled: utils.ToPTR(true),
				HealthMonitors: []PoolModelHealthMonitor{
					{
						Name: "monitor HTTP",
						Type: PoolHealthMonitorTypeHTTP,
					},
					{
						Name: "monitor TCP",
						Type: PoolHealthMonitorTypeTCP,
					},
				},
				Members:                nil,
				MemberGroupRef:         &govcdtypes.OpenApiReference{ID: urnIPSet},
				CaCertificateRefs:      nil,
				CommonNameCheckEnabled: utils.ToPTR(false), // false because CaCertificateRefs is nil.
				DomainNames:            nil,                // nil because CommonNameCheckEnabled is false.
				PersistenceProfile: &PoolModelPersistenceProfile{
					Name:  "persistence profile",
					Type:  PoolPersistenceProfileTypeClientIP,
					Value: "",
				},
				MemberCount:        1,
				UpMemberCount:      1,
				HealthMessage:      "All members are up",
				VirtualServiceRefs: nil,
				SSLEnabled:         utils.ToPTR(false),
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name: "bad-validation-algorithm",
			pool: PoolModelRequest{
				// Invalid field
				Algorithm: "LEAST CONNECTIONS", // LEAST CONNECTIONS instead of LEAST_CONNECTIONS

				// Valid fields
				Name:                     "pool1",
				Description:              "pool1 description",
				EdgeGatewayID:            urnEdgeGateway,
				Enabled:                  utils.ToPTR(true),
				DefaultPort:              utils.ToPTR(80),
				GracefulTimeoutPeriod:    utils.ToPTR(10),
				PassiveMonitoringEnabled: utils.ToPTR(true),
				HealthMonitors: []PoolModelHealthMonitor{
					{
						Name: "monitor HTTP",
						Type: PoolHealthMonitorTypeHTTP,
					},
					{
						Name: "monitor TCP",
						Type: PoolHealthMonitorTypeTCP,
					},
				},
				Members: []PoolModelMember{
					{
						Enabled:               true,
						IPAddress:             "192.168.0.1",
						Port:                  80,
						Ratio:                 utils.ToPTR(1),
						MarkedDownBy:          nil,
						HealthStatus:          "UP",
						DetailedHealthMessage: "",
					},
				},
				MemberGroupRef:         nil, // Only one of MemberGroupRef or Members should be set.
				CaCertificateRefs:      nil,
				CommonNameCheckEnabled: utils.ToPTR(false), // false because CaCertificateRefs is nil.
				DomainNames:            nil,                // nil because CommonNameCheckEnabled is false.
				PersistenceProfile: &PoolModelPersistenceProfile{
					Name:  "persistence profile",
					Type:  PoolPersistenceProfileTypeClientIP,
					Value: "",
				},
			},
			mockFunc:      func() {},
			expectedValue: nil,
			expectedErr:   true,
			err:           errors.New("Field validation for 'Algorithm' failed on the 'oneof'"),
		},
		{
			name: "bad-validation-Members",
			pool: PoolModelRequest{
				// Invalid field
				// Only one of MemberGroupRef or Members should be set.
				Members: []PoolModelMember{
					{
						Enabled:               true,
						IPAddress:             "192.168.0.1",
						Port:                  80,
						Ratio:                 utils.ToPTR(1),
						MarkedDownBy:          nil,
						HealthStatus:          "UP",
						DetailedHealthMessage: "",
					},
				},
				MemberGroupRef: &govcdtypes.OpenApiReference{ID: urnIPSet},

				// Valid fields
				Name:                     "pool1",
				Description:              "pool1 description",
				EdgeGatewayID:            urnEdgeGateway,
				Enabled:                  utils.ToPTR(true),
				DefaultPort:              utils.ToPTR(80),
				Algorithm:                "LEAST_CONNECTIONS", // LEAST CONNECTIONS instead of LEAST_CONNECTIONS
				GracefulTimeoutPeriod:    utils.ToPTR(10),
				PassiveMonitoringEnabled: utils.ToPTR(true),
				HealthMonitors: []PoolModelHealthMonitor{
					{
						Name: "monitor HTTP",
						Type: PoolHealthMonitorTypeHTTP,
					},
					{
						Name: "monitor TCP",
						Type: PoolHealthMonitorTypeTCP,
					},
				},
				CaCertificateRefs:      nil,
				CommonNameCheckEnabled: utils.ToPTR(false), // false because CaCertificateRefs is nil.
				DomainNames:            nil,                // nil because CommonNameCheckEnabled is false.
				PersistenceProfile: &PoolModelPersistenceProfile{
					Name:  "persistence profile",
					Type:  PoolPersistenceProfileTypeClientIP,
					Value: "",
				},
			},
			mockFunc:      func() {},
			expectedValue: nil,
			expectedErr:   true,
			err:           errors.New("Error:Field validation for 'Members' failed on the 'excluded_with'"),
		},
		{
			name: "error-create-pool",
			pool: PoolModelRequest{
				Name:                     "pool1",
				Description:              "pool1 description",
				EdgeGatewayID:            urnEdgeGateway,
				Enabled:                  utils.ToPTR(true),
				Algorithm:                "LEAST_CONNECTIONS",
				DefaultPort:              utils.ToPTR(80),
				GracefulTimeoutPeriod:    utils.ToPTR(10),
				PassiveMonitoringEnabled: utils.ToPTR(true),
				HealthMonitors: []PoolModelHealthMonitor{
					{
						Name: "monitor HTTP",
						Type: PoolHealthMonitorTypeHTTP,
					},
					{
						Name: "monitor TCP",
						Type: PoolHealthMonitorTypeTCP,
					},
				},
				Members: []PoolModelMember{
					{
						Enabled:               true,
						IPAddress:             "192.168.0.1",
						Port:                  80,
						Ratio:                 utils.ToPTR(1),
						MarkedDownBy:          nil,
						HealthStatus:          "UP",
						DetailedHealthMessage: "",
					},
				},
				MemberGroupRef:         nil, // Only one of MemberGroupRef or Members should be set.
				CaCertificateRefs:      nil,
				CommonNameCheckEnabled: utils.ToPTR(false), // false because CaCertificateRefs is nil.
				DomainNames:            nil,                // nil because CommonNameCheckEnabled is false.
				PersistenceProfile: &PoolModelPersistenceProfile{
					Name:  "persistence profile",
					Type:  PoolPersistenceProfileTypeClientIP,
					Value: "",
				},
			},
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().CreateNsxtAlbPool(gomock.AssignableToTypeOf(&govcdtypes.NsxtAlbPool{})).Return(nil, errors.New("error"))
			},
			expectedValue: nil,
			expectedErr:   true,
			err:           errors.New("error"),
		},
		{
			name: "refresh-error",
			pool: PoolModelRequest{
				Name:                     "pool1",
				Description:              "pool1 description",
				EdgeGatewayID:            urnEdgeGateway,
				Enabled:                  utils.ToPTR(true),
				Algorithm:                "LEAST_CONNECTIONS",
				DefaultPort:              utils.ToPTR(80),
				GracefulTimeoutPeriod:    utils.ToPTR(10),
				PassiveMonitoringEnabled: utils.ToPTR(true),
				HealthMonitors: []PoolModelHealthMonitor{
					{
						Name: "monitor HTTP",
						Type: PoolHealthMonitorTypeHTTP,
					},
					{
						Name: "monitor TCP",
						Type: PoolHealthMonitorTypeTCP,
					},
				},
				Members: []PoolModelMember{
					{
						Enabled:               true,
						IPAddress:             "192.168.0.1",
						Port:                  80,
						Ratio:                 utils.ToPTR(1),
						MarkedDownBy:          nil,
						HealthStatus:          "UP",
						DetailedHealthMessage: "",
					},
				},
				MemberGroupRef:         nil, // Only one of MemberGroupRef or Members should be set.
				CaCertificateRefs:      nil,
				CommonNameCheckEnabled: utils.ToPTR(false), // false because CaCertificateRefs is nil.
				DomainNames:            nil,                // nil because CommonNameCheckEnabled is false.
				PersistenceProfile: &PoolModelPersistenceProfile{
					Name:  "persistence profile",
					Type:  PoolPersistenceProfileTypeClientIP,
					Value: "",
				},
			},
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(errors.New("error"))
			},
			expectedErr: true,
			err:         errors.New("error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			pool, err := c.CreatePool(context.Background(), tc.pool)
			if !tc.expectedErr {
				assert.NoError(t, err)
				assert.NotNil(t, pool)
			} else {
				assert.Error(t, err)
				assert.Nil(t, pool)
				assert.Contains(t, err.Error(), tc.err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedValue, pool)
		})
	}
}

func TestClient_UpdatePool(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock client for cloudavenue.
	clientCAV := NewMockclientFake(ctrl)

	c, _ := NewFakeClient(clientCAV)

	urnEdgeGateway := urn.Gateway.String() + uuid.New().String()
	urnPool := urn.LoadBalancerPool.String() + uuid.New().String()

	tests := []struct {
		name          string
		mockFunc      func()
		expectedValue *PoolModel
		pool          PoolModelRequest
		poolID        string
		expectedErr   bool
		err           error
	}{
		{
			name:   "success",
			poolID: urnPool,
			pool: PoolModelRequest{
				Name:                     "poule 2",
				Description:              "poule description",
				EdgeGatewayID:            urnEdgeGateway,
				Enabled:                  utils.ToPTR(true),
				Algorithm:                "LEAST_CONNECTIONS",
				DefaultPort:              utils.ToPTR(80),
				GracefulTimeoutPeriod:    utils.ToPTR(10),
				PassiveMonitoringEnabled: utils.ToPTR(true),
				HealthMonitors: []PoolModelHealthMonitor{
					{
						Name: "monitor HTTP",
						Type: PoolHealthMonitorTypeHTTP,
					},
					{
						Name: "monitor TCP",
						Type: PoolHealthMonitorTypeTCP,
					},
				},
				Members: []PoolModelMember{
					{
						Enabled:               true,
						IPAddress:             "192.168.0.1",
						Port:                  80,
						Ratio:                 utils.ToPTR(1),
						MarkedDownBy:          nil,
						HealthStatus:          "UP",
						DetailedHealthMessage: "",
					},
				},
				MemberGroupRef:         nil, // Only one of MemberGroupRef or Members should be set.
				CaCertificateRefs:      nil,
				CommonNameCheckEnabled: utils.ToPTR(false), // false because CaCertificateRefs is nil.
				DomainNames:            nil,                // nil because CommonNameCheckEnabled is false.
				PersistenceProfile: &PoolModelPersistenceProfile{
					Name:  "persistence profile",
					Type:  PoolPersistenceProfileTypeClientIP,
					Value: "",
				},
			},
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbPoolById(urnPool).Return(&govcd.NsxtAlbPool{
					NsxtAlbPool: &govcdtypes.NsxtAlbPool{
						ID:                       urnPool,
						Name:                     "pool1",
						Description:              "pool1 description",
						GatewayRef:               govcdtypes.OpenApiReference{ID: urnEdgeGateway},
						Enabled:                  utils.ToPTR(true),
						Algorithm:                "LEAST_CONNECTIONS",
						DefaultPort:              utils.ToPTR(80),
						GracefulTimeoutPeriod:    utils.ToPTR(10),
						PassiveMonitoringEnabled: utils.ToPTR(true),
						HealthMonitors: []govcdtypes.NsxtAlbPoolHealthMonitor{
							{
								Name: "monitor HTTP",
								Type: "HTTP",
							},
							{
								Name: "monitor TCP",
								Type: "TCP",
							},
						},
						Members: []govcdtypes.NsxtAlbPoolMember{
							{
								Enabled:               true,
								IpAddress:             "192.168.0.1",
								Port:                  80,
								Ratio:                 utils.ToPTR(1),
								MarkedDownBy:          nil,
								HealthStatus:          "UP",
								DetailedHealthMessage: "",
							},
						},
						MemberGroupRef:         nil, // Only one of MemberGroupRef or Members should be set.
						CaCertificateRefs:      nil,
						CommonNameCheckEnabled: utils.ToPTR(false), // false because CaCertificateRefs is nil.
						DomainNames:            nil,                // nil because CommonNameCheckEnabled is false.
						PersistenceProfile:     nil,
						MemberCount:            1,
						UpMemberCount:          1,
						HealthMessage:          "All members are up",
						VirtualServiceRefs:     nil,
						SslEnabled:             utils.ToPTR(false),
					},
				}, nil)

				updatePool = func(_ fakePoolClient, _ *govcdtypes.NsxtAlbPool) (*govcd.NsxtAlbPool, error) {
					return &govcd.NsxtAlbPool{
						NsxtAlbPool: &govcdtypes.NsxtAlbPool{
							ID:                       urnPool,
							Name:                     "poule 2",
							Description:              "poule description",
							GatewayRef:               govcdtypes.OpenApiReference{ID: urnEdgeGateway},
							Enabled:                  utils.ToPTR(true),
							Algorithm:                "LEAST_CONNECTIONS",
							DefaultPort:              utils.ToPTR(80),
							GracefulTimeoutPeriod:    utils.ToPTR(10),
							PassiveMonitoringEnabled: utils.ToPTR(true),
							HealthMonitors: []govcdtypes.NsxtAlbPoolHealthMonitor{
								{
									Name: "monitor HTTP",
									Type: "HTTP",
								},
								{
									Name: "monitor TCP",
									Type: "TCP",
								},
							},
							Members: []govcdtypes.NsxtAlbPoolMember{
								{
									Enabled:               true,
									IpAddress:             "192.168.0.1",
									Port:                  80,
									Ratio:                 utils.ToPTR(1),
									MarkedDownBy:          nil,
									HealthStatus:          "UP",
									DetailedHealthMessage: "",
								},
							},
							MemberGroupRef:         nil, // Only one of MemberGroupRef or Members should be set.
							CaCertificateRefs:      nil,
							CommonNameCheckEnabled: utils.ToPTR(false), // false because CaCertificateRefs is nil.
							DomainNames:            nil,                // nil because CommonNameCheckEnabled is false.
							PersistenceProfile:     nil,
							MemberCount:            1,
							UpMemberCount:          1,
							HealthMessage:          "All members are up",
							VirtualServiceRefs:     nil,
							SslEnabled:             utils.ToPTR(false),
						},
					}, nil
				}
			},
			expectedValue: &PoolModel{
				ID:                       urnPool,
				Name:                     "poule 2",
				Description:              "poule description",
				GatewayRef:               govcdtypes.OpenApiReference{ID: urnEdgeGateway},
				Enabled:                  utils.ToPTR(true),
				Algorithm:                "LEAST_CONNECTIONS",
				DefaultPort:              utils.ToPTR(80),
				GracefulTimeoutPeriod:    utils.ToPTR(10),
				PassiveMonitoringEnabled: utils.ToPTR(true),
				HealthMonitors: []PoolModelHealthMonitor{
					{
						Name: "monitor HTTP",
						Type: PoolHealthMonitorTypeHTTP,
					},
					{
						Name: "monitor TCP",
						Type: PoolHealthMonitorTypeTCP,
					},
				},
				Members: []PoolModelMember{
					{
						Enabled:               true,
						IPAddress:             "192.168.0.1",
						Port:                  80,
						Ratio:                 utils.ToPTR(1),
						MarkedDownBy:          nil,
						HealthStatus:          "UP",
						DetailedHealthMessage: "",
					},
				},
				MemberGroupRef:         nil, // Only one of MemberGroupRef or Members should be set.
				CaCertificateRefs:      nil,
				CommonNameCheckEnabled: utils.ToPTR(false), // false because CaCertificateRefs is nil.
				DomainNames:            nil,                // nil because CommonNameCheckEnabled is false.
				PersistenceProfile: &PoolModelPersistenceProfile{
					Name:  "persistence profile",
					Type:  PoolPersistenceProfileTypeClientIP,
					Value: "",
				},
				MemberCount:        1,
				UpMemberCount:      1,
				HealthMessage:      "All members are up",
				VirtualServiceRefs: nil,
				SSLEnabled:         utils.ToPTR(false),
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name:   "error-empty-poolID",
			poolID: "",
			pool:   PoolModelRequest{},
			mockFunc: func() {
			},
			expectedValue: &PoolModel{},
			expectedErr:   true,
			err:           errors.New("poolID is empty. Please provide a valid poolID"),
		},
		{
			name:   "error-invalid-poolID",
			poolID: "1234",
			pool:   PoolModelRequest{},
			mockFunc: func() {
			},
			expectedValue: &PoolModel{},
			expectedErr:   true,
			err:           errors.New("poolID has invalid format. Please provide a valid poolID"),
		},
		{
			name:   "error-validation",
			poolID: urnPool,
			pool:   PoolModelRequest{},
			mockFunc: func() {
			},
			expectedValue: &PoolModel{},
			expectedErr:   true,
			err:           errors.New("Error:Field validation for 'Name' failed on the 'required'"),
		},
		{
			name:   "refresh-error",
			poolID: urnPool,
			pool: PoolModelRequest{
				Name:                     "pool1",
				Description:              "pool1 description",
				EdgeGatewayID:            urnEdgeGateway,
				Enabled:                  utils.ToPTR(true),
				Algorithm:                "LEAST_CONNECTIONS",
				DefaultPort:              utils.ToPTR(80),
				GracefulTimeoutPeriod:    utils.ToPTR(10),
				PassiveMonitoringEnabled: utils.ToPTR(true),
				HealthMonitors: []PoolModelHealthMonitor{
					{
						Name: "monitor HTTP",
						Type: PoolHealthMonitorTypeHTTP,
					},
					{
						Name: "monitor TCP",
						Type: PoolHealthMonitorTypeTCP,
					},
				},
				Members: []PoolModelMember{
					{
						Enabled:               true,
						IPAddress:             "192.168.0.1",
						Port:                  80,
						Ratio:                 utils.ToPTR(1),
						MarkedDownBy:          nil,
						HealthStatus:          "UP",
						DetailedHealthMessage: "",
					},
				},
				MemberGroupRef:         nil, // Only one of MemberGroupRef or Members should be set.
				CaCertificateRefs:      nil,
				CommonNameCheckEnabled: utils.ToPTR(false), // false because CaCertificateRefs is nil.
				DomainNames:            nil,                // nil because CommonNameCheckEnabled is false.
				PersistenceProfile: &PoolModelPersistenceProfile{
					Name:  "persistence profile",
					Type:  PoolPersistenceProfileTypeClientIP,
					Value: "",
				},
			},
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(errors.New("error"))
			},
			expectedErr: true,
			err:         errors.New("error"),
		},
		{
			name:   "get-pool-error",
			poolID: urnPool,
			pool: PoolModelRequest{
				Name:                     "poule 2",
				Description:              "poule description",
				EdgeGatewayID:            urnEdgeGateway,
				Enabled:                  utils.ToPTR(true),
				Algorithm:                "LEAST_CONNECTIONS",
				DefaultPort:              utils.ToPTR(80),
				GracefulTimeoutPeriod:    utils.ToPTR(10),
				PassiveMonitoringEnabled: utils.ToPTR(true),
				HealthMonitors: []PoolModelHealthMonitor{
					{
						Name: "monitor HTTP",
						Type: PoolHealthMonitorTypeHTTP,
					},
					{
						Name: "monitor TCP",
						Type: PoolHealthMonitorTypeTCP,
					},
				},
				Members: []PoolModelMember{
					{
						Enabled:               true,
						IPAddress:             "192.168.0.1",
						Port:                  80,
						Ratio:                 utils.ToPTR(1),
						MarkedDownBy:          nil,
						HealthStatus:          "UP",
						DetailedHealthMessage: "",
					},
				},
				MemberGroupRef:         nil, // Only one of MemberGroupRef or Members should be set.
				CaCertificateRefs:      nil,
				CommonNameCheckEnabled: utils.ToPTR(false), // false because CaCertificateRefs is nil.
				DomainNames:            nil,                // nil because CommonNameCheckEnabled is false.
				PersistenceProfile: &PoolModelPersistenceProfile{
					Name:  "persistence profile",
					Type:  PoolPersistenceProfileTypeClientIP,
					Value: "",
				},
			},
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbPoolById(urnPool).Return(nil, errors.New("error"))
			},
			expectedErr: true,
			err:         errors.New("error"),
		},
		{
			name:   "error-update-pool",
			poolID: urnPool,
			pool: PoolModelRequest{
				Name:                     "poule 2",
				Description:              "poule description",
				EdgeGatewayID:            urnEdgeGateway,
				Enabled:                  utils.ToPTR(true),
				Algorithm:                "LEAST_CONNECTIONS",
				DefaultPort:              utils.ToPTR(80),
				GracefulTimeoutPeriod:    utils.ToPTR(10),
				PassiveMonitoringEnabled: utils.ToPTR(true),
				HealthMonitors: []PoolModelHealthMonitor{
					{
						Name: "monitor HTTP",
						Type: PoolHealthMonitorTypeHTTP,
					},
					{
						Name: "monitor TCP",
						Type: PoolHealthMonitorTypeTCP,
					},
				},
				Members: []PoolModelMember{
					{
						Enabled:               true,
						IPAddress:             "192.168.0.1",
						Port:                  80,
						Ratio:                 utils.ToPTR(1),
						MarkedDownBy:          nil,
						HealthStatus:          "UP",
						DetailedHealthMessage: "",
					},
				},
				MemberGroupRef:         nil, // Only one of MemberGroupRef or Members should be set.
				CaCertificateRefs:      nil,
				CommonNameCheckEnabled: utils.ToPTR(false), // false because CaCertificateRefs is nil.
				DomainNames:            nil,                // nil because CommonNameCheckEnabled is false.
				PersistenceProfile: &PoolModelPersistenceProfile{
					Name:  "persistence profile",
					Type:  PoolPersistenceProfileTypeClientIP,
					Value: "",
				},
			},
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbPoolById(urnPool).Return(&govcd.NsxtAlbPool{
					NsxtAlbPool: &govcdtypes.NsxtAlbPool{
						ID:                       urnPool,
						Name:                     "pool1",
						Description:              "pool1 description",
						GatewayRef:               govcdtypes.OpenApiReference{ID: urnEdgeGateway},
						Enabled:                  utils.ToPTR(true),
						Algorithm:                "LEAST_CONNECTIONS",
						DefaultPort:              utils.ToPTR(80),
						GracefulTimeoutPeriod:    utils.ToPTR(10),
						PassiveMonitoringEnabled: utils.ToPTR(true),
						HealthMonitors: []govcdtypes.NsxtAlbPoolHealthMonitor{
							{
								Name: "monitor HTTP",
								Type: "HTTP",
							},
							{
								Name: "monitor TCP",
								Type: "TCP",
							},
						},
						Members: []govcdtypes.NsxtAlbPoolMember{
							{
								Enabled:               true,
								IpAddress:             "192.168.0.1",
								Port:                  80,
								Ratio:                 utils.ToPTR(1),
								MarkedDownBy:          nil,
								HealthStatus:          "UP",
								DetailedHealthMessage: "",
							},
						},
						MemberGroupRef:         nil, // Only one of MemberGroupRef or Members should be set.
						CaCertificateRefs:      nil,
						CommonNameCheckEnabled: utils.ToPTR(false), // false because CaCertificateRefs is nil.
						DomainNames:            nil,                // nil because CommonNameCheckEnabled is false.
						PersistenceProfile:     nil,
						MemberCount:            1,
						UpMemberCount:          1,
						HealthMessage:          "All members are up",
						VirtualServiceRefs:     nil,
						SslEnabled:             utils.ToPTR(false),
					},
				}, nil)

				updatePool = func(_ fakePoolClient, _ *govcdtypes.NsxtAlbPool) (*govcd.NsxtAlbPool, error) {
					return nil, errors.New("error")
				}
			},
			expectedErr: true,
			err:         errors.New("error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			pool, err := c.UpdatePool(context.Background(), tc.poolID, tc.pool)
			if !tc.expectedErr {
				assert.NoError(t, err)
				assert.NotNil(t, pool)
			} else {
				assert.Error(t, err)
				assert.Nil(t, pool)
				assert.Contains(t, err.Error(), tc.err.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestClient_DeletePool(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock client for cloudavenue.
	clientCAV := NewMockclientFake(ctrl)

	c, _ := NewFakeClient(clientCAV)

	urnPool := urn.LoadBalancerPool.String() + uuid.New().String()

	tests := []struct {
		name        string
		mockFunc    func()
		poolID      string
		expectedErr bool
		err         error
	}{
		{
			name:   "success",
			poolID: urnPool,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbPoolById(urnPool).Return(&govcd.NsxtAlbPool{
					NsxtAlbPool: &govcdtypes.NsxtAlbPool{
						ID: urnPool,
					},
				}, nil)

				deletePool = func(_ fakePoolClient) error {
					return nil
				}
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name:   "refresh-error",
			poolID: urnPool,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(errors.New("error"))
			},
			expectedErr: true,
			err:         errors.New("error"),
		},
		{
			name:   "error-get-pool",
			poolID: urnPool,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbPoolById(urnPool).Return(nil, errors.New("error"))
			},
			expectedErr: true,
			err:         errors.New("error retrieving Load Balancer Pool: error"),
		},
		{
			name:   "error-delete-pool",
			poolID: urnPool,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbPoolById(urnPool).Return(&govcd.NsxtAlbPool{
					NsxtAlbPool: &govcdtypes.NsxtAlbPool{
						ID: urnPool,
					},
				}, nil)
				deletePool = func(_ fakePoolClient) error {
					return errors.New("error")
				}
			},
			expectedErr: true,
			err:         errors.New("error"),
		},
		{
			name:        "param-poolID-empty",
			poolID:      "",
			mockFunc:    func() {},
			expectedErr: true,
			err:         errors.New("poolID is empty. Please provide a valid poolID"),
		},
		{
			name:   "param-poolID-invalid-id",
			poolID: "1234",
			mockFunc: func() {
			},
			expectedErr: true,
			err:         errors.New("poolID has invalid format. Please provide a valid poolID"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			err := c.DeletePool(context.Background(), tc.poolID)
			if !tc.expectedErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.err.Error())
			}
		})
	}
}
