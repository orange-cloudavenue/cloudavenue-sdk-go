package v1

import (
	"fmt"
	"net/url"
	"time"

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
		return nil, fmt.Errorf("ALB Virtual Name or ID is empty, please provide a valid name or id")
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
		nsxtALBVS, err = c.Vmware.GetAlbVirtualServiceByName(e.vcdEdge.EdgeGateway.ID, nameOrID)
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

	// Add Service Engine Group if not provided
	if vs.ServiceEngineGroupRef.Name == "" {
		// Find the first service engine group
		queryParams := url.Values{}
		queryParams.Add("filter", fmt.Sprintf("gatewayRef.id==%s", e.vcdEdge.EdgeGateway.ID)) // Filter by edge gateway ID URN
		x, err := c.Vmware.GetAllAlbServiceEngineGroupAssignments(queryParams)
		if err != nil {
			return nil, fmt.Errorf("error while fetching service engine group: %s", err.Error())
		}
		if len(x) == 0 {
			return nil, fmt.Errorf("no service engine group found for edge gateway %s", e.EdgeName)
		}
		if len(x) > 1 {
			return nil, fmt.Errorf("multiple service engine group found for edge gateway %s, please precise which one to use", e.EdgeName)
		}
		vs.ServiceEngineGroupRef = govcdtypes.OpenApiReference{Name: x[0].NsxtAlbServiceEngineGroupAssignment.ServiceEngineGroupRef.Name, ID: x[0].NsxtAlbServiceEngineGroupAssignment.ServiceEngineGroupRef.ID}
	} else {
		// Find the service engine group by name
		queryParams := url.Values{}
		queryParams.Add("filter", fmt.Sprintf("gatewayRef.id==%s", e.vcdEdge.EdgeGateway.ID)) // Filter by edge gateway ID URN
		x, err := c.Vmware.GetAllAlbServiceEngineGroupAssignments(queryParams)
		if err != nil {
			return nil, fmt.Errorf("error while fetching service engine group: %s", err.Error())
		}
		if len(x) == 0 {
			return nil, fmt.Errorf("no service engine group found for edge gateway %s", e.EdgeName)
		}
		var found bool
		for _, seGroup := range x {
			if seGroup.NsxtAlbServiceEngineGroupAssignment.ServiceEngineGroupRef.Name == vs.ServiceEngineGroupRef.Name {
				vs.ServiceEngineGroupRef = govcdtypes.OpenApiReference{Name: seGroup.NsxtAlbServiceEngineGroupAssignment.ServiceEngineGroupRef.Name, ID: seGroup.NsxtAlbServiceEngineGroupAssignment.ServiceEngineGroupRef.ID}
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("service engine group %s not found for edge gateway %s", vs.ServiceEngineGroupRef.Name, e.EdgeName)
		}
	}

	// Create the ALB Virtual Service
	albNSXTVS, err := c.Vmware.CreateNsxtAlbVirtualService(&govcdtypes.NsxtAlbVirtualService{
		Name:                  vs.Name,
		Description:           vs.Description,
		ApplicationProfile:    vs.ApplicationProfile,
		Enabled:               vs.Enabled,
		GatewayRef:            govcdtypes.OpenApiReference{ID: e.vcdEdge.EdgeGateway.ID}, // Set the Edge Gateway ID URN
		LoadBalancerPoolRef:   vs.LoadBalancerPoolRef,
		ServiceEngineGroupRef: vs.ServiceEngineGroupRef,
		CertificateRef:        vs.CertificateRef,
		ServicePorts:          vs.ServicePorts,
		VirtualIpAddress:      vs.VirtualIPAddress,
	})
	if err != nil {
		return nil, fmt.Errorf("error to Create ALB Virtual Service: %s", err.Error())
	}

	// A workaround for the issue https://github.com/vmware/go-vcloud-director/issues/729
	// Wait for 10 seconds before to get the ALB Virtual Service to retrieve the created model
	time.Sleep(10 * time.Second)
	newALBVS, err := e.GetALBVirtualService(albNSXTVS.NsxtAlbVirtualService.Name)
	if err != nil {
		return nil, fmt.Errorf("error to get ALB Virtual Service: %s", err.Error())
	}

	return &EdgeGatewayALBVirtualService{
		client:         e,
		VirtualService: newALBVS.VirtualService,
		nsxtALBVS:      newALBVS.nsxtALBVS,
	}, nil
}

// Update updates an ALB (Advanced Load Balancer) Virtual Service.
// It returns an EdgeGatewayALBVirtualService instance containing the ALB Virtual Service model,
// or an error if the ALB Virtual Service is not updated.
func (e *EdgeGatewayALBVirtualService) Update(vs *EdgeGatewayALBVirtualServiceModel) (*EdgeGatewayALBVirtualService, error) {
	// Check if the ALB Virtual Service is empty
	if vs == nil {
		return nil, fmt.Errorf("empty virtual service")
	}

	// Initialize the CloudAvenue client to call the CloudAvenue API or vmware API
	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	// TODO: Move in a new Resource/Datasource to avoid code duplication
	// Check if the ALB Virtual Service Engine Group is empty
	if vs.ServiceEngineGroupRef.Name == "" {
		// Find the first service engine group
		queryParams := url.Values{}
		queryParams.Add("filter", fmt.Sprintf("gatewayRef.id==%s", e.client.vcdEdge.EdgeGateway.ID)) // Filter by edge gateway ID URN
		x, err := c.Vmware.GetAllAlbServiceEngineGroupAssignments(queryParams)
		if err != nil {
			return nil, fmt.Errorf("error while fetching service engine group: %s", err.Error())
		}
		if len(x) == 0 {
			return nil, fmt.Errorf("no service engine group found for edge gateway %s", e.client.EdgeName)
		}
		if len(x) > 1 {
			return nil, fmt.Errorf("multiple service engine group found for edge gateway %s, please precise which one to use", e.client.EdgeName)
		}
		vs.ServiceEngineGroupRef = govcdtypes.OpenApiReference{Name: x[0].NsxtAlbServiceEngineGroupAssignment.ServiceEngineGroupRef.Name, ID: x[0].NsxtAlbServiceEngineGroupAssignment.ServiceEngineGroupRef.ID}
	} else {
		// Find the service engine group by name
		queryParams := url.Values{}
		queryParams.Add("filter", fmt.Sprintf("gatewayRef.id==%s", e.client.vcdEdge.EdgeGateway.ID)) // Filter by edge gateway ID URN
		x, err := c.Vmware.GetAllAlbServiceEngineGroupAssignments(queryParams)
		if err != nil {
			return nil, fmt.Errorf("error while fetching service engine group: %s", err.Error())
		}
		if len(x) == 0 {
			return nil, fmt.Errorf("no service engine group found for edge gateway %s", e.client.EdgeName)
		}
		var found bool
		for _, seGroup := range x {
			if seGroup.NsxtAlbServiceEngineGroupAssignment.ServiceEngineGroupRef.Name == vs.ServiceEngineGroupRef.Name {
				vs.ServiceEngineGroupRef = govcdtypes.OpenApiReference{Name: seGroup.NsxtAlbServiceEngineGroupAssignment.ServiceEngineGroupRef.Name, ID: seGroup.NsxtAlbServiceEngineGroupAssignment.ServiceEngineGroupRef.ID}
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("service engine group %s not found for edge gateway %s", vs.ServiceEngineGroupRef.Name, e.client.EdgeName)
		}
	}

	// Update the ALB Virtual Service model
	newvs := &govcdtypes.NsxtAlbVirtualService{
		ID:                    vs.ID,
		Name:                  vs.Name,
		Description:           vs.Description,
		ApplicationProfile:    vs.ApplicationProfile,
		Enabled:               vs.Enabled,
		GatewayRef:            govcdtypes.OpenApiReference{ID: e.client.vcdEdge.EdgeGateway.ID},
		LoadBalancerPoolRef:   vs.LoadBalancerPoolRef,
		ServiceEngineGroupRef: vs.ServiceEngineGroupRef,
		CertificateRef:        vs.CertificateRef,
		ServicePorts:          vs.ServicePorts,
		VirtualIpAddress:      vs.VirtualIPAddress,
	}

	// Update the ALB Virtual Service
	albVS, err := e.nsxtALBVS.Update(newvs)
	if err != nil {
		return nil, err
	}

	// A workaround for the issue https://github.com/vmware/go-vcloud-director/issues/729
	// Wait for 10 seconds before to get the ALB Virtual Service to retrieve the updated model
	time.Sleep(10 * time.Second)
	newALBVS, err := e.client.GetALBVirtualService(albVS.NsxtAlbVirtualService.Name)
	if err != nil {
		return nil, err
	}

	return &EdgeGatewayALBVirtualService{
		client:         e.client,
		VirtualService: newALBVS.VirtualService,
		nsxtALBVS:      newALBVS.nsxtALBVS,
	}, nil
}

// Delete deletes an ALB (Advanced Load Balancer) Virtual Service.
// It returns an error if the ALB Virtual Service is not deleted.
func (e *EdgeGatewayALBVirtualService) Delete() error {
	// Delete the ALB Virtual Service
	return e.nsxtALBVS.Delete()
}
