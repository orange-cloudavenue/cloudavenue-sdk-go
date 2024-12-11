package v1

import (
	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"
)

const (
	// Application Profile Types
	EdgeGatewayALBVirtualServiceModelApplicationProfileHTTP  EdgeGatewayALBVirtualServiceModelApplicationProfile = "HTTP"
	EdgeGatewayALBVirtualServiceModelApplicationProfileHTTPS EdgeGatewayALBVirtualServiceModelApplicationProfile = "HTTPS"
	EdgeGatewayALBVirtualServiceModelApplicationProfileL4    EdgeGatewayALBVirtualServiceModelApplicationProfile = "L4"
	EdgeGatewayALBVirtualServiceModelApplicationProfileL4TLS EdgeGatewayALBVirtualServiceModelApplicationProfile = "L4_TLS"
)

var EdgeGatewayALBVirtualServiceModelApplicationProfiles = []struct {
	Value       EdgeGatewayALBVirtualServiceModelApplicationProfile
	Description string
}{
	{EdgeGatewayALBVirtualServiceModelApplicationProfileHTTP, `If you choose "HTTP" you don't need to set the "port_type" and "ssl_enabled" attribute in "service_ports".`},
	{EdgeGatewayALBVirtualServiceModelApplicationProfileHTTPS, `If you choose "HTTPS", you must provide a certificate ID and you don't need to set the "port_type" attribute in "service_ports".`},
	{EdgeGatewayALBVirtualServiceModelApplicationProfileL4, `If you choose "L4", you can set a service "port_type" attribute in "service_ports.`},
	{EdgeGatewayALBVirtualServiceModelApplicationProfileL4TLS, `If you choose "L4_TLS", you must provide a certificate ID and you can set a service "port_type" attribute in "service_ports.`},
}

type (
	EdgeGatewayALBVirtualServiceModelApplicationProfile string

	// EdgeGatewayALBVirtualService represents a virtual service on an NSX-T Edge Gateway.
	// It contains the Edge Gateway Object, the Virtual Service Model and the NSX-T ALB Virtual Service.
	EdgeGatewayALBVirtualService struct {
		client         *EdgeClient
		VirtualService *EdgeGatewayALBVirtualServiceModel
		nsxtALBVS      *govcd.NsxtAlbVirtualService
	}

	EdgeGatewayALBVirtualServiceModel struct {
		ID string `json:"id,omitempty"`

		// Name contains meaningful name
		Name string `json:"name"`

		// Description is optional
		Description string `json:"description,omitempty"`

		// Enabled defines if the virtual service is enabled to accept traffic
		Enabled *bool `json:"enabled"`

		// ApplicationProfile sets protocol for load balancing by using NsxtAlbVirtualServiceApplicationProfile
		ApplicationProfile govcdtypes.NsxtAlbVirtualServiceApplicationProfile `json:"applicationProfile"`

		// LoadBalancerPoolRef contains Pool reference
		LoadBalancerPoolRef govcdtypes.OpenApiReference `json:"loadBalancerPoolRef"`

		// ServiceEngineGroupRef contains Service Engine Group reference to be used for the virtual service.
		// If not set and if more than one service engine group is assigned to the edge gateway: return an error.
		// If not set and if only one service engine group is assigned to the edge gateway it uses that service engine group.
		// If set it uses the provided service engine group.
		ServiceEngineGroupRef govcdtypes.OpenApiReference `json:"serviceEngineGroupRef,omitempty"`

		// CertificateRef contains certificate reference if serving encrypted traffic
		// If not set, the virtual service will not serve encrypted traffic (TLS/HTTPS).
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
