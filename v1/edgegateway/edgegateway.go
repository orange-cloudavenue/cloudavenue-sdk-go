/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package edgegateway

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/orange-cloudavenue/common-go/validators"

	"github.com/vmware/go-vcloud-director/v2/govcd"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/internal/endpoints"
	commoncloudavenue "github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/common/cloudavenue"
	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
)

// ListEdgeGateway fetches all edge gateways and returns them as a slice of EdgeGatewayModel.
func (c *client) ListEdgeGateway(ctx context.Context) ([]*EdgeGatewayModel, error) {
	if err := c.clientCloudavenue.Refresh(); err != nil {
		return nil, err
	}

	edgeGateways, err := c.clientGoVCDOrg.GetAllNsxtEdgeGateways(nil)
	if err != nil {
		return nil, err
	}

	edgeGatewayModels := make([]*EdgeGatewayModel, 0)

	for _, eg := range edgeGateways {
		egm := &EdgeGatewayModel{}
		egm.fromVCD(eg.EdgeGateway)

		edgeGatewayModels = append(edgeGatewayModels, egm)

		bandwidth, err := c.getBandwidth(ctx, egm)
		if err != nil {
			return nil, fmt.Errorf("error retrieving edge gateway %s bandwidth: %w", egm.ID, err)
		}

		egm.Bandwidth = bandwidth
	}

	return edgeGatewayModels, nil
}

// GetEdgeGateway retrieves an Edge Gateway by name or ID.
func (c *client) GetEdgeGateway(ctx context.Context, edgeGatewayNameOrID string) (*EdgeGateway, error) {
	if err := c.edgeGatewayNameOrIDValidation(edgeGatewayNameOrID); err != nil {
		return nil, err
	}

	if err := c.clientCloudavenue.Refresh(); err != nil {
		return nil, err
	}

	edgeGatewayModel, err := c.getEdgeGateway(ctx, edgeGatewayNameOrID)
	if err != nil {
		return nil, fmt.Errorf("error getting edge gateway: %w", err)
	}

	edge := &EdgeGateway{
		EdgeGatewayModel: edgeGatewayModel,
		internalClient:   c,
	}

	if err := edge.getNetworkServices(ctx); err != nil {
		return nil, fmt.Errorf("error getting edge gateway network services: %w", err)
	}

	return edge, nil
}

