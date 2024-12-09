package v1

import (
	"fmt"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

// GetAlbVirtualService gets an ALB (Advanced Load Balancer) Virtual Service.
// It returns an EdgeGatewayALBVirtualService instance containing the ALB Virtual Service model,
// or an error if the ALB Virtual Service is not found.
func (e *EdgeClient) GetALBVirtualService(nameOrID string) (*EdgeGatewayALBVirtualService, error) {
	// Check if the name or ID is empty
	if nameOrID == "" {
		return nil, fmt.Errorf("empty name")
	}

	// Initialize the CloudAvenue client to call the CloudAvenue API or vmware API
	var err error
	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	// Get the ALB Virtual Service by name or ID
	var nsxtALBVS *govcd.NsxtAlbVirtualService
	if !urn.IsLoadBalancerVirtualService(nameOrID) {
		nsxtALBVS, err = c.Vmware.GetAlbVirtualServiceByName(e.GetID(), nameOrID)
	} else {
		nsxtALBVS, err = c.Vmware.GetAlbVirtualServiceById(nameOrID)
	}
	if err != nil {
		return nil, err
	}

	// Set the ALB Virtual Service model returned by the CloudAvenue client
	vs := &EdgeGatewayALBVirtualServiceModel{
		ID:                    nsxtALBVS.NsxtAlbVirtualService.ID,
		Name:                  nsxtALBVS.NsxtAlbVirtualService.Name,
		Description:           nsxtALBVS.NsxtAlbVirtualService.Description,
		Enabled:               nsxtALBVS.NsxtAlbVirtualService.Enabled,
		ApplicationProfile:    nsxtALBVS.NsxtAlbVirtualService.ApplicationProfile,
		LoadBalancerPoolRef:   nsxtALBVS.NsxtAlbVirtualService.LoadBalancerPoolRef,
		ServiceEngineGroupRef: nsxtALBVS.NsxtAlbVirtualService.ServiceEngineGroupRef,
		CertificateRef:        nsxtALBVS.NsxtAlbVirtualService.CertificateRef,
		ServicePorts:          nsxtALBVS.NsxtAlbVirtualService.ServicePorts,
		VirtualIPAddress:      nsxtALBVS.NsxtAlbVirtualService.VirtualIpAddress,
		HealthStatus:          nsxtALBVS.NsxtAlbVirtualService.HealthStatus,
		HealthMessage:         nsxtALBVS.NsxtAlbVirtualService.HealthMessage,
		DetailedHealthMessage: nsxtALBVS.NsxtAlbVirtualService.DetailedHealthMessage,
	}

	return &EdgeGatewayALBVirtualService{
		client:         e,
		VirtualService: vs,
		nsxtALBVS:      nsxtALBVS,
	}, nil
}

// CreateAlbVirtualService creates an ALB (Advanced Load Balancer) Virtual Service.
// It returns an EdgeGatewayALBVirtualService instance containing the ALB Virtual Service model,
// or an error if the ALB Virtual Service is not created.
func (e *EdgeClient) CreateALBVirtualService(vs *EdgeGatewayALBVirtualServiceModel) (*EdgeGatewayALBVirtualService, error) {
	// Check if the ALB Virtual Service is empty
	if vs == nil {
		return nil, fmt.Errorf("empty virtual service")
	}

	// Initialize the CloudAvenue client to call the CloudAvenue API or vmware API
	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}
	// e.EdgeID = "urn:vcloud:gateway:" + e.GetID()

	// Create the ALB Virtual Service
	albNSXTVS, err := c.Vmware.CreateNsxtAlbVirtualService(&govcdtypes.NsxtAlbVirtualService{
		Name:                  vs.Name,
		Description:           vs.Description,
		ApplicationProfile:    vs.ApplicationProfile,
		Enabled:               vs.Enabled,
		GatewayRef:            govcdtypes.OpenApiReference{ID: e.GetID(), Name: e.GetName()},
		LoadBalancerPoolRef:   vs.LoadBalancerPoolRef,
		ServiceEngineGroupRef: vs.ServiceEngineGroupRef,
		CertificateRef:        vs.CertificateRef,
		ServicePorts:          vs.ServicePorts,
		VirtualIpAddress:      vs.VirtualIPAddress,
	})
	if err != nil {
		return nil, err
	}

	// Get the ALB Virtual Service
	newALBVS, err := e.GetALBVirtualService(albNSXTVS.NsxtAlbVirtualService.Name)
	if err != nil {
		return nil, err
	}

	return &EdgeGatewayALBVirtualService{
		client:         e,
		VirtualService: newALBVS.VirtualService,
		nsxtALBVS:      albNSXTVS,
	}, nil
}

// UpdateAlbVirtualService updates an ALB (Advanced Load Balancer) Virtual Service.
// It returns an EdgeGatewayALBVirtualService instance containing the ALB Virtual Service model,
// or an error if the ALB Virtual Service is not updated.
func (e *EdgeGatewayALBVirtualService) UpdateALBVirtualService(vs *EdgeGatewayALBVirtualServiceModel) (*EdgeGatewayALBVirtualService, error) {
	// Check if the ALB Virtual Service is empty
	if vs == nil {
		return nil, fmt.Errorf("empty virtual service")
	}

	// Get the actual ALB Virtual Service
	albVS, err := e.client.GetALBVirtualService(vs.Name)
	if err != nil {
		return nil, err
	}

	// Update the ALB Virtual Service
	_, err = e.nsxtALBVS.Update(albVS.nsxtALBVS.NsxtAlbVirtualService)
	if err != nil {
		return nil, err
	}

	return &EdgeGatewayALBVirtualService{
		client:         e.client,
		VirtualService: albVS.VirtualService,
		nsxtALBVS:      albVS.nsxtALBVS,
	}, nil
}

// DeleteAlbVirtualService deletes an ALB (Advanced Load Balancer) Virtual Service.
// It returns an error if the ALB Virtual Service is not deleted.
func (e *EdgeGatewayALBVirtualService) DeleteALBVirtualService() error {
	// Delete the ALB Virtual Service
	return e.nsxtALBVS.Delete()
}
