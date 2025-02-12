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
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/helpers/validators"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

func TestVirtualServiceRequestValidation(t *testing.T) {
	poolID := urn.LoadBalancerPool.String() + uuid.New().String()
	edgeGatewayID := urn.Gateway.String() + uuid.New().String()
	serviceEngineID := urn.ServiceEngineGroup.String() + uuid.New().String()
	// certificateID := urn.CertificateLibraryItem.String() + uuid.New().String()

	tests := []struct {
		name         string
		virtualModel VirtualServiceModelRequest
		expectedErr  bool
		err          error
	}{
		{
			name: "success",
			virtualModel: VirtualServiceModelRequest{
				Name:                 "virtualServiceName1",
				Description:          "virtualServiceDescription1",
				Enabled:              utils.ToPTR(true),
				ApplicationProfile:   VirtualServiceModelApplicationProfile("HTTP"),
				PoolID:               poolID,
				EdgeGatewayID:        edgeGatewayID,
				ServiceEngineGroupID: &serviceEngineID,
				ServicePorts: []VirtualServiceModelServicePort{
					{
						Start: utils.ToPTR(80),
						End:   nil,
						Type:  VirtualServiceModelServicePortType("TCP_PROXY"),
					},
				},
				VirtualIPAddress: "192.168.0.1",
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name: "error-app-profile",
			virtualModel: VirtualServiceModelRequest{
				Name:                 "virtualServiceName1",
				Description:          "virtualServiceDescription1",
				Enabled:              utils.ToPTR(true),
				ApplicationProfile:   VirtualServiceModelApplicationProfile("HTTE"),
				PoolID:               poolID,
				EdgeGatewayID:        edgeGatewayID,
				ServiceEngineGroupID: &serviceEngineID,
				ServicePorts: []VirtualServiceModelServicePort{
					{
						Start: utils.ToPTR(80),
						End:   nil,
						Type:  VirtualServiceModelServicePortType("TCP_PROXY"),
					},
				},
				VirtualIPAddress: "192.168.0.1",
			},
			expectedErr: true,
			err:         errors.New("Field validation for 'ApplicationProfile'"),
		},
		{
			name: "error-invalid-pool-id",
			virtualModel: VirtualServiceModelRequest{
				Name:                 "virtualServiceName1",
				Description:          "virtualServiceDescription1",
				Enabled:              utils.ToPTR(true),
				ApplicationProfile:   VirtualServiceModelApplicationProfile("HTTP"),
				PoolID:               edgeGatewayID,
				EdgeGatewayID:        edgeGatewayID,
				ServiceEngineGroupID: &serviceEngineID,
				ServicePorts: []VirtualServiceModelServicePort{
					{
						Start: utils.ToPTR(80),
						End:   nil,
						Type:  VirtualServiceModelServicePortType("TCP_PROXY"),
					},
				},
				VirtualIPAddress: "192.168.0.1",
			},
			expectedErr: true,
			err:         errors.New("Field validation for 'PoolID' failed on the 'urn'"),
		},
		{
			name: "error-invalid-edgegateway-id",
			virtualModel: VirtualServiceModelRequest{
				Name:               "virtualServiceName1",
				Description:        "virtualServiceDescription1",
				Enabled:            utils.ToPTR(true),
				ApplicationProfile: VirtualServiceModelApplicationProfile("HTTP"),
				PoolID:             poolID,
				EdgeGatewayID:      poolID,
				ServicePorts: []VirtualServiceModelServicePort{
					{
						Start: utils.ToPTR(80),
						End:   nil,
						Type:  VirtualServiceModelServicePortType("TCP_PROXY"),
					},
				},
				VirtualIPAddress: "192.168.0.1",
			},
			expectedErr: true,
			err:         errors.New("Field validation for 'EdgeGatewayID' failed on the 'urn'"),
		},
		{
			name: "error-invalid-virtual-ip",
			virtualModel: VirtualServiceModelRequest{
				Name:               "virtualServiceName1",
				Description:        "virtualServiceDescription1",
				Enabled:            utils.ToPTR(true),
				ApplicationProfile: VirtualServiceModelApplicationProfile("HTTP"),
				PoolID:             poolID,
				EdgeGatewayID:      edgeGatewayID,
				ServicePorts: []VirtualServiceModelServicePort{
					{
						Start: utils.ToPTR(80),
						End:   nil,
						Type:  VirtualServiceModelServicePortType("TCP_PROXY"),
					},
				},
				VirtualIPAddress: "192.168.0.3001",
			},
			expectedErr: true,
			err:         errors.New("Field validation for 'VirtualIPAddress' failed on the 'ip4_addr'"),
		},
		{
			name: "error-no-service-ports",
			virtualModel: VirtualServiceModelRequest{
				Name:               "virtualServiceName1",
				Description:        "virtualServiceDescription1",
				Enabled:            utils.ToPTR(true),
				ApplicationProfile: VirtualServiceModelApplicationProfile("HTTP"),
				PoolID:             poolID,
				EdgeGatewayID:      edgeGatewayID,
				ServicePorts:       []VirtualServiceModelServicePort{},
				VirtualIPAddress:   "192.168.0.1",
			},
			expectedErr: true,
			err:         errors.New("Field validation for 'ServicePorts' failed on the 'gte'"),
		},
		{
			name: "error-service-port-start",
			virtualModel: VirtualServiceModelRequest{
				Name:               "virtualServiceName1",
				Description:        "virtualServiceDescription1",
				Enabled:            utils.ToPTR(true),
				ApplicationProfile: VirtualServiceModelApplicationProfile("HTTP"),
				PoolID:             poolID,
				EdgeGatewayID:      edgeGatewayID,
				ServicePorts: []VirtualServiceModelServicePort{
					{
						Start: utils.ToPTR(85000),
						End:   nil,
						Type:  VirtualServiceModelServicePortType("TCP_PROXY"),
					},
				},
				VirtualIPAddress: "192.168.0.1",
			},
			expectedErr: true,
			err:         errors.New("Field validation for 'Start' failed on the 'lte'"),
		},
		{
			name: "error-service-port-end-not-gt-start",
			virtualModel: VirtualServiceModelRequest{
				Name:               "virtualServiceName1",
				Description:        "virtualServiceDescription1",
				Enabled:            utils.ToPTR(true),
				ApplicationProfile: VirtualServiceModelApplicationProfile("HTTP"),
				PoolID:             poolID,
				EdgeGatewayID:      edgeGatewayID,
				ServicePorts: []VirtualServiceModelServicePort{
					{
						Start: utils.ToPTR(8080),
						End:   utils.ToPTR(8070),
						Type:  VirtualServiceModelServicePortType("TCP_PROXY"),
					},
				},
				VirtualIPAddress: "192.168.0.1",
			},
			expectedErr: true,
			err:         errors.New("Field validation for 'End' failed on the 'gtfield'"),
		},
		{
			name: "error-service-port-invalid-type",
			virtualModel: VirtualServiceModelRequest{
				Name:               "virtualServiceName1",
				Description:        "virtualServiceDescription1",
				Enabled:            utils.ToPTR(true),
				ApplicationProfile: VirtualServiceModelApplicationProfile("HTTP"),
				PoolID:             poolID,
				EdgeGatewayID:      edgeGatewayID,
				ServicePorts: []VirtualServiceModelServicePort{
					{
						Start: utils.ToPTR(8080),
						End:   utils.ToPTR(8090),
						Type:  VirtualServiceModelServicePortType("TCPPROXY"),
					},
				},
				VirtualIPAddress: "192.168.0.1",
			},
			expectedErr: true,
			err:         errors.New("Field validation for 'Type' failed on the 'oneof'"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validators.New().Struct(tc.virtualModel)
			if !tc.expectedErr {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.err.Error())
				return
			}
		})
	}
}

