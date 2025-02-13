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
	"fmt"
	"net/url"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

func (c *client) ListServiceEngineGroups(ctx context.Context, edgeGatewayID string) ([]*ServiceEngineGroupModel, error) {
	if err := c.clientCloudavenue.Refresh(); err != nil {
		return nil, err
	}

	return c.listServiceEngineGroups(ctx, edgeGatewayID)
}

func (c *client) listServiceEngineGroups(_ context.Context, edgeGatewayID string) ([]*ServiceEngineGroupModel, error) {
	if edgeGatewayID == "" {
		return nil, fmt.Errorf("edgeGatewayID cannot be empty")
	}

	if !urn.IsEdgeGateway(edgeGatewayID) {
		return nil, fmt.Errorf("edgeGatewayID is not a valid URN")
	}

	// Find the service engine group by name
	queryParams := url.Values{}
	queryParams.Add("filter", fmt.Sprintf("gatewayRef.id==%s", edgeGatewayID)) // Filter

	segs, err := c.clientGoVCD.GetAllAlbServiceEngineGroupAssignments(queryParams)
	if err != nil {
		return nil, fmt.Errorf("error while fetching service engine group: %s", err.Error())
	}

	if len(segs) == 0 {
		return nil, fmt.Errorf("no service engine group found for edge gateway %s. The service Load Balancer might not be enabled on this edge gateway. Contact the support", edgeGatewayID)
	}

	// For x make it in []*ServiceEngineGroupsModel
	x := make([]*ServiceEngineGroupModel, len(segs))
	for i, seg := range segs {
		x[i] = &ServiceEngineGroupModel{
			ID:                         seg.NsxtAlbServiceEngineGroupAssignment.ServiceEngineGroupRef.ID,
			Name:                       seg.NsxtAlbServiceEngineGroupAssignment.ServiceEngineGroupRef.Name,
			GatewayRef:                 seg.NsxtAlbServiceEngineGroupAssignment.GatewayRef,
			MaxVirtualServices:         seg.NsxtAlbServiceEngineGroupAssignment.MaxVirtualServices,
			MinVirtualServices:         seg.NsxtAlbServiceEngineGroupAssignment.MinVirtualServices,
			NumDeployedVirtualServices: seg.NsxtAlbServiceEngineGroupAssignment.NumDeployedVirtualServices,
		}
	}

	return x, nil
}

// GetServiceEngineGroup return an Service Engine Group For an Edge Gateway
// The nameOrID can be either the name or the ID of the service engine group.
func (c *client) GetServiceEngineGroup(ctx context.Context, edgeGatewayID, nameOrID string) (*ServiceEngineGroupModel, error) {
	if err := c.clientCloudavenue.Refresh(); err != nil {
		return nil, err
	}

	segs, err := c.listServiceEngineGroups(ctx, edgeGatewayID)
	if err != nil {
		return nil, err
	}

	var seg *ServiceEngineGroupModel

	// Get the service engine group by name or ID
	if urn.IsServiceEngineGroup(nameOrID) {
		// Find the service engine group by ID
		for _, s := range segs {
			if s.ID == nameOrID {
				seg = s
			}
		}
	} else {
		// Find the service engine group by name
		for _, s := range segs {
			if s.Name == nameOrID {
				seg = s
			}
		}
	}

	if seg == nil {
		return nil, fmt.Errorf("the service engine group %s was not found for edge gateway %s", nameOrID, edgeGatewayID)
	}

	return seg, nil
}

// Retrieve the first service engine group for an edge gateway if one and only one is available.
func (c *client) GetFirstServiceEngineGroup(ctx context.Context, edgeGatewayID string) (*ServiceEngineGroupModel, error) {
	if err := c.clientCloudavenue.Refresh(); err != nil {
		return nil, err
	}

	segs, err := c.listServiceEngineGroups(ctx, edgeGatewayID)
	if err != nil {
		return nil, err
	}

	if len(segs) > 1 {
		return nil, fmt.Errorf("more than one service engine group available for edge gateway %s", edgeGatewayID)
	}

	return segs[0], nil
}
