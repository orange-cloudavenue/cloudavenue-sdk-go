package edgeloadbalancer

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/internal/utils"
)

type (
	fakeVirtualServiceClient interface {
		Update(*govcdtypes.NsxtAlbVirtualService) (*govcd.NsxtAlbVirtualService, error)
		Delete() error
	}

	VirtualServiceModel struct {
		ID          string
		Name        string
		Description string

		// Enabled defines if the virtual service is enabled to accept traffic
		Enabled *bool

		// ApplicationProfile sets protocol for load balancing
		ApplicationProfile VirtualServiceModelApplicationProfile

		// PoolRef contains Pool reference
		PoolRef govcdtypes.OpenApiReference

		// ServiceEngineGroupRef contains Service Engine Group reference to be used for the virtual service.
		// If not set and if more than one service engine group is assigned to the edge gateway: return an error.
		// If not set and if only one service engine group is assigned to the edge gateway it uses that service engine group.
		// If set it uses the provided service engine group.
		ServiceEngineGroupRef *govcdtypes.OpenApiReference

		// EdgeGatewayRef
		EdgeGatewayRef govcdtypes.OpenApiReference

		// CertificateRef contains certificate reference if serving encrypted traffic
		// If not set, the virtual service will not serve encrypted traffic (TLS/HTTPS).
		CertificateRef *govcdtypes.OpenApiReference

		// ServicePorts define one or more ports (or port ranges) of the virtual service
		ServicePorts []VirtualServiceModelServicePort

		// VirtualIpAddress to be used for exposing this virtual service
		VirtualIPAddress string

		// HealthStatus contains status of the Load Balancer Cloud. Possible values are:
		// VirtualServiceHealthStatusUP - The cloud is healthy and ready to enable Load Balancer for an Edge Gateway.
		// VirtualServiceHealthStatusDOWN - The cloud is in a failure state. Enabling Load balancer on an Edge Gateway may not be possible.
		// VirtualServiceHealthStatusRUNNING - The cloud is currently processing. An example is if it's enabling a Load Balancer for an Edge Gateway.
		// VirtualServiceHealthStatusUNAVAILABLE - The cloud is unavailable.
		// VirtualServiceHealthStatusUNKNOWN - The cloud state is unknown.
		HealthStatus VirtualServiceModelHealthStatus

		// HealthMessage shows a pool health status (e.g. "The pool is unassigned.")
		HealthMessage string

		// DetailedHealthMessage contains a more in depth health message
		DetailedHealthMessage string
	}

	VirtualServiceModelApplicationProfile string
	VirtualServiceModelServicePortType    string
	VirtualServiceModelHealthStatus       string

	// Type ServicePort represents a service port for the virtual service.
	VirtualServiceModelServicePort struct {
		// To make a range of ports, set the first value in Start and end value in End.
		// To make a single port, set the same value in PortStart and PortEnd.
		Start *int `validate:"required,gte=1,lte=65535"`
		End   *int `validate:"omitempty,gte=1,lte=65535,gtfield=Start"`

		// PortType can be set to
		// VirtualServiceServicePortTypeTCPProxy - "TCP_PROXY"
		// VirtualServiceServicePortTypeTCPFastPath - "TCP_FAST_PATH"
		// VirtualServiceServicePortTypeUDPFastPath - "UDP_FAST_PATH".
		Type VirtualServiceModelServicePortType `validate:"required,oneof=TCP_PROXY TCP_FAST_PATH UDP_FAST_PATH"`
	}

	VirtualServiceModelRequest struct {
		Name        string `validate:"required"`
		Description string `validate:"omitempty"`

		// Enabled defines if the virtual service is enabled to accept traffic
		Enabled *bool `validate:"required"`

		// ApplicationProfile sets protocol for load balancing
		ApplicationProfile VirtualServiceModelApplicationProfile `validate:"required,oneof=HTTP HTTPS L4 L4_TLS"`

		// PoolID contains a reference to the ELB Pool to be used for the virtual service
		PoolID string `validate:"required,urn_rfc2141,urn=loadBalancerPool"`

		// ServiceEngineGroupID contains Service Engine Group reference to be used for the virtual service.
		// If not set and if more than one service engine group is assigned to the edge gateway: return an error.
		// If not set and if only one service engine group is assigned to the edge gateway it uses that service engine group.
		// If set it uses the provided service engine group.
		ServiceEngineGroupID *string `validate:"omitempty,urn_rfc2141,urn=serviceEngineGroup"`

		// EdgeGatewayID contains a reference to the Edge Gateway where the virtual service will be created
		EdgeGatewayID string `validate:"required,urn_rfc2141,urn=gateway"`

		// CertificateRef contains certificate reference if serving encrypted traffic
		// If not set, the virtual service will not serve encrypted traffic (TLS/HTTPS).
		CertificateRef *string `validate:"omitempty,urn_rfc2141,urn=certificateLibraryItem,required_if=ApplicationProfile HTTPS,required_if=ApplicationProfile L4_TLS"`

		// ServicePorts define one or more ports (or port ranges) of the virtual service
		ServicePorts []VirtualServiceModelServicePort `validate:"required,gte=1,dive"`

		// VirtualIpAddress to be used for exposing this virtual service
		VirtualIPAddress string `validate:"required,ip4_addr"`
	}
)

