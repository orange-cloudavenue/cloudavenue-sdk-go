package v1

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"
)

type (

	// EdgeGatewayALBVirtualService represents a virtual service on an NSX-T Edge Gateway
	EdgeGatewayALBVirtualService struct {
		client         *EdgeClient
		VirtualService *EdgeGatewayALBVirtualServiceModel
		nsxtALBVS      *govcd.NsxtAlbVirtualService
	}

	EdgeGatewayALBVirtualServiceModel struct {
		ID string `json:"id,omitempty"`

		// Name contains meaningful name
		Name string `json:"name,omitempty"`

		// Description is optional
		Description string `json:"description,omitempty"`

		// Enabled defines if the virtual service is enabled to accept traffic
		Enabled *bool `json:"enabled"`

		// ApplicationProfile sets protocol for load balancing by using NsxtAlbVirtualServiceApplicationProfile
		ApplicationProfile govcdtypes.NsxtAlbVirtualServiceApplicationProfile `json:"applicationProfile"`

		// GatewayRef contains NSX-T Edge Gateway reference
		// GatewayRef govcdtypes.OpenApiReference `json:"gatewayRef"`
		// LoadBalancerPoolRef contains Pool reference
		LoadBalancerPoolRef govcdtypes.OpenApiReference `json:"loadBalancerPoolRef"`
		// ServiceEngineGroupRef points to service engine group (which must be assigned to NSX-T Edge Gateway)
		ServiceEngineGroupRef govcdtypes.OpenApiReference `json:"serviceEngineGroupRef"`

		// CertificateRef contains certificate reference if serving encrypted traffic
		CertificateRef *govcdtypes.OpenApiReference `json:"certificateRef,omitempty"`

		// ServicePorts define one or more ports (or port ranges) of the virtual service
		ServicePorts []govcdtypes.NsxtAlbVirtualServicePort `json:"servicePorts"`

		// VirtualIpAddress to be used for exposing this virtual service
		VirtualIPAddress string `json:"virtualIpAddress"`

		// HealthStatus contains status of the Load Balancer Cloud. Possible values are:
		// UP - The cloud is healthy and ready to enable Load Balancer for an Edge Gateway.
		// DOWN - The cloud is in a failure state. Enabling Load balancer on an Edge Gateway may not be possible.
		// RUNNING - The cloud is currently processing. An example is if it's enabling a Load Balancer for an Edge Gateway.
		// UNAVAILABLE - The cloud is unavailable.
		// UNKNOWN - The cloud state is unknown.
		HealthStatus string `json:"healthStatus,omitempty"`

		// HealthMessage shows a pool health status (e.g. "The pool is unassigned.")
		HealthMessage string `json:"healthMessage,omitempty"`

		// DetailedHealthMessage contains a more in depth health message
		DetailedHealthMessage string `json:"detailedHealthMessage,omitempty"`
	}
)

// ! VirtualService
