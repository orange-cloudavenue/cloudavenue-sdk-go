/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 Orange
 * SPDX-License-Identifier: Mozilla Public License 2.0
 *
 * This software is distributed under the MPL-2.0 license.
 * the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
 * or see the "LICENSE" file for more details.
 */

package v1

import (
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"
)

type (
	// EdgeGatewayALBServiceEngineGroupModel represents an ALB Service Engine Group to an Edge Gateway.
	EdgeGatewayALBServiceEngineGroupModel struct {
		ID string `json:"id,omitempty"` // urn format of the service engine group

		// Name of the service engine group
		Name string `json:"name,omitempty"`

		// GatewayRef contains reference to Edge Gateway
		GatewayRef *govcdtypes.OpenApiReference `json:"gatewayRef"`

		// MaxVirtualServices is the maximum number of virtual services that can be deployed
		MaxVirtualServices *int `json:"maxVirtualServices,omitempty"`

		// MinVirtualServices is the minimum number (reserved) of virtual services that can be deployed
		MinVirtualServices *int `json:"minVirtualServices,omitempty"`

		// NumDeployedVirtualServices is a read only value
		NumDeployedVirtualServices int `json:"numDeployedVirtualServices,omitempty"`
	}
)
