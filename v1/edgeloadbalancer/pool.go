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

	"github.com/orange-cloudavenue/common-go/validators"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/errors"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

func (c *client) ListPools(ctx context.Context, edgeGatewayID string) ([]*PoolModel, error) {
	if edgeGatewayID == "" {
		return nil, fmt.Errorf("edgeGatewayID is %w. Please provide a valid edgeGatewayID", errors.ErrEmpty)
	}

	if !urn.IsEdgeGateway(edgeGatewayID) {
		return nil, fmt.Errorf("edgeGatewayID has %w. Please provide a valid edgeGatewayID", errors.ErrInvalidFormat)
	}

	if err := c.clientCloudavenue.Refresh(); err != nil {
		return nil, err
	}

	allAlbPoolSummaries, err := c.clientGoVCD.GetAllAlbPoolSummaries(edgeGatewayID, url.Values{})
	if err != nil {
		return nil, fmt.Errorf("error retrieving all ALB Pool summaries: %w", err)
	}

	// Loop over all Summaries and retrieve complete information
	allAlbPools := make([]*PoolModel, len(allAlbPoolSummaries))
	for index := range allAlbPoolSummaries {
		allAlbPools[index], err = c.GetPool(ctx, allAlbPoolSummaries[index].NsxtAlbPool.GatewayRef.ID, allAlbPoolSummaries[index].NsxtAlbPool.ID)
		if err != nil {
			return nil, fmt.Errorf("error retrieving complete ALB Pool: %w", err)
		}
	}

	return allAlbPools, nil
}

// GetPool retrieves a pool by name or ID.
func (c *client) GetPool(ctx context.Context, edgeGatewayID, poolNameOrID string) (*PoolModel, error) {
	if edgeGatewayID == "" {
		return nil, fmt.Errorf("edgeGatewayID is %w. Please provide a valid edgeGatewayID", errors.ErrEmpty)
	}

	if !urn.IsEdgeGateway(edgeGatewayID) {
		return nil, fmt.Errorf("edgeGatewayID has %w. Please provide a valid edgeGatewayID", errors.ErrInvalidFormat)
	}

	if poolNameOrID == "" {
		return nil, fmt.Errorf("poolNameOrID is %w. Please provide a valid poolNameOrID", errors.ErrEmpty)
	}

	if err := c.clientCloudavenue.Refresh(); err != nil {
		return nil, err
	}

	albPool, err := c.getpool(ctx, edgeGatewayID, poolNameOrID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Load Balancer Pool: %w", err)
	}

	return fromVCDNsxtALBPoolToModel(albPool.NsxtAlbPool), nil
}

func (c *client) getpool(_ context.Context, edgeGatewayID, nameOrID string) (*govcd.NsxtAlbPool, error) {
	var (
		albPool *govcd.NsxtAlbPool
		err     error
	)

	if urn.IsLoadBalancerPool(nameOrID) {
		albPool, err = c.clientGoVCD.GetAlbPoolById(nameOrID)
	} else {
		albPool, err = c.clientGoVCD.GetAlbPoolByName(edgeGatewayID, nameOrID)
	}

	return albPool, err
}

func (c *client) CreatePool(ctx context.Context, pool PoolModelRequest) (*PoolModel, error) {
	if err := validators.New().StructCtx(ctx, &pool); err != nil {
		return nil, err
	}

	if err := c.clientCloudavenue.Refresh(); err != nil {
		return nil, err
	}

	poolCreated, err := c.clientGoVCD.CreateNsxtAlbPool(fromModelToGoVCDNsxtALBPool("", pool))
	if err != nil {
		return nil, fmt.Errorf("error creating Load Balancer Pool: %w", err)
	}

	return fromVCDNsxtALBPoolToModel(poolCreated.NsxtAlbPool), nil
}

var updatePool = func(poolClient fakePoolClient, pool *govcdtypes.NsxtAlbPool) (*govcd.NsxtAlbPool, error) {
	return poolClient.Update(pool)
}

func (c *client) UpdatePool(ctx context.Context, poolID string, pool PoolModelRequest) (*PoolModel, error) {
	if poolID == "" {
		return nil, fmt.Errorf("poolID is %w. Please provide a valid poolID", errors.ErrEmpty)
	}

	if !urn.IsLoadBalancerPool(poolID) {
		return nil, fmt.Errorf("poolID has %w. Please provide a valid poolID", errors.ErrInvalidFormat)
	}

	if err := validators.New().StructCtx(ctx, &pool); err != nil {
		return nil, err
	}

	if err := c.clientCloudavenue.Refresh(); err != nil {
		return nil, err
	}

	// edgegatewayID is not needed for retreiving the pool by ID
	poolToUpdate, err := c.getpool(ctx, "", poolID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Load Balancer Pool: %w", err)
	}

	poolUpdated, err := updatePool(poolToUpdate, fromModelToGoVCDNsxtALBPool(poolToUpdate.NsxtAlbPool.ID, pool))
	if err != nil {
		return nil, fmt.Errorf("error updating Load Balancer Pool: %w", err)
	}

	return fromVCDNsxtALBPoolToModel(poolUpdated.NsxtAlbPool), nil
}

var deletePool = func(poolClient fakePoolClient) error {
	return poolClient.Delete()
}

func (c *client) DeletePool(ctx context.Context, poolID string) error {
	if poolID == "" {
		return fmt.Errorf("poolID is %w. Please provide a valid poolID", errors.ErrEmpty)
	}

	if !urn.IsLoadBalancerPool(poolID) {
		return fmt.Errorf("poolID has %w. Please provide a valid poolID", errors.ErrInvalidFormat)
	}

	if err := c.clientCloudavenue.Refresh(); err != nil {
		return err
	}

	// edgegatewayID is not needed for retreiving the pool by ID
	poolToDelete, err := c.getpool(ctx, "", poolID)
	if err != nil {
		return fmt.Errorf("error retrieving Load Balancer Pool: %w", err)
	}

	return deletePool(poolToDelete)
}
