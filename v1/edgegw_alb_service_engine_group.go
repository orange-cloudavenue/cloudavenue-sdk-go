package v1

import (
	"fmt"
	"net/url"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	clientcloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/clients/cloudavenue"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

// GetALLALBServiceEngineGroups return an array of ALB Service Engine Group assignment For an Edge Gateway
func (e *EdgeClient) ListALBServiceEngineGroups() ([]*EdgeGatewayALBServiceEngineGroupModel, error) {
	// Initialize the CloudAvenue client to call the CloudAvenue API or vmware API
	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	// Find the service engine group by name
	queryParams := url.Values{}
	queryParams.Add("filter", fmt.Sprintf("gatewayRef.id==%s", e.vcdEdge.EdgeGateway.ID)) // Filter by edge gateway ID URN
	govcdSEGs, err := c.Vmware.GetAllAlbServiceEngineGroupAssignments(queryParams)
	if err != nil {
		return nil, fmt.Errorf("error while fetching service engine group: %s", err.Error())
	}
	if len(govcdSEGs) == 0 {
		return nil, fmt.Errorf("no service engine group found for edge gateway %s", e.vcdEdge.EdgeGateway.Name)
	}

	// For x make it in []*EdgeGatewayALBServiceEngineGroup
	x := make([]*EdgeGatewayALBServiceEngineGroupModel, 0)
	for _, govcdSEG := range govcdSEGs {
		x = append(x, &EdgeGatewayALBServiceEngineGroupModel{
			ID:                         govcdSEG.NsxtAlbServiceEngineGroupAssignment.ServiceEngineGroupRef.ID,
			Name:                       govcdSEG.NsxtAlbServiceEngineGroupAssignment.ServiceEngineGroupRef.Name,
			GatewayRef:                 govcdSEG.NsxtAlbServiceEngineGroupAssignment.GatewayRef,
			MaxVirtualServices:         govcdSEG.NsxtAlbServiceEngineGroupAssignment.MaxVirtualServices,
			MinVirtualServices:         govcdSEG.NsxtAlbServiceEngineGroupAssignment.MinVirtualServices,
			NumDeployedVirtualServices: govcdSEG.NsxtAlbServiceEngineGroupAssignment.NumDeployedVirtualServices,
		})
	}

	return x, nil
}

// GetALBServiceEngineGroup return an ALB Service Engine Group For an Edge Gateway
// The nameOrID can be either the name or the ID of the service engine group
func (e *EdgeClient) GetALBServiceEngineGroup(nameOrID string) (*EdgeGatewayALBServiceEngineGroupModel, error) {
	var govcdSEG *govcd.NsxtAlbServiceEngineGroupAssignment

	// Initialize the CloudAvenue client to call the CloudAvenue API or vmware API
	c, err := clientcloudavenue.New()
	if err != nil {
		return nil, err
	}

	// Get the service engine group by name or ID
	if urn.IsServiceEngineGroup(nameOrID) {
		// Find the service engine group by ID
		queryParams := url.Values{}
		queryParams.Add("filter", fmt.Sprintf("gatewayRef.id==%s;serviceEngineGroupRef.id==%s", e.vcdEdge.EdgeGateway.ID, nameOrID)) // Filter by edge gateway ID URN
		govcdSEGs, err := c.Vmware.GetAllAlbServiceEngineGroupAssignments(queryParams)
		if err != nil {
			return nil, fmt.Errorf("error while fetching service engine group: %s", err.Error())
		}
		if len(govcdSEGs) == 0 {
			return nil, fmt.Errorf("no service engine group found for edge gateway %s", e.vcdEdge.EdgeGateway.Name)
		}
		if len(govcdSEGs) > 1 {
			return nil, fmt.Errorf("more than one service engine group found for edge gateway %s", e.vcdEdge.EdgeGateway.Name)
		}
		govcdSEG = govcdSEGs[0]
	} else {
		// Find the service engine group by name
		queryParams := url.Values{}
		queryParams.Add("filter", fmt.Sprintf("gatewayRef.id==%s", e.vcdEdge.EdgeGateway.ID)) // Filter by edge gateway ID URN
		govcdSEG, err = c.Vmware.GetFilteredAlbServiceEngineGroupAssignmentByName(nameOrID, queryParams)
		if err != nil {
			return nil, fmt.Errorf("error while fetching service engine group: %s", err.Error())
		}
	}

	// Set the EdgeGatewayALBServiceEngineGroupModel struct
	return &EdgeGatewayALBServiceEngineGroupModel{
		ID:                         govcdSEG.NsxtAlbServiceEngineGroupAssignment.ServiceEngineGroupRef.ID,
		Name:                       govcdSEG.NsxtAlbServiceEngineGroupAssignment.ServiceEngineGroupRef.Name,
		GatewayRef:                 govcdSEG.NsxtAlbServiceEngineGroupAssignment.GatewayRef,
		MaxVirtualServices:         govcdSEG.NsxtAlbServiceEngineGroupAssignment.MaxVirtualServices,
		MinVirtualServices:         govcdSEG.NsxtAlbServiceEngineGroupAssignment.MinVirtualServices,
		NumDeployedVirtualServices: govcdSEG.NsxtAlbServiceEngineGroupAssignment.NumDeployedVirtualServices,
	}, nil
}
