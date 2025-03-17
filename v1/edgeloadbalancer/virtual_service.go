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

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/helpers/validators"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

// ListVirtualServices retrieves a list of virtual services for a given edge gateway.
// It returns a slice of VirtualServiceModel pointers and an error if any occurs during the process.
//
// Parameters:
//   - ctx: The context for the request.
//   - edgeGatewayID: The ID of the edge gateway for which to list virtual services.
func (c *client) ListVirtualServices(ctx context.Context, edgeGatewayID string) ([]*VirtualServiceModel, error) {
	if edgeGatewayID == "" {
		return nil, fmt.Errorf("edgeGatewayID is %w. Please provide a valid edgeGatewayID", errors.ErrEmpty)
	}

	if !urn.IsEdgeGateway(edgeGatewayID) {
		return nil, fmt.Errorf("edgeGatewayID has %w. Please provide a valid edgeGatewayID", errors.ErrInvalidFormat)
	}

	if err := c.clientCloudavenue.Refresh(); err != nil {
		return nil, err
	}

	avs, err := c.clientGoVCD.GetAllAlbVirtualServiceSummaries(edgeGatewayID, nil)
	if err != nil {
		return nil, fmt.Errorf("error retrieving list of ELB Virtual Services: %w", err)
	}

	allVirtualServices := make([]*VirtualServiceModel, len(avs))
	for index := range avs {
		allVirtualServices[index], err = c.GetVirtualService(ctx, avs[index].NsxtAlbVirtualService.GatewayRef.ID, avs[index].NsxtAlbVirtualService.ID)
		if err != nil {
			return nil, fmt.Errorf("error retrieving complete virtual service: %w", err)
		}
	}

	return allVirtualServices, nil
}

// GetVirtualService retrieves a virtual service by its name or ID from the specified edge gateway.
// It first validates the provided virtualServiceNameOrID and edgeGatewayID, ensuring they are not empty
// and have the correct format. If the virtualServiceNameOrID is not in the expected format, the edgeGatewayID
// is also validated. The function then refreshes the client session and attempts to retrieve the virtual service.
//
// Parameters:
//   - ctx: The context for the request.
//   - edgeGatewayID: The ID of the edge gateway containing the virtual service (required if virtualServiceNameOrID is a name).
//   - virtualServiceNameOrID: The name or ID of the virtual service to retrieve.
//
// Returns:
//   - *VirtualServiceModel: The retrieved virtual service model.
//   - error: An error if the retrieval fails or if any validation fails.
func (c *client) GetVirtualService(ctx context.Context, edgeGatewayID, virtualServiceNameOrID string) (*VirtualServiceModel, error) {
	if virtualServiceNameOrID == "" {
		return nil, fmt.Errorf("virtualServiceNameOrID is %w. Please provide a valid virtualServiceNameOrID", errors.ErrEmpty)
	}

	if !urn.IsLoadBalancerVirtualService(virtualServiceNameOrID) {
		if edgeGatewayID == "" {
			return nil, fmt.Errorf("edgeGatewayID is required if the provided virtual service is a name")
		}

		if !urn.IsEdgeGateway(edgeGatewayID) {
			return nil, fmt.Errorf("edgeGatewayID has %w. Please provide a valid edgeGatewayID", errors.ErrInvalidFormat)
		}
	}

	if err := c.clientCloudavenue.Refresh(); err != nil {
		return nil, err
	}

	vs, err := c.getVirtualService(ctx, edgeGatewayID, virtualServiceNameOrID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving virtual service: %w", err)
	}

	return fromVCDNsxtAlbVirtualServiceToModel(*vs.NsxtAlbVirtualService), nil
}

func (c *client) getVirtualService(_ context.Context, edgeGatewayID, virtualServiceNameOrID string) (*govcd.NsxtAlbVirtualService, error) {
	if !urn.IsLoadBalancerVirtualService(virtualServiceNameOrID) {
		return c.clientGoVCD.GetAlbVirtualServiceByName(edgeGatewayID, virtualServiceNameOrID)
	}
	return c.clientGoVCD.GetAlbVirtualServiceById(virtualServiceNameOrID)
}