func fromVCDNsxtAlbVirtualServiceToModel(vs govcdtypes.NsxtAlbVirtualService) *VirtualServiceModel {
	return &VirtualServiceModel{
		ID:                    vs.ID,
		Name:                  vs.Name,
		Description:           vs.Description,
		Enabled:               vs.Enabled,
		ApplicationProfile:    VirtualServiceModelApplicationProfile(vs.ApplicationProfile.Type),
		PoolRef:               vs.LoadBalancerPoolRef,
		EdgeGatewayRef:        vs.GatewayRef,
		ServiceEngineGroupRef: &vs.ServiceEngineGroupRef,
		CertificateRef:        vs.CertificateRef,
		ServicePorts:          fromVCDNsxtAlbVirtualServiceServicePortToModel(vs.ServicePorts),
		VirtualIPAddress:      vs.VirtualIpAddress,
		HealthStatus:          VirtualServiceModelHealthStatus(vs.HealthStatus),
		HealthMessage:         vs.HealthMessage,
		DetailedHealthMessage: vs.DetailedHealthMessage,
	}
}

func fromVCDNsxtAlbVirtualServiceServicePortToModel(sp []govcdtypes.NsxtAlbVirtualServicePort) []VirtualServiceModelServicePort {
	var servicePorts []VirtualServiceModelServicePort
	for _, port := range sp {
		servicePorts = append(servicePorts, VirtualServiceModelServicePort{
			Start: port.PortStart,
			End:   port.PortEnd,
			Type: func() VirtualServiceModelServicePortType {
				if port.TcpUdpProfile == nil {
					return ""
				}
				return VirtualServiceModelServicePortType(port.TcpUdpProfile.Type)
			}(),
		})
	}
	return servicePorts
}

func fromModelRequestToVCDNsxtAlbVirtualService(vs VirtualServiceModelRequest) *govcdtypes.NsxtAlbVirtualService {
	return &govcdtypes.NsxtAlbVirtualService{
		Name:        vs.Name,
		Description: vs.Description,
		Enabled:     vs.Enabled,
		ApplicationProfile: govcdtypes.NsxtAlbVirtualServiceApplicationProfile{
			Type: string(vs.ApplicationProfile),
		},
		GatewayRef: govcdtypes.OpenApiReference{
			ID: vs.EdgeGatewayID,
		},
		LoadBalancerPoolRef: govcdtypes.OpenApiReference{
			ID: vs.PoolID,
		},
		ServiceEngineGroupRef: func() govcdtypes.OpenApiReference {
			if vs.ServiceEngineGroupID == nil {
				return govcdtypes.OpenApiReference{}
			}
			return govcdtypes.OpenApiReference{
				ID: *vs.ServiceEngineGroupID,
			}
		}(),
		CertificateRef: func() *govcdtypes.OpenApiReference {
			if vs.CertificateRef == nil {
				return nil
			}
			return &govcdtypes.OpenApiReference{
				ID: *vs.CertificateRef,
			}
		}(),
		ServicePorts:     fromModelRequestServicePortToVCDNsxtAlbVirtualServiceServicePort(vs.ApplicationProfile, vs.ServicePorts),
		VirtualIpAddress: vs.VirtualIPAddress,
	}
}

func fromModelRequestServicePortToVCDNsxtAlbVirtualServiceServicePort(appProfile VirtualServiceModelApplicationProfile, sp []VirtualServiceModelServicePort) []govcdtypes.NsxtAlbVirtualServicePort {
	var servicePorts []govcdtypes.NsxtAlbVirtualServicePort
	for _, port := range sp {
		servicePorts = append(servicePorts, govcdtypes.NsxtAlbVirtualServicePort{
			PortStart: port.Start,
			PortEnd:   port.End,
			SslEnabled: func() *bool {
				if appProfile == VirtualServiceApplicationProfileHTTPS || appProfile == VirtualServiceApplicationProfileL4TLS {
					return utils.ToPTR(true)
				}
				return nil
			}(),
			TcpUdpProfile: &govcdtypes.NsxtAlbVirtualServicePortTcpUdpProfile{
				Type: string(port.Type),
			},
		})
	}
	return servicePorts
}
