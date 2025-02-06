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
	govcdtypes "github.com/vmware/go-vcloud-director/v2/types/v56"
)

type (
	// ServiceEngineGroupModel represents an ALB Service Engine Group to an Edge Gateway.
	ServiceEngineGroupModel struct {
		ID string // urn format of the service engine group

		// Name of the service engine group
		Name string

		// GatewayRef contains reference to Edge Gateway
		GatewayRef *govcdtypes.OpenApiReference

		// MaxVirtualServices is the maximum number of virtual services that can be deployed
		MaxVirtualServices *int

		// MinVirtualServices is the minimum number (reserved) of virtual services that can be deployed
		MinVirtualServices *int

		// NumDeployedVirtualServices is a number of deployed virtual services
		NumDeployedVirtualServices int
	}
)
