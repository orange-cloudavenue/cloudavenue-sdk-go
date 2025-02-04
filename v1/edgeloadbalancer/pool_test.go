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
		name              string
		mockFunc          func()
		expectedCertValue *PoolModel
		expectedErr       bool
		edgeGatewayID     string
		byNameOrID        string
		poolID            string
		poolName          string
		err               error
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
			expectedCertValue: &PoolModel{
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
			expectedCertValue: &PoolModel{
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
			expectedCertValue: nil,
			expectedErr:       true,
			err:               errors.New("error"),
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
			expectedCertValue: &PoolModel{},
			expectedErr:       true,
			err:               errors.New("error"),
		},
		{
			name:          "param-edgeGatewayID-empty",
			edgeGatewayID: "",
			poolID:        poolID,
			poolName:      "pool1",
			byNameOrID:    "id",
			mockFunc: func() {
			},
			expectedCertValue: &PoolModel{},
			expectedErr:       true,
			err:               errors.New("edgeGatewayID is empty. Please provide a valid edgeGatewayID"),
		},
		{
			name:          "param-edgeGatewayID-invalid-id",
			edgeGatewayID: "1234",
			poolID:        "1234",
			poolName:      "pool1",
			byNameOrID:    "id",
			mockFunc: func() {
			},
			expectedCertValue: &PoolModel{},
			expectedErr:       true,
			err:               errors.New("edgeGatewayID has invalid format. Please provide a valid edgeGatewayID"),
		},
		{
			name:          "param-poolNameOrID-empty",
			edgeGatewayID: urnEdgeGateway,
			poolID:        "",
			byNameOrID:    "name",
			mockFunc: func() {
			},
			expectedCertValue: &PoolModel{},
			expectedErr:       true,
			err:               errors.New("poolNameOrID is empty. Please provide a valid poolNameOrID"),
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
			assert.Equal(t, tc.expectedCertValue, pool)
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
		name              string
		mockFunc          func()
		expectedCertValue []*PoolModel
		expectedErr       bool
		edgeGatewayID     string
		err               error
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
			expectedCertValue: []*PoolModel{
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
			expectedCertValue: []*PoolModel{},
			expectedErr:       true,
			err:               errors.New("edgeGatewayID is empty. Please provide a valid edgeGatewayID"),
		},
		{
			name:          "param-edgeGatewayID-invalid-id",
			edgeGatewayID: "1234",
			mockFunc: func() {
			},
			expectedCertValue: []*PoolModel{},
			expectedErr:       true,
			err:               errors.New("edgeGatewayID has invalid format. Please provide a valid edgeGatewayID"),
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
			name:          "error-get-all-pools",
			edgeGatewayID: urnEdgeGateway,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAllAlbPoolSummaries(urnEdgeGateway, gomock.AssignableToTypeOf(url.Values{})).Return(nil, errors.New("error"))
			},
			expectedCertValue: []*PoolModel{},
			expectedErr:       true,
			err:               errors.New("error retrieving all ALB Pool summaries: error"),
		},
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
				clientCAV.EXPECT().GetAlbPoolById(poolID2).Return(nil, errors.New("error"))
			},
			expectedCertValue: []*PoolModel{},
			expectedErr:       true,
			err:               errors.New("error retrieving complete ALB Pool: error"),
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
				for j := range tc.expectedCertValue {
					if pools[i].ID == tc.expectedCertValue[j].ID {
						assert.Equal(t, tc.expectedCertValue[j], pools[i])
					}
				}
			}
		})
	}
}