// DeleteEdgeGateway deletes an edge gateway by name or ID.
func (c *client) DeleteEdgeGateway(ctx context.Context, edgeGatewayNameOrID string) error {
	if err := c.edgeGatewayNameOrIDValidation(edgeGatewayNameOrID); err != nil {
		return err
	}

	if err := c.clientCloudavenue.Refresh(); err != nil {
		return err
	}
	edgeGatewayModel, err := c.getEdgeGateway(ctx, edgeGatewayNameOrID)
	if err != nil {
		return fmt.Errorf("error getting edge gateway: %w", err)
	}

	// Delete the edge gateway
	r, err := c.clientCloudavenue.R().
		SetContext(ctx).
		SetResult(&commoncloudavenue.JobStatus{}).
		SetError(&commoncloudavenue.APIErrorResponse{}).
		SetPathParam("edge-id", edgeGatewayModel.getUUID()).
		Delete(endpoints.EdgeGatewayDelete)
	if err != nil {
		return err
	}

	if r.IsError() {
		return fmt.Errorf("error on delete edge gateway: %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	job := r.Result().(*commoncloudavenue.JobStatus)

	// Wait for the job to finish
	return job.WaitWithContext(ctx, 2)
}

// CreateEdgeGateway creates a new edge gateway.
func (c *client) CreateEdgeGateway(ctx context.Context, edgeGateway *EdgeGatewayModelRequest) (*EdgeGatewayModel, error) {
	if err := validators.New().Struct(edgeGateway); err != nil {
		return nil, err
	}

	if err := c.clientCloudavenue.Refresh(); err != nil {
		return nil, err
	}

	if edgeGateway.OwnerRef == nil || (edgeGateway.OwnerRef.Name == "" && edgeGateway.OwnerRef.ID == "") {
		return nil, fmt.Errorf("owner reference name or ID is required")
	}

	var isVDCGroup bool

	// If OwnerRef Name is not set get the name from the ID
	if edgeGateway.OwnerRef.Name == "" {
		switch {
		case urn.IsVDC(edgeGateway.OwnerRef.ID):
			v, err := c.clientGoVCDOrg.GetVDCById(edgeGateway.OwnerRef.ID, true)
			if err != nil {
				return nil, err
			}

			edgeGateway.OwnerRef.Name = v.Vdc.Name
		case urn.IsVDCGroup(edgeGateway.OwnerRef.ID):
			isVDCGroup = true
			vg, err := c.clientGoVCDOrg.GetVdcGroupById(edgeGateway.OwnerRef.ID)
			if err != nil {
				return nil, err
			}

			edgeGateway.OwnerRef.Name = vg.VdcGroup.Name
		default:
			return nil, fmt.Errorf("invalid owner reference ID: %s", edgeGateway.OwnerRef.ID)
		}
	} else {
		vdcOrVDCG, err := (&v1.V1{}).VDC().GetVDCOrVDCGroup(edgeGateway.OwnerRef.Name)
		if err != nil {
			return nil, fmt.Errorf("error getting VDC or VDC Group: %w", err)
		}
		isVDCGroup = vdcOrVDCG.IsVDCGroup()
	}

	// Get the list of edge gateways before creating a new one. It's used to retrieve the ID of the new edge gateway.
	edgeGateways, err := c.ListEdgeGateway(ctx)
	if err != nil {
		return nil, err
	}

	var r *resty.Response

	// * Create the edge gateway in the VDC or VDC Group
	if isVDCGroup {
		// Create the edge gateway in the VDC Group
		r, err = c.clientCloudavenue.R().
			SetContext(ctx).
			SetResult(&commoncloudavenue.JobStatus{}).
			SetError(&commoncloudavenue.APIErrorResponse{}).
			SetBody(map[string]any{
				"tier0VrfId": edgeGateway.UplinkT0,
			}).
			SetPathParam("vdc-group-name", edgeGateway.OwnerRef.Name).
			Post(endpoints.EdgeGatewayCreateFromVDCGroup)
	} else {
		// Create the edge gateway in the VDC
		r, err = c.clientCloudavenue.R().
			SetContext(ctx).
			SetResult(&commoncloudavenue.JobStatus{}).
			SetError(&commoncloudavenue.APIErrorResponse{}).
			SetBody(map[string]any{
				"tier0VrfId": edgeGateway.UplinkT0,
			}).
			SetPathParam("vdc-name", edgeGateway.OwnerRef.Name).
			Post(endpoints.EdgeGatewayCreateFromVDC)
	}
	if err != nil {
		return nil, err
	}

	if r.IsError() {
		return nil, fmt.Errorf("error on create edge gateway: %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	job := r.Result().(*commoncloudavenue.JobStatus)

	// Wait for the job to finish
	if err := job.WaitWithContext(ctx, 2); err != nil {
		return nil, err
	}

	// ReGet the list of edge gateways to get the new one
	edgeGatewaysRefreshed, err := c.ListEdgeGateway(ctx)
	if err != nil {
		return nil, err
	}

	edgeGatewayCreated := new(EdgeGatewayModel)

	for _, egRefreshed := range edgeGatewaysRefreshed {
		found := false
		for _, egOld := range edgeGateways {
			if egRefreshed.ID == egOld.ID {
				found = true
				break
			}
		}
		if !found {
			edgeGatewayCreated = egRefreshed
			break
		}
	}

	// 5Mbps is a default value for the bandwidth
	if edgeGateway.Bandwidth != 5 {
		if err := c.updateBandwidth(ctx, edgeGatewayCreated, edgeGateway.Bandwidth); err != nil {
			return nil, fmt.Errorf("error on update edge gateway bandwidth: %w", err)
		}
	}

	return edgeGatewayCreated, nil
}

func (c *client) UpdateEdgeGateway(ctx context.Context, edgeGateway *EdgeGatewayModelUpdate) error {
	if err := validators.New().Struct(edgeGateway); err != nil {
		return err
	}

	if err := c.clientCloudavenue.Refresh(); err != nil {
		return err
	}

	return c.updateBandwidth(
		ctx,
		&EdgeGatewayModel{
			ID:        edgeGateway.ID,
			Bandwidth: edgeGateway.Bandwidth,
		},
		edgeGateway.Bandwidth,
	)
}

// * Local functions

// getEdgeGateway retrieves an edge gateway by name or ID.
func (c *client) getEdgeGateway(ctx context.Context, edgeGatewayNameOrID string) (*EdgeGatewayModel, error) {
	var (
		vcdEdgeGateway   *govcd.NsxtEdgeGateway
		err              error
		edgeGatewayModel = new(EdgeGatewayModel)
	)

	// If edgeGatewayNameOrID is a URN, get edge gateway by ID (more efficient)
	if urn.IsEdgeGateway(edgeGatewayNameOrID) { // Is URN
		vcdEdgeGateway, err = c.clientGoVCDOrg.GetNsxtEdgeGatewayById(edgeGatewayNameOrID)
	} else { // Is Name
		vcdEdgeGateway, err = c.clientGoVCDOrg.GetNsxtEdgeGatewayByName(edgeGatewayNameOrID)
	}

	if err != nil {
		return nil, fmt.Errorf("error retrieving edge gateway %s: %w", edgeGatewayNameOrID, err)
	}

	edgeGatewayModel.fromVCD(vcdEdgeGateway.EdgeGateway)

	bandwidth, err := c.getBandwidth(ctx, edgeGatewayModel)
	if err != nil {
		return nil, fmt.Errorf("error retrieving edge gateway %s bandwidth: %w", edgeGatewayNameOrID, err)
	}

	edgeGatewayModel.Bandwidth = bandwidth

	return edgeGatewayModel, nil
}

// getBandwidth retrieves the bandwidth of an edge gateway.
// It returns the bandwidth in Mbps.
func (c *client) getBandwidth(ctx context.Context, edgeGateway *EdgeGatewayModel) (int, error) {
	// Create the edge gateway
	r, err := c.clientCloudavenue.R().
		SetContext(ctx).
		SetResult(&bandwidthAPI{}).
		SetError(&commoncloudavenue.APIErrorResponse{}).
		SetPathParam("edge-id", edgeGateway.getUUID()).
		Get(endpoints.EdgeGatewayGet)
	if err != nil {
		return 0, err
	}

	if r.IsError() {
		return 0, fmt.Errorf("error on get edge gateway bandwidth: %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	bandwidthResponse := r.Result().(*bandwidthAPI)
	if bandwidthResponse == nil {
		return 0, fmt.Errorf("error on get edge gateway bandwidth: response is empty")
	}

	return bandwidthResponse.RateLimit, nil
}

// updateBandwidth updates the bandwidth of an edge gateway.
// It returns the bandwidth in Mbps.
func (c *client) updateBandwidth(ctx context.Context, edgeGateway *EdgeGatewayModel, bandwidth int) error {
	// TODO calcul remaining bandwidth
	if bandwidth == 0 {
		bandwidth = 5
	}

	// Create the edge gateway
	r, err := c.clientCloudavenue.R().
		SetContext(ctx).
		SetResult(&commoncloudavenue.JobStatus{}).
		SetError(&commoncloudavenue.APIErrorResponse{}).
		SetPathParam("edge-id", edgeGateway.getUUID()).
		SetBody(bandwidthAPI{
			RateLimit: bandwidth,
		}).
		Put(endpoints.EdgeGatewayUpdate)
	if err != nil {
		return err
	}

	if r.IsError() {
		return fmt.Errorf("error on update edge gateway bandwidth: %s", r.Error().(*commoncloudavenue.APIErrorResponse).FormatError())
	}

	job := r.Result().(*commoncloudavenue.JobStatus)

	// Wait for the job to finish
	return job.WaitWithContext(ctx, 2)
}