// CreateVirtualService creates a new virtual service based on the provided VirtualServiceModelRequest.
func (c *client) CreateVirtualService(ctx context.Context, vsr VirtualServiceModelRequest) (*VirtualServiceModel, error) {
	if err := validators.New().Struct(vsr); err != nil {
		return nil, err
	}

	if err := c.clientCloudavenue.Refresh(); err != nil {
		return nil, err
	}

	model := fromModelRequestToVCDNsxtAlbVirtualService(vsr)

	if model.ServiceEngineGroupRef == (govcdtypes.OpenApiReference{}) {
		seg, err := c.GetFirstServiceEngineGroup(ctx, vsr.EdgeGatewayID)
		if err != nil {
			return nil, fmt.Errorf("error finding service engine group: %w", err)
		}

		model.ServiceEngineGroupRef = govcdtypes.OpenApiReference{
			ID:   seg.ID,
			Name: seg.Name,
		}
	}

	vs, err := c.clientGoVCD.CreateNsxtAlbVirtualService(model)
	if err != nil {
		return nil, fmt.Errorf("error creating virtual service: %w", err)
	}

	return fromVCDNsxtAlbVirtualServiceToModel(*vs.NsxtAlbVirtualService), nil
}

var updateVirtualService = func(virtualServiceClient fakeVirtualServiceClient, vs *govcdtypes.NsxtAlbVirtualService) (*govcd.NsxtAlbVirtualService, error) {
	return virtualServiceClient.Update(vs)
}

// UpdateVirtualService updates an existing virtual service identified by its ID.
func (c *client) UpdateVirtualService(ctx context.Context, virtualServiceID string, vsr VirtualServiceModelRequest) (*VirtualServiceModel, error) {
	if virtualServiceID == "" {
		return nil, fmt.Errorf("virtualServiceID is %w. Please provide a valid virtualServiceID", errors.ErrEmpty)
	}

	if !urn.IsLoadBalancerVirtualService(virtualServiceID) {
		return nil, fmt.Errorf("virtualServiceID has %w. Please provide a valid virtualServiceID", errors.ErrInvalidFormat)
	}

	if err := validators.New().Struct(vsr); err != nil {
		return nil, err
	}

	if err := c.clientCloudavenue.Refresh(); err != nil {
		return nil, err
	}

	vsToUpdate, err := c.getVirtualService(ctx, vsr.EdgeGatewayID, virtualServiceID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving virtual service: %w", err)
	}

	model := fromModelRequestToVCDNsxtAlbVirtualService(vsr)

	if model.ServiceEngineGroupRef == (govcdtypes.OpenApiReference{}) {
		seg, err := c.GetFirstServiceEngineGroup(ctx, vsr.EdgeGatewayID)
		if err != nil {
			return nil, fmt.Errorf("error finding service engine group: %w", err)
		}

		model.ServiceEngineGroupRef = govcdtypes.OpenApiReference{
			ID:   seg.ID,
			Name: seg.Name,
		}
	}

	model.ID = virtualServiceID

	vsUpdated, err := updateVirtualService(vsToUpdate, model)
	if err != nil {
		return nil, fmt.Errorf("error updating virtual service: %w", err)
	}

	return fromVCDNsxtAlbVirtualServiceToModel(*vsUpdated.NsxtAlbVirtualService), nil
}

var deleteVirtualService = func(virtualServiceClient fakeVirtualServiceClient) error {
	return virtualServiceClient.Delete()
}

// DeleteVirtualService deletes a virtual service identified by its ID.
func (c *client) DeleteVirtualService(ctx context.Context, virtualServiceID string) error {
	if virtualServiceID == "" {
		return fmt.Errorf("virtualServiceID is %w. Please provide a valid virtualServiceID", errors.ErrEmpty)
	}

	if !urn.IsLoadBalancerVirtualService(virtualServiceID) {
		return fmt.Errorf("virtualServiceID has %w. Please provide a valid virtualServiceID", errors.ErrInvalidFormat)
	}

	if err := c.clientCloudavenue.Refresh(); err != nil {
		return err
	}

	// edgegatewayID is not needed for retreiving the pool by ID
	vsToDelete, err := c.getVirtualService(ctx, "", virtualServiceID)
	if err != nil {
		return fmt.Errorf("error retrieving virtual service: %w", err)
	}

	return deleteVirtualService(vsToDelete)
}
