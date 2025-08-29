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
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"

	"github.com/orange-cloudavenue/cloudavenue-sdk-go/pkg/urn"
)

type (
	EdgeGateway struct {
		*EdgeGatewayModel

		// internal
		internalClient clientInterface
	}

	// EdgeGatewayModel represents the model of an edge gateway.
	EdgeGatewayModel struct { //nolint:revive
		ID string

		// Name of edge gateway
		Name string

		// Description of edge gateway
		Description string

		// OwnerRef defines Org VDC or VDC Group that this network belongs to.
		OwnerRef *govcdtypes.OpenApiReference

		// UplinkT0 defines the T0 router name that this edge gateway is connected to.
		UplinkT0 string

		// Bandwidth defines the bandwidth of the edge gateway.
		Bandwidth int

		Status string

		// Services is the list of network services
		// that are available on the edge gateway
		Services NetworkServicesModelSvcs
	}

	// EdgeGatewayModelRequest represents the request model for creating an edge gateway.
	EdgeGatewayModelRequest struct { //nolint:revive
		// OwnerRef defines Org VDC or VDC Group that this network belongs to.
		OwnerRef *govcdtypes.OpenApiReference `validate:"required"`

		// UplinkT0 defines the T0 router name that this edge gateway is connected to.
		UplinkT0 string `validate:"omitempty"`

		// Bandwidth defines the bandwidth of the edge gateway.
		Bandwidth int `validate:"required,min=5"`
	}

	// EdgeGatewayModelUpdate represents the update model for an edge gateway.
	EdgeGatewayModelUpdate struct { //nolint:revive
		// ID of the edge gateway.
		ID string `validate:"required"`

		// Bandwidth defines the bandwidth of the edge gateway.
		Bandwidth int `validate:"required,min=5"`
	}

	// -----.

	// Bandwidth represents the bandwidth of the edge gateway. (InfrAPI).
	bandwidthAPI struct {
		RateLimit int `json:"rateLimit"`
	}
)

// getUUID returns the UUID of the edge gateway.
// It is a specialized function to extract the UUID from the ID(URN) to call cloudavenue API.
func (m *EdgeGatewayModel) getUUID() string {
	return urn.ExtractUUID(m.ID)
}

// fromVCD converts a VCD edge gateway model to the internal EdgeGatewayModel.
func (m *EdgeGatewayModel) fromVCD(vcdEdgeGateway *govcdtypes.OpenAPIEdgeGateway) {
	if vcdEdgeGateway == nil {
		return
	}

	m.ID = vcdEdgeGateway.ID
	m.Name = vcdEdgeGateway.Name
	m.Description = vcdEdgeGateway.Description
	m.OwnerRef = vcdEdgeGateway.OwnerRef
	m.Status = vcdEdgeGateway.Status
	if len(vcdEdgeGateway.EdgeGatewayUplinks) > 0 {
		m.UplinkT0 = vcdEdgeGateway.EdgeGatewayUplinks[0].UplinkName
	}
}