func TestClient_GetVirtualService(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock client for cloudavenue.
	clientCAV := NewMockclientFake(ctrl)

	c, _ := NewFakeClient(clientCAV)

	edgeGatewayID := urn.Gateway.String() + uuid.New().String()
	poolID := urn.LoadBalancerPool.String() + uuid.New().String()
	serviceEngineID := urn.ServiceEngineGroup.String() + uuid.New().String()
	virtualServiceID := urn.LoadBalancerVirtualService.String() + uuid.New().String()
	certificateID := urn.CertificateLibraryItem.String() + uuid.New().String()

	tests := []struct {
		name               string
		mockFunc           func()
		expectedValue      *VirtualServiceModel
		edgeGatewayID      string
		virtualServiceName string
		virtualServiceID   string
		byNameOrID         string
		expectedErr        bool
		err                error
	}{
		{
			name:               "success-http-by-name",
			edgeGatewayID:      edgeGatewayID,
			virtualServiceID:   virtualServiceID,
			virtualServiceName: "virtualServiceName1",
			byNameOrID:         "name",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbVirtualServiceByName(edgeGatewayID, "virtualServiceName1").Return(&govcd.NsxtAlbVirtualService{
					NsxtAlbVirtualService: &govcdtypes.NsxtAlbVirtualService{
						ID:          virtualServiceID,
						Name:        "virtualServiceName1",
						Description: "virtualServiceDescription1",
						Enabled:     utils.ToPTR(true),
						ApplicationProfile: govcdtypes.NsxtAlbVirtualServiceApplicationProfile{
							Type: "HTTP",
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
								PortStart:  utils.ToPTR(80),
								PortEnd:    nil,
								SslEnabled: utils.ToPTR(false),
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
			},
			expectedValue: &VirtualServiceModel{
				ID:                 virtualServiceID,
				Name:               "virtualServiceName1",
				Description:        "virtualServiceDescription1",
				Enabled:            utils.ToPTR(true),
				ApplicationProfile: VirtualServiceModelApplicationProfile("HTTP"),
				PoolRef: govcdtypes.OpenApiReference{
					ID: poolID,
				},
				ServiceEngineGroupRef: &govcdtypes.OpenApiReference{
					ID: serviceEngineID,
				},
				EdgeGatewayRef: govcdtypes.OpenApiReference{
					ID: edgeGatewayID,
				},
				ServicePorts: []VirtualServiceModelServicePort{
					{
						Start: utils.ToPTR(80),
						End:   nil,
						Type:  VirtualServiceModelServicePortType("TCP_PROXY"),
					},
				},
				VirtualIPAddress:      "192.168.0.1",
				HealthStatus:          VirtualServiceModelHealthStatus("UP"),
				HealthMessage:         "OK",
				DetailedHealthMessage: "OK",
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name:               "success-https-by-id",
			edgeGatewayID:      edgeGatewayID,
			virtualServiceID:   virtualServiceID,
			virtualServiceName: "virtualServiceName2",
			byNameOrID:         "id",
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
						CertificateRef: &govcdtypes.OpenApiReference{
							ID: certificateID,
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
			},
			expectedValue: &VirtualServiceModel{
				ID:                 virtualServiceID,
				Name:               "virtualServiceName2",
				Description:        "virtualServiceDescription2",
				Enabled:            utils.ToPTR(true),
				ApplicationProfile: VirtualServiceModelApplicationProfile("HTTPS"),
				PoolRef: govcdtypes.OpenApiReference{
					ID: poolID,
				},
				ServiceEngineGroupRef: &govcdtypes.OpenApiReference{
					ID: serviceEngineID,
				},
				EdgeGatewayRef: govcdtypes.OpenApiReference{
					ID: edgeGatewayID,
				},
				CertificateRef: &govcdtypes.OpenApiReference{
					ID: certificateID,
				},
				ServicePorts: []VirtualServiceModelServicePort{
					{
						Start: utils.ToPTR(443),
						End:   nil,
						Type:  VirtualServiceModelServicePortType("TCP_PROXY"),
					},
				},
				VirtualIPAddress:      "192.168.0.1",
				HealthStatus:          VirtualServiceModelHealthStatus("UP"),
				HealthMessage:         "OK",
				DetailedHealthMessage: "OK",
			},
		},
		{
			name:               "success-http-by-name-with-service-port-type-empty",
			edgeGatewayID:      edgeGatewayID,
			virtualServiceID:   virtualServiceID,
			virtualServiceName: "virtualServiceName1",
			byNameOrID:         "name",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbVirtualServiceByName(edgeGatewayID, "virtualServiceName1").Return(&govcd.NsxtAlbVirtualService{
					NsxtAlbVirtualService: &govcdtypes.NsxtAlbVirtualService{
						ID:          virtualServiceID,
						Name:        "virtualServiceName1",
						Description: "virtualServiceDescription1",
						Enabled:     utils.ToPTR(true),
						ApplicationProfile: govcdtypes.NsxtAlbVirtualServiceApplicationProfile{
							Type: "HTTP",
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
								PortStart:     utils.ToPTR(80),
								PortEnd:       nil,
								SslEnabled:    utils.ToPTR(false),
								TcpUdpProfile: nil,
							},
						},
						VirtualIpAddress:      "192.168.0.1",
						HealthStatus:          "UP",
						HealthMessage:         "OK",
						DetailedHealthMessage: "OK",
					},
				}, nil)
			},
			expectedValue: &VirtualServiceModel{
				ID:                 virtualServiceID,
				Name:               "virtualServiceName1",
				Description:        "virtualServiceDescription1",
				Enabled:            utils.ToPTR(true),
				ApplicationProfile: VirtualServiceModelApplicationProfile("HTTP"),
				PoolRef: govcdtypes.OpenApiReference{
					ID: poolID,
				},
				ServiceEngineGroupRef: &govcdtypes.OpenApiReference{
					ID: serviceEngineID,
				},
				EdgeGatewayRef: govcdtypes.OpenApiReference{
					ID: edgeGatewayID,
				},
				ServicePorts: []VirtualServiceModelServicePort{
					{
						Start: utils.ToPTR(80),
						End:   nil,
						Type:  "",
					},
				},
				VirtualIPAddress:      "192.168.0.1",
				HealthStatus:          VirtualServiceModelHealthStatus("UP"),
				HealthMessage:         "OK",
				DetailedHealthMessage: "OK",
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name:             "refresh-error",
			edgeGatewayID:    edgeGatewayID,
			virtualServiceID: virtualServiceID,
			byNameOrID:       "id",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(errors.New("error"))
			},
			expectedValue: nil,
			expectedErr:   true,
			err:           errors.New("error"),
		},
		{
			name:             "error-get",
			edgeGatewayID:    edgeGatewayID,
			virtualServiceID: virtualServiceID,
			byNameOrID:       "id",
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbVirtualServiceById(virtualServiceID).Return(nil, errors.New("error"))
			},
			expectedValue: &VirtualServiceModel{},
			expectedErr:   true,
			err:           errors.New("error"),
		},
		{
			name:             "param-edgeGatewayID-empty",
			edgeGatewayID:    "",
			virtualServiceID: virtualServiceID,
			byNameOrID:       "id",
			mockFunc: func() {
			},
			expectedValue: &VirtualServiceModel{},
			expectedErr:   true,
			err:           errors.New("edgeGatewayID is empty. Please provide a valid edgeGatewayID"),
		},
		{
			name:             "param-edgeGatewayID-invalid-id",
			edgeGatewayID:    "1234",
			virtualServiceID: virtualServiceID,
			byNameOrID:       "id",
			mockFunc: func() {
			},
			expectedValue: &VirtualServiceModel{},
			expectedErr:   true,
			err:           errors.New("edgeGatewayID has invalid format. Please provide a valid edgeGatewayID"),
		},
		{
			name:          "param-poolNameOrID-empty",
			edgeGatewayID: edgeGatewayID,
			byNameOrID:    "name",
			mockFunc: func() {
			},
			expectedValue: &VirtualServiceModel{},
			expectedErr:   true,
			err:           errors.New("virtualServiceNameOrID is empty. Please provide a valid poolNameOrID"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			nameOrID := tc.virtualServiceName
			if tc.byNameOrID == "id" {
				nameOrID = tc.virtualServiceID
			}

			vs, err := c.GetVirtualService(context.Background(), tc.edgeGatewayID, nameOrID)
			if !tc.expectedErr {
				assert.NoError(t, err)
				assert.NotNil(t, vs)
			} else {
				assert.Error(t, err)
				assert.Nil(t, vs)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedValue, vs)
		})
	}
}

func TestClient_ListVirtualServices(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock client for cloudavenue.
	clientCAV := NewMockclientFake(ctrl)

	c, _ := NewFakeClient(clientCAV)

	edgeGatewayID := urn.Gateway.String() + uuid.New().String()
	poolID := urn.LoadBalancerPool.String() + uuid.New().String()
	serviceEngineID := urn.ServiceEngineGroup.String() + uuid.New().String()
	virtualServiceID := urn.LoadBalancerVirtualService.String() + uuid.New().String()
	virtualServiceID2 := urn.LoadBalancerVirtualService.String() + uuid.New().String()
	certificateID := urn.CertificateLibraryItem.String() + uuid.New().String()

	tests := []struct {
		name          string
		mockFunc      func()
		expectedValue []*VirtualServiceModel
		edgeGatewayID string
		expectedErr   bool
		err           error
	}{
		{
			name:          "success",
			edgeGatewayID: edgeGatewayID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil).Times(3)
				clientCAV.EXPECT().GetAllAlbVirtualServiceSummaries(edgeGatewayID, nil).Return([]*govcd.NsxtAlbVirtualService{
					{
						NsxtAlbVirtualService: &govcdtypes.NsxtAlbVirtualService{
							ID:          virtualServiceID,
							Name:        "virtualServiceName1",
							Description: "virtualServiceDescription1",
							Enabled:     utils.ToPTR(true),
							GatewayRef: govcdtypes.OpenApiReference{
								ID: edgeGatewayID,
							},
						},
					},
					{
						NsxtAlbVirtualService: &govcdtypes.NsxtAlbVirtualService{
							ID:          virtualServiceID2,
							Name:        "virtualServiceName2",
							Description: "virtualServiceDescription2",
							Enabled:     utils.ToPTR(true),
							GatewayRef: govcdtypes.OpenApiReference{
								ID: edgeGatewayID,
							},
						},
					},
				}, nil)
				clientCAV.EXPECT().GetAlbVirtualServiceById(virtualServiceID).Return(&govcd.NsxtAlbVirtualService{
					NsxtAlbVirtualService: &govcdtypes.NsxtAlbVirtualService{
						ID:          virtualServiceID,
						Name:        "virtualServiceName1",
						Description: "virtualServiceDescription1",
						Enabled:     utils.ToPTR(true),
						ApplicationProfile: govcdtypes.NsxtAlbVirtualServiceApplicationProfile{
							Type: "HTTP",
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
								PortStart:  utils.ToPTR(80),
								PortEnd:    nil,
								SslEnabled: utils.ToPTR(false),
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
				clientCAV.EXPECT().GetAlbVirtualServiceById(virtualServiceID2).Return(&govcd.NsxtAlbVirtualService{
					NsxtAlbVirtualService: &govcdtypes.NsxtAlbVirtualService{
						ID:          virtualServiceID2,
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
						CertificateRef: &govcdtypes.OpenApiReference{
							ID: certificateID,
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
						VirtualIpAddress:      "192.168.1.1",
						HealthStatus:          "UP",
						HealthMessage:         "OK",
						DetailedHealthMessage: "OK",
					},
				}, nil)
			},
			expectedValue: []*VirtualServiceModel{
				{
					ID:                 virtualServiceID,
					Name:               "virtualServiceName1",
					Description:        "virtualServiceDescription1",
					Enabled:            utils.ToPTR(true),
					ApplicationProfile: VirtualServiceModelApplicationProfile("HTTP"),
					PoolRef: govcdtypes.OpenApiReference{
						ID: poolID,
					},
					ServiceEngineGroupRef: &govcdtypes.OpenApiReference{
						ID: serviceEngineID,
					},
					EdgeGatewayRef: govcdtypes.OpenApiReference{
						ID: edgeGatewayID,
					},
					ServicePorts: []VirtualServiceModelServicePort{
						{
							Start: utils.ToPTR(80),
							End:   nil,
							Type:  VirtualServiceModelServicePortType("TCP_PROXY"),
						},
					},
					VirtualIPAddress:      "192.168.0.1",
					HealthStatus:          VirtualServiceModelHealthStatus("UP"),
					HealthMessage:         "OK",
					DetailedHealthMessage: "OK",
				},
				{
					ID:                 virtualServiceID2,
					Name:               "virtualServiceName2",
					Description:        "virtualServiceDescription2",
					Enabled:            utils.ToPTR(true),
					ApplicationProfile: VirtualServiceModelApplicationProfile("HTTPS"),
					PoolRef: govcdtypes.OpenApiReference{
						ID: poolID,
					},
					ServiceEngineGroupRef: &govcdtypes.OpenApiReference{
						ID: serviceEngineID,
					},
					EdgeGatewayRef: govcdtypes.OpenApiReference{
						ID: edgeGatewayID,
					},
					CertificateRef: &govcdtypes.OpenApiReference{
						ID: certificateID,
					},
					ServicePorts: []VirtualServiceModelServicePort{
						{
							Start: utils.ToPTR(443),
							End:   nil,
							Type:  VirtualServiceModelServicePortType("TCP_PROXY"),
						},
					},
					VirtualIPAddress:      "192.168.1.1",
					HealthStatus:          VirtualServiceModelHealthStatus("UP"),
					HealthMessage:         "OK",
					DetailedHealthMessage: "OK",
				},
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name:          "error-get-virtual-service-by-id",
			edgeGatewayID: edgeGatewayID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil).Times(3)
				clientCAV.EXPECT().GetAllAlbVirtualServiceSummaries(edgeGatewayID, nil).Return([]*govcd.NsxtAlbVirtualService{
					{
						NsxtAlbVirtualService: &govcdtypes.NsxtAlbVirtualService{
							ID:          virtualServiceID,
							Name:        "virtualServiceName1",
							Description: "virtualServiceDescription1",
							Enabled:     utils.ToPTR(true),
							GatewayRef: govcdtypes.OpenApiReference{
								ID: edgeGatewayID,
							},
						},
					},
					{
						NsxtAlbVirtualService: &govcdtypes.NsxtAlbVirtualService{
							ID:          virtualServiceID2,
							Name:        "virtualServiceName2",
							Description: "virtualServiceDescription2",
							Enabled:     utils.ToPTR(true),
							GatewayRef: govcdtypes.OpenApiReference{
								ID: edgeGatewayID,
							},
						},
					},
				}, nil)
				clientCAV.EXPECT().GetAlbVirtualServiceById(virtualServiceID).Return(&govcd.NsxtAlbVirtualService{
					NsxtAlbVirtualService: &govcdtypes.NsxtAlbVirtualService{
						ID:          virtualServiceID,
						Name:        "virtualServiceName1",
						Description: "virtualServiceDescription1",
						Enabled:     utils.ToPTR(true),
						ApplicationProfile: govcdtypes.NsxtAlbVirtualServiceApplicationProfile{
							Type: "HTTP",
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
								PortStart:  utils.ToPTR(80),
								PortEnd:    nil,
								SslEnabled: utils.ToPTR(false),
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
				clientCAV.EXPECT().GetAlbVirtualServiceById(virtualServiceID2).Return(nil, errors.New("error"))
			},
			expectedValue: []*VirtualServiceModel{},
			expectedErr:   true,
			err:           errors.New("error retrieving complete virtual service: error"),
		},
		{
			name:          "refresh-error",
			edgeGatewayID: edgeGatewayID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(errors.New("error"))
			},
			expectedValue: nil,
			expectedErr:   true,
			err:           errors.New("error"),
		},
		{
			name:          "error-get",
			edgeGatewayID: edgeGatewayID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAllAlbVirtualServiceSummaries(edgeGatewayID, nil).Return(nil, errors.New("error"))
			},
			expectedValue: nil,
			expectedErr:   true,
			err:           errors.New("error"),
		},
		{
			name:          "param-edgeGatewayID-empty",
			edgeGatewayID: "",
			mockFunc: func() {
			},
			expectedValue: nil,
			expectedErr:   true,
			err:           errors.New("edgeGatewayID is empty. Please provide a valid edgeGatewayID"),
		},
		{
			name:          "param-edgeGatewayID-invalid-id",
			edgeGatewayID: "1234",
			mockFunc: func() {
			},
			expectedValue: nil,
			expectedErr:   true,
			err:           errors.New("edgeGatewayID has invalid format. Please provide a valid edgeGatewayID"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			vs, err := c.ListVirtualServices(context.Background(), tc.edgeGatewayID)
			if !tc.expectedErr {
				assert.NoError(t, err)
				assert.NotNil(t, vs)
			} else {
				assert.Error(t, err)
				assert.Nil(t, vs)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedValue, vs)
		})
	}
}

func TestClient_CreateVirtualService(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock client for cloudavenue.
	clientCAV := NewMockclientFake(ctrl)

	c, _ := NewFakeClient(clientCAV)

	edgeGatewayID := urn.Gateway.String() + uuid.New().String()
	poolID := urn.LoadBalancerPool.String() + uuid.New().String()
	serviceEngineID := urn.ServiceEngineGroup.String() + uuid.New().String()
	virtualServiceID := urn.LoadBalancerVirtualService.String() + uuid.New().String()

	tests := []struct {
		name          string
		mockFunc      func()
		virtualModel  VirtualServiceModelRequest
		expectedValue *VirtualServiceModel
		expectedErr   bool
		err           error
	}{
		{
			name: "success-http",
			virtualModel: VirtualServiceModelRequest{
				Name:                 "virtualServiceName1",
				Description:          "virtualServiceDescription1",
				Enabled:              utils.ToPTR(true),
				ApplicationProfile:   VirtualServiceModelApplicationProfile("HTTP"),
				PoolID:               poolID,
				EdgeGatewayID:        edgeGatewayID,
				ServiceEngineGroupID: &serviceEngineID,
				ServicePorts: []VirtualServiceModelServicePort{
					{
						Start: utils.ToPTR(80),
						End:   nil,
						Type:  VirtualServiceModelServicePortType("TCP_PROXY"),
					},
				},
				VirtualIPAddress: "192.168.0.1",
			},
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().CreateNsxtAlbVirtualService(gomock.AssignableToTypeOf(&govcdtypes.NsxtAlbVirtualService{})).Return(&govcd.NsxtAlbVirtualService{
					NsxtAlbVirtualService: &govcdtypes.NsxtAlbVirtualService{
						ID:          virtualServiceID,
						Name:        "virtualServiceName1",
						Description: "virtualServiceDescription1",
						Enabled:     utils.ToPTR(true),
						ApplicationProfile: govcdtypes.NsxtAlbVirtualServiceApplicationProfile{
							Type: "HTTP",
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
								PortStart:  utils.ToPTR(80),
								PortEnd:    nil,
								SslEnabled: utils.ToPTR(false),
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
			},
			expectedValue: &VirtualServiceModel{
				ID:                 virtualServiceID,
				Name:               "virtualServiceName1",
				Description:        "virtualServiceDescription1",
				Enabled:            utils.ToPTR(true),
				ApplicationProfile: VirtualServiceModelApplicationProfile("HTTP"),
				PoolRef: govcdtypes.OpenApiReference{
					ID: poolID,
				},
				ServiceEngineGroupRef: &govcdtypes.OpenApiReference{
					ID: serviceEngineID,
				},
				EdgeGatewayRef: govcdtypes.OpenApiReference{
					ID: edgeGatewayID,
				},
				ServicePorts: []VirtualServiceModelServicePort{
					{
						Start: utils.ToPTR(80),
						End:   nil,
						Type:  VirtualServiceModelServicePortType("TCP_PROXY"),
					},
				},
				VirtualIPAddress:      "192.168.0.1",
				HealthStatus:          VirtualServiceModelHealthStatus("UP"),
				HealthMessage:         "OK",
				DetailedHealthMessage: "OK",
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name:          "error-validation",
			virtualModel:  VirtualServiceModelRequest{},
			mockFunc:      func() {},
			expectedValue: nil,
			expectedErr:   true,
			err:           errors.New("Error:Field validation"),
		},
		{
			name: "error-create-virtual-service",
			virtualModel: VirtualServiceModelRequest{
				Name:                 "virtualServiceName1",
				Description:          "virtualServiceDescription1",
				Enabled:              utils.ToPTR(true),
				ApplicationProfile:   VirtualServiceModelApplicationProfile("HTTP"),
				PoolID:               poolID,
				EdgeGatewayID:        edgeGatewayID,
				ServiceEngineGroupID: &serviceEngineID,
				ServicePorts: []VirtualServiceModelServicePort{
					{
						Start: utils.ToPTR(80),
						End:   nil,
						Type:  VirtualServiceModelServicePortType("TCP_PROXY"),
					},
				},
				VirtualIPAddress: "192.168.0.1",
			},
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().CreateNsxtAlbVirtualService(gomock.AssignableToTypeOf(&govcdtypes.NsxtAlbVirtualService{})).Return(nil, errors.New("error"))
			},
			expectedValue: nil,
			expectedErr:   true,
			err:           errors.New("error"),
		},
		{
			name: "error-refresh",
			virtualModel: VirtualServiceModelRequest{
				Name:                 "virtualServiceName1",
				Description:          "virtualServiceDescription1",
				Enabled:              utils.ToPTR(true),
				ApplicationProfile:   VirtualServiceModelApplicationProfile("HTTP"),
				PoolID:               poolID,
				EdgeGatewayID:        edgeGatewayID,
				ServiceEngineGroupID: &serviceEngineID,
				ServicePorts: []VirtualServiceModelServicePort{
					{
						Start: utils.ToPTR(80),
						End:   nil,
						Type:  VirtualServiceModelServicePortType("TCP_PROXY"),
					},
				},
				VirtualIPAddress: "192.168.0.1",
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

			vs, err := c.CreateVirtualService(context.Background(), tc.virtualModel)
			if !tc.expectedErr {
				assert.NoError(t, err)
				assert.NotNil(t, vs)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.err.Error())
				assert.Nil(t, vs)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedValue, vs)
		})
	}
}

func TestClient_UpdateVirtualService(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock client for cloudavenue.
	clientCAV := NewMockclientFake(ctrl)

	c, _ := NewFakeClient(clientCAV)

	edgeGatewayID := urn.Gateway.String() + uuid.New().String()
	poolID := urn.LoadBalancerPool.String() + uuid.New().String()
	serviceEngineID := urn.ServiceEngineGroup.String() + uuid.New().String()
	virtualServiceID := urn.LoadBalancerVirtualService.String() + uuid.New().String()
	certificateID := urn.CertificateLibraryItem.String() + uuid.New().String()

	tests := []struct {
		name             string
		mockFunc         func()
		virtualServiceID string
		virtualModel     VirtualServiceModelRequest
		expectedValue    *VirtualServiceModel
		expectedErr      bool
		err              error
	}{
		{
			name:             "success-http",
			virtualServiceID: virtualServiceID,
			virtualModel: VirtualServiceModelRequest{
				Name:                 "virtualServiceName1",
				Description:          "virtualServiceDescription1",
				Enabled:              utils.ToPTR(true),
				ApplicationProfile:   VirtualServiceModelApplicationProfile("HTTP"),
				PoolID:               poolID,
				EdgeGatewayID:        edgeGatewayID,
				ServiceEngineGroupID: &serviceEngineID,
				ServicePorts: []VirtualServiceModelServicePort{
					{
						Start: utils.ToPTR(80),
						End:   nil,
						Type:  VirtualServiceModelServicePortType("TCP_PROXY"),
					},
				},
				VirtualIPAddress: "192.168.0.1",
			},
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbVirtualServiceById(virtualServiceID).Return(&govcd.NsxtAlbVirtualService{
					NsxtAlbVirtualService: &govcdtypes.NsxtAlbVirtualService{
						ID:          virtualServiceID,
						Name:        "virtualServiceName1",
						Description: "virtualServiceDescription1",
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
						CertificateRef: &govcdtypes.OpenApiReference{
							ID: certificateID,
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

				updateVirtualService = func(_ fakeVirtualServiceClient, _ *govcdtypes.NsxtAlbVirtualService) (*govcd.NsxtAlbVirtualService, error) {
					return &govcd.NsxtAlbVirtualService{
						NsxtAlbVirtualService: &govcdtypes.NsxtAlbVirtualService{
							ID:          virtualServiceID,
							Name:        "virtualServiceName1",
							Description: "virtualServiceDescription1",
							Enabled:     utils.ToPTR(true),
							ApplicationProfile: govcdtypes.NsxtAlbVirtualServiceApplicationProfile{
								Type: "HTTP",
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
									PortStart:  utils.ToPTR(80),
									PortEnd:    nil,
									SslEnabled: utils.ToPTR(false),
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
					}, nil
				}
			},
			expectedValue: &VirtualServiceModel{
				ID:                 virtualServiceID,
				Name:               "virtualServiceName1",
				Description:        "virtualServiceDescription1",
				Enabled:            utils.ToPTR(true),
				ApplicationProfile: VirtualServiceModelApplicationProfile("HTTP"),
				PoolRef: govcdtypes.OpenApiReference{
					ID: poolID,
				},
				ServiceEngineGroupRef: &govcdtypes.OpenApiReference{
					ID: serviceEngineID,
				},
				EdgeGatewayRef: govcdtypes.OpenApiReference{
					ID: edgeGatewayID,
				},
				ServicePorts: []VirtualServiceModelServicePort{
					{
						Start: utils.ToPTR(80),
						End:   nil,
						Type:  VirtualServiceModelServicePortType("TCP_PROXY"),
					},
				},
				VirtualIPAddress:      "192.168.0.1",
				HealthStatus:          VirtualServiceModelHealthStatus("UP"),
				HealthMessage:         "OK",
				DetailedHealthMessage: "OK",
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name:        "error-empty-virtualServiceID",
			expectedErr: true,
			mockFunc:    func() {},
			err:         errors.New("virtualServiceID is empty. Please provide a valid virtualServiceID"),
		},
		{
			name:             "error-invalid-virtualServiceID",
			expectedErr:      true,
			virtualServiceID: "1234",
			mockFunc:         func() {},
			err:              errors.New("virtualServiceID has invalid format. Please provide a valid virtualServiceID"),
		},
		{
			name:             "error-validation",
			expectedErr:      true,
			virtualServiceID: virtualServiceID,
			mockFunc:         func() {},
			err:              errors.New("Error:Field validation"),
		},
		{
			name:             "error-refresh",
			expectedErr:      true,
			virtualServiceID: virtualServiceID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(errors.New("error"))
			},
			virtualModel: VirtualServiceModelRequest{
				Name:                 "virtualServiceName1",
				Description:          "virtualServiceDescription1",
				Enabled:              utils.ToPTR(true),
				ApplicationProfile:   VirtualServiceModelApplicationProfile("HTTP"),
				PoolID:               poolID,
				EdgeGatewayID:        edgeGatewayID,
				ServiceEngineGroupID: &serviceEngineID,
				ServicePorts: []VirtualServiceModelServicePort{
					{
						Start: utils.ToPTR(80),
						End:   nil,
						Type:  VirtualServiceModelServicePortType("TCP_PROXY"),
					},
				},
				VirtualIPAddress: "192.168.0.1",
			},
			err: errors.New("error"),
		},
		{
			name: "error-get-virtual-service",
			virtualModel: VirtualServiceModelRequest{
				Name:                 "virtualServiceName1",
				Description:          "virtualServiceDescription1",
				Enabled:              utils.ToPTR(true),
				ApplicationProfile:   VirtualServiceModelApplicationProfile("HTTP"),
				PoolID:               poolID,
				EdgeGatewayID:        edgeGatewayID,
				ServiceEngineGroupID: &serviceEngineID,
				ServicePorts: []VirtualServiceModelServicePort{
					{
						Start: utils.ToPTR(80),
						End:   nil,
						Type:  VirtualServiceModelServicePortType("TCP_PROXY"),
					},
				},
				VirtualIPAddress: "192.168.0.1",
			},
			virtualServiceID: virtualServiceID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbVirtualServiceById(virtualServiceID).Return(nil, errors.New("error"))
			},
			expectedErr:   true,
			expectedValue: nil,
			err:           errors.New("error"),
		},
		{
			name: "error-update-virtual-service",
			virtualModel: VirtualServiceModelRequest{
				Name:                 "virtualServiceName1",
				Description:          "virtualServiceDescription1",
				Enabled:              utils.ToPTR(true),
				ApplicationProfile:   VirtualServiceModelApplicationProfile("HTTP"),
				PoolID:               poolID,
				EdgeGatewayID:        edgeGatewayID,
				ServiceEngineGroupID: &serviceEngineID,
				ServicePorts: []VirtualServiceModelServicePort{
					{
						Start: utils.ToPTR(80),
						End:   nil,
						Type:  VirtualServiceModelServicePortType("TCP_PROXY"),
					},
				},
				VirtualIPAddress: "192.168.0.1",
			},
			virtualServiceID: virtualServiceID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbVirtualServiceById(virtualServiceID).Return(nil, nil)
				updateVirtualService = func(_ fakeVirtualServiceClient, _ *govcdtypes.NsxtAlbVirtualService) (*govcd.NsxtAlbVirtualService, error) {
					return nil, errors.New("error")
				}
			},
			expectedErr:   true,
			expectedValue: nil,
			err:           errors.New("error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			vs, err := c.UpdateVirtualService(context.Background(), tc.virtualServiceID, tc.virtualModel)
			if !tc.expectedErr {
				assert.NoError(t, err)
				assert.NotNil(t, vs)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.err.Error())
				assert.Nil(t, vs)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedValue, vs)
		})
	}
}

func TestClient_DeleteVirtualService(t *testing.T) {
	// Mock controller.
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock client for cloudavenue.
	clientCAV := NewMockclientFake(ctrl)

	c, _ := NewFakeClient(clientCAV)

	virtualServiceID := urn.LoadBalancerVirtualService.String() + uuid.New().String()

	tests := []struct {
		name             string
		mockFunc         func()
		virtualServiceID string
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
						ID: virtualServiceID,
					},
				}, nil)
				deleteVirtualService = func(_ fakeVirtualServiceClient) error {
					return nil
				}
			},
			expectedErr: false,
			err:         nil,
		},
		{
			name:             "refresh-error",
			virtualServiceID: virtualServiceID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(errors.New("error"))
			},
			expectedErr: true,
			err:         errors.New("error"),
		},
		{
			name:             "error-get-virtual-service",
			virtualServiceID: virtualServiceID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbVirtualServiceById(virtualServiceID).Return(nil, errors.New("error"))
			},
			expectedErr: true,
			err:         errors.New("error"),
		},
		{
			name:             "error-delete-virtual-service",
			virtualServiceID: virtualServiceID,
			mockFunc: func() {
				clientCAV.EXPECT().Refresh().Return(nil)
				clientCAV.EXPECT().GetAlbVirtualServiceById(virtualServiceID).Return(&govcd.NsxtAlbVirtualService{
					NsxtAlbVirtualService: &govcdtypes.NsxtAlbVirtualService{
						ID: virtualServiceID,
					},
				}, nil)
				deleteVirtualService = func(_ fakeVirtualServiceClient) error {
					return errors.New("error")
				}
			},
			expectedErr: true,
			err:         errors.New("error"),
		},
		{
			name:        "error-empty-virtualServiceID",
			expectedErr: true,
			mockFunc:    func() {},
			err:         errors.New("virtualServiceID is empty. Please provide a valid virtualServiceID"),
		},
		{
			name:             "error-invalid-virtualServiceID",
			virtualServiceID: "1234",
			expectedErr:      true,
			mockFunc:         func() {},
			err:              errors.New("virtualServiceID has invalid format. Please provide a valid virtualServiceID"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()

			err := c.DeleteVirtualService(context.Background(), tc.virtualServiceID)
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
